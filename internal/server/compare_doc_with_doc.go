package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type CompareDocWithDocClient interface {
	DocumentByUUID(ctx sirius.Context, uuid string) (sirius.Document, error)
}

type comparingDocumentsData struct {
	XSRFToken         string
	Entity            string
	Document          sirius.Document
	DocumentComparing sirius.Document
}

func CompareDocWithDoc(client CompareDocWithDocClient, tmpl template.Template) Handler {
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
