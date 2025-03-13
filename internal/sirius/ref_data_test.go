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

func TestRefDataByCategory(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		setup            func()
		category         string
		expectedResponse []RefDataItem
		expectedError    func(int) error
	}{
		{
			name:     "Warning Types",
			category: WarningTypeCategory,
			setup: func() {
				pact.
					AddInteraction().
					UponReceiving("A request for warning type ref data").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String(fmt.Sprintf("/lpa-api/v1/reference-data/%s", WarningTypeCategory)),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.EachLike(map[string]interface{}{
							"handle": matchers.String("Complaint Received"),
							"label":  matchers.String("Complaint Received"),
						}, 1),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: []RefDataItem{
				{
					Handle: "Complaint Received",
					Label:  "Complaint Received",
				},
			},
		},
		{
			name:     "Payment sources",
			category: PaymentSourceCategory,
			setup: func() {
				pact.
					AddInteraction().
					UponReceiving("A request for payment source ref data").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String(fmt.Sprintf("/lpa-api/v1/reference-data/%s", PaymentSourceCategory)),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.EachLike(map[string]interface{}{
							"handle":         matchers.String("PHONE"),
							"label":          matchers.String("Paid over the phone"),
							"userSelectable": true,
						}, 1),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: []RefDataItem{
				{
					Handle:         "PHONE",
					Label:          "Paid over the phone",
					UserSelectable: true,
				},
			},
		},
		{
			name:     "Fee reduction types",
			category: FeeReductionTypeCategory,
			setup: func() {
				pact.
					AddInteraction().
					UponReceiving("A request for fee reduction type ref data").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String(fmt.Sprintf("/lpa-api/v1/reference-data/%s", FeeReductionTypeCategory)),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.EachLike(map[string]interface{}{
							"handle": matchers.String("REMISSION"),
							"label":  matchers.String("Remission"),
						}, 1),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: []RefDataItem{
				{
					Handle: "REMISSION",
					Label:  "Remission",
				},
			},
		},
		{
			name:     "Payment reference types",
			category: PaymentReferenceType,
			setup: func() {
				pact.
					AddInteraction().
					UponReceiving("A request for payment reference type ref data").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String(fmt.Sprintf("/lpa-api/v1/reference-data/%s", PaymentReferenceType)),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.EachLike(map[string]interface{}{
							"handle": matchers.String("GOVUK"),
							"label":  matchers.String("GOV.UK Pay"),
						}, 1),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: []RefDataItem{
				{
					Handle: "GOVUK",
					Label:  "GOV.UK Pay",
				},
			},
		},
		{
			name:     "Complainant category",
			category: ComplainantCategory,
			setup: func() {
				pact.
					AddInteraction().
					UponReceiving("A request for complainant category ref data").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String(fmt.Sprintf("/lpa-api/v1/reference-data/%s", ComplainantCategory)),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.EachLike(map[string]interface{}{
							"handle": matchers.String("LPA_DONOR"),
							"label":  matchers.String("LPA Donor"),
						}, 1),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: []RefDataItem{
				{
					Handle: "LPA_DONOR",
					Label:  "LPA Donor",
				},
			},
		},
		{
			name:     "Complaint origin",
			category: ComplaintOrigin,
			setup: func() {
				pact.
					AddInteraction().
					UponReceiving("A request for complaint origin ref data").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String(fmt.Sprintf("/lpa-api/v1/reference-data/%s", ComplaintOrigin)),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.EachLike(map[string]interface{}{
							"handle": matchers.String("PHONE"),
							"label":  matchers.String("Phone call"),
						}, 1),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: []RefDataItem{
				{
					Handle: "PHONE",
					Label:  "Phone call",
				},
			},
		},
		{
			name:     "Compensation type",
			category: CompensationType,
			setup: func() {
				pact.
					AddInteraction().
					UponReceiving("A request for compensation type ref data").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String(fmt.Sprintf("/lpa-api/v1/reference-data/%s", CompensationType)),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.EachLike(map[string]interface{}{
							"handle": matchers.String("COMPENSATORY"),
							"label":  matchers.String("Compensatory"),
						}, 1),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: []RefDataItem{
				{
					Handle: "COMPENSATORY",
					Label:  "Compensatory",
				},
			},
		},
		{
			name:     "Complaint category",
			category: ComplaintCategory,
			setup: func() {
				pact.
					AddInteraction().
					UponReceiving("A request for complaint category ref data").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String(fmt.Sprintf("/lpa-api/v1/reference-data/%s", ComplaintCategory)),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.EachLike(map[string]interface{}{
							"handle": matchers.String("02"),
							"label":  matchers.String("OPG Decisions"),
							"subcategories": matchers.EachLike(map[string]interface{}{
								"handle": matchers.String("18"),
								"label":  matchers.String("Fee Decision"),
							}, 1),
						}, 1),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: []RefDataItem{
				{
					Handle: "02",
					Label:  "OPG Decisions",
					Subcategories: []RefDataItem{
						{
							Handle: "18",
							Label:  "Fee Decision",
						},
					},
				},
			},
		},
		{
			name:     "Country",
			category: CountryCategory,
			setup: func() {
				pact.
					AddInteraction().
					UponReceiving("A request for country ref data").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String(fmt.Sprintf("/lpa-api/v1/reference-data/%s", CountryCategory)),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.EachLike(map[string]interface{}{
							"handle": matchers.String("GB"),
							"label":  matchers.String("Great Britain"),
						}, 1),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: []RefDataItem{
				{
					Handle: "GB",
					Label:  "Great Britain",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				types, err := client.RefDataByCategory(Context{Context: context.Background()}, tc.category)

				assert.Equal(t, tc.expectedResponse, types)
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
