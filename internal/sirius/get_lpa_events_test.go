package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

func TestGetEvents(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name           string
		donorID        string
		caseIDs        []string
		setup          func()
		assertResponse func(*testing.T, LpaEventsResponse)
		expectedError  func(int) error
	}{
		{
			name:    "OK - returns events for all cases when none selected",
			donorID: "189",
			caseIDs: nil,
			setup: func() {
				pact.
					AddInteraction().
					Given("A donor exists with multiple cases and event history").
					UponReceiving("A request for the donor event history without case filters").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/persons/189/events"),
						Query: matchers.MapMatcher{
							"sort":  matchers.String("id:desc"),
							"limit": matchers.String("999"),
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: matchers.Like(map[string]any{
							"limit": matchers.Integer(999),
							"metadata": matchers.Like(map[string]any{
								"caseIds":     matchers.EachLike(map[string]int{}, 1),
								"sourceTypes": matchers.EachLike(map[string]any{}, 1),
							}),
							"pages": matchers.Like(map[string]any{
								"current": matchers.Integer(1),
								"total":   matchers.Integer(1),
							}),
							"total":  matchers.Integer(3),
							"events": matchers.EachLike(map[string]any{}, 1),
						}),
					})
			},
			assertResponse: func(t *testing.T, resp LpaEventsResponse) {
				if !assert.NotNil(t, resp) {
					return
				}
			},
			expectedError: nil,
		},
		// add case with selected case IDs
		// do we need an error case here?
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				resp, err := client.GetEvents(Context{Context: context.Background()}, tc.donorID, tc.caseIDs)

				if tc.assertResponse != nil {
					tc.assertResponse(t, resp)
				} else {
					assert.Nil(t, resp)
				}

				if tc.expectedError == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(config.Port), err)
				}

				return nil
			}))
		})
	}
}
