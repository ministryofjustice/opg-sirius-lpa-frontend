package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type CreateEpaClient interface {
	CreateEpa(ctx sirius.Context, donorID int, epa sirius.Epa) (int, error)
	Case(ctx sirius.Context, id int) (sirius.Case, error)
}

type createEpaData struct {
	XSRFToken            string
	Epa                  sirius.Epa
	Case                 sirius.Case
	Success              bool
	Error                sirius.ValidationError
	ShowAllSections      bool
	RelationshipToDonors []sirius.RefDataItem
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
		}

		caseIdStr := r.FormValue("caseId")
		if caseIdStr != "" {
			caseId, err := strToIntOrStatusError(caseIdStr)
			if err == nil {
				caseItem, err := client.Case(ctx, caseId)
				if err == nil {
					data.Case = caseItem
				}
			}
		}

		if r.Method == http.MethodPost {
			epa := sirius.Epa{
				EpaDonorSignatureDate:           postFormDateString(r, "epaDonorSignatureDate"),
				EpaDonorNoticeGivenDate:         postFormDateString(r, "epaDonorNoticeGivenDate"),
				DonorHasOtherEpas:               postFormString(r, "donorHasOtherEpas") == "true",
				ReceiptDate:                     postFormDateString(r, "receiptDate"),
				RegistrationDate:                postFormDateString(r, "registrationDate"),
				DispatchDate:                    postFormDateString(r, "dispatchDate"),
				CaseAttorneySingular:            r.FormValue("caseAttorneySingular") == "true",
				CaseAttorneyJointlyAndSeverally: r.FormValue("caseAttorneyJointlyAndSeverally") == "true",
				CaseAttorneyJointly:             r.FormValue("caseAttorneyJointly") == "true",
				AttorneyRelationshipToDonor:     postFormString(r, "relationshipToDonors"),
			}
			data.Epa = epa

			caseID, err := client.CreateEpa(ctx, donorID, epa)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
			} else if err != nil {
				return err
			}

			return RedirectError(fmt.Sprintf("/edit-epa?caseId=%d", caseID))
		}

		return tmpl(w, data)
	}
}
