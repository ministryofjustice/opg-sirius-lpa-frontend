package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type SelectOrCreateCorrespondentClient interface {
	AssignActorAsCorrespondent(ctx sirius.Context, caseId int, personId int) error
	Epa(ctx sirius.Context, id int) (sirius.Epa, error)
}

type selectOrCreateCorrespondentData struct {
	XSRFToken     string
	DonorId       int
	CaseId        int
	Epa           sirius.Epa
	Correspondent sirius.Correspondent
	Error         sirius.ValidationError
}

func SelectOrCreateCorrespondent(client SelectOrCreateCorrespondentClient, tmpl template.Template, partialTmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		donorId, err := strToIntOrStatusError(r.FormValue("id"))
		if err != nil {
			return err
		}

		caseId, err := strToIntOrStatusError(r.FormValue("caseId"))
		if err != nil {
			return err
		}

		caseItem, err := client.Epa(ctx, caseId)
		if err != nil {
			return err
		}

		data := selectOrCreateCorrespondentData{
			XSRFToken: ctx.XSRFToken,
			DonorId:   donorId,
			CaseId:    caseId,
			Epa:       caseItem,
		}

		if r.Method == http.MethodPost {
			if postFormString(r, "attorneyId") != "new" {
				attorneyId, err := strToIntOrStatusError(postFormString(r, "attorneyId"))
				if err != nil {
					return err
				}

				var attorney *sirius.Attorney
				for _, caseAttorney := range caseItem.Attorneys {
					if caseAttorney.ID == attorneyId {
						attorney = &caseAttorney
						break
					}
				}

				err = client.AssignActorAsCorrespondent(ctx, caseId, attorney.ID)

				if ve, ok := err.(sirius.ValidationError); ok {
					w.WriteHeader(http.StatusBadRequest)
					data.Error = ve
					if r.Header.Get("HX-Request") == "true" {
						return partialTmpl(w, data)
					}
					return tmpl(w, data)
				} else if err != nil {
					return err
				} else {
					return RedirectError(fmt.Sprintf("/create-epa?id=%d&caseId=%d", donorId, caseId))
				}
			}

			return RedirectError(fmt.Sprintf("/create-correspondent?id=%d&caseId=%d", donorId, caseId))
		}

		if r.Header.Get("HX-Request") == "true" {
			return partialTmpl(w, data)
		}
		return tmpl(w, data)
	}
}
