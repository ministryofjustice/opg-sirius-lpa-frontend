package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type EditEpaAttorneyClient interface {
	UpdatePerson(ctx sirius.Context, caseId int, attorney sirius.Person) error
	Person(sirius.Context, int) (sirius.Person, error)
}

type editEpaAttorneyData struct {
	XSRFToken            string
	CaseID               int
	Attorney             sirius.Person
	Success              bool
	Error                sirius.ValidationError
	RelationshipToDonors []sirius.RefDataItem
}

func EditEpaAttorney(client EditEpaAttorneyClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		caseId, err := strToIntOrStatusError(r.FormValue("caseId"))
		if err != nil {
			return err
		}

		attorneyId, err := strconv.Atoi(r.FormValue("attorneyId"))
		if err != nil {
			return err
		}
		fmt.Println(caseId, attorneyId)

		attorney, err := client.Person(ctx, attorneyId)
		if err != nil {
			return err
		}

		data := editEpaAttorneyData{
			XSRFToken: ctx.XSRFToken,
			CaseID:    caseId,
			Attorney:  attorney,
			RelationshipToDonors: []sirius.RefDataItem{
				{Handle: "civil partner", Label: "civil partner"},
				{Handle: "child", Label: "child"},
				{Handle: "solicitor", Label: "solicitor"},
				{Handle: "other", Label: "other"},
				{Handle: "other professional", Label: "other professional"},
			},
		}

		if r.Method == http.MethodPost {
			fmt.Println(caseId, attorneyId)
			updateAttorney := sirius.Person{
				Salutation:                   postFormString(r, "salutation"),
				Firstname:                    postFormString(r, "firstname"),
				Middlenames:                  postFormString(r, "middlenames"),
				Surname:                      postFormString(r, "surname"),
				Othernames:                   postFormString(r, "otherNames"),
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
			data.Attorney = attorney

			err := client.UpdatePerson(ctx, attorneyId, updateAttorney)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
			} else if err != nil {
				return err
			}
			return RedirectError(fmt.Sprintf("/case-actors-epa?caseId=%d", caseId))
		}

		return tmpl(w, data)
	}
}
