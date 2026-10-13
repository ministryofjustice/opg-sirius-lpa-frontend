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
	flowCases := []struct {
		name      string
		flowQuery string
	}{
		{name: "Without flow"},
		{name: "Create flow", flowQuery: "&flow=create"},
	}

	for _, flowCase := range flowCases {
		t.Run(flowCase.name, func(t *testing.T) {
			for _, isHtmx := range []bool{false, true} {
				t.Run("Is Htmx: "+strconv.FormatBool(isHtmx), func(t *testing.T) {
					client := &mockCreateAttorneyClient{}
					client.
						On("RefDataByCategory", mock.Anything, sirius.RelationshipToDonorCategory).
						Return(mockRelationshipToDonorCategories, nil)
					client.On("Lpa", mock.Anything, 2).
						Return(sirius.Lpa{Case: sirius.Case{SubType: "pfa"}}, nil)

					expectedData := createAttorneyData{
						FlowQuery:            flowCase.flowQuery,
						IsPartial:            isHtmx,
						DonorId:              1,
						CaseId:               2,
						RelationshipToDonors: mockRelationshipToDonorCategories,
						Attorney:             sirius.Attorney{SystemStatus: shared.BoolPtr(true)},
						Title:                "Add an attorney",
						CaseType:             "lpa",
						CaseSubType:          "pfa",
					}
					template := &mockTemplate{}
					template.
						On("Func", mock.Anything, expectedData).
						Return(nil)

					r, _ := http.NewRequest(http.MethodGet, "/?id=1&caseId=2&caseType=lpa"+flowCase.flowQuery, nil)
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

func TestPostCreateAttorneyTransitions(t *testing.T) {
	tests := []struct {
		name         string
		caseType     string
		caseSubType  string
		formKey      string
		flowQuery    string
		redirect     string
		htmxRedirect string
		htmxSwap     string
	}{
		{
			name:         "EPA save",
			caseType:     "epa",
			redirect:     "/create-epa?id=1&caseId=2#accordion-create-epa-heading-3",
			htmxRedirect: "/create-epa?id=1&caseId=2",
			htmxSwap:     "innerHTML show:#accordion-create-epa-heading-3:top",
		},
		{
			name:         "EPA add another",
			caseType:     "epa",
			formKey:      "add-another",
			redirect:     "/create-attorney?id=1&caseId=2&caseType=epa",
			htmxRedirect: "/create-attorney?id=1&caseId=2&caseType=epa",
			htmxSwap:     "innerHTML scroll:.action-panel__content:top",
		},
		{
			name:         "LPA save without flow",
			caseType:     "lpa",
			caseSubType:  "pfa",
			redirect:     "/create-lpa?id=1&caseId=2#scroll-to-attorneys",
			htmxRedirect: "/create-lpa?id=1&caseId=2",
			htmxSwap:     "innerHTML show:#accordion-create-epa-heading-3:top",
		},
		{
			name:         "LPA save with create flow",
			caseType:     "lpa",
			caseSubType:  "pfa",
			flowQuery:    "&flow=create",
			redirect:     "/create-lpa?id=1&caseId=2&flow=create#scroll-to-attorneys",
			htmxRedirect: "/create-lpa?id=1&caseId=2&flow=create",
			htmxSwap:     "innerHTML show:#accordion-create-epa-heading-3:top",
		},
		{
			name:         "LPA add another without flow",
			caseType:     "lpa",
			caseSubType:  "pfa",
			formKey:      "add-another",
			redirect:     "/create-attorney?id=1&caseId=2&caseType=lpa",
			htmxRedirect: "/create-attorney?id=1&caseId=2&caseType=lpa",
			htmxSwap:     "innerHTML scroll:.action-panel__content:top",
		},
		{
			name:         "LPA add another with create flow",
			caseType:     "lpa",
			caseSubType:  "pfa",
			formKey:      "add-another",
			flowQuery:    "&flow=create",
			redirect:     "/create-attorney?id=1&caseId=2&caseType=lpa&flow=create",
			htmxRedirect: "/create-attorney?id=1&caseId=2&caseType=lpa&flow=create",
			htmxSwap:     "innerHTML scroll:.action-panel__content:top",
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
						SystemStatus: shared.BoolPtr(true),
					}
					if tc.caseType == "epa" {
						attorney.RelationshipToDonor = "no relation"
					}
					client := &mockCreateAttorneyClient{}
					if tc.caseType == "lpa" {
						client.
							On("Lpa", mock.Anything, 2).
							Return(sirius.Lpa{Case: sirius.Case{SubType: tc.caseSubType}}, nil)
					}
					client.
						On("CreateAttorney", mock.Anything, 2, tc.caseType, attorney).
						Return(nil).
						On("RefDataByCategory", mock.Anything, sirius.RelationshipToDonorCategory).
						Return(mockRelationshipToDonorCategories, nil)

					template := &mockTemplate{}

					if isHtmx {
						template.
							On("Func", mock.Anything, createAttorneyData{
								FlowQuery:            tc.flowQuery,
								IsPartial:            true,
								DonorId:              1,
								CaseId:               2,
								RelationshipToDonors: mockRelationshipToDonorCategories,
								Attorney:             attorney,
								HtmxRedirect:         tc.htmxRedirect,
								HtmxSwap:             tc.htmxSwap,
								Title:                "Add an attorney",
								CaseType:             tc.caseType,
								CaseSubType:          tc.caseSubType,
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

					r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2&caseType="+tc.caseType+tc.flowQuery, strings.NewReader(form.Encode()))
					r.Header.Add("Content-Type", formUrlEncoded)
					if isHtmx {
						r.Header.Add("HX-Request", "true")
					}
					w := httptest.NewRecorder()

					err := CreateAttorney(client, template.Func)(w, r)
					resp := w.Result()

					if isHtmx {
						assert.Nil(t, err)
					} else {
						assert.Equal(t, RedirectError(tc.redirect), err)
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
				On("Func", mock.Anything, createAttorneyData{
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
