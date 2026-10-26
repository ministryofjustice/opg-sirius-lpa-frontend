package server

import (
	"fmt"
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

type mockCreateCorrespondentClient struct {
	mock.Mock
}

func (m *mockCreateCorrespondentClient) CreateCorrespondent(ctx sirius.Context, caseId int, correspondent sirius.Correspondent) error {
	args := m.Called(ctx, caseId, correspondent)
	return args.Error(0)
}

func TestGetCreateCorrespondent(t *testing.T) {
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
			for _, isHtmx := range []bool{false, true} {
				t.Run("Is Htmx: "+strconv.FormatBool(isHtmx), func(t *testing.T) {
					client := &mockCreateCorrespondentClient{}
					template := &mockTemplate{}
					template.
						On("Func", mock.Anything, createCorrespondentData{
							FlowQuery: tc.flowQuery,
							IsPartial: isHtmx,
							DonorId:   1,
							CaseId:    2,
							CaseType:  tc.caseType,
							Title:     "Add a correspondent",
						}).
						Return(nil)

					r, _ := http.NewRequest(http.MethodGet, "/?id=1&caseId=2&caseType="+tc.caseType+tc.flowQuery, nil)
					if isHtmx {
						r.Header.Add("HX-Request", "true")
					}
					w := httptest.NewRecorder()

					err := CreateCorrespondent(client, template.Func)(w, r)
					resp := w.Result()

					assert.Nil(t, err)
					assert.Equal(t, http.StatusOK, resp.StatusCode)
					mock.AssertExpectationsForObjects(t, client, template)
				})
			}
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

func TestPostCreateCorrespondent(t *testing.T) {
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

			r, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/?id=1&caseId=2&caseType=%s%s", tc.caseType, tc.flowQuery), strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := CreateCorrespondent(client, nil)(w, r)
			resp := w.Result()

			assert.Equal(t, RedirectError(tc.redirect), err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client)
		})
	}
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

	r, _ := http.NewRequest(http.MethodPost, "/create-correspondent?id=1&caseId=2&caseType=epa", strings.NewReader(form.Encode()))
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
