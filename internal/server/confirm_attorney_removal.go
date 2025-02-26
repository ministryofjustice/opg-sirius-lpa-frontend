package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type ConfirmAttorneyRemovalClient interface {
	CaseSummary(sirius.Context, string) (sirius.CaseSummary, error)
}

type confirmAttorneyRemovalData struct {
	CaseSummary sirius.CaseSummary

	Attorney         sirius.LpaStoreAttorney
	SelectedAttorney string
	Error            sirius.ValidationError
	XSRFToken        string
}

func ConfirmAttorneyRemoval(client ConfirmAttorneyRemovalClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		uid := chi.URLParam(r, "uid")
		attorneyUID := chi.URLParam(r, "attorneyUID")
		ctx := getContext(r)

		caseSummary, err := client.CaseSummary(ctx, uid)

		if err != nil {
			return err
		}

		data := confirmAttorneyRemovalData{
			SelectedAttorney: postFormString(r, "selectedAttorney"),
			CaseSummary:      caseSummary,
			XSRFToken:        ctx.XSRFToken,
			Error:            sirius.ValidationError{Field: sirius.FieldErrors{}},
		}

		attorneys := caseSummary.DigitalLpa.LpaStoreData.Attorneys

		for i, att := range attorneys {
			if att.Uid == attorneyUID {
				data.Attorney = attorneys[i]
			}
		}

		if r.Method == http.MethodPost {
			if data.SelectedAttorney == "" {
				w.WriteHeader(http.StatusBadRequest)

				data.Error.Field["selectAttorney"] = map[string]string{
					"reason": "Please select an attorney for removal.",
				}
			}

			if !data.Error.Any() {
				return RedirectError(fmt.Sprintf("/lpa/%s/confirm-attorney-removal/%s", caseSummary.DigitalLpa.UID, data.SelectedAttorney))
			}
		}

		return tmpl(w, data)
	}
}
