package server

import (
	"fmt"
	"net/http"
	"slices"

	"golang.org/x/sync/errgroup"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type RemoveAnAttorneyClient interface {
	CaseSummary(sirius.Context, string) (sirius.CaseSummary, error)
	ChangeAttorneyStatus(sirius.Context, string, []sirius.AttorneyUpdatedStatus) error
	RefDataByCategory(sirius.Context, string) ([]sirius.RefDataItem, error)
	ManageAttorneyDecisions(sirius.Context, string, []sirius.AttorneyDecisions) error
}

type formRemoveAttorney struct {
	RemovedAttorneyUid    string   `form:"removedAttorney"`
	RemovedReason         string   `form:"removedReason"`
	EnabledAttorneyUids   []string `form:"enabledAttorney"`
	SkipEnableAttorney    string   `form:"skipEnableAttorney"`
	DecisionAttorneysUids []string `form:"decisionAttorney"`
	SkipDecisionAttorney  string   `form:"skipDecisionAttorney"`
}

type SelectedAttorneyDetails struct {
	SelectedAttorneyName string
	SelectedAttorneyDob  string
}

type removeAnAttorneyData struct {
	CaseSummary                  sirius.CaseSummary
	ActiveAttorneys              []sirius.LpaStoreAttorney
	InactiveAttorneys            []sirius.LpaStoreAttorney
	DecisionAttorneys            []sirius.LpaStoreAttorney
	RemovedReasons               []sirius.RefDataItem
	Form                         formRemoveAttorney
	RemovedAttorneysDetails      SelectedAttorneyDetails
	RemovedReason                sirius.RefDataItem
	EnabledAttorneysDetails      []SelectedAttorneyDetails
	DecisionAttorneysDetails     []AttorneyDetails
	ActiveAttorneyCount          int
	ReplacementAttorneyCount     int
	Success                      bool
	Error                        sirius.ValidationError
	XSRFToken                    string
	FormName                     string
	Decisions                    string
	ReplacementAttorneyDecisions string
}

func RemoveAnAttorney(client RemoveAnAttorneyClient, removeTmpl template.Template, confirmTmpl template.Template, decisionsTmpl template.Template) Handler {

	return func(w http.ResponseWriter, r *http.Request) error {
		uid := r.PathValue("uid")
		ctx := getContext(r)

		var caseSummary sirius.CaseSummary
		var allRemovedReasons []sirius.RefDataItem
		var err error

		group, groupCtx := errgroup.WithContext(ctx.Context)

		group.Go(func() error {
			caseSummary, err = client.CaseSummary(ctx.With(groupCtx), uid)
			if err != nil {
				return err
			}
			return nil
		})

		group.Go(func() error {
			allRemovedReasons, err = client.RefDataByCategory(ctx, sirius.AttorneyRemovedReasonCategory)
			if err != nil {
				return err
			}

			return nil
		})

		if err := group.Wait(); err != nil {
			return err
		}

		data := removeAnAttorneyData{
			CaseSummary:                  caseSummary,
			XSRFToken:                    ctx.XSRFToken,
			Error:                        sirius.ValidationError{Field: sirius.FieldErrors{}},
			FormName:                     "remove",
			Decisions:                    caseSummary.DigitalLpa.LpaStoreData.HowAttorneysMakeDecisions,
			ReplacementAttorneyDecisions: caseSummary.DigitalLpa.LpaStoreData.HowReplacementAttorneysMakeDecisions,
		}

		lpa := data.CaseSummary.DigitalLpa

		for _, attorney := range lpa.LpaStoreData.Attorneys {
			if (attorney.Status == shared.RemovedAttorneyStatus.String()) ||
				(attorney.AppointmentType == shared.ReplacementAppointmentType.String() &&
					attorney.Status == shared.InactiveAttorneyStatus.String()) {
				continue
			}

			data.ActiveAttorneys = append(data.ActiveAttorneys, attorney)
			data.ActiveAttorneyCount++
		}

		for _, attorney := range lpa.LpaStoreData.Attorneys {
			if attorney.Status == shared.RemovedAttorneyStatus.String() || attorney.Status == shared.ActiveAttorneyStatus.String() {
				continue
			}

			data.InactiveAttorneys = append(data.InactiveAttorneys, attorney)
			data.ReplacementAttorneyCount++
		}

		if data.ActiveAttorneyCount > 1 && data.ReplacementAttorneyCount > 1 && data.ReplacementAttorneyDecisions == "" {
			data.ReplacementAttorneyDecisions = data.Decisions
		}

		for _, removedReason := range allRemovedReasons {
			if slices.Contains(removedReason.ValidSubTypes, lpa.SiriusData.Subtype) {
				data.RemovedReasons = append(data.RemovedReasons, removedReason)
			}
		}

		if r.Method == http.MethodPost {

			err = decoder.Decode(&data.Form, r.PostForm)
			if err != nil {
				return err
			}

			submissionStep := r.PostFormValue("step")

			validateRemoveAttorneyPage(r, data.Form.RemovedAttorneyUid, data.Form.RemovedReason, data.Form.EnabledAttorneyUids, data.Error.Field)

			if submissionStep == "decision" && data.Decisions == "jointly-for-some-severally-for-others" {
				validateManageAttorneysPage(r, data.Form.DecisionAttorneysUids, data.Error.Field)
				if _, ok := data.Error.Field["decisionAttorney"]; ok {
					data.DecisionAttorneys = decisionAttorneysListAfterRemoval(lpa.LpaStoreData.Attorneys, data.Form.EnabledAttorneyUids, data.Form.RemovedAttorneyUid)
					return decisionsTmpl(w, data)
				}
			}

			if !data.Error.Any() {
				switch submissionStep {
				case "confirm":
					attorneyUpdatedStatus := updateAttorneyStatus(data.ActiveAttorneys, data.Form.RemovedAttorneyUid, data.Form.RemovedReason, data.InactiveAttorneys, data.Form.EnabledAttorneyUids)
					attorneyDecisions := updateAttorneyDecision(data.Form.SkipDecisionAttorney, data.ActiveAttorneys, data.DecisionAttorneys, data.CaseSummary.DigitalLpa.LpaStoreData.Attorneys, data.Form.EnabledAttorneyUids, data.Form.RemovedAttorneyUid, data.Form.DecisionAttorneysUids)
					return confirmStep(ctx, client, data.CaseSummary.DigitalLpa.UID, data.Error, data.Decisions, w, attorneyUpdatedStatus, attorneyDecisions)
				case "decision":
					data.RemovedAttorneysDetails = updateRemovedAttorneysDetails(data.ActiveAttorneys, data.Form.RemovedAttorneyUid)
					data.EnabledAttorneysDetails = updateEnabledAttorneysDetails(data.Form.EnabledAttorneyUids, data.InactiveAttorneys)

					for _, r := range allRemovedReasons {
						if r.Handle == data.Form.RemovedReason {
							data.RemovedReason = r
						}
					}

					if len(data.Form.DecisionAttorneysUids) > 0 {
						data.DecisionAttorneysDetails = updateDecisionAttorneyDetails(data.CaseSummary.DigitalLpa.LpaStoreData.Attorneys, data.Form.DecisionAttorneysUids)
					}

					return confirmTmpl(w, data)
				default: //"remove"
					data.RemovedAttorneysDetails = updateRemovedAttorneysDetails(data.ActiveAttorneys, data.Form.RemovedAttorneyUid)
					data.EnabledAttorneysDetails = updateEnabledAttorneysDetails(data.Form.EnabledAttorneyUids, data.InactiveAttorneys)

					for _, r := range allRemovedReasons {
						if r.Handle == data.Form.RemovedReason {
							data.RemovedReason = r
						}
					}

					if len(data.Form.DecisionAttorneysUids) > 0 {
						data.DecisionAttorneysDetails = updateDecisionAttorneyDetails(data.CaseSummary.DigitalLpa.LpaStoreData.Attorneys, data.Form.DecisionAttorneysUids)
					}

					if data.Decisions != "jointly-for-some-severally-for-others" {
						return confirmTmpl(w, data)
					}
					data.DecisionAttorneys = decisionAttorneysListAfterRemoval(lpa.LpaStoreData.Attorneys, data.Form.EnabledAttorneyUids, data.Form.RemovedAttorneyUid)
					return decisionsTmpl(w, data)
				}
			}
		}

		return removeTmpl(w, data)
	}
}

func updateAttorneyStatus(
	activeAttorneys []sirius.LpaStoreAttorney,
	removedAttorneyUid string,
	removedReason string,
	inactiveAttorneys []sirius.LpaStoreAttorney,
	enabledAttorneyUids []string,
) []sirius.AttorneyUpdatedStatus {
	var attorneyUpdatedStatus []sirius.AttorneyUpdatedStatus
	attorneyUpdatedStatus = removeAttorneyUpdateStatus(activeAttorneys, removedAttorneyUid, removedReason, attorneyUpdatedStatus)

	if len(enabledAttorneyUids) > 0 {
		attorneyUpdatedStatus = endableAttorneyUpdateStatus(inactiveAttorneys, enabledAttorneyUids, attorneyUpdatedStatus)
	}

	return attorneyUpdatedStatus
}

func removeAttorneyUpdateStatus(activeAttorneys []sirius.LpaStoreAttorney, removedAttorneyUid string, removedReason string, attorneyUpdatedStatus []sirius.AttorneyUpdatedStatus) []sirius.AttorneyUpdatedStatus {
	for _, att := range activeAttorneys {
		if att.Uid == removedAttorneyUid {
			attorneyUpdatedStatus = append(attorneyUpdatedStatus, sirius.AttorneyUpdatedStatus{
				UID:           att.Uid,
				Status:        shared.RemovedAttorneyStatus.String(),
				RemovedReason: removedReason,
			})
		}
	}
	return attorneyUpdatedStatus
}

func endableAttorneyUpdateStatus(inactiveAttorneys []sirius.LpaStoreAttorney, enabledAttorneyUids []string, attorneyUpdatedStatus []sirius.AttorneyUpdatedStatus) []sirius.AttorneyUpdatedStatus {
	for _, att := range inactiveAttorneys {
		for _, enabledAttUid := range enabledAttorneyUids {
			if att.Uid == enabledAttUid {
				attorneyUpdatedStatus = append(attorneyUpdatedStatus, sirius.AttorneyUpdatedStatus{
					UID:    att.Uid,
					Status: shared.ActiveAttorneyStatus.String(),
				})
			}
		}
	}
	return attorneyUpdatedStatus
}

func updateAttorneyDecision(
	skipDecisionAttorney string,
	activeAttorneys []sirius.LpaStoreAttorney,
	decisionAttorneys []sirius.LpaStoreAttorney,
	lpaStoreAttorneys []sirius.LpaStoreAttorney,
	enabledAttorneyUids []string,
	removedAttorneyUid string,
	decisionAttorneysUids []string,
) []sirius.AttorneyDecisions {
	var attorneyDecisions []sirius.AttorneyDecisions

	if skipDecisionAttorney == "yes" {
		attorneyDecisions = attorneyCannotMakeJointDecisionsUpdate(activeAttorneys, decisionAttorneys, attorneyDecisions)
	} else {
		decisionAttorneys = decisionAttorneysListAfterRemoval(lpaStoreAttorneys, enabledAttorneyUids, removedAttorneyUid)
		attorneyDecisions = updateSelectedAttorneysThatCannotMakeJointDecisions(decisionAttorneys, decisionAttorneysUids, attorneyDecisions)
		attorneyDecisions = updateRemovedAttorneyToCannotMakeJointDecisions(lpaStoreAttorneys, removedAttorneyUid, attorneyDecisions)
	}

	return attorneyDecisions
}

func attorneyCannotMakeJointDecisionsUpdate(activeAttorneys []sirius.LpaStoreAttorney, decisionAttorney []sirius.LpaStoreAttorney, attorneyDecisions []sirius.AttorneyDecisions) []sirius.AttorneyDecisions {
	processedAttorneys := make(map[string]bool)

	for _, att := range append(activeAttorneys, decisionAttorney...) {
		if processedAttorneys[att.Uid] {
			continue
		}
		processedAttorneys[att.Uid] = true

		attorneyDecisions = append(attorneyDecisions, sirius.AttorneyDecisions{
			UID:                      att.Uid,
			CannotMakeJointDecisions: false,
		})
	}
	return attorneyDecisions
}

func decisionAttorneysListAfterRemoval(attorneys []sirius.LpaStoreAttorney, enabledAttorneyUidsFromForm []string, removedAttorneyUid string) []sirius.LpaStoreAttorney {
	enabledAttorneyUids := make(map[string]bool)
	for _, uid := range enabledAttorneyUidsFromForm {
		enabledAttorneyUids[uid] = true
	}

	var attorneysForDecisions []sirius.LpaStoreAttorney
	for _, att := range attorneys {
		switch att.Status {
		case shared.ActiveAttorneyStatus.String():
			if att.Uid != removedAttorneyUid {
				attorneysForDecisions = append(attorneysForDecisions, att)
			}
		case shared.InactiveAttorneyStatus.String():
			if enabledAttorneyUids[att.Uid] {
				attorneysForDecisions = append(attorneysForDecisions, att)
			}
		}
	}

	return attorneysForDecisions
}

func updateSelectedAttorneysThatCannotMakeJointDecisions(decisionAttorneys []sirius.LpaStoreAttorney, decisionAttorneysUids []string, attorneyDecisions []sirius.AttorneyDecisions) []sirius.AttorneyDecisions {
	for _, att := range decisionAttorneys {
		isChecked := false
		for _, selectedUid := range decisionAttorneysUids {
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
	return attorneyDecisions
}

func updateRemovedAttorneyToCannotMakeJointDecisions(lpaStoreDecisionAttorneys []sirius.LpaStoreAttorney, removedAttorneyUid string, attorneyDecisions []sirius.AttorneyDecisions) []sirius.AttorneyDecisions {
	for _, att := range lpaStoreDecisionAttorneys {
		if att.Uid == removedAttorneyUid {
			attorneyDecisions = append(attorneyDecisions, sirius.AttorneyDecisions{
				UID:                      att.Uid,
				CannotMakeJointDecisions: false,
			})
			break
		}
	}
	return attorneyDecisions
}

func confirmStep(
	ctx sirius.Context,
	client RemoveAnAttorneyClient,
	digitalLpaUid string,
	validationError sirius.ValidationError,
	decisions string,
	w http.ResponseWriter,
	attorneyUpdatedStatus []sirius.AttorneyUpdatedStatus,
	attorneyDecisions []sirius.AttorneyDecisions,
) error {
	err := client.ChangeAttorneyStatus(ctx, digitalLpaUid, attorneyUpdatedStatus)
	if ve, ok := err.(sirius.ValidationError); ok {
		w.WriteHeader(http.StatusBadRequest)
		validationError = ve
	} else if err != nil {
		return err
	}

	if decisions == "jointly-for-some-severally-for-others" {
		err = client.ManageAttorneyDecisions(ctx, digitalLpaUid, attorneyDecisions)

		if ve, ok := err.(sirius.ValidationError); ok {
			w.WriteHeader(http.StatusBadRequest)
			validationError = ve
		} else if err != nil {
			return err
		}
	}

	SetFlash(w, FlashNotification{Title: "Update saved"})
	return RedirectError(fmt.Sprintf("/lpa/%s", digitalLpaUid))
}

func buildAttorneyDetails(
	removedReasons []sirius.RefDataItem,
	formRemovedReason string,
	removedReason sirius.RefDataItem,
) {
	for _, r := range removedReasons {
		if r.Handle == formRemovedReason {
			removedReason = r
		}
	}
}

func updateRemovedAttorneysDetails(activeAttorneys []sirius.LpaStoreAttorney, removedAttorneyUid string) SelectedAttorneyDetails {
	var removedAttorneyDetails SelectedAttorneyDetails
	for _, att := range activeAttorneys {
		if att.Uid == removedAttorneyUid {
			removedAttorneyDetails = SelectedAttorneyDetails{
				SelectedAttorneyName: att.FirstNames + " " + att.LastName,
				SelectedAttorneyDob:  att.DateOfBirth,
			}
		}
	}
	return removedAttorneyDetails
}

func updateEnabledAttorneysDetails(enabledAttorneyUids []string, inactiveAttorneys []sirius.LpaStoreAttorney) []SelectedAttorneyDetails {
	var updatedEnabledAttorneysDetails []SelectedAttorneyDetails
	if len(enabledAttorneyUids) > 0 {
		for _, att := range inactiveAttorneys {
			for _, enabledAttUid := range enabledAttorneyUids {
				if att.Uid == enabledAttUid {
					updatedEnabledAttorneysDetails = append(updatedEnabledAttorneysDetails, SelectedAttorneyDetails{
						SelectedAttorneyName: att.FirstNames + " " + att.LastName,
						SelectedAttorneyDob:  att.DateOfBirth,
					})
					break
				}
			}
		}
	}
	return updatedEnabledAttorneysDetails
}

func updateDecisionAttorneyDetails(digitalAttorneys []sirius.LpaStoreAttorney, decisionAttorneyUids []string) []AttorneyDetails {
	var updatedDecisionAttorneysDetails []AttorneyDetails
	for _, att := range digitalAttorneys {
		if slices.Contains(decisionAttorneyUids, att.Uid) {
			updatedDecisionAttorneysDetails = append(updatedDecisionAttorneysDetails, AttorneyDetails{
				AttorneyName:    att.FirstNames + " " + att.LastName,
				AttorneyDob:     att.DateOfBirth,
				AppointmentType: att.AppointmentType,
			})
		}
	}
	return updatedDecisionAttorneysDetails
}

func validateRemoveAttorneyPage(r *http.Request, removeAnAttorneyUid string, removedReason string, enabledAttorneyUids []string, errorField sirius.FieldErrors) {
	if removeAnAttorneyUid == "" {
		errorField["removeAttorney"] = map[string]string{
			"reason": "Please select an attorney for removal",
		}
	}

	if removedReason == "" {
		errorField["removedReason"] = map[string]string{
			"reason": "Please select a reason for removal",
		}
	}

	if len(enabledAttorneyUids) > 0 && postFormCheckboxChecked(r, "skipEnableAttorney", "yes") {
		errorField["enableAttorney"] = map[string]string{
			"reason": "Please do not select both a replacement attorney and the option to skip",
		}
	}

	if len(enabledAttorneyUids) == 0 && !postFormCheckboxChecked(r, "skipEnableAttorney", "yes") {
		errorField["enableAttorney"] = map[string]string{
			"reason": "Please select either the attorneys that can be enabled or skip the replacement of the attorneys",
		}
	}
}

func validateManageAttorneysPage(r *http.Request, decisionAttorneysUids []string, errorField sirius.FieldErrors) {
	if (len(decisionAttorneysUids) == 0 && !postFormCheckboxChecked(r, "skipDecisionAttorney", "yes")) ||
		(len(decisionAttorneysUids) > 0 && postFormCheckboxChecked(r, "skipDecisionAttorney", "yes")) {
		errorField["decisionAttorney"] = map[string]string{
			"reason": "Select who cannot make joint decisions, or select 'Joint decisions can be made by all attorneys'",
		}
	}
}
