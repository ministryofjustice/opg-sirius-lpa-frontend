package server

import (
	"fmt"
	"net/http"

	"golang.org/x/sync/errgroup"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type UpdateFeeDecisionClient interface {
	Case(sirius.Context, int) (sirius.Case, error)
	GetFeeDecision(ctx sirius.Context, feeDecisionID int) (sirius.FeeDecision, error)
	UpdateFeeDecision(ctx sirius.Context, caseID int, feeDecisionID int, decisionType string, decisionReason string, decisionDate sirius.DateString) error
	RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error)
}

type updateFeeDecisionData struct {
	XSRFToken      string
	Error          sirius.ValidationError
	IsPartial      bool
	IsEditing      bool
	Case           sirius.Case
	DecisionTypes  []sirius.RefDataItem
	DecisionType   string
	DecisionReason string
	DecisionDate   sirius.DateString
	FeeDecisionID  int
	ReturnUrl      string
	HtmxRedirect   string
}

func UpdateFeeDecision(client UpdateFeeDecisionClient, tmpl template.Template) Handler {
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
		group, groupCtx := errgroup.WithContext(ctx.Context)
		data := updateFeeDecisionData{
			XSRFToken:      ctx.XSRFToken,
			IsPartial:      ctx.IsPartial,
			IsEditing:      true,
			FeeDecisionID:  feeDecisionID,
		}

		var feeDecision sirius.FeeDecision

		group.Go(func() error {
			var err error
			data.Case, err = client.Case(ctx.With(groupCtx), caseID)
			return err
		})

		group.Go(func() error {
			var err error
			feeDecision, err = client.GetFeeDecision(ctx.With(groupCtx), feeDecisionID)
			return err
		})

		group.Go(func() error {
			var err error
			data.DecisionTypes, err = client.RefDataByCategory(ctx.With(groupCtx), sirius.FeeDecisionTypeCategory)
			return err
		})

		if err := group.Wait(); err != nil {
			return err
		}

		data.DecisionType = feeDecision.DecisionType.Key()
		data.DecisionReason = feeDecision.DecisionReason
		data.DecisionDate = feeDecision.DecisionDate

		if data.Case.CaseType == "DIGITAL_LPA" {
			data.ReturnUrl = fmt.Sprintf("/lpa/%s/payments", data.Case.UID)
		} else {
			data.ReturnUrl = fmt.Sprintf("/payments/%d", data.Case.ID)
		}

		if r.Method == http.MethodPost {
			data.DecisionType = postFormString(r, "decisionType")
			data.DecisionReason = postFormString(r, "decisionReason")
			data.DecisionDate = postFormDateString(r, "decisionDate")

			err := client.UpdateFeeDecision(
				ctx,
				caseID,
				feeDecisionID,
				data.DecisionType,
				data.DecisionReason,
				data.DecisionDate,
			)
			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
				return tmpl(w, data)
			} else if err != nil {
				return err
			}

			SetFlash(w, FlashNotification{Title: "Decline a fee reduction updated"})
			if ctx.IsPartial {
				data.HtmxRedirect = data.ReturnUrl
				return tmpl(w, data)
			}
			return RedirectError(data.ReturnUrl)
		}

		return tmpl(w, data)
	}
}
