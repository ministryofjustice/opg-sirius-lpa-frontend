package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type ComparingDocumentsClient interface {
	CasesByDonor(ctx sirius.Context, id int) ([]sirius.Case, error)
	GetPersonDocuments(ctx sirius.Context, personID int, caseIDs []string) (sirius.DocumentList, error)
	DocumentByUUID(ctx sirius.Context, uuid string) (sirius.Document, error)
}

type comparingDocumentsData struct {
	XSRFToken         string
	Entity            string
	DocumentList      sirius.DocumentList
	Document          sirius.Document
	DocumentComparing sirius.Document
}

func ComparingDocuments(client ComparingDocumentsClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		docUUIDs := r.URL.Query()["docUid[]"]
		ctx := getContext(r)

		documentData, err := client.DocumentByUUID(ctx, docUUIDs[0])
		if err != nil {
			return err
		}

		documentComparingData, err := client.DocumentByUUID(ctx, docUUIDs[1])
		if err != nil {
			return err
		}

		data := comparingDocumentsData{
			XSRFToken:         ctx.XSRFToken,
			Document:          documentData,
			DocumentComparing: documentComparingData,
		}

		return tmpl(w, data)
	}
}
