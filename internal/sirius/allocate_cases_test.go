package sirius

import (
	"context"
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
						Headers: dsl.MapMatcher{
							"Content-Type": dsl.String("application/json"),
						},
						Body: map[string]interface{}{
							"data": []map[string]interface{}{{
								"id":       800,
								"caseType": "LPA",
							}},
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
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
						Headers: dsl.MapMatcher{
							"Content-Type": dsl.String("application/json"),
						},
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
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
			},
			allocations: []CaseAllocation{{ID: 800, CaseType: "LPA"}, {ID: 801, CaseType: "LPA"}, {ID: 802, CaseType: "EPA"}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				err := client.AllocateCases(Context{Context: context.Background()}, 47, tc.allocations)
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
