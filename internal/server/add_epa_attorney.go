package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type AddEpaAttorneyClient interface {
	CreateAttorney(ctx sirius.Context, caseId int, attorney sirius.Attorney) error
}

type addEpaAttorneyData struct {
	XSRFToken string
	CaseID    int
	Attorney  sirius.Attorney
	Success   bool
	Error     sirius.ValidationError
}

func AddEpaAttorney(client AddEpaAttorneyClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		caseId, err := strToIntOrStatusError(r.FormValue("caseId"))
		if err != nil {
			return err
		}

		ctx := getContext(r)

		data := addEpaAttorneyData{
			XSRFToken: ctx.XSRFToken,
			CaseID:    caseId,
		}

		if r.Method == http.MethodPost {
			attorney := sirius.Attorney{
				Salutation:        postFormString(r, "salutation"),
				FirstName:         postFormString(r, "firstname"),
				MiddleNames:       postFormString(r, "middlenames"),
				Surname:           postFormString(r, "surname"),
				OtherNames:        postFormString(r, "otherNames"),
				DOB:               postFormDateString(r, "dob"),
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
			}
			data.Attorney = attorney

			err := client.CreateAttorney(ctx, caseId, attorney)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
			} else if err != nil {
				return err
			}
			return RedirectError(fmt.Sprintf("/edit-epa?caseId=%d", caseId))
		}

		return tmpl(w, data)
	}
}
