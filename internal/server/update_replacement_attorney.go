package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type UpdateReplacementAttorneyClient interface {
	Lpa(ctx sirius.Context, id int) (sirius.Lpa, error)
	UpdateReplacementAttorney(ctx sirius.Context, attorneyId int, attorney sirius.Attorney) error
}

func UpdateReplacementAttorney(client UpdateReplacementAttorneyClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		donorId, err := strToIntOrStatusError(r.FormValue("id"))
		if err != nil {
			return err
		}

		caseId, err := strToIntOrStatusError(r.FormValue("caseId"))
		if err != nil {
			return err
		}

		attorneyId, err := strToIntOrStatusError(r.FormValue("attorneyId"))
		if err != nil {
			return err
		}

		lpa, err := client.Lpa(ctx, caseId)
		if err != nil {
			return err
		}

		data := replacementAttorneyData{
			XSRFToken:   ctx.XSRFToken,
			FlowQuery:   getFlowQuery(r),
			DonorId:     donorId,
			CaseId:      caseId,
			CaseSubType: lpa.SubType,
			Title:       "Update replacement attorney details",
			IsPartial:   ctx.IsPartial,
			IsEditing:   true,
		}

		var nextPersonType string

		for _, attorney := range lpa.ReplacementAttorneys {
			if attorney.ID == attorneyId {
				data.Attorney = attorney
				break
			}
		}

		data.NextAttorneyId, nextPersonType = GetIdForNextAttorney(attorneyId, true, lpa.TrustCorporations, lpa.ReplacementAttorneys)

		if r.Method == http.MethodPost {
			attorney := sirius.Attorney{
				Person: sirius.Person{
					Salutation:        postFormString(r, "salutation"),
					Firstname:         postFormString(r, "firstname"),
					Middlenames:       postFormString(r, "middlenames"),
					Surname:           postFormString(r, "surname"),
					DateOfBirth:       postFormDateString(r, "dob"),
					PhoneNumber:       postFormString(r, "phoneNumber"),
					Email:             postFormString(r, "email"),
					AddressLine1:      postFormString(r, "addressLine1"),
					AddressLine2:      postFormString(r, "addressLine2"),
					AddressLine3:      postFormString(r, "addressLine3"),
					Town:              postFormString(r, "town"),
					County:            postFormString(r, "county"),
					Country:           postFormString(r, "country"),
					Postcode:          postFormString(r, "postcode"),
					IsAirmailRequired: postFormString(r, "isAirmailRequired") == "true",
				},
			}
			data.Attorney = attorney

			err = client.UpdateReplacementAttorney(ctx, attorneyId, attorney)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve

				return tmpl(w, data)
			} else if err != nil {
				return err
			}

			if r.FormValue("update-next-attorney") != "" {
				redirect := fmt.Sprintf("/update-replacement-attorney?id=%d&caseId=%d&attorneyId=%d%s", donorId, caseId, data.NextAttorneyId, data.FlowQuery)
				if nextPersonType == "Trust Corporation" {
					redirect = fmt.Sprintf("/update-trust-corporation?id=%d&caseId=%d&trustCorporationId=%d&replacement=true%s", donorId, caseId, data.NextAttorneyId, data.FlowQuery)
				}

				if ctx.IsPartial {
					data.HtmxRedirect = redirect
					data.HtmxSwap = "innerHTML scroll:.action-panel__content:top"
					return tmpl(w, data)
				}
				return RedirectError(redirect)
			}

			if ctx.IsPartial {
				data.HtmxRedirect = fmt.Sprintf("/create-lpa?id=%d&caseId=%d%s", donorId, caseId, data.FlowQuery)
				data.HtmxSwap = "innerHTML show:#scroll-to-replacement-attorneys:top"
				return tmpl(w, data)
			}
			return RedirectError(fmt.Sprintf("/create-lpa?id=%d&caseId=%d%s#scroll-to-replacement-attorneys", donorId, caseId, data.FlowQuery))
		}

		return tmpl(w, data)
	}
}
