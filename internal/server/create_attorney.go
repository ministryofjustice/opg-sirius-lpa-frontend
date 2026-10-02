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
}

func CreateAttorney(client CreateAttorneyClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		data, err := Attorney(r, "Add an attorney")
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

		// Default the active status to true for new attorneys
		data.Attorney.SystemStatus = shared.BoolPtr(true)

		if r.Method == http.MethodPost {
			err = client.CreateAttorney(ctx, data.CaseId, data.CaseType, data.Attorney)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
				return tmpl(w, data)
			} else if err != nil {
				return err
			}

			if r.FormValue("add-another") != "" {
				if data.IsPartial {
					data.HtmxRedirect = fmt.Sprintf("/create-attorney?id=%d&caseId=%d&caseType=%s", data.DonorId, data.CaseId, data.CaseType)
					data.HtmxSwap = "innerHTML scroll:.action-panel__content:top"
					return tmpl(w, data)
				}
				return RedirectError(fmt.Sprintf("/create-attorney?id=%d&caseId=%d&caseType=%s", data.DonorId, data.CaseId, data.CaseType))
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
