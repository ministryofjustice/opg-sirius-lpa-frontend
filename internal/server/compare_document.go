package server

import (
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type CompareDocumentClient interface {
	DocumentByUUID(ctx sirius.Context, uuid string) (sirius.Document, error)
	GetPersonDocuments(ctx sirius.Context, personID int, caseIDs []string) (sirius.DocumentList, error)
}

func CompareDocument(client CompareDocumentClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		donorID, err := strToIntOrStatusError(r.PathValue("id"))
		if err != nil {
			return err
		}

		docUUIDs := r.Form["uid[]"]
		ctx := getContext(r)

		documentData, err := client.DocumentByUUID(ctx, docUUIDs[0])
		if err != nil {
			return err
		}

		docs, err := client.GetPersonDocuments(ctx, donorID, []string{strconv.Itoa(documentData.CaseItems[0].ID)})
		if err != nil {
			return err
		}

		selected := documentData.CaseItems

		data := documentPageData{
			XSRFToken:     ctx.XSRFToken,
			DocumentList:  docs,
			Document:      documentData,
			SelectedCases: selected,
			Comparing:     true,
		}

		return tmpl(w, data)
	}
}
