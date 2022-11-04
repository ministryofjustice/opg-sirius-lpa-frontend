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

type mockInvestigationHoldClient struct {
	mock.Mock
}

func (m *mockInvestigationHoldClient) PlaceInvestigationOnHold(ctx sirius.Context, investigationID int, reason string) error {
	args := m.Called(ctx, investigationID, reason)
	return args.Error(0)
}

func (m *mockInvestigationHoldClient) TakeInvestigationOffHold(ctx sirius.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockInvestigationHoldClient) Investigation(ctx sirius.Context, id int) (sirius.Investigation, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Investigation), args.Error(1)
}

func TestGetPlaceInvestigationOnHold(t *testing.T) {
	investigation := sirius.Investigation{
		ID:       123,
		IsOnHold: false,
	}

	client := &mockInvestigationHoldClient{}
	client.
		On("Investigation", mock.Anything, 123).
		Return(investigation, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, investigationHoldData{
			Investigation: investigation,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := InvestigationHold(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetInvestigationOnHoldBadQuery(t *testing.T) {
	testCases := map[string]string{
		"no-id":  "/?test=lpa",
		"bad-id": "/?id=hello",
	}

	for name, urlParams := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, urlParams, nil)
			w := httptest.NewRecorder()

			err := InvestigationHold(nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestInvestigationOnHoldWhenGetInvestigationErrors(t *testing.T) {
	client := &mockInvestigationHoldClient{}
	client.
		On("Investigation", mock.Anything, 123).
		Return(sirius.Investigation{}, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := InvestigationHold(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetPlaceInvestigationOnHoldWhenTemplateErrors(t *testing.T) {
	investigation := sirius.Investigation{
		ID:       123,
		IsOnHold: false,
	}

	client := &mockInvestigationHoldClient{}
	client.
		On("Investigation", mock.Anything, 123).
		Return(investigation, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, investigationHoldData{
			Investigation: investigation,
		}).
		Return(expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := InvestigationHold(client, template.Func)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostPlaceInvestigationOnHold(t *testing.T) {
	investigation := sirius.Investigation{
		ID:       123,
		IsOnHold: false,
	}

	client := &mockInvestigationHoldClient{}
	client.
		On("Investigation", mock.Anything, 123).
		Return(investigation, nil).
		On("PlaceInvestigationOnHold", mock.Anything, 123, "Police Investigation").
		Return(nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, investigationHoldData{
			Success:       true,
			Investigation: investigation,
			Reason:        "Police Investigation",
		}).
		Return(nil)

	form := url.Values{
		"reason": {"Police Investigation"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := InvestigationHold(client, template.Func)(w, r)
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
		ID:       123,
		IsOnHold: false,
	}

	client := &mockInvestigationHoldClient{}
	client.
		On("Investigation", mock.Anything, 123).
		Return(investigation, nil).
		On("PlaceInvestigationOnHold", mock.Anything, 123, "Invalid Reason").
		Return(expectedError)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, investigationHoldData{
			Success:       false,
			Error:         expectedError,
			Investigation: investigation,
			Reason:        "Invalid Reason",
		}).
		Return(nil)

	form := url.Values{
		"reason": {"Invalid Reason"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := InvestigationHold(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostPlaceInvestigationOnHoldWhenOtherError(t *testing.T) {
	investigation := sirius.Investigation{
		ID:       123,
		IsOnHold: false,
	}

	client := &mockInvestigationHoldClient{}
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

	err := InvestigationHold(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetTakeInvestigationOffHold(t *testing.T) {
	investigation := sirius.Investigation{
		ID:       123,
		IsOnHold: true,
		HoldPeriods: []sirius.HoldPeriod{
			{
				ID:     1,
				Reason: "Police Investigation",
			},
		},
	}

	client := &mockInvestigationHoldClient{}
	client.
		On("Investigation", mock.Anything, 123).
		Return(investigation, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, investigationHoldData{
			Investigation: investigation,
			Reason:        "Police Investigation",
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := InvestigationHold(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetTakeInvestigationOffHoldErrorsIfNoOpenHoldPeriodOnInvestigation(t *testing.T) {
	investigation := sirius.Investigation{
		ID:       123,
		IsOnHold: true,
		HoldPeriods: []sirius.HoldPeriod{
			{
				ID:        1,
				Reason:    "Police Investigation",
				StartDate: sirius.DateString("2022-01-01"),
				EndDate:   sirius.DateString("2022-01-14"),
			},
		},
	}

	client := &mockInvestigationHoldClient{}
	client.
		On("Investigation", mock.Anything, 123).
		Return(investigation, nil)

	template := &mockTemplate{}

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := InvestigationHold(client, template.Func)(w, r)
	resp := w.Result()

	expected := errors.New("could not find open hold period on investigation")

	assert.Equal(t, expected, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostTakeInvestigationOffHold(t *testing.T) {
	investigation := sirius.Investigation{
		ID:       123,
		IsOnHold: true,
		HoldPeriods: []sirius.HoldPeriod{
			{
				ID:     1,
				Reason: "Police Investigation",
			},
		},
	}

	client := &mockInvestigationHoldClient{}
	client.
		On("Investigation", mock.Anything, 123).
		Return(investigation, nil).
		On("TakeInvestigationOffHold", mock.Anything, 1).
		Return(nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, investigationHoldData{
			Success:       true,
			Investigation: investigation,
			Reason:        "Police Investigation",
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", nil)
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := InvestigationHold(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostTakeInvestigationOnHoldUsesCorrectHoldPeriodWhenMultiple(t *testing.T) {
	investigation := sirius.Investigation{
		ID:       123,
		IsOnHold: true,
		HoldPeriods: []sirius.HoldPeriod{
			{
				ID:        1,
				Reason:    "Police Investigation",
				StartDate: sirius.DateString("2022-01-01"),
				EndDate:   sirius.DateString("2022-01-14"),
			},
			{
				ID:        2,
				Reason:    "LA Investigation",
				StartDate: sirius.DateString("2022-03-10"),
			},
		},
	}

	client := &mockInvestigationHoldClient{}
	client.
		On("Investigation", mock.Anything, 123).
		Return(investigation, nil).
		On("TakeInvestigationOffHold", mock.Anything, 2).
		Return(nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, investigationHoldData{
			Success:       true,
			Investigation: investigation,
			Reason:        "LA Investigation",
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", nil)
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := InvestigationHold(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostTakeInvestigationOffHoldWhenValidationError(t *testing.T) {
	expectedError := sirius.ValidationError{
		Field: sirius.FieldErrors{"field": {"": "problem"}},
	}

	investigation := sirius.Investigation{
		ID:       123,
		IsOnHold: true,
		HoldPeriods: []sirius.HoldPeriod{
			{
				ID:     1,
				Reason: "Police Investigation",
			},
		},
	}

	client := &mockInvestigationHoldClient{}
	client.
		On("Investigation", mock.Anything, 123).
		Return(investigation, nil).
		On("TakeInvestigationOffHold", mock.Anything, 1).
		Return(expectedError)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, investigationHoldData{
			Success:       false,
			Error:         expectedError,
			Investigation: investigation,
			Reason:        "Police Investigation",
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", nil)
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := InvestigationHold(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostTakeInvestigationOffHoldWhenOtherError(t *testing.T) {
	investigation := sirius.Investigation{
		ID:       123,
		IsOnHold: true,
		HoldPeriods: []sirius.HoldPeriod{
			{
				ID:     1,
				Reason: "Police Investigation",
			},
		},
	}

	client := &mockInvestigationHoldClient{}
	client.
		On("Investigation", mock.Anything, 123).
		Return(investigation, nil).
		On("TakeInvestigationOffHold", mock.Anything, 1).
		Return(expectedError)

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", nil)
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := InvestigationHold(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}
