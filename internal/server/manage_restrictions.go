package server

import (
	"fmt"
	"net/http"

	"github.com/go-playground/form/v4"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/sync/errgroup"
)

type ManageRestrictionsClient interface {
	CaseSummary(sirius.Context, string) (sirius.CaseSummary, error)
	ClearTask(sirius.Context, int) error
	UpdateSeveranceStatus(sirius.Context, string, sirius.SeveranceStatusData) error
	EditSeveranceApplication(sirius.Context, string, sirius.SeveranceApplication) error
}

type manageRestrictionsData struct {
	XSRFToken              string
	Error                  sirius.ValidationError
	CaseUID                string
	CaseSummary            sirius.CaseSummary
	SeveranceAction        string
	DonorConsentGiven      string
	CourtOrderedSeverance  string
	CourtOrderDecisionDate sirius.DateString
	CourtOrderReceivedDate sirius.DateString
	FormAction             string
	Success                bool
}

func ManageRestrictions(client ManageRestrictionsClient, tmpl template.Template) Handler {
	if decoder == nil {
		decoder = form.NewDecoder()
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		caseUID := r.PathValue("uid")
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

		fmt.Println(fmt.Sprintf("%#v", cs))

		data := manageRestrictionsData{
			CaseSummary:            cs,
			SeveranceAction:        postFormString(r, "severanceAction"),
			DonorConsentGiven:      postFormString(r, "donorConsentGiven"),
			CourtOrderedSeverance:  postFormString(r, "severanceOrdered"),
			CourtOrderDecisionDate: postFormDateString(r, "courtOrderDecisionMade"),
			CourtOrderReceivedDate: postFormDateString(r, "courtOrderReceived"),
			XSRFToken:              ctx.XSRFToken,
			Error:                  sirius.ValidationError{Field: sirius.FieldErrors{}},
			CaseUID:                caseUID,
		}

		data.FormAction = r.FormValue("action")
		if data.FormAction == "" && data.CaseSummary.DigitalLpa.SiriusData.Application.SeveranceStatus == "REQUIRED" {
			data.FormAction = "donor-consent"

			if data.CaseSummary.DigitalLpa.SiriusData.Application.SeveranceApplication != nil && *data.CaseSummary.DigitalLpa.SiriusData.Application.SeveranceApplication.HasDonorConsented == true {
				data.FormAction = "court-order"
			}
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

					err := client.EditSeveranceApplication(ctx, caseUID, sirius.SeveranceApplication{
						HasDonorConsented: &hasDonorConsented,
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

			if data.FormAction == "court-order" {
				severanceApplication := sirius.SeveranceApplication{}

				if data.CourtOrderDecisionDate != "" {
					severanceApplication.CourtOrderDecisionMade = data.CourtOrderDecisionDate
				}

				if data.CourtOrderReceivedDate != "" {
					severanceApplication.CourtOrderReceived = data.CourtOrderReceivedDate
				}

				switch data.CourtOrderedSeverance {
				case "severance-ordered", "severance-not-ordered":
					isSeveranceOrdered := true
					if data.CourtOrderedSeverance == "severance-not-ordered" {
						isSeveranceOrdered = false
					}

					severanceApplication.SeveranceOrdered = &isSeveranceOrdered
				}

				err := client.EditSeveranceApplication(ctx, caseUID, severanceApplication)
				if handleError(w, &data, err) {
					return err
				}

				return handleSuccess(w, &data, caseUID)
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
