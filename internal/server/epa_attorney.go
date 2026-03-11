package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type EpaAttorneyClient interface {
	CreateAttorney(ctx sirius.Context, caseId int, attorney sirius.Person) error
	UpdatePerson(ctx sirius.Context, caseId int, attorney sirius.Person) error
	Person(sirius.Context, int) (sirius.Person, error)
}

type epaAttorneyData struct {
	XSRFToken            string
	CaseID               int
	Attorney             sirius.Person
	Success              bool
	Error                sirius.ValidationError
	RelationshipToDonors []sirius.RefDataItem
	Title                string
	IsEditing            bool
}

func EpaAttorney(client EpaAttorneyClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		caseId, err := strToIntOrStatusError(r.FormValue("caseId"))
		if err != nil {
			return err
		}

		data := epaAttorneyData{
			XSRFToken: ctx.XSRFToken,
			CaseID:    caseId,
			RelationshipToDonors: []sirius.RefDataItem{
				{Handle: "civil partner", Label: "civil partner"},
				{Handle: "child", Label: "child"},
				{Handle: "solicitor", Label: "solicitor"},
				{Handle: "other", Label: "other"},
				{Handle: "other professional", Label: "other professional"},
			},
			Title: "Add an attorney",
		}

		hasAttorneyId := r.FormValue("attorneyId") != ""
		var attorneyId int

		if hasAttorneyId {
			var err error
			attorneyId, err = strToIntOrStatusError(r.FormValue("attorneyId"))
			if err != nil {
				return err
			}

			attorney, err := client.Person(ctx, attorneyId)
			if err != nil {
				return err
			}

			data.Attorney = attorney
			data.Title = "Edit attorney"
			data.IsEditing = true
		}

		if r.Method == http.MethodPost {
			attorney := sirius.Person{
				Salutation:                   postFormString(r, "salutation"),
				Firstname:                    postFormString(r, "firstname"),
				Middlenames:                  postFormString(r, "middlenames"),
				Surname:                      postFormString(r, "surname"),
				AlsoKnownAs:                  postFormString(r, "otherNames"),
				DateOfBirth:                  postFormDateString(r, "dob"),
				PhoneNumber:                  postFormString(r, "phoneNumber"),
				Email:                        postFormString(r, "email"),
				AddressLine1:                 postFormString(r, "addressLine1"),
				AddressLine2:                 postFormString(r, "addressLine2"),
				AddressLine3:                 postFormString(r, "addressLine3"),
				Town:                         postFormString(r, "town"),
				County:                       postFormString(r, "county"),
				Country:                      postFormString(r, "country"),
				Postcode:                     postFormString(r, "postcode"),
				CompanyName:                  postFormString(r, "companyName"),
				IsAirmailRequired:            postFormString(r, "isAirmailRequired") == "true",
				RelationshipToDonor:          postFormString(r, "relationshipToDonor"),
				IsAttorneyApplyingToRegister: postFormString(r, "isAttorneyApplyingToRegister") == "true",
			}

			if hasAttorneyId {
				err = client.UpdatePerson(ctx, attorneyId, attorney)
			} else {
				err = client.CreateAttorney(ctx, caseId, attorney)
			}

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
			} else if err != nil {
				return err
			}
			if hasAttorneyId {
				return RedirectError(fmt.Sprintf("/edit-epa?caseId=%d", caseId))
			}
			return RedirectError(fmt.Sprintf("/case-actors-epa?caseId=%d", caseId))
		}

		return tmpl(w, data)
	}
}
