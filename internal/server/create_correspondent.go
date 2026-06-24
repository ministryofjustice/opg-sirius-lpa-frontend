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
					AddressLine1:      postFormString(r, "addressLine1"),
					AddressLine2:      postFormString(r, "addressLine2"),
					AddressLine3:      postFormString(r, "addressLine3"),
					CompanyName:       postFormString(r, "companyName"),
					CompanyNumber:     postFormString(r, "companyNumber"),
					Country:           postFormString(r, "country"),
					County:            postFormString(r, "county"),
					Email:             postFormString(r, "email"),
					Firstname:         postFormString(r, "firstname"),
					IsAirmailRequired: postFormString(r, "isAirmailRequired") == "true",
					Middlenames:       postFormString(r, "middlenames"),
					PhoneNumber:       postFormString(r, "phoneNumber"),
					Postcode:          postFormString(r, "postcode"),
					Salutation:        postFormString(r, "salutation"),
					Surname:           postFormString(r, "surname"),
					Town:              postFormString(r, "town"),
				},
			}
			data.Correspondent = updatedCorrespondent

			if isEditing {
				err = client.UpdateCorrespondent(ctx, correspondent.ID, updatedCorrespondent)
			} else {
				err = client.CreateCorrespondent(ctx, caseId, updatedCorrespondent)
			}

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
			} else if err != nil {
				return err
			} else {
				return RedirectError(fmt.Sprintf("/create-epa?id=%d&caseId=%d#accordion-create-epa-heading-3", donorId, caseId))
			}

		}

		return tmpl(w, data)
	}
}
