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

type mockLinkPersonClient struct {
	mock.Mock
}

func (m *mockLinkPersonClient) LinkPeople(ctx sirius.Context, parentId int, childId int) error {
	args := m.Called(ctx, parentId, childId)
	return args.Error(0)
}

func (m *mockLinkPersonClient) Person(ctx sirius.Context, id int) (sirius.Person, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Person), args.Error(1)
}

func (m *mockLinkPersonClient) PersonByUid(ctx sirius.Context, uid string) (sirius.Person, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(sirius.Person), args.Error(1)
}

func TestLinkPerson(t *testing.T) {
	person := sirius.Person{Firstname: "John", Surname: "Doe"}

	client := &mockLinkPersonClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(person, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, linkPersonData{
			Entity: person,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := LinkPerson(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestLinkPersonNoID(t *testing.T) {
	person := sirius.Person{Firstname: "John", Surname: "Doe"}

	client := &mockLinkPersonClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(person, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, linkPersonData{
			Entity: person,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	err := LinkPerson(client, template.Func)(w, r)

	assert.NotNil(t, err)
}

func TestLinkPersonGetFails(t *testing.T) {
	expectedError := errors.New("err")

	client := &mockLinkPersonClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, expectedError)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, linkPersonData{}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := LinkPerson(client, template.Func)(w, r)

	assert.Equal(t, expectedError, err)
}

func TestLinkPersonTemplateErrors(t *testing.T) {
	person := sirius.Person{Firstname: "John", Surname: "Doe"}
	expectedError := errors.New("err")

	client := &mockLinkPersonClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(person, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, linkPersonData{
			Entity: person,
		}).
		Return(expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := LinkPerson(client, template.Func)(w, r)

	assert.Equal(t, expectedError, err)
}

func TestLinkPersonSearch(t *testing.T) {
	testCases := map[string]struct {
		personHasChildren      bool
		otherPersonHasChildren bool
		expectedPrimary        int
		expectedCanChange      bool
	}{
		"no children": {
			personHasChildren:      false,
			otherPersonHasChildren: false,
			expectedPrimary:        0,
			expectedCanChange:      true,
		},
		"person has children": {
			personHasChildren:      true,
			otherPersonHasChildren: false,
			expectedPrimary:        1,
			expectedCanChange:      false,
		},
		"other person has children": {
			personHasChildren:      false,
			otherPersonHasChildren: true,
			expectedPrimary:        2,
			expectedCanChange:      false,
		},
		"both have children": {
			personHasChildren:      true,
			otherPersonHasChildren: true,
			expectedPrimary:        0,
			expectedCanChange:      true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			person := sirius.Person{ID: 1, Firstname: "John", Surname: "Doe"}
			otherPerson := sirius.Person{ID: 2, Firstname: "Murray", Surname: "Stremlow"}

			if tc.personHasChildren {
				person.Children = []sirius.Person{{}}
			}

			if tc.otherPersonHasChildren {
				otherPerson.Children = []sirius.Person{{}}
			}

			client := &mockLinkPersonClient{}
			client.
				On("Person", mock.Anything, 123).
				Return(person, nil)
			client.
				On("PersonByUid", mock.Anything, "7000-0000-8293").
				Return(otherPerson, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, linkPersonData{
					Entity:           person,
					OtherPerson:      otherPerson,
					PrimaryId:        tc.expectedPrimary,
					CanChangePrimary: tc.expectedCanChange,
				}).
				Return(nil)

			form := url.Values{
				"uid": {"7000-0000-8293"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()

			err := LinkPerson(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	}
}

func TestLinkPersonSearchNotFound(t *testing.T) {
	person := sirius.Person{Firstname: "John", Surname: "Doe"}

	client := &mockLinkPersonClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(person, nil)
	client.
		On("PersonByUid", mock.Anything, "7000-0000-8293").
		Return(sirius.Person{}, sirius.StatusError{Code: 404})

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, linkPersonData{
			Entity: person,
			Error: sirius.ValidationError{
				Field: sirius.FieldErrors{
					"uid": map[string]string{
						"notFound": "A record matching the supplied uId cannot be found.",
					},
				},
			},
		}).
		Return(nil)

	form := url.Values{
		"uid": {"7000-0000-8293"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	err := LinkPerson(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestLinkPersonSave(t *testing.T) {
	person := sirius.Person{ID: 1, Firstname: "John", Surname: "Doe"}
	otherPerson := sirius.Person{ID: 2, Firstname: "Murray", Surname: "Stremlow"}

	client := &mockLinkPersonClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(person, nil)
	client.
		On("PersonByUid", mock.Anything, "7000-0000-8293").
		Return(otherPerson, nil)
	client.
		On("LinkPeople", mock.Anything, 1, 2).
		Return(nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, linkPersonData{
			Entity:           person,
			OtherPerson:      otherPerson,
			PrimaryId:        1,
			CanChangePrimary: true,
			Success:          true,
		}).
		Return(nil)

	form := url.Values{
		"uid":        {"7000-0000-8293"},
		"primary-id": {"1"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	err := LinkPerson(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestLinkPersonSaveOtherPrimary(t *testing.T) {
	person := sirius.Person{ID: 1, Firstname: "John", Surname: "Doe"}
	otherPerson := sirius.Person{ID: 2, Firstname: "Murray", Surname: "Stremlow"}

	client := &mockLinkPersonClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(person, nil)
	client.
		On("PersonByUid", mock.Anything, "7000-0000-8293").
		Return(otherPerson, nil)
	client.
		On("LinkPeople", mock.Anything, 2, 1).
		Return(nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, linkPersonData{
			Entity:           person,
			OtherPerson:      otherPerson,
			PrimaryId:        2,
			CanChangePrimary: true,
			Success:          true,
		}).
		Return(nil)

	form := url.Values{
		"uid":        {"7000-0000-8293"},
		"primary-id": {"2"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	err := LinkPerson(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestLinkPersonSaveValidationError(t *testing.T) {
	person := sirius.Person{ID: 1, Firstname: "John", Surname: "Doe"}
	otherPerson := sirius.Person{ID: 2, Firstname: "Murray", Surname: "Stremlow"}
	saveError := sirius.ValidationError{Detail: "Cannot link records"}

	client := &mockLinkPersonClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(person, nil)
	client.
		On("PersonByUid", mock.Anything, "7000-0000-8293").
		Return(otherPerson, nil)
	client.
		On("LinkPeople", mock.Anything, 1, 2).
		Return(saveError)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, linkPersonData{
			Entity:           person,
			OtherPerson:      otherPerson,
			PrimaryId:        1,
			CanChangePrimary: true,
			Error:            saveError,
		}).
		Return(nil)

	form := url.Values{
		"uid":        {"7000-0000-8293"},
		"primary-id": {"1"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	err := LinkPerson(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
