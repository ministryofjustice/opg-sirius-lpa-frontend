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
	epa := sirius.Epa{Case: sirius.Case{ID: 2}}

	client := &mockSelectOrCreateCorrespondentClient{}
	client.
		On("Epa", mock.Anything, 2).
		Return(epa, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, selectOrCreateCorrespondentData{
			DonorId:  1,
			CaseId:   2,
			CaseType: "epa",
			Epa:      epa,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=1&caseId=2&caseType=epa", nil)
	w := httptest.NewRecorder()

	err := SelectOrCreateCorrespondent(client, template.Func, nil)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetSelectOrCreateCorrespondentHtmxRequest(t *testing.T) {
	epa := sirius.Epa{Case: sirius.Case{ID: 2}}

	client := &mockSelectOrCreateCorrespondentClient{}
	client.
		On("Epa", mock.Anything, 2).
		Return(epa, nil)

	template := &mockTemplate{}
	partialTemplate := &mockTemplate{}
	partialTemplate.
		On("Func", mock.Anything, selectOrCreateCorrespondentData{
			DonorId:  1,
			CaseId:   2,
			CaseType: "epa",
			Epa:      epa,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=1&caseId=2&caseType=epa", nil)
	r.Header.Add("HX-Request", "true")
	w := httptest.NewRecorder()

	err := SelectOrCreateCorrespondent(client, template.Func, partialTemplate.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	template.AssertNotCalled(t, "Func", mock.Anything, mock.Anything)
	mock.AssertExpectationsForObjects(t, client, template, partialTemplate)
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

			err := SelectOrCreateCorrespondent(nil, nil, nil)(w, r)

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

	err := SelectOrCreateCorrespondent(client, nil, nil)(w, r)

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

	err := SelectOrCreateCorrespondent(client, nil, nil)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestPostSelectOrCreateCorrespondentNew(t *testing.T) {
	for _, caseType := range []string{"lpa", "epa"} {
		t.Run(caseType, func(t *testing.T) {
			expectedError := RedirectError("/create-correspondent?id=1&caseId=2&caseType=" + caseType)

			client := &mockSelectOrCreateCorrespondentClient{}

			if caseType == "epa" {
				client.
					On("Epa", mock.Anything, 2).
					Return(sirius.Epa{}, nil)
			} else {
				client.
					On("Lpa", mock.Anything, 2).
					Return(sirius.Lpa{}, nil)
			}

			form := url.Values{
				"actorId": {"new"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2&caseType="+caseType, strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := SelectOrCreateCorrespondent(client, nil, nil)(w, r)
			resp := w.Result()

			assert.Equal(t, err, expectedError)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client)
		})
	}
}

func TestPostSelectOrCreateCorrespondentFromAttorney(t *testing.T) {
	for _, caseType := range []string{"lpa", "epa"} {
		t.Run(caseType, func(t *testing.T) {
			var expectedError error
			if caseType == "epa" {
				expectedError = RedirectError("/create-epa?id=1&caseId=2#accordion-create-epa-heading-3")
			} else {
				expectedError = RedirectError("/create-lpa?id=1&caseId=2#accordion-create-lpa-heading-4")
			}

			correspondent := sirius.Correspondent{Person: sirius.Person{Firstname: "Rudolph", Surname: "Stotesbury"}}

			client := &mockSelectOrCreateCorrespondentClient{}
			client.
				On("CreateCorrespondent", mock.Anything, 2, correspondent).
				Return(nil)

			if caseType == "epa" {
				epa := sirius.Epa{
					Case: sirius.Case{
						Attorneys: []sirius.Attorney{
							{Person: sirius.Person{ID: 4, Firstname: "Rudolph", Surname: "Stotesbury"}},
						},
					},
				}

				client.
					On("Epa", mock.Anything, 2).
					Return(epa, nil)
			} else {
				lpa := sirius.Lpa{
					Case: sirius.Case{
						Attorneys: []sirius.Attorney{
							{Person: sirius.Person{ID: 4, Firstname: "Rudolph", Surname: "Stotesbury"}},
						},
						Donor: &sirius.Person{ID: 876},
					},
				}

				client.
					On("Lpa", mock.Anything, 2).
					Return(lpa, nil)
			}

			form := url.Values{
				"actorId": {"4"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2&caseType="+caseType, strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := SelectOrCreateCorrespondent(client, nil, nil)(w, r)
			resp := w.Result()

			assert.Equal(t, err, expectedError)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client)
		})
	}
}

func TestPostSelectOrCreateCorrespondentBadactorId(t *testing.T) {
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

	err := SelectOrCreateCorrespondent(client, template.Func, nil)(w, r)

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
			DonorId:  1,
			CaseId:   2,
			CaseType: "epa",
			Epa:      epa,
			Error:    expectedError,
		}).
		Return(nil)

	form := url.Values{
		"actorId": {"4"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2&caseType=epa", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := SelectOrCreateCorrespondent(client, template.Func, nil)(w, r)
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
	partialTemplate := &mockTemplate{}
	partialTemplate.
		On("Func", mock.Anything, selectOrCreateCorrespondentData{
			DonorId:  1,
			CaseId:   2,
			CaseType: "epa",
			Epa:      epa,
			Error:    expectedError,
		}).
		Return(nil)

	form := url.Values{
		"actorId": {"4"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2&caseType=epa", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	r.Header.Add("HX-Request", "true")
	w := httptest.NewRecorder()

	err := SelectOrCreateCorrespondent(client, template.Func, partialTemplate.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	template.AssertNotCalled(t, "Func", mock.Anything, mock.Anything)
	mock.AssertExpectationsForObjects(t, client, template, partialTemplate)
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

	err := SelectOrCreateCorrespondent(client, template.Func, nil)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetSelectedActorForLpa(t *testing.T) {
	lpa := sirius.Lpa{
		Case: sirius.Case{
			Donor: &sirius.Person{ID: 1},
			Attorneys: []sirius.Attorney{
				{Person: sirius.Person{ID: 2}},
			},
			ReplacementAttorneys: []sirius.Attorney{
				{Person: sirius.Person{ID: 3}},
			},
			CertificateProviders: []sirius.Person{
				{ID: 4},
			},
			NotifiedPersons: []sirius.Person{
				{ID: 5},
			},
			TrustCorporations: []sirius.Attorney{
				{Person: sirius.Person{ID: 6}},
			},
		},
	}

	testCases := []struct {
		actorId        int
		expectedPerson sirius.Person
	}{
		{actorId: 1, expectedPerson: *lpa.Donor},
		{actorId: 2, expectedPerson: lpa.Attorneys[0].Person},
		{actorId: 3, expectedPerson: lpa.ReplacementAttorneys[0].Person},
		{actorId: 4, expectedPerson: lpa.CertificateProviders[0]},
		{actorId: 5, expectedPerson: lpa.NotifiedPersons[0]},
		{actorId: 6, expectedPerson: lpa.TrustCorporations[0].Person},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expectedPerson, GetSelectedActorForLpa(lpa, tc.actorId))
	}
}
