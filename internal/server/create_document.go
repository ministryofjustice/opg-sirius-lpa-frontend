package server

import (
	"fmt"
	"net/http"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type CreateDocumentClient interface {
	Case(ctx sirius.Context, id int) (sirius.Case, error)
	Person(ctx sirius.Context, id int) (sirius.Person, error)
	DocumentTemplates(ctx sirius.Context, caseType sirius.CaseType) ([]sirius.DocumentTemplateData, error)
	CreateContact(ctx sirius.Context, contact sirius.Person) (sirius.Person, error)
	CreateDocument(ctx sirius.Context, caseID, correspondentID int, templateID string, inserts []string) (sirius.Document, error)
}

type createDocumentData struct {
	XSRFToken                  string
	RecipientAddedSuccess      bool
	Error                      sirius.ValidationError
	Case                       sirius.Case
	Document                   sirius.Document
	DocumentTemplates          []sirius.DocumentTemplateData
	DocumentInsertTypes        []InsertDisplayData
	TemplateSelected           sirius.DocumentTemplateData
	HasViewedInsertPage        bool
	SelectedInserts            []string
	Recipients                 []sirius.Person
	DocumentInsertKeys         []string
	ComponentDocumentData      ComponentDocumentData
	HasSelectedAddNewRecipient bool
	Back                       string
}

type InsertDisplayData struct {
	Handle string
	Label  string
	Key    string
}

type ComponentTemplateInsertData struct {
	InsertId string `json:"id"`
	Label    string `json:"label"`
}

type ComponentTemplateData struct {
	Id      string                                   `json:"id"`
	Inserts map[string][]ComponentTemplateInsertData `json:"inserts"`
}

type ComponentDocumentData struct {
	Templates    []ComponentTemplateData `json:"templates"`
	Translations map[string]string       `json:"translations"`
}

func CreateDocument(client CreateDocumentClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if err := r.ParseForm(); err != nil {
			return err
		}
		ctx := getContext(r)

		caseID, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}

		caseType, err := sirius.ParseCaseType(r.FormValue("case"))
		if err != nil {
			return err
		}

		caseItem, err := client.Case(ctx, caseID)
		if err != nil {
			return err
		}

		data := createDocumentData{
			XSRFToken: ctx.XSRFToken,
			Case:      caseItem,
		}

		switch r.Method {
		case http.MethodGet:
			documentTemplates, err := client.DocumentTemplates(ctx, caseType)
			if err != nil {
				return err
			}

			data.DocumentTemplates = sortDocumentData(documentTemplates)

			templateId := r.FormValue("templateId")
			hasSelectedSubmitTemplate := r.FormValue("hasSelectedSubmitTemplate")

			if templateId != "" {
				for _, dt := range data.DocumentTemplates {
					if dt.TemplateId == templateId {
						data.TemplateSelected = dt
						break
					}
				}
			} else if hasSelectedSubmitTemplate == "true" {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = sirius.ValidationError{
					Field: sirius.FieldErrors{
						"templateId": {"reason": "Please select a document template to continue"},
					},
				}
			}

			if data.TemplateSelected.TemplateId != "" {
				data.DocumentInsertTypes = translateInsertData(data.TemplateSelected.Inserts)
				data.DocumentInsertKeys = getSortedInsertKeys(data.TemplateSelected.Inserts)
			} else {
				data.ComponentDocumentData = buildComponentDocumentData(data.DocumentTemplates)
			}

			hasViewedInsertPage := r.FormValue("hasViewedInserts")
			hasNoInsertsToSelect := data.TemplateSelected.TemplateId != "" && len(data.DocumentInsertTypes) == 0

			if r.FormValue("skipInserts") == "" {
				uniqueInserts := removeDuplicateStr(r.Form["insert"])
				data.SelectedInserts = uniqueInserts
			}

			if hasViewedInsertPage == "true" || hasNoInsertsToSelect {
				data.HasViewedInsertPage = true
				data.Recipients, err = getRecipients(data.Case)
				if err != nil {
					return err
				}
			}

			hasSelectedAddNewRecipient := r.FormValue("hasSelectedAddNewRecipient")
			if hasSelectedAddNewRecipient == "true" {
				data.HasSelectedAddNewRecipient = true
			}

			data.Back = getBackUrl(data)

		case http.MethodPost:
			recipientControls := postFormString(r, "recipientControls")

			switch recipientControls {
			case "selectRecipients":
				selectedRecipientIDs, err := sliceAtoi(r.Form["selectRecipients"])
				if err != nil {
					return err
				}

				templateId := r.FormValue("templateId")

				uniqueInserts := []string{}
				if r.FormValue("skipInserts") == "" {
					uniqueInserts = removeDuplicateStr(r.Form["insert"])
				}

				sort.SliceStable(uniqueInserts, func(i, j int) bool {
					return strings.HasPrefix(uniqueInserts[i], "IN") && !strings.HasPrefix(uniqueInserts[j], "IN")
				})

				for _, recipientID := range selectedRecipientIDs {
					document, err := client.CreateDocument(ctx, caseID, recipientID, templateId, uniqueInserts)
					if err != nil {
						return err
					}
					data.Document = document
				}

				if ve, ok := err.(sirius.ValidationError); ok {
					w.WriteHeader(http.StatusBadRequest)
					data.Error = ve
				} else if err != nil {
					return err
				} else {
					return RedirectError(fmt.Sprintf("/edit-document?id=%d&case=%s", caseID, caseType))
				}
			case "addNewRecipient":
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
				} else if err != nil {
					return err
				} else {
					data.RecipientAddedSuccess = true
					data.TemplateSelected.TemplateId = r.FormValue("templateId")
					data.HasViewedInsertPage = true
					data.Recipients, err = getRecipients(data.Case)
					if err != nil {
						return err
					}
					data.Recipients = append(data.Recipients, createdContact)

					if r.FormValue("skipInserts") == "" {
						data.SelectedInserts = r.Form["insert"]
					}
				}
			}
		}

		return tmpl(w, data)
	}
}

