package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
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
	CaseID        int
	CaseUIDs      string
	EntityType    string
	DonorId       int
}

func CreateInvestigation(client CreateInvestigationClient, tmpl template.Template, partialTmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		caseID, err := strToIntOrStatusError(r.FormValue("id"))
		if err != nil {
			return err
		}

		caseTypeString := r.FormValue("case")
		caseType, err := sirius.ParseCaseType(caseTypeString)
		if err != nil {
			return err
		}

		ctx := getContext(r)

		caseItem, err := client.Case(ctx, caseID)
		if err != nil {
			return err
		}

		data := createInvestigationData{
			XSRFToken:  ctx.XSRFToken,
			Case:       caseItem,
			CaseID:     caseID,
			CaseUIDs:   buildUIDQueryString(r.Form["uid[]"]),
			EntityType: caseTypeString,
			DonorId:    caseItem.Donor.ID,
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
		if r.Header.Get("HX-Request") == "true" && partialTmpl != nil {
			return partialTmpl(w, data)
		}

		return tmpl(w, data)
	}
}
