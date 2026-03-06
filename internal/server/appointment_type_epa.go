package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type AppointmentTypeEpaClient interface {
	UpdateEpa(ctx sirius.Context, caseId int, epa sirius.Case) error
	Case(ctx sirius.Context, id int) (sirius.Case, error)
}

type AppointmentTypeEpaData struct {
	XSRFToken        string
	Case             sirius.Case
	Success          bool
	Error            sirius.ValidationError
	Title            string
	ButtonName       string
	IsEditing        bool
	CaseAttorneyType string
}

func AppointmentTypeEpa(client AppointmentTypeEpaClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		isEditing := r.FormValue("isEditing") == "true"
		caseId, err := strToIntOrStatusError(r.FormValue("caseId"))
		if err != nil {
			return err
		}

		ctx := getContext(r)

		caseitem, err := client.Case(ctx, caseId)
		if err != nil {
			return err
		}

		var caseAttorneyType string
		if caseitem.CaseAttorneyJointly {
			caseAttorneyType = "jointly"
		} else if caseitem.CaseAttorneySingular {
			caseAttorneyType = "singular"
		} else if caseitem.CaseAttorneyJointlyAndSeverally {
			caseAttorneyType = "jointly-and-severally"
		}

		data := AppointmentTypeEpaData{
			XSRFToken:        ctx.XSRFToken,
			Case:             caseitem,
			Title:            "Step 2: appointment type",
			ButtonName:       "Save and continue",
			IsEditing:        isEditing,
			CaseAttorneyType: caseAttorneyType,
		}

		if isEditing {
			data.Title = "Appointment type"
			data.ButtonName = "save"
		}

		if r.Method == http.MethodPost {
			epa := sirius.Case{
				CaseAttorneySingular:            r.FormValue("caseAttorney") == "singular",
				CaseAttorneyJointlyAndSeverally: r.FormValue("caseAttorney") == "jointly-and-severally",
				CaseAttorneyJointly:             r.FormValue("caseAttorney") == "jointly",
			}

			fmt.Println(epa.CaseAttorneySingular, epa.CaseAttorneyJointlyAndSeverally, epa.CaseAttorneyJointly)

			err := client.UpdateEpa(ctx, caseId, epa)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
			} else if err != nil {
				return err
			}

			if isEditing {
				return RedirectError(fmt.Sprintf("/edit-epa?caseId=%d", caseId))
			}
			return RedirectError(fmt.Sprintf("/case-actors-epa?caseId=%d", caseId))
		}

		return tmpl(w, data)
	}
}
