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

func TestChangeCertificateProviderDetails(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name          string
		changeData    ChangeCertificateProviderDetails
		setup         func()
		expectedError func(int) error
	}{
		{
			name: "OK",
			changeData: ChangeCertificateProviderDetails{
				FirstNames: "Janelle",
				LastName:   "O'Reilly",
				Address: Address{
					Line1:    "387 Quarry Lane",
					Line2:    "Harberwell",
					Line3:    "Daugherty",
					Town:     "Greater Manchester",
					Postcode: "M1 4MG",
					Country:  "GB",
				},
				Phone:    "056 8656 8956",
				Email:    "Janelle.Oreilly@example.com",
				SignedAt: "2025-02-21",
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("A digital LPA exists").
					UponReceiving("A request for changing certificate provider details").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPut,
						Path:   matchers.String("/lpa-api/v1/digital-lpas/M-1234-9876-4567/change-certificate-provider-details"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]interface{}{
							"firstNames": "Janelle",
							"lastName":   "O'Reilly",
							"address": map[string]string{
								"addressLine1": "387 Quarry Lane",
								"addressLine2": "Harberwell",
								"addressLine3": "Daugherty",
								"town":         "Greater Manchester",
								"postcode":     "M1 4MG",
								"country":      "GB",
							},
							"phoneNumber": "056 8656 8956",
							"email":       "Janelle.Oreilly@example.com",
							"signedAt":    "21/02/2025",
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

				err := client.ChangeCertificateProviderDetails(
					Context{Context: context.Background()},
					"M-1234-9876-4567",
					tc.changeData,
				)

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
