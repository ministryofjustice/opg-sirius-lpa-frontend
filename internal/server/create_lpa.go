package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type CreateLpaClient interface {
	Person(sirius.Context, int) (sirius.Person, error)
	Lpa(sirius.Context, int) (sirius.Lpa, error)
	CreateLpa(ctx sirius.Context, donorID int, lpa sirius.Lpa) (sirius.Lpa, error)
	UpdateLpa(ctx sirius.Context, caseID int, lpa sirius.Lpa) error
	UpdateAttorney(ctx sirius.Context, attorneyId int, attorney sirius.Attorney) error
	UpdateReplacementAttorney(ctx sirius.Context, attorneyId int, attorney sirius.Attorney) error
	UpdateTrustCorporation(ctx sirius.Context, trustCorporationId int, trustCorporation sirius.TrustCorporation) error
	GetUserPermissions(sirius.Context) (sirius.Permissions, error)
}

type createLpaData struct {
	AllowNewNotifiedPerson               bool
	AppointmentType                      string
	AttorneyTrustCorporations            []sirius.TrustCorporation
	CanEditReceiptDate                   bool
	CaseId                               int
	DonorId                              int
	DonorName                            string
	Error                                sirius.ValidationError
	HtmxRedirect                         string
	HtmxSwap                             string
	IsUpdate                             bool
	Lpa                                  sirius.Lpa
	ReplacementAttorneyTrustCorporations []sirius.TrustCorporation
	Success                              bool
	SuccessMessage                       string
	Title                                string
	XSRFToken                            string
}

