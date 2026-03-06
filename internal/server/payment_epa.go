package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type PaymentEpaClient interface {
	UpdateEpa(ctx sirius.Context, caseId int, epa sirius.Case) error
	Case(ctx sirius.Context, id int) (sirius.Case, error)
}

type PaymentEpaData struct {
	XSRFToken string
	Case      sirius.Case
	Success   bool
	Error     sirius.ValidationError
	Title     string
}

func PaymentEpa(client PaymentEpaClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		caseId, err := strToIntOrStatusError(r.FormValue("caseId"))
		if err != nil {
			return err
		}

		ctx := getContext(r)

		caseitem, err := client.Case(ctx, caseId)
		if err != nil {
			return err
		}

		data := PaymentEpaData{
			XSRFToken: ctx.XSRFToken,
			Title:     "Step 4: payment",
			Case:      caseitem,
		}

		if r.FormValue("isEditing") == "true" {
			data.Title = "Payment"
		}

		if r.Method == http.MethodPost {
			epa := sirius.Case{
				PaymentByCheque:  r.FormValue("paymentByCheque") == "true",
				PaymentExemption: r.FormValue("paymentExemption") == "true",
				PaymentDate:      postFormDateString(r, "paymentDate"),
			}

			err := client.UpdateEpa(ctx, caseId, epa)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
			} else if err != nil {
				return err
			}
			return RedirectError(fmt.Sprintf("/edit-epa?caseId=%d", caseId))
		}

		return tmpl(w, data)
	}
}
