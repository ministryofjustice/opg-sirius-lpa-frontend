package sirius

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestAllocateCases(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

	testCases := []struct {
		name          string
		setup         func()
		cookies       []*http.Cookie
		allocations   []CaseAllocation
		expectedError func(int) error
		file          *NoteFile
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending case assigned").
					UponReceiving("A request to change the assignee of the case").
					WithRequest(dsl.Request{
						Method: http.MethodPut,
						Path:   dsl.String("/lpa-api/v1/users/47/cases/800"),
						Body: map[string]interface{}{
							"data": []map[string]interface{}{{
								"id":       800,
								"caseType": "LPA",
							}},
						},
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
			allocations: []CaseAllocation{{ID: 800, CaseType: "LPA"}},
		},
		{
			name: "OK multiple",
			setup: func() {
				pact.
					AddInteraction().
					Given("Multiple cases exist").
					UponReceiving("A request to change the assignee of multiple cases").
					WithRequest(dsl.Request{
						Method: http.MethodPut,
						Path:   dsl.String("/lpa-api/v1/users/47/cases/800+801+802"),
						Body: map[string]interface{}{
							"data": []map[string]interface{}{{
								"id":       800,
								"caseType": "LPA",
							}, {
								"id":       801,
								"caseType": "LPA",
							}, {
								"id":       802,
								"caseType": "EPA",
							}},
						},
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
			allocations: []CaseAllocation{{ID: 800, CaseType: "LPA"}, {ID: 801, CaseType: "LPA"}, {ID: 802, CaseType: "EPA"}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				err := client.AllocateCases(getContext(tc.cookies), 47, tc.allocations)
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
