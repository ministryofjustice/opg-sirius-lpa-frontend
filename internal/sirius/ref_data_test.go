package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestRefDataByCategoryWarningTypes(t *testing.T) {
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
						Path:   dsl.String(fmt.Sprintf("/lpa-api/v1/reference-data/%s", WarningTypeCategory)),
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

func TestRefDataByCategoryPaymentSources(t *testing.T) {
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
					Given("Some payment sources exist").
					UponReceiving("A request for payment source ref data").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String(fmt.Sprintf("/lpa-api/v1/reference-data/%s", PaymentSourceCategory)),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: dsl.Like([]map[string]interface{}{
							{
								"handle":         dsl.String("PHONE"),
								"label":          dsl.String("Paid over the phone"),
								"userSelectable": true,
							},
							{
								"handle":         dsl.String("ONLINE"),
								"label":          dsl.String("Paid online"),
								"userSelectable": true,
							},
							{
								"handle":         dsl.String("MAKE"),
								"label":          dsl.String("Paid through Make an LPA"),
								"userSelectable": false,
							},
							{
								"handle":         dsl.String("OTHER"),
								"label":          dsl.String("Paid through other method"),
								"userSelectable": false,
							},
							{
								"handle":         dsl.String("MIGRATED"),
								"label":          dsl.String("Payment was migrated"),
								"userSelectable": false,
							},
						}),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
			},
			expectedResponse: []RefDataItem{
				{
					Handle:         "PHONE",
					Label:          "Paid over the phone",
					UserSelectable: true,
				},
				{
					Handle:         "ONLINE",
					Label:          "Paid online",
					UserSelectable: true,
				},
				{
					Handle:         "MAKE",
					Label:          "Paid through Make an LPA",
					UserSelectable: false,
				},
				{
					Handle:         "OTHER",
					Label:          "Paid through other method",
					UserSelectable: false,
				},
				{
					Handle:         "MIGRATED",
					Label:          "Payment was migrated",
					UserSelectable: false,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				types, err := client.RefDataByCategory(Context{Context: context.Background()}, PaymentSourceCategory)

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
