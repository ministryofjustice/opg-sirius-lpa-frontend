package server

import (
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockCreateAdditionalDraftClient struct {
	mock.Mock
}

func (m *mockCreateAdditionalDraftClient) GetUserDetails(ctx sirius.Context) (sirius.User, error) {
	args := m.Called(ctx)
	return args.Get(0).(sirius.User), args.Error(1)
}

func (m *mockCreateAdditionalDraftClient) CreateAdditionalDraft(ctx sirius.Context, donorID int, draft sirius.AdditionalDraft) (map[string]string, error) {
	args := m.Called(ctx, donorID, draft)
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *mockCreateAdditionalDraftClient) Person(ctx sirius.Context, personID int) (sirius.Person, error) {
	args := m.Called(ctx, personID)
	return args.Get(0).(sirius.Person), args.Error(1)
}

func TestGetCreateAdditionalDraft(t *testing.T) {
	client := &mockCreateAdditionalDraftClient{}
	client.
		On("GetUserDetails", mock.Anything).
		Return(sirius.User{Roles: []string{"private-mlpa"}}, nil)
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{Firstname: "John", Surname: "Doe"}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createAdditionalDraftData{

			Donor: sirius.Person{
				ID:           123,
				Firstname:    "John",
				Surname:      "Doe",
				DateOfBirth:  "18/04/1965",
				AddressLine1: "9 Mount Pleasant Drive",
				Town:         "East Harling",
				Postcode:     "NR16 2GB",
				Country:      "UK",
				PersonType:   "Donor",
			},
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/create-additional-draft-lpa", nil)
	w := httptest.NewRecorder()

	err := CreateAdditionalDraft(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetCreateAdditionalDraftForbidden(t *testing.T) {
	client := &mockCreateAdditionalDraftClient{}
	client.
		On("GetUserDetails", mock.Anything).
		Return(sirius.User{}, nil)
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{Firstname: "John", Surname: "Doe"}, nil)

	template := &mockTemplate{}

	r, _ := http.NewRequest(http.MethodGet, "/create-additional-draft-lpa", nil)
	w := httptest.NewRecorder()

	err := CreateAdditionalDraft(client, template.Func)(w, r)

	assert.Equal(t, sirius.StatusError{Code: 403}, err)

	mock.AssertExpectationsForObjects(t, client, template)
}
