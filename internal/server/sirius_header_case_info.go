package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type SiriusHeaderCaseInfoClient interface {
	Case(ctx sirius.Context, id int) (sirius.Case, error)
}

type siriusHeaderCaseInfoData struct {
	XSRFToken string
	CaseID    int
	Case      sirius.Case
}

func SiriusHeaderCaseInfo(client SiriusHeaderCaseInfoClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		caseId, err := strToIntOrStatusError(r.FormValue("id"))
		if err != nil {
			return err
		}

		ctx := getContext(r)
		data := siriusHeaderCaseInfoData{
			XSRFToken: ctx.XSRFToken,
			CaseID:    caseId,
		}

		caseItem, err := client.Case(ctx, caseId)
		if err != nil {
			return err
		}

		data.Case = caseItem

		return tmpl(w, data)
	}
}
