package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type CreateCertificateProviderClient interface {
	CreateCertificateProvider(ctx sirius.Context, caseId int, certificateProvider sirius.Person) error
}

type CreateCertificateProviderData struct {
	XSRFToken           string
	DonorId             int
	CaseId              int
	Lpa                 sirius.Lpa
	CertificateProvider sirius.Person
	Error               sirius.ValidationError
}

func CreateCertificateProvider(client CreateCertificateProviderClient, tmpl template.Template) Handler {
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

		data := CreateCertificateProviderData{
			XSRFToken: ctx.XSRFToken,
			DonorId:   donorId,
			CaseId:    caseId,
		}

		if r.Method == http.MethodPost {

			certificateProvider := sirius.Person{
				Salutation:   postFormString(r, "salutation"),
				Firstname:    postFormString(r, "firstname"),
				Middlenames:  postFormString(r, "middlenames"),
				Surname:      postFormString(r, "surname"),
				AddressLine1: postFormString(r, "addressLine1"),
				AddressLine2: postFormString(r, "addressLine2"),
				AddressLine3: postFormString(r, "addressLine3"),
				Town:         postFormString(r, "town"),
				County:       postFormString(r, "county"),
				Country:      postFormString(r, "country"),
				Postcode:     postFormString(r, "postcode"),
			}

			err = client.CreateCertificateProvider(ctx, caseId, certificateProvider)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
				//if r.Header.Get("HX-Request") == "true" {
				//	return partialTmpl(w, data)
				//}
				//return tmpl(w, data)
			} else if err != nil {
				return err
			} else {
				return RedirectError(fmt.Sprintf("/create-lpa?id=%d&caseId=%d", donorId, caseId))
			}

			//return RedirectError(fmt.Sprintf("/create-correspondent?id=%d&caseId=%d", donorId, caseId))
		}

		return tmpl(w, data)
	}
}
