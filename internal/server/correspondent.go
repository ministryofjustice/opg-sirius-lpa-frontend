package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type CorrespondentClient interface {
	CreateCorrespondent(ctx sirius.Context, correspondents []sirius.Person) error
	UpdatePerson(ctx sirius.Context, caseId int, correspondent sirius.Person) error
	Person(sirius.Context, int) (sirius.Person, error)
	Case(sirius.Context, int) (sirius.Case, error)
}

type correspondentData struct {
	XSRFToken            string
	CaseID               int
	Case                 sirius.Case
	Correspondent        sirius.Person
	Success              bool
	Error                sirius.ValidationError
	RelationshipToDonors []sirius.RefDataItem
	Title                string
	IsEditing            bool
}

func Correspondent(client CorrespondentClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		caseId, err := strToIntOrStatusError(r.FormValue("caseId"))
		if err != nil {
			return err
		}

		//get the attorneys
		caseitem, err := client.Case(ctx, caseId)
		if err != nil {
			return err
		}

		isEditing := r.FormValue("isEditing") == "true"

		data := correspondentData{
			XSRFToken: ctx.XSRFToken,
			Case:      caseitem,
			CaseID:    caseId,
			Title:     "Add an correspondent",
			IsEditing: isEditing,
		}

		if r.Method == http.MethodPost {
			if postFormString(r, "attorneyId") != "" {
				attorneyId, err := strToIntOrStatusError(postFormString(r, "attorneyId"))
				if err != nil {
					return err
				}

				var attorney *sirius.Person
				for i, a := range caseitem.Attorneys {
					if a.ID == attorneyId {
						attorney = &caseitem.Attorneys[i]
						break
					}
				}

				correspondent := sirius.Person{
					Salutation:                   attorney.Salutation,
					Firstname:                    attorney.Firstname,
					Middlenames:                  attorney.Middlenames,
					Surname:                      attorney.Surname,
					AlsoKnownAs:                  attorney.AlsoKnownAs,
					DateOfBirth:                  attorney.DateOfBirth,
					PhoneNumber:                  attorney.PhoneNumber,
					Email:                        attorney.Email,
					AddressLine1:                 attorney.AddressLine1,
					AddressLine2:                 attorney.AddressLine2,
					AddressLine3:                 attorney.AddressLine3,
					Town:                         attorney.Town,
					County:                       attorney.County,
					Country:                      attorney.Country,
					Postcode:                     attorney.Postcode,
					CompanyName:                  attorney.CompanyName,
					IsAirmailRequired:            attorney.IsAirmailRequired,
					IsAttorneyApplyingToRegister: attorney.IsAttorneyApplyingToRegister,
					PersonType:                   "Correspondent",
					CaseId:                       caseId,
				}

				err = client.CreateCorrespondent(ctx, []sirius.Person{correspondent})

				if ve, ok := err.(sirius.ValidationError); ok {
					w.WriteHeader(http.StatusBadRequest)
					data.Error = ve
				} else if err != nil {
					return err
				} else {
					return RedirectError(fmt.Sprintf("/edit-epa?caseId=%d", caseId))
				}
			}

			return RedirectError(fmt.Sprintf("/create-correnspondent?caseId=%d", caseId))
		}

		return tmpl(w, data)
	}
}
