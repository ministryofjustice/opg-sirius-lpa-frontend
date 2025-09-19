package server

import (
	"slices"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

func buildAttorneyDetails(data *removeAnAttorneyData, removedReasons []sirius.RefDataItem) {

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

	for _, r := range removedReasons {
		if r.Handle == data.Form.RemovedReason {
			data.RemovedReason = r
		}
	}

	if len(data.Form.DecisionAttorneysUids) > 0 {
		for _, att := range data.CaseSummary.DigitalLpa.LpaStoreData.Attorneys {
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
