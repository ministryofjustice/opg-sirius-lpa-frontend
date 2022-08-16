package server

import (
	"fmt"
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
	TotalPaid string
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
		data.Payments = payments

		total := 0.00
		for _, p := range payments {
			amount, err := strconv.ParseFloat(string(p.Amount), 64)
			if err != nil {
				return err
			}
			total = total + amount
		}
		data.TotalPaid = fmt.Sprintf("%.2f", total)

		return tmpl(w, data)
	}
}
