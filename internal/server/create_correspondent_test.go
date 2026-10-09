package server

import (
	"encoding/json"
	"fmt"
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

func (m *mockCreateCorrespondentClient) Lpa(ctx sirius.Context, id int) (sirius.Lpa, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Lpa), args.Error(1)
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
			DonorId:  1,
			CaseId:   2,
			CaseType: "epa",
			Title:    "Add a correspondent",
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=1&caseId=2&caseType=epa", nil)
	w := httptest.NewRecorder()

	err := CreateCorrespondent(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetCreateCorrespondentHtmxRequest(t *testing.T) {
	client := &mockCreateCorrespondentClient{}
	client.
		On("Epa", mock.Anything, 2).
		Return(sirius.Epa{}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createCorrespondentData{
			IsPartial: true,
			DonorId:   1,
			CaseId:    2,
			CaseType:  "epa",
			Title:     "Add a correspondent",
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=1&caseId=2&caseType=epa", nil)
	r.Header.Add("HX-Request", "true")
	w := httptest.NewRecorder()

	err := CreateCorrespondent(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetEditCorrespondent(t *testing.T) {
	for _, caseType := range []string{"epa", "lpa"} {
		t.Run(caseType, func(t *testing.T) {
			existingCorrespondent := sirius.Correspondent{
				Person: sirius.Person{
					ID:        7,
					Firstname: "Rudolph",
					Surname:   "Stotesbury",
				},
			}

			client := &mockCreateCorrespondentClient{}
			if caseType == "epa" {
				client.
					On("Epa", mock.Anything, 2).
					Return(sirius.Epa{Case: sirius.Case{Correspondent: &existingCorrespondent}}, nil)
			} else {
				client.
					On("Lpa", mock.Anything, 2).
					Return(sirius.Lpa{Case: sirius.Case{Correspondent: &existingCorrespondent}}, nil)
			}

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, createCorrespondentData{
					DonorId:       1,
					CaseId:        2,
					CaseType:      caseType,
					Correspondent: existingCorrespondent,
					IsEditing:     true,
					Title:         "Update correspondent details",
				}).
				Return(nil)

			r, _ := http.NewRequest(http.MethodGet, "/?id=1&caseId=2&caseType="+caseType, nil)
			w := httptest.NewRecorder()

			err := CreateCorrespondent(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
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

func TestGetCreateCorrespondentWhenFetchCaseFails(t *testing.T) {
	for _, caseType := range []string{"epa", "lpa"} {
		t.Run(caseType, func(t *testing.T) {
			client := &mockCreateCorrespondentClient{}
			if caseType == "epa" {
				client.
					On("Epa", mock.Anything, 2).
					Return(sirius.Epa{}, errExample)
			} else {
				client.
					On("Lpa", mock.Anything, 2).
					Return(sirius.Lpa{}, errExample)
			}

			r, _ := http.NewRequest(http.MethodGet, "/?id=1&caseId=2&caseType="+caseType, nil)
			w := httptest.NewRecorder()

			err := CreateCorrespondent(client, nil)(w, r)

			assert.Equal(t, errExample, err)
			mock.AssertExpectationsForObjects(t, client)
		})
	}
}

func TestPostCreateCorrespondent(t *testing.T) {
	for _, caseType := range []string{"epa", "lpa"} {
		t.Run(caseType, func(t *testing.T) {
			var expectedError error
			if caseType == "epa" {
				expectedError = RedirectError("/create-epa?id=1&caseId=2#accordion-create-epa-heading-3")
			} else {
				expectedError = RedirectError("/create-lpa?id=1&caseId=2#accordion-create-lpa-heading-4")
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

			if caseType == "epa" {
				client.
					On("Epa", mock.Anything, 2).
					Return(sirius.Epa{}, nil)
			} else {
				client.
					On("Lpa", mock.Anything, 2).
					Return(sirius.Lpa{}, nil)
			}

			client.
				On("CreateCorrespondent", mock.Anything, 2, correspondent).
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

			r, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/?id=1&caseId=2&caseType=%s", caseType), strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := CreateCorrespondent(client, nil)(w, r)
			resp := w.Result()

			assert.Equal(t, err, expectedError)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client)
		})
	}
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

	r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2&caseType=epa", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateCorrespondent(client, nil)(w, r)
	resp := w.Result()

	assert.Equal(t, err, expectedError)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client)
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
			CaseType:      "epa",
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

	r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2&caseType=epa", strings.NewReader(form.Encode()))
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

	r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2&caseType=epa", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateCorrespondent(client, nil)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetActorsForLpa(t *testing.T) {
	lpa := sirius.Lpa{
		Case: sirius.Case{
			Donor:                &sirius.Person{ID: 1, Firstname: "John", Surname: "Smith"},
			Attorneys:            []sirius.Attorney{{Person: sirius.Person{ID: 2, Firstname: "Melanie", Surname: "Vanvolkenburg"}}},
			ReplacementAttorneys: []sirius.Attorney{{Person: sirius.Person{ID: 3, Firstname: "John", Surname: "Doe"}}},
			CertificateProviders: []sirius.Person{{ID: 4, Firstname: "Steven", Surname: "Munnell"}},
			NotifiedPersons:      []sirius.NotifiedPerson{{Person: sirius.Person{ID: 5, Firstname: "Katheryn", Surname: "Collins"}}},
			TrustCorporations:    []sirius.TrustCorporation{{Attorney: sirius.Attorney{Person: sirius.Person{ID: 6, CompanyName: "ACME"}}}},
		},
	}
	actors := GetActorsForLpa(lpa)

	type actor struct {
		ID     int
		Name   string
		Label  string
		Fields string
	}

	donorFields, _ := json.Marshal(lpa.Donor)
	attorneyFields, _ := json.Marshal(lpa.Attorneys[0])
	replacementAttorneyFields, _ := json.Marshal(lpa.ReplacementAttorneys[0])
	certificateProviderFields, _ := json.Marshal(lpa.CertificateProviders[0])
	notifiedPersonFields, _ := json.Marshal(lpa.NotifiedPersons[0])
	trustCorporationFields, _ := json.Marshal(lpa.TrustCorporations[0])

	expectedActors := []actor{
		{ID: 1, Name: "John Smith", Label: "donor", Fields: string(donorFields)},
		{ID: 2, Name: "Melanie Vanvolkenburg", Label: "attorney", Fields: string(attorneyFields)},
		{ID: 3, Name: "John Doe", Label: "replacement attorney", Fields: string(replacementAttorneyFields)},
		{ID: 4, Name: "Steven Munnell", Label: "certificate provider", Fields: string(certificateProviderFields)},
		{ID: 5, Name: "Katheryn Collins", Label: "notified person", Fields: string(notifiedPersonFields)},
		{ID: 6, Name: "ACME", Label: "trust corporation", Fields: string(trustCorporationFields)},
	}

	for i, actualActor := range actors {
		assert.Equal(t, expectedActors[i].ID, actualActor.ID)
		assert.Equal(t, expectedActors[i].Name, actualActor.Name)
		assert.Equal(t, expectedActors[i].Label, actualActor.Label)
		assert.JSONEq(t, expectedActors[i].Fields, actualActor.Fields)
	}
}
