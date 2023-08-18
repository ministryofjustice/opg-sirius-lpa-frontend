package server

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type EditFeeReductionClient interface {
	RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error)
	EditPayment(ctx sirius.Context, paymentID int, payment sirius.Payment) error
	PaymentByID(ctx sirius.Context, id int) (sirius.Payment, error)
	Case(sirius.Context, int) (sirius.Case, error)
}

type editFeeReductionData struct {
	XSRFToken string
	Error     sirius.ValidationError

	Case              sirius.Case
	PaymentID         int
	FeeReductionTypes []sirius.RefDataItem
	FeeReduction      sirius.Payment
}

func EditFeeReduction(client EditFeeReductionClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		paymentID, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}

		ctx := getContext(r)
		group, groupCtx := errgroup.WithContext(ctx.Context)
		data := editFeeReductionData{
			XSRFToken: ctx.XSRFToken,
			PaymentID: paymentID,
		}

		group.Go(func() error {
			feeReduction, err := client.PaymentByID(ctx.With(groupCtx), paymentID)
			if err != nil {
				return err
			}
			data.FeeReduction = feeReduction

			data.Case, err = client.Case(ctx, feeReduction.Case.ID)
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
			data.FeeReduction.PaymentEvidence = postFormString(r, "paymentEvidence")
			data.FeeReduction.FeeReductionType = postFormString(r, "feeReductionType")
			data.FeeReduction.PaymentDate = postFormDateString(r, "paymentDate")

			err = client.EditPayment(ctx, paymentID, data.FeeReduction)
			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
			} else if err != nil {
				return err
			} else {
				SetFlash(w, FlashNotification{Title: "Fee reduction edited"})
				return RedirectError(fmt.Sprintf("/payments/%d", data.FeeReduction.Case.ID))
			}
		}

		return tmpl(w, data)
	}
}
