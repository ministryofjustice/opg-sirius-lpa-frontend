package sirius

import (
	"context"
	"fmt"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestInvestigationOffHold(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		setup            func()
		expectedResponse Investigation
		expectedError    func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a case assigned which has an investigation open").
					UponReceiving("A request for the investigation").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/investigations/300"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.Like(map[string]interface{}{
							"id":                        matchers.Like(300),
							"investigationTitle":        matchers.String("Test title"),
							"additionalInformation":     matchers.String("Some test info"),
							"type":                      matchers.String("Normal"),
							"investigationReceivedDate": matchers.String("23/01/2022"),
						}),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: Investigation{
				ID:           300,
				Title:        "Test title",
				Information:  "Some test info",
				Type:         "Normal",
				DateReceived: DateString("2022-01-23"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				investigation, err := client.Investigation(Context{Context: context.Background()}, 300)

				assert.Equal(t, tc.expectedResponse, investigation)
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

func TestInvestigationOnHold(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		setup            func()
		expectedResponse Investigation
		expectedError    func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a case assigned which has an investigation on hold").
					UponReceiving("A request for the investigation which is on hold").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/investigations/301"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.Like(map[string]interface{}{
							"id":                        matchers.Like(301),
							"investigationTitle":        matchers.String("Test title"),
							"additionalInformation":     matchers.String("Some test info"),
							"type":                      matchers.String("Normal"),
							"investigationReceivedDate": matchers.String("23/01/2022"),
							"isOnHold":                  true,
							"holdPeriods": matchers.Like([]map[string]interface{}{
								{
									"id":        matchers.Like(175),
									"reason":    matchers.String("Police Investigation"),
									"startDate": matchers.String("25/01/2022"),
								},
							}),
						}),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: Investigation{
				ID:           301,
				Title:        "Test title",
				Information:  "Some test info",
				Type:         "Normal",
				DateReceived: DateString("2022-01-23"),
				IsOnHold:     true,
				HoldPeriods: []HoldPeriod{
					{
						ID:        175,
						Reason:    "Police Investigation",
						StartDate: DateString("2022-01-25"),
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				investigation, err := client.Investigation(Context{Context: context.Background()}, 301)

				assert.Equal(t, tc.expectedResponse, investigation)
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
