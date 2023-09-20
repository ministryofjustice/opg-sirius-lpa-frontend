package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"

	"github.com/ministryofjustice/opg-go-common/template"
)

type GetDocumentsClient interface {
	DigitalLpa(ctx sirius.Context, uid string) (sirius.DigitalLpa, error)
	Documents(ctx sirius.Context, caseType sirius.CaseType, caseId int, docTypes []string, notDocTypes []string) ([]sirius.Document, error)
}

type getDocumentsData struct {
	XSRFToken string

	Lpa          sirius.DigitalLpa
	Documents    []sirius.Document
	FlashMessage FlashNotification
}

func GetDocuments(client GetDocumentsClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		uid := chi.URLParam(r, "uid")
		ctx := getContext(r)
		lpa, err := client.DigitalLpa(ctx, uid)
		if err != nil {
			return err
		}

		documents, err := client.Documents(ctx, "lpa", lpa.ID, []string{}, []string{sirius.TypeDraft, sirius.TypePreview})
		if err != nil {
			return err
		}

		data := getDocumentsData{
			XSRFToken: ctx.XSRFToken,
			Documents: documents,
			Lpa:       lpa,
		}

		data.FlashMessage, _ = GetFlash(w, r)

		return tmpl(w, data)
	}
}
