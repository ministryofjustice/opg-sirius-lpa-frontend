package sirius

import (
	"context"
	"fmt"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestChangeDraft(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

	testCases := []struct {
		name          string
		changeData    ChangeDraft
		setup         func()
		expectedError func(int) error
	}{
		{
			name: "OK",
			changeData: ChangeDraft{
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
				Phone: "12345678",
				Email: "test@test.com",
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("A digital LPA exists").
					UponReceiving("A request for changing draft").
					WithRequest(dsl.Request{
						Method: http.MethodPut,
						Path:   dsl.String("/lpa-api/v1/digital-lpas/M-1234-9876-4567/change-draft"),
						Headers: dsl.MapMatcher{
							"Content-Type": dsl.String("application/json"),
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

				err := client.ChangeDraft(Context{Context: context.Background()}, "M-1234-9876-4567", tc.changeData)

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
