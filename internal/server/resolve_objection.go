package server

import (
	"fmt"
	"github.com/go-playground/form/v4"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/sync/errgroup"
	"net/http"
)

type ResolveObjectionClient interface {
	CaseSummary(sirius.Context, string) (sirius.CaseSummary, error)
	GetObjection(sirius.Context, string) (sirius.Objection, error)
	ResolveObjection(sirius.Context, string, string, sirius.ResolutionRequest) error
}

type resolveObjectionData struct {
	XSRFToken   string
	Success     bool
	Error       sirius.ValidationError
	CaseUID     string
	ObjectionId string
	Resolution  string
	Objection   sirius.Objection
	LpaUids     []string
	Form        formResolveObjection
}

type formResolveObjection struct {
	Resolution string `form:"resolution"`
	Notes      string `form:"notes"`
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

		if r.Method == http.MethodPost {
			err := decoder.Decode(&data.Form, r.PostForm)
			if err != nil {
				return err
			}

			resolution := sirius.ResolutionRequest{
				Resolution: data.Form.Resolution,
				Notes:      data.Form.Notes,
			}

			err = client.ResolveObjection(ctx, objectionID, caseUID, resolution)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
			} else if err != nil {
				return err
			} else {
				data.Success = true

				return RedirectError(fmt.Sprintf("/lpa/%s", caseUID))
			}
		}

		return formTmpl(w, data)
	}
}
