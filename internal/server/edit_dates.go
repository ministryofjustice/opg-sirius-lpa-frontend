package server

import (
	"net/http"
	"strconv"

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

	Dates sirius.Dates
}

func EditDates(client EditDatesClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		caseID, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}

		caseType, err := sirius.ParseCaseType(r.FormValue("case"))
		if err != nil {
			return err
		}

		ctx := getContext(r)
		data := editDatesData{XSRFToken: ctx.XSRFToken}

		if r.Method == http.MethodPost {
			dates := sirius.Dates{
				CancellationDate: postFormDateString(r, "cancellationDate"),
				DispatchDate:     postFormDateString(r, "dispatchDate"),
				DueDate:          postFormDateString(r, "dueDate"),
				InvalidDate:      postFormDateString(r, "invalidDate"),
				ReceiptDate:      postFormDateString(r, "receiptDate"),
				RegistrationDate: postFormDateString(r, "registrationDate"),
				RejectedDate:     postFormDateString(r, "rejectedDate"),
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

		if r.Method != http.MethodPost || data.Success {
			data.Dates = sirius.Dates{
				CancellationDate: caseitem.CancellationDate,
				DispatchDate:     caseitem.DispatchDate,
				DueDate:          caseitem.DueDate,
				InvalidDate:      caseitem.InvalidDate,
				ReceiptDate:      caseitem.ReceiptDate,
				RegistrationDate: caseitem.RegistrationDate,
				RejectedDate:     caseitem.RejectedDate,
				WithdrawnDate:    caseitem.WithdrawnDate,
			}
		}
		data.Entity = caseitem.Summary()

		return tmpl(w, data)
	}
}
