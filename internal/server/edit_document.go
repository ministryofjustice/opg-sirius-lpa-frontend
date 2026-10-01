package server

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/sync/errgroup"
)

type EditDocumentClient interface {
	Documents(ctx sirius.Context, caseType sirius.CaseType, caseId int, docTypes []string, notDocTypes []string) ([]sirius.Document, error)
	Case(ctx sirius.Context, id int) (sirius.Case, error)
	CaseSummary(ctx sirius.Context, uid string) (sirius.CaseSummary, error)
	DocumentByUUID(ctx sirius.Context, uuid string) (sirius.Document, error)
	EditDocument(ctx sirius.Context, uuid string, content string) (sirius.Document, error)
	DeleteDocument(ctx sirius.Context, uuid string) error
	AddDocument(ctx sirius.Context, caseID int, document sirius.Document, docType string, blankSections []string) (sirius.Document, error)
	DocumentTemplates(ctx sirius.Context, caseType sirius.CaseType) ([]sirius.DocumentTemplateData, error)
	FeatureToggles(ctx sirius.Context) (sirius.FeatureToggles, error)
}

type editDocumentData struct {
	XSRFToken             string
	IsPartial             bool
	Success               bool
	Error                 sirius.ValidationError
	Case                  sirius.Case
	CaseSummary           sirius.CaseSummary
	Documents             []sirius.Document
	Document              sirius.Document
	UsesNotify            bool
	Download              string
	SaveAndExit           bool
	PreviewDraft          bool
	DownloadUUID          string
	DonorId               int
	HasBlankSections      string
	SelectedBlankSections string
	Section11Count1       string
	Section11Count2       string
	BlankSectionsEnabled  bool
	AttorneyCount         int
}

func publishDraftDocument(
	client EditDocumentClient,
	ctx sirius.Context,
	caseID int,
	documentUUID string,
	content string,
	blankSections []string,
) error {
	_, err := client.EditDocument(ctx, documentUUID, content)
	if err != nil {
		return err
	}

	// need to retrieve for correspondent information
	document, err := client.DocumentByUUID(ctx, documentUUID)
	if err != nil {
		return err
	}

	_, err = client.AddDocument(ctx, caseID, document, sirius.TypeSave, blankSections)
	if err != nil {
		return err
	}

	err = client.DeleteDocument(ctx, documentUUID)
	if err != nil {
		return err
	}

	return nil
}

func parseBlankSections(hasBlankSections bool, selectedBlankSections, section11Count1, section11Count2, caseSubType string) ([]string, error) {
	if !hasBlankSections {
		return []string{}, nil
	}

	if selectedBlankSections == "" {
		return nil, sirius.ValidationError{
			Field: sirius.FieldErrors{
				"blankSections": {"reason": "Please select sections to insert"},
			},
		}
	}

	var containsSection11 bool
	blankSections := strings.Split(selectedBlankSections, "+")
	for i, blankSection := range blankSections {
		if blankSection == "11" {
			containsSection11 = true
		}
		blankSections[i] = caseSubType + "-" + blankSection
	}

	if containsSection11 && selectedBlankSections == "11+15" && section11Count1 == "" {
		return nil, sirius.ValidationError{
			Field: sirius.FieldErrors{
				"section11Count1": {"reason": "Please select how many section 11 to insert"},
			},
		}
	}

	if containsSection11 && selectedBlankSections == "10+11+15" && section11Count2 == "" {
		return nil, sirius.ValidationError{
			Field: sirius.FieldErrors{
				"section11Count2": {"reason": "Please select how many section 11 to insert"},
			},
		}
	}

	var section11Count string
	if selectedBlankSections == "11+15" {
		section11Count = section11Count1
	} else if selectedBlankSections == "10+11+15" {
		section11Count = section11Count2
	}
	section11CountInt, _ := strconv.Atoi(section11Count)
	for range section11CountInt - 1 {
		blankSections = append(blankSections, caseSubType+"-11")
	}
	slices.Sort(blankSections)

	return blankSections, nil
}

func getAttorneyCount(caseItem sirius.Case) int {
	var count int

	for _, attorney := range caseItem.Attorneys {
		if attorney.SystemStatus != nil && *attorney.SystemStatus {
			count++
		}
	}

	for range caseItem.ReplacementAttorneys {
		count++
	}

	for _, trustCorporation := range caseItem.TrustCorporations {
		if trustCorporation.IsReplacementAttorney || (trustCorporation.SystemStatus != nil && *trustCorporation.SystemStatus) {
			count++
		}
	}

	return count
}

