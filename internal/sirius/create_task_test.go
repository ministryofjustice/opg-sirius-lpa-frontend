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
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending case assigned").
					UponReceiving("A request to create a task").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/lpa-api/v1/tasks"),
						Body: map[string]interface{}{
							"caseId":      dsl.Like(800),
							"assigneeId":  dsl.Like(1),
							"type":        dsl.String("Change of Address"),
							"name":        dsl.String("Something"),
							"description": dsl.String("More words"),
							"dueDate":     dsl.Term("04/05/2731", `^\d{1,2}/\d{1,2}/\d{4}$`),
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

				err := client.CreateTask(Context{Context: context.Background()}, TaskRequest{
					CaseID:      800,
					AssigneeID:  1,
					Type:        "Change of Address",
					Name:        "Something",
					Description: "More words",
					DueDate:     "9999-05-04",
				})
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

func TestCreateTaskWithEmptyDescription(t *testing.T) {
	client := NewClient(http.DefaultClient, "")

	err := client.CreateTask(Context{}, TaskRequest{
		CaseID:      800,
		AssigneeID:  1,
		Type:        "Change of Address",
		Name:        "Something",
		Description: "  ",
		DueDate:     "9999-05-04",
	})

	assert.Equal(t, ValidationError{Field: FieldErrors{"description": {"": "Value can't be empty"}}}, err)
}
