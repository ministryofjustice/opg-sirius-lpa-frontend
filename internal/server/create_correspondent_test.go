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

type mockCreateCorrespondentClient struct {
	mock.Mock
}

func (m *mockCreateCorrespondentClient) Epa(ctx sirius.Context, id int) (sirius.Epa, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Epa), args.Error(1)
}

func (m *mockCreateCorrespondentClient) CreateCorrespondent(ctx sirius.Context, caseId int, correspondent sirius.Correspondent) error {
	args := m.Called(ctx, caseId, correspondent)
	return args.Error(0)
}

func (m *mockCreateCorrespondentClient) UpdateCorrespondent(ctx sirius.Context, correspondentId int, correspondent sirius.Correspondent) error {
	args := m.Called(ctx, correspondentId, correspondent)
	return args.Error(0)
}

func TestGetCreateCorrespondent(t *testing.T) {
	client := &mockCreateCorrespondentClient{}
	client.
		On("Epa", mock.Anything, 2).
		Return(sirius.Epa{}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createCorrespondentData{
			DonorId: 1,
			CaseId:  2,
			Title:   "Add a correspondent",
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=1&caseId=2", nil)
	w := httptest.NewRecorder()

	err := CreateCorrespondent(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetEditCorrespondent(t *testing.T) {
	existingCorrespondent := sirius.Correspondent{
		Person: sirius.Person{
			ID:        7,
			Firstname: "Rudolph",
			Surname:   "Stotesbury",
		},
	}

	client := &mockCreateCorrespondentClient{}
	client.
		On("Epa", mock.Anything, 2).
		Return(sirius.Epa{Case: sirius.Case{Correspondent: &existingCorrespondent}}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createCorrespondentData{
			DonorId:       1,
			CaseId:        2,
			Correspondent: existingCorrespondent,
			IsEditing:     true,
			Title:         "Update correspondent details",
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=1&caseId=2", nil)
	w := httptest.NewRecorder()

	err := CreateCorrespondent(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetCreateCorrespondentBadQuery(t *testing.T) {
	testCases := map[string]string{
		"no-id":       "/",
		"bad-id":      "/?id=test",
		"bad-case-id": "/?id=123&caseId=test",
	}

	for name, query := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, query, nil)
			w := httptest.NewRecorder()

			err := CreateCorrespondent(nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestPostCreateCorrespondent(t *testing.T) {
	expectedError := RedirectError("/create-epa?id=1&caseId=2#accordion-create-epa-heading-3")
	correspondent := sirius.Correspondent{
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
			PhoneNumber:       "079876543345",
			Email:             "rm2@email.test",
		},
	}
	client := &mockCreateCorrespondentClient{}
	client.
		On("Epa", mock.Anything, 2).
		Return(sirius.Epa{}, nil).
		On("CreateCorrespondent", mock.Anything, 2, correspondent).
		Return(nil)

	template := &mockTemplate{}

	form := url.Values{
		"salutation":        {"Rev"},
		"firstname":         {"Rudolph"},
		"middlenames":       {"Modesto"},
		"surname":           {"Stotesbury"},
		"addressLine1":      {"Rotonda Gerardo 769"},
		"addressLine2":      {"Appartamento 94"},
		"addressLine3":      {"Augusto terme"},
		"town":              {"San Sabazio"},
		"county":            {"Benevento"},
		"postcode":          {"57797"},
		"country":           {"Italy"},
		"isAirmailRequired": {"true"},
		"phoneNumber":       {"079876543345"},
		"email":             {"rm2@email.test"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2#accordion-create-epa-heading-3", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateCorrespondent(client, template.Func)(w, r)
	resp := w.Result()

	assert.Equal(t, err, expectedError)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostEditCorrespondent(t *testing.T) {
	expectedError := RedirectError("/create-epa?id=1&caseId=2#accordion-create-epa-heading-3")
	existingCorrespondent := sirius.Correspondent{Person: sirius.Person{ID: 7}}
	updatedCorrespondent := sirius.Correspondent{
		Person: sirius.Person{
			ID:                7,
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
			PhoneNumber:       "079876543345",
			Email:             "rm2@email.test",
		},
	}

	client := &mockCreateCorrespondentClient{}
	client.
		On("Epa", mock.Anything, 2).
		Return(sirius.Epa{Case: sirius.Case{Correspondent: &existingCorrespondent}}, nil).
		On("UpdateCorrespondent", mock.Anything, 7, updatedCorrespondent).
		Return(nil)

	template := &mockTemplate{}

	form := url.Values{
		"salutation":        {"Rev"},
		"firstname":         {"Rudolph"},
		"middlenames":       {"Modesto"},
		"surname":           {"Stotesbury"},
		"addressLine1":      {"Rotonda Gerardo 769"},
		"addressLine2":      {"Appartamento 94"},
		"addressLine3":      {"Augusto terme"},
		"town":              {"San Sabazio"},
		"county":            {"Benevento"},
		"postcode":          {"57797"},
		"country":           {"Italy"},
		"isAirmailRequired": {"true"},
		"phoneNumber":       {"079876543345"},
		"email":             {"rm2@email.test"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2#accordion-create-epa-heading-3", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateCorrespondent(client, template.Func)(w, r)
	resp := w.Result()

	assert.Equal(t, err, expectedError)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateCorrespondentWhenValidationError(t *testing.T) {
	expectedError := sirius.ValidationError{
		Field: sirius.FieldErrors{"field": {"": "problem"}},
	}

	correspondent := sirius.Correspondent{
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
			PhoneNumber:       "079876543345",
			Email:             "rm2@email.test",
		},
	}

	client := &mockCreateCorrespondentClient{}
	client.
		On("Epa", mock.Anything, 2).
		Return(sirius.Epa{}, nil).
		On("CreateCorrespondent", mock.Anything, 2, correspondent).
		Return(expectedError)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createCorrespondentData{
			DonorId:       1,
			CaseId:        2,
			Error:         expectedError,
			Correspondent: correspondent,
			Title:         "Add a correspondent",
		}).
		Return(nil)

	form := url.Values{
		"salutation":        {"Rev"},
		"firstname":         {"Rudolph"},
		"middlenames":       {"Modesto"},
		"surname":           {"Stotesbury"},
		"addressLine1":      {"Rotonda Gerardo 769"},
		"addressLine2":      {"Appartamento 94"},
		"addressLine3":      {"Augusto terme"},
		"town":              {"San Sabazio"},
		"county":            {"Benevento"},
		"postcode":          {"57797"},
		"country":           {"Italy"},
		"isAirmailRequired": {"true"},
		"phoneNumber":       {"079876543345"},
		"email":             {"rm2@email.test"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/create-epa?id=1&caseId=2#accordion-create-epa-heading-3", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateCorrespondent(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateCorrespondentCreationFails(t *testing.T) {
	correspondent := sirius.Correspondent{
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
			PhoneNumber:       "079876543345",
			Email:             "rm2@email.test",
		},
	}

	client := &mockCreateCorrespondentClient{}
	client.
		On("Epa", mock.Anything, 2).
		Return(sirius.Epa{}, nil).
		On("CreateCorrespondent", mock.Anything, 2, correspondent).
		Return(errExample)

	template := &mockTemplate{}

	form := url.Values{
		"salutation":        {"Rev"},
		"firstname":         {"Rudolph"},
		"middlenames":       {"Modesto"},
		"surname":           {"Stotesbury"},
		"addressLine1":      {"Rotonda Gerardo 769"},
		"addressLine2":      {"Appartamento 94"},
		"addressLine3":      {"Augusto terme"},
		"town":              {"San Sabazio"},
		"county":            {"Benevento"},
		"postcode":          {"57797"},
		"country":           {"Italy"},
		"isAirmailRequired": {"true"},
		"phoneNumber":       {"079876543345"},
		"email":             {"rm2@email.test"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2#accordion-create-epa-heading-3", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateCorrespondent(client, template.Func)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client, template)
}
