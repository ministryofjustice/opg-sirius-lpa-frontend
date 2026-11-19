package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type CreateEpaClient interface {
	CreateEpa(ctx sirius.Context, donorID int, epa sirius.Epa) (sirius.Epa, error)
}

type epaData struct {
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

func CreateEpa(client CreateEpaClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		donorID, err := strToIntOrStatusError(r.FormValue("id"))
		if err != nil {
			return err
		}

		data := epaData{
			XSRFToken: ctx.XSRFToken,
			DonorId:   donorID,
			Title:     "Create an EPA",
			IsPartial: ctx.IsPartial,
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

			epa, err = client.CreateEpa(ctx, donorID, epa)
			if err == nil {
				data.Epa = epa
				data.CaseId = epa.ID
			}

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
			} else if r.FormValue("addCorrespondent") != "" {
				if len(data.Epa.Attorneys) > 0 {
					return RedirectError(fmt.Sprintf("/select-or-create-correspondent?id=%d&caseId=%d&caseType=epa", donorID, data.CaseId))
				}
				return RedirectError(fmt.Sprintf("/create-correspondent?id=%d&caseId=%d&caseType=epa", donorID, data.CaseId))
			} else {
				data.Success = true
			}
		}

		return tmpl(w, data)
	}
}