func CreateLpa(client CreateLpaClient, tmpl template.Template, partialTmpl template.Template) Handler {
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

		userPermissions, err := client.GetUserPermissions(ctx)
		if err != nil {
			return err
		}

		data := createLpaData{
			XSRFToken:              ctx.XSRFToken,
			DonorId:                donorID,
			DonorName:              donor.Firstname + " " + donor.Surname,
			Title:                  "Create an LPA",
			AllowNewNotifiedPerson: allowNewNotifiedPerson(0),
			CanEditReceiptDate:     userPermissions.Includes("v1-lpas-edit-dates", "PUT"),
		}

		caseIdStr := r.FormValue("caseId")
		isEditing := caseIdStr != ""
		if isEditing {
			data.CaseId, err = strToIntOrStatusError(caseIdStr)
			if err != nil {
				return err
			}

			data.Lpa, err = client.Lpa(ctx, data.CaseId)
			if err != nil {
				return err
			}

			data.AppointmentType = appointmentTypeFromCase(data.Lpa.Case)
			data.AllowNewNotifiedPerson = allowNewNotifiedPerson(len(data.Lpa.NotifiedPersons))
			data.Title = "Edit LPA"
			data.IsUpdate = true

			for _, trustCorporation := range data.Lpa.TrustCorporations {
				if trustCorporation.IsReplacementAttorney {
					data.ReplacementAttorneyTrustCorporations = append(data.ReplacementAttorneyTrustCorporations, trustCorporation)
				} else {
					data.AttorneyTrustCorporations = append(data.AttorneyTrustCorporations, trustCorporation)
				}
			}
		}

		if r.Method == http.MethodPost {
			caseAttorneyValue := r.FormValue("caseAttorney")

			lpa := sirius.Lpa{
				OnlineLpaId:            postFormString(r, "onlineLpaId"),
				AttorneyActDecisions:   postFormString(r, "attorneyActDecisions"),
				ApplicantType:          postFormString(r, "applicantType"),
				ApplicantSignatureDate: postFormDateString(r, "applicantSignatureDate"),
				CardPaymentContact:     postFormString(r, "cardPaymentContact"),
				Case: sirius.Case{
					SubType:                                   postFormString(r, "caseSubtype"),
					ApplicationType:                           postFormString(r, "applicationType"),
					CaseAttorneySingular:                      shared.BoolPtr(caseAttorneyValue == "singular"),
					CaseAttorneyJointly:                       shared.BoolPtr(caseAttorneyValue == "jointly"),
					CaseAttorneyJointlyAndSeverally:           shared.BoolPtr(caseAttorneyValue == "jointly-and-severally"),
					CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(caseAttorneyValue == "jointly-and-jointly-and-severally"),
					LifeSustainingTreatment:                   postFormString(r, "lifeSustainingTreatment"),
					LifeSustainingTreatmentSignatureDateA:     postFormDateString(r, "lifeSustainingTreatmentSignatureDate"),
					LpaDonorSignatureDate:                     postFormDateString(r, "lpaDonorSignatureDate"),
				},
			}

			for _, idStr := range r.PostForm["applicantIds"] {
				if id, err := strconv.Atoi(idStr); err == nil {
					lpa.ApplicantIds = append(lpa.ApplicantIds, id)
				}
			}

			lpa.LifeSustainingTreatmentSignedAndWitnessed = shared.BoolPtr(postFormString(r, "lifeSustainingTreatmentSignedAndWitnessed") == "true")

			reducedFeeSelected := postFormCheckboxChecked(r, "applicationFee", "reducedFee")
			reducedFeeType := postFormString(r, "reducedFeeType")
			lpa.PaymentByCheque = shared.BoolPtr(postFormCheckboxChecked(r, "applicationFee", "cheque"))
			lpa.PaymentByDebitCreditCard = shared.BoolPtr(postFormCheckboxChecked(r, "applicationFee", "card"))
			lpa.PaymentExemption = shared.BoolPtr(reducedFeeSelected && reducedFeeType == "exemption")
			lpa.PaymentRemission = shared.BoolPtr(reducedFeeSelected && reducedFeeType == "remission")
			lpa.RepeatApplication = shared.BoolPtr(postFormCheckboxChecked(r, "applicationFee", "repeatApplication"))

			if !*lpa.PaymentByDebitCreditCard {
				lpa.CardPaymentContact = ""
			}

			lpa.AnyOtherInfo = shared.BoolPtr(postFormString(r, "anyOtherInfo") == "true")
			lpa.AdditionalInfo = postFormString(r, "additionalInfo")
			if !*lpa.AnyOtherInfo {
				lpa.AdditionalInfo = ""
			}

			if lpa.SubType != "pfa" {
				lpa.AttorneyActDecisions = ""
			}
			if lpa.SubType != "hw" {
				lpa.LifeSustainingTreatment = ""
				lpa.LifeSustainingTreatmentSignatureDateA = ""
				lpa.LifeSustainingTreatmentSignedAndWitnessed = nil
			}

			preferencesNone := postFormCheckboxChecked(r, "preferencesAndInstructions", "none")
			lpa.ApplicationHasGuidance = shared.BoolPtr(!preferencesNone && postFormCheckboxChecked(r, "preferencesAndInstructions", "guidance"))
			lpa.ApplicationHasRestrictions = shared.BoolPtr(!preferencesNone && postFormCheckboxChecked(r, "preferencesAndInstructions", "restrictions"))

			if data.CanEditReceiptDate {
				lpa.ReceiptDate = postFormDateString(r, "receiptDate")
			} else if isEditing {
				lpa.ReceiptDate = data.Lpa.ReceiptDate
			}

			if lpa.ApplicationType != "Online" {
				lpa.OnlineLpaId = ""
			}

			data.AppointmentType = caseAttorneyValue

			if isEditing {
				err = client.UpdateLpa(ctx, data.CaseId, lpa)
				if err == nil {
					data.Lpa, _ = client.Lpa(ctx, data.CaseId)
				}
			} else {
				var createdLpa sirius.Lpa
				createdLpa, err = client.CreateLpa(ctx, donorID, lpa)
				if err == nil {
					data.Lpa = createdLpa
					data.CaseId = createdLpa.ID
				}
			}

			if err == nil {
				for _, attorney := range data.Lpa.Attorneys {
					formValue := postFormString(r, fmt.Sprintf("lpaPartCSignatureDate-%d", attorney.ID))
					if formValue != "" && formValue != string(attorney.LpaPartCSignatureDate) {
						attorney.LpaPartCSignatureDate = sirius.DateString(formValue)
						err = client.UpdateAttorney(ctx, attorney.ID, attorney)
					}
				}
				for _, replacementAttorney := range data.Lpa.ReplacementAttorneys {
					formValue := postFormString(r, fmt.Sprintf("lpaPartCSignatureDate-%d", replacementAttorney.ID))
					if formValue != "" && formValue != string(replacementAttorney.LpaPartCSignatureDate) {
						replacementAttorney.LpaPartCSignatureDate = sirius.DateString(formValue)
						err = client.UpdateReplacementAttorney(ctx, replacementAttorney.ID, replacementAttorney)
					}
				}
				for _, trustCorporation := range data.Lpa.TrustCorporations {
					formValue := postFormString(r, fmt.Sprintf("lpaPartCSignatureDate-%d", trustCorporation.ID))
					if formValue != "" && formValue != string(trustCorporation.LpaPartCSignatureDate) {
						trustCorporation.LpaPartCSignatureDate = sirius.DateString(formValue)
						err = client.UpdateTrustCorporation(ctx, trustCorporation.ID, trustCorporation)
					}
				}
			}

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
				lpa.Attorneys = data.Lpa.Attorneys
				lpa.ReplacementAttorneys = data.Lpa.ReplacementAttorneys
				data.Lpa = lpa

				if r.Header.Get("HX-Request") == "true" {
					return partialTmpl(w, data)
				}
				return tmpl(w, data)
			} else if err != nil {
				return err
			}

			if r.FormValue("addAttorney") != "" {
				return RedirectError(fmt.Sprintf("/create-attorney?id=%d&caseId=%d&caseType=lpa", donorID, data.CaseId))
			}
			if r.FormValue("addCertificateProvider") != "" {
				return RedirectError(fmt.Sprintf("/create-certificate-provider?id=%d&caseId=%d", donorID, data.CaseId))
			}
			if r.FormValue("addCorrespondent") != "" {
				return RedirectError(fmt.Sprintf("/select-or-create-correspondent?id=%d&caseId=%d&caseType=lpa", donorID, data.CaseId))
			}
			if r.FormValue("updateCorrespondent") != "" {
				return RedirectError(fmt.Sprintf("/create-correspondent?id=%d&caseId=%d&caseType=lpa", donorID, data.CaseId))
			}
			if r.FormValue("addNotifiedPerson") != "" {
				return RedirectError(fmt.Sprintf("/create-notified-person?id=%d&caseId=%d", donorID, data.CaseId))
			} else if updateNotifiedPerson := r.FormValue("updateNotifiedPerson"); updateNotifiedPerson != "" {
				notifiedPersonID, err := strToIntOrStatusError(updateNotifiedPerson)
				if err != nil {
					return err
				}
				return RedirectError(fmt.Sprintf("/create-notified-person?id=%d&caseId=%d&notifiedPersonId=%d", donorID, data.CaseId, notifiedPersonID))
			}
			if trustCorporationIdStr := r.FormValue("updateTrustCorporationAttorney"); trustCorporationIdStr != "" {
				trustCorporationId, err := strToIntOrStatusError(trustCorporationIdStr)
				if err != nil {
					return err
				}
				return RedirectError(fmt.Sprintf("/create-trust-corporation?id=%d&caseId=%d&trustCorporationId=%d&replacement=false", donorID, data.CaseId, trustCorporationId))
			}
			if trustCorporationIdStr := r.FormValue("updateTrustCorporationReplacementAttorney"); trustCorporationIdStr != "" {
				trustCorporationId, err := strToIntOrStatusError(trustCorporationIdStr)
				if err != nil {
					return err
				}
				return RedirectError(fmt.Sprintf("/create-trust-corporation?id=%d&caseId=%d&trustCorporationId=%d&replacement=true", donorID, data.CaseId, trustCorporationId))
			}

			data.Success = true
			if isEditing {
				data.SuccessMessage = "You have successfully updated an LPA."
			} else {
				data.SuccessMessage = "You have successfully created an LPA."
			}
			data.AllowNewNotifiedPerson = allowNewNotifiedPerson(len(data.Lpa.NotifiedPersons))

		}

		if r.FormValue("addReplacementAttorney") != "" {
			if r.Header.Get("HX-Request") == "true" {
				data.HtmxRedirect = fmt.Sprintf("/create-replacement-attorney?id=%d&caseId=%d", donorID, data.CaseId)
				data.HtmxSwap = "innerHTML"
				return partialTmpl(w, data)
			}
			return RedirectError(fmt.Sprintf("/create-replacement-attorney?id=%d&caseId=%d", donorID, data.CaseId))
		}

		if r.Header.Get("HX-Request") == "true" {
			return partialTmpl(w, data)
		}

		return tmpl(w, data)
	}
}

func appointmentTypeFromCase(c sirius.Case) string {
	switch {
	case c.CaseAttorneySingular != nil && *c.CaseAttorneySingular:
		return "singular"
	case c.CaseAttorneyJointly != nil && *c.CaseAttorneyJointly:
		return "jointly"
	case c.CaseAttorneyJointlyAndSeverally != nil && *c.CaseAttorneyJointlyAndSeverally:
		return "jointly-and-severally"
	case c.CaseAttorneyJointlyAndJointlyAndSeverally != nil && *c.CaseAttorneyJointlyAndJointlyAndSeverally:
		return "jointly-and-jointly-and-severally"
	default:
		return ""
	}
}
