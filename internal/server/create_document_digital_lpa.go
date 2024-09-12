package server

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/sync/errgroup"
)

type CreateDocumentDigitalLpaClient interface {
	CaseSummary(ctx sirius.Context, uid string) (sirius.CaseSummary, error)
	DocumentTemplates(ctx sirius.Context, caseType sirius.CaseType) ([]sirius.DocumentTemplateData, error)
	CreateDocument(ctx sirius.Context, caseID, correspondentID int, templateID string, inserts []string) (sirius.Document, error)
	CreateContact(ctx sirius.Context, contact sirius.Person) (sirius.Person, error)
}

type createDocumentDigitalLpaData struct {
	XSRFToken             string
	Error                 sirius.ValidationError
	CaseSummary           sirius.CaseSummary
	DocumentTemplates     []sirius.DocumentTemplateData
	ComponentDocumentData ComponentDocumentData
	Recipients            []sirius.Person
	SelectedTemplateId    string
	SelectedInserts       []string
	SelectedRecipients    []int
}

func CreateDocumentDigitalLpa(client CreateDocumentDigitalLpaClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		var err error
		if err := r.ParseForm(); err != nil {
			return err
		}
		ctx := getContext(r)

		uid := chi.URLParam(r, "uid")

		data := createDocumentDigitalLpaData{
			XSRFToken: ctx.XSRFToken,
		}

		data.CaseSummary, err = client.CaseSummary(ctx, uid)

		if err != nil {
			return err
		}

		group, _ := errgroup.WithContext(ctx.Context)

		group.Go(func() error {
			documentTemplates, err := client.DocumentTemplates(ctx, sirius.CaseTypeDigitalLpa)
			data.DocumentTemplates = sortDocumentData(documentTemplates)

			if err != nil {
				return err
			}

			return nil
		})

		if err := group.Wait(); err != nil {
			return err
		}

		data.ComponentDocumentData = buildComponentDocumentData(data.DocumentTemplates)

		lpa := data.CaseSummary.DigitalLpa

		placeholderRecipientId := -1

		donorRecipient := sirius.Person{
			ID:           placeholderRecipientId,
			UID:          lpa.LpaStoreData.Donor.Uid,
			Firstname:    lpa.LpaStoreData.Donor.FirstNames,
			Surname:      lpa.LpaStoreData.Donor.LastName,
			DateOfBirth:  sirius.DateString(lpa.LpaStoreData.Donor.DateOfBirth),
			PersonType:   "Donor",
			AddressLine1: lpa.LpaStoreData.Donor.Address.Line1,
			AddressLine2: lpa.LpaStoreData.Donor.Address.Line2,
			AddressLine3: lpa.LpaStoreData.Donor.Address.Line3,
			Town:         lpa.LpaStoreData.Donor.Address.Town,
			Postcode:     lpa.LpaStoreData.Donor.Address.Postcode,
			Country:      lpa.LpaStoreData.Donor.Address.Country,
			Email:        lpa.LpaStoreData.Donor.Email,
			AlsoKnownAs:  lpa.LpaStoreData.Donor.OtherNamesKnownBy,
		}

		data.Recipients = append(data.Recipients, donorRecipient)

		if lpa.LpaStoreData.CertificateProvider.Uid != "" {
			placeholderRecipientId--

			certificateProviderRecipient := sirius.Person{
				ID:           placeholderRecipientId,
				UID:          lpa.LpaStoreData.CertificateProvider.Uid,
				Firstname:    lpa.LpaStoreData.CertificateProvider.FirstNames,
				Surname:      lpa.LpaStoreData.CertificateProvider.LastName,
				PersonType:   "Certificate Provider",
				AddressLine1: lpa.LpaStoreData.CertificateProvider.Address.Line1,
				AddressLine2: lpa.LpaStoreData.CertificateProvider.Address.Line2,
				AddressLine3: lpa.LpaStoreData.CertificateProvider.Address.Line3,
				Town:         lpa.LpaStoreData.CertificateProvider.Address.Town,
				Postcode:     lpa.LpaStoreData.CertificateProvider.Address.Postcode,
				Country:      lpa.LpaStoreData.CertificateProvider.Address.Country,
				Email:        lpa.LpaStoreData.CertificateProvider.Email,
				PhoneNumber:  lpa.LpaStoreData.CertificateProvider.Phone,
			}

			data.Recipients = append(data.Recipients, certificateProviderRecipient)
		}

		for _, attorney := range lpa.LpaStoreData.Attorneys {
			if attorney.Status == "removed" {
				continue
			}

			personType := "Attorney"

			if attorney.Status == "replacement" {
				personType = "Replacement Attorney"
			}

			placeholderRecipientId--

			recipient := sirius.Person{
				ID:           placeholderRecipientId,
				UID:          attorney.Uid,
				Firstname:    attorney.FirstNames,
				Surname:      attorney.LastName,
				DateOfBirth:  sirius.DateString(attorney.DateOfBirth),
				PersonType:   personType,
				AddressLine1: attorney.Address.Line1,
				AddressLine2: attorney.Address.Line2,
				AddressLine3: attorney.Address.Line3,
				Town:         attorney.Address.Town,
				Postcode:     attorney.Address.Postcode,
				Country:      attorney.Address.Country,
				Email:        attorney.Email,
				PhoneNumber:  attorney.Mobile,
			}

			data.Recipients = append(data.Recipients, recipient)
		}

		if r.Method == "POST" {
			// set data
			data.SelectedTemplateId = r.FormValue("templateId")
			data.SelectedInserts = r.Form["insert"]

			for _, recipientId := range r.Form["selectRecipients"] {
				recipientIdInt, _ := strconv.Atoi(recipientId)
				data.SelectedRecipients = append(data.SelectedRecipients, recipientIdInt)
			}

			// check data
			fieldErrors := sirius.FieldErrors{}
			if data.SelectedTemplateId == "" {
				fieldErrors["templateId"] = map[string]string{"reason": "Please select a document template to continue"}
			}

			if len(data.SelectedRecipients) == 0 {
				fieldErrors["selectRecipient"] = map[string]string{"reason": "Please select a recipient to continue"}
			}

			if len(fieldErrors) > 0 {
				data.Error = sirius.ValidationError{
					Field: fieldErrors,
				}

				return tmpl(w, data)
			}

			// save draft document
			for _, recipient := range data.Recipients {
				if !slices.Contains(data.SelectedRecipients, recipient.ID) {
					continue
				}

				if recipient.ID < 0 {
					// Create a placeholder contact for each recipient as we currently have no way to link a document to
					// LPA store actors
					recipient, err = client.CreateContact(ctx, recipient)
					if err != nil {
						return err
					}
				}

				_, err = client.CreateDocument(ctx, lpa.SiriusData.ID, recipient.ID, data.SelectedTemplateId, data.SelectedInserts)
				if err != nil {
					return err
				}

				if ve, ok := err.(sirius.ValidationError); ok {
					w.WriteHeader(http.StatusBadRequest)
					data.Error = ve
				} else if err != nil {
					return err
				} else {
					return RedirectError(fmt.Sprintf("/edit-document?id=%d&case=%s", lpa.SiriusData.ID, sirius.CaseTypeDigitalLpa))
				}
			}
		}

		return tmpl(w, data)
	}
}
