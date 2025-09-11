package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

func confirmStep(ctx sirius.Context, client RemoveAnAttorneyClient, data *removeAnAttorneyData, w http.ResponseWriter) error {
	var attorneyUpdatedStatus []sirius.AttorneyUpdatedStatus

	for _, att := range data.ActiveAttorneys {
		if att.Uid == data.Form.RemovedAttorneyUid {
			attorneyUpdatedStatus = append(attorneyUpdatedStatus, sirius.AttorneyUpdatedStatus{
				UID:           att.Uid,
				Status:        shared.RemovedAttorneyStatus.String(),
				RemovedReason: data.Form.RemovedReason,
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

	var attorneyDecisions []sirius.AttorneyDecisions

	if data.Form.SkipDecisionAttorney == "yes" {
		processedAttorneys := make(map[string]bool)

		for _, att := range append(data.ActiveAttorneys, data.DecisionAttorneys...) {
			if processedAttorneys[att.Uid] {
				continue
			}
			processedAttorneys[att.Uid] = true

			attorneyDecisions = append(attorneyDecisions, sirius.AttorneyDecisions{
				UID:                      att.Uid,
				CannotMakeJointDecisions: false,
			})
		}
	} else {
		data.DecisionAttorneys = decisionAttorneysListAfterRemoval(data.CaseSummary.DigitalLpa.LpaStoreData.Attorneys, data.Form)
		for _, att := range data.DecisionAttorneys {
			isChecked := false
			for _, selectedUid := range data.Form.DecisionAttorneysUids {
				if selectedUid == att.Uid {
					isChecked = true
					break
				}
			}
			attorneyDecisions = append(attorneyDecisions, sirius.AttorneyDecisions{
				UID:                      att.Uid,
				CannotMakeJointDecisions: isChecked,
			})
		}

		for _, att := range data.CaseSummary.DigitalLpa.LpaStoreData.Attorneys {
			if att.Uid == data.Form.RemovedAttorneyUid {
				attorneyDecisions = append(attorneyDecisions, sirius.AttorneyDecisions{
					UID:                      att.Uid,
					CannotMakeJointDecisions: false,
				})
				break
			}
		}
	}

	uid := data.CaseSummary.DigitalLpa.UID

	err := client.ChangeAttorneyStatus(ctx, uid, attorneyUpdatedStatus)
	if ve, ok := err.(sirius.ValidationError); ok {
		w.WriteHeader(http.StatusBadRequest)
		data.Error = ve
	} else if err != nil {
		return err
	}

	if data.Decisions == "jointly-for-some-severally-for-others" {
		err = client.ManageAttorneyDecisions(ctx, uid, attorneyDecisions)

		if ve, ok := err.(sirius.ValidationError); ok {
			w.WriteHeader(http.StatusBadRequest)
			data.Error = ve
		} else if err != nil {
			return err
		}
	}

	SetFlash(w, FlashNotification{Title: "Update saved"})
	return RedirectError(fmt.Sprintf("/lpa/%s", uid))
}
