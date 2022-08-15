package server

import (
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"net/http"
	"strconv"
)

type GetPaymentsClient interface {
	Payments(ctx sirius.Context, id int) ([]sirius.Payment, error)
	Case(sirius.Context, int) (sirius.Case, error)
}

type getPaymentsData struct {
	XSRFToken string
	Success   bool
	Error     sirius.ValidationError

	Case      sirius.Case
	Payments  []sirius.Payment
	TotalPaid float64
}

func GetPayments(client GetPaymentsClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		caseID, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}

		ctx := getContext(r)
		data := getPaymentsData{XSRFToken: ctx.XSRFToken}

		data.Case, err = client.Case(ctx, caseID)
		if err != nil {
			return err
		}

		payments, err := client.Payments(ctx, caseID)
		if err != nil {
			return err
		}

		for _, p := range payments {
			amount, err := strconv.ParseFloat(string(p.Amount), 64)
			if err != nil {
				return err
			}
			data.TotalPaid = data.TotalPaid + amount
		}

		data.Payments = payments

		return tmpl(w, data)
	}
}
