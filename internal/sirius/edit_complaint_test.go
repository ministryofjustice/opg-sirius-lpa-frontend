package sirius

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestEditComplaint(t *testing.T) {
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
					Given("A complaint exists").
					UponReceiving("A request to edit the complaint").
					WithRequest(dsl.Request{
						Method: http.MethodPut,
						Path:   dsl.String("/lpa-api/v1/complaints/986"),
						Body: dsl.Like(map[string]interface{}{
							"category":       "01",
							"description":    "This is seriously bad",
							"receivedDate":   "05/04/2022",
							"severity":       "Major",
							"subCategory":    "07",
							"summary":        "This and that",
							"resolution":     "complaint upheld",
							"resolutionInfo": "We did stuff",
							"resolutionDate": "07/06/2022",
						}),
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
							"Content-Type":        dsl.String("application/json"),
						},
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
					Given("A complaint exists").
					UponReceiving("A request to edit the complaint without cookies").
					WithRequest(dsl.Request{
						Method: http.MethodPut,
						Path:   dsl.String("/lpa-api/v1/complaints/986"),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusUnauthorized,
					})
			},
			expectedError: func(port int) error {
				return StatusError{
					Code:   http.StatusUnauthorized,
					URL:    fmt.Sprintf("http://localhost:%d/lpa-api/v1/complaints/986", port),
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

				err := client.EditComplaint(getContext(tc.cookies), 986, Complaint{
					Category:       "01",
					Description:    "This is seriously bad",
					ReceivedDate:   DateString("2022-04-05"),
					Severity:       "Major",
					SubCategory:    "07",
					Summary:        "This and that",
					Resolution:     "complaint upheld",
					ResolutionInfo: "We did stuff",
					ResolutionDate: DateString("2022-06-07"),
				})

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
