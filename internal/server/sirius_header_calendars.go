package server

import (
	"encoding/json"
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

type WorkingDaysData struct {
	XSRFToken      string
	StartDate      string
	EndDate        string
	NumWorkingDays int
	Mode           WorkingDaysMode
	PreviousMode   WorkingDaysMode
}

type siriusHeaderCalendarData struct {
	XSRFToken        string
	BankHolidaysJSON string
	Calculator       WorkingDaysData
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

		bankHolidaysJSON, err := json.Marshal(bankHolidays)
		if err != nil {
			bankHolidaysJSON = []byte("{}")
		}

		today := now().UTC().Truncate(24 * time.Hour)

		calculator := calculateWorkingDays(today, time.Time{}, 20, WorkingDaysModeEndDate, bankHolidays)
		calculator.XSRFToken = ctx.XSRFToken // TODO: do i need this?

		data := siriusHeaderCalendarData{
			XSRFToken:        ctx.XSRFToken,
			BankHolidaysJSON: string(bankHolidaysJSON),
			Calculator:       calculator,
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
	}

	return output
}
