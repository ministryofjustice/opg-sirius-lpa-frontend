package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type ViewDocumentClient interface {
	DocumentByUUID(ctx sirius.Context, uuid string) (sirius.Document, error)
}

func ViewDocument(client ViewDocumentClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if err := r.ParseForm(); err != nil {
			return err
		}

		uuid := r.PathValue("uuid")
		ctx := getContext(r)

		documentData, err := client.DocumentByUUID(ctx, uuid)
		if err != nil {
			return err
		}

		return tmpl(w, documentData)
	}
}
