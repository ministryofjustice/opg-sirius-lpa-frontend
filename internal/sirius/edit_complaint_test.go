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

func TestEditComplaint(t *testing.T) {
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
					Given("A complaint exists").
					UponReceiving("A request to edit the complaint").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPut,
						Path:   matchers.String("/lpa-api/v1/complaints/986"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: matchers.Like(map[string]interface{}{
							"category":             "02",
							"description":          "This is seriously bad",
							"receivedDate":         "05/04/2022",
							"severity":             "Major",
							"investigatingOfficer": "Test Officer",
							"complainantName":      "Someones name",
							"subCategory":          "18",
							"complainantCategory":  "LPA_DONOR",
							"origin":               "PHONE",
							"summary":              "This and that",
							"resolution":           "complaint upheld",
							"resolutionInfo":       "Because...",
							"resolutionDate":       "07/06/2022",
							"compensationType":     "COMPENSATORY",
							"compensationAmount":   "150.00",
						}),
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				err := client.EditComplaint(Context{Context: context.Background()}, 986, Complaint{
					Category:             "02",
					Description:          "This is seriously bad",
					ReceivedDate:         DateString("2022-04-05"),
					Severity:             "Major",
					InvestigatingOfficer: "Test Officer",
					ComplainantName:      "Someones name",
					SubCategory:          "18",
					ComplainantCategory:  "LPA_DONOR",
					Origin:               "PHONE",
					Summary:              "This and that",
					Resolution:           "complaint upheld",
					ResolutionInfo:       "Because...",
					ResolutionDate:       DateString("2022-06-07"),
					CompensationType:     "COMPENSATORY",
					CompensationAmount:   "150.00",
				})

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
