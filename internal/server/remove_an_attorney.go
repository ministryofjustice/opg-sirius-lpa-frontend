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

			validateRemoveAttorneyPage(r, &data)

			if submissionStep == "decision" && data.Decisions == "jointly-for-some-severally-for-others" {
				validateManageAttorneysPage(r, &data)
				if _, ok := data.Error.Field["decisionAttorney"]; ok {
					data.DecisionAttorneys = decisionAttorneysListAfterRemoval(lpa.LpaStoreData.Attorneys, data.Form)
					return decisionsTmpl(w, data)
				}
			}

			if !data.Error.Any() {
				switch submissionStep {
				case "confirm":
					attorneyUpdatedStatus := updateAttorneyStatus(&data)
					attorneyDecisions := updateAttorneyDecision(&data)
					return confirmStep(ctx, client, &data, w, attorneyUpdatedStatus, attorneyDecisions)
				case "decision":
					buildAttorneyDetails(&data, allRemovedReasons)
					return confirmTmpl(w, data)
				default: //"remove"
					buildAttorneyDetails(&data, allRemovedReasons)
					if data.Decisions != "jointly-for-some-severally-for-others" {
						return confirmTmpl(w, data)
					}
					data.DecisionAttorneys = decisionAttorneysListAfterRemoval(lpa.LpaStoreData.Attorneys, data.Form)
					return decisionsTmpl(w, data)
				}
			}
		}

		return removeTmpl(w, data)
	}
}

func updateAttorneyStatus(data *removeAnAttorneyData) []sirius.AttorneyUpdatedStatus {
	var attorneyUpdatedStatus []sirius.AttorneyUpdatedStatus
	attorneyUpdatedStatus = removeAttorneyUpdateStatus(data, attorneyUpdatedStatus)

	if len(data.Form.EnabledAttorneyUids) > 0 {
		attorneyUpdatedStatus = endableAttorneyUpdateStatus(data, attorneyUpdatedStatus)
	}

	return attorneyUpdatedStatus
}

func endableAttorneyUpdateStatus(data *removeAnAttorneyData, attorneyUpdatedStatus []sirius.AttorneyUpdatedStatus) []sirius.AttorneyUpdatedStatus {
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
	return attorneyUpdatedStatus
}

func removeAttorneyUpdateStatus(data *removeAnAttorneyData, attorneyUpdatedStatus []sirius.AttorneyUpdatedStatus) []sirius.AttorneyUpdatedStatus {
	for _, att := range data.ActiveAttorneys {
		if att.Uid == data.Form.RemovedAttorneyUid {
			attorneyUpdatedStatus = append(attorneyUpdatedStatus, sirius.AttorneyUpdatedStatus{
				UID:           att.Uid,
				Status:        shared.RemovedAttorneyStatus.String(),
				RemovedReason: data.Form.RemovedReason,
			})
		}
	}
	return attorneyUpdatedStatus
}

func updateAttorneyDecision(data *removeAnAttorneyData) []sirius.AttorneyDecisions {
	var attorneyDecisions []sirius.AttorneyDecisions

	if data.Form.SkipDecisionAttorney == "yes" {
		attorneyDecisions = attorneyCannotMakeJointDecisionsUpdate(data, attorneyDecisions)
	} else {
		data.DecisionAttorneys = decisionAttorneysListAfterRemoval(data.CaseSummary.DigitalLpa.LpaStoreData.Attorneys, data.Form)
		attorneyDecisions = updateSelectedAttorneysThatCannotMakeJointDecisions(data.DecisionAttorneys, data.Form.DecisionAttorneysUids, attorneyDecisions)
		attorneyDecisions = updateRemovedAttorneyToCannotMakeJointDecisions(data.CaseSummary.DigitalLpa.LpaStoreData.Attorneys, data.Form.RemovedAttorneyUid, attorneyDecisions)
	}

	return attorneyDecisions
}

func attorneyCannotMakeJointDecisionsUpdate(data *removeAnAttorneyData, attorneyDecisions []sirius.AttorneyDecisions) []sirius.AttorneyDecisions {
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
	return attorneyDecisions
}

func decisionAttorneysListAfterRemoval(attorneys []sirius.LpaStoreAttorney, form formRemoveAttorney) []sirius.LpaStoreAttorney {
	enabledAttorneyUids := make(map[string]bool)
	for _, uid := range form.EnabledAttorneyUids {
		enabledAttorneyUids[uid] = true
	}

	var attorneysForDecisions []sirius.LpaStoreAttorney
	for _, att := range attorneys {
		switch att.Status {
		case shared.ActiveAttorneyStatus.String():
			if att.Uid != form.RemovedAttorneyUid {
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
	data *removeAnAttorneyData,
	w http.ResponseWriter,
	attorneyUpdatedStatus []sirius.AttorneyUpdatedStatus,
	attorneyDecisions []sirius.AttorneyDecisions,
) error {
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

func buildAttorneyDetails(data *removeAnAttorneyData, removedReasons []sirius.RefDataItem) {
	data.RemovedAttorneysDetails = updateRemovedAttorneysDetails(data.ActiveAttorneys, data.Form.RemovedAttorneyUid)
	data.EnabledAttorneysDetails = updateEnabledAttorneysDetails(data.Form.EnabledAttorneyUids, data.InactiveAttorneys)

	for _, r := range removedReasons {
		if r.Handle == data.Form.RemovedReason {
			data.RemovedReason = r
		}
	}

	if len(data.Form.DecisionAttorneysUids) > 0 {
		data.DecisionAttorneysDetails = updateDecisionAttorneyDetails(data.CaseSummary.DigitalLpa.LpaStoreData.Attorneys, data.Form.DecisionAttorneysUids)
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

func validateRemoveAttorneyPage(r *http.Request, data *removeAnAttorneyData) {
	if data.Form.RemovedAttorneyUid == "" {
		data.Error.Field["removeAttorney"] = map[string]string{
			"reason": "Please select an attorney for removal",
		}
	}

	if data.Form.RemovedReason == "" {
		data.Error.Field["removedReason"] = map[string]string{
			"reason": "Please select a reason for removal",
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
}

func validateManageAttorneysPage(r *http.Request, data *removeAnAttorneyData) {
	if (len(data.Form.DecisionAttorneysUids) == 0 && !postFormCheckboxChecked(r, "skipDecisionAttorney", "yes")) ||
		(len(data.Form.DecisionAttorneysUids) > 0 && postFormCheckboxChecked(r, "skipDecisionAttorney", "yes")) {
		data.Error.Field["decisionAttorney"] = map[string]string{
			"reason": "Select who cannot make joint decisions, or select 'Joint decisions can be made by all attorneys'",
		}
	}
}
