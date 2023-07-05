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
		{
			name:     "Payment sources",
			category: PaymentSourceCategory,
			setup: func() {
				pact.
					AddInteraction().
					UponReceiving("A request for payment source ref data").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String(fmt.Sprintf("/lpa-api/v1/reference-data/%s", PaymentSourceCategory)),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: dsl.EachLike(map[string]interface{}{
							"handle":         dsl.String("PHONE"),
							"label":          dsl.String("Paid over the phone"),
							"userSelectable": true,
						}, 1),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
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
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String(fmt.Sprintf("/lpa-api/v1/reference-data/%s", FeeReductionTypeCategory)),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: dsl.EachLike(map[string]interface{}{
							"handle": dsl.String("REMISSION"),
							"label":  dsl.String("Remission"),
						}, 1),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
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
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String(fmt.Sprintf("/lpa-api/v1/reference-data/%s", PaymentReferenceType)),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: dsl.EachLike(map[string]interface{}{
							"handle": dsl.String("GOVUK"),
							"label":  dsl.String("GOV.UK Pay"),
						}, 1),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
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
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String(fmt.Sprintf("/lpa-api/v1/reference-data/%s", ComplainantCategory)),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: dsl.EachLike(map[string]interface{}{
							"handle": dsl.String("LPA_DONOR"),
							"label":  dsl.String("LPA Donor"),
						}, 1),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
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
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String(fmt.Sprintf("/lpa-api/v1/reference-data/%s", ComplaintOrigin)),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: dsl.EachLike(map[string]interface{}{
							"handle": dsl.String("PHONE"),
							"label":  dsl.String("Phone call"),
						}, 1),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
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
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String(fmt.Sprintf("/lpa-api/v1/reference-data/%s", CompensationType)),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: dsl.EachLike(map[string]interface{}{
							"handle": dsl.String("COMPENSATORY"),
							"label":  dsl.String("Compensatory"),
						}, 1),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
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
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String(fmt.Sprintf("/lpa-api/v1/reference-data/%s", ComplaintCategory)),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: dsl.EachLike(map[string]interface{}{
							"handle": dsl.String("02"),
							"label":  dsl.String("OPG Decisions"),
							"subcategories": dsl.EachLike(map[string]interface{}{
								"handle": dsl.String("18"),
								"label":  dsl.String("Fee Decision"),
							}, 1),
						}, 1),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
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
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String(fmt.Sprintf("/lpa-api/v1/reference-data/%s", CountryCategory)),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: dsl.EachLike(map[string]interface{}{
							"handle": dsl.String("GB"),
							"label":  dsl.String("Great Britain"),
						}, 1),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
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

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				types, err := client.RefDataByCategory(Context{Context: context.Background()}, tc.category)

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
