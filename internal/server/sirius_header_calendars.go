package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type SiriusHeaderCalendarClient interface {
	BankHolidays(ctx sirius.Context) (sirius.BankHolidays, error)
}

type WorkingDaysMode string

const (
	WorkingDaysModeNumWorkingDays WorkingDaysMode = "numworkingdays"
	WorkingDaysModeStartDate      WorkingDaysMode = "startdate"
	WorkingDaysModeEndDate        WorkingDaysMode = "enddate"
)

func parseWorkingDaysMode(value string) WorkingDaysMode {
	switch WorkingDaysMode(strings.ToLower(value)) {
	case WorkingDaysModeNumWorkingDays, WorkingDaysModeStartDate, WorkingDaysModeEndDate:
		return WorkingDaysMode(strings.ToLower(value))
	default:
		return WorkingDaysModeNumWorkingDays
	}
}

func (WorkingDaysData) ModeStartDate() WorkingDaysMode {
	return WorkingDaysModeStartDate
}

func (WorkingDaysData) ModeEndDate() WorkingDaysMode {
	return WorkingDaysModeEndDate
}

func (WorkingDaysData) ModeNumWorkingDays() WorkingDaysMode {
	return WorkingDaysModeNumWorkingDays
}

type CalendarDay struct {
	Day           int
	Date          string
	IsToday       bool
	IsBankHoliday bool
}

type CalendarMonth struct {
	Name         string
	Weeks        [][]CalendarDay
	Year         int
	Month        int
	PrevMonthURL string
	NextMonthURL string
}

type WorkingDaysData struct {
	XSRFToken      string
	StartDate      string
	EndDate        string
	NumWorkingDays int
	Mode           WorkingDaysMode
	PreviousMode   WorkingDaysMode
}

type siriusHeaderCalendarData struct {
	XSRFToken  string
	Calculator WorkingDaysData
	Months     [3]CalendarMonth
}

func SiriusHeaderCalendars(client SiriusHeaderCalendarClient, tmpl template.Template) Handler {
	return siriusHeaderCalendarsWithNow(client, tmpl, time.Now)
}

func siriusHeaderCalendarsWithNow(client SiriusHeaderCalendarClient, tmpl template.Template, now func() time.Time) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		bankHolidays, err := client.BankHolidays(ctx)
		if err != nil {
			bankHolidays = sirius.BankHolidays{}
		}

		today := now().UTC().Truncate(24 * time.Hour)

		calculator := calculateWorkingDays(today, time.Time{}, 20, WorkingDaysModeEndDate, bankHolidays)
		calculator.XSRFToken = ctx.XSRFToken

		data := siriusHeaderCalendarData{
			XSRFToken:  ctx.XSRFToken,
			Calculator: calculator,
			Months:     buildCalendarMonths(bankHolidays, today),
		}

		return tmpl(w, data)
	}
}

func WorkingDays(client SiriusHeaderCalendarClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if err := r.ParseForm(); err != nil {
			return err
		}

		ctx := getContext(r)

		numWorkingDays, _ := strconv.Atoi(r.FormValue("numworkingdays"))
		startDate, _ := time.Parse(time.DateOnly, r.FormValue("startdate"))
		endDate, _ := time.Parse(time.DateOnly, r.FormValue("enddate"))
		mode := parseWorkingDaysMode(r.FormValue("mode"))
		previousMode := parseWorkingDaysMode(r.FormValue("previousmode"))

		bankHolidays, err := client.BankHolidays(ctx)
		if err != nil {
			bankHolidays = sirius.BankHolidays{}
		}

		data := calculateWorkingDays(startDate, endDate, numWorkingDays, mode, bankHolidays)
		data.PreviousMode = previousMode
		data.XSRFToken = ctx.XSRFToken

		return tmpl(w, data)
	}
}

func isWorkingDay(day time.Time, bankHolidaySet map[time.Time]bool) bool {
	weekendDays := map[time.Weekday]bool{
		time.Saturday: true,
		time.Sunday:   true,
	}
	if weekendDays[day.Weekday()] {
		return false
	}
	return !bankHolidaySet[day]
}

