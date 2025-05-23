package server

import (
	"github.com/go-playground/form/v4"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/sync/errgroup"
	"net/http"
)

type ResolveObjectionClient interface {
	CaseSummary(sirius.Context, string) (sirius.CaseSummary, error)
	GetObjection(sirius.Context, string) (sirius.Objection, error)
}

type resolveObjectionData struct {
	XSRFToken   string
	Success     bool
	Error       sirius.ValidationError
	CaseUID     string
	ObjectionId string
	Outcome     string
	Objection   sirius.Objection
	LpaUids     []string
}

func ResolveObjection(client ResolveObjectionClient, formTmpl template.Template) Handler {
	if decoder == nil {
		decoder = form.NewDecoder()
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		caseUID := r.PathValue("uid")
		objectionID := r.PathValue("id")

		ctx := getContext(r)

		var obj sirius.Objection
		var err error

		group, groupCtx := errgroup.WithContext(ctx.Context)

		group.Go(func() error {
			obj, err = client.GetObjection(ctx.With(groupCtx), objectionID)
			if err != nil {
				return err
			}
			return nil
		})

		if err := group.Wait(); err != nil {
			return err
		}

		data := resolveObjectionData{
			XSRFToken:   ctx.XSRFToken,
			CaseUID:     caseUID,
			ObjectionId: objectionID,
			Objection:   obj,
			LpaUids:     obj.LpaUids,
		}

		return formTmpl(w, data)
	}
}
