package server

import (
	"fmt"
	"net/http"

	"github.com/go-playground/form/v4"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type UpdateDecisionsClient interface {
	CaseSummary(sirius.Context, string) (sirius.CaseSummary, error)
	UpdateDecisions(sirius.Context, string, sirius.UpdateDecisions) error
}

type updateDecisionsData struct {
	XSRFToken   string
	Success     bool
	Error       sirius.ValidationError
	Form        formDecisionsDetails
	CaseSummary sirius.CaseSummary
}

type formDecisionsDetails struct {
	WhenTheLpaCanBeUsed                         string `form:"whenTheLpaCanBeUsed"`
	LifeSustainingTreatmentOption               string `form:"lifeSustainingTreatmentOption"`
	HowAttorneysMakeDecisions                   string `form:"howAttorneysMakeDecisions"`
	HowAttorneysMakeDecisionsDetails            string `form:"howAttorneysMakeDecisionsDetails"`
	HowReplacementAttorneysStepIn               string `form:"howReplacementAttorneysStepIn"`
	HowReplacementAttorneysStepInDetails        string `form:"howReplacementAttorneysStepInDetails"`
	HowReplacementAttorneysMakeDecisions        string `form:"howReplacementAttorneysMakeDecisions"`
	HowReplacementAttorneysMakeDecisionsDetails string `form:"howReplacementAttorneysMakeDecisionsDetails"`
}

func UpdateDecisions(client UpdateDecisionsClient, tmpl template.Template) Handler {
	if decoder == nil {
		decoder = form.NewDecoder()
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		caseUID := r.PathValue("uid")
		ctx := getContext(r)

		cs, err := client.CaseSummary(ctx, caseUID)
		if err != nil {
			return err
		}

		lpaStoreData := cs.DigitalLpa.LpaStoreData

		data := updateDecisionsData{
			XSRFToken:   ctx.XSRFToken,
			CaseSummary: cs,
			Form: formDecisionsDetails{
				WhenTheLpaCanBeUsed:                         lpaStoreData.WhenTheLpaCanBeUsed,
				LifeSustainingTreatmentOption:               lpaStoreData.LifeSustainingTreatmentOption,
				HowAttorneysMakeDecisions:                   lpaStoreData.HowAttorneysMakeDecisions,
				HowAttorneysMakeDecisionsDetails:            lpaStoreData.HowAttorneysMakeDecisionsDetails,
				HowReplacementAttorneysStepIn:               lpaStoreData.HowReplacementAttorneysStepIn,
				HowReplacementAttorneysStepInDetails:        lpaStoreData.HowReplacementAttorneysStepInDetails,
				HowReplacementAttorneysMakeDecisions:        lpaStoreData.HowReplacementAttorneysMakeDecisions,
				HowReplacementAttorneysMakeDecisionsDetails: lpaStoreData.HowReplacementAttorneysMakeDecisionsDetails,
			},
		}

		if r.Method == http.MethodPost {
			if err := decoder.Decode(&data.Form, r.PostForm); err != nil {
				return err
			}

			err := client.UpdateDecisions(ctx, caseUID, sirius.UpdateDecisions{
				WhenTheLpaCanBeUsed:                         data.Form.WhenTheLpaCanBeUsed,
				LifeSustainingTreatmentOption:               data.Form.LifeSustainingTreatmentOption,
				HowAttorneysMakeDecisions:                   data.Form.HowAttorneysMakeDecisions,
				HowAttorneysMakeDecisionsDetails:            data.Form.HowAttorneysMakeDecisionsDetails,
				HowReplacementAttorneysStepIn:               data.Form.HowReplacementAttorneysStepIn,
				HowReplacementAttorneysStepInDetails:        data.Form.HowReplacementAttorneysStepInDetails,
				HowReplacementAttorneysMakeDecisions:        data.Form.HowReplacementAttorneysMakeDecisions,
				HowReplacementAttorneysMakeDecisionsDetails: data.Form.HowReplacementAttorneysMakeDecisionsDetails,
			})

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
			} else if err != nil {
				return err
			} else {
				data.Success = true

				SetFlash(w, FlashNotification{
					Title: "Update saved",
				})

				return RedirectError(fmt.Sprintf("/lpa/%s/lpa-details", caseUID))
			}
		}

		return tmpl(w, data)
	}
}
