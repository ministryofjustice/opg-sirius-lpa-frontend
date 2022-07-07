package sirius

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestAssignTask(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

	testCases := []struct {
		name          string
		setup         func()
		cookies       []*http.Cookie
		expectedError func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a case with an open task assigned").
					UponReceiving("A request to assign a task").
					WithRequest(dsl.Request{
						Method: http.MethodPut,
						Path:   dsl.String("/api/v1/tasks/990"),
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
							"Content-Type":        dsl.String("application/json"),
						},
						Body: dsl.Like(map[string]interface{}{
							"id":         dsl.Like(990),
							"assigneeId": dsl.Like(47),
						}),
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
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
					Given("I have a case with an open task assigned").
					UponReceiving("A request to assign a task without cookies").
					WithRequest(dsl.Request{
						Method: http.MethodPut,
						Path:   dsl.String("/api/v1/tasks/990"),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusUnauthorized,
					})
			},
			expectedError: func(port int) error {
				return StatusError{
					Code:   http.StatusUnauthorized,
					URL:    fmt.Sprintf("http://localhost:%d/api/v1/tasks/990", port),
					Method: http.MethodPut,
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				err := client.AssignTask(getContext(tc.cookies), 47, 990)

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
