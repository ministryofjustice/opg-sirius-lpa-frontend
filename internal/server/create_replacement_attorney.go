package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type CreateReplacementAttorneyClient interface {
	Lpa(ctx sirius.Context, id int) (sirius.Lpa, error)
	CreateReplacementAttorney(ctx sirius.Context, caseId int, attorney sirius.Attorney) (sirius.Attorney, error)
	UpdateReplacementAttorney(ctx sirius.Context, attorneyId int, attorney sirius.Attorney) error
}

type createReplacementAttorneyData struct {
	XSRFToken              string
	Attorney               sirius.Attorney
	Error                  sirius.ValidationError
	DonorId                int
	CaseId                 int
	IsEditing              bool
	Title                  string
	NextAttorneyId         int
	HtmxRedirect           string
	HtmxSwap               string
	IsPartial              bool
	NextTrustCorporationId int
	AppointedAs            string
	IsReplacementAttorney  bool
}

func CreateReplacementAttorney(client CreateReplacementAttorneyClient, tmpl template.Template) Handler {
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

		lpa, err := client.Lpa(ctx, caseId)
		if err != nil {
			return err
		}

		data := createReplacementAttorneyData{
			XSRFToken: ctx.XSRFToken,
			DonorId:   donorId,
			CaseId:    caseId,
			Title:     "Add a replacement attorney",
			IsPartial: ctx.IsPartial,
		}

		if lpa.ReceiptDate == "" {
			data.Error = sirius.ValidationError{
				Field: sirius.FieldErrors{
					"receiptDate": {"receiptDate": "A receipt date must be added to the LPA before a replacement attorney can be added"},
				},
			}

			w.WriteHeader(http.StatusBadRequest)
			return tmpl(w, data)
		}

		var attorneyId int
		attorneyIdStr := r.FormValue("attorneyId")
		isEditing := attorneyIdStr != ""
		if isEditing {
			attorneyId, err = strToIntOrStatusError(attorneyIdStr)
			if err != nil {
				return err
			}

			for _, attorney := range lpa.ReplacementAttorneys {
				if attorney.ID == attorneyId {
					data.Attorney = attorney
					break
				}
			}

			data.Title = "Update replacement attorney details"
			data.IsEditing = true

			data.NextAttorneyId = GetNextAttorneyId(attorneyId, lpa.ReplacementAttorneys)

			if data.NextAttorneyId == 0 && len(lpa.TrustCorporations) > 0 {
				data.NextTrustCorporationId = 0
				for _, trustCorporation := range lpa.TrustCorporations {
					if data.NextTrustCorporationId == 0 || trustCorporation.ID < data.NextTrustCorporationId {
						data.NextTrustCorporationId = trustCorporation.ID
						data.IsReplacementAttorney = trustCorporation.IsReplacementAttorney
						if trustCorporation.IsReplacementAttorney {
							data.AppointedAs = "Replacement attorney"
						} else {
							data.AppointedAs = "Attorney"
						}
					}
				}
			}
		}

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

			if isEditing {
				err = client.UpdateReplacementAttorney(ctx, attorneyId, attorney)
			} else {
				_, err = client.CreateReplacementAttorney(ctx, caseId, attorney)
			}

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve

				return tmpl(w, data)
			} else if err != nil {
				return err
			}

			if r.FormValue("next-trust-corporation") != "" {
				return RedirectError(fmt.Sprintf("/update-trust-corporation?id=%d&caseId=%d&trustCorporationId=%d&replacement=%s", donorId, data.CaseId, data.NextTrustCorporationId, strconv.FormatBool(data.IsReplacementAttorney)))
			}

			if r.FormValue("add-another") != "" {
				if ctx.IsPartial {
					data.HtmxRedirect = fmt.Sprintf("/create-replacement-attorney?id=%d&caseId=%d", donorId, caseId)
					data.HtmxSwap = "innerHTML scroll:.action-panel__content:top"
					return tmpl(w, data)
				}
				return RedirectError(fmt.Sprintf("/create-replacement-attorney?id=%d&caseId=%d", donorId, caseId))
			}

			if r.FormValue("next-attorney") != "" {
				if ctx.IsPartial {
					data.HtmxRedirect = fmt.Sprintf("/create-replacement-attorney?id=%d&caseId=%d&attorneyId=%d", donorId, caseId, data.NextAttorneyId)
					data.HtmxSwap = "innerHTML scroll:.action-panel__content:top"
					return tmpl(w, data)
				}
				return RedirectError(fmt.Sprintf("/create-replacement-attorney?id=%d&caseId=%d&attorneyId=%d", donorId, caseId, data.NextAttorneyId))
			}

			if ctx.IsPartial {
				data.HtmxRedirect = fmt.Sprintf("/create-lpa?id=%d&caseId=%d", donorId, caseId)
				data.HtmxSwap = "innerHTML show:#scroll-to-replacement-attorneys:top"
				return tmpl(w, data)
			}
			return RedirectError(fmt.Sprintf("/create-lpa?id=%d&caseId=%d#scroll-to-replacement-attorneys", donorId, caseId))
		}

		return tmpl(w, data)
	}
}
