package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type ManageAttorneysClient interface {
	CaseSummary(sirius.Context, string, bool) (sirius.CaseSummary, error)
}

type manageAttorneysData struct {
	CaseSummary sirius.CaseSummary

	AttorneyAction string
	Error          sirius.ValidationError
	XSRFToken      string
}

func ManageAttorneys(client ManageAttorneysClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		uid := chi.URLParam(r, "uid")
		ctx := getContext(r)

		caseSummary, err := client.CaseSummary(ctx, uid, false)

		if err != nil {
			return err
		}

		data := manageAttorneysData{
			AttorneyAction: postFormString(r, "attorneyAction"),
			CaseSummary:    caseSummary,
			XSRFToken:      ctx.XSRFToken,
			Error:          sirius.ValidationError{Field: sirius.FieldErrors{}},
		}

		if r.Method == http.MethodPost {
			var redirectUrl string

			switch data.AttorneyAction {
			case "remove-an-attorney":
				redirectUrl = fmt.Sprintf("/lpa/%s/remove-an-attorney", caseSummary.DigitalLpa.UID)

			case "enable-replacement-attorney":
				redirectUrl = fmt.Sprintf("/lpa/%s/enable-replacement-attorney", caseSummary.DigitalLpa.UID)

			default:
				w.WriteHeader(http.StatusBadRequest)

				data.Error.Field["attorneyAction"] = map[string]string{
					"reason": "Please select an option to manage attorneys.",
				}
			}

			if !data.Error.Any() && redirectUrl != "" {
				return RedirectError(redirectUrl)
			}
		}

		return tmpl(w, data)
	}
}
