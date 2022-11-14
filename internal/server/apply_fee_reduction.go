package server

import (
	"fmt"
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
		data := applyFeeReductionData{
			XSRFToken:        ctx.XSRFToken,
			PaymentEvidence:  postFormString(r, "paymentEvidence"),
			FeeReductionType: postFormString(r, "feeReductionType"),
			PaymentDate:      postFormDateString(r, "paymentDate"),
		}

		data.Case, err = client.Case(ctx, caseID)
		if err != nil {
			return err
		}

		data.FeeReductionTypes, err = client.RefDataByCategory(ctx, sirius.FeeReductionTypeCategory)
		if err != nil {
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
				return RedirectError(fmt.Sprintf("/payments?id=%d", caseID))
			}
		}

		return tmpl(w, data)
	}
}
