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
						Body: matchers.EachLike(map[string]interface{}{
							"id":            matchers.Like(105),
							"notes":         matchers.String("Test"),
							"objectionType": matchers.String("factual"),
							"receivedDate":  matchers.String("05/09/2024"),
							"lpaUids":       []string{"M-9999-9999-9999"},
						}, 1),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: []Objection{
				{
					ID:            105,
					Notes:         "Test",
					ObjectionType: "factual",
					ReceivedDate:  "05/09/2024",
					LpaUids:       []string{"M-9999-9999-9999"},
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

func TestAddObjection(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name           string
		objectionsData ObjectionRequest
		setup          func()
		expectedError  func(int) error
	}{
		{
			name: "OK",
			objectionsData: ObjectionRequest{
				LpaUids:       []string{"M-1234-9876-4567"},
				ReceivedDate:  "2025-01-02",
				ObjectionType: "factual",
				Notes:         "test",
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("A digital LPA exists").
					UponReceiving("A request to add an objection").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPost,
						Path:   matchers.String("/lpa-api/v1/objections"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: matchers.Like(map[string]interface{}{
							"lpaUids":       []string{"M-1234-9876-4567"},
							"receivedDate":  matchers.Like("02/01/2025"),
							"objectionType": matchers.Like("factual"),
							"notes":         matchers.Like("test"),
						}),
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusCreated,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				err := client.AddObjection(Context{Context: context.Background()}, tc.objectionsData)
				if (tc.expectedError) == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(config.Port), err)
				}
				return nil
			}))
		})
	}
}

func TestUpdateObjection(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name           string
		objectionsData ObjectionRequest
		setup          func()
		expectedError  func(int) error
	}{
		{
			name: "OK",
			objectionsData: ObjectionRequest{
				LpaUids:       []string{"M-9999-9999-9999"},
				ReceivedDate:  "2025-01-02",
				ObjectionType: "prescribed",
				Notes:         "test",
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a digital LPA with an objection").
					UponReceiving("A request to update an objection").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPut,
						Path:   matchers.String("/lpa-api/v1/objections/3"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: matchers.Like(map[string]interface{}{
							"lpaUids":       matchers.EachLike(matchers.String("M-9999-9999-9999"), 1),
							"receivedDate":  matchers.Like("02/01/2025"),
							"objectionType": matchers.Like("prescribed"),
							"notes":         matchers.Like("test"),
						}),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusNoContent,
					})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				err := client.UpdateObjection(Context{Context: context.Background()}, "3", tc.objectionsData)
				if (tc.expectedError) == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(config.Port), err)
				}
				return nil
			}))
		})
	}
}

func TestGetObjection(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		setup            func()
		expectedResponse Objection2
		expectedError    func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a digital LPA with an objection").
					UponReceiving("A request for the objection").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/objections/3"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.Like(map[string]interface{}{
							"id":            matchers.Like(3),
							"notes":         matchers.String("Test"),
							"objectionType": matchers.String("factual"),
							"receivedDate":  matchers.String("05/09/2024"),
							"lpaUids":       []string{"M-1234-9876-4567"},
							"resolutions": []map[string]interface{}{
								{
									"resolution":      matchers.Like("not upheld"),
									"resolutionNotes": matchers.Like("Everything is fine"),
									"resolutionDate":  matchers.Like("01/01/2025"),
								},
							},
						}),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: Objection2{
				ID:            3,
				Notes:         "Test",
				ObjectionType: "factual",
				ReceivedDate:  "05/09/2024",
				LpaUids:       []string{"M-1234-9876-4567"},
				Resolutions: []ObjectionResolution{
					{
						Resolution:      "not upheld",
						ResolutionNotes: "Everything is fine",
						ResolutionDate:  "2025-01-01",
					},
				},
			},
		},
		{
			name: "404",
			setup: func() {
				pact.
					AddInteraction().
					Given("There is no objection with the specified ID").
					UponReceiving("A request for a non-existent objection").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/objections/3"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusNotFound,
					})
			},
			expectedError: func(port int) error {
				return StatusError{
					Code:          404,
					URL:           fmt.Sprintf("http://127.0.0.1:%d/lpa-api/v1/objections/3", port),
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

				objection, err := client.GetObjection(Context{Context: context.Background()}, "3")

				assert.Equal(t, tc.expectedResponse, objection)
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
