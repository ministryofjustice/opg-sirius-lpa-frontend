package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type CreateLpaClient interface {
	Person(sirius.Context, int) (sirius.Person, error)
	Lpa(sirius.Context, int) (sirius.Lpa, error)
	CreateLpa(ctx sirius.Context, donorID int, lpa sirius.Lpa) (sirius.Lpa, error)
}

type lpaData struct {
	AllowNewNotifiedPerson               bool
	AppointmentType                      string
	ApplicantIds                         []int
	AttorneyTrustCorporations            []sirius.TrustCorporation
	CaseId                               int
	DonorId                              int
	DonorName                            string
	Error                                sirius.ValidationError
	HtmxRedirect                         string
	HtmxSwap                             string
	IsCreating                           bool
	IsUpdate                             bool
	FlowQuery                            string
	Lpa                                  sirius.Lpa
	ReplacementAttorneyTrustCorporations []sirius.TrustCorporation
	Success                              bool
	SuccessMessage                       string
	Title                                string
	XSRFToken                            string
	IsPartial                            bool
	AttorneyApplicants                   []sirius.Attorney
	ReplacementAttorneys                 []sirius.Attorney
}

func CreateLpa(client CreateLpaClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		donorID, err := strToIntOrStatusError(r.FormValue("id"))
		if err != nil {
			return err
		}

		donor, err := client.Person(ctx, donorID)
		if err != nil {
			return err
		}

		data := lpaData{
			AllowNewNotifiedPerson: allowNewNotifiedPerson(0),
			DonorId:                donorID,
			DonorName:              donor.Firstname + " " + donor.Surname,
			FlowQuery:              "&flow=create",
			IsCreating:             true,
			IsPartial:              ctx.IsPartial,
			Title:                  "Create an LPA",
			XSRFToken:              ctx.XSRFToken,
		}

		if r.Method == http.MethodPost {
			caseAttorneyValue := r.FormValue("caseAttorney")
			preferencesNone := postFormCheckboxChecked(r, "preferencesAndInstructions", "none")

			lpa := sirius.Lpa{
				AdditionalInfo:                   postFormString(r, "additionalInfo"),
				ApplicantSignatureDate:           postFormDateString(r, "applicantSignatureDate"),
				ApplicantType:                    postFormString(r, "applicantType"),
				ApplicationHasContinuationSheets: shared.BoolPtr(postFormCheckboxChecked(r, "continuationSheets", "continuationSheets")),
				ApplicationHasGuidance:           shared.BoolPtr(!preferencesNone && postFormCheckboxChecked(r, "preferencesAndInstructions", "guidance")),
				ApplicationHasRestrictions:       shared.BoolPtr(!preferencesNone && postFormCheckboxChecked(r, "preferencesAndInstructions", "restrictions")),
				AnyOtherInfo:                     shared.BoolPtr(postFormString(r, "anyOtherInfo") == "true"),
				AttorneyActDecisions:             postFormString(r, "attorneyActDecisions"),
				CardPaymentContact:               postFormString(r, "cardPaymentContact"),
				CertificateProviderSignatureDate: postFormDateString(r, "certificateProviderSignatureDate"),
				Case: sirius.Case{
					ApplicationType:                           postFormString(r, "applicationType"),
					CaseAttorneyJointly:                       shared.BoolPtr(caseAttorneyValue == "jointly"),
					CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(caseAttorneyValue == "jointly-and-jointly-and-severally"),
					CaseAttorneyJointlyAndSeverally:           shared.BoolPtr(caseAttorneyValue == "jointly-and-severally"),
					CaseAttorneySingular:                      shared.BoolPtr(caseAttorneyValue == "singular"),
					LifeSustainingTreatment:                   postFormString(r, "lifeSustainingTreatment"),
					LifeSustainingTreatmentSignatureDateA:     postFormDateString(r, "lifeSustainingTreatmentSignatureDate"),
					LpaDonorSignatureDate:                     postFormDateString(r, "lpaDonorSignatureDate"),
					PaymentByCheque:                           shared.BoolPtr(postFormCheckboxChecked(r, "applicationFee", "cheque")),
					ReceiptDate:                               postFormDateString(r, "receiptDate"),
					SubType:                                   postFormString(r, "caseSubtype"),
				},
				OnlineLpaId:              postFormString(r, "onlineLpaId"),
				PaymentByDebitCreditCard: shared.BoolPtr(postFormCheckboxChecked(r, "applicationFee", "card")),
				PaymentRemission:         shared.BoolPtr(postFormCheckboxChecked(r, "applicationFee", "reducedFee")),
				RepeatApplication:        shared.BoolPtr(postFormCheckboxChecked(r, "applicationFee", "repeatApplication")),
			}

			// Omit unanswered appointment fields during creation or intermediate actor saves.
			if caseAttorneyValue == "" {
				lpa.CaseAttorneySingular = nil
				lpa.CaseAttorneyJointly = nil
				lpa.CaseAttorneyJointlyAndSeverally = nil
				lpa.CaseAttorneyJointlyAndJointlyAndSeverally = nil
			}

			if lpa.ApplicantType == "donor" {
				lpa.ApplicantIds = append(lpa.ApplicantIds, donorID)
			} else {
				lpa.ApplicantIds = nil
			}

			if lpa.ApplicationType != "Online" {
				lpa.OnlineLpaId = ""
			}

			if lpa.SubType != "pfa" {
				lpa.AttorneyActDecisions = ""
			}
			if lpa.SubType != "hw" {
				lpa.LifeSustainingTreatment = ""
				lpa.LifeSustainingTreatmentSignatureDateA = ""
			}

			if !*lpa.PaymentByDebitCreditCard {
				lpa.CardPaymentContact = ""
			}

			if !*lpa.AnyOtherInfo {
				lpa.AdditionalInfo = ""
			}

			var createdLpa sirius.Lpa
			createdLpa, err = client.CreateLpa(ctx, donorID, lpa)
			if err == nil {
				data.Lpa = createdLpa
				data.CaseId = createdLpa.ID
			}

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
				data.AppointmentType = caseAttorneyValue
				data.ApplicantIds = lpa.ApplicantIds
				data.Lpa = lpa

				return tmpl(w, data)
			} else if err != nil {
				return err
			}

			if r.FormValue("addAttorney") != "" {
				return RedirectError(fmt.Sprintf("/create-attorney?id=%d&caseId=%d&caseType=lpa%s", donorID, data.CaseId, data.FlowQuery))
			}
			if r.FormValue("addCertificateProvider") != "" {
				return RedirectError(fmt.Sprintf("/create-certificate-provider?id=%d&caseId=%d%s", donorID, data.CaseId, data.FlowQuery))
			}
			if r.FormValue("addCorrespondent") != "" {
				return RedirectError(fmt.Sprintf("/select-or-create-correspondent?id=%d&caseId=%d&caseType=lpa%s", donorID, data.CaseId, data.FlowQuery))
			}
			if r.FormValue("addNotifiedPerson") != "" {
				return RedirectError(fmt.Sprintf("/create-notified-person?id=%d&caseId=%d%s", donorID, data.CaseId, data.FlowQuery))
			}
			if r.FormValue("addReplacementAttorney") != "" {
				return RedirectError(fmt.Sprintf("/create-replacement-attorney?id=%d&caseId=%d%s", donorID, data.CaseId, data.FlowQuery))
			}
			data.Success = true
			data.SuccessMessage = "You have successfully created an LPA."
		}

		return tmpl(w, data)
	}
}
