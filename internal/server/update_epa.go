package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type UpdateEpaClient interface {
	Epa(ctx sirius.Context, id int) (sirius.Epa, error)
	UpdateEpa(ctx sirius.Context, caseID int, epa sirius.Epa) error
}

type updateEpaData struct {
	XSRFToken       string
	Success         bool
	IsUpdate        bool
	IsPartial       bool
	Error           sirius.ValidationError
	Epa             sirius.Epa
	DonorId         int
	AppointmentType string
	Title           string
	CaseId          int
}

func UpdateEpa(client UpdateEpaClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		donorID, err := strToIntOrStatusError(r.FormValue("id"))
		if err != nil {
			return err
		}

		caseID, err := strToIntOrStatusError(r.FormValue("caseId"))
		if err != nil {
			return err
		}

		data := updateEpaData{
			XSRFToken: ctx.XSRFToken,
			DonorId:   donorID,
			Title:     "Edit EPA",
			IsPartial: ctx.IsPartial,
			IsUpdate:  true,
			CaseId:    caseID,
		}

		caseItem, err := client.Epa(ctx, data.CaseId)
		if err != nil {
			return err
		}

		data.Epa = caseItem

		if caseItem.CaseAttorneyJointly != nil && *caseItem.CaseAttorneyJointly {
			data.AppointmentType = "jointly"
		} else if caseItem.CaseAttorneySingular != nil && *caseItem.CaseAttorneySingular {
			data.AppointmentType = "singular"
		} else if caseItem.CaseAttorneyJointlyAndSeverally != nil && *caseItem.CaseAttorneyJointlyAndSeverally {
			data.AppointmentType = "jointly-and-severally"
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

			if !*epa.DonorHasOtherEpas {
				epa.OtherEpaInfo = ""
			}

			err = client.UpdateEpa(ctx, data.CaseId, epa)
			data.Epa, _ = client.Epa(ctx, data.CaseId)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
				if data.Error.Field["receiptDate"] != nil {
					if r.FormValue("addAttorney") != "" || r.FormValue("addCorrespondent") != "" {
						data.Error.Field["receiptDate"]["receiptDate"] = "Enter or select a receipt date to continue to step 3"
					} else {
						data.Error.Field["receiptDate"]["receiptDate"] = "Enter or select a receipt date to save and exit"
					}
				}

				return tmpl(w, data)
			} else if err != nil {
				return err
			}

			if r.FormValue("addAttorney") != "" {
				return RedirectError(fmt.Sprintf("/create-attorney?id=%d&caseId=%d&caseType=epa", donorID, data.CaseId))
			} else if updateAttorney := r.FormValue("updateAttorney"); updateAttorney != "" {
				attorneyID, err := strToIntOrStatusError(updateAttorney)
				if err != nil {
					return err
				}

				return RedirectError(fmt.Sprintf("/update-attorney?id=%d&caseId=%d&caseType=epa&attorneyId=%d", donorID, data.CaseId, attorneyID))
			} else if r.FormValue("addCorrespondent") != "" {
				if len(data.Epa.Attorneys) > 0 {
					return RedirectError(fmt.Sprintf("/select-or-create-correspondent?id=%d&caseId=%d&caseType=epa", donorID, data.CaseId))
				}
				return RedirectError(fmt.Sprintf("/create-correspondent?id=%d&caseId=%d&caseType=epa", donorID, data.CaseId))
			} else if r.FormValue("updateCorrespondent") != "" {
				return RedirectError(fmt.Sprintf("/edit-correspondent?id=%d&caseId=%d&caseType=epa", donorID, data.CaseId))
			} else {
				data.Success = true
			}
		}

		return tmpl(w, data)
	}
}
