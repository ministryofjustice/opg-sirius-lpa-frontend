package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type UpdateAttorneyClient interface {
	Epa(ctx sirius.Context, id int) (sirius.Epa, error)
	Lpa(ctx sirius.Context, id int) (sirius.Lpa, error)
	RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error)
	UpdateAttorney(ctx sirius.Context, attorneyId int, attorney sirius.Attorney) error
}

func UpdateAttorney(client UpdateAttorneyClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		data, err := Attorney(r, "Update attorney details")
		if err != nil {
			return err
		}

		var lpa sirius.Lpa
		if data.CaseType == "lpa" {
			lpa, err = client.Lpa(ctx, data.CaseId)
			if err != nil {
				return err
			}

			data.CaseSubType = lpa.SubType
		}

		data.RelationshipToDonors, err = client.RefDataByCategory(ctx, sirius.RelationshipToDonorCategory)
		if err != nil {
			return err
		}

		var nextPersonType string
		var attorneyId int
		attorneyIdStr := r.FormValue("attorneyId")
		attorneyId, err = strToIntOrStatusError(attorneyIdStr)
		if err != nil {
			return err
		}

		var attorneys []sirius.Attorney
		if data.CaseType == "epa" {
			epa, err := client.Epa(ctx, data.CaseId)
			if err != nil {
				return err
			}

			attorneys = epa.Attorneys
		} else {
			attorneys = lpa.Attorneys
		}
		var existingAttorney sirius.Attorney
		for _, attorney := range attorneys {
			if attorney.ID == attorneyId {
				existingAttorney = attorney
				break
			}
		}

		// Only overwrite for GET requests
		if r.Method == http.MethodGet {
			data.Attorney = existingAttorney
		}

		data.Title = "Update attorney details"
		data.IsEditing = true
		data.NextAttorneyId, nextPersonType = GetIdForNextAttorney(attorneyId, false, lpa.TrustCorporations, lpa.Attorneys)

		if r.Method == http.MethodPost {
			err = client.UpdateAttorney(ctx, attorneyId, data.Attorney)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
				return tmpl(w, data)
			} else if err != nil {
				return err
			}

			if r.FormValue("update-next-attorney") != "" {
				redirect := fmt.Sprintf("/update-attorney?id=%d&caseId=%d&caseType=%s&attorneyId=%d", data.DonorId, data.CaseId, data.CaseType, data.NextAttorneyId)
				if nextPersonType == "Trust Corporation" {
					redirect = fmt.Sprintf("/create-trust-corporation?id=%d&caseId=%d&trustCorporationId=%d&replacement=false", data.DonorId, data.CaseId, data.NextAttorneyId)
				}

				if data.IsPartial {
					data.HtmxRedirect = redirect
					data.HtmxSwap = "innerHTML scroll:.action-panel__content:top"
					return tmpl(w, data)
				}
				return RedirectError(redirect)
			}

			if data.IsPartial {
				data.HtmxRedirect = fmt.Sprintf("/create-%s?id=%d&caseId=%d", data.CaseType, data.DonorId, data.CaseId)
				data.HtmxSwap = "innerHTML show:#accordion-create-epa-heading-3:top"
				return tmpl(w, data)
			}

			if data.CaseType == "epa" {
				return RedirectError(fmt.Sprintf("/create-epa?id=%d&caseId=%d#accordion-create-epa-heading-3", data.DonorId, data.CaseId))
			}
			return RedirectError(fmt.Sprintf("/create-lpa?id=%d&caseId=%d#scroll-to-attorneys", data.DonorId, data.CaseId))
		}

		return tmpl(w, data)
	}
}
