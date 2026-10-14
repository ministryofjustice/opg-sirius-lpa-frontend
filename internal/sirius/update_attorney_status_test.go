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

func TestChangeAttorneyStatus(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name          string
		setup         func()
		expectedError func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("A digital LPA exists").
					UponReceiving("A request for changing the status of the digital LPA attorney").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPut,
						Path:   matchers.String("/lpa-api/v1/digital-lpas/M-1234-9876-4567/attorney-status"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]interface{}{
							"attorneyStatuses": []map[string]interface{}{{
								"uid":           "cf128305-37c8-4ceb-bedf-89ed5f4ae661",
								"status":        "removed",
								"removedReason": "BANKRUPT",
							}},
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusNoContent,
					})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", config.Port))

				err := client.ChangeAttorneyStatus(Context{Context: context.Background()}, "M-1234-9876-4567", []AttorneyUpdatedStatus{{UID: "cf128305-37c8-4ceb-bedf-89ed5f4ae661", Status: "removed", RemovedReason: "BANKRUPT"}})

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
