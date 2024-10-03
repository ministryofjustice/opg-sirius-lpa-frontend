package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"net/http"
)

type ChangeCaseStatusClient interface {
	CaseSummary(ctx sirius.Context, uid string) (sirius.CaseSummary, error)
	EditDigitalLPAStatus(sirius.Context, string, sirius.CaseStatusData) error
}

type changeCaseStatusData struct {
	XSRFToken string
	Entity    string
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

		data := changeCaseStatusData{
			XSRFToken: ctx.XSRFToken,
			Entity:    fmt.Sprintf("%s %s", cs.DigitalLpa.SiriusData.Subtype, caseUID),
			OldStatus: cs.DigitalLpa.SiriusData.Status,
			NewStatus: postFormString(r, "status"),
		}

		if r.Method == http.MethodPost {
			caseDetails := sirius.CaseStatusData{
				Status: data.NewStatus,
			}

			err = client.EditDigitalLPAStatus(ctx, caseUID, caseDetails)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
			} else if err != nil {
				return err
			} else {
				data.Success = true
			}
		}

		return tmpl(w, data)
	}
}
