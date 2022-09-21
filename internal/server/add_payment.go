package server

import (
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"net/http"
	"strconv"
)

type AddPaymentClient interface {
	RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error)
	AddPayment(ctx sirius.Context, caseID int, amount int, source string, paymentDate sirius.DateString, feeReductionType string, paymentEvidence string, appliedDate sirius.DateString) error
	Case(sirius.Context, int) (sirius.Case, error)
}

type addPaymentData struct {
	XSRFToken string
	Success   bool
	Error     sirius.ValidationError

	Case              sirius.Case
	Amount            string
	Source            string
	PaymentEvidence   string
	FeeReductionType  string
	PaymentDate       sirius.DateString
	AppliedDate       sirius.DateString
	PaymentSources    []sirius.RefDataItem
	FeeReductionTypes []sirius.RefDataItem
}

const FeeReductionSource = "FEE_REDUCTION"

func AddPayment(client AddPaymentClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		caseID, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}

		ctx := getContext(r)
		data := addPaymentData{
			XSRFToken:        ctx.XSRFToken,
			Amount:           postFormString(r, "amount"),
			Source:           postFormString(r, "source"),
			PaymentDate:      postFormDateString(r, "paymentDate"),
			PaymentEvidence:  postFormString(r, "paymentEvidence"),
			FeeReductionType: postFormString(r, "feeReductionType"),
			AppliedDate:      postFormDateString(r, "appliedDate"),
		}

		data.Case, err = client.Case(ctx, caseID)
		if err != nil {
			return err
		}

		data.PaymentSources, err = client.RefDataByCategory(ctx, sirius.PaymentSourceCategory)
		if err != nil {
			return err
		}

		data.FeeReductionTypes, err = client.RefDataByCategory(ctx, sirius.FeeReductionTypeCategory)
		if err != nil {
			return err
		}

		if r.Method == http.MethodPost {
			if r.URL.Path == "/apply-fee-reduction" {
				data.Source = FeeReductionSource
			}

			if !sirius.IsAmountValid(data.Amount) && data.Source != FeeReductionSource {
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

			var amountInPence int
			if data.Source != FeeReductionSource {
				amountFloat, err := strconv.ParseFloat(data.Amount, 64)
				if err != nil {
					return err
				}

				amountInPence = sirius.PoundsToPence(amountFloat)
			}

			err = client.AddPayment(ctx, caseID, amountInPence, data.Source, data.PaymentDate, data.FeeReductionType, data.PaymentEvidence, data.AppliedDate)
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
