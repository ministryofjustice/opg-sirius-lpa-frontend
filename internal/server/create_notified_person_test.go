package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCreateNotifiedPersonClient struct {
	mock.Mock
}

func (m *mockCreateNotifiedPersonClient) CreateNotifiedPerson(ctx sirius.Context, caseId int, notifiedPerson sirius.NotifiedPerson) error {
	args := m.Called(ctx, caseId, notifiedPerson)
	return args.Error(0)
}

func (m *mockCreateNotifiedPersonClient) Lpa(ctx sirius.Context, caseId int) (sirius.Lpa, error) {
	args := m.Called(ctx, caseId)
	return args.Get(0).(sirius.Lpa), args.Error(1)
}

func (m *mockCreateNotifiedPersonClient) RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error) {
	args := m.Called(ctx, category)
	if args.Get(0) != nil {
		return args.Get(0).([]sirius.RefDataItem), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockCreateNotifiedPersonClient) UpdateNotifiedPerson(ctx sirius.Context, notifiedPersonId int, notifiedPerson sirius.NotifiedPerson) error {
	args := m.Called(ctx, notifiedPersonId, notifiedPerson)
	return args.Error(0)
}

func TestGetCreateNotifiedPerson(t *testing.T) {
	for _, isHtmx := range []bool{false, true} {
		t.Run("Is Htmx: "+strconv.FormatBool(isHtmx), func(t *testing.T) {
			client := &mockCreateNotifiedPersonClient{}
			client.
				On("RefDataByCategory", mock.Anything, sirius.RelationshipToDonorCategory).
				Return(mockRelationshipToDonorCategories, nil).
				On("Lpa", mock.Anything, 2).
				Return(sirius.Lpa{}, nil)

			expectedData := createNotifiedPersonData{
				IsPartial:              isHtmx,
				DonorId:                1,
				CaseId:                 2,
				RelationshipToDonors:   mockRelationshipToDonorCategories,
				Title:                  "Add a notified person",
				AllowNewNotifiedPerson: true,
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

			err := CreateNotifiedPerson(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestGetEditNotifiedPerson(t *testing.T) {
	for _, isHtmx := range []bool{false, true} {
		t.Run("Is Htmx: "+strconv.FormatBool(isHtmx), func(t *testing.T) {
			existingNotifiedPerson := sirius.NotifiedPerson{
				Person: sirius.Person{
					ID:        4,
					Firstname: "Rudolph",
					Surname:   "Stotesbury",
				},
				NoticeGivenDate: sirius.DateString("2022-04-05"),
			}

			client := &mockCreateNotifiedPersonClient{}
			client.
				On("RefDataByCategory", mock.Anything, sirius.RelationshipToDonorCategory).
				Return(mockRelationshipToDonorCategories, nil).
				On("Lpa", mock.Anything, 2).
				Return(sirius.Lpa{Case: sirius.Case{NotifiedPersons: []sirius.NotifiedPerson{existingNotifiedPerson}}}, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, createNotifiedPersonData{
					IsPartial:              isHtmx,
					DonorId:                1,
					CaseId:                 2,
					RelationshipToDonors:   mockRelationshipToDonorCategories,
					NotifiedPerson:         existingNotifiedPerson,
					IsEditing:              true,
					Title:                  "Update notified person details",
					AllowNewNotifiedPerson: true,
				}).
				Return(nil)

			r, _ := http.NewRequest(http.MethodGet, "/?id=1&caseId=2&notifiedPersonId=4", nil)
			w := httptest.NewRecorder()

			if isHtmx {
				r.Header.Add("HX-Request", "true")
			}

			err := CreateNotifiedPerson(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestGetCreateNotifiedPersonBadQuery(t *testing.T) {
	testCases := map[string]string{
		"no-id":       "/",
		"bad-id":      "/?id=test",
		"bad-case-id": "/?id=123&caseId=test",
	}

	for name, query := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, query, nil)
			w := httptest.NewRecorder()

			err := CreateNotifiedPerson(nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestGetCreateNotifiedPersonWhenRefDataErrors(t *testing.T) {
	client := &mockCreateNotifiedPersonClient{}
	client.
		On("RefDataByCategory", mock.Anything, sirius.RelationshipToDonorCategory).
		Return([]sirius.RefDataItem{}, errExample)

	r, _ := http.NewRequest(http.MethodGet, "/?id=1&caseId=2", nil)
	w := httptest.NewRecorder()

	err := CreateNotifiedPerson(client, nil)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetCreateNotifiedPersonWhenLpaErrors(t *testing.T) {
	client := &mockCreateNotifiedPersonClient{}
	client.
		On("RefDataByCategory", mock.Anything, sirius.RelationshipToDonorCategory).
		Return(mockRelationshipToDonorCategories, nil).
		On("Lpa", mock.Anything, 2).
		Return(sirius.Lpa{}, errExample)

	r, _ := http.NewRequest(http.MethodGet, "/?id=1&caseId=2", nil)
	w := httptest.NewRecorder()

	err := CreateNotifiedPerson(client, nil)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestPostCreateNotifiedPerson(t *testing.T) {
	for _, isHtmx := range []bool{false, true} {
		t.Run("Is Htmx: "+strconv.FormatBool(isHtmx), func(t *testing.T) {
			dateString := "2022-04-05"
			notifiedPerson := sirius.NotifiedPerson{
				Person: sirius.Person{
					Salutation:        "Rev",
					Firstname:         "Rudolph",
					Middlenames:       "Modesto",
					Surname:           "Stotesbury",
					AddressLine1:      "Rotonda Gerardo 769",
					AddressLine2:      "Appartamento 94",
					AddressLine3:      "Augusto terme",
					Town:              "San Sabazio",
					County:            "Benevento",
					Postcode:          "57797",
					Country:           "Italy",
					IsAirmailRequired: true,
				},
				NoticeGivenDate: sirius.DateString(dateString),
			}
			client := &mockCreateNotifiedPersonClient{}
			client.
				On("RefDataByCategory", mock.Anything, sirius.RelationshipToDonorCategory).
				Return(mockRelationshipToDonorCategories, nil).
				On("Lpa", mock.Anything, 2).
				Return(sirius.Lpa{}, nil).
				On("CreateNotifiedPerson", mock.Anything, 2, notifiedPerson).
				Return(nil)

			template := &mockTemplate{}
			if isHtmx {
				template.
					On("Func", mock.Anything, createNotifiedPersonData{
						IsPartial:              true,
						DonorId:                1,
						CaseId:                 2,
						RelationshipToDonors:   mockRelationshipToDonorCategories,
						NotifiedPerson:         notifiedPerson,
						Title:                  "Add a notified person",
						AllowNewNotifiedPerson: true,
						HtmxRedirect:           "/create-lpa?id=1&caseId=2",
						HtmxSwap:               "innerHTML show:#accordion-create-lpa-heading-2:top",
					}).
					Return(nil)
			}

			form := url.Values{
				"salutation":        {"Rev"},
				"firstname":         {"Rudolph"},
				"middlenames":       {"Modesto"},
				"surname":           {"Stotesbury"},
				"noticeGivenDate":   {dateString},
				"addressLine1":      {"Rotonda Gerardo 769"},
				"addressLine2":      {"Appartamento 94"},
				"addressLine3":      {"Augusto terme"},
				"town":              {"San Sabazio"},
				"county":            {"Benevento"},
				"postcode":          {"57797"},
				"country":           {"Italy"},
				"isAirmailRequired": {"true"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			if isHtmx {
				r.Header.Add("HX-Request", "true")
			}
			w := httptest.NewRecorder()

			err := CreateNotifiedPerson(client, template.Func)(w, r)
			resp := w.Result()

			if !isHtmx {
				expectedError := RedirectError("/create-lpa?id=1&caseId=2#accordion-create-lpa-heading-2")
				assert.Equal(t, expectedError, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostEditNotifiedPerson(t *testing.T) {
	for _, isHtmx := range []bool{false, true} {
		t.Run("Is Htmx: "+strconv.FormatBool(isHtmx), func(t *testing.T) {
			dateString := "2022-04-05"
			existingNotifiedPerson := sirius.NotifiedPerson{Person: sirius.Person{ID: 4}}
			updatedNotifiedPerson := sirius.NotifiedPerson{
				Person: sirius.Person{
					Salutation:        "Rev",
					Firstname:         "Rudolph",
					Middlenames:       "Modesto",
					Surname:           "Stotesbury",
					AddressLine1:      "Rotonda Gerardo 769",
					AddressLine2:      "Appartamento 94",
					AddressLine3:      "Augusto terme",
					Town:              "San Sabazio",
					County:            "Benevento",
					Postcode:          "57797",
					Country:           "Italy",
					IsAirmailRequired: true,
				},
				NoticeGivenDate: sirius.DateString(dateString),
			}

			client := &mockCreateNotifiedPersonClient{}
			client.
				On("RefDataByCategory", mock.Anything, sirius.RelationshipToDonorCategory).
				Return(mockRelationshipToDonorCategories, nil).
				On("Lpa", mock.Anything, 2).
				Return(sirius.Lpa{Case: sirius.Case{NotifiedPersons: []sirius.NotifiedPerson{existingNotifiedPerson}}}, nil).
				On("UpdateNotifiedPerson", mock.Anything, 4, updatedNotifiedPerson).
				Return(nil)

			template := &mockTemplate{}
			if isHtmx {
				template.
					On("Func", mock.Anything, createNotifiedPersonData{
						IsPartial:              true,
						DonorId:                1,
						CaseId:                 2,
						RelationshipToDonors:   mockRelationshipToDonorCategories,
						NotifiedPerson:         updatedNotifiedPerson,
						IsEditing:              true,
						Title:                  "Update notified person details",
						AllowNewNotifiedPerson: true,
						HtmxRedirect:           "/create-lpa?id=1&caseId=2",
						HtmxSwap:               "innerHTML show:#accordion-create-lpa-heading-2:top",
					}).
					Return(nil)
			}

			form := url.Values{
				"salutation":        {"Rev"},
				"firstname":         {"Rudolph"},
				"middlenames":       {"Modesto"},
				"surname":           {"Stotesbury"},
				"noticeGivenDate":   {dateString},
				"addressLine1":      {"Rotonda Gerardo 769"},
				"addressLine2":      {"Appartamento 94"},
				"addressLine3":      {"Augusto terme"},
				"town":              {"San Sabazio"},
				"county":            {"Benevento"},
				"postcode":          {"57797"},
				"country":           {"Italy"},
				"isAirmailRequired": {"true"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2&notifiedPersonId=4", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			if isHtmx {
				r.Header.Add("HX-Request", "true")
			}
			w := httptest.NewRecorder()

			err := CreateNotifiedPerson(client, template.Func)(w, r)
			resp := w.Result()

			if !isHtmx {
				expectedError := RedirectError("/create-lpa?id=1&caseId=2#accordion-create-lpa-heading-2")
				assert.Equal(t, expectedError, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostCreateNotifiedPersonAddAnother(t *testing.T) {
	for _, isHtmx := range []bool{false, true} {
		t.Run("Is Htmx: "+strconv.FormatBool(isHtmx), func(t *testing.T) {
			dateString := "2022-04-05"
			notifiedPerson := sirius.NotifiedPerson{
				Person: sirius.Person{
					Salutation:        "Rev",
					Firstname:         "Rudolph",
					Middlenames:       "Modesto",
					Surname:           "Stotesbury",
					AddressLine1:      "Rotonda Gerardo 769",
					AddressLine2:      "Appartamento 94",
					AddressLine3:      "Augusto terme",
					Town:              "San Sabazio",
					County:            "Benevento",
					Postcode:          "57797",
					Country:           "Italy",
					IsAirmailRequired: true,
				},
				NoticeGivenDate: sirius.DateString(dateString),
			}
			client := &mockCreateNotifiedPersonClient{}
			client.
				On("RefDataByCategory", mock.Anything, sirius.RelationshipToDonorCategory).
				Return(mockRelationshipToDonorCategories, nil).
				On("Lpa", mock.Anything, 2).
				Return(sirius.Lpa{}, nil).
				On("CreateNotifiedPerson", mock.Anything, 2, notifiedPerson).
				Return(nil)

			template := &mockTemplate{}
			if isHtmx {
				template.
					On("Func", mock.Anything, createNotifiedPersonData{
						IsPartial:              true,
						DonorId:                1,
						CaseId:                 2,
						RelationshipToDonors:   mockRelationshipToDonorCategories,
						NotifiedPerson:         notifiedPerson,
						Title:                  "Add a notified person",
						AllowNewNotifiedPerson: true,
						HtmxRedirect:           "/create-notified-person?id=1&caseId=2",
						HtmxSwap:               "innerHTML scroll:.action-panel__content:top",
					}).
					Return(nil)
			}

			form := url.Values{
				"salutation":                  {"Rev"},
				"firstname":                   {"Rudolph"},
				"middlenames":                 {"Modesto"},
				"surname":                     {"Stotesbury"},
				"noticeGivenDate":             {dateString},
				"addressLine1":                {"Rotonda Gerardo 769"},
				"addressLine2":                {"Appartamento 94"},
				"addressLine3":                {"Augusto terme"},
				"town":                        {"San Sabazio"},
				"county":                      {"Benevento"},
				"postcode":                    {"57797"},
				"country":                     {"Italy"},
				"isAirmailRequired":           {"true"},
				"add-another-notified-person": {"true"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			if isHtmx {
				r.Header.Add("HX-Request", "true")
			}
			w := httptest.NewRecorder()

			err := CreateNotifiedPerson(client, template.Func)(w, r)
			resp := w.Result()

			if !isHtmx {
				expectedError := RedirectError("/create-notified-person?id=1&caseId=2")
				assert.Equal(t, expectedError, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostEditNotifiedPersonNextAnother(t *testing.T) {
	for _, isHtmx := range []bool{false, true} {
		t.Run("Is Htmx: "+strconv.FormatBool(isHtmx), func(t *testing.T) {
			dateString := "2022-04-05"
			existingNotifiedPerson := sirius.NotifiedPerson{Person: sirius.Person{ID: 4}}
			updatedNotifiedPerson := sirius.NotifiedPerson{
				Person: sirius.Person{
					Salutation:        "Rev",
					Firstname:         "Rudolph",
					Middlenames:       "Modesto",
					Surname:           "Stotesbury",
					AddressLine1:      "Rotonda Gerardo 769",
					AddressLine2:      "Appartamento 94",
					AddressLine3:      "Augusto terme",
					Town:              "San Sabazio",
					County:            "Benevento",
					Postcode:          "57797",
					Country:           "Italy",
					IsAirmailRequired: true,
				},
				NoticeGivenDate: sirius.DateString(dateString),
			}

			client := &mockCreateNotifiedPersonClient{}
			client.
				On("RefDataByCategory", mock.Anything, sirius.RelationshipToDonorCategory).
				Return(mockRelationshipToDonorCategories, nil).
				On("Lpa", mock.Anything, 2).
				Return(sirius.Lpa{Case: sirius.Case{NotifiedPersons: []sirius.NotifiedPerson{
					existingNotifiedPerson,
					{Person: sirius.Person{ID: 5}},
				}}}, nil).
				On("UpdateNotifiedPerson", mock.Anything, 4, updatedNotifiedPerson).
				Return(nil)

			template := &mockTemplate{}
			if isHtmx {
				template.
					On("Func", mock.Anything, createNotifiedPersonData{
						IsPartial:              true,
						DonorId:                1,
						CaseId:                 2,
						RelationshipToDonors:   mockRelationshipToDonorCategories,
						NotifiedPerson:         updatedNotifiedPerson,
						IsEditing:              true,
						Title:                  "Update notified person details",
						NextNotifiedPersonId:   5,
						AllowNewNotifiedPerson: true,
						HtmxRedirect:           "/create-notified-person?id=1&caseId=2&notifiedPersonId=5",
						HtmxSwap:               "innerHTML scroll:.action-panel__content:top",
					}).
					Return(nil)
			}

			form := url.Values{
				"salutation":           {"Rev"},
				"firstname":            {"Rudolph"},
				"middlenames":          {"Modesto"},
				"surname":              {"Stotesbury"},
				"noticeGivenDate":      {dateString},
				"addressLine1":         {"Rotonda Gerardo 769"},
				"addressLine2":         {"Appartamento 94"},
				"addressLine3":         {"Augusto terme"},
				"town":                 {"San Sabazio"},
				"county":               {"Benevento"},
				"postcode":             {"57797"},
				"country":              {"Italy"},
				"isAirmailRequired":    {"true"},
				"next-notified-person": {"true"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2&notifiedPersonId=4", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			if isHtmx {
				r.Header.Add("HX-Request", "true")
			}
			w := httptest.NewRecorder()

			err := CreateNotifiedPerson(client, template.Func)(w, r)
			resp := w.Result()

			if !isHtmx {
				expectedRedirect := RedirectError("/create-notified-person?id=1&caseId=2&notifiedPersonId=5")
				assert.Equal(t, expectedRedirect, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostCreateNotifiedPersonWhenValidationError(t *testing.T) {
	for _, isHtmx := range []bool{false, true} {
		t.Run("Is Htmx: "+strconv.FormatBool(isHtmx), func(t *testing.T) {
			expectedError := sirius.ValidationError{
				Field: sirius.FieldErrors{"field": {"": "problem"}},
			}

			dateString := "2022-04-05"
			notifiedPerson := sirius.NotifiedPerson{
				Person: sirius.Person{
					Salutation:        "Rev",
					Firstname:         "Rudolph",
					Middlenames:       "Modesto",
					Surname:           "Stotesbury",
					AddressLine1:      "Rotonda Gerardo 769",
					AddressLine2:      "Appartamento 94",
					AddressLine3:      "Augusto terme",
					Town:              "San Sabazio",
					County:            "Benevento",
					Postcode:          "57797",
					Country:           "Italy",
					IsAirmailRequired: true,
				},
				NoticeGivenDate: sirius.DateString(dateString),
			}

			client := &mockCreateNotifiedPersonClient{}
			client.
				On("RefDataByCategory", mock.Anything, sirius.RelationshipToDonorCategory).
				Return(mockRelationshipToDonorCategories, nil).
				On("Lpa", mock.Anything, 2).
				Return(sirius.Lpa{}, nil).
				On("CreateNotifiedPerson", mock.Anything, 2, notifiedPerson).
				Return(expectedError)

			expectedData := createNotifiedPersonData{
				IsPartial:              isHtmx,
				NotifiedPerson:         notifiedPerson,
				DonorId:                1,
				CaseId:                 2,
				Error:                  expectedError,
				RelationshipToDonors:   mockRelationshipToDonorCategories,
				Title:                  "Add a notified person",
				AllowNewNotifiedPerson: true,
			}

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, expectedData).
				Return(nil)

			form := url.Values{
				"salutation":        {"Rev"},
				"firstname":         {"Rudolph"},
				"middlenames":       {"Modesto"},
				"surname":           {"Stotesbury"},
				"noticeGivenDate":   {dateString},
				"addressLine1":      {"Rotonda Gerardo 769"},
				"addressLine2":      {"Appartamento 94"},
				"addressLine3":      {"Augusto terme"},
				"town":              {"San Sabazio"},
				"county":            {"Benevento"},
				"postcode":          {"57797"},
				"country":           {"Italy"},
				"isAirmailRequired": {"true"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			if isHtmx {
				r.Header.Add("HX-Request", "true")
			}
			w := httptest.NewRecorder()

			err := CreateNotifiedPerson(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestAllowNewNotifiedPerson(t *testing.T) {
	assert.True(t, allowNewNotifiedPerson(0))
	assert.True(t, allowNewNotifiedPerson(3))
	assert.False(t, allowNewNotifiedPerson(4))
}

func TestGetNextNotifiedPersonIdWillReturnNextNumber(t *testing.T) {
	result := GetNextNotifiedPersonId(2, []sirius.NotifiedPerson{
		{Person: sirius.Person{ID: 1}},
		{Person: sirius.Person{ID: 2}},
		{Person: sirius.Person{ID: 3}},
	})
	expected := 3
	assert.Equal(t, expected, result)
}

func TestGetNextNotifiedPersonIdWillReturnNextHigherIdWhenSequenceHasGaps(t *testing.T) {
	result := GetNextNotifiedPersonId(2, []sirius.NotifiedPerson{
		{Person: sirius.Person{ID: 1}},
		{Person: sirius.Person{ID: 5}},
		{Person: sirius.Person{ID: 2}},
		{Person: sirius.Person{ID: 9}},
	})
	expected := 5
	assert.Equal(t, expected, result)
}

func TestGetNextNotifiedPersonIdLastAttorneyWillReturnZero(t *testing.T) {
	result := GetNextNotifiedPersonId(3, []sirius.NotifiedPerson{
		{Person: sirius.Person{ID: 1}},
		{Person: sirius.Person{ID: 2}},
		{Person: sirius.Person{ID: 3}},
	})
	expected := 0
	assert.Equal(t, expected, result)
}
