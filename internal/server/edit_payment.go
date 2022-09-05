package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"net/http"
	"regexp"
	"strconv"
)

type EditPaymentClient interface {
	EditPayment(ctx sirius.Context, payment sirius.Payment) error
	Case(sirius.Context, int) (sirius.Case, error)
	PaymentByID(ctx sirius.Context, id int) (sirius.Payment, error)
}

type editPaymentData struct {
	XSRFToken string
	Success   bool
	Error     sirius.ValidationError

	Case        sirius.Case
	PaymentID   int
	Amount      string
	Source      string
	PaymentDate sirius.DateString
}

func EditPayment(client EditPaymentClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		caseID, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}
		paymentID, err := strconv.Atoi(r.FormValue("payment"))
		if err != nil {
			return err
		}

		ctx := getContext(r)
		p, err := client.PaymentByID(ctx, paymentID)
		if err != nil {
			return err
		}

		data := editPaymentData{
			XSRFToken:   ctx.XSRFToken,
			PaymentID:   paymentID,
			Amount:      fmt.Sprintf("%.2f", sirius.PenceToPounds(p.Amount)),
			Source:      p.Source,
			PaymentDate: p.PaymentDate,
		}

		data.Case, err = client.Case(ctx, caseID)
		if err != nil {
			return err
		}

		if r.Method == http.MethodPost {
			data.Amount = postFormString(r, "amount")
			data.Source = postFormString(r, "source")
			data.PaymentDate = postFormDateString(r, "paymentDate")

			m, err := regexp.Match(`^\d+\.\d{2}$`, []byte(data.Amount))
			if err != nil {
				return err
			}

			if !m {
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

			paymentEdit := sirius.Payment{
				ID:          paymentID,
				Amount:      amountInPence,
				Source:      data.Source,
				PaymentDate: data.PaymentDate,
			}

			err = client.EditPayment(ctx, paymentEdit)
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
