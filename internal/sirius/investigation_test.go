package sirius

import (
	"context"
	"fmt"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestInvestigationOffHold(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

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
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/lpa-api/v1/investigations/300"),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: dsl.Like(map[string]interface{}{
							"id":                        dsl.Like(300),
							"investigationTitle":        dsl.String("Test title"),
							"additionalInformation":     dsl.String("Some test info"),
							"type":                      dsl.String("Normal"),
							"investigationReceivedDate": dsl.String("23/01/2022"),
						}),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
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

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				investigation, err := client.Investigation(Context{Context: context.Background()}, 300)

				assert.Equal(t, tc.expectedResponse, investigation)
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

func TestInvestigationOnHold(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

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
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/lpa-api/v1/investigations/301"),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: dsl.Like(map[string]interface{}{
							"id":                        dsl.Like(301),
							"investigationTitle":        dsl.String("Test title"),
							"additionalInformation":     dsl.String("Some test info"),
							"type":                      dsl.String("Normal"),
							"investigationReceivedDate": dsl.String("23/01/2022"),
							"isOnHold":                  true,
							"holdPeriods": dsl.Like([]map[string]interface{}{
								{
									"id":        dsl.Like(175),
									"reason":    dsl.String("Police Investigation"),
									"startDate": dsl.String("25/01/2022"),
								},
							}),
						}),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
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

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				investigation, err := client.Investigation(Context{Context: context.Background()}, 301)

				assert.Equal(t, tc.expectedResponse, investigation)
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
