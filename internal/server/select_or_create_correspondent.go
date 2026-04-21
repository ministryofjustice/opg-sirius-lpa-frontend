package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type SelectOrCreateCorrespondentClient interface {
	CreateCorrespondent(ctx sirius.Context, caseId int, correspondent sirius.Correspondent) error
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

func SelectOrCreateCorrespondent(client SelectOrCreateCorrespondentClient, tmpl template.Template) Handler {
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

				correspondent := sirius.Correspondent{
					Person: sirius.Person{
						Salutation:        attorney.Salutation,
						Firstname:         attorney.Firstname,
						Middlenames:       attorney.Middlenames,
						Surname:           attorney.Surname,
						AlsoKnownAs:       attorney.AlsoKnownAs,
						DateOfBirth:       attorney.DateOfBirth,
						PhoneNumber:       attorney.PhoneNumber,
						Email:             attorney.Email,
						AddressLine1:      attorney.AddressLine1,
						AddressLine2:      attorney.AddressLine2,
						AddressLine3:      attorney.AddressLine3,
						Town:              attorney.Town,
						County:            attorney.County,
						Country:           attorney.Country,
						Postcode:          attorney.Postcode,
						CompanyName:       attorney.CompanyName,
						IsAirmailRequired: attorney.IsAirmailRequired,
					},
				}

				err = client.CreateCorrespondent(ctx, caseId, correspondent)

				if ve, ok := err.(sirius.ValidationError); ok {
					w.WriteHeader(http.StatusBadRequest)
					data.Error = ve
				} else if err != nil {
					return err
				} else {
					return RedirectError(fmt.Sprintf("/create-epa?id=%d&caseId=%d", donorId, caseId))
				}
			}

			return RedirectError(fmt.Sprintf("/create-correspondent?id=%d&caseId=%d", donorId, caseId))
		}

		return tmpl(w, data)
	}
}
