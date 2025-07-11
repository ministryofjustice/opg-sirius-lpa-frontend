package server

import (
	"fmt"
	"net/http"
	"slices"

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
	CaseSummary              sirius.CaseSummary
	ActiveAttorneys          []sirius.LpaStoreAttorney
	InactiveAttorneys        []sirius.LpaStoreAttorney
	RemovedReasons           []sirius.RefDataItem
	Form                     formRemoveAttorney
	RemovedAttorneysDetails  SelectedAttorneyDetails
	RemovedReason            sirius.RefDataItem
	EnabledAttorneysDetails  []SelectedAttorneyDetails
	DecisionAttorneysDetails []AttorneyDetails
	Success                  bool
	Error                    sirius.ValidationError
	XSRFToken                string
	FormName                 string
	Decisions                string
}

func RemoveAnAttorney(client RemoveAnAttorneyClient, removeTmpl template.Template, confirmTmpl template.Template, decisionsTmpl template.Template) Handler {

	return func(w http.ResponseWriter, r *http.Request) error {
		uid := r.PathValue("uid")
		ctx := getContext(r)

		caseSummary, err := client.CaseSummary(ctx, uid)

		if err != nil {
			return err
		}

		data := removeAnAttorneyData{
			CaseSummary: caseSummary,
			XSRFToken:   ctx.XSRFToken,
			Error:       sirius.ValidationError{Field: sirius.FieldErrors{}},
			FormName:    "remove",
			Decisions:   caseSummary.DigitalLpa.LpaStoreData.HowAttorneysMakeDecisions,
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

		allRemovedReasons, err := client.RefDataByCategory(ctx, sirius.AttorneyRemovedReasonCategory)
		if err != nil {
			return err
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

			validateRemoveAttorneyPage(r, &data)

			if !data.Error.Any() {
				submissionStep := r.PostFormValue("step")

				switch submissionStep {
				case "confirm":
					return confirmStep(ctx, client, &data, w)
				case "decision":
					buildAttorneyDetails(&data, allRemovedReasons, true)
					return confirmTmpl(w, data)
				case "remove":
					buildAttorneyDetails(&data, allRemovedReasons, false)
					if lpa.LpaStoreData.HowAttorneysMakeDecisions == "jointly-for-some-severally-for-others" {
						data.ActiveAttorneys = decisionAttorneysListAfterRemoval(lpa.LpaStoreData.Attorneys, data.Form)
						return decisionsTmpl(w, data)
					} else {
						return confirmTmpl(w, data)
					}
				default:
					return removeTmpl(w, data)
				}

			}
		}

		return removeTmpl(w, data)
	}
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

func buildAttorneyDetails(data *removeAnAttorneyData, removedReasons []sirius.RefDataItem, decisionAttorneys bool) {

	for _, att := range data.ActiveAttorneys {
		if att.Uid == data.Form.RemovedAttorneyUid {
			data.RemovedAttorneysDetails = SelectedAttorneyDetails{
				SelectedAttorneyName: att.FirstNames + " " + att.LastName,
				SelectedAttorneyDob:  att.DateOfBirth,
			}
		}
	}

	for _, att := range data.InactiveAttorneys {
		if slices.Contains(data.Form.EnabledAttorneyUids, att.Uid) {
			data.EnabledAttorneysDetails = append(data.EnabledAttorneysDetails, SelectedAttorneyDetails{
				SelectedAttorneyName: att.FirstNames + " " + att.LastName,
				SelectedAttorneyDob:  att.DateOfBirth,
			})
		}
	}

	for _, r := range removedReasons {
		if r.Handle == data.Form.RemovedReason {
			data.RemovedReason = r
		}
	}

	if decisionAttorneys {
		for _, att := range data.ActiveAttorneys {
			if slices.Contains(data.Form.DecisionAttorneysUids, att.Uid) {
				data.DecisionAttorneysDetails = append(data.DecisionAttorneysDetails, AttorneyDetails{
					AttorneyName:    att.FirstNames + " " + att.LastName,
					AttorneyDob:     att.DateOfBirth,
					AppointmentType: att.AppointmentType,
				})
			}
		}
	}
}

func decisionAttorneysListAfterRemoval(attorneys []sirius.LpaStoreAttorney, form formRemoveAttorney) []sirius.LpaStoreAttorney {
	enabled := map[string]bool{}
	for _, uid := range form.EnabledAttorneyUids {
		enabled[uid] = true
	}

	//TODO: Test for when no attorneys enabled
	var attorneysForDecisions []sirius.LpaStoreAttorney
	for _, att := range attorneys {
		switch att.Status {
		case shared.ActiveAttorneyStatus.String():
			if att.Uid != form.RemovedAttorneyUid {
				attorneysForDecisions = append(attorneysForDecisions, att)
			}
		case shared.InactiveAttorneyStatus.String():
			if enabled[att.Uid] {
				attorneysForDecisions = append(attorneysForDecisions, att)
			}
		}
	}

	return attorneysForDecisions
}

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
		for _, att := range data.ActiveAttorneys {
			attorneyDecisions = append(attorneyDecisions, sirius.AttorneyDecisions{
				UID:                      att.Uid,
				CannotMakeJointDecisions: false,
			})
		}
	} else {
		for _, att := range data.ActiveAttorneys {
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
	}

	uid := data.CaseSummary.DigitalLpa.UID

	err := client.ChangeAttorneyStatus(ctx, uid, attorneyUpdatedStatus)
	if ve, ok := err.(sirius.ValidationError); ok {
		w.WriteHeader(http.StatusBadRequest)
		data.Error = ve
	} else if err != nil {
		return err
	}

	err = client.ManageAttorneyDecisions(ctx, uid, attorneyDecisions)

	if ve, ok := err.(sirius.ValidationError); ok {
		w.WriteHeader(http.StatusBadRequest)
		data.Error = ve
	} else if err != nil {
		return err
	}

	SetFlash(w, FlashNotification{Title: "Update saved"})
	return RedirectError(fmt.Sprintf("/lpa/%s", uid))
}
