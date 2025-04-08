package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/form/v4"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/sync/errgroup"
	"net/http"
)

type ManageRestrictionsClient interface {
	CaseSummary(sirius.Context, string) (sirius.CaseSummary, error)
	ClearTask(sirius.Context, int) error
	UpdateSeveranceStatus(sirius.Context, string, sirius.SeveranceStatusData) error
	EditSeveranceApplication(sirius.Context, string, sirius.SeveranceApplicationDetails) error
}

type manageRestrictionsData struct {
	XSRFToken         string
	Error             sirius.ValidationError
	CaseUID           string
	CaseSummary       sirius.CaseSummary
	SeveranceAction   string
	DonorConsentGiven string
	FormAction        string
	Success           bool
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
			CaseSummary:       cs,
			SeveranceAction:   postFormString(r, "severanceAction"),
			DonorConsentGiven: postFormString(r, "donorConsentGiven"),
			XSRFToken:         ctx.XSRFToken,
			Error:             sirius.ValidationError{Field: sirius.FieldErrors{}},
			CaseUID:           caseUID,
		}

		data.FormAction = r.FormValue("action")
		if data.FormAction == "" && data.CaseSummary.DigitalLpa.SiriusData.Application.SeveranceStatus == "REQUIRED" {
			data.FormAction = "donor-consent"
		}

		if r.Method == http.MethodPost {
			taskList := data.CaseSummary.TaskList
			var taskID int
			for _, task := range taskList {
				if task.Name == "Review restrictions and conditions" && task.Status != "Completed" {
					taskID = task.ID
				}
			}

			if data.FormAction == "" || data.FormAction == "change-severance-required" {
				switch data.SeveranceAction {
				case "severance-application-required", "severance-application-not-required":
					severanceStatus := "REQUIRED"
					if data.SeveranceAction == "severance-application-not-required" {
						severanceStatus = "NOT_REQUIRED"
					}
					err = client.UpdateSeveranceStatus(ctx, caseUID, sirius.SeveranceStatusData{
						SeveranceStatus: severanceStatus,
					})
					if handleError(w, &data, err) {
						return err
					}

					if data.SeveranceAction == "severance-application-not-required" {
						err := client.ClearTask(ctx, taskID)
						if handleError(w, &data, err) {
							return err
						}
					}

					return handleSuccess(w, &data, caseUID)

				default:
					w.WriteHeader(http.StatusBadRequest)
					data.Error.Field["severanceAction"] = map[string]string{
						"reason": "Please select an option",
					}
				}
			}

			if data.FormAction == "donor-consent" {
				switch data.DonorConsentGiven {
				case "donor-consent-given", "donor-consent-not-given":
					hasDonorConsented := true
					if data.DonorConsentGiven == "donor-consent-not-given" {
						hasDonorConsented = false
					}

					err := client.EditSeveranceApplication(ctx, caseUID, sirius.SeveranceApplicationDetails{
						HasDonorConsented: hasDonorConsented,
					})
					if handleError(w, &data, err) {
						return err
					}

					return handleSuccess(w, &data, caseUID)

				default:
					w.WriteHeader(http.StatusBadRequest)
					data.Error.Field["donorConsentAction"] = map[string]string{
						"reason": "Please select an option",
					}
				}
			}
		}
		return tmpl(w, data)
	}
}

func handleError(w http.ResponseWriter, data *manageRestrictionsData, err error) bool {
	if err == nil {
		return false
	}
	if ve, ok := err.(sirius.ValidationError); ok {
		w.WriteHeader(http.StatusBadRequest)
		data.Error = ve
		return true
	}
	return true
}

func handleSuccess(w http.ResponseWriter, data *manageRestrictionsData, caseUID string) RedirectError {
	data.Success = true
	SetFlash(w, FlashNotification{
		Title: "Update saved",
	})
	return RedirectError(fmt.Sprintf("/lpa/%s/lpa-details", caseUID))
}
