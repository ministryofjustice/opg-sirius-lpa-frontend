package server

import (
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"net/http"
	"strconv"
)

type EditFeeReductionClient interface {
	RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error)
	EditPayment(ctx sirius.Context, paymentID int, payment sirius.Payment) error
	PaymentByID(ctx sirius.Context, id int) (sirius.Payment, error)
	Case(sirius.Context, int) (sirius.Case, error)
}

type editFeeReductionData struct {
	XSRFToken string
	Success   bool
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
		feeReduction, err := client.PaymentByID(ctx, paymentID)
		if err != nil {
			return err
		}

		data := editFeeReductionData{
			XSRFToken:    ctx.XSRFToken,
			PaymentID:    paymentID,
			FeeReduction: feeReduction,
		}

		data.Case, err = client.Case(ctx, feeReduction.Case.ID)
		if err != nil {
			return err
		}

		data.FeeReductionTypes, err = client.RefDataByCategory(ctx, sirius.FeeReductionTypeCategory)
		if err != nil {
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
				data.Success = true
			}
		}

		return tmpl(w, data)
	}
}
