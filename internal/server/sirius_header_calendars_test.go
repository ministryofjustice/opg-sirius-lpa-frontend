package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockSiriusCalendarsClient struct {
	mock.Mock
}

func (m *mockSiriusCalendarsClient) BankHolidays(ctx sirius.Context) (sirius.BankHolidays, error) {
	args := m.Called(ctx)
	return args.Get(0).(sirius.BankHolidays), args.Error(1)
}

func TestGetSiriusCalendars(t *testing.T) {
	bankHolidays := sirius.BankHolidays{
		"2025": {
			"New Year": "2025-01-01T00:00:00+00:00",
		},
	}

	client := &mockSiriusCalendarsClient{}
	client.
		On("BankHolidays", mock.Anything).
		Return(bankHolidays, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, siriusHeaderCalendarData{
			BankHolidaysJSON: `{"2025":{"New Year":"2025-01-01T00:00:00+00:00"}}`,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/lpa-api/v1/dates/bank-holidays", nil)
	w := httptest.NewRecorder()

	err := SiriusHeaderCalendars(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetSiriusCalendarsWhenBankHolidaysErrors(t *testing.T) {
	client := &mockSiriusCalendarsClient{}
	client.
		On("BankHolidays", mock.Anything).
		Return(sirius.BankHolidays{}, errExample)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, siriusHeaderCalendarData{
			BankHolidaysJSON: `{}`,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/lpa-api/v1/dates/bank-holidays", nil)
	w := httptest.NewRecorder()

	err := SiriusHeaderCalendars(client, template.Func)(w, r)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}
