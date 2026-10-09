package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type actor struct {
	ID     int
	Name   string
	Label  string
	Fields string
}

type CreateCorrespondentClient interface {
	Epa(ctx sirius.Context, id int) (sirius.Epa, error)
	Lpa(ctx sirius.Context, id int) (sirius.Lpa, error)
	CreateCorrespondent(ctx sirius.Context, caseId int, correspondent sirius.Correspondent) error
	UpdateCorrespondent(ctx sirius.Context, correspondentId int, correspondent sirius.Correspondent) error
}

type createCorrespondentData struct {
	XSRFToken     string
	IsPartial     bool
	DonorId       int
	CaseId        int
	CaseType      string
	Correspondent sirius.Correspondent
	Actors        []actor
	Error         sirius.ValidationError
	IsEditing     bool
	Title         string
}

func CreateCorrespondent(client CreateCorrespondentClient, tmpl template.Template) Handler {
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

		data := createCorrespondentData{
			XSRFToken: ctx.XSRFToken,
			IsPartial: ctx.IsPartial,
			DonorId:   donorId,
			CaseId:    caseId,
			CaseType:  r.FormValue("caseType"),
			Title:     "Add a correspondent",
		}

		var correspondent *sirius.Correspondent
		if data.CaseType == "epa" {
			epa, err := client.Epa(ctx, caseId)
			if err != nil {
				return err
			}
			correspondent = epa.Correspondent
		} else {
			lpa, err := client.Lpa(ctx, caseId)
			if err != nil {
				return err
			}
			correspondent = lpa.Correspondent
			data.Actors = GetActorsForLpa(lpa)
		}

		if correspondent != nil {
			data.Correspondent = *correspondent
			data.Title = "Update correspondent details"
			data.IsEditing = true
		}

		if r.Method == http.MethodPost {
			updatedCorrespondent := sirius.Correspondent{
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
			data.Correspondent = updatedCorrespondent

			if data.IsEditing {
				updatedCorrespondent.ID = correspondent.ID
				data.Correspondent = updatedCorrespondent
				err = client.UpdateCorrespondent(ctx, data.Correspondent.ID, updatedCorrespondent)
			} else {
				err = client.CreateCorrespondent(ctx, caseId, updatedCorrespondent)
			}

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
			} else if err != nil {
				return err
			} else {
				if data.CaseType == "epa" {
					return RedirectError(fmt.Sprintf("/create-epa?id=%d&caseId=%d#accordion-create-epa-heading-3", donorId, caseId))
				} else {
					return RedirectError(fmt.Sprintf("/create-lpa?id=%d&caseId=%d#accordion-create-lpa-heading-4", donorId, caseId))
				}
			}

		}

		return tmpl(w, data)
	}
}

func GetActorsForLpa(lpa sirius.Lpa) []actor {
	var actors []actor

	if lpa.Donor != nil {
		fields, err := json.Marshal(lpa.Donor)
		if err == nil {
			actors = append(actors, actor{
				ID:     lpa.Donor.ID,
				Name:   lpa.Donor.Firstname + " " + lpa.Donor.Surname,
				Label:  "donor",
				Fields: string(fields),
			})
		}
	}

	for _, attorney := range lpa.Attorneys {
		fields, err := json.Marshal(attorney)
		if err != nil {
			continue
		}
		actors = append(actors, actor{
			ID:     attorney.ID,
			Name:   attorney.Firstname + " " + attorney.Surname,
			Label:  "attorney",
			Fields: string(fields),
		})
	}

	for _, attorney := range lpa.ReplacementAttorneys {
		fields, err := json.Marshal(attorney)
		if err != nil {
			continue
		}
		actors = append(actors, actor{
			ID:     attorney.ID,
			Name:   attorney.Firstname + " " + attorney.Surname,
			Label:  "replacement attorney",
			Fields: string(fields),
		})
	}

	for _, certificateProvider := range lpa.CertificateProviders {
		fields, err := json.Marshal(certificateProvider)
		if err != nil {
			continue
		}
		actors = append(actors, actor{
			ID:     certificateProvider.ID,
			Name:   certificateProvider.Firstname + " " + certificateProvider.Surname,
			Label:  "certificate provider",
			Fields: string(fields),
		})
	}

	for _, notifiedPerson := range lpa.NotifiedPersons {
		fields, err := json.Marshal(notifiedPerson)
		if err != nil {
			continue
		}
		actors = append(actors, actor{
			ID:     notifiedPerson.ID,
			Name:   notifiedPerson.Firstname + " " + notifiedPerson.Surname,
			Label:  "notified person",
			Fields: string(fields),
		})
	}

	for _, trustCorporation := range lpa.TrustCorporations {
		fields, err := json.Marshal(trustCorporation)
		if err != nil {
			continue
		}
		actors = append(actors, actor{
			ID:     trustCorporation.ID,
			Name:   trustCorporation.CompanyName,
			Label:  "trust corporation",
			Fields: string(fields),
		})
	}

	return actors
}
