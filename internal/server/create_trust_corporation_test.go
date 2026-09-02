package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCreateTrustCorporationClient struct {
	mock.Mock
}

func (m *mockCreateTrustCorporationClient) Lpa(ctx sirius.Context, id int) (sirius.Lpa, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Lpa), args.Error(1)
}

func (m *mockCreateTrustCorporationClient) CreateTrustCorporation(ctx sirius.Context, caseId int, trustCorporation sirius.TrustCorporation) error {
	args := m.Called(ctx, caseId, trustCorporation)
	return args.Error(0)
}

func (m *mockCreateTrustCorporationClient) UpdateTrustCorporation(ctx sirius.Context, trustCorporationId int, trustCorporation sirius.TrustCorporation) error {
	args := m.Called(ctx, trustCorporationId, trustCorporation)
	return args.Error(0)
}

func TestGetCreateTrustCorporation(t *testing.T) {
	for _, isReplacementAttorney := range []string{"false", "true"} {
		t.Run("Is Replacement Attorney: "+isReplacementAttorney, func(t *testing.T) {
			expectedData := createTrustCorporationData{
				DonorId:  1,
				CaseId:   2,
				CaseType: "lpa",
				TrustCorporation: sirius.TrustCorporation{
					IsReplacementAttorney: isReplacementAttorney == "true",
					Attorney:              sirius.Attorney{SystemStatus: shared.BoolPtr(true)},
				},
				Title:    "Add a trust corporation",
				HtmxPost: "/create-trust-corporation?id=1&caseId=2&replacement=" + isReplacementAttorney,
			}

			if isReplacementAttorney == "true" {
				expectedData.AppointedAs = "Replacement attorney"
			} else {
				expectedData.AppointedAs = "Attorney"
			}

			client := &mockCreateTrustCorporationClient{}
			client.
				On("Lpa", mock.Anything, 2).
				Return(sirius.Lpa{Case: sirius.Case{
					CaseType: "LPA",
				}}, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, expectedData).
				Return(nil)

			r, _ := http.NewRequest(http.MethodGet, "/?id=1&caseId=2&replacement="+isReplacementAttorney, nil)
			w := httptest.NewRecorder()

			err := CreateTrustCorporation(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, template)
		})
	}
}

