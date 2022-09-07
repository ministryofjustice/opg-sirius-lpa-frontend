package sirius

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestAvailableStatuses(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

	testCases := []struct {
		name             string
		setup            func()
		cookies          []*http.Cookie
		caseType         CaseType
		expectedResponse []string
		expectedError    func(int) error
	}{
		{
			name: "OK LPA",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending case assigned").
					UponReceiving("A request for the LPA's available statuses").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/lpa-api/v1/lpas/800/available-statuses"),
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Body:    dsl.EachLike(dsl.String("Perfect"), 1),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			caseType:         CaseTypeLpa,
			expectedResponse: []string{"Perfect"},
		},
		{
			name: "OK EPA",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending EPA assigned").
					UponReceiving("A request for the EPA's available statuses").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/lpa-api/v1/epas/800/available-statuses"),
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Body:    dsl.EachLike(dsl.String("Perfect"), 1),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			caseType:         CaseTypeEpa,
			expectedResponse: []string{"Perfect"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				caseitem, err := client.AvailableStatuses(getContext(tc.cookies), 800, tc.caseType)

				assert.Equal(t, tc.expectedResponse, caseitem)
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
