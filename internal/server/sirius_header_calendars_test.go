package server

import (
	"fmt"
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

	today := fixedNow.UTC().Truncate(24 * time.Hour)
	expectedCalculator := calculateWorkingDays(
		today,
		time.Time{},
		20,
		WorkingDaysModeEndDate,
		bankHolidays,
	)

	expectedMonths := buildCalendarMonths(bankHolidays, today)

	client := &mockSiriusCalendarsClient{}
	client.
		On("BankHolidays", mock.Anything).
		Return(bankHolidays, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, siriusHeaderCalendarData{
			Calculator: expectedCalculator,
			Months:     expectedMonths,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/lpa-api/v1/dates/bank-holidays", nil)
	w := httptest.NewRecorder()

	err := siriusHeaderCalendarsWithNow(client, template.Func, now)(PageVars{}, w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetSiriusCalendarsWhenBankHolidaysErrors(t *testing.T) {
	fixedNow := time.Date(2026, 6, 26, 0, 0, 0, 0, time.UTC)
	now := func() time.Time { return fixedNow }

	today := fixedNow.UTC().Truncate(24 * time.Hour)
	expectedCalculator := calculateWorkingDays(
		today,
		time.Time{},
		20,
		WorkingDaysModeEndDate,
		sirius.BankHolidays{},
	)

	expectedMonths := buildCalendarMonths(sirius.BankHolidays{}, today)

	client := &mockSiriusCalendarsClient{}
	client.
		On("BankHolidays", mock.Anything).
		Return(sirius.BankHolidays{}, errExample)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, siriusHeaderCalendarData{
			Calculator: expectedCalculator,
			Months:     expectedMonths,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/lpa-api/v1/dates/bank-holidays", nil)
	w := httptest.NewRecorder()

	err := siriusHeaderCalendarsWithNow(client, template.Func, now)(PageVars{}, w, r)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetSiriusCalendarsWithEmptyBankHolidaysJSON(t *testing.T) {
	fixedNow := time.Date(2026, 6, 26, 0, 0, 0, 0, time.UTC)
	now := func() time.Time { return fixedNow }

	bankHolidays := sirius.BankHolidays{}

	today := fixedNow.UTC().Truncate(24 * time.Hour)
	expectedCalculator := calculateWorkingDays(
		today,
		time.Time{},
		20,
		WorkingDaysModeEndDate,
		bankHolidays,
	)

	expectedMonths := buildCalendarMonths(bankHolidays, today)

	client := &mockSiriusCalendarsClient{}
	client.
		On("BankHolidays", mock.Anything).
		Return(bankHolidays, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, siriusHeaderCalendarData{
			Calculator: expectedCalculator,
			Months:     expectedMonths,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/lpa-api/v1/dates/bank-holidays", nil)
	w := httptest.NewRecorder()

	err := siriusHeaderCalendarsWithNow(client, template.Func, now)(PageVars{}, w, r)
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

	err := WorkingDays(client, template.Func)(PageVars{}, w, r)
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

func TestCalendarMonthPartialWithValidParameters(t *testing.T) {
	bankHolidays := sirius.BankHolidays{
		"2026": {
			"Good Friday":   "2026-04-03T00:00:00+00:00",
			"Easter Monday": "2026-04-06T00:00:00+00:00",
		},
	}

	client := &mockSiriusCalendarsClient{}
	client.
		On("BankHolidays", mock.Anything).
		Return(bankHolidays, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.MatchedBy(func(data CalendarMonth) bool {
			return data.Year == 2026 && data.Month == 4
		})).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/calendar-month?year=2026&month=4", nil)
	w := httptest.NewRecorder()

	err := CalendarMonthPartial(client, template.Func)(PageVars{}, w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestCalendarMonthPartialWithMissingYear(t *testing.T) {
	bankHolidays := sirius.BankHolidays{}

	client := &mockSiriusCalendarsClient{}
	client.
		On("BankHolidays", mock.Anything).
		Return(bankHolidays, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.MatchedBy(func(data CalendarMonth) bool {
			// Should use current date when year is missing
			return data.Year == time.Now().Year()
		})).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/calendar-month?month=6", nil)
	w := httptest.NewRecorder()

	err := CalendarMonthPartial(client, template.Func)(PageVars{}, w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestCalendarMonthPartialWithMissingMonth(t *testing.T) {
	bankHolidays := sirius.BankHolidays{}

	client := &mockSiriusCalendarsClient{}
	client.
		On("BankHolidays", mock.Anything).
		Return(bankHolidays, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.MatchedBy(func(data CalendarMonth) bool {
			// Should use current date when month is missing
			return data.Month == int(time.Now().Month())
		})).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/calendar-month?year=2026", nil)
	w := httptest.NewRecorder()

	err := CalendarMonthPartial(client, template.Func)(PageVars{}, w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestCalendarMonthPartialWithMissingBoth(t *testing.T) {
	bankHolidays := sirius.BankHolidays{}

	client := &mockSiriusCalendarsClient{}
	client.
		On("BankHolidays", mock.Anything).
		Return(bankHolidays, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.MatchedBy(func(data CalendarMonth) bool {
			// Should use current date when both are missing
			now := time.Now().UTC()
			return data.Year == now.Year() && data.Month == int(now.Month())
		})).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/calendar-month", nil)
	w := httptest.NewRecorder()

	err := CalendarMonthPartial(client, template.Func)(PageVars{}, w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestCalendarMonthPartialWithInvalidYear(t *testing.T) {
	client := &mockSiriusCalendarsClient{}
	template := &mockTemplate{}

	r, _ := http.NewRequest(http.MethodGet, "/calendar-month?year=notanumber&month=6", nil)
	w := httptest.NewRecorder()

	err := CalendarMonthPartial(client, template.Func)(PageVars{}, w, r)

	assert.NotNil(t, err)
}

func TestCalendarMonthPartialWithInvalidMonth(t *testing.T) {
	client := &mockSiriusCalendarsClient{}
	template := &mockTemplate{}

	r, _ := http.NewRequest(http.MethodGet, "/calendar-month?year=2026&month=notanumber", nil)
	w := httptest.NewRecorder()

	err := CalendarMonthPartial(client, template.Func)(PageVars{}, w, r)

	assert.NotNil(t, err)
}

func TestCalendarMonthPartialWithMonthTooSmall(t *testing.T) {
	client := &mockSiriusCalendarsClient{}
	template := &mockTemplate{}

	r, _ := http.NewRequest(http.MethodGet, "/calendar-month?year=2026&month=0", nil)
	w := httptest.NewRecorder()

	err := CalendarMonthPartial(client, template.Func)(PageVars{}, w, r)

	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "invalid month")
}

func TestCalendarMonthPartialWithMonthTooLarge(t *testing.T) {
	client := &mockSiriusCalendarsClient{}
	template := &mockTemplate{}

	r, _ := http.NewRequest(http.MethodGet, "/calendar-month?year=2026&month=13", nil)
	w := httptest.NewRecorder()

	err := CalendarMonthPartial(client, template.Func)(PageVars{}, w, r)

	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "invalid month")
}

func TestCalendarMonthPartialWhenBankHolidaysErrors(t *testing.T) {
	client := &mockSiriusCalendarsClient{}
	client.
		On("BankHolidays", mock.Anything).
		Return(sirius.BankHolidays{}, errExample)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.MatchedBy(func(data CalendarMonth) bool {
			// Should continue with empty bank holidays when client errors
			return data.Year == 2026 && data.Month == 4
		})).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/calendar-month?year=2026&month=4", nil)
	w := httptest.NewRecorder()

	err := CalendarMonthPartial(client, template.Func)(PageVars{}, w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestCalendarMonthPartialWhenTemplateErrors(t *testing.T) {
	bankHolidays := sirius.BankHolidays{}

	client := &mockSiriusCalendarsClient{}
	client.
		On("BankHolidays", mock.Anything).
		Return(bankHolidays, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.Anything).
		Return(errExample)

	r, _ := http.NewRequest(http.MethodGet, "/calendar-month?year=2026&month=4", nil)
	w := httptest.NewRecorder()

	err := CalendarMonthPartial(client, template.Func)(PageVars{}, w, r)

	assert.NotNil(t, err)
	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestBuildCalendarMonths(t *testing.T) {
	bankHolidays := sirius.BankHolidays{
		"2025": {
			"New Year": "2025-01-01T00:00:00+00:00",
		},
	}

	today := time.Date(2026, 6, 26, 0, 0, 0, 0, time.UTC)
	result := buildCalendarMonths(bankHolidays, today)

	assert.Equal(t, 3, len(result))

	assert.Equal(t, 2026, result[0].Year)
	assert.Equal(t, 5, result[0].Month)
	assert.Equal(t, "May 2026", result[0].Name)
	assert.NotEmpty(t, result[0].Weeks)
	assert.Equal(t, 7, len(result[0].Weeks[0]), "each week should have 7 days")

	assert.Equal(t, 2026, result[1].Year)
	assert.Equal(t, 6, result[1].Month)
	assert.Equal(t, "June 2026", result[1].Name)
	assert.NotEmpty(t, result[1].Weeks)
	assert.Equal(t, 7, len(result[1].Weeks[0]), "each week should have 7 days")

	todayFound := false
	for _, week := range result[1].Weeks {
		for _, day := range week {
			if day.Day == 26 {
				todayFound = true
				assert.True(t, day.IsToday, "June 26 should be marked as today")
			}
		}
	}
	assert.True(t, todayFound, "June 26 should exist in the calendar")

	assert.Equal(t, 2026, result[2].Year)
	assert.Equal(t, 7, result[2].Month)
	assert.Equal(t, "July 2026", result[2].Name)
	assert.NotEmpty(t, result[2].Weeks)
	assert.Equal(t, 7, len(result[2].Weeks[0]), "each week should have 7 days")
}

func TestBuildCalendarMonth(t *testing.T) {
	bankHolidays := sirius.BankHolidays{
		"2026": {
			"Good Friday":   "2026-04-03T00:00:00+00:00",
			"Easter Monday": "2026-04-06T00:00:00+00:00",
		},
	}

	todayStr := "2026-04-15"
	result := buildCalendarMonth(2026, time.April, bankHolidays, todayStr)

	assert.Equal(t, "April 2026", result.Name)
	assert.Equal(t, 2026, result.Year)
	assert.Equal(t, 4, result.Month)
	assert.NotEmpty(t, result.Weeks)

	assert.Equal(t, "/calendar-month?year=2026&month=3", result.PrevMonthURL)
	assert.Equal(t, "/calendar-month?year=2026&month=5", result.NextMonthURL)

	for _, week := range result.Weeks {
		assert.Equal(t, 7, len(week), "each week should have exactly 7 days")
	}

	assert.Equal(t, 5, len(result.Weeks))

	totalDays := 0
	for _, week := range result.Weeks {
		for _, day := range week {
			if day.Day > 0 {
				totalDays++
			}
		}
	}
	assert.Equal(t, 30, totalDays, "April should have 30 days")

	todayFound := false
	for _, week := range result.Weeks {
		for _, day := range week {
			if day.Day == 15 {
				todayFound = true
				assert.True(t, day.IsToday, "April 15 should be marked as today")
				assert.Equal(t, "2026-04-15", day.Date)
			}
		}
	}
	assert.True(t, todayFound, "April 15 should exist in calendar")

	bankHolidaysFound := map[int]bool{3: false, 6: false}
	for _, week := range result.Weeks {
		for _, day := range week {
			if day.Day == 3 || day.Day == 6 {
				if day.Day == 3 {
					assert.True(t, day.IsBankHoliday, "Good Friday (April 3) should be marked as bank holiday")
				}
				if day.Day == 6 {
					assert.True(t, day.IsBankHoliday, "Easter Monday (April 6) should be marked as bank holiday")
				}
				bankHolidaysFound[day.Day] = true
			}
		}
	}
	assert.True(t, bankHolidaysFound[3], "Good Friday should be found and marked")
	assert.True(t, bankHolidaysFound[6], "Easter Monday should be found and marked")
}

func TestBuildCalendarMonthFebruary(t *testing.T) {
	bankHolidays := sirius.BankHolidays{}
	todayStr := "2026-02-14"
	result := buildCalendarMonth(2026, time.February, bankHolidays, todayStr)

	assert.Equal(t, "February 2026", result.Name)
	assert.Equal(t, 2026, result.Year)
	assert.Equal(t, 2, result.Month)

	totalDays := 0
	for _, week := range result.Weeks {
		for _, day := range week {
			if day.Day > 0 {
				totalDays++
			}
		}
	}
	assert.Equal(t, 28, totalDays, "February 2026 should have 28 days")

	for _, week := range result.Weeks {
		for _, day := range week {
			if day.Day == 14 {
				assert.True(t, day.IsToday, "February 14 should be marked as today")
			}
		}
	}
}

func TestBuildCalendarMonthLeapYear(t *testing.T) {
	bankHolidays := sirius.BankHolidays{}
	todayStr := "2024-02-29"
	result := buildCalendarMonth(2024, time.February, bankHolidays, todayStr)

	// Verify February in leap year
	assert.Equal(t, "February 2024", result.Name)

	totalDays := 0
	for _, week := range result.Weeks {
		for _, day := range week {
			if day.Day > 0 {
				totalDays++
			}
		}
	}
	assert.Equal(t, 29, totalDays, "February 2024 should have 29 days (leap year)")

	for _, week := range result.Weeks {
		for _, day := range week {
			if day.Day == 29 {
				assert.True(t, day.IsToday, "February 29 should be marked as today")
			}
		}
	}
}

func TestBuildCalendarMonthWithoutBankHolidays(t *testing.T) {
	bankHolidays := sirius.BankHolidays{}
	todayStr := "2026-06-15"
	result := buildCalendarMonth(2026, time.June, bankHolidays, todayStr)

	for _, week := range result.Weeks {
		for _, day := range week {
			assert.False(t, day.IsBankHoliday, "no days should be marked as bank holiday when none provided")
		}
	}
}

func TestBuildCalendarMonthYearChange(t *testing.T) {
	bankHolidays := sirius.BankHolidays{}

	decemberResult := buildCalendarMonth(2025, time.December, bankHolidays, "2025-12-15")
	assert.Equal(t, "/calendar-month?year=2025&month=11", decemberResult.PrevMonthURL)
	assert.Equal(t, "/calendar-month?year=2026&month=1", decemberResult.NextMonthURL)

	januaryResult := buildCalendarMonth(2026, time.January, bankHolidays, "2026-01-15")
	assert.Equal(t, "/calendar-month?year=2025&month=12", januaryResult.PrevMonthURL)
	assert.Equal(t, "/calendar-month?year=2026&month=2", januaryResult.NextMonthURL)
}

func TestBuildCalendarMonthDayOfWeekAlignment(t *testing.T) {
	bankHolidays := sirius.BankHolidays{}
	todayStr := "2026-06-01"
	result := buildCalendarMonth(2026, time.June, bankHolidays, todayStr)

	firstWeek := result.Weeks[0]
	assert.Equal(t, 1, firstWeek[0].Day, "June 1 should be on Monday (first position in week)")
	assert.Equal(t, "2026-06-01", firstWeek[0].Date)

	for i := 0; i < 0; i++ {
		assert.Equal(t, 0, firstWeek[0].Day, "days before June 1 should be empty")
	}
}

func TestBuildCalendarMonthAllMonths(t *testing.T) {
	expectedDays := map[time.Month]int{
		time.January:   31,
		time.February:  28,
		time.March:     31,
		time.April:     30,
		time.May:       31,
		time.June:      30,
		time.July:      31,
		time.August:    31,
		time.September: 30,
		time.October:   31,
		time.November:  30,
		time.December:  31,
	}

	bankHolidays := sirius.BankHolidays{}

	for month, expectedCount := range expectedDays {
		result := buildCalendarMonth(2026, month, bankHolidays, "2026-01-01")

		totalDays := 0
		for _, week := range result.Weeks {
			for _, day := range week {
				if day.Day > 0 {
					totalDays++
				}
			}
		}
		assert.Equal(t, expectedCount, totalDays, fmt.Sprintf("%v 2026 should have %d days", month, expectedCount))
	}
}

func TestBuildCalendarMonthBankHolidayFromPreviousYear(t *testing.T) {
	bankHolidays := sirius.BankHolidays{
		"2025": {
			"Previous Year Holiday": "2025-12-25T00:00:00+00:00",
		},
		"2026": {
			"New Year": "2026-01-01T00:00:00+00:00",
		},
	}

	todayStr := "2026-01-15"
	result := buildCalendarMonth(2026, time.January, bankHolidays, todayStr)

	for _, week := range result.Weeks {
		for _, day := range week {
			if day.Day == 1 {
				assert.True(t, day.IsBankHoliday, "January 1 should be marked as bank holiday")
			}
		}
	}
}
