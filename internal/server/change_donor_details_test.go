package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func (m *mockChangeDonorDetailsClient) ProgressIndicatorsForDigitalLpa(ctx sirius.Context, uid string) ([]sirius.ProgressIndicator, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).([]sirius.ProgressIndicator), args.Error(1)
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

var newSignedOn = dob{9, 10, 2024}
var newSignedOnTime = newSignedOn.toTime()

func TestGetChangeDonorDetails(t *testing.T) {
	client := &mockChangeDonorDetailsClient{}
	client.
		On("CaseSummary", mock.Anything, "M-AAAA-1111-BBBB").
		Return(testCaseSummary, nil)
	client.
		On("ProgressIndicatorsForDigitalLpa", mock.Anything, "M-AAAA-1111-BBBB").
		Return([]sirius.ProgressIndicator{{
			Indicator: "DONOR_ID",
			Status:    "IN_PROGRESS",
		}}, nil)
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
					PhoneNumber:               "1234567890",
					LpaSignedOn:               dob{11, 2, 2024},
					AuthorisedSignatory:       "",
					SignedByWitnessOne:        "No",
					SignedByWitnessTwo:        "No",
					IndependentWitnessName:    "",
					IndependentWitnessAddress: sirius.Address{},
				},
				DonorIdentityCheckComplete: false,
				DonorDobString:             "1965-04-18",
				SignedByWitnessTwoLabel:    "Signed by witness 2",
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
		Return(sirius.CaseSummary{}, errExample)

	assertChangeDonorDetailsErrors(t, client, "M-EEEE-EEEE-EEEE", errExample)
}

func TestGetChangeDonorDetailsWhenProgressIndicatorsForDigitalLpaErrors(t *testing.T) {
	client := &mockChangeDonorDetailsClient{}
	client.
		On("CaseSummary", mock.Anything, "M-AAAA-1111-BBBB").
		Return(testCaseSummary, nil)
	client.
		On("ProgressIndicatorsForDigitalLpa", mock.Anything, "M-AAAA-1111-BBBB").
		Return([]sirius.ProgressIndicator{}, errExample)

	assertChangeDonorDetailsErrors(t, client, "M-AAAA-1111-BBBB", errExample)
}

func TestGetChangeDonorDetailsWhenRefDataByCategoryErrors(t *testing.T) {
	client := &mockChangeDonorDetailsClient{}
	client.
		On("CaseSummary", mock.Anything, "M-AAAA-1111-BBBB").
		Return(testCaseSummary, nil)
	client.
		On("ProgressIndicatorsForDigitalLpa", mock.Anything, "M-AAAA-1111-BBBB").
		Return([]sirius.ProgressIndicator{{
			Indicator: "DONOR_ID",
			Status:    "IN_PROGRESS",
		}}, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.CountryCategory).
		Return([]sirius.RefDataItem{}, errExample)

	assertChangeDonorDetailsErrors(t, client, "M-AAAA-1111-BBBB", errExample)
}

