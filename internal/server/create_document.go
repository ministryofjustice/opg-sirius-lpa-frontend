package server

import (
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
)

type CreateDocumentClient interface {
	Case(ctx sirius.Context, id int) (sirius.Case, error)
	RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error)
	DocumentTemplates(ctx sirius.Context, caseType sirius.CaseType) ([]sirius.DocumentTemplateData, error)
	CreateContact(ctx sirius.Context, contact sirius.Person) (sirius.Person, error)
	CreateDocument(ctx sirius.Context, caseID, correspondentID int, templateID string, inserts []string) (sirius.DocumentData, error)
}

type createDocumentData struct {
	XSRFToken               string
	Success                 bool
	Error                   sirius.ValidationError
	Case                    sirius.Case
	Document                sirius.Document
	DocumentTemplates       []sirius.DocumentTemplateData
	DocumentTemplateTypes   []sirius.RefDataItem
	DocumentTemplateRefData []sirius.RefDataItem
	DocumentInsertTypes     []InsertDisplayData
	TemplateSelected        sirius.DocumentTemplateData
	HasViewedInsertPage     bool
	SelectedInserts         []string
	Recipients              []sirius.Person
}

type InsertDisplayData struct {
	Handle string
	Label  string
	Key    string
}

func CreateDocument(client CreateDocumentClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if err := r.ParseForm(); err != nil {
			return err
		}
		ctx := getContext(r)
		data := createDocumentData{
			XSRFToken: ctx.XSRFToken,
		}

		caseID, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}

		caseType, err := sirius.ParseCaseType(r.FormValue("case"))
		if err != nil {
			return err
		}

		group, groupCtx := errgroup.WithContext(ctx.Context)

		group.Go(func() error {
			caseItem, err := client.Case(ctx.With(groupCtx), caseID)
			if err != nil {
				return err
			}

			data.Case = caseItem
			return nil
		})

		group.Go(func() error {
			documentTemplates, err := client.DocumentTemplates(ctx.With(groupCtx), caseType)
			if err != nil {
				return err
			}

			data.DocumentTemplates = documentTemplates
			return nil
		})

		group.Go(func() error {
			documentTemplateRefData, err := client.RefDataByCategory(ctx, sirius.DocumentTemplateIdCategory)
			if err != nil {
				return err
			}
			data.DocumentTemplateRefData = documentTemplateRefData
			return nil
		})

		if err := group.Wait(); err != nil {
			return err
		}

		data.DocumentTemplateTypes = translateDocumentData(data.DocumentTemplates, data.DocumentTemplateRefData)

		templateId := r.FormValue("templateId")
		if templateId != "" {
			for _, dt := range data.DocumentTemplates {
				if dt.TemplateId == templateId {
					data.TemplateSelected = dt
					break
				}
			}
		}
		if data.TemplateSelected.TemplateId != "" {
			data.DocumentInsertTypes = translateInsertData(data.TemplateSelected.Inserts, data.DocumentTemplateRefData)
		}

		hasViewedInsertPage := r.FormValue("hasViewedInserts")
		if hasViewedInsertPage == "true" {
			data.HasViewedInsertPage = true
			inserts := r.Form["insert"]
			if len(inserts) == 0 {
				data.SelectedInserts = []string{}
			} else {
				data.SelectedInserts = inserts
			}
			data.Recipients = getRecipients(data.Case)
		}

		switch r.Method {
		case http.MethodPost:
			recipientControls := postFormString(r, "recipientControls")

			switch recipientControls {
			case "select":
				selectedRecipientID, err := strconv.Atoi(postFormString(r, "selectRecipient"))
				if err != nil {
					return err
				}

				_, err = client.CreateDocument(ctx, caseID, selectedRecipientID, data.TemplateSelected.TemplateId, data.SelectedInserts)
				if err != nil {
					return err
				}

				if ve, ok := err.(sirius.ValidationError); ok {
					w.WriteHeader(http.StatusBadRequest)
					data.Error = ve
				} else if err != nil {
					return err
				} else {
					data.Success = true
					// redirect
				}

			case "generate":
				contact := sirius.Person{
					Salutation:            postFormString(r, "salutation"),
					Firstname:             postFormString(r, "firstname"),
					Middlenames:           postFormString(r, "middlenames"),
					Surname:               postFormString(r, "surname"),
					CompanyName:           postFormString(r, "companyName"),
					CompanyReference:      postFormString(r, "companyReference"),
					AddressLine1:          postFormString(r, "addressLine1"),
					AddressLine2:          postFormString(r, "addressLine2"),
					AddressLine3:          postFormString(r, "addressLine3"),
					Town:                  postFormString(r, "town"),
					County:                postFormString(r, "county"),
					Postcode:              postFormString(r, "postcode"),
					IsAirmailRequired:     postFormString(r, "isAirmailRequired") == "Yes",
					PhoneNumber:           postFormString(r, "phoneNumber"),
					Email:                 postFormString(r, "email"),
					CorrespondenceByPost:  postFormCheckboxChecked(r, "correspondenceBy", "post"),
					CorrespondenceByEmail: postFormCheckboxChecked(r, "correspondenceBy", "email"),
					CorrespondenceByPhone: postFormCheckboxChecked(r, "correspondenceBy", "phone"),
					CorrespondenceByWelsh: postFormCheckboxChecked(r, "correspondenceBy", "welsh"),
				}

				createdContact, err := client.CreateContact(ctx, contact)

				if ve, ok := err.(sirius.ValidationError); ok {
					w.WriteHeader(http.StatusBadRequest)
					data.Error = ve
					data.Recipients = append(data.Recipients, contact)
				} else if err != nil {
					return err
				} else {
					data.Success = true
					data.Recipients = append(data.Recipients, createdContact)
				}
			}
		}

		return tmpl(w, data)
	}
}

func getRecipients(caseItem sirius.Case) []sirius.Person {
	var recipients []sirius.Person
	recipients = append(recipients, *caseItem.Donor)
	recipients = append(recipients, caseItem.TrustCorporations...)
	recipients = append(recipients, caseItem.Attorneys...)

	return recipients
}

func translateDocumentData(documentTemplateData []sirius.DocumentTemplateData, documentTemplateRefData []sirius.RefDataItem) []sirius.RefDataItem {
	var documentTemplateTypes []sirius.RefDataItem
	for _, dt := range documentTemplateData {
		for _, refData := range documentTemplateRefData {
			if refData.Handle == dt.OnScreenSummary {
				translatedRefDataItem := sirius.RefDataItem{
					Handle:         dt.TemplateId,
					Label:          refData.Label,
					UserSelectable: false,
				}
				documentTemplateTypes = append(documentTemplateTypes, translatedRefDataItem)
			}
		}
	}

	return documentTemplateTypes
}

func translateInsertData(selectedTemplateInserts []sirius.Insert, documentTemplateRefData []sirius.RefDataItem) []InsertDisplayData {
	var documentTemplateInserts []InsertDisplayData
	for _, in := range selectedTemplateInserts {
		for _, refData := range documentTemplateRefData {
			if refData.Handle == in.OnScreenSummary {
				translatedRefDataItem := InsertDisplayData{
					Handle: in.InsertId,
					Label:  refData.Label,
					Key:    in.Key,
				}
				documentTemplateInserts = append(documentTemplateInserts, translatedRefDataItem)
			}
		}
	}
	return documentTemplateInserts
}
