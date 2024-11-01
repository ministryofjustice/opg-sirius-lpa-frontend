package server

import (
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type mockChangeDonorDetailsClient struct {
	mock.Mock
}

func (m *mockChangeDonorDetailsClient) CaseSummary(ctx sirius.Context, uid string) (sirius.CaseSummary, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(sirius.CaseSummary), args.Error(1)
}

func (m *mockChangeDonorDetailsClient) ChangeDonorDetails(ctx sirius.Context, caseUID string, donorDetailsData sirius.ChangeDonorDetails) error {
	return m.Called(ctx, caseUID, donorDetailsData).Error(0)

}

func (m *mockChangeDonorDetailsClient) RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error) {
	args := m.Called(ctx, category)
	if args.Get(0) != nil {
		return args.Get(0).([]sirius.RefDataItem), args.Error(1)
	}
	return nil, args.Error(1)
}

var testCaseSummary = sirius.CaseSummary{
	DigitalLpa: sirius.DigitalLpa{
		UID: "M-AAAA-1111-BBBB",
		SiriusData: sirius.SiriusData{
			ID: 15,
			Application: sirius.Draft{
				DonorFirstNames: "Zackary",
				DonorLastName:   "Lemmonds",
				DonorAddress: sirius.Address{
					Line1:    "9 Mount Pleasant Drive",
					Town:     "East Harling",
					Postcode: "NR16 2GB",
					Country:  "UK",
				},
				PhoneNumber: "1234567890",
			},
		},
		LpaStoreData: sirius.LpaStoreData{
			Donor: sirius.LpaStoreDonor{
				LpaStorePerson: sirius.LpaStorePerson{
					Uid:        "302b05c7-896c-4290-904e-2005e4f1e81e",
					FirstNames: "Zackary",
					LastName:   "Lemmonds",
					Address: sirius.LpaStoreAddress{
						Line1:    "9 Mount Pleasant Drive",
						Town:     "East Harling",
						Postcode: "NR16 2GB",
						Country:  "UK",
					},
				},
				DateOfBirth: "1965-04-18",
			},
			SignedAt: "2024-02-11T15:04:05Z",
		},
	},
}

