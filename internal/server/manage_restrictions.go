package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/form/v4"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/sync/errgroup"
	"net/http"
)

type ManageRestrictionsClient interface {
	CaseSummary(sirius.Context, string) (sirius.CaseSummary, error)
}

type manageRestrictionsData struct {
	XSRFToken   string
	Success     bool
	Error       sirius.ValidationError
	CaseUID     string
	CaseSummary sirius.CaseSummary
}

func ManageRestrictions(client ManageRestrictionsClient, tmpl template.Template) Handler {
	if decoder == nil {
		decoder = form.NewDecoder()
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		caseUID := chi.URLParam(r, "uid")
		ctx := getContext(r)

		var cs sirius.CaseSummary
		var err error

		group, groupCtx := errgroup.WithContext(ctx.Context)

		group.Go(func() error {
			cs, err = client.CaseSummary(ctx.With(groupCtx), caseUID)
			if err != nil {
				return err
			}
			return nil
		})

		if err := group.Wait(); err != nil {
			return err
		}

		data := manageRestrictionsData{
			CaseSummary: cs,
		}

		return tmpl(w, data)
	}
}
