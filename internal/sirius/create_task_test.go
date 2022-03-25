package sirius

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	pact := newPact()
	defer pact.Teardown()

	testCases := []struct {
		name          string
		setup         func()
		cookies       []*http.Cookie
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
						Path:   dsl.String("/api/v1/tasks"),
						Body: dsl.Like(map[string]interface{}{
							"caseId":      800,
							"assigneeId":  1,
							"type":        "Change of Address",
							"name":        "Something",
							"description": "More words",
							"dueDate":     "04/05/9999",
						}),
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
							"Content-Type":        dsl.String("application/json"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusCreated,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
		},
		{
			name: "Unauthorized",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending case assigned").
					UponReceiving("A request to create a task without cookies").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/api/v1/tasks"),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusUnauthorized,
					})
			},
			expectedError: func(port int) error {
				return StatusError{
					Code:   http.StatusUnauthorized,
					URL:    fmt.Sprintf("http://localhost:%d/api/v1/tasks", port),
					Method: http.MethodPost,
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				err := client.CreateTask(getContext(tc.cookies), Task{
					CaseID:      800,
					AssigneeID:  1,
					Type:        "Change of Address",
					Name:        "Something",
					Description: "More words",
					DueDate:     "04/05/9999",
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

	err := client.CreateTask(Context{}, Task{
		CaseID:      800,
		AssigneeID:  1,
		Type:        "Change of Address",
		Name:        "Something",
		Description: "",
		DueDate:     "04/05/9999",
	})

	assert.Equal(t, ValidationError{Errors: ValidationErrors{"description": {"": "Value is required"}}}, err)
}
