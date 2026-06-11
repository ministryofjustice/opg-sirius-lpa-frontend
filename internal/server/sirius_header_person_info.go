package server

import (
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type SiriusHeaderPeopleInfoClient interface {
	Case(ctx sirius.Context, id int) (sirius.Case, error)
}

type siriusHeaderPeopleInfoData struct {
	XSRFToken      string
	CaseID         int
	SelectedID     int
	Case           sirius.Case
	SelectedPerson sirius.Recipient
}

func SiriusHeaderPeopleInfo(client SiriusHeaderPeopleInfoClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		caseId, err := strToIntOrStatusError(r.FormValue("id"))
		if err != nil {
			return err
		}

		ctx := getContext(r)
		data := siriusHeaderPeopleInfoData{
			XSRFToken: ctx.XSRFToken,
			CaseID:    caseId,
		}

		caseItem, err := client.Case(ctx, caseId)
		if err != nil {
			return err
		}

		selectedId := r.FormValue("selected")
		selectedIdInt := caseItem.Donor.ID
		if selectedId != "" {
			selectedIdInt, err = strconv.Atoi(selectedId)
			if err != nil {
				return err
			}
		}

		selectedPerson := getSelectedPerson(caseItem, selectedIdInt)
		data.Case = caseItem
		data.SelectedPerson = selectedPerson
		data.SelectedID = selectedIdInt

		return tmpl(w, data)
	}
}

func getSelectedPerson(caseItem sirius.Case, selectedId int) sirius.Recipient {
	if selectedId != 0 {
		if caseItem.Donor.ID == selectedId {
			return *caseItem.Donor
		}
		for _, attorney := range caseItem.Attorneys {
			if attorney.ID == selectedId {
				return attorney
			}
		}
		for _, replacementAttorney := range caseItem.ReplacementAttorneys {
			if replacementAttorney.ID == selectedId {
				return replacementAttorney
			}
		}
		for _, trustCorporation := range caseItem.TrustCorporations {
			if trustCorporation.ID == selectedId {
				return trustCorporation
			}
		}
		for _, certificateProvider := range caseItem.CertificateProviders {
			if certificateProvider.ID == selectedId {
				return certificateProvider
			}
		}
		for _, notifiedPerson := range caseItem.NotifiedPersons {
			if notifiedPerson.ID == selectedId {
				return notifiedPerson
			}
		}
		if caseItem.Correspondent != nil && caseItem.Correspondent.ID == selectedId {
			return caseItem.Correspondent
		}
	}
	return *caseItem.Donor
}
