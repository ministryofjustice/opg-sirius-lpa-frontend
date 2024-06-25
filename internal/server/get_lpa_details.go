package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/sync/errgroup"
)

type GetLpaDetailsClient interface {
	CaseSummary(ctx sirius.Context, uid string) (sirius.CaseSummary, error)
	AnomaliesForDigitalLpa(ctx sirius.Context, uid string) ([]sirius.Anomaly, error)
}

type getLpaDetails struct {
	CaseSummary             sirius.CaseSummary
	DigitalLpa              sirius.DigitalLpa
	AnomalyDisplay          sirius.AnomalyDisplay
	ReplacementAttorneys    []sirius.LpaStoreAttorney
	NonReplacementAttorneys []sirius.LpaStoreAttorney
	FlashMessage            FlashNotification
}

func GetLpaDetails(client GetLpaDetailsClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		uid := chi.URLParam(r, "uid")
		ctx := getContext(r)

		var err error
		var data getLpaDetails
		var anomalies []sirius.Anomaly

		group, groupCtx := errgroup.WithContext(ctx.Context)

		group.Go(func() error {
			data.CaseSummary, err = client.CaseSummary(ctx.With(groupCtx), uid)
			if err != nil {
				return err
			}
			return nil
		})

		group.Go(func() error {
			// ignore errors: there may be no LPA in the store yet for a given UID (e.g. if it's still a draft)
			// TODO is there a better way to be selective about ignored errors?
			anomalies, _ = client.AnomaliesForDigitalLpa(ctx.With(groupCtx), uid)
			return nil
		})

		if err := group.Wait(); err != nil {
			return err
		}

		data.AnomalyDisplay = sirius.AnomalyDisplay{}
		if anomalies != nil {
			data.AnomalyDisplay.Group(&data.CaseSummary.DigitalLpa.LpaStoreData, anomalies)
		}

		data.DigitalLpa = data.CaseSummary.DigitalLpa
		data.FlashMessage, _ = GetFlash(w, r)

		var replacementAttorneys []sirius.LpaStoreAttorney
		var nonReplacementAttorneys []sirius.LpaStoreAttorney
		for _, attorney := range data.DigitalLpa.LpaStoreData.Attorneys {
			switch attorney.Status {
			case "replacement":
				replacementAttorneys = append(replacementAttorneys, attorney)
			case "active":
				nonReplacementAttorneys = append(nonReplacementAttorneys, attorney)
			}
		}

		data.ReplacementAttorneys = replacementAttorneys
		data.NonReplacementAttorneys = nonReplacementAttorneys

		return tmpl(w, data)
	}
}
