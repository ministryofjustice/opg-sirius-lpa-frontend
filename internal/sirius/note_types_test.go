package sirius

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestNoteTypes(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

	testCases := []struct {
		name             string
		setup            func()
		cookies          []*http.Cookie
		expectedResponse []string
		expectedError    func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("Some note types exist").
					UponReceiving("A request for note types").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/lpa-api/v1/note-types/lpa"),
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Body:    dsl.EachLike(dsl.String("Application processing"), 1),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			expectedResponse: []string{"Application processing"},
		},
		{
			name: "Unauthorized",
			setup: func() {
				pact.
					AddInteraction().
					Given("Some note types exist").
					UponReceiving("A request for note types without cookies").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/lpa-api/v1/note-types/lpa"),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusUnauthorized,
					})
			},
			expectedError: func(port int) error {
				return StatusError{
					Code:   http.StatusUnauthorized,
					URL:    fmt.Sprintf("http://localhost:%d/lpa-api/v1/note-types/lpa", port),
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

				types, err := client.NoteTypes(getContext(tc.cookies))

				assert.Equal(t, tc.expectedResponse, types)
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
