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

type mockUpdateTrustCorporationClient struct {
	mock.Mock
}

func (m *mockUpdateTrustCorporationClient) Lpa(ctx sirius.Context, id int) (sirius.Lpa, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Lpa), args.Error(1)
}

func (m *mockUpdateTrustCorporationClient) UpdateTrustCorporation(ctx sirius.Context, caseId int, trustCorporation sirius.TrustCorporation) error {
	args := m.Called(ctx, caseId, trustCorporation)
	return args.Error(0)
}

func TestGetEditTrustCorporation(t *testing.T) {
	for _, isReplacementAttorney := range []string{"false", "true"} {
		t.Run("Is Replacement Attorney: "+isReplacementAttorney, func(t *testing.T) {
			expectedTrustCorporation := sirius.TrustCorporation{
				Attorney:              sirius.Attorney{Person: sirius.Person{ID: 3}},
				IsReplacementAttorney: isReplacementAttorney == "true",
			}
			expectedData := TrustCorporationData{
				DonorId:          1,
				CaseId:           2,
				TrustCorporation: expectedTrustCorporation,
				Title:            "Update trust corporation details",
				HtmxPost:         "/update-trust-corporation?id=1&caseId=2&trustCorporationId=3&replacement=" + isReplacementAttorney,
				IsEditing:        true,
			}

			if isReplacementAttorney == "true" {
				expectedData.AppointedAs = "Replacement attorney"
			} else {
				expectedData.AppointedAs = "Attorney"
			}

			client := &mockUpdateTrustCorporationClient{}
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

			err := UpdateTrustCorporation(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestUpdateTrustCorporationWhenLpaFails(t *testing.T) {
	client := &mockUpdateTrustCorporationClient{}
	client.
		On("Lpa", mock.Anything, 2).
		Return(sirius.Lpa{}, errExample)

	r, _ := http.NewRequest(http.MethodGet, "/update-trust-corporation/?id=1&caseId=2&trustCorporationId=3", nil)
	w := httptest.NewRecorder()

	err := UpdateTrustCorporation(client, nil)(w, r)
	assert.Equal(t, errExample, err)
}

func TestPostUpdateTrustCorporation(t *testing.T) {
	existingTrustCorporation := sirius.TrustCorporation{
		Attorney:              sirius.Attorney{Person: sirius.Person{ID: 3}},
		IsReplacementAttorney: false,
	}
	updatedTrustCorporation := sirius.TrustCorporation{
		Attorney: sirius.Attorney{
			Person: sirius.Person{
				ID:                3,
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

	client := &mockUpdateTrustCorporationClient{}
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

	r, _ := http.NewRequest(http.MethodPost, "update-trust-corporation/?id=1&caseId=2&trustCorporationId=3&replacement=false", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := UpdateTrustCorporation(client, nil)(w, r)

	assert.Equal(t, RedirectError("/create-lpa?id=1&caseId=2#scroll-to-attorneys-corporation"), err)
}

func TestPostUpdateTrustCorporationFails(t *testing.T) {
	existingTrustCorporation := sirius.TrustCorporation{
		Attorney:              sirius.Attorney{Person: sirius.Person{ID: 3}},
		IsReplacementAttorney: false,
	}
	expectedTrustCorporation := sirius.TrustCorporation{
		Attorney: sirius.Attorney{
			Person: sirius.Person{
				ID:          3,
				CompanyName: "ACME",
			},
			SystemStatus: shared.BoolPtr(true),
		},
		IsReplacementAttorney:       false,
		TrustCorporationAppointedAs: "Attorney",
	}

	client := &mockUpdateTrustCorporationClient{}
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

	r, _ := http.NewRequest(http.MethodPost, "update-trust-corporation/?id=1&caseId=2&trustCorporationId=3&replacement=false", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := UpdateTrustCorporation(client, nil)(w, r)

	assert.Equal(t, errExample, err)
}

func TestPostUpdateTrustCorporationUpdateNextRedirects(t *testing.T) {
	lpa := sirius.Lpa{Case: sirius.Case{TrustCorporations: []sirius.TrustCorporation{
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
				ID:          3,
				CompanyName: "ACME",
			},
			SystemStatus: shared.BoolPtr(true),
		},
		IsReplacementAttorney:       false,
		TrustCorporationAppointedAs: "Attorney",
	}

	client := &mockUpdateTrustCorporationClient{}
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
		"editNextTrustCorporation": {"true"},
	}

	r, _ := http.NewRequest(http.MethodPost, "update-trust-corporation/?id=1&caseId=2&trustCorporationId=3&replacement=false", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := UpdateTrustCorporation(client, nil)(w, r)

	assert.Equal(t, RedirectError("/update-trust-corporation?id=1&caseId=2&trustCorporationId=4&replacement=false"), err)
}

func TestGetNextTrustCorporationIdWillReturnNextNumber(t *testing.T) {
	result := GetNextTrustCorporationId(2, false, []sirius.TrustCorporation{
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
	})
	assert.Equal(t, 3, result)
}

func TestGetNextTrustCorporationIdWillReturnNextHigherIdWhenSequenceHasGaps(t *testing.T) {
	result := GetNextTrustCorporationId(2, false, []sirius.TrustCorporation{
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
	})
	assert.Equal(t, 5, result)
}

func TestGetNextTrustCorporationIdWillReturnZero(t *testing.T) {
	result := GetNextTrustCorporationId(3, false, []sirius.TrustCorporation{
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
	})
	assert.Equal(t, 0, result)
}

func TestGetNextTrustCorporationIdWillReturnNextNumberWithSameAppointedType(t *testing.T) {
	result := GetNextTrustCorporationId(2, false, []sirius.TrustCorporation{
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
	})
	assert.Equal(t, 4, result)
}

func TestGetEditTrustCorporationBadQuery(t *testing.T) {
	testCases := map[string]struct {
		query        string
		expectedCode int
	}{
		"no-id":                    {"/update-trust-corporation/", http.StatusNotFound},
		"bad-id":                   {"/update-trust-corporation/?id=test", http.StatusBadRequest},
		"no-case-id":               {"/update-trust-corporation/?id=123", http.StatusNotFound},
		"bad-case-id":              {"/update-trust-corporation/?id=123&caseId=test", http.StatusBadRequest},
		"bad-trust-corporation-id": {"/update-trust-corporation/?id=123&caseId=123&trustCorporationId=test", http.StatusBadRequest},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, tc.query, nil)
			w := httptest.NewRecorder()

			err := UpdateTrustCorporation(nil, nil)(w, r)

			// We expect an error originating from TrustCorporation (strToIntOrStatusError).
			if assert.Error(t, err) {
				if se, ok := err.(sirius.StatusError); ok {
					assert.Equal(t, tc.expectedCode, se.Code)
				}
			}
		})
	}
}
