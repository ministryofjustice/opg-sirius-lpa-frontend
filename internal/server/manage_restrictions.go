package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"net/http"
)

type ManageRestrictionsClient interface {
	CaseSummary(sirius.Context, string) (sirius.CaseSummary, error)
}

type manageRestrictionsData struct {
	XSRFToken       string
	Error           sirius.ValidationError
	CaseUID         string
	CaseSummary     sirius.CaseSummary
	SeveranceAction string
}

func ManageRestrictions(client ManageRestrictionsClient, tmpl template.Template) Handler {
	//if decoder == nil {
	//	decoder = form.NewDecoder()
	//}

	return func(w http.ResponseWriter, r *http.Request) error {
		caseUID := chi.URLParam(r, "uid")
		ctx := getContext(r)

		cs, err := client.CaseSummary(ctx, caseUID)
		if err != nil {
			return err
		}

		//var cs sirius.CaseSummary
		//var err error
		//
		//group, groupCtx := errgroup.WithContext(ctx.Context)
		//
		//group.Go(func() error {
		//	cs, err = client.CaseSummary(ctx.With(groupCtx), caseUID)
		//	if err != nil {
		//		return err
		//	}
		//	return nil
		//})
		//
		//if err := group.Wait(); err != nil {
		//	return err
		//}

		data := manageRestrictionsData{
			CaseSummary:     cs,
			SeveranceAction: postFormString(r, "severanceAction"),
			XSRFToken:       ctx.XSRFToken,
			Error:           sirius.ValidationError{Field: sirius.FieldErrors{}},
			CaseUID:         caseUID,
		}

		if r.Method == http.MethodPost {
			var redirectUrl string

			switch data.SeveranceAction {
			case "severance-application-not-required":
				redirectUrl = fmt.Sprintf("/lpa/%s/remove-an-attorney", caseUID)

			case "severance-application-required":
				redirectUrl = fmt.Sprintf("/lpa/%s/enable-replacement-attorney", caseUID)

			default:
				w.WriteHeader(http.StatusBadRequest)

				data.Error.Field["severanceAction"] = map[string]string{
					"reason": "Please select an option",
				}
			}

			if !data.Error.Any() && redirectUrl != "" {
				return RedirectError(redirectUrl)
			}
		}

		return tmpl(w, data)
	}
}
