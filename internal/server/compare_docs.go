package server

import (
	"fmt"
	"net/http"

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
	View1            *viewingDocumentData
	View2            *viewingDocumentData
}

type viewingDocumentData struct {
	Document sirius.Document
	Pane     int
	BackURL  string
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
		baseURL := fmt.Sprintf("/compare/%d/%s", donorID, caseID)

		data := compareDocsData{
			Pane1: "list",
			Pane2: "list",
			DocListPane1Data: documentPageData{
				XSRFToken:     ctx.XSRFToken,
				DocumentList:  docs,
				SelectedCases: selected,
				Comparing:     true,
			},
			DocListPane2Data: documentPageData{
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
		}

		pane1UUID := r.URL.Query().Get("pane1")
		pane2UUID := r.URL.Query().Get("pane2")

		data.DocListPane1Data.CompareURLs = make(map[string]string)
		data.DocListPane2Data.CompareURLs = make(map[string]string)

		for _, doc := range docs.Documents {
			panel1Url := baseURL + "?pane1=" + doc.UUID
			panel2Url := baseURL + "?pane2=" + doc.UUID
			if pane1UUID != "" {
				panel2Url += "&pane1=" + pane1UUID
			}
			if pane2UUID != "" {
				panel1Url += "&pane2=" + pane2UUID
			}
			data.DocListPane1Data.CompareURLs[doc.UUID] = panel1Url
			data.DocListPane2Data.CompareURLs[doc.UUID] = panel2Url
		}

		if pane1UUID != "" {
			doc, err := client.DocumentByUUID(ctx, pane1UUID)
			if err != nil {
				return err
			}

			backURL := baseURL
			if pane2UUID != "" {
				backURL = baseURL + "?pane2=" + pane2UUID
			}

			data.Pane1 = "doc"
			data.View1 = &viewingDocumentData{
				Document: doc,
				Pane:     1,
				BackURL:  backURL,
			}
		}

		if pane2UUID != "" {
			doc, err := client.DocumentByUUID(ctx, pane2UUID)
			if err != nil {
				return err
			}

			backURL := baseURL
			if pane1UUID != "" {
				backURL = baseURL + "?pane1=" + pane1UUID
			}

			data.Pane2 = "doc"
			data.View2 = &viewingDocumentData{
				Document: doc,
				Pane:     2,
				BackURL:  backURL,
			}
		}

		return tmpl(w, data)
	}
}
