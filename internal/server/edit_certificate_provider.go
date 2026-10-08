package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type EditCertificateProviderClient interface {
	UpdateCertificateProvider(ctx sirius.Context, certificateProviderId int, certificateProvider sirius.Person) error
	Person(sirius.Context, int) (sirius.Person, error)
}

func EditCertificateProvider(client EditCertificateProviderClient, tmpl template.Template) Handler {
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

		personId, err := strToIntOrStatusError(r.FormValue("personId"))
		if err != nil {
			return err
		}

		certificateProvider, err := client.Person(ctx, personId)
		if err != nil {
			return err
		}

		data := CertificateProviderData{
			XSRFToken:           ctx.XSRFToken,
			FlowQuery:           getFlowQuery(r),
			DonorId:             donorId,
			CaseId:              caseId,
			CanAddActor:         false,
			CertificateProvider: certificateProvider,
			Title:               "Edit a certificate provider",
			PostURL:             fmt.Sprintf("/edit-certificate-provider?id=%d&caseId=%d&personId=%d%s", donorId, caseId, personId, getFlowQuery(r)),
			IsPartial:           r.Header.Get("HX-Request") == "true",
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

			err = client.UpdateCertificateProvider(ctx, personId, certificateProvider)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
			} else if err != nil {
				return err
			} else {
				if data.IsPartial {
					data.HtmxRedirect = fmt.Sprintf("/create-lpa?id=%d&caseId=%d%s#accordion-create-lpa-heading-3", donorId, caseId, data.FlowQuery)
					data.HtmxSwap = "innerHTML show:#accordion-create-lpa-heading-3:top"
					return tmpl(w, data)
				}

				return RedirectError(fmt.Sprintf("/create-lpa?id=%d&caseId=%d%s#accordion-create-lpa-heading-3", donorId, caseId, data.FlowQuery))
			}
		}

		return tmpl(w, data)
	}
}
