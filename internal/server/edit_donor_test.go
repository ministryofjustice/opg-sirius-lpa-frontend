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

type mockEditDonorClient struct {
	mock.Mock
}

func (m *mockEditDonorClient) EditDonor(ctx sirius.Context, personID int, personData sirius.Person) error {
	args := m.Called(ctx, personID, personData)
	return args.Error(0)
}

func (m *mockEditDonorClient) Person(ctx sirius.Context, personID int) (sirius.Person, error) {
	args := m.Called(ctx, personID)
	return args.Get(0).(sirius.Person), args.Error(1)
}

func TestGetEditDonor(t *testing.T) {
	person := sirius.Person{
		Firstname: "Wanda",
		Surname:   "Bratu",
	}

	client := &mockEditDonorClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(person, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, donorData{
			Donor: person,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/edit-donor?id=123", nil)
	w := httptest.NewRecorder()

	err := EditDonor(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostEditDonor(t *testing.T) {
	newPerson := sirius.Person{
		Salutation:            "Rev",
		Firstname:             "Rudolph",
		Middlenames:           "Modesto",
		Surname:               "Stotesbury",
		DateOfBirth:           sirius.DateString("1981-10-03"),
		PreviouslyKnownAs:     "Loma",
		AlsoKnownAs:           "Modesto",
		AddressLine1:          "Rotonda Gerardo 769",
		AddressLine2:          "Appartamento 94",
		AddressLine3:          "Augusto terme",
		Town:                  "San Sabazio",
		County:                "Benevento",
		Postcode:              "57797",
		Country:               "Italy",
		IsAirmailRequired:     true,
		PhoneNumber:           "079876543345",
		Email:                 "rm2@email.test",
		CorrespondenceByPost:  true,
		CorrespondenceByEmail: true,
		CorrespondenceByPhone: false,
		CorrespondenceByWelsh: false,
		ResearchOptOut:        false,
	}

	client := &mockEditDonorClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)

	client.
		On("EditDonor", mock.Anything, 123, newPerson).
		Return(nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, donorData{
			Donor:   newPerson,
			Success: true,
		}).
		Return(nil)

	form := url.Values{
		"salutation":        {"Rev"},
		"firstname":         {"Rudolph"},
		"middlenames":       {"Modesto"},
		"surname":           {"Stotesbury"},
		"dob":               {"1981-10-03"},
		"previousNames":     {"Loma"},
		"otherNames":        {"Modesto"},
		"addressLine1":      {"Rotonda Gerardo 769"},
		"addressLine2":      {"Appartamento 94"},
		"addressLine3":      {"Augusto terme"},
		"town":              {"San Sabazio"},
		"county":            {"Benevento"},
		"postcode":          {"57797"},
		"country":           {"Italy"},
		"isAirmailRequired": {"Yes"},
		"phoneNumber":       {"079876543345"},
		"email":             {"rm2@email.test"},
		"sageId":            {""},
		"correspondenceBy":  {"post", "email"},
		"researchOptOut":    {"No"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/edit-donor?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := EditDonor(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostEditDonorWhenAPIFails(t *testing.T) {
	client := &mockEditDonorClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)

	client.
		On("EditDonor", mock.Anything, 123, sirius.Person{
			Firstname: "Rudolph",
			Surname:   "Stotesbury",
		}).
		Return(expectedError)

	template := &mockTemplate{}

	form := url.Values{
		"firstname": {"Rudolph"},
		"surname":   {"Stotesbury"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/edit-donor?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := EditDonor(client, template.Func)(w, r)
	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostEditDonorWhenValidationError(t *testing.T) {
	client := &mockEditDonorClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)

	client.
		On("EditDonor", mock.Anything, 123, sirius.Person{
			Firstname: "Rudolph",
		}).
		Return(sirius.ValidationError{Field: sirius.FieldErrors{
			"surname": {"required": "This field is required"},
		}})

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, donorData{
			Donor: sirius.Person{
				Firstname: "Rudolph",
			},
			Error: sirius.ValidationError{
				Field: sirius.FieldErrors{
					"surname": {"required": "This field is required"},
				},
			},
		}).
		Return(nil)

	form := url.Values{
		"firstname": {"Rudolph"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/edit-donor?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := EditDonor(client, template.Func)(w, r)
	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}
