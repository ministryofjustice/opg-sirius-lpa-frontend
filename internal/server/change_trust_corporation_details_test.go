package server

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockChangeTrustCorporationDetailsClient struct {
	mock.Mock
}

func (m *mockChangeTrustCorporationDetailsClient) CaseSummary(ctx sirius.Context, uid string) (sirius.CaseSummary, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(sirius.CaseSummary), args.Error(1)
}

func (m *mockChangeTrustCorporationDetailsClient) RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error) {
	args := m.Called(ctx, category)
	if args.Get(0) != nil {
		return args.Get(0).([]sirius.RefDataItem), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockChangeTrustCorporationDetailsClient) ChangeTrustCorporationDetails(ctx sirius.Context, caseUID string, trustCorpUID string, trustCorpDetailsData sirius.ChangeTrustCorporationDetails) error {
	return m.Called(ctx, caseUID, trustCorpUID, trustCorpDetailsData).Error(0)
}

const caseUID = "M-TCTC-TCTC-TCTC"

var testChangeTrustCorpDetailsCaseSummary = sirius.CaseSummary{
	DigitalLpa: sirius.DigitalLpa{
		UID: caseUID,
		LpaStoreData: sirius.LpaStoreData{
			TrustCorporations: []sirius.LpaStoreTrustCorporation{
				{
					Name:          "Trust Me Once Ltd.",
					CompanyNumber: "123456789",
					Signatories: []sirius.Signatory{
						{
							FirstNames: "First",
						},
					},
					LpaStoreAttorney: sirius.LpaStoreAttorney{
						Decisions:       true,
						Status:          shared.ActiveAttorneyStatus.String(),
						AppointmentType: shared.OriginalAppointmentType.String(),
						Mobile:          "077577575757",
						Email:           "trust.me.once@does.not.exist",
						LpaStorePerson: sirius.LpaStorePerson{
							Uid: "302b05c7-896c-4290-904e-2005e4f1e81e",
							Address: sirius.LpaStoreAddress{
								Line1:    "9 Mount Pleasant Drive",
								Town:     "East Harling",
								Postcode: "NR16 2GB",
								Country:  "UK",
							},
						},
					},
				},
				{
					Name:          "Trust Me Twice",
					CompanyNumber: "987654321",
					Signatories: []sirius.Signatory{
						{
							FirstNames: "Second",
						},
					},
					LpaStoreAttorney: sirius.LpaStoreAttorney{
						Decisions:       false,
						Status:          shared.InactiveAttorneyStatus.String(),
						AppointmentType: shared.ReplacementAppointmentType.String(),
						Mobile:          "07122121212",
						Email:           "trust.me.twice@does.not.exist",
						LpaStorePerson: sirius.LpaStorePerson{
							Uid: "123a01b1-456d-5391-813d-2010d3e2d72d",
							Address: sirius.LpaStoreAddress{
								Line1:    "29 Grange Road",
								Town:     "Birmingham",
								Postcode: "B29 6BL",
								Country:  "UK",
							},
						},
					},
				},
				{
					Name:          "Dont Trust Me",
					CompanyNumber: "987654321",
					LpaStoreAttorney: sirius.LpaStoreAttorney{
						Decisions:       false,
						Status:          shared.RemovedAttorneyStatus.String(),
						AppointmentType: shared.OriginalAppointmentType.String(),
						Email:           "dont.trust.me@does.not.exist",
						LpaStorePerson: sirius.LpaStorePerson{
							Uid: "638f049f-c01f-4ab2-973a-2ea763b3cf7a",
						},
					},
				},
			},
		},
	},
}

func TestGetChangeTrustCorpDetails(t *testing.T) {
	tests := []struct {
		name            string
		trustCorpUID    string
		status          string
		appointmentType string
		form            formTrustCorporationDetails
		errorReturned   error
	}{
		{
			name:            "Change Active Original Trust Corporation Details",
			trustCorpUID:    "302b05c7-896c-4290-904e-2005e4f1e81e",
			status:          shared.ActiveAttorneyStatus.String(),
			appointmentType: shared.OriginalAppointmentType.String(),
			form: formTrustCorporationDetails{
				Name: "Trust Me Once Ltd.",
				Address: sirius.Address{
					Line1:    "9 Mount Pleasant Drive",
					Town:     "East Harling",
					Postcode: "NR16 2GB",
					Country:  "UK",
				},
				Email:         "trust.me.once@does.not.exist",
				PhoneNumber:   "077577575757",
				CompanyNumber: "123456789",
			},
			errorReturned: nil,
		},
		{
			name:            "Change Inactive Replacement Trust Corporation Details",
			trustCorpUID:    "123a01b1-456d-5391-813d-2010d3e2d72d",
			status:          shared.InactiveAttorneyStatus.String(),
			appointmentType: shared.ReplacementAppointmentType.String(),
			form: formTrustCorporationDetails{
				Name: "Trust Me Twice",
				Address: sirius.Address{
					Line1:    "29 Grange Road",
					Town:     "Birmingham",
					Postcode: "B29 6BL",
					Country:  "UK",
				},
				Email:         "trust.me.twice@does.not.exist",
				PhoneNumber:   "07122121212",
				CompanyNumber: "987654321",
			},
			errorReturned: nil,
		},
		{
			name:            "Template Error Returned",
			trustCorpUID:    "302b05c7-896c-4290-904e-2005e4f1e81e",
			status:          shared.ActiveAttorneyStatus.String(),
			appointmentType: shared.OriginalAppointmentType.String(),
			form: formTrustCorporationDetails{
				Name: "Trust Me Once Ltd.",
				Address: sirius.Address{
					Line1:    "9 Mount Pleasant Drive",
					Town:     "East Harling",
					Postcode: "NR16 2GB",
					Country:  "UK",
				},
				Email:         "trust.me.once@does.not.exist",
				PhoneNumber:   "077577575757",
				CompanyNumber: "123456789",
			},
			errorReturned: errExample,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockChangeTrustCorporationDetailsClient{}
			client.
				On("CaseSummary", mock.Anything, caseUID).
				Return(testChangeTrustCorpDetailsCaseSummary, nil)
			client.
				On("RefDataByCategory", mock.Anything, sirius.CountryCategory).
				Return([]sirius.RefDataItem{{Handle: "GB", Label: "Great Britain"}}, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything,
					changeTrustCorporationDetailsData{
						Countries:       []sirius.RefDataItem{{Handle: "GB", Label: "Great Britain"}},
						CaseUID:         caseUID,
						Status:          tc.status,
						AppointmentType: tc.appointmentType,
						Form:            tc.form,
					}).
				Return(tc.errorReturned)

			server := newMockServer("/lpa/{uid}/trust-corporation/{trustCorporationUID}/change-details", ChangeTrustCorporationDetails(client, template.Func))

			r, _ := http.NewRequest(http.MethodGet, "/lpa/"+caseUID+"/trust-corporation/"+tc.trustCorpUID+"/change-details", nil)
			_, err := server.serve(r)

			if tc.errorReturned != nil {
				assert.Equal(t, tc.errorReturned, err)
			} else {
				assert.Nil(t, err)
			}

			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestGetChangeTrustCorpDetailsWhenCaseSummaryErrors(t *testing.T) {

	client := &mockChangeTrustCorporationDetailsClient{}
	client.
		On("CaseSummary", mock.Anything, caseUID).
		Return(sirius.CaseSummary{}, errExample)

	assertChangeTrustCorporationDetailsErrors(t, client, caseUID, errExample)
}

func TestGetChangeTrustCorporationDetailsWhenRefDataByCategoryErrors(t *testing.T) {
	client := &mockChangeTrustCorporationDetailsClient{}
	client.
		On("CaseSummary", mock.Anything, caseUID).
		Return(testChangeTrustCorpDetailsCaseSummary, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.CountryCategory).
		Return([]sirius.RefDataItem{}, errExample)

	assertChangeTrustCorporationDetailsErrors(t, client, caseUID, errExample)
}

func assertChangeTrustCorporationDetailsErrors(t *testing.T, client *mockChangeTrustCorporationDetailsClient, uid string, expectedError error) {
	server := newMockServer("/lpa/{uid}/trust-corporation/{trustCorporationUID}/change-details", ChangeTrustCorporationDetails(client, nil))

	r, _ := http.NewRequest(http.MethodGet, "/lpa/"+uid+"/trust-corporation/123a01b1-456d-5391-813d-2010d3e2d72d/change-details", nil)
	_, err := server.serve(r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestPostChangeTrustCorporationDetails(t *testing.T) {
	tests := []struct {
		name          string
		apiError      error
		expectedError error
	}{
		{
			name:          "Post form successfully submits",
			apiError:      nil,
			expectedError: RedirectError("/lpa/M-TCTC-TCTC-TCTC/lpa-details"),
		},
		{
			name:          "Post form returns an API failure",
			apiError:      errExample,
			expectedError: errExample,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockChangeTrustCorporationDetailsClient{}
			client.
				On("CaseSummary", mock.Anything, caseUID).
				Return(testChangeTrustCorpDetailsCaseSummary, nil)
			client.
				On("RefDataByCategory", mock.Anything, sirius.CountryCategory).
				Return([]sirius.RefDataItem{{Handle: "GB", Label: "Great Britain"}}, nil)
			client.
				On("ChangeTrustCorporationDetails", mock.Anything, caseUID, "302b05c7-896c-4290-904e-2005e4f1e81e", sirius.ChangeTrustCorporationDetails{
					Name: "Trust Ltd.",
					Address: sirius.Address{
						Line1:    "9 Mount",
						Line2:    "Pleasant Drive",
						Line3:    "Norwich",
						Town:     "East Harling",
						Postcode: "NR16 2GB",
						Country:  "GB",
					},
					Phone:         "123456789",
					Email:         "test@test.com",
					CompanyNumber: "20241009",
				}).
				Return(tc.apiError)

			template := &mockTemplate{}

			server := newMockServer("/lpa/{uid}/trust-corporation/{trustCorporationUID}/change-details", ChangeTrustCorporationDetails(client, template.Func))

			form := url.Values{
				"name":             {"Trust Ltd."},
				"address.Line1":    {"9 Mount"},
				"address.Line2":    {"Pleasant Drive"},
				"address.Line3":    {"Norwich"},
				"address.Town":     {"East Harling"},
				"address.Postcode": {"NR16 2GB"},
				"address.Country":  {"GB"},
				"phoneNumber":      {"123456789"},
				"email":            {"test@test.com"},
				"companyNumber":    {"20241009"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/lpa/M-TCTC-TCTC-TCTC/trust-corporation/302b05c7-896c-4290-904e-2005e4f1e81e/change-details", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			_, err := server.serve(r)

			assert.Equal(t, tc.expectedError, err)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostChangeTrustCorporationDetailsWhenValidationError(t *testing.T) {
	client := &mockChangeTrustCorporationDetailsClient{}
	client.
		On("CaseSummary", mock.Anything, caseUID).
		Return(testChangeTrustCorpDetailsCaseSummary, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.CountryCategory).
		Return([]sirius.RefDataItem{{Handle: "GB", Label: "Great Britain"}}, nil)
	client.
		On("ChangeTrustCorporationDetails", mock.Anything, caseUID, "302b05c7-896c-4290-904e-2005e4f1e81e", sirius.ChangeTrustCorporationDetails{
			Address: sirius.Address{
				Line1:    "9 Mount Pleasant Drive",
				Town:     "East Harling",
				Postcode: "NR16 2GB",
				Country:  "UK",
			},
			Phone:         "077577575757",
			Email:         "trust.me.once@does.not.exist",
			CompanyNumber: "123456789",
		}).
		Return(sirius.ValidationError{Field: sirius.FieldErrors{
			"name": {"required": "Value required and can't be empty"},
		}})

	template := &mockTemplate{}

	template.
		On("Func", mock.Anything,
			mock.MatchedBy(func(data changeTrustCorporationDetailsData) bool {
				return data.Error.Field["name"]["required"] == "Value required and can't be empty"
			}),
		).
		Return(nil)

	server := newMockServer("/lpa/{uid}/trust-corporation/{trustCorporationUID}/change-details", ChangeTrustCorporationDetails(client, template.Func))

	form := url.Values{
		"name": {""},
	}

	r, _ := http.NewRequest(http.MethodPost, "/lpa/M-TCTC-TCTC-TCTC/trust-corporation/302b05c7-896c-4290-904e-2005e4f1e81e/change-details", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	_, err := server.serve(r)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}
