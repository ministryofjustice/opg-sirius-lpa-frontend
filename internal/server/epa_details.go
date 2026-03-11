package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type EpaDetailsClient interface {
	UpdateEpa(ctx sirius.Context, caseId int, epa sirius.Case) error
	Case(sirius.Context, int) (sirius.Case, error)
}
type EpaDetailsData struct {
	XSRFToken            string
	CaseID               int
	Case                 sirius.Case
	Success              bool
	Error                sirius.ValidationError
	ShowAllSections      bool
	RelationshipToDonors []sirius.RefDataItem
}

func EpaDetails(client EpaDetailsClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		caseId, err := strToIntOrStatusError(r.FormValue("caseId"))
		if err != nil {
			return err
		}

		ctx := getContext(r)
		caseitem, err := client.Case(ctx, caseId)
		if err != nil {
			return err
		}

		data := EpaDetailsData{
			XSRFToken: ctx.XSRFToken,
			CaseID:    caseId,
			Case:      caseitem,
			RelationshipToDonors: []sirius.RefDataItem{
				{Handle: "civil partner", Label: "civil partner"},
				{Handle: "child", Label: "child"},
				{Handle: "solicitor", Label: "solicitor"},
				{Handle: "other", Label: "other"},
				{Handle: "other professional", Label: "other professional"},
			},
		}

		return tmpl(w, data)
	}
}
