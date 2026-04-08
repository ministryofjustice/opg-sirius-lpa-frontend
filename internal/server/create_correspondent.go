package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type CreateCorrespondentClient interface {
	CreateCorrespondent(ctx sirius.Context, correspondents []sirius.Person) error
	UpdatePerson(ctx sirius.Context, caseId int, correspondent sirius.Person) error
	Person(sirius.Context, int) (sirius.Person, error)
	Case(sirius.Context, int) (sirius.Case, error)
}

type createCorrespondentData struct {
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

func CreateCorrespondent(client CreateCorrespondentClient, tmpl template.Template) Handler {
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

		data := createCorrespondentData{
			XSRFToken: ctx.XSRFToken,
			Case:      caseitem,
			CaseID:    caseId,
			Title:     "Add a new correspondent",
			IsEditing: isEditing,
		}

		if r.Method == http.MethodPost {
			correspondent := sirius.Person{
				Salutation:        postFormString(r, "salutation"),
				Firstname:         postFormString(r, "firstname"),
				Middlenames:       postFormString(r, "middlenames"),
				Surname:           postFormString(r, "surname"),
				AlsoKnownAs:       postFormString(r, "otherNames"),
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
				PersonType:        "Correspondent",
				CaseId:            caseId,
			}

			err = client.CreateCorrespondent(ctx, []sirius.Person{correspondent})
			fmt.Println(err)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
			} else if err != nil {
				return err
			} else {
				return RedirectError(fmt.Sprintf("/edit-epa?caseId=%d", caseId))
			}
		}

		return tmpl(w, data)
	}
}
