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
	client.
		On("Person", mock.Anything, 189).
		Return(person, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, unlinkPersonData{
			Person: person,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=189", nil)
	w := httptest.NewRecorder()

	err := UnlinkPerson(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestUnlinkPersonNoID(t *testing.T) {
	testCases := map[string]string{
		"no-id":  "/?id=123",
		"bad-id": "/?id=test",
	}

	for name, testUrl := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, testUrl, nil)
			w := httptest.NewRecorder()

			err := EditDates(nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestUnlinkPersonsWhenFailure(t *testing.T) {
	expectedError := errors.New("err")

	client := &mockUnlinkPerson{}
	client.
		On("Person", mock.Anything, 189).
		Return(sirius.Person{}, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=189", nil)
	w := httptest.NewRecorder()

	err := UnlinkPerson(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
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
			Person: person,
		}).
		Return(expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=189", nil)
	w := httptest.NewRecorder()

	err := UnlinkPerson(client, template.Func)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostUnlinkPersonWhenChildNotSelected(t *testing.T) {
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
			Person:  person,
			Error:   validationError,
			Success: false,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodPost, "/?id=189", nil)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	err := UnlinkPerson(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestUnlinkPersonWhenValidationError(t *testing.T) {
	parent := sirius.Person{
		ID:         189,
		Salutation: "Mrs",
		Firstname:  "Parent",
		Surname:    "Person",
	}

	validationError := sirius.ValidationError{Detail: "Cannot unlink records"}

	client := &mockUnlinkPerson{}
	client.
		On("Person", mock.Anything, 189).
		Return(parent, nil)

	client.
		On("UnlinkPerson", mock.Anything, 189, 5).
		Return(validationError)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, unlinkPersonData{
			Person:  parent,
			Error:   validationError,
			Success: false,
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
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
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
			Person:  parent,
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
	mock.AssertExpectationsForObjects(t, client, template)
}
