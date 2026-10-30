package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUpdateCorrespondentClient struct {
	mock.Mock
}

func (m *mockUpdateCorrespondentClient) Epa(ctx sirius.Context, id int) (sirius.Epa, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Epa), args.Error(1)
}

func (m *mockUpdateCorrespondentClient) Lpa(ctx sirius.Context, id int) (sirius.Lpa, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Lpa), args.Error(1)
}

func (m *mockUpdateCorrespondentClient) UpdateCorrespondent(ctx sirius.Context, correspondentId int, correspondent sirius.Correspondent) error {
	args := m.Called(ctx, correspondentId, correspondent)
	return args.Error(0)
}

func TestGetUpdateCorrespondent(t *testing.T) {
	cases := []struct {
		name      string
		caseType  string
		flowQuery string
	}{
		{name: "EPA", caseType: "epa"},
		{name: "LPA without flow", caseType: "lpa"},
		{name: "LPA create flow", caseType: "lpa", flowQuery: "&flow=create"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			existingCorrespondent := sirius.Correspondent{
				Person: sirius.Person{
					ID:        7,
					Firstname: "Rudolph",
					Surname:   "Stotesbury",
				},
			}

			client := &mockUpdateCorrespondentClient{}
			if tc.caseType == "epa" {
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
				On("Func", mock.Anything, correspondentData{
					FlowQuery:     tc.flowQuery,
					DonorId:       1,
					CaseId:        2,
					CaseType:      tc.caseType,
					Correspondent: existingCorrespondent,
					IsEditing:     true,
					Title:         "Update correspondent details",
				}).
				Return(nil)

			r, _ := http.NewRequest(http.MethodGet, "/edit-correspondent?id=1&caseId=2&caseType="+tc.caseType+tc.flowQuery, nil)
			w := httptest.NewRecorder()

			err := UpdateCorrespondent(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestGetUpdateCorrespondentWhenFetchCaseFails(t *testing.T) {
	for _, caseType := range []string{"epa", "lpa"} {
		t.Run(caseType, func(t *testing.T) {
			client := &mockUpdateCorrespondentClient{}
			if caseType == "epa" {
				client.
					On("Epa", mock.Anything, 2).
					Return(sirius.Epa{}, errExample)
			} else {
				client.
					On("Lpa", mock.Anything, 2).
					Return(sirius.Lpa{}, errExample)
			}

			r, _ := http.NewRequest(http.MethodGet, "/edit-correspondent?id=1&caseId=2&caseType="+caseType, nil)
			w := httptest.NewRecorder()

			err := UpdateCorrespondent(client, nil)(w, r)

			assert.Equal(t, errExample, err)
			mock.AssertExpectationsForObjects(t, client)
		})
	}
}

func TestGetUpdateCorrespondentBadQuery(t *testing.T) {
	testCases := map[string]string{
		"no-id":       "/",
		"bad-id":      "/?id=test",
		"bad-case-id": "/?id=123&caseId=test",
	}

	for name, query := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, query, nil)
			w := httptest.NewRecorder()

			err := UpdateCorrespondent(nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestPostUpdateCorrespondent(t *testing.T) {
	cases := []struct {
		name      string
		caseType  string
		flowQuery string
		redirect  string
	}{
		{
			name:     "EPA",
			caseType: "epa",
			redirect: "/create-epa?id=1&caseId=2#accordion-create-epa-heading-3",
		},
		{
			name:     "LPA without flow",
			caseType: "lpa",
			redirect: "/create-lpa?id=1&caseId=2#accordion-create-lpa-heading-4",
		},
		{
			name:      "LPA create flow",
			caseType:  "lpa",
			flowQuery: "&flow=create",
			redirect:  "/create-lpa?id=1&caseId=2&flow=create#accordion-create-lpa-heading-4",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
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

			client := &mockUpdateCorrespondentClient{}
			if tc.caseType == "epa" {
				client.
					On("Epa", mock.Anything, 2).
					Return(sirius.Epa{Case: sirius.Case{Correspondent: &existingCorrespondent}}, nil)
			} else {
				client.
					On("Lpa", mock.Anything, 2).
					Return(sirius.Lpa{Case: sirius.Case{Correspondent: &existingCorrespondent}}, nil)
			}
			client.
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

			r, _ := http.NewRequest(http.MethodPost, "/edit-correspondent?id=1&caseId=2&caseType="+tc.caseType+tc.flowQuery, strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := UpdateCorrespondent(client, nil)(w, r)
			resp := w.Result()

			assert.Equal(t, RedirectError(tc.redirect), err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client)
		})
	}
}

func TestPostUpdateCorrespondentValidationErrors(t *testing.T) {
	cases := []struct {
		name        string
		caseType    string
		expectedErr sirius.ValidationError
	}{
		{
			name:     "EPA validation errors",
			caseType: "epa",
			expectedErr: sirius.ValidationError{
				Field: sirius.FieldErrors{"field": {"": "problem"}},
			},
		},
		{
			name:     "LPA validation errors",
			caseType: "lpa",
			expectedErr: sirius.ValidationError{
				Field: sirius.FieldErrors{"field": {"": "problem"}},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
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

			client := &mockUpdateCorrespondentClient{}
			if tc.caseType == "epa" {
				client.
					On("Epa", mock.Anything, 2).
					Return(sirius.Epa{Case: sirius.Case{Correspondent: &existingCorrespondent}}, nil)
			} else {
				client.
					On("Lpa", mock.Anything, 2).
					Return(sirius.Lpa{Case: sirius.Case{Correspondent: &existingCorrespondent}}, nil)
			}
			client.
				On("UpdateCorrespondent", mock.Anything, 7, updatedCorrespondent).
				Return(tc.expectedErr)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, correspondentData{
					DonorId:       1,
					CaseId:        2,
					CaseType:      tc.caseType,
					Error:         tc.expectedErr,
					Correspondent: updatedCorrespondent,
					Title:         "Update correspondent details",
					IsEditing:     true,
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

			r, _ := http.NewRequest(http.MethodPost, "/edit-correspondent?id=1&caseId=2&caseType="+tc.caseType, strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := UpdateCorrespondent(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostUpdateCorrespondentApiErrors(t *testing.T) {
	cases := []struct {
		name        string
		caseType    string
		expectedErr error
	}{
		{
			name:        "EPA errors",
			caseType:    "epa",
			expectedErr: errExample,
		},
		{
			name:        "LPA errors",
			caseType:    "lpa",
			expectedErr: errExample,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
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

			client := &mockUpdateCorrespondentClient{}
			if tc.caseType == "epa" {
				client.
					On("Epa", mock.Anything, 2).
					Return(sirius.Epa{Case: sirius.Case{Correspondent: &existingCorrespondent}}, nil)
			} else {
				client.
					On("Lpa", mock.Anything, 2).
					Return(sirius.Lpa{Case: sirius.Case{Correspondent: &existingCorrespondent}}, nil)
			}
			client.
				On("UpdateCorrespondent", mock.Anything, 7, updatedCorrespondent).
				Return(tc.expectedErr)

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

			r, _ := http.NewRequest(http.MethodPost, "/edit-correspondent?id=1&caseId=2&caseType="+tc.caseType, strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := UpdateCorrespondent(client, nil)(w, r)

			assert.Equal(t, tc.expectedErr, err)
			mock.AssertExpectationsForObjects(t, client)
		})
	}
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
