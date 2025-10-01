package server

import (
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
					return confirmStep(ctx, client, &data, w)
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
