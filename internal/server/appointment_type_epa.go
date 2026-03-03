package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type AppointmentTypeEpaClient interface {
	UpdateEpaPut(ctx sirius.Context, caseId int, epa sirius.Case) error
	Case(ctx sirius.Context, id int) (sirius.Case, error)
}

type AppointmentTypeEpaData struct {
	XSRFToken string
	Case      sirius.Case
	Success   bool
	Error     sirius.ValidationError
	Title     string
}

func AppointmentTypeEpa(client AppointmentTypeEpaClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		caseId, err := strToIntOrStatusError(r.FormValue("caseId"))
		if err != nil {
			return err
		}

		ctx := getContext(r)

		data := AppointmentTypeEpaData{
			XSRFToken: ctx.XSRFToken,
			Title:     "Create EPA details",
		}

		epa, err := client.Case(ctx, caseId)
		if err != nil {
			return err
		}

		if r.Method == http.MethodPost {
			epa = sirius.Case{
				CaseAttorneySingular:            r.FormValue("caseAttorney") == "singular",
				CaseAttorneyJointlyAndSeverally: r.FormValue("caseAttorney") == "jointly-and-severally",
				CaseAttorneyJointly:             r.FormValue("caseAttorney") == "jointly",
			}

			err := client.UpdateEpaPut(ctx, caseId, epa)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
			} else if err != nil {
				return err
			} else {
				return RedirectError(fmt.Sprintf("/case-actors-epa?caseId=%d", caseId))
			}
		}

		return tmpl(w, data)
	}
}
