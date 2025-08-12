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

type AttorneyDetails struct {
	AttorneyName    string
	AttorneyDob     string
	AppointmentType string
}

type formManageAttorneyDecisions struct {
	DecisionAttorneysUids []string `form:"decisionAttorney"`
	SkipDecisionAttorney  string   `form:"skipDecisionAttorney"`
}

type manageAttorneyDecisionsData struct {
	CaseSummary                  sirius.CaseSummary
	DecisionAttorneys            []sirius.LpaStoreAttorney
	Form                         formManageAttorneyDecisions
	DecisionAttorneysDetails     []AttorneyDetails
	ActiveAttorneyCount          int
	ReplacementAttorneyCount     int
	Decisions                    string
	ReplacementAttorneyDecisions string
	Success                      bool
	Error                        sirius.ValidationError
	XSRFToken                    string
	FormName                     string
}

func AttorneyDecisions(client AttorneyDecisionsClient, decisionTmpl template.Template, confirmTmpl template.Template) Handler {

	return func(w http.ResponseWriter, r *http.Request) error {
		uid := r.PathValue("uid")
		ctx := getContext(r)

		caseSummary, err := client.CaseSummary(ctx, uid)

		if err != nil {
			return err
		}

		data := manageAttorneyDecisionsData{
			CaseSummary:                  caseSummary,
			XSRFToken:                    ctx.XSRFToken,
			Error:                        sirius.ValidationError{Field: sirius.FieldErrors{}},
			FormName:                     "decisions",
			Decisions:                    caseSummary.DigitalLpa.LpaStoreData.HowAttorneysMakeDecisions,
			ReplacementAttorneyDecisions: caseSummary.DigitalLpa.LpaStoreData.HowReplacementAttorneysMakeDecisions,
		}

		lpa := data.CaseSummary.DigitalLpa

		for _, attorney := range lpa.LpaStoreData.Attorneys {
			if attorney.Status == shared.ActiveAttorneyStatus.String() {
				data.DecisionAttorneys = append(data.DecisionAttorneys, attorney)
				data.ActiveAttorneyCount++
			}

			if attorney.AppointmentType == shared.ReplacementAppointmentType.String() &&
				attorney.Status == shared.InactiveAttorneyStatus.String() {
				data.ReplacementAttorneyCount++
			}
		}

		if data.ActiveAttorneyCount > 1 && data.ReplacementAttorneyCount > 1 && data.ReplacementAttorneyDecisions == "" {
			data.ReplacementAttorneyDecisions = data.Decisions
		}

		if r.Method == http.MethodPost {

			err = decoder.Decode(&data.Form, r.PostForm)
			if err != nil {
				return err
			}

			if (len(data.Form.DecisionAttorneysUids) == 0 && !postFormCheckboxChecked(r, "skipDecisionAttorney", "yes")) ||
				(len(data.Form.DecisionAttorneysUids) > 0 && postFormCheckboxChecked(r, "skipDecisionAttorney", "yes")) {
				data.Error.Field["decisionAttorney"] = map[string]string{
					"reason": "Select who cannot make joint decisions, or select 'Joint decisions can be made by all attorneys'",
				}
			}

			if !data.Error.Any() {
				if !postFormKeySet(r, "confirmDecisions") {

					if len(data.Form.DecisionAttorneysUids) > 0 {
						for _, att := range data.DecisionAttorneys {
							for _, enabledAttUid := range data.Form.DecisionAttorneysUids {
								if att.Uid == enabledAttUid {
									data.DecisionAttorneysDetails = append(data.DecisionAttorneysDetails, AttorneyDetails{
										AttorneyName:    att.FirstNames + " " + att.LastName,
										AttorneyDob:     att.DateOfBirth,
										AppointmentType: att.AppointmentType,
									})
									break
								}
							}
						}
					}

					return confirmTmpl(w, data)
				} else {
					var attorneyDecisions []sirius.AttorneyDecisions

					if data.Form.SkipDecisionAttorney == "yes" {
						for _, att := range data.DecisionAttorneys {
							attorneyDecisions = append(attorneyDecisions, sirius.AttorneyDecisions{
								UID:                      att.Uid,
								CannotMakeJointDecisions: false,
							})
						}
					} else {
						for _, att := range data.DecisionAttorneys {
							isChecked := false
							for _, selectedUid := range data.Form.DecisionAttorneysUids {
								if selectedUid == att.Uid {
									isChecked = true
									break
								}
							}
							attorneyDecisions = append(attorneyDecisions, sirius.AttorneyDecisions{
								UID:                      att.Uid,
								CannotMakeJointDecisions: isChecked,
							})
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
			}

		}

		return decisionTmpl(w, data)
	}
}
