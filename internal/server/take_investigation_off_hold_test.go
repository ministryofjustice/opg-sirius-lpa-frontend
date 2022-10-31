package server

import (
	"errors"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockTakeInvestigationOffHoldClient struct {
	mock.Mock
}

func (m *mockTakeInvestigationOffHoldClient) TakeInvestigationOffHold(ctx sirius.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockTakeInvestigationOffHoldClient) HoldPeriod(ctx sirius.Context, id int) (sirius.HoldPeriod, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.HoldPeriod), args.Error(1)
}

func TestGetTakeInvestigationOffHold(t *testing.T) {
	holdPeriod := sirius.HoldPeriod{
		ID:            123,
		Reason:        "Police Investigation",
		Investigation: sirius.Investigation{Type: "Normal"},
	}

	client := &mockTakeInvestigationOffHoldClient{}
	client.
		On("HoldPeriod", mock.Anything, 123).
		Return(holdPeriod, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, takeInvestigationOffHoldData{
			HoldPeriod:    holdPeriod,
			Investigation: holdPeriod.Investigation,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := TakeInvestigationOffHold(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetTakeInvestigationOffHoldBadQuery(t *testing.T) {
	testCases := map[string]string{
		"no-id":  "/?test=lpa",
		"bad-id": "/?id=hello",
	}

	for name, url := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			err := TakeInvestigationOffHold(nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestGetTakeInvestigationOffHoldWhenHoldPeriodErrors(t *testing.T) {
	expectedError := errors.New("err")

	client := &mockTakeInvestigationOffHoldClient{}
	client.
		On("HoldPeriod", mock.Anything, 123).
		Return(sirius.HoldPeriod{}, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := TakeInvestigationOffHold(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetTakeInvestigationOffHoldWhenTemplateErrors(t *testing.T) {
	expectedError := errors.New("err")

	holdPeriod := sirius.HoldPeriod{
		ID:            123,
		Reason:        "Police Investigation",
		Investigation: sirius.Investigation{Type: "Normal"},
	}

	client := &mockTakeInvestigationOffHoldClient{}
	client.
		On("HoldPeriod", mock.Anything, 123).
		Return(holdPeriod, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, takeInvestigationOffHoldData{
			HoldPeriod:    holdPeriod,
			Investigation: holdPeriod.Investigation,
		}).
		Return(expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := TakeInvestigationOffHold(client, template.Func)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostTakeInvestigationOffHold(t *testing.T) {
	holdPeriod := sirius.HoldPeriod{
		ID:            123,
		Reason:        "Police Investigation",
		Investigation: sirius.Investigation{Type: "Normal"},
	}

	client := &mockTakeInvestigationOffHoldClient{}
	client.
		On("HoldPeriod", mock.Anything, 123).
		Return(holdPeriod, nil).
		On("TakeInvestigationOffHold", mock.Anything, 123).
		Return(nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, takeInvestigationOffHoldData{
			Success:       true,
			HoldPeriod:    holdPeriod,
			Investigation: holdPeriod.Investigation,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", nil)
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := TakeInvestigationOffHold(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostTakeInvestigationOffHoldWhenValidationError(t *testing.T) {
	expectedError := sirius.ValidationError{
		Field: sirius.FieldErrors{"field": {"": "problem"}},
	}

	holdPeriod := sirius.HoldPeriod{
		ID:            123,
		Reason:        "Police Investigation",
		Investigation: sirius.Investigation{Type: "Normal"},
	}

	client := &mockTakeInvestigationOffHoldClient{}
	client.
		On("HoldPeriod", mock.Anything, 123).
		Return(holdPeriod, nil).
		On("TakeInvestigationOffHold", mock.Anything, 123).
		Return(expectedError)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, takeInvestigationOffHoldData{
			Success:       false,
			Error:         expectedError,
			HoldPeriod:    holdPeriod,
			Investigation: holdPeriod.Investigation,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", nil)
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := TakeInvestigationOffHold(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostTakeInvestigationOffHoldWhenOtherError(t *testing.T) {
	expectedError := errors.New("err")

	holdPeriod := sirius.HoldPeriod{
		ID:            123,
		Reason:        "Police Investigation",
		Investigation: sirius.Investigation{Type: "Normal"},
	}

	client := &mockTakeInvestigationOffHoldClient{}
	client.
		On("HoldPeriod", mock.Anything, 123).
		Return(holdPeriod, nil).
		On("TakeInvestigationOffHold", mock.Anything, 123).
		Return(expectedError)

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", nil)
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := TakeInvestigationOffHold(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}
