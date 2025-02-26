package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"net/http"
)

type RemoveAnAttorneyClient interface {
	CaseSummary(sirius.Context, string) (sirius.CaseSummary, error)
}

type removeAnAttorneyData struct {
	CaseSummary sirius.CaseSummary

	ActiveAttorneys      []sirius.LpaStoreAttorney
	SelectedAttorneyUid  string
	SelectedAttorneyName string
	SelectedAttorneyDob  string
	Error                sirius.ValidationError
	XSRFToken            string
}

func RemoveAnAttorney(client RemoveAnAttorneyClient, removeTmpl template.Template, confirmTmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		uid := chi.URLParam(r, "uid")
		ctx := getContext(r)

		caseSummary, err := client.CaseSummary(ctx, uid)

		if err != nil {
			return err
		}

		data := removeAnAttorneyData{
			SelectedAttorneyUid: postFormString(r, "selectedAttorney"),
			CaseSummary:         caseSummary,
			XSRFToken:           ctx.XSRFToken,
			Error:               sirius.ValidationError{Field: sirius.FieldErrors{}},
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
			if data.SelectedAttorneyUid == "" {
				w.WriteHeader(http.StatusBadRequest)

				data.Error.Field["selectAttorney"] = map[string]string{
					"reason": "Please select an attorney for removal.",
				}
			}

			if !data.Error.Any() {
				if !postFormKeySet(r, "confirmRemoval") {
					for _, att := range data.ActiveAttorneys {
						if att.Uid == data.SelectedAttorneyUid {
							data.SelectedAttorneyName = att.FirstNames + " " + att.LastName
							data.SelectedAttorneyDob = att.DateOfBirth
						}
					}

					return confirmTmpl(w, data)
				} else {
					return RedirectError(fmt.Sprintf("/lpa/%s", caseSummary.DigitalLpa.UID))
				}
			}
		}

		return removeTmpl(w, data)
	}
}
