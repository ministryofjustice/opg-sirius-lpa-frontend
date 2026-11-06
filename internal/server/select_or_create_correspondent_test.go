package server

import (
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

type mockSelectOrCreateCorrespondentClient struct {
	mock.Mock
}

func (m *mockSelectOrCreateCorrespondentClient) CreateCorrespondent(ctx sirius.Context, caseId int, correspondent sirius.Correspondent) error {
	args := m.Called(ctx, caseId, correspondent)
	return args.Error(0)
}

func (m *mockSelectOrCreateCorrespondentClient) Epa(ctx sirius.Context, id int) (sirius.Epa, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Epa), args.Error(1)
}

func (m *mockSelectOrCreateCorrespondentClient) Lpa(ctx sirius.Context, id int) (sirius.Lpa, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Lpa), args.Error(1)
}

func TestGetSelectOrCreateCorrespondent(t *testing.T) {
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
					client := &mockSelectOrCreateCorrespondentClient{}
					data := selectOrCreateCorrespondentData{
						FlowQuery: tc.flowQuery,
						DonorId:   1,
						CaseId:    2,
						CaseType:  tc.caseType,
						IsPartial: isHtmx,
					}
					if tc.caseType == "epa" {
						data.Epa = sirius.Epa{Case: sirius.Case{ID: 2}}
						client.
							On("Epa", mock.Anything, 2).
							Return(data.Epa, nil)
					} else {
						data.Lpa = sirius.Lpa{Case: sirius.Case{ID: 2}}
						client.
							On("Lpa", mock.Anything, 2).
							Return(data.Lpa, nil)
					}

					template := &mockTemplate{}
					template.
						On("Func", mock.Anything, data).
						Return(nil)

					r, _ := http.NewRequest(http.MethodGet, "/?id=1&caseId=2&caseType="+tc.caseType+tc.flowQuery, nil)
					if isHtmx {
						r.Header.Add("HX-Request", "true")
					}
					w := httptest.NewRecorder()

					err := SelectOrCreateCorrespondent(client, template.Func)(w, r)
					resp := w.Result()

					assert.Nil(t, err)
					assert.Equal(t, http.StatusOK, resp.StatusCode)
					mock.AssertExpectationsForObjects(t, client, template)
				})
			}
		})
	}
}

