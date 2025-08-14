package server

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockChangeCertificateProviderDetailsClient struct {
	mock.Mock
}

func (m *mockChangeCertificateProviderDetailsClient) CaseSummary(ctx sirius.Context, uid string) (sirius.CaseSummary, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(sirius.CaseSummary), args.Error(1)
}

func (m *mockChangeCertificateProviderDetailsClient) ChangeCertificateProviderDetails(ctx sirius.Context, caseUID string, certificateProviderDetailsData sirius.ChangeCertificateProviderDetails) error {
	return m.Called(ctx, caseUID, certificateProviderDetailsData).Error(0)
}

func (m *mockChangeCertificateProviderDetailsClient) RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error) {
	args := m.Called(ctx, category)
	if args.Get(0) != nil {
		return args.Get(0).([]sirius.RefDataItem), args.Error(1)
	}
	return nil, args.Error(1)
}

var testChangeCertificateProviderCaseSummary = sirius.CaseSummary{
	DigitalLpa: sirius.DigitalLpa{
		LpaStoreData: sirius.LpaStoreData{
			CertificateProvider: sirius.LpaStoreCertificateProvider{
				LpaStorePerson: sirius.LpaStorePerson{
					Uid:        "f9982b9a-c40c-4aee-85df-06f95c92bd12",
					FirstNames: "Josefa",
					LastName:   "Kihn",
					Address: sirius.LpaStoreAddress{
						Line1:    "37 Davonte Grange",
						Line2:    "Nether Raynor",
						Line3:    "New Fahey",
						Town:     "Surrey",
						Postcode: "KV94 9CD",
						Country:  "GB",
					},
					Email: "Kyra.Schowalter@example.com",
				},
				Phone:                     "01452 927995",
				Channel:                   "online",
				ContactLanguagePreference: "email",
				SignedAt:                  "2024-01-12T10:09:09Z",
			},
		},
	},
}

func TestGetChangeCertificateProviderDetails(t *testing.T) {
	caseUid := "M-1111-1111-1111"

	client := &mockChangeCertificateProviderDetailsClient{}
	client.
		On("CaseSummary", mock.Anything, caseUid).
		Return(testChangeCertificateProviderCaseSummary, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.CountryCategory).
		Return([]sirius.RefDataItem{{Handle: "GB", Label: "Great Britain"}}, nil)

	form := formCertificateProviderDetails{
		FirstNames: "Josefa",
		LastName:   "Kihn",
		Address: sirius.Address{
			Line1:    "37 Davonte Grange",
			Line2:    "Nether Raynor",
			Line3:    "New Fahey",
			Town:     "Surrey",
			Postcode: "KV94 9CD",
			Country:  "GB",
		},
		Channel:  "online",
		Email:    "Kyra.Schowalter@example.com",
		Phone:    "01452 927995",
		SignedAt: dob{12, 1, 2024},
	}

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, changeCertificateProviderDetailsData{
			Countries: []sirius.RefDataItem{{Handle: "GB", Label: "Great Britain"}},
			CaseUid:   caseUid,
			Form:      form,
		}).
		Return(nil)

	server := newMockServer(
		"/lpa/{uid}/certificate-provider/change-details",
		ChangeCertificateProviderDetails(client, template.Func),
	)

	req, _ := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/lpa/%s/certificate-provider/change-details", caseUid),
		nil,
	)
	_, err := server.serve(req)

	assert.Nil(t, err)

	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetChangeCertificateProviderDetailsCaseSummaryError(t *testing.T) {
	caseUid := "M-1111-1111-1111"

	caseSummary := sirius.CaseSummary{}

	client := &mockChangeCertificateProviderDetailsClient{}
	client.
		On("CaseSummary", mock.Anything, caseUid).
		Return(caseSummary, errExample)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, changeCertificateProviderDetailsData{}).
		Return(nil)

	server := newMockServer(
		"/lpa/{uid}/certificate-provider/change-details",
		ChangeCertificateProviderDetails(client, template.Func),
	)

	req, _ := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/lpa/%s/certificate-provider/change-details", caseUid),
		nil,
	)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
}

