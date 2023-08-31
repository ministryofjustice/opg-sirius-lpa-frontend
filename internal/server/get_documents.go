package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
)

type GetDocumentsClient interface {
	DigitalLpa(ctx sirius.Context, uid string) (sirius.DigitalLpa, error)
	Documents(ctx sirius.Context, caseType sirius.CaseType, caseId int, docType string) ([]sirius.Document, error)
}

type getDocumentsData struct {
	XSRFToken string

	Lpa       sirius.DigitalLpa
	Documents []sirius.Document
}

func GetDocuments(client GetDocumentsClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		uid := chi.URLParam(r, "uid")
		ctx := getContext(r)
		lpa, err := client.DigitalLpa(ctx, uid)
		if err != nil {
			return err
		}

		documents, err := client.Documents(ctx, "lpa", lpa.ID, sirius.TypeSave)
		if err != nil {
			return err
		}

		data := getDocumentsData{
			XSRFToken: ctx.XSRFToken,
			Documents: documents,
			Lpa:       lpa,
		}

		return tmpl(w, data)
	}
}
