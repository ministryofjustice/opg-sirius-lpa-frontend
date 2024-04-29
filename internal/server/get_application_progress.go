package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"golang.org/x/sync/errgroup"
)

type GetApplicationProgressClient interface {
	CaseSummary(siriusCtx sirius.Context, uid string) (sirius.CaseSummary, error)
	ProgressIndicatorsForDigitalLpa(siriusCtx sirius.Context, uid string) ([]sirius.ProgressIndicator, error)
}

type getApplicationProgressDetails struct {
	CaseSummary        sirius.CaseSummary
	ProgressIndicators []sirius.ProgressIndicator
	FlashMessage       FlashNotification
}

func GetApplicationProgressDetails(client GetApplicationProgressClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		var data getApplicationProgressDetails

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

		data.FlashMessage, _ = GetFlash(w, r)

		return tmpl(w, data)
	}
}
