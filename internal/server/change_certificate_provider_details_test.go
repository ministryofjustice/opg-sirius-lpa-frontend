package server

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockChangeCertificateDetailsClient struct {
	mock.Mock
}

func (m *mockChangeCertificateDetailsClient) CaseSummary(ctx sirius.Context, uid string) (sirius.CaseSummary, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(sirius.CaseSummary), args.Error(1)
}

func (m *mockChangeCertificateDetailsClient) RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error) {
	args := m.Called(ctx, category)
	if args.Get(0) != nil {
		return args.Get(0).([]sirius.RefDataItem), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestGetChangeCertificateProviderDetails(t *testing.T) {
	caseUid := "M-1111-1111-1111"

	caseSummary := sirius.CaseSummary{
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

	client := &mockChangeCertificateDetailsClient{}
	client.
		On("CaseSummary", mock.Anything, caseUid).
		Return(caseSummary, nil)
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

	client := &mockChangeCertificateDetailsClient{}
	client.
		On("CaseSummary", mock.Anything, caseUid).
		Return(caseSummary, expectedError)

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

	assert.Equal(t, expectedError, err)
}

func TestGetChangeCertificateProviderDetailsRefDataByCategoryError(t *testing.T) {
	caseUid := "M-1111-1111-1111"

	caseSummary := sirius.CaseSummary{}

	client := &mockChangeCertificateDetailsClient{}
	client.
		On("CaseSummary", mock.Anything, caseUid).
		Return(caseSummary, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.CountryCategory).
		Return([]sirius.RefDataItem{}, expectedError)

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

	assert.Equal(t, expectedError, err)
}

func TestPostChangeCertificateProviderDetailsRedirectReturned(t *testing.T) {
	caseUid := "M-1111-1111-1111"

	caseSummary := sirius.CaseSummary{}

	client := &mockChangeCertificateDetailsClient{}
	client.
		On("CaseSummary", mock.Anything, caseUid).
		Return(caseSummary, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.CountryCategory).
		Return([]sirius.RefDataItem{}, nil)

	template := &mockTemplate{}

	server := newMockServer(
		"/lpa/{uid}/certificate-provider/change-details",
		ChangeCertificateProviderDetails(client, template.Func),
	)

	req, _ := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("/lpa/%s/certificate-provider/change-details", caseUid),
		nil,
	)
	_, err := server.serve(req)

	assert.Equal(t, RedirectError(fmt.Sprintf("/lpa/%s/lpa-details#certificate-provider", caseUid)), err)
	mock.AssertExpectationsForObjects(t, client, template)
}
