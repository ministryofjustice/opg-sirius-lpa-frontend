package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type CreateEpaClient interface {
	Epa(ctx sirius.Context, id int) (sirius.Epa, error)
	CreateEpa(ctx sirius.Context, donorID int, epa sirius.Epa) (sirius.Epa, error)
	UpdateEpa(ctx sirius.Context, caseID int, epa sirius.Epa) error
}

type createEpaData struct {
	XSRFToken       string
	Success         bool
	Error           sirius.ValidationError
	Epa             sirius.Epa
	AppointmentType string
	Title           string
}

func CreateEpa(client CreateEpaClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		donorID, err := strToIntOrStatusError(r.FormValue("id"))
		if err != nil {
			return err
		}

		data := createEpaData{
			XSRFToken: ctx.XSRFToken,
			Title:     "Create an EPA",
		}

		var caseId int
		caseIdStr := r.FormValue("caseId")
		isEditing := caseIdStr != ""
		if isEditing {
			caseId, err = strToIntOrStatusError(caseIdStr)
			if err != nil {
				return err
			}

			caseItem, err := client.Epa(ctx, caseId)
			if err != nil {
				return err
			}
			data.Epa = caseItem
			data.Title = "Edit EPA"

			if caseItem.CaseAttorneyJointly != nil && *caseItem.CaseAttorneyJointly {
				data.AppointmentType = "jointly"
			} else if caseItem.CaseAttorneySingular != nil && *caseItem.CaseAttorneySingular {
				data.AppointmentType = "singular"
			} else if caseItem.CaseAttorneyJointlyAndSeverally != nil && *caseItem.CaseAttorneyJointlyAndSeverally {
				data.AppointmentType = "jointly-and-severally"
			}
		}

		if r.Method == http.MethodPost {
			caseAttorneyValue := r.FormValue("caseAttorney")

			epa := sirius.Epa{
				EpaDonorSignatureDate:   postFormDateString(r, "epaDonorSignatureDate"),
				EpaDonorNoticeGivenDate: postFormDateString(r, "epaDonorNoticeGivenDate"),
				DonorHasOtherEpas:       shared.BoolPtr(postFormString(r, "donorHasOtherEpas") == "true"),
				OtherEpaInfo:            postFormString(r, "otherEpaInfo"),
				Case: sirius.Case{
					ReceiptDate:                     postFormDateString(r, "receiptDate"),
					CaseAttorneySingular:            shared.BoolPtr(caseAttorneyValue == "singular"),
					CaseAttorneyJointlyAndSeverally: shared.BoolPtr(caseAttorneyValue == "jointly-and-severally"),
					CaseAttorneyJointly:             shared.BoolPtr(caseAttorneyValue == "jointly"),
					PaymentByCheque:                 shared.BoolPtr(r.FormValue("paymentByCheque") == "true"),
					PaymentExemption:                shared.BoolPtr(r.FormValue("paymentExemption") == "true"),
					PaymentDate:                     postFormDateString(r, "paymentDate"),
				},
			}
			data.AppointmentType = caseAttorneyValue

			if isEditing {
				err = client.UpdateEpa(ctx, caseId, epa)
			} else {
				epa, err = client.CreateEpa(ctx, donorID, epa)
				if err == nil {
					caseId = epa.ID
				}
			}

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
				return tmpl(w, data)
			} else if err != nil {
				return err
			}

			if r.FormValue("addAttorney") != "" {
				return RedirectError(fmt.Sprintf("/create-attorney?id=%d&caseId=%d", donorID, caseId))
			} else if r.FormValue("addCorrespondent") != "" {
				if len(data.Epa.Attorneys) > 0 {
					return RedirectError(fmt.Sprintf("/select-or-create-correspondent?id=%d&caseId=%d", donorID, caseId))
				}
				return RedirectError(fmt.Sprintf("/create-correspondent?id=%d&caseId=%d", donorID, caseId))
			} else {
				data.Success = true
			}
		}

		return tmpl(w, data)
	}
}
