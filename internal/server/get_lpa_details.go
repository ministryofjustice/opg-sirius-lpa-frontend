package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/sync/errgroup"
)

type GetLpaDetailsClient interface {
	CaseSummaryWithImages(ctx sirius.Context, uid string) (sirius.CaseSummary, error)
	AnomaliesForDigitalLpa(ctx sirius.Context, uid string) ([]sirius.Anomaly, error)
}

type getLpaDetails struct {
	CaseSummary                     sirius.CaseSummary
	DigitalLpa                      sirius.DigitalLpa
	AnomalyDisplay                  *sirius.AnomalyDisplay
	ReviewRestrictions              bool
	CheckedForSeverance             bool
	SeveranceType                   string
	ReplacementAttorneys            []sirius.LpaStoreAttorney
	NonReplacementAttorneys         []sirius.LpaStoreAttorney
	RemovedAttorneys                []sirius.LpaStoreAttorney
	DecisionAttorneys               []sirius.LpaStoreAttorney
	ReplacementTrustCorporations    []sirius.LpaStoreTrustCorporation
	NonReplacementTrustCorporations []sirius.LpaStoreTrustCorporation
	RemovedTrustCorporations        []sirius.LpaStoreTrustCorporation
	DecisionTrustCorporations       []sirius.LpaStoreTrustCorporation
	FlashMessage                    FlashNotification
}

func GetLpaDetails(client GetLpaDetailsClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		uid := r.PathValue("uid")
		ctx := getContext(r)

		var err error
		var data getLpaDetails
		var anomalies []sirius.Anomaly

		group, groupCtx := errgroup.WithContext(ctx.Context)

		group.Go(func() error {
			data.CaseSummary, err = client.CaseSummaryWithImages(ctx.With(groupCtx), uid)
			if err != nil {
				return err
			}
			return nil
		})

		group.Go(func() error {
			// ignore errors: there may be no LPA in the store yet for a given UID (e.g. if it's still a draft)
			// TODO is there a better way to be selective about ignored client errors?
			anomalies, _ = client.AnomaliesForDigitalLpa(ctx.With(groupCtx), uid)

			return nil
		})

		if err := group.Wait(); err != nil {
			return err
		}

		taskList := data.CaseSummary.TaskList
		data.ReviewRestrictions = false
		for _, task := range taskList {
			if task.Name == "Review restrictions and conditions" && task.Status != "Completed" {
				data.ReviewRestrictions = true
			}
		}

		data.CheckedForSeverance = false
		data.SeveranceType = ""
		if data.CaseSummary.DigitalLpa.SiriusData.Application.SeveranceApplication != nil && data.CaseSummary.DigitalLpa.SiriusData.Application.SeveranceApplication.SeveranceOrdered != nil {
			data.CheckedForSeverance = true

			if *data.CaseSummary.DigitalLpa.SiriusData.Application.SeveranceApplication.SeveranceOrdered && data.CaseSummary.DigitalLpa.LpaStoreData.RestrictionsAndConditions == "" {
				data.SeveranceType = "All restrictions and conditions have been removed"
			} else if *data.CaseSummary.DigitalLpa.SiriusData.Application.SeveranceApplication.SeveranceOrdered && data.CaseSummary.DigitalLpa.LpaStoreData.RestrictionsAndConditions != "" {
				data.SeveranceType = "Some words have been removed"
			} else if !*data.CaseSummary.DigitalLpa.SiriusData.Application.SeveranceApplication.SeveranceOrdered {
				data.SeveranceType = "No words have been removed"
			}
		}

		data.AnomalyDisplay = &sirius.AnomalyDisplay{}
		if anomalies != nil {
			data.AnomalyDisplay.Group(&data.CaseSummary.DigitalLpa.LpaStoreData, anomalies)
		}

		data.DigitalLpa = data.CaseSummary.DigitalLpa
		data.FlashMessage, _ = GetFlash(w, r)

		var replacementAttorneys []sirius.LpaStoreAttorney
		var nonReplacementAttorneys []sirius.LpaStoreAttorney
		var removedAttorneys []sirius.LpaStoreAttorney
		var decisionAttorneys []sirius.LpaStoreAttorney
		for _, attorney := range data.DigitalLpa.LpaStoreData.Attorneys {
			if attorney.Status == shared.ActiveAttorneyStatus.String() {
				nonReplacementAttorneys = append(nonReplacementAttorneys, attorney)
			} else if attorney.Status == shared.InactiveAttorneyStatus.String() &&
				attorney.AppointmentType == shared.ReplacementAppointmentType.String() {
				replacementAttorneys = append(replacementAttorneys, attorney)
			} else if attorney.Status == shared.RemovedAttorneyStatus.String() {
				removedAttorneys = append(removedAttorneys, attorney)
			}
			if attorney.Decisions {
				decisionAttorneys = append(decisionAttorneys, attorney)
			}
		}

		data.ReplacementAttorneys = replacementAttorneys
		data.NonReplacementAttorneys = nonReplacementAttorneys
		data.RemovedAttorneys = removedAttorneys
		data.DecisionAttorneys = decisionAttorneys

		var replacementTrustCorporations []sirius.LpaStoreTrustCorporation
		var nonReplacementTrustCorporations []sirius.LpaStoreTrustCorporation
		var removedTrustCorporations []sirius.LpaStoreTrustCorporation
		var decisionTrustCorporations []sirius.LpaStoreTrustCorporation
		for _, trustCorporation := range data.DigitalLpa.LpaStoreData.TrustCorporations {
			if trustCorporation.Status == shared.ActiveAttorneyStatus.String() {
				nonReplacementTrustCorporations = append(nonReplacementTrustCorporations, trustCorporation)
			} else if trustCorporation.Status == shared.InactiveAttorneyStatus.String() &&
				trustCorporation.AppointmentType == shared.ReplacementAppointmentType.String() {
				replacementTrustCorporations = append(replacementTrustCorporations, trustCorporation)
			} else if trustCorporation.Status == shared.RemovedAttorneyStatus.String() {
				removedTrustCorporations = append(removedTrustCorporations, trustCorporation)
			}
			if trustCorporation.Decisions {
				decisionTrustCorporations = append(decisionTrustCorporations, trustCorporation)
			}
		}

		data.ReplacementTrustCorporations = replacementTrustCorporations
		data.NonReplacementTrustCorporations = nonReplacementTrustCorporations
		data.RemovedTrustCorporations = removedTrustCorporations
		data.DecisionTrustCorporations = decisionTrustCorporations

		return tmpl(w, data)
	}
}
