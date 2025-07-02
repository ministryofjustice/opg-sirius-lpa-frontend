package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type AttorneyDecisionsClient interface {
	CaseSummary(sirius.Context, string) (sirius.CaseSummary, error)
	ManageAttorneyDecisions(sirius.Context, string, []sirius.AttorneyDecisions) error
}

type formManageAttorneyDecisions struct {
	DecisionAttorneysUids                 []string `form:"decisionAttorney"`
	JointDecisionsCanBeMadeByAllAttorneys string   `form:"allAttorneysCanMakeDecisions"`
	SkipDecisionAttorney                  string   `form:"skipDecisionAttorney"`
}

type manageAttorneyDecisionsData struct {
	CaseSummary              sirius.CaseSummary
	ActiveAttorneys          []sirius.LpaStoreAttorney
	Form                     formManageAttorneyDecisions
	DecisionAttorneysDetails []SelectedAttorneyDetails
	Success                  bool
	Error                    sirius.ValidationError
	XSRFToken                string
}

func AttorneyDecisions(client AttorneyDecisionsClient, removeTmpl template.Template) Handler {

	return func(w http.ResponseWriter, r *http.Request) error {
		uid := r.PathValue("uid")
		ctx := getContext(r)

		caseSummary, err := client.CaseSummary(ctx, uid)

		if err != nil {
			return err
		}

		data := manageAttorneyDecisionsData{
			CaseSummary: caseSummary,
			XSRFToken:   ctx.XSRFToken,
			Error:       sirius.ValidationError{Field: sirius.FieldErrors{}},
		}

		lpa := data.CaseSummary.DigitalLpa

		for _, attorney := range lpa.LpaStoreData.Attorneys {
			if attorney.Status == shared.ActiveAttorneyStatus.String() {
				data.ActiveAttorneys = append(data.ActiveAttorneys, attorney)
			}
		}

		if r.Method == http.MethodPost {

			err = decoder.Decode(&data.Form, r.PostForm)
			if err != nil {
				return err
			}

			var attorneyDecisions []sirius.AttorneyDecisions

			if len(data.Form.DecisionAttorneysUids) > 0 {
				for _, att := range data.ActiveAttorneys {
					for _, enabledAttUid := range data.Form.DecisionAttorneysUids {
						if att.Uid == enabledAttUid {
							attorneyDecisions = append(attorneyDecisions, sirius.AttorneyDecisions{
								UID:                      att.Uid,
								CannotMakeJointDecisions: true,
							})
						}
					}
				}
			}

			err = client.ManageAttorneyDecisions(ctx, uid, attorneyDecisions)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
			} else if err != nil {
				return err
			} else {
				data.Success = true

				SetFlash(w, FlashNotification{Title: "Update saved"})
				return RedirectError(fmt.Sprintf("/lpa/%s", uid))
			}

		}

		return removeTmpl(w, data)
	}
}
