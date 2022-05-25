package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockSearchPersonsClient struct {
	mock.Mock
}

func (m *mockSearchPersonsClient) SearchPersons(ctx sirius.Context, term string) ([]sirius.Person, error) {
	args := m.Called(ctx, term)
	return args.Get(0).([]sirius.Person), args.Error(1)
}

func TestGetSearchPersons(t *testing.T) {
	expectedPersons := []sirius.Person{{ID: 1, Firstname: "John"}}

	client := &mockSearchPersonsClient{}
	client.
		On("SearchPersons", mock.Anything, "something").
		Return(expectedPersons, nil)

	req, _ := http.NewRequest(http.MethodGet, "/?q=something", nil)

	w := httptest.NewRecorder()
	err := SearchPersons(client)(w, req)

	assert.Nil(t, err)
	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var persons []sirius.Person
	_ = json.NewDecoder(resp.Body).Decode(&persons)

	assert.Equal(t, expectedPersons, persons)
}

func TestGetSearchPersonsWhenError(t *testing.T) {
	expectedError := errors.New("hmm")

	client := &mockSearchPersonsClient{}
	client.
		On("SearchPersons", mock.Anything, "something").
		Return([]sirius.Person{}, expectedError)

	req, _ := http.NewRequest(http.MethodGet, "/?q=something", nil)

	w := httptest.NewRecorder()
	err := SearchPersons(client)(w, req)

	assert.Equal(t, expectedError, err)
}
