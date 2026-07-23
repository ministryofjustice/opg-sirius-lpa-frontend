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

type mockCreateDonorClient struct {
	mock.Mock
}

func (m *mockCreateDonorClient) CreateDonor(ctx sirius.Context, personData sirius.Person) (sirius.Person, error) {
	args := m.Called(ctx, personData)
	return args.Get(0).(sirius.Person), args.Error(1)
}

func TestGetCreateDonor(t *testing.T) {
	client := &mockCreateDonorClient{}

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, donorData{
			IsNew: true,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/create-donor", nil)
	w := httptest.NewRecorder()

	err := CreateDonor(client, template.Func, nil)(PageVars{}, w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetCreateDonorHtmxRequest(t *testing.T) {
	client := &mockCreateDonorClient{}

	partialTemplate := &mockTemplate{}
	partialTemplate.
		On("Func", mock.Anything, donorData{
			IsNew:      true,
			DonorId:    123,
			CaseUids:   "&uid[]=7000-1234-1234",
			EntityType: "person",
		}).
		Return(nil)

	template := &mockTemplate{}

	r, _ := http.NewRequest(http.MethodGet, "/create-donor?id=123&entity=person&uid[]=7000-1234-1234", nil)
	r.Header.Add("HX-Request", "true")
	w := httptest.NewRecorder()

	err := CreateDonor(client, template.Func, partialTemplate.Func)(PageVars{}, w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	template.AssertNotCalled(t, "Func")
	mock.AssertExpectationsForObjects(t, client, partialTemplate)
}

func TestPostCreateDonor(t *testing.T) {
	client := &mockCreateDonorClient{}
	client.
		On("CreateDonor", mock.Anything, sirius.Person{
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
		On("Func", mock.Anything, donorData{
			IsNew: true,
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

	r, _ := http.NewRequest(http.MethodPost, "/create-donor", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateDonor(client, template.Func, nil)(PageVars{}, w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateDonorHtmxRequest(t *testing.T) {
	client := &mockCreateDonorClient{}
	client.
		On("CreateDonor", mock.Anything, sirius.Person{
			Firstname:   "Rudolph",
			Middlenames: "Modesto",
		}).
		Return(sirius.Person{ID: 809, UID: "7123-4567-8901"}, nil)

	partialTemplate := &mockTemplate{}
	partialTemplate.
		On("Func", mock.Anything, donorData{
			IsNew: true,
			Donor: sirius.Person{
				ID:  809,
				UID: "7123-4567-8901",
			},
			Success:    true,
			DonorId:    123,
			CaseUids:   "&uid[]=7000-1234-1234",
			EntityType: "person",
		}).
		Return(nil)

	template := &mockTemplate{}

	form := url.Values{
		"firstname":   {"Rudolph"},
		"middlenames": {"Modesto"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/create-donor?id=123&entity=person&uid[]=7000-1234-1234", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	r.Header.Add("HX-Request", "true")
	w := httptest.NewRecorder()

	err := CreateDonor(client, template.Func, partialTemplate.Func)(PageVars{}, w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	template.AssertNotCalled(t, "Func")
	mock.AssertExpectationsForObjects(t, client, partialTemplate)
}

func TestPostCreateDonorWhenAPIFails(t *testing.T) {
	client := &mockCreateDonorClient{}
	client.
		On("CreateDonor", mock.Anything, sirius.Person{
			Firstname: "Rudolph",
			Surname:   "Stotesbury",
		}).
		Return(sirius.Person{}, errExample)

	template := &mockTemplate{}

	form := url.Values{
		"firstname": {"Rudolph"},
		"surname":   {"Stotesbury"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/create-donor", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateDonor(client, template.Func, nil)(PageVars{}, w, r)
	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateDonorWhenValidationError(t *testing.T) {
	client := &mockCreateDonorClient{}
	client.
		On("CreateDonor", mock.Anything, sirius.Person{
			Firstname: "Rudolph",
		}).
		Return(sirius.Person{}, sirius.ValidationError{Field: sirius.FieldErrors{
			"surname": {"required": "This field is required"},
		}})

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, donorData{
			IsNew: true,
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

	err := CreateDonor(client, template.Func, nil)(PageVars{}, w, r)
	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}
