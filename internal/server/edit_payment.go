package server

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type EditPaymentClient interface {
	EditPayment(ctx sirius.Context, paymentID int, payment sirius.Payment) error
	Case(sirius.Context, int) (sirius.Case, error)
	PaymentByID(ctx sirius.Context, id int) (sirius.Payment, error)
	RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error)
}

type editPaymentData struct {
	XSRFToken string
	Error     sirius.ValidationError

	Case           sirius.Case
	PaymentID      int
	Amount         string
	Source         string
	PaymentDate    sirius.DateString
	PaymentSources []sirius.RefDataItem
}

func EditPayment(client EditPaymentClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		paymentID, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}

		ctx := getContext(r)
		group, groupCtx := errgroup.WithContext(ctx.Context)
		data := editPaymentData{
			XSRFToken: ctx.XSRFToken,
		}

		group.Go(func() error {
			p, err := client.PaymentByID(ctx, paymentID)
			if err != nil {
				return err
			}

			data.PaymentID = paymentID
			data.Amount = fmt.Sprintf("%.2f", sirius.PenceToPounds(p.Amount))
			data.Source = p.Source
			data.PaymentDate = p.PaymentDate

			data.Case, err = client.Case(ctx.With(groupCtx), p.Case.ID)
			if err != nil {
				return err
			}
			return nil
		})

		group.Go(func() error {
			data.PaymentSources, err = client.RefDataByCategory(ctx.With(groupCtx), sirius.PaymentSourceCategory)
			if err != nil {
				return err
			}
			return nil
		})

		if err := group.Wait(); err != nil {
			return err
		}

		if r.Method == http.MethodPost {
			data.Amount = postFormString(r, "amount")
			data.Source = postFormString(r, "source")
			data.PaymentDate = postFormDateString(r, "paymentDate")

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

			paymentEdit := sirius.Payment{
				Amount:      amountInPence,
				Source:      data.Source,
				PaymentDate: data.PaymentDate,
			}

			err = client.EditPayment(ctx, paymentID, paymentEdit)
			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
			} else if err != nil {
				return err
			} else {
				SetFlash(w, FlashNotification{Title: "Payment saved"})
				return RedirectError(fmt.Sprintf("/payments?id=%d", data.Case.ID))
			}
		}

		return tmpl(w, data)
	}
}
