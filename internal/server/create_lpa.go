package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type CreateLpaClient interface {
	Person(sirius.Context, int) (sirius.Person, error)
	Case(sirius.Context, int) (sirius.Case, error)
}

type createLpaData struct {
	XSRFToken string
	Success   bool
	Error     sirius.ValidationError
	DonorId   int
	DonorName string
	Title     string
	CaseId    int
	CaseItem  sirius.Case
}

func CreateLpa(client CreateLpaClient, tmpl template.Template, partialTmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		donorID, err := strToIntOrStatusError(r.FormValue("id"))
		if err != nil {
			return err
		}

		donor, err := client.Person(ctx, donorID)
		if err != nil {
			return err
		}

		data := createLpaData{
			XSRFToken: ctx.XSRFToken,
			DonorId:   donorID,
			DonorName: donor.Firstname + " " + donor.Surname,
			Title:     "Create an LPA",
		}

		caseIdStr := r.FormValue("caseId")
		isEditing := caseIdStr != ""
		if isEditing {
			data.CaseId, err = strToIntOrStatusError(caseIdStr)
			if err != nil {
				return err
			}

			data.CaseItem, err = client.Case(ctx, data.CaseId)
			if err != nil {
				return err
			}

			data.Title = "Edit LPA"
		}

		if r.Header.Get("HX-Request") == "true" {
			return partialTmpl(w, data)
		}

		return tmpl(w, data)
	}
}
