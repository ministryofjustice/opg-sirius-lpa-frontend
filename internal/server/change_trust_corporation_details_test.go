package server

import (
	"net/http"
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

var testChangeTrustCorpDetailsCaseSummary = sirius.CaseSummary{
	DigitalLpa: sirius.DigitalLpa{
		UID: "M-TCTC-TCTC-TCTC",
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
		name          string
		caseUID       string
		trustCorpUID  string
		form          formTrustCorporationDetails
		errorReturned error
	}{
		{
			name:         "Change Active Original Trust Corporation Details",
			caseUID:      "M-TCTC-TCTC-TCTC",
			trustCorpUID: "302b05c7-896c-4290-904e-2005e4f1e81e",
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
			name:         "Change Inactive Replacement Trust Corporation Details",
			caseUID:      "M-TCTC-TCTC-TCTC",
			trustCorpUID: "123a01b1-456d-5391-813d-2010d3e2d72d",
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
			name:         "Template Error Returned",
			caseUID:      "M-TCTC-TCTC-TCTC",
			trustCorpUID: "302b05c7-896c-4290-904e-2005e4f1e81e",
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
				On("CaseSummary", mock.Anything, tc.caseUID).
				Return(testChangeTrustCorpDetailsCaseSummary, nil)
			client.
				On("RefDataByCategory", mock.Anything, sirius.CountryCategory).
				Return([]sirius.RefDataItem{{Handle: "GB", Label: "Great Britain"}}, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything,
					changeTrustCorporationDetailsData{
						Countries: []sirius.RefDataItem{{Handle: "GB", Label: "Great Britain"}},
						CaseUID:   tc.caseUID,
						Form:      tc.form,
					}).
				Return(tc.errorReturned)

			server := newMockServer("/lpa/{uid}/trust-corporation/{trustCorporationUID}/change-details", ChangeTrustCorporationDetails(client, template.Func))

			r, _ := http.NewRequest(http.MethodGet, "/lpa/"+tc.caseUID+"/trust-corporation/"+tc.trustCorpUID+"/change-details", nil)
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
		On("CaseSummary", mock.Anything, "M-TCTC-TCTC-TCTC").
		Return(sirius.CaseSummary{}, errExample)

	assertChangeTrustCorporationDetailsErrors(t, client, "M-TCTC-TCTC-TCTC", errExample)
}

func TestGetChangeTrustCorporationDetailsWhenRefDataByCategoryErrors(t *testing.T) {
	client := &mockChangeTrustCorporationDetailsClient{}
	client.
		On("CaseSummary", mock.Anything, "M-TCTC-TCTC-TCTC").
		Return(testChangeTrustCorpDetailsCaseSummary, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.CountryCategory).
		Return([]sirius.RefDataItem{}, errExample)

	assertChangeTrustCorporationDetailsErrors(t, client, "M-TCTC-TCTC-TCTC", errExample)
}

func assertChangeTrustCorporationDetailsErrors(t *testing.T, client *mockChangeTrustCorporationDetailsClient, uid string, expectedError error) {
	server := newMockServer("/lpa/{uid}/trust-corporation/{trustCorporationUID}/change-details", ChangeTrustCorporationDetails(client, nil))

	r, _ := http.NewRequest(http.MethodGet, "/lpa/"+uid+"/trust-corporation/123a01b1-456d-5391-813d-2010d3e2d72d/change-details", nil)
	_, err := server.serve(r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}
