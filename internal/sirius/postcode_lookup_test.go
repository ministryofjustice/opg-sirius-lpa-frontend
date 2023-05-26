package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestPostcodeLookup(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

	testCases := []struct {
		name             string
		setup            func()
		expectedResponse []PostcodeLookupAddress
		expectedError    func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					UponReceiving("A postcode search").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/lpa-api/v1/postcode-lookup"),
						Query: dsl.MapMatcher{
							"postcode": dsl.String("SW1A 1AA"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body: dsl.EachLike(map[string]interface{}{
							"addressLine1": dsl.Like("Office of the Public Guardian"),
							"addressLine2": dsl.Like("1 Something Street"),
							"addressLine3": dsl.Like("Someborough"),
							"town":         dsl.Like("Someton"),
							"postcode":     dsl.Like("SW1A 1AA"),
							"description":  dsl.Like("Office of the Public Guardian, 1 Something Street, Someborough"),
						}, 1),
					})
			},
			expectedResponse: []PostcodeLookupAddress{
				{
					Line1:       "Office of the Public Guardian",
					Line2:       "1 Something Street",
					Line3:       "Someborough",
					Town:        "Someton",
					Postcode:    "SW1A 1AA",
					Description: "Office of the Public Guardian, 1 Something Street, Someborough",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				addresses, err := client.PostcodeLookup(Context{Context: context.Background()}, "SW1A 1AA")
				assert.Equal(t, tc.expectedResponse, addresses)
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
