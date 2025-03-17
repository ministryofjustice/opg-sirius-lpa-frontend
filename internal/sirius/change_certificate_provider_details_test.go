package sirius

import (
	"context"
	"fmt"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestChangeCertificateProviderDetails(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

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
					WithRequest(dsl.Request{
						Method: http.MethodPut,
						Path:   dsl.String("/lpa-api/v1/digital-lpas/M-1234-9876-4567/change-certificate-provider-details"),
						Headers: dsl.MapMatcher{
							"Content-Type": dsl.String("application/json"),
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
					WillRespondWith(dsl.Response{
						Status: http.StatusNoContent,
					})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				err := client.ChangeCertificateProviderDetails(
					Context{Context: context.Background()},
					"M-1234-9876-4567",
					tc.changeData,
				)

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
