package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type CreateReplacementAttorneyClient interface {
	Lpa(ctx sirius.Context, id int) (sirius.Lpa, error)
	CreateReplacementAttorney(ctx sirius.Context, caseId int, attorney sirius.Attorney) (sirius.Attorney, error)
	UpdateReplacementAttorney(ctx sirius.Context, attorneyId int, attorney sirius.Attorney) error
}

type createReplacementAttorneyData struct {
	XSRFToken      string
	Attorney       sirius.Attorney
	Error          sirius.ValidationError
	DonorId        int
	CaseId         int
	IsEditing      bool
	Title          string
	NextAttorneyId int
	HtmxRedirect   string
	HtmxSwap       string
}

func CreateReplacementAttorney(client CreateReplacementAttorneyClient, tmpl template.Template, partialTmpl template.Template) Handler {
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
		}

		if lpa.ReceiptDate == "" {
			data.Error = sirius.ValidationError{
				Field: sirius.FieldErrors{
					"receiptDate": {"receiptDate": "A receipt date must be added to the LPA before a replacement attorney can be added"},
				},
			}

			w.WriteHeader(http.StatusBadRequest)
			if r.Header.Get("HX-Request") == "true" {
				return partialTmpl(w, data)
			}
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

				if r.Header.Get("HX-Request") == "true" {
					return partialTmpl(w, data)
				}
				return tmpl(w, data)
			} else if err != nil {
				return err
			}

			if r.FormValue("add-another") != "" {
				if r.Header.Get("HX-Request") == "true" {
					data.HtmxRedirect = fmt.Sprintf("/create-replacement-attorney?id=%d&caseId=%d", donorId, caseId)
					data.HtmxSwap = "innerHTML scroll:.action-panel__content:top"
					return partialTmpl(w, data)
				}
				return RedirectError(fmt.Sprintf("/create-replacement-attorney?id=%d&caseId=%d", donorId, caseId))
			}

			if r.FormValue("next-attorney") != "" {
				if r.Header.Get("HX-Request") == "true" {
					data.HtmxRedirect = fmt.Sprintf("/create-replacement-attorney?id=%d&caseId=%d&attorneyId=%d", donorId, caseId, data.NextAttorneyId)
					data.HtmxSwap = "innerHTML scroll:.action-panel__content:top"
					return partialTmpl(w, data)
				}
				return RedirectError(fmt.Sprintf("/create-replacement-attorney?id=%d&caseId=%d&attorneyId=%d", donorId, caseId, data.NextAttorneyId))
			}

			if r.Header.Get("HX-Request") == "true" {
				data.HtmxRedirect = fmt.Sprintf("/create-lpa?id=%d&caseId=%d", donorId, caseId)
				data.HtmxSwap = "innerHTML show:#scroll-to-replacement-attorneys:top"
				return partialTmpl(w, data)
			}
			return RedirectError(fmt.Sprintf("/create-lpa?id=%d&caseId=%d#scroll-to-replacement-attorneys", donorId, caseId))
		}

		if r.Header.Get("HX-Request") == "true" {
			return partialTmpl(w, data)
		}

		return tmpl(w, data)
	}
}
