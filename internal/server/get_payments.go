package server

import (
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type GetPaymentsClient interface {
	RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error)
	Payments(ctx sirius.Context, id int) ([]sirius.Payment, error)
	Case(sirius.Context, int) (sirius.Case, error)
}

type getPaymentsData struct {
	XSRFToken string

	Case           sirius.Case
	Payments       []sirius.Payment
	Refunds        []sirius.Payment
	PaymentSources []sirius.RefDataItem
	TotalPaid      int
	TotalRefunds   int
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

		data.PaymentSources, err = client.RefDataByCategory(ctx, sirius.PaymentSourceCategory)
		if err != nil {
			return err
		}

		for _, p := range payments {
			if p.Amount < 0 {
				data.Refunds = append(data.Refunds, p)
				data.TotalRefunds = data.TotalRefunds + (p.Amount * -1)
			} else {
				data.Payments = append(data.Payments, p)
				data.TotalPaid = data.TotalPaid + p.Amount
			}
		}

		return tmpl(w, data)
	}
}
