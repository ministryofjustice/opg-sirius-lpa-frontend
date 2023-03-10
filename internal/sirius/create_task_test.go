package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

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
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/lpa-api/v1/cases/800/tasks"),
						Headers: dsl.MapMatcher{
							"Content-Type": dsl.String("application/json"),
						},
						Body: map[string]interface{}{
							"assigneeId":  47,
							"type":        "Check Application",
							"name":        "Something",
							"description": "More words",
							"dueDate":     dsl.Term("04/03/2035", `^\d{1,2}/\d{1,2}/\d{4}$`),
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusCreated,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
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
					Given("LPA team with members exists").
					UponReceiving("A request to create a task for a team").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/lpa-api/v1/cases/800/tasks"),
						Headers: dsl.MapMatcher{
							"Content-Type": dsl.String("application/json"),
						},
						Body: map[string]interface{}{
							"assigneeId":  23,
							"type":        "Check Application",
							"name":        "A title",
							"description": "A description",
							"dueDate":     dsl.Term("04/03/2035", `^\d{1,2}/\d{1,2}/\d{4}$`),
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusCreated,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				err := client.CreateTask(Context{Context: context.Background()}, 800, tc.task)
				if (tc.expectedError) == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(pact.Server.Port), err)
				}
				return nil
			}))
		})
	}
}
