package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type CreateAttorneyClient interface {
	Epa(ctx sirius.Context, id int) (sirius.Epa, error)
	Lpa(ctx sirius.Context, id int) (sirius.Lpa, error)
	CreateAttorney(ctx sirius.Context, caseId int, caseTyp string, attorney sirius.Attorney) error
	RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error)
	UpdateAttorney(ctx sirius.Context, attorneyId int, attorney sirius.Attorney) error
}

type createAttorneyData struct {
	XSRFToken            string
	IsPartial            bool
	Attorney             sirius.Attorney
	Error                sirius.ValidationError
	RelationshipToDonors []sirius.RefDataItem
	DonorId              int
	CaseId               int
	CaseType             string
	CaseSubType          string
	IsEditing            bool
	Title                string
	NextAttorneyId       int
	HtmxRedirect         string
	HtmxSwap             string
}

func CreateAttorney(client CreateAttorneyClient, tmpl template.Template) Handler {
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

		caseType := r.FormValue("caseType")

		data := createAttorneyData{
			XSRFToken: ctx.XSRFToken,
			IsPartial: ctx.IsPartial,
			DonorId:   donorId,
			CaseId:    caseId,
			CaseType:  caseType,
			Title:     "Add an attorney",
		}

		var lpa sirius.Lpa
		if caseType == "lpa" {
			lpa, err = client.Lpa(ctx, caseId)
			if err != nil {
				return err
			}

			data.CaseSubType = lpa.SubType
		}

		data.RelationshipToDonors, err = client.RefDataByCategory(ctx, sirius.RelationshipToDonorCategory)
		if err != nil {
			return err
		}

		// Default the active status to true for new attorneys
		data.Attorney.SystemStatus = shared.BoolPtr(true)

		var attorneyId int
		attorneyIdStr := r.FormValue("attorneyId")
		isEditing := attorneyIdStr != ""
		if isEditing {
			attorneyId, err = strToIntOrStatusError(attorneyIdStr)
			if err != nil {
				return err
			}

			var attorneys []sirius.Attorney
			if caseType == "epa" {
				epa, err := client.Epa(ctx, caseId)
				if err != nil {
					return err
				}

				attorneys = epa.Attorneys
			} else {
				attorneys = lpa.Attorneys
			}

			for _, attorney := range attorneys {
				if attorney.ID == attorneyId {
					data.Attorney = attorney
					break
				}
			}

			data.Title = "Update attorney details"
			data.IsEditing = true
			data.NextAttorneyId = GetNextAttorneyId(attorneyId, attorneys)
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
				SystemStatus: shared.BoolPtr(postFormString(r, "isAttorneyActive") == "true"),
			}

			if caseType == "epa" {
				attorney.CompanyName = postFormString(r, "companyName")
				attorney.RelationshipToDonor = postFormString(r, "relationshipToDonor")
			}
			data.Attorney = attorney

			if isEditing {
				err = client.UpdateAttorney(ctx, attorneyId, attorney)
			} else {
				err = client.CreateAttorney(ctx, caseId, caseType, attorney)
			}

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
				return tmpl(w, data)
			} else if err != nil {
				return err
			}

			if r.FormValue("add-another") != "" {
				if data.IsPartial {
					data.HtmxRedirect = fmt.Sprintf("/create-attorney?id=%d&caseId=%d&caseType=%s", donorId, caseId, caseType)
					data.HtmxSwap = "innerHTML scroll:.action-panel__content:top"
					return tmpl(w, data)
				}
				return RedirectError(fmt.Sprintf("/create-attorney?id=%d&caseId=%d&caseType=%s", donorId, caseId, caseType))
			}

			if r.FormValue("next-attorney") != "" {
				if data.IsPartial {
					data.HtmxRedirect = fmt.Sprintf("/create-attorney?id=%d&caseId=%d&caseType=%s&attorneyId=%d", donorId, caseId, caseType, data.NextAttorneyId)
					data.HtmxSwap = "innerHTML scroll:.action-panel__content:top"
					return tmpl(w, data)
				}
				return RedirectError(fmt.Sprintf("/create-attorney?id=%d&caseId=%d&caseType=%s&attorneyId=%d", donorId, caseId, caseType, data.NextAttorneyId))
			}

			if data.IsPartial {
				data.HtmxRedirect = fmt.Sprintf("/create-%s?id=%d&caseId=%d", caseType, donorId, caseId)
				data.HtmxSwap = "innerHTML show:#accordion-create-epa-heading-3:top"
				return tmpl(w, data)
			}

			if caseType == "epa" {
				return RedirectError(fmt.Sprintf("/create-epa?id=%d&caseId=%d#accordion-create-epa-heading-3", donorId, caseId))
			}
			return RedirectError(fmt.Sprintf("/create-lpa?id=%d&caseId=%d#scroll-to-attorneys", donorId, caseId))
		}

		return tmpl(w, data)
	}
}

func GetNextAttorneyId(id int, attorneys []sirius.Attorney) int {
	nextAttorneyId := 0
	for _, attorney := range attorneys {
		if attorney.ID > id && (nextAttorneyId == 0 || attorney.ID < nextAttorneyId) {
			nextAttorneyId = attorney.ID
		}
	}
	return nextAttorneyId
}
