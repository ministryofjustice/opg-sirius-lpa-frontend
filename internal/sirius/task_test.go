package sirius

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestTask(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

	testCases := []struct {
		name             string
		setup            func()
		cookies          []*http.Cookie
		expectedResponse Task
		expectedError    func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a case with an open task assigned").
					UponReceiving("A request for a task").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/lpa-api/v1/tasks/990"),
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: dsl.Like(map[string]interface{}{
							"id":      dsl.Like(990),
							"status":  dsl.String("Not Started"),
							"dueDate": dsl.String("10/01/2022"),
							"name":    dsl.String("Create physical case file"),
							"caseItems": dsl.EachLike(map[string]interface{}{
								"uId":      dsl.String("7000-0000-0001"),
								"caseType": dsl.String("LPA"),
							}, 1),
						}),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			expectedResponse: Task{
				ID:      990,
				Status:  "Not Started",
				DueDate: DateString("2022-01-10"),
				Name:    "Create physical case file",
				CaseItems: []Case{{
					UID:      "7000-0000-0001",
					CaseType: "LPA",
				}},
			},
		},
		{
			name: "Unauthorized",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a case with an open task assigned").
					UponReceiving("A request for a task without cookies").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/lpa-api/v1/tasks/990"),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusUnauthorized,
					})
			},
			expectedError: func(port int) error {
				return StatusError{
					Code:   http.StatusUnauthorized,
					URL:    fmt.Sprintf("http://localhost:%d/lpa-api/v1/tasks/990", port),
					Method: http.MethodGet,
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				task, err := client.Task(getContext(tc.cookies), 990)

				assert.Equal(t, tc.expectedResponse, task)
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
