package server

import (
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"net/http"
	"strconv"
)

type AddPaymentClient interface {
	AddPayment(ctx sirius.Context, caseID int, amount int, source string, paymentDate sirius.DateString) error
	Case(sirius.Context, int) (sirius.Case, error)
}

type addPaymentData struct {
	XSRFToken string
	Success   bool
	Error     sirius.ValidationError

	Case        sirius.Case
	Amount      string
	Source      string
	PaymentDate sirius.DateString
}

func AddPayment(client AddPaymentClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		caseID, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}

		ctx := getContext(r)
		data := addPaymentData{
			XSRFToken:   ctx.XSRFToken,
			Amount:      postFormString(r, "amount"),
			Source:      postFormString(r, "source"),
			PaymentDate: postFormDateString(r, "paymentDate"),
		}

		data.Case, err = client.Case(ctx, caseID)
		if err != nil {
			return err
		}

		if r.Method == http.MethodPost {
			if !sirius.IsAmountValid(data.Amount) {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = sirius.ValidationError{
					Field: sirius.FieldErrors{
						"amount": {"reason": "Please enter the amount to 2 decimal places"},
					},
				}
				return tmpl(w, data)
			}

			amountFloat, err := strconv.ParseFloat(data.Amount, 64)
			if err != nil {
				return err
			}

			amountInPence := sirius.PoundsToPence(amountFloat)

			err = client.AddPayment(ctx, caseID, amountInPence, data.Source, data.PaymentDate)
			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
			} else if err != nil {
				return err
			} else {
				data.Success = true
			}
		}

		return tmpl(w, data)
	}
}
