package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"net/http"
)

type LpaClient interface {
	CaseSummary(ctx sirius.Context, uid string) (sirius.CaseSummary, error)
}

type lpaData struct {
	DigitalLpa   sirius.DigitalLpa
	CaseSummary  sirius.CaseSummary
	FlashMessage FlashNotification
}

func Lpa(client LpaClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		uid := chi.URLParam(r, "uid")
		ctx := getContext(r)

		caseSummary, err := client.CaseSummary(ctx, uid)

		if err != nil {
			return err
		}

		data := lpaData{
			DigitalLpa:  caseSummary.DigitalLpa,
			CaseSummary: caseSummary,
		}

		data.FlashMessage, _ = GetFlash(w, r)

		return tmpl(w, data)
	}
}
