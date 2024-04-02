package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
)

type GetLpaDetailsClient interface {
	CaseSummary(ctx sirius.Context, uid string) (sirius.CaseSummary, error)
}

type getLpaDetails struct {
	CaseSummary  sirius.CaseSummary
	DigitalLpa   sirius.DigitalLpa
	FlashMessage FlashNotification
}

func GetLpaDetails(client GetLpaDetailsClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		uid := chi.URLParam(r, "uid")
		ctx := getContext(r)

		caseSummary, err := client.CaseSummary(ctx, uid)
		if err != nil {
			return err
		}

		data := getLpaDetails{
			CaseSummary: caseSummary,
			DigitalLpa:  caseSummary.DigitalLpa,
		}

		data.FlashMessage, _ = GetFlash(w, r)

		return tmpl(w, data)
	}
}