func calculateWorkingDays(startDate time.Time, endDate time.Time, numWorkingDays int, mode WorkingDaysMode, bankHolidays sirius.BankHolidays) WorkingDaysData {
	bankHolidaySet := make(map[time.Time]bool)
	for _, years := range bankHolidays {
		for _, date := range years {
			parsedDate, err := time.Parse(time.RFC3339, date)
			if err == nil {
				y, m, d := parsedDate.In(time.UTC).Date()
				bankHolidaySet[time.Date(y, m, d, 0, 0, 0, 0, time.UTC)] = true
			}
		}
	}

	output := WorkingDaysData{
		StartDate:      startDate.Format(time.DateOnly),
		EndDate:        endDate.Format(time.DateOnly),
		NumWorkingDays: numWorkingDays,
		Mode:           mode,
	}

	switch mode {
	case WorkingDaysModeNumWorkingDays:
		if !startDate.Before(endDate) {
			output.StartDate = endDate.AddDate(0, 0, -1).Format(time.DateOnly)
			output.NumWorkingDays = 1
			break
		}
		workingDaysCount, d := 0, startDate
		for d.Before(endDate) {
			if isWorkingDay(d, bankHolidaySet) {
				workingDaysCount++
			}
			d = d.AddDate(0, 0, 1)
		}
		output.NumWorkingDays = workingDaysCount

	case WorkingDaysModeStartDate:
		if numWorkingDays < 0 {
			output.StartDate = endDate.Format(time.DateOnly)
			output.NumWorkingDays = 0
			break
		}
		remaining, d := numWorkingDays, endDate
		for remaining > 0 {
			d = d.AddDate(0, 0, -1)
			if isWorkingDay(d, bankHolidaySet) {
				remaining--
			}
		}
		output.StartDate = d.Format(time.DateOnly)

	case WorkingDaysModeEndDate:
		if numWorkingDays < 0 {
			output.EndDate = startDate.Format(time.DateOnly)
			output.NumWorkingDays = 0
			break
		}
		remaining, d := numWorkingDays, startDate
		for remaining > 0 {
			d = d.AddDate(0, 0, 1)
			if isWorkingDay(d, bankHolidaySet) {
				remaining--
			}
		}
		output.EndDate = d.Format(time.DateOnly)

	default:

	}

	return output
}

func buildCalendarMonth(year int, month time.Month, bankHolidays sirius.BankHolidays, todayStr string) CalendarMonth {
	// Build a map of bank holiday dates for quick lookup
	bhSet := make(map[string]bool)
	for _, years := range bankHolidays {
		for _, dateStr := range years {
			parsedDate, err := time.Parse(time.RFC3339, dateStr)
			if err == nil {
				y, m, d := parsedDate.In(time.UTC).Date()
				bhSet[fmt.Sprintf("%04d-%02d-%02d", y, int(m), d)] = true
			}
		}
	}

	first := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	last := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC)

	// ISO weekday Mon=1…Sun=7; convert to 0-based Monday offset for grid.
	startOffset := int(first.Weekday())
	if startOffset == 0 {
		startOffset = 7
	}
	startOffset--

	var days []CalendarDay
	for i := 0; i < startOffset; i++ {
		days = append(days, CalendarDay{})
	}
	for d := 1; d <= last.Day(); d++ {
		dateStr := fmt.Sprintf("%04d-%02d-%02d", year, int(month), d)
		days = append(days, CalendarDay{
			Day:           d,
			Date:          dateStr,
			IsToday:       dateStr == todayStr,
			IsBankHoliday: bhSet[dateStr],
		})
	}
	for len(days)%7 != 0 {
		days = append(days, CalendarDay{})
	}

	var weeks [][]CalendarDay
	for i := 0; i < len(days); i += 7 {
		row := make([]CalendarDay, 7)
		copy(row, days[i:i+7])
		weeks = append(weeks, row)
	}

	// Calculate previous and next month
	prevMonth := first.AddDate(0, -1, 0)
	nextMonth := first.AddDate(0, 1, 0)
	prevMonthURL := fmt.Sprintf("/calendar-month?year=%d&month=%d", prevMonth.Year(), prevMonth.Month())
	nextMonthURL := fmt.Sprintf("/calendar-month?year=%d&month=%d", nextMonth.Year(), nextMonth.Month())

	return CalendarMonth{
		Name:         first.Format("January 2006"),
		Weeks:        weeks,
		Year:         year,
		Month:        int(month),
		PrevMonthURL: prevMonthURL,
		NextMonthURL: nextMonthURL,
	}
}

func buildCalendarMonths(bankHolidays sirius.BankHolidays, today time.Time) [3]CalendarMonth {
	todayStr := today.Format("2006-01-02")
	prev := today.AddDate(0, -1, 0)
	next := today.AddDate(0, 1, 0)
	return [3]CalendarMonth{
		buildCalendarMonth(prev.Year(), prev.Month(), bankHolidays, todayStr),
		buildCalendarMonth(today.Year(), today.Month(), bankHolidays, todayStr),
		buildCalendarMonth(next.Year(), next.Month(), bankHolidays, todayStr),
	}
}

func CalendarMonthPartial(client SiriusHeaderCalendarClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		// Parse year and month from query parameters
		yearStr := r.URL.Query().Get("year")
		monthStr := r.URL.Query().Get("month")

		if yearStr == "" || monthStr == "" {
			// If not provided, use current date
			now := time.Now().UTC()
			yearStr = fmt.Sprintf("%d", now.Year())
			monthStr = fmt.Sprintf("%d", now.Month())
		}

		year, err := strconv.Atoi(yearStr)
		if err != nil {
			return err
		}

		month, err := strconv.Atoi(monthStr)
		if err != nil {
			return err
		}

		// Validate month (1-12)
		if month < 1 || month > 12 {
			return fmt.Errorf("invalid month: %d", month)
		}

		bankHolidays, err := client.BankHolidays(ctx)
		if err != nil {
			bankHolidays = sirius.BankHolidays{}
		}

		today := time.Now().UTC().Truncate(24 * time.Hour)
		todayStr := today.Format("2006-01-02")

		// Build the calendar month for the requested date
		calMonth := buildCalendarMonth(year, time.Month(month), bankHolidays, todayStr)

		return tmpl(w, calMonth)
	}
}
