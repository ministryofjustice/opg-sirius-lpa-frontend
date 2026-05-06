package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// CalendarDay represents one cell in the calendar grid.
// Day == 0 indicates a padding cell before the first of the month.
type CalendarDay struct {
	Day           int
	Date          string // "YYYY-MM-DD", empty for padding cells
	IsToday       bool
	IsBankHoliday bool
}

// CalendarMonth is a pre-computed month ready for template rendering.
// Weeks contains rows of 7 days ordered Mon–Sun.
type CalendarMonth struct {
	Name  string // e.g. "April 2026"
	Weeks [][]CalendarDay
}

// WorkingDaysVars holds the state of the working-days calculator.
// It is passed to the workingDaysCalculator template both when embedded in
// the full page (via .Calculator) and directly in HTMX partial responses.
type WorkingDaysVars struct {
	XSRFToken  string
	DonorID    int
	Date1      string // DD/MM/YYYY
	Date2      string // DD/MM/YYYY
	Difference int
	Mode       string // "date1" | "date2" | "difference"
}

// DonorHeaderData is the template data for the donor case record page.
type DonorHeaderData struct {
	XSRFToken  string
	Person     sirius.Person
	Cases      []sirius.Case
	Months     [3]CalendarMonth
	Calculator WorkingDaysVars
	DonorID    int
}

// ---------------------------------------------------------------------------
// Bank holidays — cached fetch from GOV.UK API
// ---------------------------------------------------------------------------

type bankHolidayCache struct {
	mu        sync.Mutex
	holidays  []time.Time
	fetchedAt time.Time
	ttl       time.Duration
}

type govukBankHolidayResponse struct {
	EnglandAndWales struct {
		Events []struct {
			Date string `json:"date"`
		} `json:"events"`
	} `json:"england-and-wales"`
}

var bhCache = &bankHolidayCache{ttl: 24 * time.Hour}

