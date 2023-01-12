package server

import (
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
)

type EditDocumentClient interface {
	Documents(ctx sirius.Context, caseType sirius.CaseType, caseId int) ([]sirius.Document, error)
	Case(ctx sirius.Context, id int) (sirius.Case, error)
	DocumentByUUID(ctx sirius.Context, uuid string) (sirius.Document, error)
	EditDocument(ctx sirius.Context, uuid string, content string) (sirius.Document, error)
	DeleteDocument(ctx sirius.Context, uuid string) error
	AddDocument(ctx sirius.Context, caseID int, document sirius.Document, docType string) (sirius.Document, error)
}

type editDocumentData struct {
	XSRFToken    string
	Success      bool
	Error        sirius.ValidationError
	Case         sirius.Case
	Documents    []sirius.Document
	Document     sirius.Document
	Download     string
	SaveAndExit  bool
	PreviewDraft bool
	DownloadUUID string
}

func EditDocument(client EditDocumentClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		caseID, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}

		caseType, err := sirius.ParseCaseType(r.FormValue("case"))
		if err != nil {
			return err
		}

		data := editDocumentData{
			XSRFToken: ctx.XSRFToken,
		}

		switch r.Method {
		case http.MethodGet:
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
				documents, err := client.Documents(ctx.With(groupCtx), caseType, caseID)
				if err != nil {
					return err
				}
				draftDocuments := getDraftDocuments(documents)
				data.Documents = draftDocuments
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

			switch documentControls {
			case "save":
				document, err := client.EditDocument(ctx, documentUUID, content)
				if err != nil {
					return err
				}
				data.Document = document

			case "preview":
				_, err := client.EditDocument(ctx, documentUUID, content)
				if err != nil {
					return err
				}

				// need to retrieve for correspondent information
				document, err := client.DocumentByUUID(ctx, documentUUID)
				if err != nil {
					return err
				}

				previewDocument, err := client.AddDocument(ctx, caseID, document, sirius.TypePreview)
				if err != nil {
					return err
				}

				data.Document = document
				data.PreviewDraft = true
				data.DownloadUUID = previewDocument.UUID

			case "delete":
				err := client.DeleteDocument(ctx, documentUUID)
				if err != nil {
					return err
				}

			case "publish":
				_, err := client.EditDocument(ctx, documentUUID, content)
				if err != nil {
					return err
				}

				// need to retrieve for correspondent information
				document, err := client.DocumentByUUID(ctx, documentUUID)
				if err != nil {
					return err
				}

				_, err = client.AddDocument(ctx, caseID, document, sirius.TypeSave)
				if err != nil {
					return err
				}

				err = client.DeleteDocument(ctx, documentUUID)
				if err != nil {
					return err
				}
			case "saveAndExit":
				_, err := client.EditDocument(ctx, documentUUID, content)
				if err != nil {
					return err
				}
				data.SaveAndExit = true
			}

			if !data.SaveAndExit {
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
					documents, err := client.Documents(ctx.With(groupCtx), caseType, caseID)
					if err != nil {
						return err
					}

					data.Documents = getDraftDocuments(documents)
					return nil
				})

				if err := group.Wait(); err != nil {
					return err
				}

				if documentControls == "delete" || documentControls == "publish" {
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

		return tmpl(w, data)
	}
}

func getDraftDocuments(documents []sirius.Document) []sirius.Document {
	var draftDocuments []sirius.Document
	for _, document := range documents {
		if document.Type == sirius.TypeDraft {
			draftDocuments = append(draftDocuments, document)
		}
	}
	return draftDocuments
}
