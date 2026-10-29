package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type DeleteFeeDecisionClient interface {
	DeleteFeeDecision(ctx sirius.Context, feeDecisionID int) error
	GetFeeDecision(ctx sirius.Context, feeDecisionID int) (sirius.FeeDecision, error)
}

type deleteFeeDecisionData struct {
	XSRFToken       string
	Error           sirius.ValidationError
	CaseID          int
	FeeDecision     sirius.FeeDecision
	ReturnUrl       string
	IsPartial       bool
	HtmxRedirect    string
	FeeDecisionID   int
	FeeDecisionType string
}

func DeleteFeeDecision(client DeleteFeeDecisionClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		caseID, err := strToIntOrStatusError(r.FormValue("caseId"))
		if err != nil {
			return err
		}

		feeDecisionID, err := strToIntOrStatusError(r.FormValue("feeDecisionId"))
		if err != nil {
			return err
		}

		ctx := getContext(r)
		data := deleteFeeDecisionData{
			XSRFToken: ctx.XSRFToken,
			CaseID:    caseID,
			IsPartial: ctx.IsPartial,
			ReturnUrl: fmt.Sprintf("/payments/%d", caseID),
		}

		data.FeeDecision, err = client.GetFeeDecision(ctx, feeDecisionID)
		if err != nil {
			return err
		}

		if data.FeeDecision.DecisionType == shared.DecisionTypeExemption {
			data.FeeDecisionType = "declined exemption"
		} else if data.FeeDecision.DecisionType == shared.DecisionTypeRemission {
			data.FeeDecisionType = "declined remission"
		}

		if r.Method == http.MethodPost {
			err = client.DeleteFeeDecision(ctx, feeDecisionID)
			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
			} else if err != nil {
				return err
			} else {
				SetFlash(w, FlashNotification{Title: "Successfully  deleted a declined fee reduction"})

				if ctx.IsPartial {
					data.HtmxRedirect = data.ReturnUrl
					return tmpl(w, data)
				}
				return RedirectError(data.ReturnUrl)
			}
		}

		return tmpl(w, data)
	}
}
