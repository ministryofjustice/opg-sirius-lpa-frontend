package server

import (
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"net/http"
	"strconv"
)

type DeletePaymentClient interface {
	PaymentByID(ctx sirius.Context, id int) (sirius.Payment, error)
	Case(sirius.Context, int) (sirius.Case, error)
	DeletePayment(ctx sirius.Context, paymentID int) error
}

type deletePaymentData struct {
	XSRFToken string
	Success   bool
	Payment   sirius.Payment
	Case      sirius.Case
}

func DeletePayment(client DeletePaymentClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}

		ctx := getContext(r)

		p, err := client.PaymentByID(ctx, id)
		if err != nil {
			return err
		}

		data := deletePaymentData{
			XSRFToken: ctx.XSRFToken,
			Payment:   p,
		}

		data.Case, err = client.Case(ctx, p.Case.ID)
		if err != nil {
			return err
		}

		if r.Method == http.MethodPost {
			err = client.DeletePayment(ctx, id)
			if err != nil {
				return err
			}

			data.Success = true
		}

		return tmpl(w, data)
	}
}
