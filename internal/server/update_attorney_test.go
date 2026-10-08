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

type mockUpdateAttorneyClient struct {
	mock.Mock
}

func (m *mockUpdateAttorneyClient) Epa(ctx sirius.Context, id int) (sirius.Epa, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Epa), args.Error(1)
}

func (m *mockUpdateAttorneyClient) Lpa(ctx sirius.Context, id int) (sirius.Lpa, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Lpa), args.Error(1)
}

func (m *mockUpdateAttorneyClient) RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error) {
	args := m.Called(ctx, category)
	if args.Get(0) != nil {
		return args.Get(0).([]sirius.RefDataItem), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockUpdateAttorneyClient) UpdateAttorney(ctx sirius.Context, attorneyId int, attorney sirius.Attorney) error {
	args := m.Called(ctx, attorneyId, attorney)
	return args.Error(0)
}

func TestGetEditAttorney(t *testing.T) {
	for _, caseType := range []string{"lpa", "epa"} {
		t.Run("Case Type: "+caseType, func(t *testing.T) {
			existingAttorney := sirius.Attorney{
				Person: sirius.Person{
					ID:        4,
					Firstname: "Rudolph",
					Surname:   "Stotesbury",
				},
				RelationshipToDonor: "no relation",
				SystemStatus:        shared.BoolPtr(true),
			}

			client := &mockUpdateAttorneyClient{}
			client.
				On("RefDataByCategory", mock.Anything, sirius.RelationshipToDonorCategory).
				Return(mockRelationshipToDonorCategories, nil)

			if caseType == "lpa" {
				client.On("Lpa", mock.Anything, 2).
					Return(sirius.Lpa{Case: sirius.Case{Attorneys: []sirius.Attorney{existingAttorney}}}, nil)
			} else {
				client.On("Epa", mock.Anything, 2).
					Return(sirius.Epa{Case: sirius.Case{Attorneys: []sirius.Attorney{existingAttorney}}}, nil)
			}

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, createAttorneyData{
					Attorney:             existingAttorney,
					CaseId:               2,
					CaseType:             caseType,
					DonorId:              1,
					IsEditing:            true,
					RelationshipToDonors: mockRelationshipToDonorCategories,
					Title:                "Update attorney details",
				}).
				Return(nil)

			r, _ := http.NewRequest(http.MethodGet, "/?id=1&caseId=2&attorneyId=4&caseType="+caseType, nil)
			w := httptest.NewRecorder()

			err := UpdateAttorney(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestGetUpdateAttorneyBadQueries(t *testing.T) {
	tests := map[string]struct {
		query  string
		client bool
	}{
		"bad-attorney-id-editing": {"/?id=1&caseId=2&attorneyId=bad", true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			client := &mockUpdateAttorneyClient{}
			if tc.client {
				client.
					On("RefDataByCategory", mock.Anything, sirius.RelationshipToDonorCategory).
					Return(mockRelationshipToDonorCategories, nil)
			}

			r, _ := http.NewRequest(http.MethodGet, tc.query, nil)
			w := httptest.NewRecorder()

			err := UpdateAttorney(client, nil)(w, r)

			assert.NotNil(t, err)
			if tc.client {
				mock.AssertExpectationsForObjects(t, client)
			}
		})
	}
}

func TestGetUpdateAttorneyWhenRefDataErrors(t *testing.T) {
	client := &mockUpdateAttorneyClient{}
	client.
		On("RefDataByCategory", mock.Anything, sirius.RelationshipToDonorCategory).
		Return([]sirius.RefDataItem{}, errExample)

	r, _ := http.NewRequest(http.MethodGet, "/?id=1&caseId=2", nil)
	w := httptest.NewRecorder()

	err := UpdateAttorney(client, nil)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetUpdateAttorneyWhenLpaErrors(t *testing.T) {
	client := &mockUpdateAttorneyClient{}
	client.
		On("Lpa", mock.Anything, 2).
		Return(sirius.Lpa{}, errExample)

	r, _ := http.NewRequest(http.MethodGet, "/?id=1&caseId=2&caseType=lpa", nil)
	w := httptest.NewRecorder()

	err := UpdateAttorney(client, nil)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetUpdateAttorneyEpaErrors(t *testing.T) {
	client := &mockUpdateAttorneyClient{}
	client.
		On("RefDataByCategory", mock.Anything, sirius.RelationshipToDonorCategory).
		Return(mockRelationshipToDonorCategories, nil).
		On("Epa", mock.Anything, 2).
		Return(sirius.Epa{}, errExample)

	r, _ := http.NewRequest(http.MethodGet, "/?id=1&caseId=2&caseType=epa&attorneyId=3", nil)
	w := httptest.NewRecorder()

	err := UpdateAttorney(client, nil)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestPostEditAttorney(t *testing.T) {
	for _, isHtmx := range []bool{false, true} {
		t.Run("Is Htmx: "+strconv.FormatBool(isHtmx), func(t *testing.T) {
			dateString := "2022-04-05"
			existingAttorney := sirius.Attorney{Person: sirius.Person{ID: 4}}
			updatedAttorney := sirius.Attorney{
				Person: sirius.Person{
					AddressLine1:      "Rotonda Gerardo 769",
					AddressLine2:      "Appartamento 94",
					AddressLine3:      "Augusto terme",
					Country:           "Italy",
					County:            "Benevento",
					DateOfBirth:       sirius.DateString(dateString),
					Email:             "rm2@email.test",
					Firstname:         "Rudolph",
					IsAirmailRequired: true,
					Middlenames:       "Modesto",
					PhoneNumber:       "079876543345",
					Postcode:          "57797",
					Salutation:        "Rev",
					Surname:           "Stotesbury",
					Town:              "San Sabazio",
				},
				RelationshipToDonor: "no relation",
				SystemStatus:        shared.BoolPtr(true),
			}

			client := &mockUpdateAttorneyClient{}
			client.
				On("RefDataByCategory", mock.Anything, sirius.RelationshipToDonorCategory).
				Return(mockRelationshipToDonorCategories, nil).
				On("Epa", mock.Anything, 2).
				Return(sirius.Epa{Case: sirius.Case{Attorneys: []sirius.Attorney{existingAttorney}}}, nil).
				On("UpdateAttorney", mock.Anything, 4, updatedAttorney).
				Return(nil)

			template := &mockTemplate{}

			if isHtmx {
				template.
					On("Func", mock.Anything, createAttorneyData{
						Attorney:             updatedAttorney,
						CaseId:               2,
						CaseType:             "epa",
						DonorId:              1,
						HtmxRedirect:         "/create-epa?id=1&caseId=2",
						HtmxSwap:             "innerHTML show:#accordion-create-epa-heading-3:top",
						IsEditing:            true,
						IsPartial:            true,
						RelationshipToDonors: mockRelationshipToDonorCategories,
						Title:                "Update attorney details",
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

			r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2&attorneyId=4&caseType=epa", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			if isHtmx {
				r.Header.Add("HX-Request", "true")
			}
			w := httptest.NewRecorder()

			err := UpdateAttorney(client, template.Func)(w, r)
			resp := w.Result()

			if !isHtmx {
				expectedError := RedirectError("/create-epa?id=1&caseId=2#accordion-create-epa-heading-3")
				assert.Equal(t, err, expectedError)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostEditAttorneyReturnLpa(t *testing.T) {
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
					existingAttorney := sirius.Attorney{Person: sirius.Person{ID: 4}}
					updatedAttorney := sirius.Attorney{
						Person: sirius.Person{
							Firstname: "Rudolph",
							Surname:   "Stotesbury",
						},
						SystemStatus: shared.BoolPtr(true),
					}

					client := &mockUpdateAttorneyClient{}
					client.
						On("Lpa", mock.Anything, 2).
						Return(sirius.Lpa{
							Case: sirius.Case{
								SubType:   "pfa",
								Attorneys: []sirius.Attorney{existingAttorney},
							},
						}, nil).
						On("UpdateAttorney", mock.Anything, 4, updatedAttorney).
						Return(nil).
						On("RefDataByCategory", mock.Anything, sirius.RelationshipToDonorCategory).
						Return(mockRelationshipToDonorCategories, nil)

					template := &mockTemplate{}
					if isHtmx {
						template.
							On("Func", mock.Anything, createAttorneyData{
								FlowQuery:            flowCase.flowQuery,
								IsPartial:            true,
								DonorId:              1,
								CaseId:               2,
								RelationshipToDonors: mockRelationshipToDonorCategories,
								Attorney:             updatedAttorney,
								IsEditing:            true,
								Title:                "Update attorney details",
								HtmxRedirect:         "/create-lpa?id=1&caseId=2" + flowCase.flowQuery,
								HtmxSwap:             "innerHTML show:#accordion-create-epa-heading-3:top",
								CaseType:             "lpa",
								CaseSubType:          "pfa",
							}).
							Return(nil)
					}

					form := url.Values{
						"firstname":        {"Rudolph"},
						"surname":          {"Stotesbury"},
						"isAttorneyActive": {"true"},
					}

					r, _ := http.NewRequest(http.MethodPost, "/update-attorney?id=1&caseId=2&attorneyId=4&caseType=lpa"+flowCase.flowQuery, strings.NewReader(form.Encode()))
					r.Header.Add("Content-Type", formUrlEncoded)
					if isHtmx {
						r.Header.Add("HX-Request", "true")
					}
					w := httptest.NewRecorder()

					err := UpdateAttorney(client, template.Func)(w, r)
					resp := w.Result()

					if isHtmx {
						assert.Nil(t, err)
					} else {
						expectedRedirect := RedirectError("/create-lpa?id=1&caseId=2" + flowCase.flowQuery + "#scroll-to-attorneys")
						assert.Equal(t, expectedRedirect, err)
					}
					assert.Equal(t, http.StatusOK, resp.StatusCode)
					mock.AssertExpectationsForObjects(t, client, template)
				})
			}
		})
	}
}

func TestPostUpdateLpaAttorneyNextActor(t *testing.T) {
	flowCases := []struct {
		name      string
		flowQuery string
	}{
		{name: "Without flow"},
		{name: "Create flow", flowQuery: "&flow=create"},
	}
	for _, flowCase := range flowCases {
		t.Run(flowCase.name, func(t *testing.T) {
			for _, nextActor := range []struct {
				name     string
				redirect string
			}{
				{name: "Attorney", redirect: "/update-attorney?id=1&caseId=2&caseType=lpa&attorneyId=3"},
				{name: "Trust Corporation", redirect: "/create-trust-corporation?id=1&caseId=2&trustCorporationId=3&replacement=false"},
			} {
				t.Run(nextActor.name, func(t *testing.T) {
					for _, isHtmx := range []bool{false, true} {
						t.Run("Is Htmx: "+strconv.FormatBool(isHtmx), func(t *testing.T) {
							dateString := "2022-04-05"
							trustCorp := sirius.TrustCorporation{
								Attorney: sirius.Attorney{
									Person: sirius.Person{
										ID:               3,
										CompanyName:      "ACME",
										CompanyReference: "testing",
										DateOfBirth:      sirius.DateString(dateString),
										AddressLine1:     "Rotonda Gerardo 769",
										Town:             "San Sabazio",
										County:           "Benevento",
										Postcode:         "57797",
										PersonType:       "Trust Corporation",
									},
									SystemStatus: shared.BoolPtr(true),
								},
							}
							lpa := sirius.Lpa{
								Case: sirius.Case{
									Attorneys:         []sirius.Attorney{{Person: sirius.Person{ID: 2}, SystemStatus: shared.BoolPtr(true)}},
									TrustCorporations: []sirius.TrustCorporation{trustCorp},
								},
							}
							if nextActor.name == "Attorney" {
								lpa.TrustCorporations = nil
								lpa.Attorneys = append(lpa.Attorneys, sirius.Attorney{Person: sirius.Person{ID: 3}, SystemStatus: shared.BoolPtr(true)})
							}

							updatedAttorney := sirius.Attorney{
								Person: sirius.Person{
									Salutation:   "Rev",
									Firstname:    "Rudolph",
									Surname:      "Stotesbury",
									DateOfBirth:  sirius.DateString(dateString),
									AddressLine1: "Rotonda Gerardo 769",
									County:       "Benevento",
									Postcode:     "57797",
									Country:      "Italy",
									Email:        "hello@example.com",
									PhoneNumber:  "0123456789",
								},
								SystemStatus: shared.BoolPtr(true),
							}

							client := &mockUpdateAttorneyClient{}
							client.
								On("UpdateAttorney", mock.Anything, 2, updatedAttorney).
								Return(nil).
								On("RefDataByCategory", mock.Anything, sirius.RelationshipToDonorCategory).
								Return(mockRelationshipToDonorCategories, nil)

							client.
								On("Lpa", mock.Anything, 2).
								Return(lpa, nil)

							template := &mockTemplate{}

							if isHtmx {
								template.
									On("Func", mock.Anything, createAttorneyData{
										FlowQuery:            flowCase.flowQuery,
										IsPartial:            true,
										IsEditing:            true,
										DonorId:              1,
										CaseId:               2,
										RelationshipToDonors: mockRelationshipToDonorCategories,
										Attorney:             updatedAttorney,
										Title:                "Update attorney details",
										HtmxRedirect:         nextActor.redirect + flowCase.flowQuery,
										HtmxSwap:             "innerHTML scroll:.action-panel__content:top",
										CaseType:             "lpa",
										NextAttorneyId:       3,
									}).
									Return(nil)
							}

							form := url.Values{
								"salutation":           {"Rev"},
								"firstname":            {"Rudolph"},
								"surname":              {"Stotesbury"},
								"dob":                  {dateString},
								"addressLine1":         {"Rotonda Gerardo 769"},
								"county":               {"Benevento"},
								"postcode":             {"57797"},
								"country":              {"Italy"},
								"email":                {"hello@example.com"},
								"phoneNumber":          {"0123456789"},
								"isAttorneyActive":     {"true"},
								"update-next-attorney": {"true"},
							}

							r, _ := http.NewRequest(http.MethodPost, "/update-attorney?id=1&caseId=2&caseType=lpa&attorneyId=2"+flowCase.flowQuery, strings.NewReader(form.Encode()))
							r.Header.Add("Content-Type", formUrlEncoded)
							if isHtmx {
								r.Header.Add("HX-Request", "true")
							}
							w := httptest.NewRecorder()

							err := UpdateAttorney(client, template.Func)(w, r)
							resp := w.Result()

							if !isHtmx {
								expectedRedirect := RedirectError(nextActor.redirect + flowCase.flowQuery)
								assert.Equal(t, err, expectedRedirect)
							} else {
								assert.Nil(t, err)
							}
							assert.Equal(t, http.StatusOK, resp.StatusCode)
							mock.AssertExpectationsForObjects(t, client, template)
						})
					}
				})
			}
		})
	}
}
