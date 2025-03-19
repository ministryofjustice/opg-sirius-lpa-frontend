package sirius

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newPact() (*consumer.V4HTTPMockProvider, error) {
	return consumer.NewV4Pact(consumer.MockHTTPProviderConfig{
		Consumer: "sirius-lpa-frontend",
		Provider: "sirius",
		Host:     "127.0.0.1",
		LogDir:   "../../logs",
		PactDir:  "../../pacts",
	})
}

func TestStatusError(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "/some/url", nil)

	resp := &http.Response{
		StatusCode: http.StatusTeapot,
		Request:    req,
	}

	err := newStatusError(resp)

	assert.Equal(t, "POST /some/url returned 418", err.Error())
	assert.Equal(t, "unexpected response from Sirius", err.Title())
	assert.Equal(t, err, err.Data())
	assert.False(t, err.IsUnauthorized())
}

func TestStatusErrorIsUnauthorized(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "/some/url", nil)

	resp := &http.Response{
		StatusCode: http.StatusUnauthorized,
		Request:    req,
	}

	err := newStatusError(resp)

	assert.Equal(t, "POST /some/url returned 401", err.Error())
	assert.Equal(t, "unexpected response from Sirius", err.Title())
	assert.Equal(t, err, err.Data())
	assert.True(t, err.IsUnauthorized())
}

func TestToFieldErrors(t *testing.T) {
	var unformattedErr flexibleFieldErrors
	err := json.Unmarshal([]byte(`{"riskAssessmentDate":["This field is required"],"reportApprovalDate":["This field is required"]}`), &unformattedErr)
	if err != nil {
		return
	}
	result, err := unformattedErr.toFieldErrors()
	formattedErr := FieldErrors{"riskAssessmentDate": {"": "This field is required"}, "reportApprovalDate": {"": "This field is required"}}

	assert.Equal(t, formattedErr, result)
	assert.Nil(t, err)
}

func TestToFieldErrorsThrowsError(t *testing.T) {
	var unformattedErr flexibleFieldErrors
	err := json.Unmarshal([]byte(`{"test":123}`), &unformattedErr)
	if err != nil {
		return
	}
	result, err := unformattedErr.toFieldErrors()

	assert.Equal(t, err, errors.New("could not parse field validation_errors"))
	assert.Nil(t, result)
}

func TestValidationErrorSummary(t *testing.T) {
	emptyValidationError := ValidationError{}

	assert.Equal(t, "validation error", emptyValidationError.Error())

	descriptiveValidationError := ValidationError{
		Detail: "This case is in a Registered status and cannot be deleted",
	}

	assert.Equal(t, "This case is in a Registered status and cannot be deleted", descriptiveValidationError.Error())
}

type mockHTTPClient struct {
	mock.Mock
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestRequestInternalError(t *testing.T) {
	testCases := map[string]struct {
		fn func(client *Client) error
	}{
		"get": {
			fn: func(client *Client) error {
				return client.get(Context{Context: context.Background()}, "/resource/14", nil)
			},
		},
		"post": {
			fn: func(client *Client) error {
				return client.post(Context{Context: context.Background()}, "/resource/14", nil, nil)
			},
		},
		"put": {
			fn: func(client *Client) error {
				return client.put(Context{Context: context.Background()}, "/resource/14", nil, nil)
			},
		},
		"delete": {
			fn: func(client *Client) error {
				return client.delete(Context{Context: context.Background()}, "/resource/14")
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			expectedErr := errors.New("something went wrong")

			mockHttpClient := &mockHTTPClient{}
			mockHttpClient.
				On("Do", mock.Anything).
				Return(&http.Response{}, expectedErr)

			client := NewClient(mockHttpClient, "https://host.example")

			err := tc.fn(client)

			assert.Equal(t, expectedErr, err)
		})
	}
}

func TestRequestStatusError(t *testing.T) {
	testCases := map[string]struct {
		fn func(client *Client) error
	}{
		"get": {
			fn: func(client *Client) error {
				return client.get(Context{Context: context.Background()}, "/resource/14", nil)
			},
		},
		"post": {
			fn: func(client *Client) error {
				return client.post(Context{Context: context.Background()}, "/resource/14", nil, nil)
			},
		},
		"put": {
			fn: func(client *Client) error {
				return client.put(Context{Context: context.Background()}, "/resource/14", nil, nil)
			},
		},
		"delete": {
			fn: func(client *Client) error {
				return client.delete(Context{Context: context.Background()}, "/resource/14")
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mockHttpClient := &mockHTTPClient{}
			mockCall := mockHttpClient.On("Do", mock.Anything)

			mockCall.RunFn = func(args mock.Arguments) {
				mockCall.ReturnArguments = mock.Arguments{&http.Response{
					Request:    args.Get(0).(*http.Request),
					StatusCode: http.StatusServiceUnavailable,
					Body:       io.NopCloser(strings.NewReader(`{"id":492, "name":"policy_4"}`)),
				}, nil}
			}

			client := NewClient(mockHttpClient, "https://host.example")

			err := tc.fn(client)

			statusError, ok := err.(StatusError)
			if !ok {
				t.Error("expected error to be StatusError")
			}

			assert.Equal(t, http.StatusServiceUnavailable, statusError.Code)
		})
	}
}

func TestRequestMarshalBodyError(t *testing.T) {
	testCases := map[string]struct {
		fn func(client *Client) error
	}{
		"post": {
			fn: func(client *Client) error {
				return client.post(Context{Context: context.Background()}, "/resource/14", complex(1, 16), nil)
			},
		},
		"put": {
			fn: func(client *Client) error {
				return client.put(Context{Context: context.Background()}, "/resource/14", complex(1, 16), nil)
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mockHttpClient := &mockHTTPClient{}

			client := NewClient(mockHttpClient, "https://host.example")

			err := tc.fn(client)

			assert.Equal(t, "json: unsupported type: complex128", err.Error())
		})
	}
}

func TestRequestUnmarshalError(t *testing.T) {
	testCases := map[string]struct {
		status int
		fn     func(client *Client) error
	}{
		"get": {
			fn: func(client *Client) error {
				resp := ""
				return client.get(Context{Context: context.Background()}, "/resource/14", &resp)
			},
			status: http.StatusOK,
		},
		"post": {
			fn: func(client *Client) error {
				resp := ""
				return client.post(Context{Context: context.Background()}, "/resource/14", nil, &resp)
			},
			status: http.StatusCreated,
		},
		"put": {
			fn: func(client *Client) error {
				resp := ""
				return client.put(Context{Context: context.Background()}, "/resource/14", nil, &resp)
			},
			status: http.StatusNoContent,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mockHttpClient := &mockHTTPClient{}
			mockCall := mockHttpClient.On("Do", mock.Anything)

			mockCall.RunFn = func(args mock.Arguments) {
				mockCall.ReturnArguments = mock.Arguments{&http.Response{
					Request:    args.Get(0).(*http.Request),
					StatusCode: tc.status,
					Body:       io.NopCloser(strings.NewReader(`{"id":492, "name":"policy_4"}`)),
				}, nil}
			}

			client := NewClient(mockHttpClient, "https://host.example")

			err := tc.fn(client)

			assert.Equal(t, "json: cannot unmarshal object into Go value of type string", err.Error())
		})
	}
}
