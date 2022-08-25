package server

import (
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"net/http"
	"regexp"
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
	Amount      int
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
			Source:      postFormString(r, "source"),
			PaymentDate: postFormDateString(r, "paymentDate"),
		}

		data.Case, err = client.Case(ctx, caseID)
		if err != nil {
			return err
		}

		if r.Method == http.MethodPost {
			amountString := r.FormValue("amount")
			m, err := regexp.Match(`^\d*\.\d{2}$`, []byte(amountString))
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

			amountFloat, err := strconv.ParseFloat(amountString, 64)
			if err != nil {
				return err
			}

			data.Amount = sirius.PoundsToPence(amountFloat)

			err = client.AddPayment(ctx, caseID, data.Amount, data.Source, data.PaymentDate)
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
