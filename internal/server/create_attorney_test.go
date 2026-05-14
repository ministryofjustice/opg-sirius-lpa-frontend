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

func (m *mockCreateAttorneyClient) CreateAttorney(ctx sirius.Context, caseId int, attorney sirius.Attorney) error {
	args := m.Called(ctx, caseId, attorney)
	return args.Error(0)
}

func (m *mockCreateAttorneyClient) RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error) {
	args := m.Called(ctx, category)
	if args.Get(0) != nil {
		return args.Get(0).([]sirius.RefDataItem), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockCreateAttorneyClient) UpdateAttorney(ctx sirius.Context, attorneyId int, attorney sirius.Attorney) error {
	args := m.Called(ctx, attorneyId, attorney)
	return args.Error(0)
}

func (m *mockCreateAttorneyClient) UpdateCorrespondent(ctx sirius.Context, correspondentId int, correspondent sirius.Correspondent) error {
	args := m.Called(ctx, correspondentId, correspondent)
	return args.Error(0)
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

			expectedData := createAttorneyData{
				DonorId:              1,
				CaseId:               2,
				RelationshipToDonors: mockRelationshipToDonorCategories,
				Attorney:             sirius.Attorney{SystemStatus: shared.BoolPtr(true)},
				Title:                "Add an attorney",
			}
			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, expectedData).
				Return(nil)

			r, _ := http.NewRequest(http.MethodGet, "/?id=1&caseId=2", nil)
			w := httptest.NewRecorder()

			if isHtmx {
				r.Header.Add("HX-Request", "true")
			}

			err := CreateAttorney(client, template.Func, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestGetEditAttorney(t *testing.T) {
	existingAttorney := sirius.Attorney{
		Person: sirius.Person{
			ID:        4,
			Firstname: "Rudolph",
			Surname:   "Stotesbury",
		},
		RelationshipToDonor: "no relation",
		SystemStatus:        shared.BoolPtr(true),
	}

	client := &mockCreateAttorneyClient{}
	client.
		On("RefDataByCategory", mock.Anything, sirius.RelationshipToDonorCategory).
		Return(mockRelationshipToDonorCategories, nil).
		On("Epa", mock.Anything, 2).
		Return(sirius.Epa{Case: sirius.Case{Attorneys: []sirius.Attorney{existingAttorney}}}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createAttorneyData{
			DonorId:              1,
			CaseId:               2,
			RelationshipToDonors: mockRelationshipToDonorCategories,
			Attorney:             existingAttorney,
			IsEditing:            true,
			Title:                "Update attorney details",
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=1&caseId=2&attorneyId=4", nil)
	w := httptest.NewRecorder()

	err := CreateAttorney(client, template.Func, nil)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetCreateAttorneyBadQuery(t *testing.T) {
	testCases := map[string]string{
		"no-id":       "/",
		"bad-id":      "/?id=test",
		"bad-case-id": "/?id=123&caseId=test",
	}

	for name, query := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, query, nil)
			w := httptest.NewRecorder()

			err := CreateAttorney(nil, nil, nil)(w, r)

			assert.NotNil(t, err)
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

	err := CreateAttorney(client, nil, nil)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestPostCreateAttorney(t *testing.T) {
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
				On("CreateAttorney", mock.Anything, 2, attorney).
				Return(nil).
				On("RefDataByCategory", mock.Anything, sirius.RelationshipToDonorCategory).
				Return(mockRelationshipToDonorCategories, nil)

			template := &mockTemplate{}
			partialTmpl := &mockTemplate{}

			if isHtmx {
				partialTmpl.
					On("Func", mock.Anything, createAttorneyData{
						DonorId:              1,
						CaseId:               2,
						RelationshipToDonors: mockRelationshipToDonorCategories,
						Attorney:             attorney,
						IsEditing:            false,
						Title:                "Add an attorney",
						HtmxRedirect:         "/create-epa?id=1&caseId=2",
						HtmxSwap:             "innerHTML show:#accordion-create-epa-heading-3:top",
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

			r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			if isHtmx {
				r.Header.Add("HX-Request", "true")
			}
			w := httptest.NewRecorder()

			err := CreateAttorney(client, template.Func, partialTmpl.Func)(w, r)
			resp := w.Result()

			if !isHtmx {
				expectedError := RedirectError("/create-epa?id=1&caseId=2#accordion-create-epa-heading-3")
				assert.Equal(t, err, expectedError)
			}
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template, partialTmpl)
		})
	}
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

			client := &mockCreateAttorneyClient{}
			client.
				On("RefDataByCategory", mock.Anything, sirius.RelationshipToDonorCategory).
				Return(mockRelationshipToDonorCategories, nil).
				On("Epa", mock.Anything, 2).
				Return(sirius.Epa{Case: sirius.Case{Attorneys: []sirius.Attorney{existingAttorney}}}, nil).
				On("UpdateAttorney", mock.Anything, 4, updatedAttorney).
				Return(nil)

			template := &mockTemplate{}
			partialTemplate := &mockTemplate{}

			if isHtmx {
				partialTemplate.
					On("Func", mock.Anything, createAttorneyData{
						DonorId:              1,
						CaseId:               2,
						RelationshipToDonors: mockRelationshipToDonorCategories,
						Attorney:             updatedAttorney,
						IsEditing:            true,
						Title:                "Update attorney details",
						HtmxRedirect:         "/create-epa?id=1&caseId=2",
						HtmxSwap:             "innerHTML show:#accordion-create-epa-heading-3:top",
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

			r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2&attorneyId=4", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			if isHtmx {
				r.Header.Add("HX-Request", "true")
			}
			w := httptest.NewRecorder()

			err := CreateAttorney(client, template.Func, partialTemplate.Func)(w, r)
			resp := w.Result()

			if !isHtmx {
				expectedError := RedirectError("/create-epa?id=1&caseId=2#accordion-create-epa-heading-3")
				assert.Equal(t, err, expectedError)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template, partialTemplate)
		})
	}
}

func TestPostCreateAttorneyAddAnother(t *testing.T) {
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
				On("CreateAttorney", mock.Anything, 2, attorney).
				Return(nil).
				On("RefDataByCategory", mock.Anything, sirius.RelationshipToDonorCategory).
				Return(mockRelationshipToDonorCategories, nil)

			template := &mockTemplate{}
			partialTemplate := &mockTemplate{}

			if isHtmx {
				partialTemplate.
					On("Func", mock.Anything, createAttorneyData{
						DonorId:              1,
						CaseId:               2,
						RelationshipToDonors: mockRelationshipToDonorCategories,
						Attorney:             attorney,
						Title:                "Add an attorney",
						HtmxRedirect:         "/create-attorney?id=1&caseId=2",
						HtmxSwap:             "innerHTML scroll:.action-panel__content:top",
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
				"add-another":         {"true"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			if isHtmx {
				r.Header.Add("HX-Request", "true")
			}
			w := httptest.NewRecorder()

			err := CreateAttorney(client, template.Func, partialTemplate.Func)(w, r)
			resp := w.Result()

			if !isHtmx {
				expectedError := RedirectError("/create-attorney?id=1&caseId=2")
				assert.Equal(t, err, expectedError)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template, partialTemplate)
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
				On("CreateAttorney", mock.Anything, 2, attorney).
				Return(expectedError).
				On("RefDataByCategory", mock.Anything, sirius.RelationshipToDonorCategory).
				Return(mockRelationshipToDonorCategories, nil)

			template := &mockTemplate{}
			partialTemplate := &mockTemplate{}

			expectedData := createAttorneyData{
				Attorney:             attorney,
				DonorId:              1,
				CaseId:               2,
				Error:                expectedError,
				RelationshipToDonors: mockRelationshipToDonorCategories,
				Title:                "Add an attorney",
			}

			if isHtmx {
				partialTemplate.
					On("Func", mock.Anything, expectedData).
					Return(nil)
			} else {
				template.
					On("Func", mock.Anything, expectedData).
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

			r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			if isHtmx {
				r.Header.Add("HX-Request", "true")
			}
			w := httptest.NewRecorder()

			err := CreateAttorney(client, template.Func, partialTemplate.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template, partialTemplate)
		})
	}
}

func TestPostCreateAttorneyNextAnother(t *testing.T) {
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
				On("CreateAttorney", mock.Anything, 2, attorney).
				Return(nil).
				On("RefDataByCategory", mock.Anything, sirius.RelationshipToDonorCategory).
				Return(mockRelationshipToDonorCategories, nil)

			template := &mockTemplate{}
			partialTemplate := &mockTemplate{}

			if isHtmx {
				partialTemplate.
					On("Func", mock.Anything, createAttorneyData{
						DonorId:              1,
						CaseId:               2,
						RelationshipToDonors: mockRelationshipToDonorCategories,
						Attorney:             attorney,
						Title:                "Add an attorney",
						HtmxRedirect:         "/create-attorney?id=1&caseId=2&attorneyId=0",
						HtmxSwap:             "innerHTML scroll:.action-panel__content:top",
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
				"add-another":         {""},
				"next-attorney":       {"true"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			if isHtmx {
				r.Header.Add("HX-Request", "true")
			}
			w := httptest.NewRecorder()

			err := CreateAttorney(client, template.Func, partialTemplate.Func)(w, r)
			resp := w.Result()

			if !isHtmx {
				expectedRedirect := RedirectError("/create-attorney?id=1&caseId=2&attorneyId=0")
				assert.Equal(t, err, expectedRedirect)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template, partialTemplate)
		})
	}
}

func TestGetNextAttorneyIdAttorneyWillReturnNextNumber(t *testing.T) {
	result := GetNextAttorneyId(2, []sirius.Attorney{
		{
			Person: sirius.Person{ID: 1},
		},
		{
			Person: sirius.Person{ID: 2},
		},
		{
			Person: sirius.Person{ID: 3},
		},
	})
	expected := 3
	assert.Equal(t, expected, result)
}

func TestGetNextAttorneyIdAttorneyWillReturnNextHigherIdWhenSequenceHasGaps(t *testing.T) {
	result := GetNextAttorneyId(2, []sirius.Attorney{
		{
			Person: sirius.Person{ID: 1},
		},
		{
			Person: sirius.Person{ID: 5},
		},
		{
			Person: sirius.Person{ID: 2},
		},
		{
			Person: sirius.Person{ID: 9},
		},
	})
	expected := 5
	assert.Equal(t, expected, result)
}

func TestGetNextAttorneyIdLastAttorneyWillReturnZero(t *testing.T) {
	result := GetNextAttorneyId(3, []sirius.Attorney{
		{
			Person: sirius.Person{ID: 1},
		},
		{
			Person: sirius.Person{ID: 2},
		},
		{
			Person: sirius.Person{ID: 3},
		},
	})
	expected := 0
	assert.Equal(t, expected, result)
}
