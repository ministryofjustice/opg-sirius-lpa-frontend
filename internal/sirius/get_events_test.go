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

func TestGetEvents(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		setup            func()
		expectedResponse any
		expectedError    func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("A donor with a digital LPA with events exists").
					UponReceiving("A request for the donor's event history").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.Like("/lpa-api/v1/persons/33/events?filter=case%3A66&sort=id%3Adesc"),
						Query: matchers.MapMatcher{
							"filter": matchers.Like("case:66"),
							"sort":   matchers.Like("id:desc"),
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: matchers.Like(map[string]interface{}{
							"events": matchers.EachLike(struct{}{}, 0),
						}),
					})
			},
			expectedResponse: []interface{}{map[string]interface{}{}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				events, err := client.GetEvents(Context{Context: context.Background()}, 33, 66)

				assert.Equal(t, tc.expectedResponse, events)
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
