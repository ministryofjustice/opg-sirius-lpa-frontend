package server

import (
	"github.com/go-chi/chi/v5"
	"golang.org/x/sync/errgroup"
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
	CaseSummary(ctx sirius.Context, uid string) (sirius.CaseSummary, error)
}

type getPaymentsData struct {
	XSRFToken string

	CaseSummary       sirius.CaseSummary
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
		ctx := getContext(r)
		group, groupCtx := errgroup.WithContext(ctx.Context)
		data := getPaymentsData{
			XSRFToken:   ctx.XSRFToken,
			CaseSummary: sirius.CaseSummary{},
		}

		var caseID int
		var err error

		uid := chi.URLParam(r, "uid")
		if uid != "" {
			data.CaseSummary, err = client.CaseSummary(ctx, uid)
			if err != nil {
				return err
			}
			caseID = data.CaseSummary.DigitalLpa.SiriusData.ID
		} else {
			caseID, err = strconv.Atoi(chi.URLParam(r, "id"))
			if err != nil {
				return err
			}
		}

		group.Go(func() error {
			data.Case, err = client.Case(ctx.With(groupCtx), caseID)
			if err != nil {
				return err
			}
			return nil
		})

		group.Go(func() error {
			payments, err := client.Payments(ctx, caseID)
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
						data.FeeReductions = append(data.FeeReductions, p)
					} else {
						data.Payments = append(data.Payments, p)
						data.TotalPaid = data.TotalPaid + p.Amount
					}
				}
				totalPaidAndReductions = totalPaidAndReductions + p.Amount
			}

			outstandingFeeOrRefund := 8200 - totalPaidAndReductions
			if outstandingFeeOrRefund < 0 {
				data.RefundAmount = outstandingFeeOrRefund * -1 /*convert to pos num for display*/
			} else {
				data.OutstandingFee = outstandingFeeOrRefund
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

		group.Go(func() error {
			data.ReferenceTypes, err = client.RefDataByCategory(ctx.With(groupCtx), sirius.PaymentReferenceType)
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

		group.Go(func() error {
			user, err := client.GetUserDetails(ctx.With(groupCtx))
			if err != nil {
				return err
			}
			data.IsReducedFeesUser = user.HasRole("Reduced Fees User")
			return nil
		})

		if err := group.Wait(); err != nil {
			return err
		}

		data.FlashMessage, _ = GetFlash(w, r)

		return tmpl(w, data)
	}
}