func getRecipients(caseItem sirius.Case) ([]sirius.Person, error) {
	caseItem = caseItem.FilterInactiveAttorneys()

	var recipients []sirius.Person
	recipients = append(recipients, *caseItem.Donor)

	if caseItem.Correspondent != nil {
		recipients = append(recipients, *caseItem.Correspondent)
	}

	sort.Slice(caseItem.Attorneys, func(i, j int) bool {
		if caseItem.Attorneys[i].Surname == caseItem.Attorneys[j].Surname {
			return caseItem.Attorneys[i].Firstname < caseItem.Attorneys[j].Firstname
		}
		return caseItem.Attorneys[i].Surname < caseItem.Attorneys[j].Surname
	})

	recipients = append(recipients, caseItem.Attorneys...)
	recipients = append(recipients, caseItem.TrustCorporations...)

	return recipients, nil
}

func sortDocumentData(documentTemplateData []sirius.DocumentTemplateData) []sirius.DocumentTemplateData {
	sort.Slice(documentTemplateData, func(i, j int) bool {
		return documentTemplateData[i].TemplateId < documentTemplateData[j].TemplateId
	})

	return documentTemplateData
}

func translateInsertData(selectedTemplateInserts []sirius.Insert) []InsertDisplayData {
	var documentTemplateInserts []InsertDisplayData
	for _, in := range selectedTemplateInserts {
		translatedRefDataItem := InsertDisplayData{
			Handle: in.InsertId,
			Label:  in.Label,
			Key:    in.Key,
		}
		documentTemplateInserts = append(documentTemplateInserts, translatedRefDataItem)
	}

	return documentTemplateInserts
}

func getSortedInsertKeys(selectedTemplateInserts []sirius.Insert) []string {
	var documentInsertKeys []string
	for _, in := range selectedTemplateInserts {
		if !slices.Contains(documentInsertKeys, in.Key) {
			documentInsertKeys = append(documentInsertKeys, in.Key)
		}
	}
	// insert api response infrequently includes the key all
	// note insert keys are lowercase so this check is case-sensitive
	if !slices.Contains(documentInsertKeys, "all") {
		documentInsertKeys = append(documentInsertKeys, "all")
	}

	sort.Strings(documentInsertKeys)
	return documentInsertKeys
}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]struct{})
	var list []string
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = struct{}{}
			list = append(list, item)
		}
	}
	return list
}

func sliceAtoi(strSlice []string) ([]int, error) {
	intSlice := make([]int, len(strSlice))

	for idx, s := range strSlice {
		intEl, err := strconv.Atoi(s)
		if err != nil {
			return intSlice, err
		}
		intSlice[idx] = intEl
	}

	return intSlice, nil
}

func getBackUrl(data createDocumentData) string {
	url := fmt.Sprintf("/create-document?id=%d&case=%s", data.Case.ID, data.Case.CaseType)
	if data.HasSelectedAddNewRecipient {
		url = fmt.Sprintf("%s&templateId=%s&hasViewedInserts=true", url, data.TemplateSelected.TemplateId)
	} else if data.HasViewedInsertPage {
		if len(data.DocumentInsertTypes) != 0 {
			url = fmt.Sprintf("%s&templateId=%s", url, data.TemplateSelected.TemplateId)
		}
	}
	if len(data.SelectedInserts) > 0 {
		for _, i := range data.SelectedInserts {
			url = fmt.Sprintf("%s&insert=%s", url, i)
		}
	}
	return url
}

func buildComponentDocumentData(templates []sirius.DocumentTemplateData) ComponentDocumentData {
	data := ComponentDocumentData{
		Translations: map[string]string{},
	}

	for _, template := range templates {
		formed := ComponentTemplateData{
			Id:      template.TemplateId,
			Inserts: map[string][]ComponentTemplateInsertData{},
		}

		for _, insert := range template.Inserts {
			formedInsert := ComponentTemplateInsertData{
				InsertId: insert.InsertId,
				Label:    insert.Label,
			}

			formed.Inserts[insert.Key] = append(formed.Inserts[insert.Key], formedInsert)
		}

		data.Templates = append(data.Templates, formed)
	}

	return data
}
