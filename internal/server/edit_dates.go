package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type EditDatesClient interface {
	Case(sirius.Context, int) (sirius.Case, error)
	EditDates(sirius.Context, int, sirius.CaseType, sirius.Dates) error
}

type editDatesData struct {
	XSRFToken string
	Entity    string
	Success   bool
	Error     sirius.ValidationError

	Dates            sirius.Dates
	DonorId          int
	CaseUid          string
	CaseType         string
	CaseId           int
	ReceiptDate      dob
	PaymentDate      dob
	FilingDate       dob
	DueDate          dob
	RegistrationDate dob
}

func EditDates(client EditDatesClient, tmpl template.Template, partialTmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		caseID, err := strToIntOrStatusError(r.FormValue("id"))
		if err != nil {
			return err
		}

		caseType, err := sirius.ParseCaseType(r.FormValue("case"))
		if err != nil {
			return err
		}

		ctx := getContext(r)
		data := editDatesData{
			XSRFToken: ctx.XSRFToken,
			CaseId:    caseID,
		}

		if r.Method == http.MethodPost {
			dates := sirius.Dates{
				CancellationDate: postFormDateString(r, "cancellationDate"),
				DispatchDate:     postFormDateString(r, "dispatchDate"),
				DueDate:          postFormDayMonthYear(r, "dueDate"),
				FilingDate:       postFormDayMonthYear(r, "filingDate"),
				InvalidDate:      postFormDateString(r, "invalidDate"),
				PaymentDate:      postFormDayMonthYear(r, "paymentDate"),
				ReceiptDate:      postFormDayMonthYear(r, "receiptDate"),
				RegistrationDate: postFormDayMonthYear(r, "registrationDate"),
				RejectedDate:     postFormDateString(r, "rejectedDate"),
				RevokedDate:      postFormDateString(r, "revokedDate"),
				WithdrawnDate:    postFormDateString(r, "withdrawnDate"),
			}

			err = client.EditDates(ctx, caseID, caseType, dates)
			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
				data.Dates = dates
			} else if err != nil {
				return err
			} else {
				data.Success = true
			}
		}

		caseitem, err := client.Case(ctx, caseID)
		if err != nil {
			return err
		}

		data.CaseUid = caseitem.UID
		data.CaseType = caseitem.CaseType
		if caseitem.Donor != nil {
			data.DonorId = caseitem.Donor.ID
		}

		if r.Method != http.MethodPost || data.Success {
			receiptDate, err := checkDobValue(caseitem.ReceiptDate)
			if err != nil {
				return err
			}
			paymentDate, err := checkDobValue(caseitem.PaymentDate)
			if err != nil {
				return err
			}
			filingDate, err := checkDobValue(caseitem.FilingDate)
			if err != nil {
				return err
			}
			dueDate, err := checkDobValue(caseitem.DueDate)
			if err != nil {
				return err
			}
			registrationDate, err := checkDobValue(caseitem.RegistrationDate)
			if err != nil {
				return err
			}
			data.Dates = sirius.Dates{
				CancellationDate: caseitem.CancellationDate,
				DispatchDate:     caseitem.DispatchDate,
				//DueDate:          caseitem.DueDate,
				InvalidDate: caseitem.InvalidDate,
				//PaymentDate:      caseitem.PaymentDate,
				//FilingDate:       caseitem.FilingDate,
				//ReceiptDate:      caseitem.ReceiptDate,
				//RegistrationDate: caseitem.RegistrationDate,
				RejectedDate:  caseitem.RejectedDate,
				RevokedDate:   caseitem.RevokedDate,
				WithdrawnDate: caseitem.WithdrawnDate,
			}
			data.ReceiptDate = receiptDate
			data.PaymentDate = paymentDate
			data.FilingDate = filingDate
			data.DueDate = dueDate
			data.RegistrationDate = registrationDate
		}
		data.Entity = caseitem.Summary()

		if r.Header.Get("HX-Request") == "true" {
			return partialTmpl(w, data)
		}

		return tmpl(w, data)
	}
}

func checkDobValue(date sirius.DateString) (dob, error) {
	if date == "" {
		return dob{}, nil
	}

	dateString, err := date.ToHyphenateDates()
	if err != nil {
		return dob{}, err
	}

	return parseDate(dateString)
}
