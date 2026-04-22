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

func TestGetSelectOrCreateCorrespondent(t *testing.T) {
	epa := sirius.Epa{Case: sirius.Case{ID: 2}}

	client := &mockSelectOrCreateCorrespondentClient{}
	client.
		On("Epa", mock.Anything, 2).
		Return(epa, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, selectOrCreateCorrespondentData{
			DonorId: 1,
			CaseId:  2,
			Epa:     epa,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=1&caseId=2", nil)
	w := httptest.NewRecorder()

	err := SelectOrCreateCorrespondent(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
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

	r, _ := http.NewRequest(http.MethodGet, "/?id=1&caseId=2", nil)
	w := httptest.NewRecorder()

	err := SelectOrCreateCorrespondent(client, nil)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestPostSelectOrCreateCorrespondentNew(t *testing.T) {
	expectedError := RedirectError("/create-correspondent?id=1&caseId=2")

	client := &mockSelectOrCreateCorrespondentClient{}
	client.
		On("Epa", mock.Anything, 2).
		Return(sirius.Epa{}, nil)

	template := &mockTemplate{}

	form := url.Values{
		"attorneyId": {"new"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := SelectOrCreateCorrespondent(client, template.Func)(w, r)
	resp := w.Result()

	assert.Equal(t, err, expectedError)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostSelectOrCreateCorrespondentFromAttorney(t *testing.T) {
	expectedError := RedirectError("/create-epa?id=1&caseId=2")

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
		Return(nil)

	template := &mockTemplate{}

	form := url.Values{
		"attorneyId": {"4"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := SelectOrCreateCorrespondent(client, template.Func)(w, r)
	resp := w.Result()

	assert.Equal(t, err, expectedError)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostSelectOrCreateCorrespondentBadAttorneyId(t *testing.T) {
	expectedErr := sirius.StatusError{Code: http.StatusBadRequest}
	client := &mockSelectOrCreateCorrespondentClient{}
	client.
		On("Epa", mock.Anything, 2).
		Return(sirius.Epa{}, nil)

	template := &mockTemplate{}

	form := url.Values{
		"attorneyId": {"bad"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2", strings.NewReader(form.Encode()))
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
			DonorId: 1,
			CaseId:  2,
			Epa:     epa,
			Error:   expectedError,
		}).
		Return(nil)

	form := url.Values{
		"attorneyId": {"4"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=1&caseId=2", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := SelectOrCreateCorrespondent(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}
