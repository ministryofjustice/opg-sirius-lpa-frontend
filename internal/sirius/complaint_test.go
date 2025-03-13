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

func TestComplaint(t *testing.T) {
	t.Parallel()

	pact, err := newPact2()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		setup            func()
		expectedResponse Complaint
		expectedError    func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("A complaint exists").
					UponReceiving("A request for the complaint").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/complaints/986"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.Like(map[string]interface{}{
							"category":             matchers.String("01"),
							"description":          matchers.String("This is seriously bad"),
							"receivedDate":         matchers.String("05/04/2022"),
							"severity":             matchers.String("Major"),
							"investigatingOfficer": matchers.String("Test Officer"),
							"subCategory":          matchers.String("07"),
							"summary":              matchers.String("This and that"),
						}),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: Complaint{
				Category:             "01",
				Description:          "This is seriously bad",
				ReceivedDate:         DateString("2022-04-05"),
				Severity:             "Major",
				InvestigatingOfficer: "Test Officer",
				SubCategory:          "07",
				Summary:              "This and that",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				complaint, err := client.Complaint(Context{Context: context.Background()}, 986)

				assert.Equal(t, tc.expectedResponse, complaint)
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
