package server

import (
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"net/http"
	"strconv"
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
	PaymentSources    []sirius.RefDataItem
	User              sirius.User
	IsReducedFeesUser bool
	TotalPaid         int
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
		data.Payments = payments

		data.PaymentSources, err = client.RefDataByCategory(ctx, sirius.PaymentSourceCategory)
		if err != nil {
			return err
		}

		total := 0
		for _, p := range payments {
			total = total + p.Amount
		}
		data.TotalPaid = total

		user, err := client.GetUserDetails(ctx)
		if err != nil {
			return err
		}

		data.IsReducedFeesUser = user.HasRole("Reduced Fees User")

		return tmpl(w, data)
	}
}
