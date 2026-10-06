package server

import (
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

type mockCreateAttorneyClient struct {
	mock.Mock
}

func (m *mockCreateAttorneyClient) Epa(ctx sirius.Context, id int) (sirius.Epa, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Epa), args.Error(1)
}

func (m *mockCreateAttorneyClient) Lpa(ctx sirius.Context, id int) (sirius.Lpa, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Lpa), args.Error(1)
}

func (m *mockCreateAttorneyClient) CreateAttorney(ctx sirius.Context, caseId int, caseType string, attorney sirius.Attorney) error {
	args := m.Called(ctx, caseId, caseType, attorney)
	return args.Error(0)
}

func (m *mockCreateAttorneyClient) RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error) {
	args := m.Called(ctx, category)
	if args.Get(0) != nil {
		return args.Get(0).([]sirius.RefDataItem), args.Error(1)
	}
	return nil, args.Error(1)
}

var mockRelationshipToDonorCategories = []sirius.RefDataItem{
	{
		Handle: "LPA_DONOR",
		Label:  "LPA Donor",
	},
}

func TestGetCreateAttorney(t *testing.T) {
	for _, isHtmx := range []bool{false, true} {
		t.Run("Is Htmx: "+strconv.FormatBool(isHtmx), func(t *testing.T) {
			client := &mockCreateAttorneyClient{}
			client.
				On("RefDataByCategory", mock.Anything, sirius.RelationshipToDonorCategory).
				Return(mockRelationshipToDonorCategories, nil)
			client.On("Lpa", mock.Anything, 2).
				Return(sirius.Lpa{Case: sirius.Case{SubType: "pfa"}}, nil)

			expectedData := AttorneyData{
				Attorney:             sirius.Attorney{SystemStatus: shared.BoolPtr(true)},
				CaseId:               2,
				CaseType:             "lpa",
				CaseSubType:          "pfa",
				DonorId:              1,
				IsPartial:            isHtmx,
				RelationshipToDonors: mockRelationshipToDonorCategories,
				Title:                "Add an attorney",
			}

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, expectedData).
				Return(nil)

			r, _ := http.NewRequest(http.MethodGet, "/?id=1&caseId=2&caseType=lpa", nil)
			w := httptest.NewRecorder()

			if isHtmx {

				r.Header.Add("HX-Request", "true")
			}

			err := CreateAttorney(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestGetCreateAttorneyBadQueries(t *testing.T) {
	tests := map[string]struct {
		query  string
		client bool
	}{
		"no-id":       {"/", false},
		"bad-id":      {"/?id=test", false},
		"bad-case-id": {"/?id=123&caseId=test", false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			client := &mockCreateAttorneyClient{}
			if tc.client {
				client.
					On("RefDataByCategory", mock.Anything, sirius.RelationshipToDonorCategory).
					Return(mockRelationshipToDonorCategories, nil)
			}

			r, _ := http.NewRequest(http.MethodGet, tc.query, nil)
			w := httptest.NewRecorder()

			err := CreateAttorney(client, nil)(w, r)

			assert.NotNil(t, err)
			if tc.client {
				mock.AssertExpectationsForObjects(t, client)
			}
		})
	}
}

func TestGetCreateAttorneyWhenRefDataErrors(t *testing.T) {
	client := &mockCreateAttorneyClient{}
	client.
		On("RefDataByCategory", mock.Anything, sirius.RelationshipToDonorCategory).
		Return([]sirius.RefDataItem{}, errExample)

	r, _ := http.NewRequest(http.MethodGet, "/?id=1&caseId=2", nil)
	w := httptest.NewRecorder()

	err := CreateAttorney(client, nil)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetCreateAttorneyWhenLpaErrors(t *testing.T) {
	client := &mockCreateAttorneyClient{}
	client.
		On("Lpa", mock.Anything, 2).
		Return(sirius.Lpa{}, errExample)

	r, _ := http.NewRequest(http.MethodGet, "/?id=1&caseId=2&caseType=lpa", nil)
	w := httptest.NewRecorder()

	err := CreateAttorney(client, nil)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestPostCreateAttorneyAddOrNextAnotherEpa(t *testing.T) {
	tests := []struct {
		name         string
		formKey      string
		htmxSwap     string
		redirect     string
		htmxRedirect string
	}{
		{
			name:         "Post create attorney",
			formKey:      "",
			htmxSwap:     "innerHTML show:#accordion-create-epa-heading-3:top",
			redirect:     "/create-epa?id=1&caseId=2#accordion-create-epa-heading-3",
			htmxRedirect: "/create-epa?id=1&caseId=2",
		},
		{
			name:         "Post create attorney - add another",
			formKey:      "add-another",
			htmxSwap:     "innerHTML scroll:.action-panel__content:top",
			redirect:     "/create-attorney?id=1&caseId=2&caseType=epa",
			htmxRedirect: "/create-attorney?id=1&caseId=2&caseType=epa",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for _, isHtmx := range []bool{false, true} {
				t.Run("Is Htmx: "+strconv.FormatBool(isHtmx), func(t *testing.T) {
					dateString := "2022-04-05"
					attorney := sirius.Attorney{
						Person: sirius.Person{
							Salutation:        "Rev",
							Firstname:         "Rudolph",
							Middlenames:       "Modesto",
							Surname:           "Stotesbury",
							DateOfBirth:       sirius.DateString(dateString),
							AddressLine1:      "Rotonda Gerardo 769",
							AddressLine2:      "Appartamento 94",
							AddressLine3:      "Augusto terme",
							Town:              "San Sabazio",
							County:            "Benevento",
							Postcode:          "57797",
							Country:           "Italy",
							IsAirmailRequired: true,
							PhoneNumber:       "079876543345",
							Email:             "rm2@email.test",
						},
						RelationshipToDonor: "no relation",
						SystemStatus:        shared.BoolPtr(true),
					}
					client := &mockCreateAttorneyClient{}
					client.
						On("CreateAttorney", mock.Anything, 2, "epa", attorney).
						Return(nil).
						On("RefDataByCategory", mock.Anything, sirius.RelationshipToDonorCategory).
						Return(mockRelationshipToDonorCategories, nil)

					template := &mockTemplate{}

					if isHtmx {
						template.
							On("Func", mock.Anything, AttorneyData{
								Attorney:             attorney,
								CaseId:               2,
								CaseType:             "epa",
								DonorId:              1,
								HtmxRedirect:         tc.htmxRedirect,
								HtmxSwap:             tc.htmxSwap,
								IsPartial:            true,
								RelationshipToDonors: mockRelationshipToDonorCategories,
								Title:                "Add an attorney",
							}).
							Return(nil)
					}

					form := url.Values{
						"addressLine1":        {"Rotonda Gerardo 769"},
						"addressLine2":        {"Appartamento 94"},
						"addressLine3":        {"Augusto terme"},
						"country":             {"Italy"},
						"county":              {"Benevento"},
						"dob":                 {dateString},
						"email":               {"rm2@email.test"},
						"firstname":           {"Rudolph"},
						"isAirmailRequired":   {"true"},
						"isAttorneyActive":    {"true"},
						"middlenames":         {"Modesto"},
						"phoneNumber":         {"079876543345"},
						"postcode":            {"57797"},
						"relationshipToDonor": {"no relation"},
						"salutation":          {"Rev"},
						"surname":             {"Stotesbury"},
						"town":                {"San Sabazio"},
					}
					if tc.formKey != "" {
						form.Set(tc.formKey, "true")
					}

					r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2&caseType=epa", strings.NewReader(form.Encode()))
					r.Header.Add("Content-Type", formUrlEncoded)
					if isHtmx {
						r.Header.Add("HX-Request", "true")
					}
					w := httptest.NewRecorder()

					err := CreateAttorney(client, template.Func)(w, r)
					resp := w.Result()

					if !isHtmx {
						assert.Equal(t, RedirectError(tc.redirect), err)
					} else {
						assert.Nil(t, err)
					}
					assert.Equal(t, http.StatusOK, resp.StatusCode)
					mock.AssertExpectationsForObjects(t, client, template)
				})
			}
		})
	}
}

func TestPostCreateAttorneyWhenValidationError(t *testing.T) {
	for _, isHtmx := range []bool{false, true} {
		t.Run("Is Htmx: "+strconv.FormatBool(isHtmx), func(t *testing.T) {
			expectedError := sirius.ValidationError{
				Field: sirius.FieldErrors{"field": {"": "problem"}},
			}

			dateString := "2022-04-05"
			attorney := sirius.Attorney{
				Person: sirius.Person{
					Salutation:        "Rev",
					Firstname:         "Rudolph",
					Middlenames:       "Modesto",
					Surname:           "Stotesbury",
					DateOfBirth:       sirius.DateString(dateString),
					AddressLine1:      "Rotonda Gerardo 769",
					AddressLine2:      "Appartamento 94",
					AddressLine3:      "Augusto terme",
					Town:              "San Sabazio",
					County:            "Benevento",
					Postcode:          "57797",
					Country:           "Italy",
					IsAirmailRequired: true,
					PhoneNumber:       "079876543345",
					Email:             "rm2@email.test",
				},
				RelationshipToDonor: "no relation",
				SystemStatus:        shared.BoolPtr(true),
			}

			client := &mockCreateAttorneyClient{}
			client.
				On("CreateAttorney", mock.Anything, 2, "epa", attorney).
				Return(expectedError).
				On("RefDataByCategory", mock.Anything, sirius.RelationshipToDonorCategory).
				Return(mockRelationshipToDonorCategories, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, AttorneyData{
					Attorney:             attorney,
					CaseId:               2,
					CaseType:             "epa",
					DonorId:              1,
					Error:                expectedError,
					IsPartial:            isHtmx,
					RelationshipToDonors: mockRelationshipToDonorCategories,
					Title:                "Add an attorney",
				}).
				Return(nil)

			form := url.Values{
				"addressLine1":        {"Rotonda Gerardo 769"},
				"addressLine2":        {"Appartamento 94"},
				"addressLine3":        {"Augusto terme"},
				"country":             {"Italy"},
				"county":              {"Benevento"},
				"dob":                 {dateString},
				"email":               {"rm2@email.test"},
				"firstname":           {"Rudolph"},
				"isAirmailRequired":   {"true"},
				"isAttorneyActive":    {"true"},
				"middlenames":         {"Modesto"},
				"phoneNumber":         {"079876543345"},
				"postcode":            {"57797"},
				"relationshipToDonor": {"no relation"},
				"salutation":          {"Rev"},
				"surname":             {"Stotesbury"},
				"town":                {"San Sabazio"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2&caseType=epa", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			if isHtmx {
				r.Header.Add("HX-Request", "true")
			}
			w := httptest.NewRecorder()

			err := CreateAttorney(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}
