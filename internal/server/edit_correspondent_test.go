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

type mockEditCorrespondentClient struct {
	mock.Mock
}

func (m *mockEditCorrespondentClient) Epa(ctx sirius.Context, id int) (sirius.Epa, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Epa), args.Error(1)
}

func (m *mockEditCorrespondentClient) Lpa(ctx sirius.Context, id int) (sirius.Lpa, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Lpa), args.Error(1)
}

func (m *mockEditCorrespondentClient) UpdateCorrespondent(ctx sirius.Context, correspondentId int, correspondent sirius.Correspondent) error {
	args := m.Called(ctx, correspondentId, correspondent)
	return args.Error(0)
}

func TestGetEditCorrespondent(t *testing.T) {
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

			client := &mockEditCorrespondentClient{}
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
				On("Func", mock.Anything, editCorrespondentData{
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

			err := EditCorrespondent(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestGetEditCorrespondentWhenFetchCaseFails(t *testing.T) {
	for _, caseType := range []string{"epa", "lpa"} {
		t.Run(caseType, func(t *testing.T) {
			client := &mockEditCorrespondentClient{}
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

			err := EditCorrespondent(client, nil)(w, r)

			assert.Equal(t, errExample, err)
			mock.AssertExpectationsForObjects(t, client)
		})
	}
}

func TestPostEditCorrespondent(t *testing.T) {
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

			client := &mockEditCorrespondentClient{}
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

			r, _ := http.NewRequest(http.MethodPost, "edit-correspondent?id=1&caseId=2&caseType="+tc.caseType+tc.flowQuery, strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := EditCorrespondent(client, nil)(w, r)
			resp := w.Result()

			assert.Equal(t, RedirectError(tc.redirect), err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client)
		})
	}
}
