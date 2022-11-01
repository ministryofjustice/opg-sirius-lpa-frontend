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

type mockPlaceInvestigationOnHoldClient struct {
	mock.Mock
}

func (m *mockPlaceInvestigationOnHoldClient) PlaceInvestigationOnHold(ctx sirius.Context, investigationID int, reason string) error {
	args := m.Called(ctx, investigationID, reason)
	return args.Error(0)
}

func (m *mockPlaceInvestigationOnHoldClient) Investigation(ctx sirius.Context, id int) (sirius.Investigation, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Investigation), args.Error(1)
}

func TestGetPlaceInvestigationOnHold(t *testing.T) {
	investigation := sirius.Investigation{
		ID:           123,
		Title:        "Test title",
		Information:  "Some test info",
		Type:         "Aspect",
		DateReceived: sirius.DateString("2022-01-01"),
	}

	client := &mockPlaceInvestigationOnHoldClient{}
	client.
		On("Investigation", mock.Anything, 123).
		Return(investigation, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, placeInvestigationOnHoldData{
			Investigation: investigation,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := PlaceInvestigationOnHold(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetPlaceInvestigationOnHoldBadQuery(t *testing.T) {
	testCases := map[string]string{
		"no-id":  "/?test=lpa",
		"bad-id": "/?id=hello",
	}

	for name, url := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			err := PlaceInvestigationOnHold(nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestGetPlaceInvestigationOnHoldWhenInvestigationErrors(t *testing.T) {
	expectedError := errors.New("err")

	client := &mockPlaceInvestigationOnHoldClient{}
	client.
		On("Investigation", mock.Anything, 123).
		Return(sirius.Investigation{}, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := PlaceInvestigationOnHold(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetPlaceInvestigationOnHoldWhenTemplateErrors(t *testing.T) {
	expectedError := errors.New("err")

	investigation := sirius.Investigation{
		ID:           123,
		Title:        "Test title",
		Information:  "Some test info",
		Type:         "Aspect",
		DateReceived: sirius.DateString("2022-01-01"),
	}

	client := &mockPlaceInvestigationOnHoldClient{}
	client.
		On("Investigation", mock.Anything, 123).
		Return(investigation, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, placeInvestigationOnHoldData{
			Investigation: investigation,
		}).
		Return(expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := PlaceInvestigationOnHold(client, template.Func)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostPlaceInvestigationOnHold(t *testing.T) {
	investigation := sirius.Investigation{
		ID:           123,
		Title:        "Test title",
		Information:  "Some test info",
		Type:         "Aspect",
		DateReceived: sirius.DateString("2022-01-01"),
	}

	client := &mockPlaceInvestigationOnHoldClient{}
	client.
		On("Investigation", mock.Anything, 123).
		Return(investigation, nil).
		On("PlaceInvestigationOnHold", mock.Anything, 123, "Police Investigation").
		Return(nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, placeInvestigationOnHoldData{
			Success:       true,
			Investigation: investigation,
		}).
		Return(nil)

	form := url.Values{
		"reason": {"Police Investigation"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := PlaceInvestigationOnHold(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostPlaceInvestigationOnHoldWhenValidationError(t *testing.T) {
	expectedError := sirius.ValidationError{
		Field: sirius.FieldErrors{"field": {"": "problem"}},
	}

	investigation := sirius.Investigation{
		Type: "Aspect",
	}

	client := &mockPlaceInvestigationOnHoldClient{}
	client.
		On("Investigation", mock.Anything, 123).
		Return(investigation, nil).
		On("PlaceInvestigationOnHold", mock.Anything, 123, "Invalid Reason").
		Return(expectedError)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, placeInvestigationOnHoldData{
			Success:       false,
			Error:         expectedError,
			Investigation: investigation,
		}).
		Return(nil)

	form := url.Values{
		"reason": {"Invalid Reason"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := PlaceInvestigationOnHold(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostPlaceInvestigationOnHoldWhenOtherError(t *testing.T) {
	expectedError := errors.New("err")

	investigation := sirius.Investigation{
		Type: "Priority",
	}

	client := &mockPlaceInvestigationOnHoldClient{}
	client.
		On("Investigation", mock.Anything, 123).
		Return(investigation, nil).
		On("PlaceInvestigationOnHold", mock.Anything, 123, "Police Investigation").
		Return(expectedError)

	form := url.Values{
		"reason": {"Police Investigation"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := PlaceInvestigationOnHold(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}