func EditDocument(client EditDocumentClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		caseID, err := strToIntOrStatusError(r.FormValue("id"))
		if err != nil {
			return err
		}

		caseType, err := sirius.ParseCaseType(r.FormValue("case"))
		if err != nil {
			return err
		}

		data := editDocumentData{
			XSRFToken: ctx.XSRFToken,
			IsPartial: ctx.IsPartial,
		}

		featureToggles, err := client.FeatureToggles(ctx)
		if err == nil {
			data.BlankSectionsEnabled = featureToggles.Enabled("poasBlankSections")
		}

		caseItem, err := client.Case(ctx, caseID)
		if err != nil {
			return err
		}
		data.Case = caseItem
		data.DonorId = caseItem.Donor.ID
		data.AttorneyCount = getAttorneyCount(caseItem)

		var documentTemplates []sirius.DocumentTemplateData

		switch r.Method {
		case http.MethodGet:
			group, groupCtx := errgroup.WithContext(ctx.Context)

			if caseType == sirius.CaseTypeDigitalLpa {
				group.Go(func() error {
					data.CaseSummary, err = client.CaseSummary(ctx, data.Case.UID)
					if err != nil {
						return err
					}
					return nil
				})
			}

			group.Go(func() error {
				documents, err := client.Documents(ctx.With(groupCtx), caseType, caseID, []string{sirius.TypeDraft}, []string{})
				if err != nil {
					return err
				}
				data.Documents = documents
				return nil
			})

			group.Go(func() error {
				documentTemplates, err = client.DocumentTemplates(ctx, caseType)
				if err != nil {
					return err
				}

				return nil
			})

			if err := group.Wait(); err != nil {
				return err
			}

			if len(data.Documents) > 0 {
				defaultDocumentUUID := data.Documents[0].UUID
				selectedDocumentUUID := r.FormValue("document")
				if selectedDocumentUUID != "" {
					defaultDocumentUUID = selectedDocumentUUID
				}
				document, err := client.DocumentByUUID(ctx, defaultDocumentUUID)
				if err != nil {
					return err
				}
				data.Document = document
			}
		case http.MethodPost:
			documentControls := postFormString(r, "documentControls")
			content := r.FormValue("documentTextEditor")
			documentUUID := r.FormValue("documentUUID")

			data.HasBlankSections = r.FormValue("hasBlankSections")
			data.SelectedBlankSections = r.FormValue("blankSections")
			data.Section11Count1 = r.FormValue("section11Count1")
			data.Section11Count2 = r.FormValue("section11Count2")
			hasBlankSections := data.HasBlankSections == "true"

			switch documentControls {
			case "save":
				document, err := client.EditDocument(ctx, documentUUID, content)
				if err != nil {
					return err
				}
				data.Document = document

			case "preview":
				blankSections, err := parseBlankSections(hasBlankSections, data.SelectedBlankSections, data.Section11Count1, data.Section11Count2, caseItem.SubType)
				if ve, ok := err.(sirius.ValidationError); ok {
					w.WriteHeader(http.StatusBadRequest)
					data.Error = ve
				} else {
					_, err = client.EditDocument(ctx, documentUUID, content)
					if err != nil {
						return err
					}

					// need to retrieve for correspondent information
					document, err := client.DocumentByUUID(ctx, documentUUID)
					if err != nil {
						return err
					}

					previewDocument, err := client.AddDocument(ctx, caseID, document, sirius.TypePreview, blankSections)
					if err != nil {
						return err
					}

					data.Document = document
					data.PreviewDraft = true
					data.DownloadUUID = previewDocument.UUID
				}

			case "delete":
				err := client.DeleteDocument(ctx, documentUUID)
				if err != nil {
					return err
				}

				if caseType == sirius.CaseTypeDigitalLpa {
					SetFlash(w, FlashNotification{
						Title: "Draft document deleted",
					})

					return RedirectError(fmt.Sprintf("/lpa/%s/documents", caseItem.UID))
				}

			case "publish":
				blankSections, err := parseBlankSections(hasBlankSections, data.SelectedBlankSections, data.Section11Count1, data.Section11Count2, caseItem.SubType)
				if err == nil {
					err = publishDraftDocument(client, ctx, caseID, documentUUID, content, blankSections)
				}
				if ve, ok := err.(sirius.ValidationError); ok {
					w.WriteHeader(http.StatusBadRequest)
					data.Error = ve
				} else if err != nil {
					return err
				} else {
					if caseType == sirius.CaseTypeDigitalLpa {
						SetFlash(w, FlashNotification{
							Title: "Document published",
						})

						return RedirectError(fmt.Sprintf("/lpa/%s/documents", caseItem.UID))
					}

					data.Success = true
				}

			case "saveAndExit":
				_, err := client.EditDocument(ctx, documentUUID, content)
				if err != nil {
					return err
				}

				if caseType == sirius.CaseTypeDigitalLpa {
					SetFlash(w, FlashNotification{
						Title: "Document saved",
					})

					return RedirectError(fmt.Sprintf("/lpa/%s/documents", caseItem.UID))
				}

				data.SaveAndExit = true
			}

			if !data.SaveAndExit {
				group, groupCtx := errgroup.WithContext(ctx.Context)

				group.Go(func() error {
					if caseType == sirius.CaseTypeDigitalLpa {
						data.CaseSummary, err = client.CaseSummary(ctx.With(groupCtx), data.Case.UID)
						if err != nil {
							return err
						}
					}

					return nil
				})

				group.Go(func() error {
					data.Documents, err = client.Documents(ctx.With(groupCtx), caseType, caseID, []string{sirius.TypeDraft}, []string{})
					if err != nil {
						return err
					}
					return nil
				})

				group.Go(func() error {
					documentTemplates, err = client.DocumentTemplates(ctx, caseType)
					if err != nil {
						return err
					}
					return nil
				})

				if err := group.Wait(); err != nil {
					return err
				}

				if documentControls == "delete" || documentControls == "publish" || documentControls == "preview" {
					if len(data.Documents) > 0 {
						defaultDocumentUUID := data.Documents[0].UUID
						documentToDisplay, err := client.DocumentByUUID(ctx, defaultDocumentUUID)
						if err != nil {
							return err
						}
						data.Document = documentToDisplay
					}
				}
			}
		}

		if len(documentTemplates) > 0 && data.Document.SystemType != "" {
			for _, dt := range documentTemplates {
				if dt.TemplateId == data.Document.SystemType {
					data.UsesNotify = dt.UsesNotify

					break
				}
			}
		}

		return tmpl(w, data)
	}
}
