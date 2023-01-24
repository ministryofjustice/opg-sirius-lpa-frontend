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

func TestRefDataByCategoryFeeReductionType(t *testing.T) {
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
					Given("Some fee reduction types exist").
					UponReceiving("A request for fee reduction ref data").
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				types, err := client.RefDataByCategory(Context{Context: context.Background()}, FeeReductionTypeCategory)

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

func TestRefDataByCategoryPaymentReferenceType(t *testing.T) {
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
					Given("Some payment reference types exist").
					UponReceiving("A request for payment reference ref data").
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				types, err := client.RefDataByCategory(Context{Context: context.Background()}, PaymentReferenceType)

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

func TestRefDataByCategoryDocumentTemplateId(t *testing.T) {
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
					Given("Some document template ids types exist").
					UponReceiving("A request for document template id ref data").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String(fmt.Sprintf("/lpa-api/v1/reference-data/%s", DocumentTemplateIdCategory)),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: dsl.EachLike(map[string]interface{}{
							"handle": dsl.String("DDONSCREENSUMMARY"),
							"label":  dsl.String("Donor deceased: Blank template"),
						}, 1),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
			},
			expectedResponse: []RefDataItem{
				{
					Handle: "DDONSCREENSUMMARY",
					Label:  "Donor deceased: Blank template",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				types, err := client.RefDataByCategory(Context{Context: context.Background()}, DocumentTemplateIdCategory)

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

func TestRefDataByCategoryComplainantCategory(t *testing.T) {
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				types, err := client.RefDataByCategory(Context{Context: context.Background()}, ComplainantCategory)

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

func TestRefDataByCategoryComplainantOrigin(t *testing.T) {
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				types, err := client.RefDataByCategory(Context{Context: context.Background()}, ComplaintOrigin)

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