func TestGetEditTrustCorporation(t *testing.T) {
	for _, isReplacementAttorney := range []string{"false", "true"} {
		t.Run("Is Replacement Attorney: "+isReplacementAttorney, func(t *testing.T) {
			expectedTrustCorporation := sirius.TrustCorporation{
				Attorney:              sirius.Attorney{Person: sirius.Person{ID: 3}},
				IsReplacementAttorney: isReplacementAttorney == "true",
			}
			expectedData := createTrustCorporationData{
				DonorId:          1,
				CaseId:           2,
				TrustCorporation: expectedTrustCorporation,
				Title:            "Update trust corporation details",
				HtmxPost:         "/create-trust-corporation?id=1&caseId=2&trustCorporationId=3&replacement=" + isReplacementAttorney,
				IsEditing:        true,
			}

			if isReplacementAttorney == "true" {
				expectedData.AppointedAs = "Replacement attorney"
			} else {
				expectedData.AppointedAs = "Attorney"
			}

			client := &mockCreateTrustCorporationClient{}
			client.
				On("Lpa", mock.Anything, 2).
				Return(sirius.Lpa{Case: sirius.Case{
					TrustCorporations: []sirius.TrustCorporation{expectedTrustCorporation},
				}}, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, expectedData).
				Return(nil)

			r, _ := http.NewRequest(http.MethodGet, "/?id=1&caseId=2&trustCorporationId=3&replacement="+isReplacementAttorney, nil)
			w := httptest.NewRecorder()

			err := CreateTrustCorporation(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestGetCreateTrustCorporationBadQuery(t *testing.T) {
	testCases := map[string]string{
		"no-id":       "/",
		"bad-id":      "/?id=test",
		"no-case-id":  "/?id=123",
		"bad-case-id": "/?id=123&caseId=test",
	}

	for name, query := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, query, nil)
			w := httptest.NewRecorder()

			err := CreateTrustCorporation(nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestGetCreateTrustCorporationBadTrustCorpQuery(t *testing.T) {
	client := &mockCreateTrustCorporationClient{}
	client.
		On("Lpa", mock.Anything, 123).
		Return(sirius.Lpa{}, nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123&caseId=123&trustCorporationId=test", nil)
	w := httptest.NewRecorder()

	err := CreateTrustCorporation(client, nil)(w, r)

	assert.NotNil(t, err)
}

func TestGetEditTrustCorporationWhenLpaFails(t *testing.T) {
	client := &mockCreateTrustCorporationClient{}
	client.
		On("Lpa", mock.Anything, 2).
		Return(sirius.Lpa{}, errExample)

	r, _ := http.NewRequest(http.MethodGet, "/create-trust-corporation/?id=1&caseId=2&trustCorporationId=3", nil)
	w := httptest.NewRecorder()

	err := CreateTrustCorporation(client, nil)(w, r)
	assert.Equal(t, errExample, err)
}

func TestPostCreateTrustCorporation(t *testing.T) {
	for _, isReplacementAttorney := range []string{"false", "true"} {
		t.Run("Is Replacement Attorney: "+isReplacementAttorney, func(t *testing.T) {
			expectedTrustCorporation := sirius.TrustCorporation{
				Attorney: sirius.Attorney{
					Person: sirius.Person{
						CompanyName:       "ACME",
						Email:             "test@test.com",
						PhoneNumber:       "1234",
						AddressLine1:      "221B",
						AddressLine2:      "Baker Street",
						AddressLine3:      "Marylebone",
						Town:              "London",
						County:            "Greater London",
						Postcode:          "NW1 6XE",
						Country:           "United Kingdom",
						IsAirmailRequired: false,
					},
					CompanyNumber: "123",
				},
				IsReplacementAttorney: isReplacementAttorney == "true",
			}

			if isReplacementAttorney == "true" {
				expectedTrustCorporation.TrustCorporationAppointedAs = "Replacement Attorney"
				expectedTrustCorporation.SystemStatus = shared.BoolPtr(false)
			} else {
				expectedTrustCorporation.TrustCorporationAppointedAs = "Attorney"
				expectedTrustCorporation.SystemStatus = shared.BoolPtr(true)
			}

			client := &mockCreateTrustCorporationClient{}
			client.
				On("CreateTrustCorporation", mock.Anything, 2, expectedTrustCorporation).
				Return(nil)
			client.
				On("Lpa", mock.Anything, 2).
				Return(sirius.Lpa{Case: sirius.Case{CaseType: "LPA", ID: 2}}, nil)

			form := url.Values{
				"companyName":              {"ACME"},
				"companyNumber":            {"123"},
				"email":                    {"test@test.com"},
				"phoneNumber":              {"1234"},
				"addressLine1":             {"221B"},
				"addressLine2":             {"Baker Street"},
				"addressLine3":             {"Marylebone"},
				"town":                     {"London"},
				"county":                   {"Greater London"},
				"postcode":                 {"NW1 6XE"},
				"country":                  {"United Kingdom"},
				"isAirmailRequired":        {"false"},
				"isReplacementAttorney":    {isReplacementAttorney},
				"isTrustCorporationActive": {"true"},
			}

			r, _ := http.NewRequest(http.MethodPost, "create-trust-corporation/?id=1&caseId=2&replacement="+isReplacementAttorney, strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := CreateTrustCorporation(client, nil)(w, r)

			assert.Equal(t, RedirectError("/create-lpa?id=1&caseId=2#scroll-to-attorneys-corporation"), err)
		})
	}
}

func TestPostEditTrustCorporation(t *testing.T) {
	existingTrustCorporation := sirius.TrustCorporation{
		Attorney:              sirius.Attorney{Person: sirius.Person{ID: 3}},
		IsReplacementAttorney: false,
	}
	updatedTrustCorporation := sirius.TrustCorporation{
		Attorney: sirius.Attorney{
			Person: sirius.Person{
				CompanyName:       "ACME",
				Email:             "test@test.com",
				PhoneNumber:       "1234",
				AddressLine1:      "221B",
				AddressLine2:      "Baker Street",
				AddressLine3:      "Marylebone",
				Town:              "London",
				County:            "Greater London",
				Postcode:          "NW1 6XE",
				Country:           "United Kingdom",
				IsAirmailRequired: false,
			},
			SystemStatus:  shared.BoolPtr(true),
			CompanyNumber: "123",
		},
		IsReplacementAttorney:       false,
		TrustCorporationAppointedAs: "Attorney",
	}

	client := &mockCreateTrustCorporationClient{}
	client.
		On("Lpa", mock.Anything, 2).
		Return(sirius.Lpa{Case: sirius.Case{TrustCorporations: []sirius.TrustCorporation{existingTrustCorporation}}}, nil)
	client.
		On("UpdateTrustCorporation", mock.Anything, 3, updatedTrustCorporation).
		Return(nil)

	form := url.Values{
		"companyName":              {"ACME"},
		"companyNumber":            {"123"},
		"email":                    {"test@test.com"},
		"phoneNumber":              {"1234"},
		"addressLine1":             {"221B"},
		"addressLine2":             {"Baker Street"},
		"addressLine3":             {"Marylebone"},
		"town":                     {"London"},
		"county":                   {"Greater London"},
		"postcode":                 {"NW1 6XE"},
		"country":                  {"United Kingdom"},
		"isAirmailRequired":        {"false"},
		"isReplacementAttorney":    {"false"},
		"isTrustCorporationActive": {"true"},
	}

	r, _ := http.NewRequest(http.MethodPost, "create-trust-corporation/?id=1&caseId=2&trustCorporationId=3&replacement=false", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateTrustCorporation(client, nil)(w, r)

	assert.Equal(t, RedirectError("/create-lpa?id=1&caseId=2#scroll-to-attorneys-corporation"), err)
}

func TestPostCreateTrustCorporationWhenCreateFails(t *testing.T) {
	expectedTrustCorporation := sirius.TrustCorporation{
		Attorney: sirius.Attorney{
			Person: sirius.Person{
				CompanyName: "ACME",
			},
			SystemStatus: shared.BoolPtr(true),
		},
		IsReplacementAttorney:       false,
		TrustCorporationAppointedAs: "Attorney",
	}

	client := &mockCreateTrustCorporationClient{}
	client.
		On("Lpa", mock.Anything, 2).
		Return(sirius.Lpa{Case: sirius.Case{CaseType: "LPA", ID: 2}}, nil)
	client.
		On("CreateTrustCorporation", mock.Anything, 2, expectedTrustCorporation).
		Return(errExample)

	form := url.Values{
		"companyName":              {"ACME"},
		"isReplacementAttorney":    {"false"},
		"isTrustCorporationActive": {"true"},
	}

	r, _ := http.NewRequest(http.MethodPost, "create-trust-corporation/?id=1&caseId=2&replacement=false", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateTrustCorporation(client, nil)(w, r)

	assert.Equal(t, errExample, err)
}

func TestPostEditTrustCorporationWhenUpdateFails(t *testing.T) {
	existingTrustCorporation := sirius.TrustCorporation{
		Attorney:              sirius.Attorney{Person: sirius.Person{ID: 3}},
		IsReplacementAttorney: false,
	}
	expectedTrustCorporation := sirius.TrustCorporation{
		Attorney: sirius.Attorney{
			Person: sirius.Person{
				CompanyName: "ACME",
			},
			SystemStatus: shared.BoolPtr(true),
		},
		IsReplacementAttorney:       false,
		TrustCorporationAppointedAs: "Attorney",
	}

	client := &mockCreateTrustCorporationClient{}
	client.
		On("Lpa", mock.Anything, 2).
		Return(sirius.Lpa{Case: sirius.Case{TrustCorporations: []sirius.TrustCorporation{existingTrustCorporation}}}, nil)
	client.
		On("UpdateTrustCorporation", mock.Anything, 3, expectedTrustCorporation).
		Return(errExample)

	form := url.Values{
		"companyName":              {"ACME"},
		"isReplacementAttorney":    {"false"},
		"isTrustCorporationActive": {"true"},
	}

	r, _ := http.NewRequest(http.MethodPost, "create-trust-corporation/?id=1&caseId=2&trustCorporationId=3&replacement=false", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateTrustCorporation(client, nil)(w, r)

	assert.Equal(t, errExample, err)
}

func TestPostCreateTrustCorporationWhenValidationError(t *testing.T) {
	expectedError := sirius.ValidationError{
		Field: sirius.FieldErrors{"field": {"": "problem"}},
	}

	client := &mockCreateTrustCorporationClient{}
	client.
		On("Lpa", mock.Anything, 2).
		Return(sirius.Lpa{Case: sirius.Case{CaseType: "LPA", ID: 2}}, nil)
	client.
		On("CreateTrustCorporation", mock.Anything, 2, mock.Anything).
		Return(expectedError)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createTrustCorporationData{
			AppointedAs: "Attorney",
			DonorId:     1,
			CaseId:      2,
			CaseType:    "lpa",
			TrustCorporation: sirius.TrustCorporation{
				IsReplacementAttorney:       false,
				TrustCorporationAppointedAs: "Attorney",
				Attorney: sirius.Attorney{
					Person: sirius.Person{
						CompanyName: "ACME",
					},
					SystemStatus: shared.BoolPtr(true),
				},
			},
			Title:    "Add a trust corporation",
			HtmxPost: "/create-trust-corporation?id=1&caseId=2&replacement=false",
			Error:    expectedError,
		}).
		Return(nil)

	form := url.Values{
		"companyName":              {"ACME"},
		"isReplacementAttorney":    {"false"},
		"isTrustCorporationActive": {"true"},
	}

	r, _ := http.NewRequest(http.MethodPost, "create-trust-corporation/?id=1&caseId=2&replacement=false", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateTrustCorporation(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateTrustCorporationAddAnotherRedirects(t *testing.T) {
	expectedTrustCorporation := sirius.TrustCorporation{
		Attorney: sirius.Attorney{
			Person: sirius.Person{
				CompanyName: "ACME",
			},
			SystemStatus: shared.BoolPtr(true),
		},
		IsReplacementAttorney:       false,
		TrustCorporationAppointedAs: "Attorney",
	}

	client := &mockCreateTrustCorporationClient{}
	client.
		On("Lpa", mock.Anything, 2).
		Return(sirius.Lpa{Case: sirius.Case{CaseType: "LPA", ID: 2}}, nil)
	client.
		On("CreateTrustCorporation", mock.Anything, 2, expectedTrustCorporation).
		Return(nil)

	form := url.Values{
		"companyName":              {"ACME"},
		"isReplacementAttorney":    {"false"},
		"isTrustCorporationActive": {"true"},
		"add-another":              {"true"},
	}

	r, _ := http.NewRequest(http.MethodPost, "create-trust-corporation/?id=1&caseId=2&replacement=false", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateTrustCorporation(client, nil)(w, r)

	assert.Equal(t, RedirectError("/create-attorney?id=1&caseId=2&caseType=lpa"), err)
}

func TestPostEditTrustCorporationEditNextRedirects(t *testing.T) {
	lpa := sirius.Lpa{Case: sirius.Case{
		CaseType: "LPA",
		TrustCorporations: []sirius.TrustCorporation{
			{
				Attorney:              sirius.Attorney{Person: sirius.Person{ID: 3}},
				IsReplacementAttorney: false,
			},
			{
				Attorney:              sirius.Attorney{Person: sirius.Person{ID: 4}},
				IsReplacementAttorney: false,
			},
		}}}
	expectedTrustCorporation := sirius.TrustCorporation{
		Attorney: sirius.Attorney{
			Person: sirius.Person{
				CompanyName: "ACME",
			},
			SystemStatus: shared.BoolPtr(true),
		},
		IsReplacementAttorney:       false,
		TrustCorporationAppointedAs: "Attorney",
	}

	client := &mockCreateTrustCorporationClient{}
	client.
		On("Lpa", mock.Anything, 2).
		Return(lpa, nil)
	client.
		On("UpdateTrustCorporation", mock.Anything, 3, expectedTrustCorporation).
		Return(nil)

	form := url.Values{
		"companyName":              {"ACME"},
		"isReplacementAttorney":    {"false"},
		"isTrustCorporationActive": {"true"},
		"next-attorney":            {"true"},
	}

	r, _ := http.NewRequest(http.MethodPost, "create-trust-corporation/?id=1&caseId=2&trustCorporationId=3&replacement=false", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateTrustCorporation(client, nil)(w, r)

	assert.Equal(t, RedirectError("/create-attorney?id=1&caseId=2&caseType=lpa&attorneyId=4"), err)
}

func TestGetNextTrustCorporationIdWillReturnNextNumber(t *testing.T) {
	result, _ := GetIdForNextAttorney(2, false, []sirius.TrustCorporation{
		{
			Attorney:              sirius.Attorney{Person: sirius.Person{ID: 1}},
			IsReplacementAttorney: false,
		},
		{
			Attorney:              sirius.Attorney{Person: sirius.Person{ID: 2}},
			IsReplacementAttorney: false,
		},
		{
			Attorney:              sirius.Attorney{Person: sirius.Person{ID: 3}},
			IsReplacementAttorney: false,
		},
	}, []sirius.Attorney{})
	assert.Equal(t, 3, result)
}

func TestGetNextTrustCorporationIdWillReturnNextHigherIdWhenSequenceHasGaps(t *testing.T) {
	result, _ := GetIdForNextAttorney(2, false, []sirius.TrustCorporation{
		{
			Attorney:              sirius.Attorney{Person: sirius.Person{ID: 1}},
			IsReplacementAttorney: false,
		},
		{
			Attorney:              sirius.Attorney{Person: sirius.Person{ID: 5}},
			IsReplacementAttorney: false,
		},
		{
			Attorney:              sirius.Attorney{Person: sirius.Person{ID: 2}},
			IsReplacementAttorney: false,
		},
		{
			Attorney:              sirius.Attorney{Person: sirius.Person{ID: 9}},
			IsReplacementAttorney: false,
		},
	}, []sirius.Attorney{})
	assert.Equal(t, 5, result)
}

func TestGetNextTrustCorporationIdWillReturnZero(t *testing.T) {
	result, _ := GetIdForNextAttorney(3, false, []sirius.TrustCorporation{
		{
			Attorney:              sirius.Attorney{Person: sirius.Person{ID: 1}},
			IsReplacementAttorney: false,
		},
		{
			Attorney:              sirius.Attorney{Person: sirius.Person{ID: 2}},
			IsReplacementAttorney: false,
		},
		{
			Attorney:              sirius.Attorney{Person: sirius.Person{ID: 3}},
			IsReplacementAttorney: false,
		},
	}, []sirius.Attorney{})
	assert.Equal(t, 0, result)
}

func TestGetNextTrustCorporationIdWillReturnNextNumberWithSameAppointedType(t *testing.T) {
	result, _ := GetIdForNextAttorney(2, false, []sirius.TrustCorporation{
		{
			Attorney:              sirius.Attorney{Person: sirius.Person{ID: 3}},
			IsReplacementAttorney: true,
		},
		{
			Attorney:              sirius.Attorney{Person: sirius.Person{ID: 4}},
			IsReplacementAttorney: false,
		},
		{
			Attorney:              sirius.Attorney{Person: sirius.Person{ID: 5}},
			IsReplacementAttorney: true,
		},
	}, []sirius.Attorney{})
	assert.Equal(t, 4, result)
}
