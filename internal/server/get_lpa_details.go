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
	CaseSummary        sirius.CaseSummary
	DigitalLpa         sirius.DigitalLpa
	LpaStoreDataLegacy map[string]interface{}
}

// TODO move to digital_lpa.go
type LpaStoreData struct {
	Donor LpaStoreDonor `json:"donor"`
}

type LpaStoreDonor struct {
	FirstNames string `json:"firstNames"`
	LastName   string `json:"lastName"`
}

func GetLpaDetails(client GetLpaDetailsClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		uid := chi.URLParam(r, "uid")
		ctx := getContext(r)

		caseSummary, err := client.CaseSummary(ctx, uid)
		if err != nil {
			return err
		}

		var lpaStoreDataLegacy map[string]interface{}
		err = json.Unmarshal(caseSummary.DigitalLpa.LpaStoreData, &lpaStoreDataLegacy)
		if err != nil {
			return err
		}

		data := getLpaDetails{
			CaseSummary:        caseSummary,
			DigitalLpa:         caseSummary.DigitalLpa,
			LpaStoreDataLegacy: lpaStoreDataLegacy,
		}

		return tmpl(w, data)
	}
}
