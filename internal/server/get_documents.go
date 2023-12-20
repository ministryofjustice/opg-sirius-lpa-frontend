package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/ministryofjustice/opg-go-common/template"
)

type GetDocumentsClient interface {
	CaseSummary(ctx sirius.Context, uid string) (sirius.CaseSummary, error)
	Documents(ctx sirius.Context, caseType sirius.CaseType, caseId int, docTypes []string, notDocTypes []string) ([]sirius.Document, error)
}

type getDocumentsData struct {
	XSRFToken string

	CaseSummary  sirius.CaseSummary
	Documents    []sirius.Document
	FlashMessage FlashNotification
}

func GetDocuments(client GetDocumentsClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		var err error

		uid := chi.URLParam(r, "uid")
		ctx := getContext(r)

		data := getDocumentsData{
			XSRFToken: ctx.XSRFToken,
		}

		data.CaseSummary, err = client.CaseSummary(ctx, uid)

		if err != nil {
			return err
		}

		data.Documents, err = client.Documents(
			ctx,
			"lpa",
			data.CaseSummary.Lpa.ID,
			[]string{}, []string{sirius.TypeDraft, sirius.TypePreview})

		if err != nil {
			return err
		}

		data.FlashMessage, _ = GetFlash(w, r)

		return tmpl(w, data)
	}
}
