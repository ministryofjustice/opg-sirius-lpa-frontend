package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type CreateAttorneyClient interface {
	Epa(ctx sirius.Context, id int) (sirius.Epa, error)
	CreateAttorney(ctx sirius.Context, caseId int, attorney sirius.Attorney) error
	RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error)
	UpdateAttorney(ctx sirius.Context, attorneyId int, attorney sirius.Attorney) error
}

type createAttorneyData struct {
	XSRFToken            string
	Attorney             sirius.Attorney
	Error                sirius.ValidationError
	RelationshipToDonors []sirius.RefDataItem
	DonorId              int
	CaseId               int
	IsEditing            bool
	Title                string
}

func CreateAttorney(client CreateAttorneyClient, tmpl template.Template) Handler {
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

		data := createAttorneyData{
			XSRFToken: ctx.XSRFToken,
			DonorId:   donorId,
			CaseId:    caseId,
			Title:     "Add an attorney",
		}

		data.RelationshipToDonors, err = client.RefDataByCategory(ctx, sirius.RelationshipToDonorCategory)
		if err != nil {
			return err
		}

		// Default the active status to true for new attorneys
		data.Attorney.SystemStatus = shared.BoolPtr(true)

		var attorneyId int
		attorneyIdStr := r.FormValue("attorneyId")
		isEditing := attorneyIdStr != ""
		if isEditing {
			attorneyId, err = strToIntOrStatusError(attorneyIdStr)
			if err != nil {
				return err
			}
			epa, err := client.Epa(ctx, caseId)
			if err != nil {
				return err
			}

			for _, attorney := range epa.Attorneys {
				if attorney.ID == attorneyId {
					data.Attorney = attorney
					break
				}
			}

			data.Title = "Update attorney details"
			data.IsEditing = true
		}

		if r.Method == http.MethodPost {
			attorney := sirius.Attorney{
				Person: sirius.Person{
					Salutation:        postFormString(r, "salutation"),
					Firstname:         postFormString(r, "firstname"),
					Middlenames:       postFormString(r, "middlenames"),
					Surname:           postFormString(r, "surname"),
					DateOfBirth:       postFormDateString(r, "dob"),
					PhoneNumber:       postFormString(r, "phoneNumber"),
					Email:             postFormString(r, "email"),
					AddressLine1:      postFormString(r, "addressLine1"),
					AddressLine2:      postFormString(r, "addressLine2"),
					AddressLine3:      postFormString(r, "addressLine3"),
					Town:              postFormString(r, "town"),
					County:            postFormString(r, "county"),
					Country:           postFormString(r, "country"),
					Postcode:          postFormString(r, "postcode"),
					CompanyName:       postFormString(r, "companyName"),
					IsAirmailRequired: postFormString(r, "isAirmailRequired") == "true",
				},
				RelationshipToDonor: postFormString(r, "relationshipToDonor"),
				SystemStatus:        shared.BoolPtr(postFormString(r, "isAttorneyActive") == "true"),
			}
			data.Attorney = attorney

			if isEditing {
				err = client.UpdateAttorney(ctx, attorneyId, attorney)
			} else {
				err = client.CreateAttorney(ctx, caseId, attorney)
			}

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
				return tmpl(w, data)
			} else if err != nil {
				return err
			}

			if r.FormValue("add-another") != "" {
				return RedirectError(fmt.Sprintf("/create-attorney?id=%d&caseId=%d", donorId, caseId))
			}

			return RedirectError(fmt.Sprintf("/create-epa?id=%d&caseId=%d", donorId, caseId))

		}

		return tmpl(w, data)
	}
}
