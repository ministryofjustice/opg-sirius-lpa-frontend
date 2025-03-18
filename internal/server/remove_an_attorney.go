package server

import (
	"fmt"
	"github.com/go-playground/form/v4"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type RemoveAnAttorneyClient interface {
	CaseSummary(sirius.Context, string) (sirius.CaseSummary, error)
	ChangeAttorneyStatus(sirius.Context, string, []sirius.AttorneyUpdatedStatus) error
}

type formRemoveAttorney struct {
	RemovedAttorneyUid  string   `form:"removedAttorney"`
	EnabledAttorneyUids []string `form:"enabledAttorney"`
	SkipEnableAttorney  string   `form:"skipEnableAttorney"`
}

type SelectedAttorneyDetails struct {
	SelectedAttorneyName string
	SelectedAttorneyDob  string
}

type removeAnAttorneyData struct {
	CaseSummary             sirius.CaseSummary
	ActiveAttorneys         []sirius.LpaStoreAttorney
	InactiveAttorneys       []sirius.LpaStoreAttorney
	Form                    formRemoveAttorney
	RemovedAttorneysDetails SelectedAttorneyDetails
	EnabledAttorneysDetails []SelectedAttorneyDetails
	Success                 bool
	Error                   sirius.ValidationError
	XSRFToken               string
}

func RemoveAnAttorney(client RemoveAnAttorneyClient, removeTmpl template.Template, confirmTmpl template.Template) Handler {
	if decoder == nil {
		decoder = form.NewDecoder()
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		uid := chi.URLParam(r, "uid")
		ctx := getContext(r)

		caseSummary, err := client.CaseSummary(ctx, uid)

		if err != nil {
			return err
		}

		data := removeAnAttorneyData{
			CaseSummary: caseSummary,
			XSRFToken:   ctx.XSRFToken,
			Error:       sirius.ValidationError{Field: sirius.FieldErrors{}},
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

		for _, attorney := range lpa.LpaStoreData.Attorneys {
			if attorney.Status == shared.RemovedAttorneyStatus.String() || attorney.Status == shared.ActiveAttorneyStatus.String() {
				continue
			}

			data.InactiveAttorneys = append(data.InactiveAttorneys, attorney)
		}

		if r.Method == http.MethodPost {

			err = decoder.Decode(&data.Form, r.PostForm)
			if err != nil {
				return err
			}

			if data.Form.RemovedAttorneyUid == "" {
				data.Error.Field["removeAttorney"] = map[string]string{
					"reason": "Please select an attorney for removal",
				}
			}

			if len(data.Form.EnabledAttorneyUids) > 0 && postFormCheckboxChecked(r, "skipEnableAttorney", "yes") {
				data.Error.Field["enableAttorney"] = map[string]string{
					"reason": "Please do not select both a replacement attorney and the option to skip",
				}
			}

			if len(data.Form.EnabledAttorneyUids) == 0 && !postFormCheckboxChecked(r, "skipEnableAttorney", "yes") {
				data.Error.Field["enableAttorney"] = map[string]string{
					"reason": "Please select either the attorneys that can be enabled or skip the replacement of the attorneys",
				}
			}

			if !data.Error.Any() {
				if !postFormKeySet(r, "confirmRemoval") {
					for _, att := range data.ActiveAttorneys {
						if att.Uid == data.Form.RemovedAttorneyUid {
							data.RemovedAttorneysDetails = SelectedAttorneyDetails{
								SelectedAttorneyName: att.FirstNames + " " + att.LastName,
								SelectedAttorneyDob:  att.DateOfBirth,
							}
						}
					}

					if len(data.Form.EnabledAttorneyUids) > 0 {
						for _, att := range data.InactiveAttorneys {
							for _, enabledAttUid := range data.Form.EnabledAttorneyUids {
								if att.Uid == enabledAttUid {
									data.EnabledAttorneysDetails = append(data.EnabledAttorneysDetails, SelectedAttorneyDetails{
										SelectedAttorneyName: att.FirstNames + " " + att.LastName,
										SelectedAttorneyDob:  att.DateOfBirth,
									})
									break
								}
							}
						}
					}

					return confirmTmpl(w, data)
				} else {
					var attorneyUpdatedStatus []sirius.AttorneyUpdatedStatus

					for _, att := range data.ActiveAttorneys {
						if att.Uid == data.Form.RemovedAttorneyUid {
							attorneyUpdatedStatus = append(attorneyUpdatedStatus, sirius.AttorneyUpdatedStatus{
								UID:    att.Uid,
								Status: shared.RemovedAttorneyStatus.String(),
							})
						}
					}

					if len(data.Form.EnabledAttorneyUids) > 0 {
						for _, att := range data.InactiveAttorneys {
							for _, enabledAttUid := range data.Form.EnabledAttorneyUids {
								if att.Uid == enabledAttUid {
									attorneyUpdatedStatus = append(attorneyUpdatedStatus, sirius.AttorneyUpdatedStatus{
										UID:    att.Uid,
										Status: shared.ActiveAttorneyStatus.String(),
									})
								}
							}
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

						SetFlash(w, FlashNotification{Title: "Update saved"})
						return RedirectError(fmt.Sprintf("/lpa/%s", uid))
					}
				}
			}
		}

		return removeTmpl(w, data)
	}
}
