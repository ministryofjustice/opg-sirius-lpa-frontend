package server

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type UpdateLpaClient interface {
	Person(sirius.Context, int) (sirius.Person, error)
	Lpa(sirius.Context, int) (sirius.Lpa, error)
	UpdateLpa(ctx sirius.Context, caseID int, lpa sirius.Lpa) error
	UpdateAttorney(ctx sirius.Context, attorneyId int, attorney sirius.Attorney) error
	UpdateReplacementAttorney(ctx sirius.Context, attorneyId int, attorney sirius.Attorney) error
	UpdateTrustCorporation(ctx sirius.Context, trustCorporationId int, trustCorporation sirius.TrustCorporation) error
}

func UpdateLpa(client UpdateLpaClient, tmpl template.Template) Handler {
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

		caseId, err := strToIntOrStatusError(r.FormValue("caseId"))
		if err != nil {
			return err
		}

		lpa, err := client.Lpa(ctx, caseId)
		if err != nil {
			return err
		}

		isCreating := r.FormValue("flow") == "create"

		data := lpaData{
			AllowNewNotifiedPerson: allowNewNotifiedPerson(len(lpa.NotifiedPersons)),
			ApplicantIds:           lpa.GetApplicantIds(),
			AppointmentType:        appointmentTypeFromCase(lpa.Case),
			CaseId:                 caseId,
			DonorId:                donorID,
			DonorName:              donor.Firstname + " " + donor.Surname,
			IsCreating:             isCreating,
			IsPartial:              ctx.IsPartial,
			Lpa:                    lpa,
			Title:                  "Edit LPA",
			XSRFToken:              ctx.XSRFToken,
		}

		if isCreating {
			data.FlowQuery = "&flow=create"
			data.Title = "Create an LPA"
		}

		for _, trustCorporation := range data.Lpa.TrustCorporations {
			if trustCorporation.IsReplacementAttorney {
				data.ReplacementAttorneyTrustCorporations = append(data.ReplacementAttorneyTrustCorporations, trustCorporation)
			} else {
				data.AttorneyTrustCorporations = append(data.AttorneyTrustCorporations, trustCorporation)
			}
		}

		data.AttorneyApplicants = applicants(lpa.Attorneys, data.AttorneyTrustCorporations)
		data.ReplacementAttorneys = applicants(lpa.ReplacementAttorneys, data.ReplacementAttorneyTrustCorporations)

		if r.Method == http.MethodPost {
			caseAttorneyValue := r.FormValue("caseAttorney")
			isSaveAndExit := postFormString(r, "saveAndExit") == "saveAndExit"
			isActorAction := isLpaActorAction(r)
			preferencesNone := postFormCheckboxChecked(r, "preferencesAndInstructions", "none")

			lpa := sirius.Lpa{
				AdditionalInfo:                   postFormString(r, "additionalInfo"),
				AnyOtherInfo:                     shared.BoolPtr(postFormString(r, "anyOtherInfo") == "true"),
				ApplicantSignatureDate:           postFormDateString(r, "applicantSignatureDate"),
				ApplicantType:                    postFormString(r, "applicantType"),
				ApplicationHasContinuationSheets: shared.BoolPtr(postFormCheckboxChecked(r, "continuationSheets", "continuationSheets")),
				ApplicationHasGuidance:           shared.BoolPtr(!preferencesNone && postFormCheckboxChecked(r, "preferencesAndInstructions", "guidance")),
				ApplicationHasRestrictions:       shared.BoolPtr(!preferencesNone && postFormCheckboxChecked(r, "preferencesAndInstructions", "restrictions")),
				AttorneyActDecisions:             postFormString(r, "attorneyActDecisions"),
				CardPaymentContact:               postFormString(r, "cardPaymentContact"),
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
					SubType:                                   postFormString(r, "caseSubtype"),
				},
				CertificateProviderSignatureDate: postFormDateString(r, "certificateProviderSignatureDate"),
				OnlineLpaId:                      postFormString(r, "onlineLpaId"),
				PaymentByDebitCreditCard:         shared.BoolPtr(postFormCheckboxChecked(r, "applicationFee", "card")),
				PaymentRemission:                 shared.BoolPtr(postFormCheckboxChecked(r, "applicationFee", "reducedFee")),
				RepeatApplication:                shared.BoolPtr(postFormCheckboxChecked(r, "applicationFee", "repeatApplication")),
				RepeatApplicationReference:       postFormString(r, "repeatApplicationReference"),
			}

			// Omit unanswered appointment fields during creation or intermediate actor saves.
			if (isCreating || (isActorAction && !isSaveAndExit)) && caseAttorneyValue == "" {
				lpa.CaseAttorneySingular = nil
				lpa.CaseAttorneyJointly = nil
				lpa.CaseAttorneyJointlyAndSeverally = nil
				lpa.CaseAttorneyJointlyAndJointlyAndSeverally = nil
			}

			hasAttorneyIdBeenSelected := len(r.PostForm["applicantIds"]) > 0

			switch lpa.ApplicantType {
			case "donor":
				lpa.ApplicantIds = append(lpa.ApplicantIds, donorID)
			case "attorney":
				if !hasAttorneyIdBeenSelected {
					data.Error.Field = make(sirius.FieldErrors)
					data.Error.Field["applicantType"] = map[string]string{
						"": "Select an applicant",
					}
				} else {
					for _, idStr := range r.PostForm["applicantIds"] {
						if id, err := strconv.Atoi(idStr); err == nil {
							lpa.ApplicantIds = append(lpa.ApplicantIds, id)
						}
					}
				}
			default:
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

			if !*lpa.RepeatApplication {
				lpa.RepeatApplicationReference = ""
			}

			if !*lpa.AnyOtherInfo {
				lpa.AdditionalInfo = ""
			}

			data.AppointmentType = caseAttorneyValue

			if !hasAttorneyIdBeenSelected && lpa.ApplicantType == "attorney" {
				lpa.Attorneys = data.Lpa.Attorneys
				lpa.ReplacementAttorneys = data.Lpa.ReplacementAttorneys
				data.ApplicantIds = data.Lpa.GetApplicantIds()
				data.Lpa = lpa

				w.WriteHeader(http.StatusBadRequest)
				return tmpl(w, data)
			}

			err = client.UpdateLpa(ctx, data.CaseId, lpa)
			if err == nil {
				data.Lpa, _ = client.Lpa(ctx, data.CaseId)
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
			if updateAttorney := r.FormValue("updateAttorney"); updateAttorney != "" {
				attorneyID, err := strToIntOrStatusError(updateAttorney)
				if err != nil {
					return err
				}
				return RedirectError(fmt.Sprintf("/update-attorney?id=%d&caseId=%d&caseType=lpa&attorneyId=%d%s", donorID, data.CaseId, attorneyID, data.FlowQuery))
			}
			if updateCertificateProvider := r.FormValue("updateCertificateProvider"); updateCertificateProvider != "" {
				certificateProviderID, err := strToIntOrStatusError(updateCertificateProvider)
				if err != nil {
					return err
				}
				return RedirectError(fmt.Sprintf("/update-certificate-provider?id=%d&caseId=%d&personId=%d%s", donorID, data.CaseId, certificateProviderID, data.FlowQuery))
			}
			if r.FormValue("updateCorrespondent") != "" {
				return RedirectError(fmt.Sprintf("/edit-correspondent?id=%d&caseId=%d&caseType=lpa%s", donorID, data.CaseId, data.FlowQuery))
			}
			if updateNotifiedPerson := r.FormValue("updateNotifiedPerson"); updateNotifiedPerson != "" {
				notifiedPersonID, err := strToIntOrStatusError(updateNotifiedPerson)
				if err != nil {
					return err
				}
				return RedirectError(fmt.Sprintf("/update-notified-person?id=%d&caseId=%d&notifiedPersonId=%d%s", donorID, data.CaseId, notifiedPersonID, data.FlowQuery))
			}
			if updateReplacementAttorney := r.FormValue("updateReplacementAttorney"); updateReplacementAttorney != "" {
				attorneyID, err := strToIntOrStatusError(updateReplacementAttorney)
				if err != nil {
					return err
				}
				return RedirectError(fmt.Sprintf("/update-replacement-attorney?id=%d&caseId=%d&attorneyId=%d%s", donorID, data.CaseId, attorneyID, data.FlowQuery))
			}
			if trustCorporationIdStr := r.FormValue("updateTrustCorporationAttorney"); trustCorporationIdStr != "" {
				trustCorporationId, err := strToIntOrStatusError(trustCorporationIdStr)
				if err != nil {
					return err
				}
				return RedirectError(fmt.Sprintf("/update-trust-corporation?id=%d&caseId=%d&trustCorporationId=%d&replacement=false%s", donorID, data.CaseId, trustCorporationId, data.FlowQuery))
			}
			if trustCorporationIdStr := r.FormValue("updateTrustCorporationReplacementAttorney"); trustCorporationIdStr != "" {
				trustCorporationId, err := strToIntOrStatusError(trustCorporationIdStr)
				if err != nil {
					return err
				}
				return RedirectError(fmt.Sprintf("/update-trust-corporation?id=%d&caseId=%d&trustCorporationId=%d&replacement=true%s", donorID, data.CaseId, trustCorporationId, data.FlowQuery))
			}

			data.Success = true
			if !isCreating {
				data.SuccessMessage = "You have successfully updated an LPA."
			} else {
				data.SuccessMessage = "You have successfully created an LPA."
			}
		}

		return tmpl(w, data)
	}
}

func getFlowQuery(r *http.Request) string {
	if r.FormValue("flow") == "create" {
		return "&flow=create"
	}
	return ""
}

func isLpaActorAction(r *http.Request) bool {
	for _, name := range []string{
		"addAttorney", "updateAttorney",
		"addReplacementAttorney", "updateReplacementAttorney",
		"updateTrustCorporationAttorney",
		"updateTrustCorporationReplacementAttorney",
		"addCertificateProvider", "updateCertificateProvider",
		"addNotifiedPerson", "updateNotifiedPerson",
		"addCorrespondent", "updateCorrespondent",
	} {
		if postFormString(r, name) != "" {
			return true
		}
	}
	return false
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

func applicants(attorneys []sirius.Attorney, trustCorporations []sirius.TrustCorporation) []sirius.Attorney {
	applicants := make([]sirius.Attorney, 0, len(attorneys)+len(trustCorporations))

	applicants = append(applicants, attorneys...)

	for _, tc := range trustCorporations {
		applicants = append(applicants, tc.Attorney)
	}

	sort.Slice(applicants, func(i, j int) bool {
		return applicants[i].ID < applicants[j].ID
	})

	return applicants
}
