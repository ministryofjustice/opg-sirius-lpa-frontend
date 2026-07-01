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

type mockRelationshipClient struct {
	mock.Mock
}

func (m *mockRelationshipClient) CreatePersonReference(ctx sirius.Context, personID int, referencedUID, reason string) error {
	args := m.Called(ctx, personID, referencedUID, reason)
	return args.Error(0)
}

func (m *mockRelationshipClient) Person(ctx sirius.Context, id int) (sirius.Person, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Person), args.Error(1)
}

func TestGetRelationship(t *testing.T) {
	client := &mockRelationshipClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{Firstname: "John", Surname: "Doe"}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, relationshipData{
			Entity:  "John Doe",
			DonorID: 123,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := Relationship(client, template.Func, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetRelationshipNoID(t *testing.T) {
	client := &mockRelationshipClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{Firstname: "John", Surname: "Doe"}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, relationshipData{
			Entity: "John Doe",
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	err := Relationship(client, template.Func, nil)(w, r)

	assert.NotNil(t, err)
}

func TestGetCreateRelationshipWithHXRequest(t *testing.T) {
	client := &mockRelationshipClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{Firstname: "John", Surname: "Doe"}, nil)

	partialTemplate := &mockTemplate{}
	partialTemplate.
		On("Func", mock.Anything, relationshipData{
			Entity:     "John Doe",
			DonorID:    123,
			EntityType: "person",
		}).
		Return(nil)

	template := &mockTemplate{}

	r, _ := http.NewRequest(http.MethodGet, "/?id=123&entity=person", nil)
	r.Header.Add("HX-Request", "true")
	w := httptest.NewRecorder()

	err := Relationship(client, template.Func, partialTemplate.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, partialTemplate)
	mock.AssertExpectationsForObjects(t, client, template)
	template.AssertNotCalled(t, "Func")
	partialTemplate.AssertCalled(t, "Func", mock.Anything, mock.Anything)
}

func TestGetRelationshipWhenPersonErrors(t *testing.T) {
	client := &mockRelationshipClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, errExample)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, relationshipData{
			Entity: "John Doe",
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := Relationship(client, template.Func, nil)(w, r)

	assert.Equal(t, errExample, err)
}

func TestGetRelationshipWhenTemplateErrors(t *testing.T) {
	client := &mockRelationshipClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{Firstname: "John", Surname: "Doe"}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, relationshipData{
			Entity:  "John Doe",
			DonorID: 123,
		}).
		Return(errExample)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := Relationship(client, template.Func, template.Func)(w, r)

	assert.Equal(t, errExample, err)
}

func TestPostRelationship(t *testing.T) {
	client := &mockRelationshipClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{Firstname: "John", Surname: "Doe"}, nil)
	client.
		On("CreatePersonReference", mock.Anything, 123, "7000-1000-1111", "Father").
		Return(nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, relationshipData{
			Success: true,
			Entity:  "John Doe",
			DonorID: 123,
		}).
		Return(nil)

	form := url.Values{
		"search": {"7000-1000-1111:Some Person (7000-1000-1111)"},
		"reason": {"Father"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := Relationship(client, template.Func, nil)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostRelationshipWhenCreatePersonReferenceValidationError(t *testing.T) {
	expectedError := sirius.ValidationError{
		Field: sirius.FieldErrors{"field": {"": "problem"}},
	}

	client := &mockRelationshipClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{Firstname: "John", Surname: "Doe"}, nil)
	client.
		On("CreatePersonReference", mock.Anything, 123, "7000-1000-1111", "Father").
		Return(expectedError)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, relationshipData{
			Success:    false,
			Error:      expectedError,
			DonorID:    123,
			Entity:     "John Doe",
			SearchUID:  "7000-1000-1111",
			SearchName: "Some Person (7000-1000-1111)",
			Reason:     "Father",
		}).
		Return(nil)

	form := url.Values{
		"search": {"7000-1000-1111:Some Person (7000-1000-1111)"},
		"reason": {"Father"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := Relationship(client, template.Func, nil)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestPostRelationshipWhenCreatePersonReferenceOtherError(t *testing.T) {
	client := &mockRelationshipClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{Firstname: "John", Surname: "Doe"}, nil)
	client.
		On("CreatePersonReference", mock.Anything, 123, "7000-1000-1111", "Father").
		Return(errExample)

	form := url.Values{
		"search": {"7000-1000-1111:Some Person (7000-1000-1111)"},
		"reason": {"Father"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := Relationship(client, nil, nil)(w, r)

	assert.Equal(t, errExample, err)
}
