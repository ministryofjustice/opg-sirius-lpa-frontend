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

type mockWarningHttpClient struct {
	mock.Mock
}

func (m *mockWarningHttpClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestWarningsForCaseBadContext(t *testing.T) {
	client := NewClient(http.DefaultClient, "http://localhost")
	_, err := client.WarningsForCase(Context{Context: nil}, 990)
	assert.Equal(t, "net/http: nil Context", err.Error())
}

func TestWarningsForCaseNetworkError(t *testing.T) {
	mockWarningHttpClient := &mockWarningHttpClient{}

	mockWarningHttpClient.On("Do", mock.Anything).Return(&http.Response{}, errors.New("Networking issue"))

	client := NewClient(mockWarningHttpClient, "http://localhost")
	_, err := client.WarningsForCase(Context{Context: context.Background()}, 990)
	assert.Equal(t, "Networking issue", err.Error())
}

func TestWarningsForCaseBadResponseJSON(t *testing.T) {
	mockWarningHttpClient := &mockWarningHttpClient{}

	badJsonResponse := http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString("[bad json time")),
	}

	mockWarningHttpClient.On("Do", mock.Anything).Return(&badJsonResponse, nil)

	client := NewClient(mockWarningHttpClient, "http://localhost")
	_, err := client.WarningsForCase(Context{Context: context.Background()}, 990)

	assert.Equal(t, "invalid character 'b' looking for beginning of value", err.Error())
}

func TestWarningsForCase(t *testing.T) {
	t.Parallel()

	pact, err := newPact2()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		setup            func()
		expectedResponse []Warning
		expectedError    func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a case with a warning").
					UponReceiving("A request for the warnings for a case").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/cases/990/warnings"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.EachLike(map[string]interface{}{
							"id":          matchers.Like(9901),
							"dateAdded":   matchers.Regex("07/01/2023 10:10:10", `\d{2}\/\d{2}\/\d{4} \d{2}:\d{2}:\d{2}`),
							"warningType": matchers.String("Attorney removed"),
							"warningText": matchers.String("Attorney was removed"),
							"caseItems": matchers.EachLike(map[string]interface{}{
								"uId":      matchers.String("M-TTTT-RRRR-EEEE"),
								"caseType": matchers.String("DIGITAL_LPA"),
							}, 1),
						}, 1),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: []Warning{
				Warning{
					ID:          9901,
					DateAdded:   "07/01/2023 10:10:10",
					WarningType: "Attorney removed",
					WarningText: "Attorney was removed",
					CaseItems: []Case{
						Case{
							UID:      "M-TTTT-RRRR-EEEE",
							CaseType: "DIGITAL_LPA",
						},
					},
				},
			},
		},
		{
			name: "404",
			setup: func() {
				pact.
					AddInteraction().
					Given("There is no case with warnings with the specified ID").
					UponReceiving("A request for the warnings for a non-existent case").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/cases/990/warnings"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusNotFound,
					})
			},
			expectedError: func(port int) error {
				return StatusError{
					Code:          404,
					URL:           fmt.Sprintf("http://127.0.0.1:%d/lpa-api/v1/cases/990/warnings", port),
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

				warnings, err := client.WarningsForCase(Context{Context: context.Background()}, 990)

				assert.Equal(t, tc.expectedResponse, warnings)
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
