package sirius

import (
	"context"
	"fmt"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestEditInvestigation(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

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
					Given("I have a case assigned which has an investigation open").
					UponReceiving("A request to edit the investigation").
					WithRequest(dsl.Request{
						Method: http.MethodPut,
						Path:   dsl.String("/lpa-api/v1/investigations/300"),
						Headers: dsl.MapMatcher{
							"Content-Type": dsl.String("application/json"),
						},
						Body: dsl.Like(map[string]interface{}{
							"investigationTitle":        "Test title",
							"additionalInformation":     "Some test info",
							"type":                      "Normal",
							"investigationReceivedDate": "23/01/2022",
							"reportApprovalDate":        "05/04/2022",
							"riskAssessmentDate":        "05/04/2022",
							"reportApprovalOutcome":     "Court Application",
							"investigationClosureDate":  "05/04/2022",
						}),
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				err := client.EditInvestigation(Context{Context: context.Background()}, 300, Investigation{
					Title:                    "Test title",
					Information:              "Some test info",
					Type:                     "Normal",
					DateReceived:             DateString("2022-01-23"),
					ApprovalDate:             DateString("2022-04-05"),
					RiskAssessmentDate:       DateString("2022-04-05"),
					ApprovalOutcome:          "Court Application",
					InvestigationClosureDate: DateString("2022-04-05"),
				})

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
