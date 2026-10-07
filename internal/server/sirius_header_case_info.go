package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type SiriusHeaderCaseInfoClient interface {
	Case(ctx sirius.Context, id int) (sirius.Case, error)
	Lpa(ctx sirius.Context, id int) (sirius.Lpa, error)
}

type siriusHeaderCaseInfoData struct {
	XSRFToken                    string
	CaseID                       int
	Case                         sirius.Case
	HasInstructionsOrPreferences bool
}

func SiriusHeaderCaseInfo(client SiriusHeaderCaseInfoClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		caseId, err := strToIntOrStatusError(r.FormValue("id"))
		if err != nil {
			return err
		}

		caseType := r.FormValue("caseType")

		ctx := getContext(r)
		data := siriusHeaderCaseInfoData{
			XSRFToken: ctx.XSRFToken,
			CaseID:    caseId,
		}

		if caseType == "LPA" {
			caseItem, err := client.Lpa(ctx, caseId)
			if err != nil {
				return err
			}
			data.Case = caseItem.Case
			data.HasInstructionsOrPreferences = (caseItem.ApplicationHasGuidance != nil && *caseItem.ApplicationHasGuidance) || (caseItem.ApplicationHasRestrictions != nil && *caseItem.ApplicationHasRestrictions)
		} else {
			caseItem, err := client.Case(ctx, caseId)
			if err != nil {
				return err
			}
			data.Case = caseItem
		}

		return tmpl(w, data)
	}
}
