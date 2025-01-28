package server

import (
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type mockChangeAttorneyDetailsClient struct {
	mock.Mock
}

func (m *mockChangeAttorneyDetailsClient) CaseSummary(ctx sirius.Context, uid string) (sirius.CaseSummary, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(sirius.CaseSummary), args.Error(1)
}

func (m *mockChangeAttorneyDetailsClient) RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error) {
	args := m.Called(ctx, category)
	if args.Get(0) != nil {
		return args.Get(0).([]sirius.RefDataItem), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockChangeAttorneyDetailsClient) ChangeAttorneyDetails(ctx sirius.Context, caseUID string, attorneyUID string, attorneyDetailsData sirius.ChangeAttorneyDetails) error {
	return m.Called(ctx, caseUID, attorneyUID, attorneyDetailsData).Error(0)
}

var testChangeAttorneyDetailsCaseSummary = sirius.CaseSummary{
	DigitalLpa: sirius.DigitalLpa{
		UID: "M-DDDD-DDDD-DDDD",
		LpaStoreData: sirius.LpaStoreData{
			Attorneys: []sirius.LpaStoreAttorney{
				{
					LpaStorePerson: sirius.LpaStorePerson{
						Uid:        "302b05c7-896c-4290-904e-2005e4f1e81e",
						FirstNames: "Jack",
						LastName:   "Black",
						Address: sirius.LpaStoreAddress{
							Line1:    "9 Mount Pleasant Drive",
							Town:     "East Harling",
							Postcode: "NR16 2GB",
							Country:  "UK",
						},
					},
					DateOfBirth:     "1990-02-22",
					Status:          "active",
					AppointmentType: "original",
					Email:           "a@example.com",
					Mobile:          "077577575757",
					SignedAt:        "2024-01-12T10:09:09Z",
				},
				{
					LpaStorePerson: sirius.LpaStorePerson{
						Uid:        "123a01b1-456d-5391-813d-2010d3e2d72d",
						FirstNames: "Jack",
						LastName:   "White",
						Address: sirius.LpaStoreAddress{
							Line1:    "29 Grange Road",
							Town:     "Birmingham",
							Postcode: "B29 6BL",
							Country:  "UK",
						},
					},
					DateOfBirth:     "1990-02-22",
					Status:          "inactive",
					AppointmentType: "replacement",
					Email:           "b@example.com",
					Mobile:          "07122121212",
					SignedAt:        "2024-11-28T19:22:11Z",
				},
			},
		},
	},
}

func TestGetChangeAttorneyDetails(t *testing.T) {
	tests := []struct {
		name           string
		caseUID        string
		attorneyUID    string
		attorneyStatus string
		form           formAttorneyDetails
		errorReturned  error
	}{
		{
			name:           "Change Regular Attorney Details",
			caseUID:        "M-DDDD-DDDD-DDDD",
			attorneyUID:    "302b05c7-896c-4290-904e-2005e4f1e81e",
			attorneyStatus: "active",
			form: formAttorneyDetails{
				FirstNames:  "Jack",
				LastName:    "Black",
				DateOfBirth: dob{Day: 22, Month: 2, Year: 1990},
				Address: sirius.Address{
					Line1:    "9 Mount Pleasant Drive",
					Town:     "East Harling",
					Postcode: "NR16 2GB",
					Country:  "UK",
				},
				Email:       "a@example.com",
				PhoneNumber: "077577575757",
				SignedAt:    dob{Day: 12, Month: 1, Year: 2024},
			},
			errorReturned: nil,
		},
		{
			name:           "Change Replacement Attorney Details",
			caseUID:        "M-DDDD-DDDD-DDDD",
			attorneyUID:    "123a01b1-456d-5391-813d-2010d3e2d72d",
			attorneyStatus: "replacement",
			form: formAttorneyDetails{
				FirstNames:  "Jack",
				LastName:    "White",
				DateOfBirth: dob{Day: 22, Month: 2, Year: 1990},
				Address: sirius.Address{
					Line1:    "29 Grange Road",
					Town:     "Birmingham",
					Postcode: "B29 6BL",
					Country:  "UK",
				},
				Email:       "b@example.com",
				PhoneNumber: "07122121212",
				SignedAt:    dob{Day: 28, Month: 11, Year: 2024},
			},
			errorReturned: nil,
		},
		{
			name:           "Template Error Returned",
			caseUID:        "M-DDDD-DDDD-DDDD",
			attorneyUID:    "302b05c7-896c-4290-904e-2005e4f1e81e",
			attorneyStatus: "active",
			form: formAttorneyDetails{
				FirstNames:  "Jack",
				LastName:    "Black",
				DateOfBirth: dob{Day: 22, Month: 2, Year: 1990},
				Address: sirius.Address{
					Line1:    "9 Mount Pleasant Drive",
					Town:     "East Harling",
					Postcode: "NR16 2GB",
					Country:  "UK",
				},
				Email:       "a@example.com",
				PhoneNumber: "077577575757",
				SignedAt:    dob{Day: 12, Month: 1, Year: 2024},
			},
			errorReturned: expectedError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockChangeAttorneyDetailsClient{}
			client.
				On("CaseSummary", mock.Anything, tc.caseUID).
				Return(testChangeAttorneyDetailsCaseSummary, nil)
			client.
				On("RefDataByCategory", mock.Anything, sirius.CountryCategory).
				Return([]sirius.RefDataItem{{Handle: "GB", Label: "Great Britain"}}, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything,
					changeAttorneyDetailsData{
						Countries:      []sirius.RefDataItem{{Handle: "GB", Label: "Great Britain"}},
						CaseUID:        tc.caseUID,
						Form:           tc.form,
						AttorneyStatus: tc.attorneyStatus,
					}).
				Return(tc.errorReturned)

			server := newMockServer("/lpa/{uid}/attorney/{attorneyUID}/change-details", ChangeAttorneyDetails(client, template.Func))

			r, _ := http.NewRequest(http.MethodGet, "/lpa/"+tc.caseUID+"/attorney/"+tc.attorneyUID+"/change-details", nil)
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

func TestGetChangeAttorneyDetailsWhenCaseSummaryErrors(t *testing.T) {

	client := &mockChangeAttorneyDetailsClient{}
	client.
		On("CaseSummary", mock.Anything, "M-EEEE-EEEE-EEEE").
		Return(sirius.CaseSummary{}, expectedError)

	assertChangeAttorneyDetailsErrors(t, client, "M-EEEE-EEEE-EEEE", expectedError)
}

func TestGetChangeAttorneyDetailsWhenRefDataByCategoryErrors(t *testing.T) {
	client := &mockChangeAttorneyDetailsClient{}
	client.
		On("CaseSummary", mock.Anything, "M-DDDD-DDDD-DDDD").
		Return(testChangeAttorneyDetailsCaseSummary, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.CountryCategory).
		Return([]sirius.RefDataItem{}, expectedError)

	assertChangeAttorneyDetailsErrors(t, client, "M-DDDD-DDDD-DDDD", expectedError)
}

func assertChangeAttorneyDetailsErrors(t *testing.T, client *mockChangeAttorneyDetailsClient, uid string, expectedError error) {
	server := newMockServer("/lpa/{uid}/attorney/{attorneyUID}/change-details", ChangeAttorneyDetails(client, nil))

	r, _ := http.NewRequest(http.MethodGet, "/lpa/"+uid+"/attorney/123a01b1-456d-5391-813d-2010d3e2d72d/change-details", nil)
	_, err := server.serve(r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestPostChangeAttorneyDetails(t *testing.T) {
	tests := []struct {
		name          string
		apiError      error
		expectedError error
	}{
		{
			name:          "Post form successfully submits",
			apiError:      nil,
			expectedError: RedirectError("/lpa/M-DDDD-DDDD-DDDD/lpa-details"),
		},
		{
			name:          "Post form returns an API failure",
			apiError:      expectedError,
			expectedError: expectedError,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockChangeAttorneyDetailsClient{}
			client.
				On("CaseSummary", mock.Anything, "M-DDDD-DDDD-DDDD").
				Return(testChangeAttorneyDetailsCaseSummary, nil)
			client.
				On("RefDataByCategory", mock.Anything, sirius.CountryCategory).
				Return([]sirius.RefDataItem{{Handle: "GB", Label: "Great Britain"}}, nil)
			client.
				On("ChangeAttorneyDetails", mock.Anything, "M-DDDD-DDDD-DDDD", "302b05c7-896c-4290-904e-2005e4f1e81e", sirius.ChangeAttorneyDetails{
					FirstNames:  "Samuel",
					LastName:    "Smith",
					DateOfBirth: "1991-01-01",
					Address: sirius.Address{
						Line1:    "9 Mount",
						Line2:    "Pleasant Drive",
						Line3:    "Norwich",
						Town:     "East Harling",
						Postcode: "NR16 2GB",
						Country:  "GB",
					},
					Phone:    "123456789",
					Email:    "test@test.com",
					SignedAt: "2024-10-09",
				}).
				Return(tc.apiError)

			template := &mockTemplate{}

			server := newMockServer("/lpa/{uid}/attorney/{attorneyUID}/change-details", ChangeAttorneyDetails(client, template.Func))

			form := url.Values{
				"firstNames":       {"Samuel"},
				"lastName":         {"Smith"},
				"dob.day":          {"1"},
				"dob.month":        {"1"},
				"dob.year":         {"1991"},
				"address.Line1":    {"9 Mount"},
				"address.Line2":    {"Pleasant Drive"},
				"address.Line3":    {"Norwich"},
				"address.Town":     {"East Harling"},
				"address.Postcode": {"NR16 2GB"},
				"address.Country":  {"GB"},
				"phoneNumber":      {"123456789"},
				"email":            {"test@test.com"},
				"signedAt.day":     {"9"},
				"signedAt.month":   {"10"},
				"signedAt.year":    {"2024"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/lpa/M-DDDD-DDDD-DDDD/attorney/302b05c7-896c-4290-904e-2005e4f1e81e/change-details", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			_, err := server.serve(r)

			assert.Equal(t, tc.expectedError, err)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostChangeAttorneyDetailsWhenValidationError(t *testing.T) {
	client := &mockChangeAttorneyDetailsClient{}
	client.
		On("CaseSummary", mock.Anything, "M-DDDD-DDDD-DDDD").
		Return(testChangeAttorneyDetailsCaseSummary, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.CountryCategory).
		Return([]sirius.RefDataItem{{Handle: "GB", Label: "Great Britain"}}, nil)
	client.
		On("ChangeAttorneyDetails", mock.Anything, "M-DDDD-DDDD-DDDD", "302b05c7-896c-4290-904e-2005e4f1e81e", sirius.ChangeAttorneyDetails{
			LastName:    "Black",
			DateOfBirth: "1990-02-22",
			Address: sirius.Address{
				Line1:    "9 Mount Pleasant Drive",
				Town:     "East Harling",
				Postcode: "NR16 2GB",
				Country:  "UK",
			},
			Phone:    "077577575757",
			Email:    "a@example.com",
			SignedAt: "2024-01-12",
		}).
		Return(sirius.ValidationError{Field: sirius.FieldErrors{
			"firstNames": {"required": "Value required and can't be empty"},
		}})

	template := &mockTemplate{}

	template.
		On("Func", mock.Anything,
			mock.MatchedBy(func(data changeAttorneyDetailsData) bool {
				return data.Error.Field["firstNames"]["required"] == "Value required and can't be empty"
			}),
		).
		Return(nil)

	server := newMockServer("/lpa/{uid}/attorney/{attorneyUID}/change-details", ChangeAttorneyDetails(client, template.Func))

	form := url.Values{
		"firstNames": {""},
	}

	r, _ := http.NewRequest(http.MethodPost, "/lpa/M-DDDD-DDDD-DDDD/attorney/302b05c7-896c-4290-904e-2005e4f1e81e/change-details", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	_, err := server.serve(r)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}