func assertChangeDonorDetailsErrors(t *testing.T, client *mockChangeDonorDetailsClient, uid string, expectedError error) {
	r, _ := http.NewRequest(http.MethodGet, "/change-donor-details/?uid="+uid, nil)
	w := httptest.NewRecorder()

	err := ChangeDonorDetails(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetChangeDonorDetailsWithDonorIdentityCheck(t *testing.T) {
	client := &mockChangeDonorDetailsClient{}
	client.
		On("CaseSummary", mock.Anything, "M-AAAA-1111-BBBB").
		Return(testCaseSummary, nil)
	client.
		On("ProgressIndicatorsForDigitalLpa", mock.Anything, "M-AAAA-1111-BBBB").
		Return([]sirius.ProgressIndicator{{
			Indicator: "DONOR_ID",
			Status:    "COMPLETE",
		}}, nil)
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
					PhoneNumber:               "1234567890",
					LpaSignedOn:               dob{11, 2, 2024},
					AuthorisedSignatory:       "",
					SignedByWitnessOne:        "No",
					SignedByWitnessTwo:        "No",
					IndependentWitnessName:    "",
					IndependentWitnessAddress: sirius.Address{},
				},
				DonorIdentityCheckComplete: true,
				DonorDobString:             "1965-04-18",
				SignedByWitnessTwoLabel:    "Signed by witness 2",
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

func TestGetChangeDonorDetailsWhenTemplateErrors(t *testing.T) {
	client := &mockChangeDonorDetailsClient{}
	client.
		On("CaseSummary", mock.Anything, "M-AAAA-1111-BBBB").
		Return(testCaseSummary, nil)
	client.
		On("ProgressIndicatorsForDigitalLpa", mock.Anything, "M-AAAA-1111-BBBB").
		Return([]sirius.ProgressIndicator{{
			Indicator: "DONOR_ID",
			Status:    "IN_PROGRESS",
		}}, nil)
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
					PhoneNumber:               "1234567890",
					LpaSignedOn:               dob{11, 2, 2024},
					AuthorisedSignatory:       "",
					SignedByWitnessOne:        "No",
					SignedByWitnessTwo:        "No",
					IndependentWitnessName:    "",
					IndependentWitnessAddress: sirius.Address{},
				},
				DonorIdentityCheckComplete: false,
				DonorDobString:             "1965-04-18",
				SignedByWitnessTwoLabel:    "Signed by witness 2",
			}).
		Return(errExample)

	r, _ := http.NewRequest(http.MethodGet, "/change-donor-details/?uid=M-AAAA-1111-BBBB", nil)
	w := httptest.NewRecorder()

	err := ChangeDonorDetails(client, template.Func)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostChangeDonorDetails(t *testing.T) {
	client := &mockChangeDonorDetailsClient{}
	client.
		On("CaseSummary", mock.Anything, "M-AAAA-1111-BBBB").
		Return(testCaseSummary, nil)
	client.
		On("ProgressIndicatorsForDigitalLpa", mock.Anything, "M-AAAA-1111-BBBB").
		Return([]sirius.ProgressIndicator{{
			Indicator: "DONOR_ID",
			Status:    "IN_PROGRESS",
		}}, nil)
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
			Phone:                            "",
			Email:                            "test@test.com",
			LpaSignedOn:                      newSignedOn.toDateString(),
			AuthorisedSignatory:              "Lucy Mueller",
			WitnessedByCertificateProviderAt: &newSignedOnTime,
			WitnessedByIndependentWitnessAt:  &newSignedOnTime,
			IndependentWitnessName:           "Ora Reagan",
			IndependentWitnessAddress: sirius.Address{
				Line1:    "6 Poplar Close",
				Line2:    "Swift",
				Line3:    "Schneider",
				Town:     "Durham",
				Postcode: "CC9 1GF",
				Country:  "GB",
			},
		}).
		Return(nil)

	template := &mockTemplate{}

	form := url.Values{
		"firstNames":                         {"Samuel"},
		"lastName":                           {"Smith"},
		"otherNamesKnownBy":                  {"Sam"},
		"dob.day":                            {"1"},
		"dob.month":                          {"1"},
		"dob.year":                           {"1991"},
		"address.Line1":                      {"9 Mount"},
		"address.Line2":                      {"Pleasant Drive"},
		"address.Line3":                      {"Norwich"},
		"address.Town":                       {"East Harling"},
		"address.Postcode":                   {"NR16 2GB"},
		"address.Country":                    {"GB"},
		"phoneNumber":                        {""},
		"email":                              {"test@test.com"},
		"lpaSignedOn.day":                    {"9"},
		"lpaSignedOn.month":                  {"10"},
		"lpaSignedOn.year":                   {"2024"},
		"authorisedSignatory":                {"Lucy Mueller"},
		"signedByWitnessOne":                 {"Yes"},
		"signedByWitnessTwo":                 {"Yes"},
		"independentWitnessAddress.Line1":    {"6 Poplar Close"},
		"independentWitnessAddress.Line2":    {"Swift"},
		"independentWitnessAddress.Line3":    {"Schneider"},
		"independentWitnessAddress.Town":     {"Durham"},
		"independentWitnessAddress.Postcode": {"CC9 1GF"},
		"independentWitnessAddress.Country":  {"GB"},
		"independentWitnessName":             {"Ora Reagan"},
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
		On("ProgressIndicatorsForDigitalLpa", mock.Anything, "M-AAAA-1111-BBBB").
		Return([]sirius.ProgressIndicator{{
			Indicator: "DONOR_ID",
			Status:    "IN_PROGRESS",
		}}, nil)
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
			Phone:                            "",
			Email:                            "test@test.com",
			LpaSignedOn:                      newSignedOn.toDateString(),
			AuthorisedSignatory:              "Lucy Mueller",
			WitnessedByCertificateProviderAt: &newSignedOnTime,
			WitnessedByIndependentWitnessAt:  &newSignedOnTime,
			IndependentWitnessName:           "Ora Reagan",
			IndependentWitnessAddress: sirius.Address{
				Line1:    "6 Poplar Close",
				Line2:    "Swift",
				Line3:    "Schneider",
				Town:     "Durham",
				Postcode: "CC9 1GF",
				Country:  "GB",
			},
		}).
		Return(errExample)

	template := &mockTemplate{}

	form := url.Values{
		"firstNames":                         {"Samuel"},
		"lastName":                           {"Smith"},
		"otherNamesKnownBy":                  {"Sam"},
		"dob.day":                            {"1"},
		"dob.month":                          {"1"},
		"dob.year":                           {"1991"},
		"address.Line1":                      {"9 Mount"},
		"address.Line2":                      {"Pleasant Drive"},
		"address.Line3":                      {"Norwich"},
		"address.Town":                       {"East Harling"},
		"address.Postcode":                   {"NR16 2GB"},
		"address.Country":                    {"GB"},
		"phoneNumber":                        {""},
		"email":                              {"test@test.com"},
		"lpaSignedOn.day":                    {"9"},
		"lpaSignedOn.month":                  {"10"},
		"lpaSignedOn.year":                   {"2024"},
		"authorisedSignatory":                {"Lucy Mueller"},
		"signedByWitnessOne":                 {"Yes"},
		"signedByWitnessTwo":                 {"Yes"},
		"independentWitnessAddress.Line1":    {"6 Poplar Close"},
		"independentWitnessAddress.Line2":    {"Swift"},
		"independentWitnessAddress.Line3":    {"Schneider"},
		"independentWitnessAddress.Town":     {"Durham"},
		"independentWitnessAddress.Postcode": {"CC9 1GF"},
		"independentWitnessAddress.Country":  {"GB"},
		"independentWitnessName":             {"Ora Reagan"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/change-donor-details/?uid=M-AAAA-1111-BBBB", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := ChangeDonorDetails(client, template.Func)(w, r)
	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostChangeDonorDetailsWhenValidationError(t *testing.T) {
	client := &mockChangeDonorDetailsClient{}
	client.
		On("CaseSummary", mock.Anything, "M-AAAA-1111-BBBB").
		Return(testCaseSummary, nil)
	client.
		On("ProgressIndicatorsForDigitalLpa", mock.Anything, "M-AAAA-1111-BBBB").
		Return([]sirius.ProgressIndicator{{
			Indicator: "DONOR_ID",
			Status:    "IN_PROGRESS",
		}}, nil)
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

func TestParseDateTime(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected dob
		wantErr  bool
	}{
		{"valid RFC3339", "2024-04-18T09:10:11Z", dob{Day: 18, Month: 4, Year: 2024}, false},
		{"empty string", "", dob{}, false},
		{"null date", "0001-01-01T00:00:00Z", dob{}, false},
		{"invalid format", "18/04/2024", dob{}, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := parseDateTime(tc.input)

			if tc.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expected, actual)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
