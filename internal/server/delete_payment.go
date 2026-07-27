package server

import (
	"fmt"
	"net/http"

	"golang.org/x/sync/errgroup"

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
	ReturnUrl         string
	HtmxRedirect      string
}

func DeletePayment(client DeletePaymentClient, tmpl template.Template, partialTmpl template.Template) Handler {
	return func(pageVars PageVars, w http.ResponseWriter, r *http.Request) error {
		id, err := strToIntOrStatusError(r.FormValue("id"))
		if err != nil {
			return err
		}

		ctx := getContext(r)
		group, groupCtx := errgroup.WithContext(ctx.Context)

		p, err := client.PaymentByID(ctx, id)
		if err != nil {
			return err
		}

		data := deletePaymentData{
			XSRFToken: ctx.XSRFToken,
			Payment:   p,
		}

		group.Go(func() error {
			data.Case, err = client.Case(ctx.With(groupCtx), p.Case.ID)
			if err != nil {
				return err
			}

			return nil
		})

		group.Go(func() error {
			data.FeeReductionTypes, err = client.RefDataByCategory(ctx.With(groupCtx), sirius.FeeReductionTypeCategory)
			if err != nil {
				return err
			}

			return nil
		})

		if err := group.Wait(); err != nil {
			return err
		}

		if data.Case.CaseType == "DIGITAL_LPA" {
			data.ReturnUrl = fmt.Sprintf("/lpa/%s/payments", data.Case.UID)
		} else {
			data.ReturnUrl = fmt.Sprintf("/payments/%d", p.Case.ID)
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
				Title: fmt.Sprintf("%s deleted", item),
			})
			if r.Header.Get("HX-Request") == "true" {
				data.HtmxRedirect = data.ReturnUrl
				return partialTmpl(w, data)
			}
			return RedirectError(data.ReturnUrl)
		}

		if r.Header.Get("HX-Request") == "true" {
			return partialTmpl(w, data)
		}

		return tmpl(w, data)
	}
}
