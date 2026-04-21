package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCreateAttorneyClient struct {
	mock.Mock
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

var mockRelationshipToDonorCategories = []sirius.RefDataItem{
	{
		Handle: "LPA_DONOR",
		Label:  "LPA Donor",
	},
}

func TestGetCreateAttorney(t *testing.T) {
	client := &mockCreateAttorneyClient{}
	client.
		On("RefDataByCategory", mock.Anything, sirius.RelationshipToDonorCategory).
		Return(mockRelationshipToDonorCategories, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createAttorneyData{
			DonorId:              1,
			CaseId:               2,
			RelationshipToDonors: mockRelationshipToDonorCategories,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=1&caseId=2", nil)
	w := httptest.NewRecorder()

	err := CreateAttorney(client, template.Func)(w, r)
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

			err := CreateAttorney(nil, nil)(w, r)

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

	err := CreateAttorney(client, nil)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestPostCreateAttorney(t *testing.T) {
	expectedError := RedirectError("/create-epa?id=1&caseId=2")
	truePtr := true
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
		SystemStatus:        &truePtr,
	}
	client := &mockCreateAttorneyClient{}
	client.
		On("CreateAttorney", mock.Anything, 2, attorney).
		Return(nil).
		On("RefDataByCategory", mock.Anything, sirius.RelationshipToDonorCategory).
		Return(mockRelationshipToDonorCategories, nil)

	template := &mockTemplate{}

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
	w := httptest.NewRecorder()

	err := CreateAttorney(client, template.Func)(w, r)
	resp := w.Result()

	assert.Equal(t, err, expectedError)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateAttorneyWhenValidationError(t *testing.T) {
	expectedError := sirius.ValidationError{
		Field: sirius.FieldErrors{"field": {"": "problem"}},
	}

	truePtr := true
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
		SystemStatus:        &truePtr,
	}

	client := &mockCreateAttorneyClient{}
	client.
		On("CreateAttorney", mock.Anything, 2, attorney).
		Return(expectedError).
		On("RefDataByCategory", mock.Anything, sirius.RelationshipToDonorCategory).
		Return(mockRelationshipToDonorCategories, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createAttorneyData{
			Attorney:             attorney,
			DonorId:              1,
			CaseId:               2,
			Error:                expectedError,
			RelationshipToDonors: mockRelationshipToDonorCategories,
		}).
		Return(nil)

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
	w := httptest.NewRecorder()

	err := CreateAttorney(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}
