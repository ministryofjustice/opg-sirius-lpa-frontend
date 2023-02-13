package server

import (
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type GetPaymentsClient interface {
	RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error)
	Payments(ctx sirius.Context, id int) ([]sirius.Payment, error)
	Case(sirius.Context, int) (sirius.Case, error)
	GetUserDetails(sirius.Context) (sirius.User, error)
}

type getPaymentsData struct {
	XSRFToken string

	Case              sirius.Case
	Payments          []sirius.Payment
	FeeReductions     []sirius.Payment
	Refunds           []sirius.Payment
	PaymentSources    []sirius.RefDataItem
	ReferenceTypes    []sirius.RefDataItem
	FeeReductionTypes []sirius.RefDataItem
	IsReducedFeesUser bool
	TotalPaid         int
	TotalRefunds      int
	OutstandingFee    int
	RefundAmount      int
	FlashMessage      FlashNotification
}

func GetPayments(client GetPaymentsClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		caseID, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}

		ctx := getContext(r)
		data := getPaymentsData{XSRFToken: ctx.XSRFToken}

		data.Case, err = client.Case(ctx, caseID)
		if err != nil {
			return err
		}

		payments, err := client.Payments(ctx, caseID)
		if err != nil {
			return err
		}

		data.PaymentSources, err = client.RefDataByCategory(ctx, sirius.PaymentSourceCategory)
		if err != nil {
			return err
		}

		data.ReferenceTypes, err = client.RefDataByCategory(ctx, sirius.PaymentReferenceType)
		if err != nil {
			return err
		}

		data.FeeReductionTypes, err = client.RefDataByCategory(ctx, sirius.FeeReductionTypeCategory)
		if err != nil {
			return err
		}

		totalPaidAndReductions := 0
		for _, p := range payments {
			if p.Amount < 0 {
				data.Refunds = append(data.Refunds, p)
				data.TotalRefunds = data.TotalRefunds + (p.Amount * -1)
			} else {
				if p.Source == sirius.FeeReductionSource {
					data.FeeReductions = append(data.FeeReductions, p) // multiple fee reductions?
				} else {
					data.Payments = append(data.Payments, p)
					data.TotalPaid = data.TotalPaid + p.Amount
				}
				totalPaidAndReductions = totalPaidAndReductions + p.Amount
			}
		}

		outstandingFeeOrRefund := 8200 - totalPaidAndReductions + data.TotalRefunds
		if outstandingFeeOrRefund < 0 {
			data.RefundAmount = outstandingFeeOrRefund * -1 /*convert to pos num for display*/
		} else {
			data.OutstandingFee = outstandingFeeOrRefund
		}

		user, err := client.GetUserDetails(ctx)
		if err != nil {
			return err
		}

		data.IsReducedFeesUser = user.HasRole("Reduced Fees User")

		data.FlashMessage, _ = GetFlash(w, r)

		return tmpl(w, data)
	}
}
