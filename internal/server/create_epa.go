package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type CreateEpaClient interface {
	CreateEpa(ctx sirius.Context, donorID int, epa sirius.Case) (int, error)
	UpdateEpa(ctx sirius.Context, caseID int, epa sirius.Case) error
	Case(ctx sirius.Context, id int) (sirius.Case, error)
}

type createEpaData struct {
	XSRFToken            string
	Case                 sirius.Case
	Success              bool
	Error                sirius.ValidationError
	ShowAllSections      bool
	RelationshipToDonors []sirius.RefDataItem
	Title                string
}

func CreateEpa(client CreateEpaClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		donorID, err := strToIntOrStatusError(r.FormValue("id"))
		if err != nil {
			return err
		}

		ctx := getContext(r)

		data := createEpaData{
			XSRFToken: ctx.XSRFToken,
			RelationshipToDonors: []sirius.RefDataItem{
				{Handle: "civil partner", Label: "civil partner"},
				{Handle: "child", Label: "child"},
				{Handle: "solicitor", Label: "solicitor"},
				{Handle: "other", Label: "other"},
				{Handle: "other professional", Label: "other professional"},
			},
			Title: "Step 1: EPA details",
		}

		caseIdStr := r.FormValue("caseId")
		isEditing := caseIdStr != ""
		if isEditing {
			caseId, err := strToIntOrStatusError(caseIdStr)
			if err == nil {
				caseItem, err := client.Case(ctx, caseId)
				if err == nil {
					data.Case = caseItem
				}
			}
			data.Title = "EPA details"
		}

		if r.Method == http.MethodPost {
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
			data.Case = epa

			var caseId int
			if isEditing {
				caseId, err = strToIntOrStatusError(caseIdStr)
				if err == nil {
					err = client.UpdateEpa(ctx, caseId, epa)
				}
			} else {
				caseId, err = client.CreateEpa(ctx, donorID, epa)
			}

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
				return tmpl(w, data)
			} else if err != nil {
				return err
			}

			if r.FormValue("submit-continue") == "" {
				return RedirectError(fmt.Sprintf("/edit-epa?caseId=%d", caseId))
			}

			if isEditing {
				return RedirectError(fmt.Sprintf("/appointment-epa?caseId=%d&isEditing=true", caseId))
			}
			return RedirectError(fmt.Sprintf("/appointment-epa?caseId=%d", caseId))
		}

		return tmpl(w, data)
	}
}
