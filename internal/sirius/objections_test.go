package sirius

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockObjectionHttpClient struct {
	mock.Mock
}

func (m *mockObjectionHttpClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestObjectionsForCaseBadContext(t *testing.T) {
	client := NewClient(http.DefaultClient, "http://localhost")
	_, err := client.ObjectionsForCase(Context{Context: nil}, "M-9999-9999-9999")
	assert.Equal(t, "net/http: nil Context", err.Error())
}

func TestObjectionsForCaseNetworkError(t *testing.T) {
	mockObjectionHttpClient := &mockObjectionHttpClient{}

	mockObjectionHttpClient.On("Do", mock.Anything).Return(&http.Response{}, errors.New("Networking issue"))

	client := NewClient(mockObjectionHttpClient, "http://localhost")
	_, err := client.ObjectionsForCase(Context{Context: context.Background()}, "M-9999-9999-9999")
	assert.Equal(t, "Networking issue", err.Error())
}

func TestObjectionsForCaseBadResponseJSON(t *testing.T) {
	mockObjectionHttpClient := &mockObjectionHttpClient{}

	badJsonResponse := http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString("[bad json time")),
	}

	mockObjectionHttpClient.On("Do", mock.Anything).Return(&badJsonResponse, nil)

	client := NewClient(mockObjectionHttpClient, "http://localhost")
	_, err := client.ObjectionsForCase(Context{Context: context.Background()}, "M-9999-9999-9999")

	assert.Equal(t, "invalid character 'b' looking for beginning of value", err.Error())
}

func TestObjectionsForCase(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		setup            func()
		expectedResponse []Objection
		expectedError    func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a digital LPA with an objection").
					UponReceiving("A request for the objections for a case").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/digital-lpas/M-9999-9999-9999/objections"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.Like(map[string]interface{}{
							"uid": matchers.Regex("M-9999-9999-9999", `^M(-[0-9A-Z]{4}){3}$`),
							"objections": matchers.Like([]map[string]interface{}{
								{
									"id":            matchers.Like(105),
									"notes":         matchers.String("Test"),
									"objectionType": matchers.String("factual"),
									"receivedDate":  matchers.String("05/09/2024"),
								},
							}),
						}),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: []Objection{
				Objection{
					ID:            105,
					Notes:         "Test",
					ObjectionType: "factual",
					ReceivedDate:  "05/09/2024",
				},
			},
		},
		{
			name: "404",
			setup: func() {
				pact.
					AddInteraction().
					Given("There is no case with objections with the specified caseUID").
					UponReceiving("A request for the objections for a non-existent case").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/digital-lpas/M-9999-9999-9999/objections"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusNotFound,
					})
			},
			expectedError: func(port int) error {
				return StatusError{
					Code:          404,
					URL:           fmt.Sprintf("http://127.0.0.1:%d/lpa-api/v1/digital-lpas/M-9999-9999-9999/objections", port),
					Method:        "GET",
					CorrelationId: "",
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				objectionList, err := client.ObjectionsForCase(Context{Context: context.Background()}, "M-9999-9999-9999")

				assert.Equal(t, tc.expectedResponse, objectionList)
				if tc.expectedError == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(config.Port), err)
				}

				return nil
			}))
		})
	}
}
