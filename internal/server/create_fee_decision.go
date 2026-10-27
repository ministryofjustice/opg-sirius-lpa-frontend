package server

import (
	"fmt"
	"net/http"

	"golang.org/x/sync/errgroup"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type CreateFeeDecisionClient interface {
	RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error)
	CreateFeeDecision(ctx sirius.Context, caseID int, decisionType string, decisionReason string, decisionDate sirius.DateString) error
	Case(sirius.Context, int) (sirius.Case, error)
}

type createFeeDecisionData struct {
	XSRFToken      string
	Error          sirius.ValidationError
	IsEditing      bool
	Case           sirius.Case
	DecisionTypes  []sirius.RefDataItem
	DecisionType   string
	DecisionReason string
	DecisionDate   sirius.DateString
	ReturnUrl      string
}

func CreateFeeDecision(client CreateFeeDecisionClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		caseID, err := strToIntOrStatusError(r.FormValue("id"))
		if err != nil {
			return err
		}

		ctx := getContext(r)
		group, groupCtx := errgroup.WithContext(ctx.Context)
		data := createFeeDecisionData{
			XSRFToken:      ctx.XSRFToken,
			DecisionType:   postFormString(r, "decisionType"),
			DecisionReason: postFormString(r, "decisionReason"),
			DecisionDate:   postFormDateString(r, "decisionDate"),
		}

		group.Go(func() error {
			data.Case, err = client.Case(ctx.With(groupCtx), caseID)
			if err != nil {
				return err
			}

			return nil
		})

		group.Go(func() error {
			data.DecisionTypes, err = client.RefDataByCategory(ctx.With(groupCtx), sirius.FeeDecisionTypeCategory)
			if err != nil {
				return err
			}

			return nil
		})

		if err := group.Wait(); err != nil {
			return err
		}

		isDigitalLpa := data.Case.CaseType == "DIGITAL_LPA"

		if isDigitalLpa {
			data.ReturnUrl = fmt.Sprintf("/lpa/%s/payments", data.Case.UID)
		} else {
			data.ReturnUrl = fmt.Sprintf("/payments/%d", caseID)
		}

		if r.Method == http.MethodPost {
			err = client.CreateFeeDecision(ctx, caseID, data.DecisionType, data.DecisionReason, data.DecisionDate)
			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
			} else if err != nil {
				return err
			} else {

				title := "Decline a fee reduction recorded"

				if isDigitalLpa {
					title = "Fee decision added"
				}

				SetFlash(w, FlashNotification{
					Title: title,
				})

				return RedirectError(data.ReturnUrl)
			}
		}

		return tmpl(w, data)
	}
}
