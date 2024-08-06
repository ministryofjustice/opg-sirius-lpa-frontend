package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"

	"github.com/ministryofjustice/opg-go-common/template"
)

type GetApplicationProgressClient interface {
	CaseSummary(siriusCtx sirius.Context, uid string) (sirius.CaseSummary, error)
	ProgressIndicatorsForDigitalLpa(siriusCtx sirius.Context, uid string) ([]sirius.ProgressIndicator, error)
}

type IndicatorView struct {
	UID                        string
	CertificateProviderName    string
	CertificateProviderChannel string
	ApplicationSource          string
	sirius.ProgressIndicator
}

type getApplicationProgressDetails struct {
	CaseSummary        sirius.CaseSummary
	ProgressIndicators []IndicatorView
	FlashMessage       FlashNotification
}

func GetApplicationProgressDetails(client GetApplicationProgressClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		var data getApplicationProgressDetails

		uid := chi.URLParam(r, "uid")
		ctx := getContext(r)

		var cpName string

		cs, err := client.CaseSummary(ctx, uid)
		if err != nil {
			return err
		}
		data.CaseSummary = cs

		cpName = cs.DigitalLpa.LpaStoreData.CertificateProvider.FirstNames + " " + cs.DigitalLpa.LpaStoreData.CertificateProvider.LastName

		inds, err := client.ProgressIndicatorsForDigitalLpa(ctx, uid)
		if err != nil {
			return err
		}
		for _, v := range inds {
			data.ProgressIndicators = append(data.ProgressIndicators, IndicatorView{
				UID:                        uid,
				CertificateProviderName:    cpName,
				ProgressIndicator:          v,
				CertificateProviderChannel: cs.DigitalLpa.LpaStoreData.CertificateProvider.Channel,
				ApplicationSource:          cs.DigitalLpa.SiriusData.Application.Source,
			})
		}

		data.FlashMessage, _ = GetFlash(w, r)

		return tmpl(w, data)
	}
}
