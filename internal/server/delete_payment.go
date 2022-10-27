package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type DeletePaymentClient interface {
	PaymentByID(ctx sirius.Context, id int) (sirius.Payment, error)
	Case(sirius.Context, int) (sirius.Case, error)
	DeletePayment(ctx sirius.Context, paymentID int) error
	RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error)
}

type deletePaymentData struct {
	XSRFToken         string
	Payment           sirius.Payment
	Case              sirius.Case
	FeeReductionTypes []sirius.RefDataItem
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

		data.FeeReductionTypes, err = client.RefDataByCategory(ctx, sirius.FeeReductionTypeCategory)
		if err != nil {
			return err
		}

		if r.Method == http.MethodPost {
			err = client.DeletePayment(ctx, id)
			if err != nil {
				return err
			}

			item := "Payment"
			if p.Source == sirius.FeeReductionSource {
				item = translateRefData(data.FeeReductionTypes, p.FeeReductionType)
			}

			SetFlash(w, FlashNotification{
				Title:       fmt.Sprintf("%s deleted", item),
				Description: "Please clear the task if you have completed it",
			})
			return RedirectError(fmt.Sprintf("/payments?id=%d", p.Case.ID))
		}

		return tmpl(w, data)
	}
}
