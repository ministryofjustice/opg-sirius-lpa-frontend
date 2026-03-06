package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type EditEpaClient interface {
	UpdateEpa(ctx sirius.Context, caseId int, epa sirius.Case) error
	Case(sirius.Context, int) (sirius.Case, error)
}

type editEpaData struct {
	XSRFToken            string
	Case                 sirius.Case
	Success              bool
	Error                sirius.ValidationError
	ShowAllSections      bool
	RelationshipToDonors []sirius.RefDataItem
}

func EditEpa(client EditEpaClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		caseId, err := strToIntOrStatusError(r.FormValue("caseId"))
		if err != nil {
			return err
		}

		ctx := getContext(r)
		caseitem, err := client.Case(ctx, caseId)
		if err != nil {
			return err
		}

		data := editEpaData{
			XSRFToken: ctx.XSRFToken,
			Case:      caseitem,
			RelationshipToDonors: []sirius.RefDataItem{
				{Handle: "civil partner", Label: "civil partner"},
				{Handle: "child", Label: "child"},
				{Handle: "solicitor", Label: "solicitor"},
				{Handle: "other", Label: "other"},
				{Handle: "other professional", Label: "other professional"},
			},
		}

		if r.Method == http.MethodPost {
			if r.FormValue("showAllSections") == "true" {
				data.ShowAllSections = true
				return tmpl(w, data)
			}
			if r.FormValue("showAllSections") == "false" {
				data.ShowAllSections = false
				return tmpl(w, data)
			}

			epa := sirius.Case{
				EpaDonorSignatureDate:           postFormDateString(r, "epaDonorSignatureDate"),
				EpaDonorNoticeGivenDate:         postFormDateString(r, "epaDonorNoticeGivenDate"),
				DonorHasOtherEpas:               postFormString(r, "donorHasOtherEpas") == "true",
				ReceiptDate:                     postFormDateString(r, "receiptDate"),
				CaseAttorneySingular:            r.FormValue("caseAttorneySingular") == "true",
				CaseAttorneyJointlyAndSeverally: r.FormValue("caseAttorneyJointlyAndSeverally") == "true",
				CaseAttorneyJointly:             r.FormValue("caseAttorneyJointly") == "true",
				AttorneyRelationshipToDonor:     postFormString(r, "relationshipToDonors"),
			}

			err = client.UpdateEpa(ctx, caseId, epa)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
			} else if err != nil {
				return err
			} else {
				data.Success = true
			}
		}

		if r.FormValue("isEditing") == "true" {
			return RedirectError(fmt.Sprintf("/edit-epa?caseId=%d", caseId))
		}

		return tmpl(w, data)
	}
}
