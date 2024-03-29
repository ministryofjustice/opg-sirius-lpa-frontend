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

		donorRecipient := sirius.Person{
			ID:           lpa.SiriusData.Donor.ID,
			Firstname:    lpa.SiriusData.Donor.Firstname,
			Surname:      lpa.SiriusData.Donor.Surname,
			PersonType:   lpa.SiriusData.Donor.PersonType,
			AddressLine1: lpa.SiriusData.Donor.AddressLine1,
			AddressLine2: lpa.SiriusData.Donor.AddressLine2,
			AddressLine3: lpa.SiriusData.Donor.AddressLine3,
			Town:         lpa.SiriusData.Donor.Town,
			Postcode:     lpa.SiriusData.Donor.Postcode,
			Country:      lpa.SiriusData.Donor.Country,
		}
		data.Recipients = append(data.Recipients, donorRecipient)

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
