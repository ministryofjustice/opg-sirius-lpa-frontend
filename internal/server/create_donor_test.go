package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCreateDonorClient struct {
	mock.Mock
}

func (m *mockCreateDonorClient) CreatePerson(ctx sirius.Context, personType sirius.PersonType, personData sirius.Person) (sirius.Person, error) {
	args := m.Called(ctx, personType, personData)
	return args.Get(0).(sirius.Person), args.Error(1)
}

func TestGetCreateDonor(t *testing.T) {
	client := &mockCreateDonorClient{}

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createDonorData{}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/create-donor", nil)
	w := httptest.NewRecorder()

	err := CreateDonor(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateDonor(t *testing.T) {
	client := &mockCreateDonorClient{}
	client.
		On("CreatePerson", mock.Anything, sirius.PersonTypeDonor, sirius.Person{
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
		}).
		Return(sirius.Person{ID: 809, UID: "7123-4567-8901"}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createDonorData{
			Donor: sirius.Person{
				ID:  809,
				UID: "7123-4567-8901",
			},
			Success: true,
		}).
		Return(nil)

	form := url.Values{
		"salutation":        {"Rev"},
		"firstname":         {"Rudolph"},
		"middlenames":       {"Modesto"},
		"surname":           {"Stotesbury"},
		"dob":               {"1981-10-03"},
		"previouslyKnownAs": {"Loma"},
		"alsoKnownAs":       {"Modesto"},
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

	r, _ := http.NewRequest(http.MethodPost, "/create-donor", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateDonor(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateDonorWhenAPIFails(t *testing.T) {
	expectedError := errors.New("failed to create donor")

	client := &mockCreateDonorClient{}
	client.
		On("CreatePerson", mock.Anything, sirius.PersonTypeDonor, sirius.Person{
			Firstname: "Rudolph",
			Surname:   "Stotesbury",
		}).
		Return(sirius.Person{}, expectedError)

	template := &mockTemplate{}

	form := url.Values{
		"firstname": {"Rudolph"},
		"surname":   {"Stotesbury"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/create-donor", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateDonor(client, template.Func)(w, r)
	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateDonorWhenValidationError(t *testing.T) {
	client := &mockCreateDonorClient{}
	client.
		On("CreatePerson", mock.Anything, sirius.PersonTypeDonor, sirius.Person{
			Firstname: "Rudolph",
		}).
		Return(sirius.Person{}, sirius.ValidationError{Field: sirius.FieldErrors{
			"surname": {"required": "This field is required"},
		}})

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createDonorData{
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

	r, _ := http.NewRequest(http.MethodPost, "/create-donor", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateDonor(client, template.Func)(w, r)
	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}