func TestGetSelectOrCreateCorrespondentBadQuery(t *testing.T) {
	testCases := map[string]string{
		"no-id":       "/",
		"bad-id":      "/?id=test",
		"bad-case-id": "/?id=123&caseId=test",
	}

	for name, query := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, query, nil)
			w := httptest.NewRecorder()

			err := SelectOrCreateCorrespondent(nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestGetSelectOrCreateCorrespondentWhenEpaErrors(t *testing.T) {
	client := &mockSelectOrCreateCorrespondentClient{}
	client.
		On("Epa", mock.Anything, 2).
		Return(sirius.Epa{}, errExample)

	r, _ := http.NewRequest(http.MethodGet, "/?id=1&caseId=2&caseType=epa", nil)
	w := httptest.NewRecorder()

	err := SelectOrCreateCorrespondent(client, nil)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetSelectOrCreateCorrespondentWhenLpaErrors(t *testing.T) {
	client := &mockSelectOrCreateCorrespondentClient{}
	client.
		On("Lpa", mock.Anything, 2).
		Return(sirius.Lpa{}, errExample)

	r, _ := http.NewRequest(http.MethodGet, "/?id=1&caseId=2&caseType=lpa", nil)
	w := httptest.NewRecorder()

	err := SelectOrCreateCorrespondent(client, nil)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestPostSelectOrCreateCorrespondentNew(t *testing.T) {
	cases := []struct {
		name      string
		caseType  string
		flowQuery string
		redirect  string
	}{
		{name: "EPA", caseType: "epa", redirect: "/create-correspondent?id=1&caseId=2&caseType=epa"},
		{name: "LPA without flow", caseType: "lpa", redirect: "/create-lpa?id=1&caseId=2#accordion-create-lpa-heading-4"},
		{name: "LPA create flow", caseType: "lpa", flowQuery: "&flow=create", redirect: "/create-lpa?id=1&caseId=2&flow=create#accordion-create-lpa-heading-4"},
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

			client := &mockSelectOrCreateCorrespondentClient{}

			if tc.caseType == "epa" {
				client.
					On("Epa", mock.Anything, 2).
					Return(sirius.Epa{}, nil)
			} else {
				client.
					On("CreateCorrespondent", mock.Anything, 2, correspondent).
					Return(nil)
				client.
					On("Lpa", mock.Anything, 2).
					Return(sirius.Lpa{}, nil)
			}

			form := url.Values{
				"actorId":           {"new"},
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

			r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2&caseType="+tc.caseType+tc.flowQuery, strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := SelectOrCreateCorrespondent(client, nil)(w, r)
			resp := w.Result()

			assert.Equal(t, RedirectError(tc.redirect), err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client)
		})
	}
}

func TestPostSelectOrCreateCorrespondentFromAttorney(t *testing.T) {
	correspondent := sirius.Correspondent{Person: sirius.Person{Firstname: "Rudolph", Surname: "Stotesbury"}}
	epa := sirius.Epa{
		Case: sirius.Case{
			Attorneys: []sirius.Attorney{
				{Person: sirius.Person{ID: 4, Firstname: "Rudolph", Surname: "Stotesbury"}},
			},
		},
	}

	client := &mockSelectOrCreateCorrespondentClient{}
	client.
		On("CreateCorrespondent", mock.Anything, 2, correspondent).
		Return(nil)
	client.
		On("Epa", mock.Anything, 2).
		Return(epa, nil)

	form := url.Values{
		"actorId": {"4"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2&caseType=epa", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := SelectOrCreateCorrespondent(client, nil)(w, r)
	resp := w.Result()

	assert.Equal(t, RedirectError("/create-epa?id=1&caseId=2#accordion-create-epa-heading-3"), err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client)
}

func TestPostSelectOrCreateCorrespondentBadActorId(t *testing.T) {
	expectedErr := sirius.StatusError{Code: http.StatusBadRequest}
	client := &mockSelectOrCreateCorrespondentClient{}
	client.
		On("Epa", mock.Anything, 2).
		Return(sirius.Epa{}, nil)

	template := &mockTemplate{}

	form := url.Values{
		"actorId": {"bad"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2&caseType=epa", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := SelectOrCreateCorrespondent(client, template.Func)(w, r)

	assert.Equal(t, err, expectedErr)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostSelectOrCreateCorrespondentValidationError(t *testing.T) {
	expectedError := sirius.ValidationError{
		Field: sirius.FieldErrors{"field": {"": "problem"}},
	}

	epa := sirius.Epa{
		Case: sirius.Case{
			Attorneys: []sirius.Attorney{
				{Person: sirius.Person{ID: 4, Firstname: "Rudolph", Surname: "Stotesbury"}},
			},
		},
	}
	correspondent := sirius.Correspondent{Person: sirius.Person{Firstname: "Rudolph", Surname: "Stotesbury"}}

	client := &mockSelectOrCreateCorrespondentClient{}
	client.
		On("Epa", mock.Anything, 2).
		Return(epa, nil).
		On("CreateCorrespondent", mock.Anything, 2, correspondent).
		Return(expectedError)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, selectOrCreateCorrespondentData{
			DonorId:       1,
			CaseId:        2,
			CaseType:      "epa",
			Epa:           epa,
			Error:         expectedError,
			Correspondent: correspondent,
		}).
		Return(nil)

	form := url.Values{
		"actorId": {"4"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2&caseType=epa", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := SelectOrCreateCorrespondent(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostSelectOrCreateCorrespondentValidationErrorHtmxRequest(t *testing.T) {
	expectedError := sirius.ValidationError{
		Field: sirius.FieldErrors{"field": {"": "problem"}},
	}

	epa := sirius.Epa{
		Case: sirius.Case{
			Attorneys: []sirius.Attorney{
				{Person: sirius.Person{ID: 4, Firstname: "Rudolph", Surname: "Stotesbury"}},
			},
		},
	}
	correspondent := sirius.Correspondent{Person: sirius.Person{Firstname: "Rudolph", Surname: "Stotesbury"}}

	client := &mockSelectOrCreateCorrespondentClient{}
	client.
		On("Epa", mock.Anything, 2).
		Return(epa, nil).
		On("CreateCorrespondent", mock.Anything, 2, correspondent).
		Return(expectedError)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, selectOrCreateCorrespondentData{
			DonorId:       1,
			CaseId:        2,
			CaseType:      "epa",
			Epa:           epa,
			Error:         expectedError,
			IsPartial:     true,
			Correspondent: correspondent,
		}).
		Return(nil)

	form := url.Values{
		"actorId": {"4"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2&caseType=epa", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	r.Header.Add("HX-Request", "true")
	w := httptest.NewRecorder()

	err := SelectOrCreateCorrespondent(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostSelectOrCreateCorrespondentCreationFails(t *testing.T) {
	epa := sirius.Epa{
		Case: sirius.Case{
			Attorneys: []sirius.Attorney{
				{Person: sirius.Person{ID: 4, Firstname: "Rudolph", Surname: "Stotesbury"}},
			},
		},
	}
	correspondent := sirius.Correspondent{Person: sirius.Person{Firstname: "Rudolph", Surname: "Stotesbury"}}

	client := &mockSelectOrCreateCorrespondentClient{}
	client.
		On("Epa", mock.Anything, 2).
		Return(epa, nil).
		On("CreateCorrespondent", mock.Anything, 2, correspondent).
		Return(errExample)

	template := &mockTemplate{}

	form := url.Values{
		"actorId": {"4"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2&caseType=epa", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := SelectOrCreateCorrespondent(client, template.Func)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client, template)
}
