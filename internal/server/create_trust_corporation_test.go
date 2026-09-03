package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
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

func (m *mockCreateTrustCorporationClient) CreateTrustCorporation(ctx sirius.Context, caseId int, trustCorporation sirius.TrustCorporation) error {
	args := m.Called(ctx, caseId, trustCorporation)
	return args.Error(0)
}

func TestGetCreateTrustCorporation(t *testing.T) {
	for _, isReplacementAttorney := range []string{"false", "true"} {
		t.Run("Is Replacement Attorney: "+isReplacementAttorney, func(t *testing.T) {
			expectedData := TrustCorporationData{
				DonorId: 1,
				CaseId:  2,
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

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, expectedData).
				Return(nil)

			r, _ := http.NewRequest(http.MethodGet, "/?id=1&caseId=2&replacement="+isReplacementAttorney, nil)
			w := httptest.NewRecorder()

			err := CreateTrustCorporation(nil, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, template)
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

func TestPostCreateTrustCorporationWhenValidationError(t *testing.T) {
	expectedError := sirius.ValidationError{
		Field: sirius.FieldErrors{"field": {"": "problem"}},
	}

	client := &mockCreateTrustCorporationClient{}
	client.
		On("CreateTrustCorporation", mock.Anything, 2, mock.Anything).
		Return(expectedError)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, TrustCorporationData{
			AppointedAs: "Attorney",
			DonorId:     1,
			CaseId:      2,
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
	tests := []struct {
		name                  string
		isReplacementAttorney bool
		redirect              string
		attorneyType          string
		systemStatus          *bool
	}{
		{
			name:                  "Redirects when attorney",
			isReplacementAttorney: false,
			redirect:              "/create-attorney?id=1&caseId=2&caseType=lpa",
			attorneyType:          "Attorney",
			systemStatus:          shared.BoolPtr(true),
		},
		{
			name:                  "Redirects when replacement attorney",
			isReplacementAttorney: true,
			redirect:              "/create-replacement-attorney?id=1&caseId=2",
			attorneyType:          "Replacement Attorney",
			systemStatus:          shared.BoolPtr(false),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			expectedTrustCorporation := sirius.TrustCorporation{
				Attorney: sirius.Attorney{
					Person: sirius.Person{
						CompanyName: "ACME",
					},
					SystemStatus: tc.systemStatus,
				},
				IsReplacementAttorney:       tc.isReplacementAttorney,
				TrustCorporationAppointedAs: tc.attorneyType,
			}

			client := &mockCreateTrustCorporationClient{}
			client.
				On("CreateTrustCorporation", mock.Anything, 2, expectedTrustCorporation).
				Return(nil)

			form := url.Values{
				"companyName":              {"ACME"},
				"isReplacementAttorney":    {strconv.FormatBool(tc.isReplacementAttorney)},
				"isTrustCorporationActive": {"true"},
				"add-another":              {"true"},
			}

			r, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("create-trust-corporation/?id=1&caseId=2&replacement=%t", tc.isReplacementAttorney), strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := CreateTrustCorporation(client, nil)(w, r)

			assert.Equal(t, RedirectError(tc.redirect), err)
		})
	}
}
