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

type mockSearchUsersClient struct {
	mock.Mock
}

func (m *mockSearchUsersClient) SearchUsers(ctx sirius.Context, term string) ([]sirius.User, error) {
	args := m.Called(ctx, term)
	return args.Get(0).([]sirius.User), args.Error(1)
}

func TestGetSearchUsers(t *testing.T) {
	expectedUsers := []sirius.User{{ID: 1, DisplayName: "Some something person"}}

	client := &mockSearchUsersClient{}
	client.
		On("SearchUsers", mock.Anything, "something").
		Return(expectedUsers, nil)

	req, _ := http.NewRequest(http.MethodGet, "/?q=something", nil)

	w := httptest.NewRecorder()
	err := SearchUsers(client)(w, req)

	assert.Nil(t, err)
	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var users []sirius.User
	_ = json.NewDecoder(resp.Body).Decode(&users)

	assert.Equal(t, expectedUsers, users)
}

func TestGetSearchUsersWhenError(t *testing.T) {
	client := &mockSearchUsersClient{}
	client.
		On("SearchUsers", mock.Anything, "something").
		Return([]sirius.User{}, expectedError)

	req, _ := http.NewRequest(http.MethodGet, "/?q=something", nil)

	w := httptest.NewRecorder()
	err := SearchUsers(client)(w, req)

	assert.Equal(t, expectedError, err)
}
