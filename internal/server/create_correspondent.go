package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type CreateCorrespondentClient interface {
	Epa(ctx sirius.Context, id int) (sirius.Epa, error)
	CreateCorrespondent(ctx sirius.Context, caseId int, correspondent sirius.Correspondent) error
	UpdateCorrespondent(ctx sirius.Context, correspondentId int, correspondent sirius.Correspondent) error
}

type createCorrespondentData struct {
	XSRFToken     string
	DonorId       int
	CaseId        int
	Correspondent sirius.Correspondent
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
			DonorId:   donorId,
			CaseId:    caseId,
			Title:     "Add a correspondent",
		}

		epa, err := client.Epa(ctx, caseId)
		if err != nil {
			return err
		}
		correspondent := epa.Correspondent
		isEditing := correspondent != nil

		if isEditing {
			data.Correspondent = *correspondent
			data.Title = "Update correspondent details"
			data.IsEditing = true
		}

		if r.Method == http.MethodPost {
			updatedCorrespondent := sirius.Correspondent{
				Person: sirius.Person{
					Salutation:        postFormString(r, "salutation"),
					Firstname:         postFormString(r, "firstname"),
					Middlenames:       postFormString(r, "middlenames"),
					Surname:           postFormString(r, "surname"),
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
			}
			data.Correspondent = updatedCorrespondent

			if isEditing {
				updatedCorrespondent.ID = correspondent.ID
				data.Correspondent = updatedCorrespondent
				err = client.UpdateCorrespondent(ctx, updatedCorrespondent.ID, updatedCorrespondent)
			} else {
				err = client.CreateCorrespondent(ctx, caseId, updatedCorrespondent)
			}

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
			} else if err != nil {
				return err
			} else {
				return RedirectError(fmt.Sprintf("/create-epa?id=%d&caseId=%d", donorId, caseId))
			}
		}

		return tmpl(w, data)
	}
}
