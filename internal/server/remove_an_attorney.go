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
	ChangeAttorneyStatus(sirius.Context, string, []sirius.AttorneyUpdatedStatus) error
}

type removeAnAttorneyData struct {
	CaseSummary sirius.CaseSummary

	ActiveAttorneys      []sirius.LpaStoreAttorney
	SelectedAttorneyUid  string
	SelectedAttorneyName string
	SelectedAttorneyDob  string
	Success              bool
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
					var attorneyUpdatedStatus []sirius.AttorneyUpdatedStatus

					for _, att := range data.ActiveAttorneys {
						if att.Uid == data.SelectedAttorneyUid {
							attorneyUpdatedStatus = append(attorneyUpdatedStatus, sirius.AttorneyUpdatedStatus{
								UID:    att.Uid,
								Status: shared.RemovedAttorneyStatus.String(),
							})
						}
					}

					err = client.ChangeAttorneyStatus(ctx, uid, attorneyUpdatedStatus)

					if ve, ok := err.(sirius.ValidationError); ok {
						w.WriteHeader(http.StatusBadRequest)
						data.Error = ve
					} else if err != nil {
						return err
					} else {
						data.Success = true

						SetFlash(w, FlashNotification{Title: "Attorney statuses updated"})
						return RedirectError(fmt.Sprintf("/lpa/%s", uid))
					}
				}
			}
		}

		return removeTmpl(w, data)
	}
}
