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

func TestGetCombinedEvents(t *testing.T) {
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
			name: "OK - returns combined events from both sources",
			setup: func() {
				pact.
					AddInteraction().
					Given("A digital LPA exists").
					UponReceiving("A request for the combined event history").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.Like("/lpa-api/v1/digital-lpas/M-1234-9876-4567/events"),
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: matchers.EachLike(map[string]interface{}{
							"source": matchers.Like("sirius"),
						}, 1),
					})
			},
			expectedResponse: APIEvent{
				Event{
					ChangeSet:  []interface{}(nil),
					CreatedOn:  "",
					Entity:     interface{}(nil),
					Source:     "sirius",
					SourceType: "",
					Type:       "",
					User:       EventUser{DisplayName: ""},
					UUID:       "",
					Applied:    "",
					DateTime:   "",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				events, err := client.GetCombinedEvents(Context{Context: context.Background()}, "M-1234-5678-9012")

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
