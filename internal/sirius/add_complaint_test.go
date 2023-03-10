package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestAddComplaint(t *testing.T) {
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
					Given("I have a pending case assigned").
					UponReceiving("A request to add a complaint to the case").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/lpa-api/v1/lpas/800/complaints"),
						Headers: dsl.MapMatcher{
							"Content-Type": dsl.String("application/json"),
						},
						Body: dsl.Like(map[string]interface{}{
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
					WillRespondWith(dsl.Response{
						Status:  http.StatusCreated,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body:    dsl.Like(map[string]interface{}{"id": dsl.Integer()}),
					})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

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
					assert.Equal(t, tc.expectedError(pact.Server.Port), err)
				}
				return nil
			}))
		})
	}
}
