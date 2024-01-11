package server

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
)

type GetLpaDetailsClient interface {
	CaseSummary(ctx sirius.Context, uid string) (sirius.CaseSummary, error)
}

type getLpaDetails struct {
	XSRFToken string

	LpaStoreData map[string]interface{}
	CaseSummary  sirius.CaseSummary
}

func GetLpaDetails(client GetLpaDetailsClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		uid := chi.URLParam(r, "uid")
		ctx := getContext(r)

		caseSummary, err := client.CaseSummary(ctx, uid)

		if err != nil {
			return err
		}

		var lpaStoreData map[string]interface{}
		_ = json.Unmarshal(caseSummary.DigitalLpa.LpaStoreData, &lpaStoreData)

		data := getLpaDetails{
			LpaStoreData: lpaStoreData,
			CaseSummary:  caseSummary,
		}

		return tmpl(w, data)
	}
}
