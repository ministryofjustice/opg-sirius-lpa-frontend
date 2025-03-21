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

func TestAddComplaint(t *testing.T) {
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
					Given("I have a pending case assigned").
					UponReceiving("A request to add a complaint to the case").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPost,
						Path:   matchers.String("/lpa-api/v1/lpas/800/complaints"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: matchers.Like(map[string]interface{}{
							"category":             "02",
							"description":          "A description",
							"receivedDate":         "05/04/2022",
							"severity":             "Major",
							"investigatingOfficer": "Test Officer",
							"complainantName":      "Someones name",
							"subCategory":          "18",
							"complainantCategory":  "LPA_DONOR",
							"origin":               "PHONE",
							"summary":              "A title",
						}),
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusCreated,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body:    matchers.Like(map[string]interface{}{"id": matchers.Integer(1)}),
					})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				err := client.AddComplaint(Context{Context: context.Background()}, 800, CaseTypeLpa, Complaint{
					Category:             "02",
					Description:          "A description",
					ReceivedDate:         DateString("2022-04-05"),
					Severity:             "Major",
					InvestigatingOfficer: "Test Officer",
					ComplainantName:      "Someones name",
					SubCategory:          "18",
					ComplainantCategory:  "LPA_DONOR",
					Origin:               "PHONE",
					Summary:              "A title",
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