func TestGetChangeCertificateProviderDetailsRefDataByCategoryError(t *testing.T) {
	caseUid := "M-1111-1111-1111"

	caseSummary := sirius.CaseSummary{}

	client := &mockChangeCertificateProviderDetailsClient{}
	client.
		On("CaseSummary", mock.Anything, caseUid).
		Return(caseSummary, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.CountryCategory).
		Return([]sirius.RefDataItem{}, errExample)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, changeCertificateProviderDetailsData{}).
		Return(nil)

	server := newMockServer(
		"/lpa/{uid}/certificate-provider/change-details",
		ChangeCertificateProviderDetails(client, template.Func),
	)

	req, _ := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/lpa/%s/certificate-provider/change-details", caseUid),
		nil,
	)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
}

func TestPostChangeCertificateProviderDetailsValidationError(t *testing.T) {
	caseUid := "M-1111-1111-1111"

	client := &mockChangeCertificateProviderDetailsClient{}
	client.
		On("CaseSummary", mock.Anything, caseUid).
		Return(sirius.CaseSummary{}, nil)
	client.
		On("ChangeCertificateProviderDetails", mock.Anything, caseUid, sirius.ChangeCertificateProviderDetails{
			FirstNames: "",
			LastName:   "",
			Address:    sirius.Address{},
			Phone:      "",
			Email:      "",
			SignedAt:   "",
		}).
		Return(sirius.ValidationError{Field: sirius.FieldErrors{
			"firstNames":   {"required": "Value required and can't be empty"},
			"lastName":     {"required": "Value required and can't be empty"},
			"addressLine1": {"required": "Enter line 1 of address"},
			"postcode":     {"required": "Must provide postcode for UK addresses"},
		}})
	client.
		On("RefDataByCategory", mock.Anything, sirius.CountryCategory).
		Return([]sirius.RefDataItem{{Handle: "GB", Label: "Great Britain"}}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything,
			mock.MatchedBy(func(data changeCertificateProviderDetailsData) bool {
				return data.Error.Field["firstNames"]["required"] == "Value required and can't be empty" &&
					data.Error.Field["lastName"]["required"] == "Value required and can't be empty" &&
					data.Error.Field["addressLine1"]["required"] == "Enter line 1 of address" &&
					data.Error.Field["postcode"]["required"] == "Must provide postcode for UK addresses"
			}),
		).
		Return(nil)

	form := url.Values{
		"firstNames":       {""},
		"lastName":         {""},
		"address.Line1":    {""},
		"address.Line2":    {""},
		"address.Line3":    {""},
		"address.Town":     {""},
		"address.Postcode": {""},
		"address.Country":  {""},
		"phone":            {""},
		"email":            {""},
		"signedAt.day":     {""},
		"signedAt.month":   {""},
		"signedAt.year":    {""},
	}

	server := newMockServer(
		"/lpa/{uid}/certificate-provider/change-details",
		ChangeCertificateProviderDetails(client, template.Func),
	)

	req, _ := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("/lpa/%s/certificate-provider/change-details", caseUid),
		strings.NewReader(form.Encode()),
	)
	req.Header.Add("Content-Type", formUrlEncoded)
	_, err := server.serve(req)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostChangeCertificateProviderDetailsRedirectReturned(t *testing.T) {
	caseUid := "M-1111-1111-1111"

	client := &mockChangeCertificateProviderDetailsClient{}
	client.
		On("CaseSummary", mock.Anything, caseUid).
		Return(testChangeCertificateProviderCaseSummary, nil)
	client.
		On("ChangeCertificateProviderDetails", mock.Anything, caseUid, sirius.ChangeCertificateProviderDetails{
			FirstNames: "Johathan",
			LastName:   "Hammes",
			Address: sirius.Address{
				Line1:    "4 Edyth Place",
				Line2:    "Leannonthorpe",
				Line3:    "Gwent",
				Town:     "Heidenreichwick",
				Postcode: "HR6 9YN",
				Country:  "GB",
			},
			Phone:    "01697 722252",
			Email:    "Johathan.Hammes@example.com",
			SignedAt: "2025-01-19",
		}).
		Return(nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.CountryCategory).
		Return([]sirius.RefDataItem{}, nil)

	template := &mockTemplate{}

	form := url.Values{
		"firstNames":       {"Johathan"},
		"lastName":         {"Hammes"},
		"address.Line1":    {"4 Edyth Place"},
		"address.Line2":    {"Leannonthorpe"},
		"address.Line3":    {"Gwent"},
		"address.Town":     {"Heidenreichwick"},
		"address.Postcode": {"HR6 9YN"},
		"address.Country":  {"GB"},
		"phone":            {"01697 722252"},
		"email":            {"Johathan.Hammes@example.com"},
		"signedAt.day":     {"19"},
		"signedAt.month":   {"01"},
		"signedAt.year":    {"2025"},
	}

	server := newMockServer(
		"/lpa/{uid}/certificate-provider/change-details",
		ChangeCertificateProviderDetails(client, template.Func),
	)

	req, _ := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("/lpa/%s/certificate-provider/change-details", caseUid),
		strings.NewReader(form.Encode()),
	)
	req.Header.Add("Content-Type", formUrlEncoded)
	_, err := server.serve(req)

	assert.Equal(t, RedirectError(fmt.Sprintf("/lpa/%s/lpa-details#certificate-provider", caseUid)), err)
	mock.AssertExpectationsForObjects(t, client, template)
}

