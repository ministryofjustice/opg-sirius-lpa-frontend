package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

func boolPtr(b bool) *bool {
	return &b
}

type CreateEpaClient interface {
	CreateEpa(ctx sirius.Context, donorID int, epa sirius.Case) error
}

type createEpaData struct {
	XSRFToken       string
	Success         bool
	Error           sirius.ValidationError
	Case            sirius.Case
	AppointmentType string
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
		}

		if r.Method == http.MethodPost {
			caseAttorneyValue := r.FormValue("caseAttorney")

			epa := sirius.Case{
				EpaDonorSignatureDate:           postFormDateString(r, "epaDonorSignatureDate"),
				EpaDonorNoticeGivenDate:         postFormDateString(r, "epaDonorNoticeGivenDate"),
				DonorHasOtherEpas:               boolPtr(postFormString(r, "donorHasOtherEpas") == "true"),
				OtherEpaInfo:                    postFormString(r, "otherEpaInfo"),
				ReceiptDate:                     postFormDateString(r, "receiptDate"),
				CaseAttorneySingular:            boolPtr(caseAttorneyValue == "singular"),
				CaseAttorneyJointlyAndSeverally: boolPtr(caseAttorneyValue == "jointly-and-severally"),
				CaseAttorneyJointly:             boolPtr(caseAttorneyValue == "jointly"),
				PaymentByCheque:                 boolPtr(r.FormValue("paymentByCheque") == "true"),
				PaymentExemption:                boolPtr(r.FormValue("paymentExemption") == "true"),
				PaymentDate:                     postFormDateString(r, "paymentDate"),
			}
			data.Case = epa
			data.AppointmentType = caseAttorneyValue

			err := client.CreateEpa(ctx, donorID, epa)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
				return tmpl(w, data)
			} else if err != nil {
				return err
			}

			data.Success = true
		}

		return tmpl(w, data)
	}
}
