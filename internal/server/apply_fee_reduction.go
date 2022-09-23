package server

import (
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"net/http"
	"strconv"
)

type ApplyFeeReductionClient interface {
	RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error)
	ApplyFeeReduction(ctx sirius.Context, caseID int, source string, feeReductionType string, paymentEvidence string, paymentDate sirius.DateString) error
	Case(sirius.Context, int) (sirius.Case, error)
}

type applyFeeReductionData struct {
	XSRFToken string
	Success   bool
	Error     sirius.ValidationError

	Case              sirius.Case
	Source            string
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
			data.Source = "FEE_REDUCTION"

			err = client.ApplyFeeReduction(ctx, caseID, data.Source, data.FeeReductionType, data.PaymentEvidence, data.PaymentDate)
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