var testChangeCertificateProviderCaseSummaryWithEligibilityConfirmed = sirius.CaseSummary{
	DigitalLpa: sirius.DigitalLpa{
		LpaStoreData: sirius.LpaStoreData{
			CertificateProvider: sirius.LpaStoreCertificateProvider{
				LpaStorePerson: sirius.LpaStorePerson{
					Uid:        "f9982b9a-c40c-4aee-85df-06f95c92bd12",
					FirstNames: "Josefa",
					LastName:   "Smith",
					Address: sirius.LpaStoreAddress{
						Line1:    "37 Davonte Grange",
						Line2:    "Nether Raynor",
						Line3:    "New Fahey",
						Town:     "Surrey",
						Postcode: "KV94 9CD",
						Country:  "GB",
					},
					Email: "Kyra.Schowalter@example.com",
				},
				Phone:                     "01452 927995",
				Channel:                   "online",
				ContactLanguagePreference: "email",
				SignedAt:                  "2024-01-12T10:09:09Z",
			},
			Donor: sirius.LpaStoreDonor{
				LpaStorePerson: sirius.LpaStorePerson{
					Uid:        "donor-uid",
					FirstNames: "John",
					LastName:   "Smith",
				},
			},
			CertificateProviderNotRelatedConfirmedAt: "2024-01-15T10:30:00Z",
		},
	},
}

func TestGetChangeCertificateProviderDetailsWithEligibilityConfirmed(t *testing.T) {
	caseUid := "M-1111-1111-1111"

	client := &mockChangeCertificateProviderDetailsClient{}
	client.
		On("CaseSummary", mock.Anything, caseUid).
		Return(testChangeCertificateProviderCaseSummaryWithEligibilityConfirmed, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.CountryCategory).
		Return([]sirius.RefDataItem{{Handle: "GB", Label: "Great Britain"}}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.MatchedBy(func(data changeCertificateProviderDetailsData) bool {
			return data.CaseUid == caseUid
		})).
		Return(nil)

	server := newMockServer(
		"/lpa/{uid}/certificate-provider/change-details",
		ChangeCertificateProviderDetails(client, template.Func),
	)

	req, _ := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/lpa/%s/certificate-provider/change-details", caseUid),
		nil,
	)
	_, err := server.serve(req)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}
