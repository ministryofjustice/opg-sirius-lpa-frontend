package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/templatefn"
	"net/http"
)

type ChangeCaseStatusClient interface {
	CaseSummary(sirius.Context, string) (sirius.CaseSummary, error)
	EditDigitalLPAStatus(sirius.Context, string, sirius.CaseStatusData) error
}

type changeCaseStatusData struct {
	XSRFToken string
	Entity    string
	CaseUID   string
	Success   bool
	Error     sirius.ValidationError

	OldStatus string
	NewStatus string
}

func ChangeCaseStatus(client ChangeCaseStatusClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		caseUID := r.FormValue("uid")

		ctx := getContext(r)
		cs, err := client.CaseSummary(ctx, caseUID)
		if err != nil {
			return err
		}

		status := "draft"

		if cs.DigitalLpa.LpaStoreData.Status != "" {
			status = cs.DigitalLpa.LpaStoreData.Status
		}

		data := changeCaseStatusData{
			XSRFToken: ctx.XSRFToken,
			Entity:    fmt.Sprintf("%s %s", cs.DigitalLpa.SiriusData.Subtype, caseUID),
			CaseUID:   caseUID,
			OldStatus: status,
			NewStatus: postFormString(r, "status"),
		}

		if r.Method == http.MethodPost {
			caseStatusData := sirius.CaseStatusData{
				Status: data.NewStatus,
			}

			err = client.EditDigitalLPAStatus(ctx, caseUID, caseStatusData)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
			} else if err != nil {
				return err
			} else {
				data.Success = true
				data.OldStatus = data.NewStatus

				SetFlash(w, FlashNotification{
					Title: fmt.Sprintf("Status changed to %s", templatefn.StatusLabelFormat(data.NewStatus)),
				})
				return RedirectError(fmt.Sprintf("/lpa/%s", data.CaseUID))
			}
		}

		return tmpl(w, data)
	}
}
