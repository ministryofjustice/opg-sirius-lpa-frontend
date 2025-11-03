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

func TestChangeTrustCorporationDetails(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	trustCorporationUid := "302b05c7-896c-4290-904e-2005e4f1e81e"

	testCases := []struct {
		name          string
		changeData    ChangeTrustCorporationDetails
		setup         func()
		expectedError func(int) error
	}{
		{
			name: "OK",
			changeData: ChangeTrustCorporationDetails{
				Name: "Trust Ltd.",
				Address: Address{
					Line1:    "Flat Number",
					Line2:    "Building",
					Line3:    "Road Name",
					Town:     "South Bend",
					Postcode: "AI1 6VW",
					Country:  "GB",
				},
				Phone:         "12345678",
				Email:         "test@test.com",
				CompanyNumber: "123456789",
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("A digital LPA in statutory waiting period").
					UponReceiving("A request for changing trust corporation details").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPut,
						Path:   matchers.String("/lpa-api/v1/digital-lpas/M-1111-2222-3333/trust-corporation/" + trustCorporationUid + "/change-details"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]interface{}{
							"name": "Trust Ltd.",
							"address": map[string]string{
								"addressLine1": "Flat Number",
								"addressLine2": "Building",
								"addressLine3": "Road Name",
								"town":         "South Bend",
								"postcode":     "AI1 6VW",
								"country":      "GB",
							},
							"phoneNumber":   "12345678",
							"email":         "test@test.com",
							"companyNumber": "123456789",
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
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				err := client.ChangeTrustCorporationDetails(Context{Context: context.Background()}, "M-1111-2222-3333", trustCorporationUid, tc.changeData)

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
