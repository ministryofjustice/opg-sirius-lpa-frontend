package server

import (
	"errors"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type mockUnlinkPerson struct {
	mock.Mock
}

func (m *mockUnlinkPerson) Person(ctx sirius.Context, id int) (sirius.Person, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Person), args.Error(1)
}

func (m *mockUnlinkPerson) UnlinkPerson(ctx sirius.Context, parentId int, childId int) error {
	return m.Called(ctx, parentId, childId).Error(0)
}

func TestUnlinkPerson(t *testing.T) {
	children := []sirius.Person{
		{
			ID:         5,
			Salutation: "Mr",
			Firstname:  "First",
			Surname:    "Child",
			Children:   nil,
		},
	}

	person := sirius.Person{
		ID:           189,
		UID:          "700000001234",
		Firstname:    "John",
		Surname:      "Doe",
		DateOfBirth:  "1998-09-02",
		AddressLine1: "123 Somewhere Street",
		Children:     children,
	}

	client := &mockUnlinkPerson{}
	client.On("UnlinkPerson", mock.Anything, 189, 5).Return(nil)

	client.
		On("Person", mock.Anything, 189).
		Return(person, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, unlinkPersonData{
			Entity: person,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=189", nil)
	w := httptest.NewRecorder()

	err := UnlinkPerson(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUnlinkPersonNoID(t *testing.T) {
	person := sirius.Person{ID: 189, Firstname: "John", Surname: "Doe"}

	client := &mockUnlinkPerson{}
	client.
		On("Person", mock.Anything, 189).
		Return(person, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, unlinkPersonData{
			Entity: person,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	err := UnlinkPerson(client, template.Func)(w, r)

	assert.NotNil(t, err)
}

func TestUnlinkPersonsWhenFailure(t *testing.T) {
	expectedError := errors.New("err")

	client := &mockUnlinkPerson{}
	client.
		On("Person", mock.Anything, 189).
		Return(sirius.Person{}, expectedError)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, unlinkPersonData{}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=189", nil)
	w := httptest.NewRecorder()

	err := UnlinkPerson(client, template.Func)(w, r)

	assert.Equal(t, expectedError, err)
}

func TestUnlinkPersonWhenTemplateErrors(t *testing.T) {
	person := sirius.Person{ID: 189, Firstname: "John", Surname: "Doe"}
	expectedError := errors.New("err")

	client := &mockUnlinkPerson{}
	client.
		On("Person", mock.Anything, 189).
		Return(person, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, unlinkPersonData{
			Entity: person,
		}).
		Return(expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=189", nil)
	w := httptest.NewRecorder()

	err := UnlinkPerson(client, template.Func)(w, r)

	assert.Equal(t, expectedError, err)
}

func TestUnlinkPersonWhenValidationError(t *testing.T) {
	person := sirius.Person{ID: 189, Firstname: "John", Surname: "Doe"}
	validationError := sirius.ValidationError{
		Field: sirius.FieldErrors{
			"child": {"reason": "Please select the record to be unlinked"},
		},
	}

	client := &mockUnlinkPerson{}
	client.
		On("Person", mock.Anything, 189).
		Return(person, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, unlinkPersonData{
			Entity: person,
			Error:  validationError,
		}).
		Return(nil)

	form := url.Values{
		"childIds": {"5"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=189", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	err := UnlinkPerson(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestPostUnlinkPerson(t *testing.T) {
	children := []sirius.Person{
		{
			ID:         5,
			Salutation: "Mr",
			Firstname:  "First",
			Surname:    "Child",
			Children:   nil,
		},
	}

	parent := sirius.Person{
		ID:         189,
		Salutation: "Mrs",
		Firstname:  "Parent",
		Surname:    "Person",
		Children:   children,
	}

	client := &mockUnlinkPerson{}
	client.
		On("Person", mock.Anything, 189).
		Return(parent, nil)

	client.
		On("UnlinkPerson", mock.Anything, 189, 5).
		Return(nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, unlinkPersonData{
			Entity:  parent,
			Success: true,
		}).
		Return(nil)

	form := url.Values{
		"child-id": {"5"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=189", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	err := UnlinkPerson(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
