package server

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type ApplyFeeReductionClient interface {
	RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error)
	ApplyFeeReduction(ctx sirius.Context, caseID int, feeReductionType string, paymentEvidence string, paymentDate sirius.DateString) error
	Case(sirius.Context, int) (sirius.Case, error)
}

type applyFeeReductionData struct {
	XSRFToken string
	Error     sirius.ValidationError

	Case              sirius.Case
	PaymentEvidence   string
	FeeReductionType  string
	PaymentDate       sirius.DateString
	FeeReductionTypes []sirius.RefDataItem
}

func ApplyFeeReduction(client ApplyFeeReductionClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		caseID, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}

		ctx := getContext(r)
		group, groupCtx := errgroup.WithContext(ctx.Context)
		data := applyFeeReductionData{
			XSRFToken:        ctx.XSRFToken,
			PaymentEvidence:  postFormString(r, "paymentEvidence"),
			FeeReductionType: postFormString(r, "feeReductionType"),
			PaymentDate:      postFormDateString(r, "paymentDate"),
		}

		group.Go(func() error {
			data.Case, err = client.Case(ctx.With(groupCtx), caseID)
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
			err = client.ApplyFeeReduction(ctx, caseID, data.FeeReductionType, data.PaymentEvidence, data.PaymentDate)
			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
			} else if err != nil {
				return err
			} else {
				SetFlash(w, FlashNotification{
					Title: fmt.Sprintf("%s approved", translateRefData(data.FeeReductionTypes, data.FeeReductionType)),
				})
				return RedirectError(fmt.Sprintf("/payments/%d", caseID))
			}
		}

		return tmpl(w, data)
	}
}
