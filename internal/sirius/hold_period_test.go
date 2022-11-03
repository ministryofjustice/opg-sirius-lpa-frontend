package sirius

import (
	"context"
	"fmt"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestHoldPeriod(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

	testCases := []struct {
		name             string
		setup            func()
		expectedResponse HoldPeriod
		expectedError    func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a case assigned which has an investigation on hold").
					UponReceiving("A request for the hold period").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/lpa-api/v1/hold-periods/175"),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: dsl.Like(map[string]interface{}{
							"id": dsl.Like(175),
							"investigation": dsl.Like(map[string]interface{}{
								"id":                        1,
								"investigationTitle":        dsl.String("Test title"),
								"type":                      dsl.String("Normal"),
								"investigationReceivedDate": dsl.String("23/01/2022"),
							}),
							"reason": "Police Investigation",
						}),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
			},
			expectedResponse: HoldPeriod{
				ID: 175,
				Investigation: Investigation{
					ID:           1,
					Title:        "Test title",
					Type:         "Normal",
					DateReceived: DateString("2022-01-23"),
				},
				Reason: "Police Investigation",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				hp, err := client.HoldPeriod(Context{Context: context.Background()}, 175)

				assert.Equal(t, tc.expectedResponse, hp)
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
