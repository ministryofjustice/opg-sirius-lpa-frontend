package server

import (
	"fmt"
	"net/http"

	"golang.org/x/sync/errgroup"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type DeletePaymentAlertClient interface {
	PaymentByID(ctx sirius.Context, id int) (sirius.Payment, error)
	RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error)
}

type deletePaymentAlertData struct {
	XSRFToken      string
	Payment        sirius.Payment
	PaymentSources []sirius.RefDataItem
}

func DeletePaymentAlert(client DeletePaymentAlertClient, tmpl template.Template, fragmentTmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
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

		data := deletePaymentAlertData{
			XSRFToken: ctx.XSRFToken,
			Payment:   p,
		}

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

		if r.Header.Get("HX-Request") == "true" {
			return fragmentTmpl(w, data)
		}

		// Test 4: redirect back to payments page with confirm param so alert shows inline
		if r.FormValue("mode") == "inline" {
			return RedirectError(fmt.Sprintf("/payments/%d?confirm=%d#alert-container-%d", p.Case.ID, p.ID, p.ID))
		}

		return tmpl(w, data)
	}
}

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
			return RedirectError(fmt.Sprintf("/payments/%d", p.Case.ID))
		}

		return tmpl(w, data)
	}
}
