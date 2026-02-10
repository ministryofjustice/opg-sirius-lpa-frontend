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
	DocListData documentPageData
	Pane1       string
	Pane2       string
	Document1   sirius.Document
	Document2   sirius.Document
	View1       *viewingDocumentData
	View2       *viewingDocumentData
	docUUIDs    []string
}

type viewingDocumentData struct {
	Document sirius.Document
	Pane     int
	BackURL  string
	DocUUIDs []string
}

func compareURL(donorID int, caseID string, uuids []string) string {
	base := fmt.Sprintf("/compare/%d/%s", donorID, caseID)

	if len(uuids) == 0 {
		return base
	}

	q := url.Values{}
	for _, u := range uuids {
		q.Add("docUid[]", u)
	}

	return base + "?" + q.Encode()
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
			DocListData: documentPageData{
				XSRFToken:     ctx.XSRFToken,
				DocumentList:  docs,
				SelectedCases: selected,
				Comparing:     true,
			},
			Pane1: "list",
			Pane2: "list",
		}

		docUUIDs := r.URL.Query()["docUid[]"]
		data.docUUIDs = docUUIDs

		data.DocListData.CompareBaseURL = fmt.Sprintf(
			"/compare/%d/%s",
			donorID,
			caseID,
		)
		data.DocListData.DocUUIDs = docUUIDs

		switch len(docUUIDs) {
		case 0:
		//	remains the same (list and list already set)

		case 1:
			doc, err := client.DocumentByUUID(ctx, docUUIDs[0])
			if err != nil {
				return err
			}

			data.Pane1 = "doc"
			data.View1 = &viewingDocumentData{
				Document: doc,
				Pane:     1,
				BackURL:  compareURL(donorID, caseID, nil),
				DocUUIDs: docUUIDs,
			}

		default:
			doc1, err := client.DocumentByUUID(ctx, docUUIDs[0])
			if err != nil {
				return err
			}

			doc2, err := client.DocumentByUUID(ctx, docUUIDs[1])
			if err != nil {
				return err
			}

			data.Pane1 = "doc"
			data.Pane2 = "doc"
			data.View1 = &viewingDocumentData{
				Document: doc1,
				Pane:     1,
				BackURL:  compareURL(donorID, caseID, []string{docUUIDs[0]}),
				DocUUIDs: docUUIDs,
			}
			data.View2 = &viewingDocumentData{
				Document: doc2,
				Pane:     2,
				BackURL:  compareURL(donorID, caseID, []string{docUUIDs[1]}),
				DocUUIDs: docUUIDs,
			}
		}

		return tmpl(w, data)
	}
}
