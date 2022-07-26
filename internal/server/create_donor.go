package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type CreateDonorClient interface {
	CreatePerson(ctx sirius.Context, personType sirius.PersonType, donor sirius.Person) (sirius.Person, error)
}

type createDonorData struct {
	XSRFToken string
	Success   bool
	Error     sirius.ValidationError
	Donor     sirius.Person
}

func CreateDonor(client CreateDonorClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		data := createDonorData{
			XSRFToken: ctx.XSRFToken,
		}

		if r.Method == http.MethodPost {
			donor := sirius.Person{
				Salutation:            postFormString(r, "salutation"),
				Firstname:             postFormString(r, "firstname"),
				Middlenames:           postFormString(r, "middlenames"),
				Surname:               postFormString(r, "surname"),
				DateOfBirth:           postFormDateString(r, "dob"),
				PreviouslyKnownAs:     postFormString(r, "previouslyKnownAs"),
				AlsoKnownAs:           postFormString(r, "alsoKnownAs"),
				AddressLine1:          postFormString(r, "addressLine1"),
				AddressLine2:          postFormString(r, "addressLine2"),
				AddressLine3:          postFormString(r, "addressLine3"),
				Town:                  postFormString(r, "town"),
				County:                postFormString(r, "county"),
				Postcode:              postFormString(r, "postcode"),
				Country:               postFormString(r, "country"),
				IsAirmailRequired:     postFormString(r, "isAirmailRequired") == "Yes",
				PhoneNumber:           postFormString(r, "phoneNumber"),
				Email:                 postFormString(r, "email"),
				SageId:                postFormString(r, "sageId"),
				CorrespondenceByPost:  postFormCheckboxChecked(r, "correspondenceBy", "post"),
				CorrespondenceByEmail: postFormCheckboxChecked(r, "correspondenceBy", "email"),
				CorrespondenceByPhone: postFormCheckboxChecked(r, "correspondenceBy", "phone"),
				CorrespondenceByWelsh: postFormCheckboxChecked(r, "correspondenceBy", "welsh"),
				ResearchOptOut:        postFormString(r, "researchOptOut") == "Yes",
			}

			createdDonor, err := client.CreatePerson(ctx, sirius.PersonTypeDonor, donor)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
				data.Donor = donor
			} else if err != nil {
				return err
			} else {
				data.Success = true
				data.Donor = createdDonor
			}
		}

		return tmpl(w, data)
	}
}
