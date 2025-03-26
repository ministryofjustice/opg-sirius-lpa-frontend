package sirius

import (
	"context"
	"errors"
	"fmt"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

type mockClearTaskHttpClient struct {
	mock.Mock
}

func (m *mockClearTaskHttpClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestClearTaskForCaseBadContext(t *testing.T) {
	client := NewClient(http.DefaultClient, "http://localhost")
	err := client.ClearTask(Context{Context: nil}, 990)
	assert.Equal(t, "net/http: nil Context", err.Error())
}

func TestClearTaskForCaseNetworkError(t *testing.T) {
	mockClearTaskHttpClient := &mockClearTaskHttpClient{}

	mockClearTaskHttpClient.On("Do", mock.Anything).Return(&http.Response{}, errors.New("Networking issue"))

	client := NewClient(mockClearTaskHttpClient, "http://localhost")
	err := client.ClearTask(Context{Context: context.Background()}, 990)
	assert.Equal(t, "Networking issue", err.Error())
}

func TestClearTask(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name          string
		setup         func()
		expectedError func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a case with an open task assigned").
					UponReceiving("A request to clear the task").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPut,
						Path:   matchers.String("/lpa-api/v1/tasks/990/mark-as-completed"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
					})
			},
		},
		{
			name: "404",
			setup: func() {
				pact.
					AddInteraction().
					Given("There is no tasks with the specified ID").
					UponReceiving("A request to clear a non-existent task").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPut,
						Path:   matchers.Like("/lpa-api/v1/tasks/990/mark-as-completed"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusNotFound,
					})
			},
			expectedError: func(port int) error {
				return StatusError{
					Code:          404,
					URL:           fmt.Sprintf("http://127.0.0.1:%d/lpa-api/v1/tasks/990/mark-as-completed", port),
					Method:        "PUT",
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

				err := client.ClearTask(Context{Context: context.Background()}, 990)

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
