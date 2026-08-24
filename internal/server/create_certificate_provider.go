package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type CreateCertificateProviderClient interface {
	CreateCertificateProvider(ctx sirius.Context, caseId int, certificateProvider sirius.Person) error
	Lpa(sirius.Context, int) (sirius.Lpa, error)
}

type CertificateProviderData struct {
	XSRFToken           string
	CanAddActor         bool
	CaseId              int
	CertificateProvider sirius.Person
	DonorId             int
	Error               sirius.ValidationError
	HtmxRedirect        string
	HtmxSwap            string
	Title               string
	PostURL             string
	IsPartial           bool
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

		caseItem, err := client.Lpa(ctx, caseId)
		if err != nil {
			return err
		}

		data := CertificateProviderData{
			XSRFToken:   ctx.XSRFToken,
			DonorId:     donorId,
			CaseId:      caseId,
			CanAddActor: len(caseItem.CertificateProviders) < 1,
			Title:       "Add a certificate provider",
			PostURL:     fmt.Sprintf("/create-certificate-provider?id=%d&caseId=%d", donorId, caseId),
			IsPartial:   ctx.IsPartial,
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
			} else if err != nil {
				return err
			} else {
				var redirect, swap string

				if r.FormValue("add-another") != "" {
					redirect = fmt.Sprintf("/create-certificate-provider?id=%d&caseId=%d", donorId, caseId)
					swap = "innerHTML scroll:.action-panel__content:top"
				} else {
					redirect = fmt.Sprintf("/create-lpa?id=%d&caseId=%d#accordion-create-lpa-heading-3", donorId, caseId)
					swap = "innerHTML show:#accordion-create-lpa-heading-3:top"
				}

				if data.IsPartial {
					data.HtmxRedirect = redirect
					data.HtmxSwap = swap
					return tmpl(w, data)
				}

				return RedirectError(redirect)
			}
		}

		return tmpl(w, data)
	}
}
