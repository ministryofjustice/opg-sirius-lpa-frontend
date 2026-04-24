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

func (m *mockCreateCorrespondentClient) CreateCorrespondent(ctx sirius.Context, caseId int, correspondent sirius.Correspondent) error {
	args := m.Called(ctx, caseId, correspondent)
	return args.Error(0)
}

func TestGetCreateCorrespondent(t *testing.T) {
	client := &mockCreateCorrespondentClient{}

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createCorrespondentData{
			DonorId: 1,
			CaseId:  2,
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
	expectedError := RedirectError("/create-epa?id=1&caseId=2")
	dateString := "2022-04-05"
	correspondent := sirius.Correspondent{
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
	}
	client := &mockCreateCorrespondentClient{}
	client.
		On("CreateCorrespondent", mock.Anything, 2, correspondent).
		Return(nil)

	template := &mockTemplate{}

	form := url.Values{
		"salutation":        {"Rev"},
		"firstname":         {"Rudolph"},
		"middlenames":       {"Modesto"},
		"surname":           {"Stotesbury"},
		"dob":               {dateString},
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

	r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2", strings.NewReader(form.Encode()))
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

	dateString := "2022-04-05"
	correspondent := sirius.Correspondent{
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
	}

	client := &mockCreateCorrespondentClient{}
	client.
		On("CreateCorrespondent", mock.Anything, 2, correspondent).
		Return(expectedError)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createCorrespondentData{
			DonorId:       1,
			CaseId:        2,
			Error:         expectedError,
			Correspondent: correspondent,
		}).
		Return(nil)

	form := url.Values{
		"salutation":        {"Rev"},
		"firstname":         {"Rudolph"},
		"middlenames":       {"Modesto"},
		"surname":           {"Stotesbury"},
		"dob":               {dateString},
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

	r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateCorrespondent(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateCorrespondentCreationFails(t *testing.T) {
	dateString := "2022-04-05"
	correspondent := sirius.Correspondent{
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
	}

	client := &mockCreateCorrespondentClient{}
	client.
		On("CreateCorrespondent", mock.Anything, 2, correspondent).
		Return(errExample)

	template := &mockTemplate{}

	form := url.Values{
		"salutation":        {"Rev"},
		"firstname":         {"Rudolph"},
		"middlenames":       {"Modesto"},
		"surname":           {"Stotesbury"},
		"dob":               {dateString},
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

	r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateCorrespondent(client, template.Func)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client, template)
}
