package sirius

import (
	"context"
	"fmt"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name          string
		setup         func()
		expectedError func(int) error
		file          *NoteFile
		task          TaskRequest
	}{
		{
			name: "OK for user",
			task: TaskRequest{
				AssigneeID:  47,
				Type:        "Check Application",
				Name:        "Something",
				Description: "More words",
				DueDate:     "2035-03-04",
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending case assigned").
					UponReceiving("A request to create a task").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPost,
						Path:   matchers.String("/lpa-api/v1/cases/800/tasks"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]interface{}{
							"assigneeId":  47,
							"type":        "Check Application",
							"name":        "Something",
							"description": "More words",
							"dueDate":     matchers.Term("04/03/2035", `^\d{1,2}/\d{1,2}/\d{4}$`),
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusCreated,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
		},
		{
			name: "OK for team",
			task: TaskRequest{
				AssigneeID:  23,
				Type:        "Check Application",
				Name:        "A title",
				Description: "A description",
				DueDate:     "2035-03-04",
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("An LPA and a team exists").
					UponReceiving("A request to create a task for a team").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPost,
						Path:   matchers.String("/lpa-api/v1/cases/800/tasks"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]interface{}{
							"assigneeId":  23,
							"type":        "Check Application",
							"name":        "A title",
							"description": "A description",
							"dueDate":     matchers.Term("04/03/2035", `^\d{1,2}/\d{1,2}/\d{4}$`),
						},
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

				err := client.CreateTask(Context{Context: context.Background()}, 800, tc.task)
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
