package server

import (
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
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
					DateOfBirth: "1990-02-22",
					Status:      "active",
					Email:       "a@example.com",
					Mobile:      "077577575757",
					SignedAt:    "2024-01-12T10:09:09Z",
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
					DateOfBirth: "1990-02-22",
					Status:      "replacement",
					Email:       "b@example.com",
					Mobile:      "07122121212",
					SignedAt:    "2024-11-28T19:22:11Z",
				},
			},
		},
	},
}

func TestGetChangeAttorneyDetails(t *testing.T) {
	client := &mockChangeAttorneyDetailsClient{}
	client.
		On("CaseSummary", mock.Anything, "M-DDDD-DDDD-DDDD").
		Return(testChangeAttorneyDetailsCaseSummary, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.CountryCategory).
		Return([]sirius.RefDataItem{{Handle: "GB", Label: "Great Britain"}}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything,
			changeAttorneyDetailsData{
				Countries: []sirius.RefDataItem{{Handle: "GB", Label: "Great Britain"}},
				CaseUID:   "M-DDDD-DDDD-DDDD",
				Form: formAttorneyDetails{
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
					SignedAt:    dob{12, 1, 2024},
				},
			}).
		Return(nil)

	server := newMockServer("/lpa/{uid}/attorney/{attorneyUID}/change-details", ChangeAttorneyDetails(client, template.Func))

	r, _ := http.NewRequest(http.MethodGet, "/lpa/M-DDDD-DDDD-DDDD/attorney/302b05c7-896c-4290-904e-2005e4f1e81e/change-details", nil)
	_, err := server.serve(r)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetChangeReplacementAttorneyDetails(t *testing.T) {
	client := &mockChangeAttorneyDetailsClient{}
	client.
		On("CaseSummary", mock.Anything, "M-DDDD-DDDD-DDDD").
		Return(testChangeAttorneyDetailsCaseSummary, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.CountryCategory).
		Return([]sirius.RefDataItem{{Handle: "GB", Label: "Great Britain"}}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything,
			changeAttorneyDetailsData{
				Countries: []sirius.RefDataItem{{Handle: "GB", Label: "Great Britain"}},
				CaseUID:   "M-DDDD-DDDD-DDDD",
				Form: formAttorneyDetails{
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
					SignedAt:    dob{28, 11, 2024},
				},
			}).
		Return(nil)

	server := newMockServer("/lpa/{uid}/attorney/{attorneyUID}/change-details", ChangeAttorneyDetails(client, template.Func))

	r, _ := http.NewRequest(http.MethodGet, "/lpa/M-DDDD-DDDD-DDDD/attorney/123a01b1-456d-5391-813d-2010d3e2d72d/change-details", nil)
	_, err := server.serve(r)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
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

func TestGetChangeAttorneyDetailsWhenTemplateErrors(t *testing.T) {
	client := &mockChangeAttorneyDetailsClient{}
	client.
		On("CaseSummary", mock.Anything, "M-DDDD-DDDD-DDDD").
		Return(testChangeAttorneyDetailsCaseSummary, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.CountryCategory).
		Return([]sirius.RefDataItem{{Handle: "GB", Label: "Great Britain"}}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, changeAttorneyDetailsData{
			Countries: []sirius.RefDataItem{{Handle: "GB", Label: "Great Britain"}},
			CaseUID:   "M-DDDD-DDDD-DDDD",
			Form: formAttorneyDetails{
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
				SignedAt:    dob{12, 1, 2024},
			},
		}).
		Return(expectedError)

	server := newMockServer("/lpa/{uid}/attorney/{attorneyUID}/change-details", ChangeAttorneyDetails(client, template.Func))

	r, _ := http.NewRequest(http.MethodGet, "/lpa/M-DDDD-DDDD-DDDD/attorney/302b05c7-896c-4290-904e-2005e4f1e81e/change-details", nil)
	_, err := server.serve(r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client, template)
}
