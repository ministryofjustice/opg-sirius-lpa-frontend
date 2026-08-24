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
	IsPartial     bool
	DonorId       int
	CaseId        int
	CaseType      string
	Epa           sirius.Epa
	Lpa           sirius.Lpa
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

		data := selectOrCreateCorrespondentData{
			XSRFToken: ctx.XSRFToken,
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
		}

		if r.Method == http.MethodPost {
			if postFormString(r, "actorId") != "new" {
				actorId, err := strToIntOrStatusError(postFormString(r, "actorId"))
				if err != nil {
					return err
				}

				var person sirius.Person
				if data.CaseType == "epa" {
					for _, caseAttorney := range data.Epa.Attorneys {
						if caseAttorney.ID == actorId {
							person = caseAttorney.Person
							break
						}
					}
				} else {
					person = GetSelectedActorForLpa(data.Lpa, actorId)
				}

				correspondent := sirius.Correspondent{
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

				err = client.CreateCorrespondent(ctx, caseId, correspondent)

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
					return RedirectError(fmt.Sprintf("/create-lpa?id=%d&caseId=%d#accordion-create-lpa-heading-4", donorId, caseId))
				}
			}

			return RedirectError(fmt.Sprintf("/create-correspondent?id=%d&caseId=%d&caseType=%s", donorId, caseId, data.CaseType))
		}

		return tmpl(w, data)
	}
}

func GetSelectedActorForLpa(lpa sirius.Lpa, actorId int) sirius.Person {
	if actorId == lpa.Donor.ID {
		return *lpa.Donor
	}

	for _, attorney := range lpa.Attorneys {
		if attorney.ID == actorId {
			return attorney.Person
		}
	}

	for _, attorney := range lpa.ReplacementAttorneys {
		if attorney.ID == actorId {
			return attorney.Person
		}
	}

	for _, certifiedProvider := range lpa.CertificateProviders {
		if certifiedProvider.ID == actorId {
			return certifiedProvider
		}
	}

	for _, notifiedPerson := range lpa.NotifiedPersons {
		if notifiedPerson.ID == actorId {
			return notifiedPerson.Person
		}
	}

	for _, trustCorporation := range lpa.TrustCorporations {
		if trustCorporation.ID == actorId {
			return trustCorporation.Person
		}
	}

	return sirius.Person{}
}