func (c *bankHolidayCache) get() ([]time.Time, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if time.Since(c.fetchedAt) < c.ttl && len(c.holidays) > 0 {
		return c.holidays, nil
	}

	// TODO: If a bank holidays endpoint already exists on *sirius.Client,
	// replace this direct GOV.UK call. Using http.DefaultClient is
	// acceptable for a POC; production should inject the client.
	resp, err := http.DefaultClient.Get("https://www.gov.uk/bank-holidays.json")
	if err != nil {
		return nil, fmt.Errorf("fetching bank holidays: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	var data govukBankHolidayResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("decoding bank holidays: %w", err)
	}

	var holidays []time.Time
	for _, event := range data.EnglandAndWales.Events {
		t, err := time.Parse("2006-01-02", event.Date)
		if err == nil {
			holidays = append(holidays, t)
		}
	}

	c.holidays = holidays
	c.fetchedAt = time.Now()
	return holidays, nil
}

// ---------------------------------------------------------------------------
// Calendar building
// ---------------------------------------------------------------------------

func buildCalendarMonth(year int, month time.Month, bhSet map[string]bool, todayStr string) CalendarMonth {
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

	return CalendarMonth{Name: first.Format("January 2006"), Weeks: weeks}
}

func buildCalendarMonths(bankHolidays []time.Time, today time.Time) [3]CalendarMonth {
	bhSet := make(map[string]bool, len(bankHolidays))
	for _, bh := range bankHolidays {
		bhSet[bh.Format("2006-01-02")] = true
	}
	todayStr := today.Format("2006-01-02")
	prev := today.AddDate(0, -1, 0)
	next := today.AddDate(0, 1, 0)
	return [3]CalendarMonth{
		buildCalendarMonth(prev.Year(), prev.Month(), bhSet, todayStr),
		buildCalendarMonth(today.Year(), today.Month(), bhSet, todayStr),
		buildCalendarMonth(next.Year(), next.Month(), bhSet, todayStr),
	}
}

// ---------------------------------------------------------------------------
// Working days calculation
//
// Mirrors PanelCalendarCtrl exactly:
//   "difference" — count working days from date1 up to (not including) date2
//   "date2"      — step forward from date1 until <difference> working days elapsed
//   "date1"      — step backward from date2 until <difference> working days elapsed
// ---------------------------------------------------------------------------

const dateLayout = "02/01/2006"

func isWorkingDay(t time.Time, bhSet map[string]bool) bool {
	if t.Weekday() == time.Saturday || t.Weekday() == time.Sunday {
		return false
	}
	return !bhSet[t.Format("2006-01-02")]
}

func calculateWorkingDays(req WorkingDaysVars, bankHolidays []time.Time) WorkingDaysVars {
	bhSet := make(map[string]bool, len(bankHolidays))
	for _, bh := range bankHolidays {
		bhSet[bh.Format("2006-01-02")] = true
	}

	switch req.Mode {
	case "difference":
		d1, err1 := time.Parse(dateLayout, req.Date1)
		d2, err2 := time.Parse(dateLayout, req.Date2)
		if err1 != nil || err2 != nil {
			break
		}
		if !d1.Before(d2) {
			d2 = d1.AddDate(0, 0, 1)
			req.Date2 = d2.Format(dateLayout)
		}
		count, cur := 0, d1
		for cur.Before(d2) {
			if isWorkingDay(cur, bhSet) {
				count++
			}
			cur = cur.AddDate(0, 0, 1)
		}
		req.Difference = count

	case "date2":
		d1, err := time.Parse(dateLayout, req.Date1)
		if err != nil || req.Difference < 0 {
			break
		}
		remaining, cur := req.Difference, d1
		for remaining > 0 {
			cur = cur.AddDate(0, 0, 1)
			if isWorkingDay(cur, bhSet) {
				remaining--
			}
		}
		req.Date2 = cur.Format(dateLayout)

	case "date1":
		d2, err := time.Parse(dateLayout, req.Date2)
		if err != nil || req.Difference < 0 {
			break
		}
		remaining, cur := req.Difference, d2
		for remaining > 0 {
			cur = cur.AddDate(0, 0, -1)
			if isWorkingDay(cur, bhSet) {
				remaining--
			}
		}
		req.Date1 = cur.Format(dateLayout)
	}

	return req
}

// ---------------------------------------------------------------------------
// Client interface
// ---------------------------------------------------------------------------

type DonorHeaderClient interface {
	Person(ctx sirius.Context, id int) (sirius.Person, error)
	CasesByDonor(ctx sirius.Context, id int) ([]sirius.Case, error)
}

// ---------------------------------------------------------------------------
// Handlers
// ---------------------------------------------------------------------------

// DonorHeader renders the donor case record page with the Go header bar.
//
// Wire in server.go:
//
//	mux.Handle("/donor/{id}", wrap(DonorHeader(client, templates.Get("donor-header-wrapper.gohtml"))))
func DonorHeader(client DonorHeaderClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		donorID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			return sirius.StatusError{Code: http.StatusBadRequest}
		}

		ctx := getContext(r)

		person, err := client.Person(ctx, donorID)
		if err != nil {
			return err
		}
		if person.Parent != nil {
			return RedirectError(fmt.Sprintf("/donor/%d", person.Parent.ID))
		}

		cases, err := client.CasesByDonor(ctx, donorID)
		if err != nil {
			return err
		}

		// Non-fatal: calendar renders without bank holiday markers if fetch fails.
		bankHolidays, _ := bhCache.get()
		today := time.Now().UTC().Truncate(24 * time.Hour)

		// Default calculator state mirrors PanelCalendarCtrl initialisation.
		calc := WorkingDaysVars{
			XSRFToken: ctx.XSRFToken,
			DonorID:   donorID,
			Date1:     today.Format(dateLayout),
			Date2:     today.AddDate(0, 0, 7).Format(dateLayout),
			Mode:      "difference",
		}
		calc = calculateWorkingDays(calc, bankHolidays)

		return tmpl(w, DonorHeaderData{
			XSRFToken:  ctx.XSRFToken,
			Person:     person,
			Cases:      cases,
			Months:     buildCalendarMonths(bankHolidays, today),
			Calculator: calc,
			DonorID:    donorID,
		})
	}
}

// DonorHeaderWorkingDays handles HTMX POST requests from the working-days
// calculator. Always returns the workingDaysCalculator partial fragment.
//
// Wire in server.go:
//
//	mux.Handle("/working-days", wrap(DonorHeaderWorkingDays(templates.Get("working-days-partial.gohtml"))))
func DonorHeaderWorkingDays(tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if err := r.ParseForm(); err != nil {
			return err
		}

		ctx := getContext(r)
		donorID, _ := strconv.Atoi(r.FormValue("donorId"))
		diff, _ := strconv.Atoi(r.FormValue("difference"))

		req := WorkingDaysVars{
			XSRFToken:  ctx.XSRFToken,
			DonorID:    donorID,
			Date1:      r.FormValue("date1"),
			Date2:      r.FormValue("date2"),
			Mode:       r.FormValue("mode"),
			Difference: diff,
		}

		bankHolidays, _ := bhCache.get()
		return tmpl(w, calculateWorkingDays(req, bankHolidays))
	}
}
