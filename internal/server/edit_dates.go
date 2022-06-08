package server

import (
	"fmt"
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

	Case sirius.Case
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
				CancellationDate: sirius.DateString(r.FormValue("cancellationDate")),
				DispatchDate:     sirius.DateString(r.FormValue("dispatchDate")),
				DueDate:          sirius.DateString(r.FormValue("dueDate")),
				InvalidDate:      sirius.DateString(r.FormValue("invalidDate")),
				ReceiptDate:      sirius.DateString(r.FormValue("receiptDate")),
				RegistrationDate: sirius.DateString(r.FormValue("registrationDate")),
				RejectedDate:     sirius.DateString(r.FormValue("rejectedDate")),
				WithdrawnDate:    sirius.DateString(r.FormValue("withdrawnDate")),
			}

			err = client.EditDates(ctx, caseID, caseType, dates)
			if err != nil {
				return err
			}

			data.Success = true
		}

		caseitem, err := client.Case(ctx, caseID)
		if err != nil {
			return err
		}

		data.Case = caseitem
		data.Entity = fmt.Sprintf("%s %s", caseitem.CaseType, caseitem.UID)

		return tmpl(w, data)
	}
}
