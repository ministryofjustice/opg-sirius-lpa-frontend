package server

import (
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"net/http"
	"strconv"
)

type EditDocumentClient interface {
	Documents(ctx sirius.Context, caseType sirius.CaseType, caseId int) ([]sirius.Document, error)
	Case(ctx sirius.Context, id int) (sirius.Case, error)
	DocumentByUUID(ctx sirius.Context, uuid string) (sirius.Document, error)
}

type editDocumentData struct {
	XSRFToken string
	Success   bool
	Error     sirius.ValidationError
	Case      sirius.Case
	Documents []sirius.Document
	Document  sirius.Document
}

func EditDocument(client EditDocumentClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		data := editDocumentData{
			XSRFToken: ctx.XSRFToken,
		}

		switch r.Method {
		case http.MethodGet:
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

			documents, err := client.Documents(ctx, caseType, caseID)
			if err != nil {
				return err
			}

			data.Case = caseItem
			data.Documents = documents

			defaultDocumentUUID := documents[0].UUID
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

		return tmpl(w, data)
	}
}
