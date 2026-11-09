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
	Lpa(ctx sirius.Context, id int) (sirius.Lpa, error)
}

type selectOrCreateCorrespondentData struct {
	XSRFToken     string
	FlowQuery     string
	IsPartial     bool
	DonorId       int
	CaseId        int
	CaseType      string
	Epa           sirius.Epa
	Lpa           sirius.Lpa
	Correspondent sirius.Correspondent
	Error         sirius.ValidationError
	Actors        []actor
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

		data := selectOrCreateCorrespondentData{
			XSRFToken: ctx.XSRFToken,
			FlowQuery: getFlowQuery(r),
			DonorId:   donorId,
			CaseId:    caseId,
			CaseType:  r.FormValue("caseType"),
			IsPartial: ctx.IsPartial,
		}

		if data.CaseType == "epa" {
			caseItem, err := client.Epa(ctx, caseId)
			if err != nil {
				return err
			}
			data.Epa = caseItem
		} else {
			caseItem, err := client.Lpa(ctx, caseId)
			if err != nil {
				return err
			}
			data.Lpa = caseItem
			data.Actors = GetActorsForLpa(caseItem)
		}

		if r.Method == http.MethodPost {
			if data.CaseType == "lpa" {
				data.Correspondent = sirius.Correspondent{
					Person: sirius.Person{
						AddressLine1:          postFormString(r, "addressLine1"),
						AddressLine2:          postFormString(r, "addressLine2"),
						AddressLine3:          postFormString(r, "addressLine3"),
						CompanyName:           postFormString(r, "companyName"),
						Country:               postFormString(r, "country"),
						County:                postFormString(r, "county"),
						Email:                 postFormString(r, "email"),
						Firstname:             postFormString(r, "firstname"),
						IsAirmailRequired:     postFormString(r, "isAirmailRequired") == "true",
						Middlenames:           postFormString(r, "middlenames"),
						PhoneNumber:           postFormString(r, "phoneNumber"),
						Postcode:              postFormString(r, "postcode"),
						Salutation:            postFormString(r, "salutation"),
						Surname:               postFormString(r, "surname"),
						Town:                  postFormString(r, "town"),
						CorrespondenceByPost:  postFormCheckboxChecked(r, "correspondenceBy", "post"),
						CorrespondenceByEmail: postFormCheckboxChecked(r, "correspondenceBy", "email"),
						CorrespondenceByPhone: postFormCheckboxChecked(r, "correspondenceBy", "phone"),
						CorrespondenceByWelsh: postFormCheckboxChecked(r, "correspondenceBy", "welsh"),
					},
					CompanyNumber: postFormString(r, "companyNumber"),
				}
			} else if postFormString(r, "actorId") != "new" {
				actorId, err := strToIntOrStatusError(postFormString(r, "actorId"))
				if err != nil {
					return err
				}

				var person sirius.Person
				for _, caseAttorney := range data.Epa.Attorneys {
					if caseAttorney.ID == actorId {
						person = caseAttorney.Person
						break
					}
				}

				data.Correspondent = sirius.Correspondent{
					Person: sirius.Person{
						Salutation:        person.Salutation,
						Firstname:         person.Firstname,
						Middlenames:       person.Middlenames,
						Surname:           person.Surname,
						DateOfBirth:       person.DateOfBirth,
						PhoneNumber:       person.PhoneNumber,
						Email:             person.Email,
						AddressLine1:      person.AddressLine1,
						AddressLine2:      person.AddressLine2,
						AddressLine3:      person.AddressLine3,
						Town:              person.Town,
						County:            person.County,
						Country:           person.Country,
						Postcode:          person.Postcode,
						CompanyName:       person.CompanyName,
						IsAirmailRequired: person.IsAirmailRequired,
					},
				}
			} else {
				return RedirectError(fmt.Sprintf("/create-correspondent?id=%d&caseId=%d&caseType=%s%s", donorId, caseId, data.CaseType, data.FlowQuery))
			}

			err = client.CreateCorrespondent(ctx, caseId, data.Correspondent)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
				return tmpl(w, data)
			} else if err != nil {
				return err
			} else {
				if data.CaseType == "epa" {
					return RedirectError(fmt.Sprintf("/create-epa?id=%d&caseId=%d#accordion-create-epa-heading-3", donorId, caseId))
				}
				return RedirectError(fmt.Sprintf("/create-lpa?id=%d&caseId=%d%s#accordion-create-lpa-heading-4", donorId, caseId, data.FlowQuery))
			}
		}

		return tmpl(w, data)
	}
}
