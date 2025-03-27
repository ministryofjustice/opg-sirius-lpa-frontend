package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockSearchDonorsClient struct {
	mock.Mock
}

func (m *mockSearchDonorsClient) SearchDonors(ctx sirius.Context, term string) ([]sirius.Person, error) {
	args := m.Called(ctx, term)
	return args.Get(0).([]sirius.Person), args.Error(1)
}

func TestGetSearchDonors(t *testing.T) {
	expectedDonors := []sirius.Person{{ID: 1, Firstname: "John"}}

	client := &mockSearchDonorsClient{}
	client.
		On("SearchDonors", mock.Anything, "something").
		Return(expectedDonors, nil)

	req, _ := http.NewRequest(http.MethodGet, "/?q=something", nil)

	w := httptest.NewRecorder()
	err := SearchDonors(client)(w, req)

	assert.Nil(t, err)
	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var donors []sirius.Person
	_ = json.NewDecoder(resp.Body).Decode(&donors)

	assert.Equal(t, expectedDonors, donors)
}

func TestGetSearchDonorsWhenError(t *testing.T) {
	client := &mockSearchDonorsClient{}
	client.
		On("SearchDonors", mock.Anything, "something").
		Return([]sirius.Person{}, errExample)

	req, _ := http.NewRequest(http.MethodGet, "/?q=something", nil)

	w := httptest.NewRecorder()
	err := SearchDonors(client)(w, req)

	assert.Equal(t, errExample, err)
}
