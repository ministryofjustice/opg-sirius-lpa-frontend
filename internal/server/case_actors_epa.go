package server

import (
	"fmt"
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
	Epa       sirius.Case
	Success   bool
	Error     sirius.ValidationError
	Title     string
}

func CaseActorsEpa(client CaseActorsEpaClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		caseId, err := strToIntOrStatusError(r.FormValue("caseId"))
		if err != nil {
			return err
		}

		ctx := getContext(r)

		epa, err := client.Case(ctx, caseId)
		if err != nil {
			return err
		}

		data := CaseActorsEpaData{
			XSRFToken: ctx.XSRFToken,
			CaseID:    caseId,
			Title:     "Create EPA details",
			Epa:       epa,
		}

		if r.Method == http.MethodPost {
			epa = sirius.Case{}

			err := client.UpdateEpa(ctx, caseId, epa)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
			} else if err != nil {
				return err
			} else {
				return RedirectError(fmt.Sprintf("/payment-epa?caseId=%d", caseId))
			}
		}

		return tmpl(w, data)
	}
}
