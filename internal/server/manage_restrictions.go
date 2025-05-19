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

type CourtOrderRestrictionDetails struct {
	SelectedCourtOrderDecisionDate sirius.DateString
	SelectedCourtOrderReceivedDate sirius.DateString
	SelectedSeveranceAction        string
	SelectedSeveranceType          string
	SelectedSeveranceActionDetail  string
	RemovedWords                   string
	ChangedRestrictions            string
	SelectedFormAction             string
}

type manageRestrictionsData struct {
	XSRFToken                 string
	Error                     sirius.ValidationError
	CaseUID                   string
	CaseSummary               sirius.CaseSummary
	SeveranceAction           string
	DonorConsentGiven         string
	SeveranceOrderedByCourt   string
	CourtOrderDecisionDate    sirius.DateString
	CourtOrderReceivedDate    sirius.DateString
	SeveranceType             string
	WordsToBeRemoved          string
	AmendedRestrictions       string
	FormAction                string
	ConfirmRestrictionDetails CourtOrderRestrictionDetails
	Success                   bool
}

func ManageRestrictions(client ManageRestrictionsClient, manageTmpl template.Template, confirmTmpl template.Template) Handler {
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

		data := manageRestrictionsData{
			CaseSummary:             cs,
			SeveranceAction:         postFormString(r, "severanceAction"),
			DonorConsentGiven:       postFormString(r, "donorConsentGiven"),
			SeveranceOrderedByCourt: postFormString(r, "severanceOrdered"),
			CourtOrderDecisionDate:  postFormDateString(r, "courtOrderDecisionMade"),
			CourtOrderReceivedDate:  postFormDateString(r, "courtOrderReceived"),
			SeveranceType:           postFormString(r, "severanceType"),
			WordsToBeRemoved:        postFormString(r, "removedWords"),
			AmendedRestrictions:     postFormString(r, "updatedRestrictions"),
			XSRFToken:               ctx.XSRFToken,
			Error:                   sirius.ValidationError{Field: sirius.FieldErrors{}},
			CaseUID:                 caseUID,
		}

		data.FormAction = r.FormValue("action")
		if data.FormAction == "" && data.CaseSummary.DigitalLpa.SiriusData.Application.SeveranceStatus == "REQUIRED" {
			data.FormAction = "donor-consent"

			if data.CaseSummary.DigitalLpa.SiriusData.Application.SeveranceApplication != nil && *data.CaseSummary.DigitalLpa.SiriusData.Application.SeveranceApplication.HasDonorConsented {
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

					if data.SeveranceAction == "severance-application-not-required" && taskID != 0 {
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

			if data.FormAction == "court-order" || data.FormAction == "edit-restrictions" {
				if data.CourtOrderDecisionDate == "" {
					w.WriteHeader(http.StatusBadRequest)
					data.Error.Field["courtOrderDecisionDate"] = map[string]string{
						"reason": "Enter or select the date court order was made",
					}
				}

				if data.CourtOrderReceivedDate == "" {
					w.WriteHeader(http.StatusBadRequest)
					data.Error.Field["courtOrderReceivedDate"] = map[string]string{
						"reason": "Enter or select the date court order was issued",
					}
				}

				data.ConfirmRestrictionDetails = CourtOrderRestrictionDetails{
					SelectedCourtOrderDecisionDate: data.CourtOrderDecisionDate,
					SelectedCourtOrderReceivedDate: data.CourtOrderReceivedDate,
				}

				switch data.SeveranceOrderedByCourt {
				case "severance-not-ordered":
					data.ConfirmRestrictionDetails.SelectedSeveranceActionDetail = "No words are to be removed"
				case "severance-ordered":
					switch data.SeveranceType {
					case "severance-partial":
						data.ConfirmRestrictionDetails.SelectedSeveranceActionDetail = "Some words are to be removed"
						data.ConfirmRestrictionDetails.RemovedWords = data.WordsToBeRemoved
						data.ConfirmRestrictionDetails.ChangedRestrictions = data.AmendedRestrictions
					case "severance-not-partial":
						data.ConfirmRestrictionDetails.SelectedSeveranceActionDetail = "All restrictions and conditions are to be removed"

					}

				default:
					w.WriteHeader(http.StatusBadRequest)
					data.Error.Field["severanceOrderedAction"] = map[string]string{
						"reason": "Select if severance of the restrictions and conditions has been ordered",
					}
				}

				data.ConfirmRestrictionDetails.SelectedSeveranceAction = data.SeveranceOrderedByCourt
				data.ConfirmRestrictionDetails.SelectedSeveranceType = data.SeveranceType
				data.ConfirmRestrictionDetails.SelectedFormAction = data.FormAction

				if !data.Error.Any() {
					if data.FormAction != "edit-restrictions" && data.SeveranceOrderedByCourt == "severance-ordered" && data.SeveranceType == "severance-partial" {
						data.FormAction = "edit-restrictions"
						return manageTmpl(w, data)
					} else if !postFormKeySet(r, "confirmRestrictions") {

						if data.SeveranceOrderedByCourt == "severance-ordered" && data.SeveranceType == "severance-partial" {
							if data.WordsToBeRemoved == "" {
								w.WriteHeader(http.StatusBadRequest)
								data.Error.Field["wordsToBeRemoved"] = map[string]string{
									"reason": "Enter words to be removed",
								}
							}

							if data.AmendedRestrictions == "" {
								w.WriteHeader(http.StatusBadRequest)
								data.Error.Field["amendedRestrictions"] = map[string]string{
									"reason": "Enter the updated restrictions and conditions",
								}
							}
						}

						return confirmTmpl(w, data)
					} else {
						severanceApplication := sirius.SeveranceApplication{}

						if data.CourtOrderDecisionDate != "" {
							severanceApplication.CourtOrderDecisionMade = data.CourtOrderDecisionDate
						}

						if data.CourtOrderReceivedDate != "" {
							severanceApplication.CourtOrderReceived = data.CourtOrderReceivedDate
						}

						if data.AmendedRestrictions != "" {
							severanceApplication.UpdatedRestrictions = data.AmendedRestrictions
						}

						switch data.SeveranceOrderedByCourt {
						case "severance-ordered", "severance-not-ordered":
							isSeveranceOrdered := true
							if data.SeveranceOrderedByCourt == "severance-not-ordered" {
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
			}
		}
		return manageTmpl(w, data)
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
