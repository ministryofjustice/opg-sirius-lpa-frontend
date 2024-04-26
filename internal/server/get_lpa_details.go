package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"golang.org/x/sync/errgroup"
)

type GetLpaDetailsClient interface {
	CaseSummary(siriusCtx sirius.Context, uid string) (sirius.CaseSummary, error)
	ProgressIndicatorsForDigitalLpa(siriusCtx sirius.Context, uid string) ([]sirius.ProgressIndicator, error)
}

type getLpaDetails struct {
	CaseSummary             sirius.CaseSummary
	DigitalLpa              sirius.DigitalLpa
	ReplacementAttorneys    []sirius.LpaStoreAttorney
	NonReplacementAttorneys []sirius.LpaStoreAttorney
	ProgressIndicators      []sirius.ProgressIndicator
	FlashMessage            FlashNotification
}

func GetLpaDetails(client GetLpaDetailsClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		var data getLpaDetails

		uid := chi.URLParam(r, "uid")
		ctx := getContext(r)

		group, groupCtx := errgroup.WithContext(ctx.Context)

		group.Go(func() error {
			cs, err := client.CaseSummary(ctx.With(groupCtx), uid)
			if err != nil {
				return err
			}
			data.CaseSummary = cs
			return nil
		})

		group.Go(func() error {
			inds, err := client.ProgressIndicatorsForDigitalLpa(ctx.With(groupCtx), uid)
			if err != nil {
				return err
			}
			data.ProgressIndicators = inds
			return nil
		})

		if err := group.Wait(); err != nil {
			return err
		}

		// to prevent lots of changes to template structure
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
