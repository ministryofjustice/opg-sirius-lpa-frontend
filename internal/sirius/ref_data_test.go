package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestRefDataByCategory(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

	testCases := []struct {
		name             string
		setup            func()
		expectedResponse []RefDataItem
		expectedError    func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("Some warning types exist").
					UponReceiving("A request for warning types").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/lpa-api/v1/reference-data/warningType"),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: dsl.EachLike(map[string]interface{}{
							"handle": dsl.String("Complaint Received"),
							"label":  dsl.String("Complaint Received"),
						}, 1),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
			},
			expectedResponse: []RefDataItem{
				{
					Handle: "Complaint Received",
					Label:  "Complaint Received",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				types, err := client.RefDataByCategory(Context{Context: context.Background()}, WarningTypeCategory)

				assert.Equal(t, tc.expectedResponse, types)
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
