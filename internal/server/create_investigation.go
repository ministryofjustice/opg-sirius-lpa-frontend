package server

import (
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"net/http"
	"strconv"
)

type CreateInvestigationClient interface {
	CreateInvestigation(ctx sirius.Context, caseID int, caseType sirius.CaseType, investigation sirius.Investigation) error
	Case(ctx sirius.Context, id int) (sirius.Case, error)
}

type createInvestigationData struct {
	XSRFToken     string
	Success       bool
	Error         sirius.ValidationError
	Case          sirius.Case
	Investigation sirius.Investigation
}

func CreateInvestigation(client CreateInvestigationClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		caseID, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}

		caseType, err := sirius.ParseCaseType(r.FormValue("case"))
		if err != nil {
			return err
		}

		ctx := getContext(r)

		caseItem, err := client.Case(ctx, caseID)
		if err != nil {
			return err
		}

		data := createInvestigationData{
			XSRFToken: ctx.XSRFToken,
			Case:      caseItem,
		}

		if r.Method == http.MethodPost {
			investigation := sirius.Investigation{
				Title:        postFormString(r, "title"),
				Information:  postFormString(r, "information"),
				Type:         postFormString(r, "type"),
				DateReceived: postFormDateString(r, "dateReceived"),
			}

			err = client.CreateInvestigation(ctx, caseID, caseType, investigation)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
				data.Investigation = investigation
			} else if err != nil {
				return err
			} else {
				data.Success = true
			}
		}

		return tmpl(w, data)
	}
}