func TestGetChangeDonorDetails(t *testing.T) {
	client := &mockChangeDonorDetailsClient{}
	client.
		On("CaseSummary", mock.Anything, "M-AAAA-1111-BBBB").
		Return(testCaseSummary, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.CountryCategory).
		Return([]sirius.RefDataItem{{Handle: "GB", Label: "Great Britain"}}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything,
			changeDonorDetailsData{
				Countries: []sirius.RefDataItem{{Handle: "GB", Label: "Great Britain"}},
				CaseUID:   "M-AAAA-1111-BBBB",
				Form: formDonorDetails{
					FirstNames:        "Zackary",
					LastName:          "Lemmonds",
					OtherNamesKnownBy: "",
					DateOfBirth:       dob{Day: 18, Month: 4, Year: 1965},
					Address: sirius.Address{
						Line1:    "9 Mount Pleasant Drive",
						Town:     "East Harling",
						Postcode: "NR16 2GB",
						Country:  "UK",
					},
					PhoneNumber: "1234567890",
					LpaSignedOn: dob{11, 2, 2024},
				},
			}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/change-donor-details/?uid=M-AAAA-1111-BBBB", nil)
	w := httptest.NewRecorder()

	err := ChangeDonorDetails(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetChangeDonorDetailsWhenCaseSummaryErrors(t *testing.T) {

	client := &mockChangeDonorDetailsClient{}
	client.
		On("CaseSummary", mock.Anything, "M-EEEE-EEEE-EEEE").
		Return(sirius.CaseSummary{}, expectedError)

	assertChangeDonorDetailsErrors(t, client, "M-EEEE-EEEE-EEEE", expectedError)
}

func TestGetChangeDonorDetailsWhenRefDataByCategoryErrors(t *testing.T) {
	client := &mockChangeDonorDetailsClient{}
	client.
		On("CaseSummary", mock.Anything, "M-AAAA-1111-BBBB").
		Return(testCaseSummary, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.CountryCategory).
		Return([]sirius.RefDataItem{}, expectedError)

	assertChangeDonorDetailsErrors(t, client, "M-AAAA-1111-BBBB", expectedError)
}

func assertChangeDonorDetailsErrors(t *testing.T, client *mockChangeDonorDetailsClient, uid string, expectedError error) {
	r, _ := http.NewRequest(http.MethodGet, "/change-donor-details/?uid="+uid, nil)
	w := httptest.NewRecorder()

	err := ChangeDonorDetails(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetChangeDonorDetailsWhenTemplateErrors(t *testing.T) {
	client := &mockChangeDonorDetailsClient{}
	client.
		On("CaseSummary", mock.Anything, "M-AAAA-1111-BBBB").
		Return(testCaseSummary, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.CountryCategory).
		Return([]sirius.RefDataItem{{Handle: "GB", Label: "Great Britain"}}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything,
			changeDonorDetailsData{
				Countries: []sirius.RefDataItem{{Handle: "GB", Label: "Great Britain"}},
				CaseUID:   "M-AAAA-1111-BBBB",
				Form: formDonorDetails{
					FirstNames:        "Zackary",
					LastName:          "Lemmonds",
					OtherNamesKnownBy: "",
					DateOfBirth:       dob{Day: 18, Month: 4, Year: 1965},
					Address: sirius.Address{
						Line1:    "9 Mount Pleasant Drive",
						Town:     "East Harling",
						Postcode: "NR16 2GB",
						Country:  "UK",
					},
					PhoneNumber: "1234567890",
					LpaSignedOn: dob{11, 2, 2024},
				},
			}).
		Return(expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/change-donor-details/?uid=M-AAAA-1111-BBBB", nil)
	w := httptest.NewRecorder()

	err := ChangeDonorDetails(client, template.Func)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostChangeDonorDetails(t *testing.T) {
	client := &mockChangeDonorDetailsClient{}
	client.
		On("CaseSummary", mock.Anything, "M-AAAA-1111-BBBB").
		Return(testCaseSummary, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.CountryCategory).
		Return([]sirius.RefDataItem{{Handle: "GB", Label: "Great Britain"}}, nil)
	client.
		On("ChangeDonorDetails", mock.Anything, "M-AAAA-1111-BBBB", sirius.ChangeDonorDetails{
			FirstNames:        "Samuel",
			LastName:          "Smith",
			OtherNamesKnownBy: "Sam",
			DateOfBirth:       "1991-01-01",
			Address: sirius.Address{
				Line1:    "9 Mount",
				Line2:    "Pleasant Drive",
				Line3:    "Norwich",
				Town:     "East Harling",
				Postcode: "NR16 2GB",
				Country:  "GB",
			},
			Phone:       "",
			Email:       "test@test.com",
			LpaSignedOn: "2024-10-09",
		}).
		Return(nil)

	template := &mockTemplate{}

	form := url.Values{
		"firstNames":        {"Samuel"},
		"lastName":          {"Smith"},
		"otherNamesKnownBy": {"Sam"},
		"dob.day":           {"1"},
		"dob.month":         {"1"},
		"dob.year":          {"1991"},
		"address.Line1":     {"9 Mount"},
		"address.Line2":     {"Pleasant Drive"},
		"address.Line3":     {"Norwich"},
		"address.Town":      {"East Harling"},
		"address.Postcode":  {"NR16 2GB"},
		"address.Country":   {"GB"},
		"phoneNumber":       {""},
		"email":             {"test@test.com"},
		"lpaSignedOn.day":   {"9"},
		"lpaSignedOn.month": {"10"},
		"lpaSignedOn.year":  {"2024"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/change-donor-details/?uid=M-AAAA-1111-BBBB", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := ChangeDonorDetails(client, template.Func)(w, r)
	resp := w.Result()

	assert.Equal(t, RedirectError("/lpa/M-AAAA-1111-BBBB/lpa-details"), err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostChangeDonorDetailsWhenAPIFails(t *testing.T) {
	client := &mockChangeDonorDetailsClient{}
	client.
		On("CaseSummary", mock.Anything, "M-AAAA-1111-BBBB").
		Return(testCaseSummary, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.CountryCategory).
		Return([]sirius.RefDataItem{{Handle: "GB", Label: "Great Britain"}}, nil)
	client.
		On("ChangeDonorDetails", mock.Anything, "M-AAAA-1111-BBBB", sirius.ChangeDonorDetails{
			FirstNames:        "Samuel",
			LastName:          "Smith",
			OtherNamesKnownBy: "Sam",
			DateOfBirth:       "1991-01-01",
			Address: sirius.Address{
				Line1:    "9 Mount",
				Line2:    "Pleasant Drive",
				Line3:    "Norwich",
				Town:     "East Harling",
				Postcode: "NR16 2GB",
				Country:  "GB",
			},
			Phone:       "",
			Email:       "test@test.com",
			LpaSignedOn: "2024-10-09",
		}).
		Return(expectedError)

	template := &mockTemplate{}

	form := url.Values{
		"firstNames":        {"Samuel"},
		"lastName":          {"Smith"},
		"otherNamesKnownBy": {"Sam"},
		"dob.day":           {"1"},
		"dob.month":         {"1"},
		"dob.year":          {"1991"},
		"address.Line1":     {"9 Mount"},
		"address.Line2":     {"Pleasant Drive"},
		"address.Line3":     {"Norwich"},
		"address.Town":      {"East Harling"},
		"address.Postcode":  {"NR16 2GB"},
		"address.Country":   {"GB"},
		"phoneNumber":       {""},
		"email":             {"test@test.com"},
		"lpaSignedOn.day":   {"9"},
		"lpaSignedOn.month": {"10"},
		"lpaSignedOn.year":  {"2024"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/change-donor-details/?uid=M-AAAA-1111-BBBB", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := ChangeDonorDetails(client, template.Func)(w, r)
	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostChangeDonorDetailsWhenValidationError(t *testing.T) {
	client := &mockChangeDonorDetailsClient{}
	client.
		On("CaseSummary", mock.Anything, "M-AAAA-1111-BBBB").
		Return(testCaseSummary, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.CountryCategory).
		Return([]sirius.RefDataItem{{Handle: "GB", Label: "Great Britain"}}, nil)
	client.
		On("ChangeDonorDetails", mock.Anything, "M-AAAA-1111-BBBB", sirius.ChangeDonorDetails{
			LastName:    "Smith",
			DateOfBirth: "1965-04-18",
			Address: sirius.Address{
				Line1:    "9 Mount",
				Town:     "East Harling",
				Postcode: "NR16 2GB",
				Country:  "GB",
			},
			Phone:       "1234567890",
			LpaSignedOn: "2024-10-09",
		}).
		Return(sirius.ValidationError{Field: sirius.FieldErrors{
			"firstName": {"required": "Value required and can't be empty"},
		}})

	template := &mockTemplate{}

	template.
		On("Func", mock.Anything,
			mock.MatchedBy(func(data changeDonorDetailsData) bool {
				return data.Error.Field["firstName"]["required"] == "Value required and can't be empty"
			}),
		).
		Return(nil)

	form := url.Values{
		"firstNames":        {""},
		"lastName":          {"Smith"},
		"address.Line1":     {"9 Mount"},
		"address.Town":      {"East Harling"},
		"address.Postcode":  {"NR16 2GB"},
		"address.Country":   {"GB"},
		"lpaSignedOn.day":   {"9"},
		"lpaSignedOn.month": {"10"},
		"lpaSignedOn.year":  {"2024"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/change-donor-details/?uid=M-AAAA-1111-BBBB", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := ChangeDonorDetails(client, template.Func)(w, r)
	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}
