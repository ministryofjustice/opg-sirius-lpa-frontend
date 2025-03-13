package sirius

import (
	"context"
	"fmt"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllocateCases(t *testing.T) {
	t.Parallel()

	pact, err := newPact2()
	assert.NoError(t, err)

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
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPut,
						Path:   matchers.String("/lpa-api/v1/users/47/cases/800"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]interface{}{
							"data": []map[string]interface{}{{
								"id":       800,
								"caseType": "LPA",
							}},
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
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
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPut,
						Path:   matchers.String("/lpa-api/v1/users/47/cases/800+801+802"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
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
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			allocations: []CaseAllocation{{ID: 800, CaseType: "LPA"}, {ID: 801, CaseType: "LPA"}, {ID: 802, CaseType: "EPA"}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				err := client.AllocateCases(Context{Context: context.Background()}, 47, tc.allocations)
				if (tc.expectedError) == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(config.Port), err)
				}
				return nil
			}))
		})
	}
}
