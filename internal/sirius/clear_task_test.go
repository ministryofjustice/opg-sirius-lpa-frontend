package sirius

import (
	"context"
	"errors"
	"fmt"
	"github.com/pact-foundation/pact-go/dsl"
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

func TestClearTaskForCaseNetworkError(t *testing.T) {
	mockClearTaskHttpClient := &mockClearTaskHttpClient{}

	mockClearTaskHttpClient.On("Do", mock.Anything).Return(&http.Response{}, errors.New("Networking issue"))

	client := NewClient(mockClearTaskHttpClient, "http://localhost")
	err := client.ClearTask(Context{Context: context.Background()}, 990)
	assert.Equal(t, "Networking issue", err.Error())
}

func TestClearTask(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

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
					WithRequest(dsl.Request{
						Method: http.MethodPut,
						Path:   dsl.String("/lpa-api/v1/tasks/990/mark-as-completed"),
					}).
					WillRespondWith(dsl.Response{
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
					WithRequest(dsl.Request{
						Method: http.MethodPut,
						Path:   dsl.Like("/lpa-api/v1/tasks/990/mark-as-completed"),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusNotFound,
					})
			},
			expectedError: func(port int) error {
				return StatusError{
					Code:          404,
					URL:           fmt.Sprintf("http://localhost:%d/lpa-api/v1/tasks/990/mark-as-completed", port),
					Method:        "PUT",
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

				err := client.ClearTask(Context{Context: context.Background()}, 990)

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
