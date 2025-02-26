package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type RemoveAnAttorneyClient interface {
	CaseSummary(sirius.Context, string) (sirius.CaseSummary, error)
}

type removeAnAttorneyData struct {
	CaseSummary sirius.CaseSummary

	ActiveAttorneys  []sirius.LpaStoreAttorney
	SelectedAttorney string
	Error            sirius.ValidationError
	XSRFToken        string
}

func RemoveAnAttorney(client RemoveAnAttorneyClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		uid := chi.URLParam(r, "uid")
		ctx := getContext(r)

		caseSummary, err := client.CaseSummary(ctx, uid)

		if err != nil {
			return err
		}

		data := removeAnAttorneyData{
			SelectedAttorney: postFormString(r, "selectedAttorney"),
			CaseSummary:      caseSummary,
			XSRFToken:        ctx.XSRFToken,
			Error:            sirius.ValidationError{Field: sirius.FieldErrors{}},
		}

		lpa := data.CaseSummary.DigitalLpa

		for _, attorney := range lpa.LpaStoreData.Attorneys {
			if (attorney.Status == shared.RemovedAttorneyStatus.String()) ||
				(attorney.AppointmentType == shared.ReplacementAppointmentType.String() &&
					attorney.Status == shared.InactiveAttorneyStatus.String()) {
				continue
			}

			data.ActiveAttorneys = append(data.ActiveAttorneys, attorney)
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
