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

func TestChangeAttorneyDetails(t *testing.T) {
	t.Parallel()

	pact, err := newPact2()
	assert.NoError(t, err)

	attorneyUid := "302b05c7-896c-4290-904e-2005e4f1e81e"

	testCases := []struct {
		name          string
		changeData    ChangeAttorneyDetails
		setup         func()
		expectedError func(int) error
	}{
		{
			name: "OK",
			changeData: ChangeAttorneyDetails{
				FirstNames:  "Jake",
				LastName:    "Sullivan",
				DateOfBirth: "2000-11-12",
				Address: Address{
					Line1:    "Flat Number",
					Line2:    "Building",
					Line3:    "Road Name",
					Town:     "South Bend",
					Postcode: "AI1 6VW",
					Country:  "GB",
				},
				Phone:    "12345678",
				Email:    "test@test.com",
				SignedAt: "2024-10-01",
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("A digital LPA in statutory waiting period").
					UponReceiving("A request for changing attorney details").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPut,
						Path:   matchers.String("/lpa-api/v1/digital-lpas/M-1111-2222-3333/attorney/" + attorneyUid + "/change-details"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]interface{}{
							"firstNames":  "Jake",
							"lastName":    "Sullivan",
							"dateOfBirth": "12/11/2000",
							"address": map[string]string{
								"addressLine1": "Flat Number",
								"addressLine2": "Building",
								"addressLine3": "Road Name",
								"town":         "South Bend",
								"postcode":     "AI1 6VW",
								"country":      "GB",
							},
							"phoneNumber": "12345678",
							"email":       "test@test.com",
							"signedAt":    "01/10/2024",
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

				err := client.ChangeAttorneyDetails(Context{Context: context.Background()}, "M-1111-2222-3333", attorneyUid, tc.changeData)

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
