package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/exp/slices"
	"golang.org/x/sync/errgroup"
	"net/http"
	"sort"
	"strconv"
	"sync"
)

type CreateDocumentClient interface {
	Case(ctx sirius.Context, id int) (sirius.Case, error)
	Person(ctx sirius.Context, id int) (sirius.Person, error)
	RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error)
	DocumentTemplates(ctx sirius.Context, caseType sirius.CaseType) ([]sirius.DocumentTemplateData, error)
	CreateContact(ctx sirius.Context, contact sirius.Person) (sirius.Person, error)
	CreateDocument(ctx sirius.Context, caseID, correspondentID int, templateID string, inserts []string) (sirius.Document, error)
}

type createDocumentData struct {
	XSRFToken                  string
	Success                    bool
	Error                      sirius.ValidationError
	Case                       sirius.Case
	Document                   sirius.Document
	DocumentTemplates          []sirius.DocumentTemplateData
	DocumentTemplateTypes      []sirius.RefDataItem
	DocumentTemplateRefData    []sirius.RefDataItem
	DocumentInsertTypes        []InsertDisplayData
	TemplateSelected           sirius.DocumentTemplateData
	HasViewedInsertPage        bool
	SelectedInserts            []string
	Recipients                 []sirius.Person
	DocumentInsertKeys         []string
	HasSelectedAddNewRecipient bool
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
			group, groupCtx := errgroup.WithContext(ctx.Context)

			group.Go(func() error {
				documentTemplates, err := client.DocumentTemplates(ctx.With(groupCtx), caseType)
				if err != nil {
					return err
				}

				data.DocumentTemplates = documentTemplates
				return nil
			})

			group.Go(func() error {
				documentTemplateRefData, err := client.RefDataByCategory(ctx.With(groupCtx), sirius.DocumentTemplateIdCategory)
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
				data.DocumentInsertTypes = translateInsertData(data.TemplateSelected.Inserts, data.DocumentTemplateRefData)
				data.DocumentInsertKeys = getSortedInsertKeys(data.TemplateSelected.Inserts)
			}

			hasViewedInsertPage := r.FormValue("hasViewedInserts")
			hasNoInsertsToSelect := data.TemplateSelected.TemplateId != "" && len(data.DocumentInsertTypes) == 0
			uniqueInserts := removeDuplicateStr(r.Form["insert"])
			data.SelectedInserts = uniqueInserts

			if hasViewedInsertPage == "true" || hasNoInsertsToSelect {
				data.HasViewedInsertPage = true
				data.Recipients, err = getRecipients(ctx, client, data.Case)
				if err != nil {
					return err
				}
			}

			hasSelectedAddNewRecipient := r.FormValue("hasSelectedAddNewRecipient")
			if hasSelectedAddNewRecipient == "true" {
				data.HasSelectedAddNewRecipient = true
			}

		case http.MethodPost:
			selectedRecipientIDs, err := sliceAtoi(r.Form["selectRecipients"])

			if len(selectedRecipientIDs) > 0 {
				templateId := r.FormValue("templateId")
				uniqueInserts := removeDuplicateStr(r.Form["insert"])

				if len(uniqueInserts) == 0 {
					uniqueInserts = []string{}
				}

				if err != nil {
					return err
				}
				if len(selectedRecipientIDs) == 0 {
					w.WriteHeader(http.StatusBadRequest)
					data.Error = sirius.ValidationError{
						Field: sirius.FieldErrors{
							"selectRecipient": {"reason": "Value is required and can't be empty"},
						},
					}
					data.TemplateSelected.TemplateId = r.FormValue("templateId")
					data.SelectedInserts = uniqueInserts
					data.HasViewedInsertPage = true
					data.Recipients, err = getRecipients(ctx, client, data.Case)
					if err != nil {
						return err
					}
					return tmpl(w, data)
				}

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
					data.Success = true
					return RedirectError(fmt.Sprintf("/edit-document?id=%d&case=%s", caseID, caseType))
				}
			} else {
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
					data.Success = true
					data.TemplateSelected.TemplateId = r.FormValue("templateId")
					data.SelectedInserts = r.Form["insert"]
					data.HasViewedInsertPage = true
					data.Recipients, err = getRecipients(ctx, client, data.Case)
					if err != nil {
						return err
					}
					data.Recipients = append(data.Recipients, createdContact)
				}
			}
		}

		return tmpl(w, data)
	}
}

func getRecipients(ctx sirius.Context, client CreateDocumentClient, caseItem sirius.Case) ([]sirius.Person, error) {
	var recipientIds []int
	donor := *caseItem.Donor
	recipientIds = append(recipientIds, donor.ID)
	recipientIds = append(recipientIds, getPersonIds(caseItem.TrustCorporations)...)
	recipientIds = append(recipientIds, getPersonIds(caseItem.Attorneys)...)

	var recipients []sirius.Person
	group, groupCtx := errgroup.WithContext(ctx.Context)
	var personsMu sync.Mutex

	for _, recipientID := range recipientIds {
		recipientID := recipientID

		group.Go(func() error {
			person, err := client.Person(ctx.With(groupCtx), recipientID)
			if err != nil {
				return err
			}

			personsMu.Lock()
			recipients = append(recipients, person)
			personsMu.Unlock()
			return nil
		})
	}

	if err := group.Wait(); err != nil {
		return nil, err
	}

	return recipients, nil
}

func getPersonIds(persons []sirius.Person) []int {
	var personIds []int
	for _, person := range persons {
		personIds = append(personIds, person.ID)
	}
	return personIds
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
				break
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
				break
			}
		}
	}

	sort.Slice(documentTemplateInserts, func(i, j int) bool {
		return documentTemplateInserts[i].Handle < documentTemplateInserts[j].Handle
	})

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
