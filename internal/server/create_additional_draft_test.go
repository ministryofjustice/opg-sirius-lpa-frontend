package server

import (
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
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
		Return(sirius.Person{ID: 123, Firstname: "John", Surname: "Doe"}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createAdditionalDraftData{
			Donor: sirius.Person{
				ID:        123,
				Firstname: "John",
				Surname:   "Doe",
			},
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/create-additional-draft-lpa/?id=123", nil)
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
		Return(sirius.Person{ID: 123, Firstname: "John", Surname: "Doe"}, nil)

	template := &mockTemplate{}

	r, _ := http.NewRequest(http.MethodGet, "/create-additional-draft-lpa/?id=123", nil)
	w := httptest.NewRecorder()

	err := CreateAdditionalDraft(client, template.Func)(w, r)

	assert.Equal(t, sirius.StatusError{Code: 403}, err)

	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateAdditionalDraft(t *testing.T) {
	client := &mockCreateAdditionalDraftClient{}
	client.
		On("GetUserDetails", mock.Anything).
		Return(sirius.User{Roles: []string{"private-mlpa"}}, nil)
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{ID: 123, Firstname: "John", Surname: "Doe"}, nil)
	client.
		On("CreateAdditionalDraft", mock.Anything, 123, sirius.AdditionalDraft{
			CaseType:                  []string{"property-and-affairs", "personal-welfare"},
			CorrespondenceByWelsh:     true,
			CorrespondenceLargeFormat: false,
		}).
		Return(map[string]string{
			"property-and-affairs": "M-0123-4567-8901",
		}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createAdditionalDraftData{
			Form: formAdditionalDraft{
				SubTypes:                  []string{"property-and-affairs", "personal-welfare"},
				CorrespondenceByWelsh:     true,
				CorrespondenceLargeFormat: false,
			},
			Success: true,
			Uids: []createAdditionalDraftResult{
				{Subtype: "property-and-affairs", Uid: "M-0123-4567-8901"},
			},
			Donor: sirius.Person{
				ID:        123,
				Firstname: "John",
				Surname:   "Doe",
			},
		}).
		Return(nil)

	form := url.Values{
		"subtype":                   {"property-and-affairs", "personal-welfare"},
		"correspondenceByWelsh":     {"true"},
		"correspondenceLargeFormat": {"false"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/create-additional-draft-lpa/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateAdditionalDraft(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateAdditionalDraftWhenAPIFails(t *testing.T) {
	client := &mockCreateAdditionalDraftClient{}
	client.
		On("GetUserDetails", mock.Anything).
		Return(sirius.User{Roles: []string{"private-mlpa"}}, nil)
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{ID: 123, Firstname: "John", Surname: "Doe"}, nil)
	client.
		On("CreateAdditionalDraft", mock.Anything, 123, sirius.AdditionalDraft{
			CaseType: []string{"property-and-affairs", "personal-welfare"},
		}).
		Return(map[string]string{}, expectedError)

	template := &mockTemplate{}

	form := url.Values{
		"subtype": {"property-and-affairs", "personal-welfare"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/create-additional-draft-lpa/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateAdditionalDraft(client, template.Func)(w, r)
	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateAdditionalDraftWhenValidationError(t *testing.T) {
	client := &mockCreateAdditionalDraftClient{}
	client.
		On("GetUserDetails", mock.Anything).
		Return(sirius.User{Roles: []string{"private-mlpa"}}, nil)
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{ID: 123, Firstname: "John", Surname: "Doe"}, nil)
	client.
		On("CreateAdditionalDraft", mock.Anything, 123, sirius.AdditionalDraft{
			CorrespondenceByWelsh: false,
		}).
		Return(map[string]string{}, sirius.ValidationError{Field: sirius.FieldErrors{
			"subtype": {"required": "Value required and can't be empty"},
		}})

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createAdditionalDraftData{
			Form: formAdditionalDraft{},
			Error: sirius.ValidationError{
				Field: sirius.FieldErrors{
					"subtype": {"required": "Value required and can't be empty"},
				},
			},
			Donor: sirius.Person{
				ID:        123,
				Firstname: "John",
				Surname:   "Doe",
			},
		}).
		Return(nil)

	form := url.Values{}

	r, _ := http.NewRequest(http.MethodPost, "/create-additional-draft-lpa/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateAdditionalDraft(client, template.Func)(w, r)
	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}