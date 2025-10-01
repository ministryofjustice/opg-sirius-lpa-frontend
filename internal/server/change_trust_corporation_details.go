package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type ChangeTrustCorporationDetailsClient interface {
	CaseSummary(sirius.Context, string) (sirius.CaseSummary, error)
	RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error)
}

type changeTrustCorporationDetailsData struct {
	XSRFToken string
	Success   bool
	Error     sirius.ValidationError
	CaseUID   string
	TrustCorp sirius.LpaStoreTrustCorporation
}

func ChangeTrustCorporationDetails(client ChangeTrustCorporationDetailsClient, tmpl template.Template) Handler {

	return func(w http.ResponseWriter, r *http.Request) error {
		caseUID := r.PathValue("uid")
		trustCorpUID := r.PathValue("trustCorporationUID")

		ctx := getContext(r)

		cs, err := client.CaseSummary(ctx, caseUID)
		if err != nil {
			return err
		}

		trustCorps := cs.DigitalLpa.LpaStoreData.TrustCorporations

		var trustCorporation sirius.LpaStoreTrustCorporation
		for i, trustCorp := range trustCorps {
			if trustCorp.Uid == trustCorpUID {
				trustCorporation = trustCorps[i]
			}
		}

		data := changeTrustCorporationDetailsData{
			XSRFToken: ctx.XSRFToken,
			CaseUID:   caseUID,
			TrustCorp: trustCorporation,
		}

		return tmpl(w, data)
	}
}
