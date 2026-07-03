package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

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
	fixedNow := time.Date(2026, 6, 26, 12, 0, 0, 0, time.UTC)
	now := func() time.Time { return fixedNow }

	bankHolidays := sirius.BankHolidays{
		"2025": {
			"New Year": "2025-01-01T00:00:00+00:00",
		},
	}

	expectedCalculator := calculateWorkingDays(
		fixedNow.UTC().Truncate(24*time.Hour),
		time.Time{},
		20,
		WorkingDaysModeEndDate,
		bankHolidays,
	)

	client := &mockSiriusCalendarsClient{}
	client.
		On("BankHolidays", mock.Anything).
		Return(bankHolidays, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, siriusHeaderCalendarData{
			BankHolidaysJSON: `{"2025":{"New Year":"2025-01-01T00:00:00+00:00"}}`,
			Calculator:       expectedCalculator,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/lpa-api/v1/dates/bank-holidays", nil)
	w := httptest.NewRecorder()

	err := siriusHeaderCalendarsWithNow(client, template.Func, now)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetSiriusCalendarsWhenBankHolidaysErrors(t *testing.T) {
	fixedNow := time.Date(2026, 6, 26, 0, 0, 0, 0, time.UTC)
	now := func() time.Time { return fixedNow }

	expectedCalculator := calculateWorkingDays(
		fixedNow.UTC().Truncate(24*time.Hour),
		time.Time{},
		20,
		WorkingDaysModeEndDate,
		sirius.BankHolidays{},
	)

	client := &mockSiriusCalendarsClient{}
	client.
		On("BankHolidays", mock.Anything).
		Return(sirius.BankHolidays{}, errExample)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, siriusHeaderCalendarData{
			BankHolidaysJSON: `{}`,
			Calculator:       expectedCalculator,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/lpa-api/v1/dates/bank-holidays", nil)
	w := httptest.NewRecorder()

	err := siriusHeaderCalendarsWithNow(client, template.Func, now)(w, r)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetSiriusCalendarsWithEmptyBankHolidaysJSON(t *testing.T) {
	fixedNow := time.Date(2026, 6, 26, 0, 0, 0, 0, time.UTC)
	now := func() time.Time { return fixedNow }

	bankHolidays := sirius.BankHolidays{}

	expectedCalculator := calculateWorkingDays(
		fixedNow.UTC().Truncate(24*time.Hour),
		time.Time{},
		20,
		WorkingDaysModeEndDate,
		bankHolidays,
	)

	client := &mockSiriusCalendarsClient{}
	client.
		On("BankHolidays", mock.Anything).
		Return(bankHolidays, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, siriusHeaderCalendarData{
			BankHolidaysJSON: `{}`,
			Calculator:       expectedCalculator,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/lpa-api/v1/dates/bank-holidays", nil)
	w := httptest.NewRecorder()

	err := siriusHeaderCalendarsWithNow(client, template.Func, now)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestWorkingDays(t *testing.T) {
	bankHolidays := sirius.BankHolidays{
		"2025": {
			"New Year": "2025-01-01T00:00:00+00:00",
		},
	}

	expectedCalculator := calculateWorkingDays(
		time.Date(2026, 6, 26, 0, 0, 0, 0, time.UTC),
		time.Time{},
		20,
		WorkingDaysModeEndDate,
		bankHolidays,
	)
	expectedCalculator.PreviousMode = WorkingDaysModeStartDate

	client := &mockSiriusCalendarsClient{}
	client.
		On("BankHolidays", mock.Anything).
		Return(bankHolidays, nil)

	form := url.Values{
		"startdate":      {"2026-06-26"},
		"enddate":        {""},
		"numworkingdays": {"20"},
		"mode":           {"enddate"},
		"previousmode":   {"startdate"},
	}

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, expectedCalculator).
		Return(nil)

	r, _ := http.NewRequest(http.MethodPost, "/lpa-api/v1/dates/bank-holidays", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := WorkingDays(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestCalculateWorkingDays(t *testing.T) {
	testCases := []struct {
		name           string
		startDate      time.Time
		endDate        time.Time
		numWorkingDays int
		mode           WorkingDaysMode
		expected       WorkingDaysData
	}{
		{
			name:      "numworkingdays: between two dates with no weekends or bank holidays",
			startDate: time.Date(2026, 4, 7, 0, 0, 0, 0, time.UTC),
			endDate:   time.Date(2026, 4, 10, 0, 0, 0, 0, time.UTC),
			mode:      WorkingDaysModeNumWorkingDays,
			expected: WorkingDaysData{
				StartDate:      "2026-04-07",
				EndDate:        "2026-04-10",
				Mode:           WorkingDaysModeNumWorkingDays,
				NumWorkingDays: 3,
			},
		},
		{
			name:      "numworkingdays: between two dates with weekend",
			startDate: time.Date(2026, 4, 7, 0, 0, 0, 0, time.UTC),
			endDate:   time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC),
			mode:      WorkingDaysModeNumWorkingDays,
			expected: WorkingDaysData{
				StartDate:      "2026-04-07",
				EndDate:        "2026-04-13",
				Mode:           WorkingDaysModeNumWorkingDays,
				NumWorkingDays: 4,
			},
		},
		{
			name:      "numworkingdays: between two dates with bank holidays",
			startDate: time.Date(2026, 4, 2, 0, 0, 0, 0, time.UTC),
			endDate:   time.Date(2026, 4, 7, 0, 0, 0, 0, time.UTC),
			mode:      WorkingDaysModeNumWorkingDays,
			expected: WorkingDaysData{
				StartDate:      "2026-04-02",
				EndDate:        "2026-04-07",
				Mode:           WorkingDaysModeNumWorkingDays,
				NumWorkingDays: 1,
			},
		},
		{
			name:      "numberworkingdays: enddate before startdate, expect startdate to reset to day before enddate",
			startDate: time.Date(2026, 4, 25, 0, 0, 0, 0, time.UTC),
			endDate:   time.Date(2026, 4, 7, 0, 0, 0, 0, time.UTC),
			mode:      WorkingDaysModeNumWorkingDays,
			expected: WorkingDaysData{
				StartDate:      "2026-04-06",
				EndDate:        "2026-04-07",
				Mode:           WorkingDaysModeNumWorkingDays,
				NumWorkingDays: 1,
			},
		},
		{
			name:           "startdate: no weekends or bank holidays",
			endDate:        time.Date(2026, 4, 10, 0, 0, 0, 0, time.UTC),
			mode:           WorkingDaysModeStartDate,
			numWorkingDays: 3,
			expected: WorkingDaysData{
				StartDate:      "2026-04-07",
				EndDate:        "2026-04-10",
				Mode:           WorkingDaysModeStartDate,
				NumWorkingDays: 3,
			},
		},
		{
			name:           "startdate: weekends and bank holidays",
			endDate:        time.Date(2026, 4, 7, 0, 0, 0, 0, time.UTC),
			mode:           WorkingDaysModeStartDate,
			numWorkingDays: 1,
			expected: WorkingDaysData{
				StartDate:      "2026-04-02",
				EndDate:        "2026-04-07",
				Mode:           WorkingDaysModeStartDate,
				NumWorkingDays: 1,
			},
		},
		{
			name:           "startdate: numworkingdays is negative, expect startdate to reset to enddate",
			endDate:        time.Date(2026, 4, 7, 0, 0, 0, 0, time.UTC),
			mode:           WorkingDaysModeStartDate,
			numWorkingDays: -1,
			expected: WorkingDaysData{
				StartDate:      "2026-04-07",
				EndDate:        "2026-04-07",
				Mode:           WorkingDaysModeStartDate,
				NumWorkingDays: 0,
			},
		},
		{
			name:           "enddate: no weekends or bank holidays",
			startDate:      time.Date(2026, 4, 7, 0, 0, 0, 0, time.UTC),
			mode:           WorkingDaysModeEndDate,
			numWorkingDays: 3,
			expected: WorkingDaysData{
				StartDate:      "2026-04-07",
				EndDate:        "2026-04-10",
				Mode:           WorkingDaysModeEndDate,
				NumWorkingDays: 3,
			},
		},
		{
			name:           "enddate: weekends and bank holidays",
			startDate:      time.Date(2026, 4, 2, 0, 0, 0, 0, time.UTC),
			mode:           WorkingDaysModeEndDate,
			numWorkingDays: 1,
			expected: WorkingDaysData{
				StartDate:      "2026-04-02",
				EndDate:        "2026-04-07",
				Mode:           WorkingDaysModeEndDate,
				NumWorkingDays: 1,
			},
		},
		{
			name:           "enddate: numworkingdays is negative, expect enddate to reset to startdate",
			startDate:      time.Date(2026, 4, 7, 0, 0, 0, 0, time.UTC),
			mode:           WorkingDaysModeEndDate,
			numWorkingDays: -1,
			expected: WorkingDaysData{
				StartDate:      "2026-04-07",
				EndDate:        "2026-04-07",
				Mode:           WorkingDaysModeEndDate,
				NumWorkingDays: 0,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bankHolidays := sirius.BankHolidays{
				"2026": {
					"Good Friday":   "2026-04-03T00:00:00+00:00",
					"Easter Monday": "2026-04-06T00:00:00+00:00",
				},
			}
			result := calculateWorkingDays(tc.startDate, tc.endDate, tc.numWorkingDays, tc.mode, bankHolidays)
			assert.Equal(t, tc.expected, result)
		})
	}
}
