package server

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type CompareDocsClient interface {
	DocumentByUUID(ctx sirius.Context, uuid string) (sirius.Document, error)
	GetPersonDocuments(ctx sirius.Context, personID int, caseIDs []string) (sirius.DocumentList, error)
}

type compareDocsData struct {
	DocListPane1Data documentPageData
	DocListPane2Data documentPageData
	Pane1            string
	Pane2            string
	Document1        sirius.Document
	Document2        sirius.Document
	View1            *viewingDocumentData
	View2            *viewingDocumentData
	docUUIDs         []string
}

type viewingDocumentData struct {
	Document sirius.Document
	Pane     int
	BackURL  string
	DocUUIDs []string
}

func compareURL(donorID int, caseID string, panes map[string]string) string {
	baseUrl := fmt.Sprintf("/compare/%d/%s", donorID, caseID)

	query := url.Values{}
	for key, value := range panes {
		if value != "" {
			query.Set(key, value)
		}
	}

	if len(query) == 0 {
		return baseUrl
	}

	return baseUrl + "?" + query.Encode()
}

func CompareDocs(client CompareDocsClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		donorID, err := strToIntOrStatusError(r.PathValue("id"))
		if err != nil {
			return err
		}

		caseID := r.PathValue("caseId")
		ctx := getContext(r)

		docs, err := client.GetPersonDocuments(ctx, donorID, []string{caseID})
		if err != nil {
			return err
		}

		selected := docs.Documents[0].CaseItems

		data := compareDocsData{
			DocListPane1Data: documentPageData{
				XSRFToken:     ctx.XSRFToken,
				DocumentList:  docs,
				SelectedCases: selected,
				Comparing:     true,
				TargetPane:    1,
			},
			DocListPane2Data: documentPageData{
				XSRFToken:     ctx.XSRFToken,
				DocumentList:  docs,
				SelectedCases: selected,
				Comparing:     true,
				TargetPane:    2,
			},
			Pane1: "list",
			Pane2: "list",
		}

		pane1UUID := r.URL.Query().Get("pane1")
		pane2UUID := r.URL.Query().Get("pane2")

		data.DocListPane1Data.CompareBaseURL = fmt.Sprintf(
			"/compare/%d/%s",
			donorID,
			caseID,
		)
		data.DocListPane2Data.CompareBaseURL = fmt.Sprintf(
			"/compare/%d/%s",
			donorID,
			caseID,
		)

		data.DocListPane1Data.CompareURLs = make(map[string]string)
		for _, doc := range docs.Documents {
			data.DocListPane1Data.CompareURLs[doc.UUID] = compareURL(donorID, caseID, map[string]string{
				"pane1": doc.UUID,
				"pane2": pane2UUID,
			})
		}

		data.DocListPane2Data.CompareURLs = make(map[string]string)
		for _, doc := range docs.Documents {
			data.DocListPane2Data.CompareURLs[doc.UUID] = compareURL(donorID, caseID, map[string]string{
				"pane1": pane1UUID,
				"pane2": doc.UUID,
			})
		}

		if pane1UUID != "" && pane2UUID == "" {
			doc, err := client.DocumentByUUID(ctx, pane1UUID)
			if err != nil {
				return err
			}

			data.Pane1 = "doc"
			data.View1 = &viewingDocumentData{
				Document: doc,
				Pane:     1,
				BackURL:  compareURL(donorID, caseID, map[string]string{}),
			}
		}

		if pane1UUID == "" && pane2UUID != "" {
			doc, err := client.DocumentByUUID(ctx, pane2UUID)
			if err != nil {
				return err
			}

			data.Pane2 = "doc"
			data.View2 = &viewingDocumentData{
				Document: doc,
				Pane:     2,
				BackURL:  compareURL(donorID, caseID, map[string]string{}),
			}
		}

		if pane1UUID != "" && pane2UUID != "" {
			doc1, err := client.DocumentByUUID(ctx, pane1UUID)
			if err != nil {
				return err
			}
			doc2, err := client.DocumentByUUID(ctx, pane2UUID)
			if err != nil {
				return err
			}

			data.Pane1 = "doc"
			data.View1 = &viewingDocumentData{
				Document: doc1,
				Pane:     1,
				BackURL: compareURL(donorID, caseID, map[string]string{
					"pane2": pane2UUID,
				}),
			}

			data.Pane2 = "doc"
			data.View2 = &viewingDocumentData{
				Document: doc2,
				Pane:     2,
				BackURL: compareURL(donorID, caseID, map[string]string{
					"pane1": pane1UUID,
				}),
			}
		}

		return tmpl(w, data)
	}
}
