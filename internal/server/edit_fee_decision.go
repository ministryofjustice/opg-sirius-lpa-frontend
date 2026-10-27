package server

import (
	"fmt"
	"net/http"

	"golang.org/x/sync/errgroup"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type EditFeeDecisionClient interface {
	Case(sirius.Context, int) (sirius.Case, error)
	DeclineFeeReductions(ctx sirius.Context, id int) ([]sirius.DeclineFeeReductions, error)
	EditFeeDecision(ctx sirius.Context, caseID int, feeDecisionID int, decisionType string, decisionReason string, decisionDate sirius.DateString) error
	RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error)
}

type editFeeDecisionData struct {
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

func EditFeeDecision(client EditFeeDecisionClient, tmpl template.Template) Handler {
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
		data := editFeeDecisionData{
			XSRFToken:      ctx.XSRFToken,
			IsPartial:      ctx.IsPartial,
			IsEditing:      true,
			FeeDecisionID:  feeDecisionID,
			DecisionType:   postFormString(r, "decisionType"),
			DecisionReason: postFormString(r, "decisionReason"),
			DecisionDate:   postFormDateString(r, "decisionDate"),
		}

		var decisions []sirius.DeclineFeeReductions

		group.Go(func() error {
			var err error
			data.Case, err = client.Case(ctx.With(groupCtx), caseID)
			return err
		})

		group.Go(func() error {
			var err error
			decisions, err = client.DeclineFeeReductions(ctx.With(groupCtx), caseID)
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

		var found bool
		for _, decision := range decisions {
			if decision.ID == feeDecisionID {
				if r.Method != http.MethodPost {
					data.DecisionType = decision.DecisionType.Key()
					data.DecisionReason = decision.DecisionReason
					data.DecisionDate = decision.DecisionDate
				}
				found = true
				break
			}
		}
		if !found {
			return sirius.StatusError{Code: http.StatusNotFound}
		}

		if data.Case.CaseType == "DIGITAL_LPA" {
			data.ReturnUrl = fmt.Sprintf("/lpa/%s/payments", data.Case.UID)
		} else {
			data.ReturnUrl = fmt.Sprintf("/payments/%d", data.Case.ID)
		}

		if r.Method == http.MethodPost {
			err := client.EditFeeDecision(
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
