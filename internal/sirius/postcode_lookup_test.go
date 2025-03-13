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

func TestPostcodeLookup(t *testing.T) {
	t.Parallel()

	pact, err := newPact2()
	assert.NoError(t, err)

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
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/postcode-lookup"),
						Query: matchers.MapMatcher{
							"postcode": matchers.String("SW1A 1AA"),
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: matchers.EachLike(map[string]interface{}{
							"addressLine1": matchers.Like("Office of the Public Guardian"),
							"addressLine2": matchers.Like("1 Something Street"),
							"addressLine3": matchers.Like("Someborough"),
							"town":         matchers.Like("Someton"),
							"postcode":     matchers.Like("SW1A 1AA"),
							"description":  matchers.Like("Office of the Public Guardian, 1 Something Street, Someborough"),
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

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				addresses, err := client.PostcodeLookup(Context{Context: context.Background()}, "SW1A 1AA")
				assert.Equal(t, tc.expectedResponse, addresses)
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
