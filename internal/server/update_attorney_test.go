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

//var mockRelationshipToDonorCategories = []sirius.RefDataItem{
//	{
//		Handle: "LPA_DONOR",
//		Label:  "LPA Donor",
//	},
//}

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
				On("Func", mock.Anything, AttorneyData{
					DonorId:              1,
					CaseId:               2,
					RelationshipToDonors: mockRelationshipToDonorCategories,
					Attorney:             existingAttorney,
					IsEditing:            true,
					Title:                "Update attorney details",
					CaseType:             caseType,
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
					On("Func", mock.Anything, AttorneyData{
						IsPartial:            true,
						DonorId:              1,
						CaseId:               2,
						RelationshipToDonors: mockRelationshipToDonorCategories,
						Attorney:             updatedAttorney,
						IsEditing:            true,
						Title:                "Update attorney details",
						HtmxRedirect:         "/create-epa?id=1&caseId=2",
						HtmxSwap:             "innerHTML show:#accordion-create-epa-heading-3:top",
						CaseType:             "epa",
					}).
					Return(nil)
			}

			form := url.Values{
				"salutation":          {"Rev"},
				"firstname":           {"Rudolph"},
				"middlenames":         {"Modesto"},
				"surname":             {"Stotesbury"},
				"dob":                 {dateString},
				"addressLine1":        {"Rotonda Gerardo 769"},
				"addressLine2":        {"Appartamento 94"},
				"addressLine3":        {"Augusto terme"},
				"town":                {"San Sabazio"},
				"county":              {"Benevento"},
				"postcode":            {"57797"},
				"country":             {"Italy"},
				"isAirmailRequired":   {"true"},
				"phoneNumber":         {"079876543345"},
				"email":               {"rm2@email.test"},
				"relationshipToDonor": {"no relation"},
				"isAttorneyActive":    {"true"},
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
					On("Func", mock.Anything, AttorneyData{
						IsPartial:            true,
						IsEditing:            true,
						DonorId:              1,
						CaseId:               2,
						RelationshipToDonors: mockRelationshipToDonorCategories,
						Attorney:             updatedAttorney,
						Title:                "Update attorney details",
						HtmxRedirect:         "/create-trust-corporation?id=1&caseId=2&trustCorporationId=3&replacement=false",
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
