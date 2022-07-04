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

type mockDeleteRelationshipClient struct {
	mock.Mock
}

func (m *mockDeleteRelationshipClient) PersonReferences(ctx sirius.Context, personID int) ([]sirius.PersonReference, error) {
	args := m.Called(ctx, personID)
	return args.Get(0).([]sirius.PersonReference), args.Error(1)
}

func (m *mockDeleteRelationshipClient) DeletePersonReference(ctx sirius.Context, referenceID int) error {
	return m.Called(ctx, referenceID).Error(0)
}

func (m *mockDeleteRelationshipClient) Person(ctx sirius.Context, id int) (sirius.Person, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Person), args.Error(1)
}

func TestGetDeleteRelationship(t *testing.T) {
	client := &mockDeleteRelationshipClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{Firstname: "John", Surname: "Doe"}, nil)
	client.
		On("PersonReferences", mock.Anything, 123).
		Return([]sirius.PersonReference{{ReferenceID: 1}}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, deleteRelationshipData{
			Entity:           "John Doe",
			PersonReferences: []sirius.PersonReference{{ReferenceID: 1}},
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := DeleteRelationship(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetDeleteRelationshipNoID(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	err := DeleteRelationship(nil, nil)(w, r)

	assert.NotNil(t, err)
}

func TestGetDeleteRelationshipWhenPersonErrors(t *testing.T) {
	expectedError := errors.New("hmm")

	client := &mockDeleteRelationshipClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{Firstname: "John", Surname: "Doe"}, expectedError)
	client.
		On("PersonReferences", mock.Anything, 123).
		Return([]sirius.PersonReference{{ReferenceID: 1}}, nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := DeleteRelationship(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetDeleteRelationshipWhenPersonReferencesErrors(t *testing.T) {
	expectedError := errors.New("hmm")

	client := &mockDeleteRelationshipClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{Firstname: "John", Surname: "Doe"}, nil)
	client.
		On("PersonReferences", mock.Anything, 123).
		Return([]sirius.PersonReference{{ReferenceID: 1}}, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := DeleteRelationship(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetDeleteRelationshipWhenTemplateErrors(t *testing.T) {
	expectedError := errors.New("hmm")

	client := &mockDeleteRelationshipClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{Firstname: "John", Surname: "Doe"}, nil)
	client.
		On("PersonReferences", mock.Anything, 123).
		Return([]sirius.PersonReference{{ReferenceID: 1}}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, deleteRelationshipData{
			Entity:           "John Doe",
			PersonReferences: []sirius.PersonReference{{ReferenceID: 1}},
		}).
		Return(expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := DeleteRelationship(client, template.Func)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostDeleteRelationship(t *testing.T) {
	client := &mockDeleteRelationshipClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{Firstname: "John", Surname: "Doe"}, nil)
	client.
		On("PersonReferences", mock.Anything, 123).
		Return([]sirius.PersonReference{{ReferenceID: 1}}, nil)
	client.
		On("DeletePersonReference", mock.Anything, 1).
		Return(nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, deleteRelationshipData{
			Entity:           "John Doe",
			Success:          true,
			PersonReferences: []sirius.PersonReference{{ReferenceID: 1}},
		}).
		Return(nil)

	form := url.Values{
		"reference-id": {"1"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := DeleteRelationship(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostDeleteRelationshipWhenDeletePersonReferenceErrors(t *testing.T) {
	expectedError := errors.New("hmm")

	client := &mockDeleteRelationshipClient{}
	client.
		On("DeletePersonReference", mock.Anything, 1).
		Return(expectedError)

	form := url.Values{
		"reference-id": {"1"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := DeleteRelationship(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}
