package sirius

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
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
	_, err := client.TasksForCase(Context{Context: context.Background()}, 990)
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
	_, err := client.TasksForCase(Context{Context: context.Background()}, 990)

	assert.Equal(t, "invalid character 'b' looking for beginning of value", err.Error())
}

func TestWarningsForCase(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

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
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/lpa-api/v1/cases/990/warnings"),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: dsl.EachLike(map[string]interface{}{
							"id":           dsl.Like(9901),
							"dateAdded":    dsl.String("07/01/2023"),
							"warningType":  dsl.String("Donor Deceased"),
							"warningText":  dsl.String("Donor died"),
							"caseItems":    dsl.EachLike(map[string]interface{}{
								"uId":      dsl.String("M-TTTT-RRRR-EEEE"),
								"caseType": dsl.String("DIGITAL_LPA"),
							}, 1),
						}, 1),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
			},
			expectedResponse: []Warning{
				Warning{
					ID: 9901,
					DateAdded: "07/01/2023",
					WarningType: "Donor Deceased",
					WarningText: "Donor died",
					CaseItems: []Case{
						Case{
							UID: "M-TTTT-RRRR-EEEE",
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
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/lpa-api/v1/cases/990/warnings"),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusNotFound,
					})
			},
			expectedError: func(port int) error {
				return StatusError{
					Code:          404,
					URL:           fmt.Sprintf("http://localhost:%d/lpa-api/v1/cases/990/warnings", port),
					Method:        "GET",
					CorrelationId: "",
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				warnings, err := client.WarningsForCase(Context{Context: context.Background()}, 990)

				assert.Equal(t, tc.expectedResponse, warnings)
				if tc.expectedError == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(pact.Server.Port), err)
				}

				return nil
			}))
		})
	}
}
