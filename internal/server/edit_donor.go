package server

import (
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type EditDonorClient interface {
	EditDonor(ctx sirius.Context, personID int, donor sirius.Person) error
	Person(ctx sirius.Context, personID int) (sirius.Person, error)
}

func EditDonor(client EditDonorClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}

		ctx := getContext(r)

		donor, err := client.Person(ctx, id)
		if err != nil {
			return err
		}

		data := donorData{
			XSRFToken: ctx.XSRFToken,
			Donor:     donor,
		}

		if r.Method == http.MethodPost {
			data.Donor.Salutation = postFormString(r, "salutation")
			data.Donor.Firstname = postFormString(r, "firstname")
			data.Donor.Middlenames = postFormString(r, "middlenames")
			data.Donor.Surname = postFormString(r, "surname")
			data.Donor.DateOfBirth = postFormDateString(r, "dob")
			data.Donor.PreviouslyKnownAs = postFormString(r, "previousNames")
			data.Donor.AlsoKnownAs = postFormString(r, "otherNames")
			data.Donor.AddressLine1 = postFormString(r, "addressLine1")
			data.Donor.AddressLine2 = postFormString(r, "addressLine2")
			data.Donor.AddressLine3 = postFormString(r, "addressLine3")
			data.Donor.Town = postFormString(r, "town")
			data.Donor.County = postFormString(r, "county")
			data.Donor.Postcode = postFormString(r, "postcode")
			data.Donor.Country = postFormString(r, "country")
			data.Donor.IsAirmailRequired = postFormString(r, "isAirmailRequired") == "Yes"
			data.Donor.PhoneNumber = postFormString(r, "phoneNumber")
			data.Donor.Email = postFormString(r, "email")
			data.Donor.SageId = postFormString(r, "sageId")
			data.Donor.CorrespondenceByPost = postFormCheckboxChecked(r, "correspondenceBy", "post")
			data.Donor.CorrespondenceByEmail = postFormCheckboxChecked(r, "correspondenceBy", "email")
			data.Donor.CorrespondenceByPhone = postFormCheckboxChecked(r, "correspondenceBy", "phone")
			data.Donor.CorrespondenceByWelsh = postFormCheckboxChecked(r, "correspondenceBy", "welsh")
			data.Donor.ResearchOptOut = postFormString(r, "researchOptOut") == "Yes"

			err := client.EditDonor(ctx, id, data.Donor)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
			} else if err != nil {
				return err
			} else {
				data.Success = true
			}
		}

		return tmpl(w, data)
	}
}
