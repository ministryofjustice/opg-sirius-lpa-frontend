package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type CaseActorsEpaClient interface {
	UpdateEpa(ctx sirius.Context, caseId int, epa sirius.Case) error
	Case(ctx sirius.Context, id int) (sirius.Case, error)
}

type CaseActorsEpaData struct {
	XSRFToken string
	CaseID    int
	Case      sirius.Case
	Success   bool
	Error     sirius.ValidationError
	Title     string
	IsEditing bool
}

func CaseActorsEpa(client CaseActorsEpaClient, tmpl template.Template) Handler {
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

		isEditing := r.FormValue("isEditing") == "true"

		data := CaseActorsEpaData{
			XSRFToken: ctx.XSRFToken,
			CaseID:    caseId,
			Title:     "Create EPA details",
			Case:      caseitem,
			IsEditing: isEditing,
		}

		return tmpl(w, data)
	}
}
