package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type UpdateEpaClient interface {
	UpdateEpa(ctx sirius.Context, caseId int, epa sirius.Epa) error
	Case(sirius.Context, int) (sirius.Case, error)
}

type updateEpaData struct {
	XSRFToken            string
	Epa                  sirius.Epa
	Case                 sirius.Case
	Success              bool
	Error                sirius.ValidationError
	ShowAllSections      bool
	RelationshipToDonors []sirius.RefDataItem
}

func UpdateEpa(client UpdateEpaClient, tmpl template.Template) Handler {
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

		data := updateEpaData{
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

			if r.FormValue("showAllSections") == "true" {
				data.ShowAllSections = true
				return tmpl(w, data)
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

		return tmpl(w, data)
	}
}
