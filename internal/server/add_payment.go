package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type AddPaymentClient interface {
	RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error)
	AddPayment(ctx sirius.Context, caseID int, amount int, source string, paymentDate sirius.DateString) error
	Case(sirius.Context, int) (sirius.Case, error)
}

type addPaymentData struct {
	XSRFToken string
	Error     sirius.ValidationError

	Case           sirius.Case
	Amount         string
	Source         string
	PaymentDate    sirius.DateString
	PaymentSources []sirius.RefDataItem
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

		data.PaymentSources, err = client.RefDataByCategory(ctx, sirius.PaymentSourceCategory)
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
				if data.Source == "" {
					data.Error.Field["source"] = map[string]string{
						"reason": "Value is required and can't be empty",
					}
				}
				if data.PaymentDate == "" {
					data.Error.Field["paymentDate"] = map[string]string{
						"reason": "Value is required and can't be empty",
					}
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
				SetFlash(w, FlashNotification{
					Title: "Payment added",
				})
				return RedirectError(fmt.Sprintf("/payments?id=%d", caseID))
			}
		}

		return tmpl(w, data)
	}
}
