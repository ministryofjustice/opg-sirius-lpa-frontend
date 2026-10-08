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

func TestPostUpdateAttorneyNextAnotherLpa(t *testing.T) {
	for _, isHtmx := range []bool{false, true} {
		t.Run("Is Htmx: "+strconv.FormatBool(isHtmx), func(t *testing.T) {
			dateString := "2022-04-05"
			trustCorp := sirius.TrustCorporation{
				Attorney: sirius.Attorney{
					Person: sirius.Person{
						AddressLine1:     "Rotonda Gerardo 769",
						CompanyName:      "ACME",
						CompanyReference: "testing",
						County:           "Benevento",
						DateOfBirth:      sirius.DateString(dateString),
						ID:               3,
						PersonType:       "Trust Corporation",
						Postcode:         "57797",
						Town:             "San Sabazio",
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
			updatedAttorney := sirius.Attorney{
				Person: sirius.Person{
					AddressLine1: "Rotonda Gerardo 769",
					Country:      "Italy",
					County:       "Benevento",
					DateOfBirth:  sirius.DateString(dateString),
					Email:        "hello@example.com",
					Firstname:    "Rudolph",
					PhoneNumber:  "0123456789",
					Postcode:     "57797",
					Salutation:   "Rev",
					Surname:      "Stotesbury",
				},
				SystemStatus: shared.BoolPtr(true),
			}

			client := &mockUpdateAttorneyClient{}

			client.
				On("Lpa", mock.Anything, 2).
				Return(lpa, nil)
			client.
				On("RefDataByCategory", mock.Anything, sirius.RelationshipToDonorCategory).
				Return(mockRelationshipToDonorCategories, nil)
			client.
				On("UpdateAttorney", mock.Anything, 2, updatedAttorney).
				Return(nil)

			template := &mockTemplate{}

			if isHtmx {
				template.
					On("Func", mock.Anything, createAttorneyData{
						Attorney:             updatedAttorney,
						CaseId:               2,
						CaseType:             "lpa",
						DonorId:              1,
						HtmxRedirect:         "/create-trust-corporation?id=1&caseId=2&trustCorporationId=3&replacement=false",
						HtmxSwap:             "innerHTML scroll:.action-panel__content:top",
						IsEditing:            true,
						IsPartial:            true,
						NextAttorneyId:       3,
						RelationshipToDonors: mockRelationshipToDonorCategories,
						Title:                "Update attorney details",
					}).
					Return(nil)
			}

			form := url.Values{
				"addressLine1":         {"Rotonda Gerardo 769"},
				"country":              {"Italy"},
				"county":               {"Benevento"},
				"dob":                  {dateString},
				"email":                {"hello@example.com"},
				"firstname":            {"Rudolph"},
				"isAttorneyActive":     {"true"},
				"phoneNumber":          {"0123456789"},
				"postcode":             {"57797"},
				"salutation":           {"Rev"},
				"surname":              {"Stotesbury"},
				"update-next-attorney": {"true"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2&caseType=lpa&attorneyId=2", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			if isHtmx {
				r.Header.Add("HX-Request", "true")
			}
			w := httptest.NewRecorder()

			err := UpdateAttorney(client, template.Func)(w, r)
			resp := w.Result()

			if !isHtmx {
				expectedRedirect := RedirectError("/create-trust-corporation?id=1&caseId=2&trustCorporationId=3&replacement=false")
				assert.Equal(t, err, expectedRedirect)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}
