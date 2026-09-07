package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type CreateTrustCorporationClient interface {
	CreateTrustCorporation(ctx sirius.Context, caseId int, trustCorporation sirius.TrustCorporation) error
}

func CreateTrustCorporation(client CreateTrustCorporationClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		data, err := TrustCorporation(r, "Add a trust corporation")
		if err != nil {
			return err
		}

		if r.Method == http.MethodPost {
			err = client.CreateTrustCorporation(ctx, data.CaseId, data.TrustCorporation)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve

				return tmpl(w, data)
			} else if err != nil {
				return err
			}

			if r.FormValue("add-another") != "" {
				if data.TrustCorporation.IsReplacementAttorney {
					return RedirectError(fmt.Sprintf("/create-replacement-attorney?id=%d&caseId=%d", data.DonorId, data.CaseId))
				} else {
					return RedirectError(fmt.Sprintf("/create-attorney?id=%d&caseId=%d&caseType=lpa", data.DonorId, data.CaseId))
				}
			}

			return RedirectError(fmt.Sprintf("/create-lpa?id=%d&caseId=%d#scroll-to-attorneys-corporation", data.DonorId, data.CaseId))
		}

		return tmpl(w, data)
	}
}
