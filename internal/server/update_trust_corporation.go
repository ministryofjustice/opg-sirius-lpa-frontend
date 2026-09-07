package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type UpdateTrustCorporationClient interface {
	Lpa(ctx sirius.Context, caseId int) (sirius.Lpa, error)
	UpdateTrustCorporation(ctx sirius.Context, trustCorporationId int, trustCorporation sirius.TrustCorporation) error
}

func UpdateTrustCorporation(client UpdateTrustCorporationClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		data, err := TrustCorporation(r, "Update trust corporation details")
		if err != nil {
			return err
		}

		var trustCorporationId int
		trustCorporationIdStr := r.FormValue("trustCorporationId")
		trustCorporationId, err = strToIntOrStatusError(trustCorporationIdStr)
		if err != nil {
			return err
		}

		lpa, err := client.Lpa(ctx, data.CaseId)
		if err != nil {
			return err
		}

		for _, trustCorporation := range lpa.TrustCorporations {
			if trustCorporation.ID == trustCorporationId {
				if r.Method == http.MethodPost {
					data.TrustCorporation.ID = trustCorporation.ID
				} else {
					data.TrustCorporation = trustCorporation
				}
				break
			}
		}

		data.IsEditing = true
		data.NextTrustCorporationId = GetNextTrustCorporationId(trustCorporationId, data.TrustCorporation.IsReplacementAttorney, lpa.TrustCorporations)
		data.HtmxPost = fmt.Sprintf("/update-trust-corporation?id=%d&caseId=%d&trustCorporationId=%d&replacement=%s", data.DonorId, data.CaseId, trustCorporationId, strconv.FormatBool(data.TrustCorporation.IsReplacementAttorney))

		if data.TrustCorporation.IsReplacementAttorney {
			data.AppointedAs = "Replacement attorney"
		} else {
			data.AppointedAs = "Attorney"
		}

		if r.Method == http.MethodPost {
			err = client.UpdateTrustCorporation(ctx, trustCorporationId, data.TrustCorporation)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve

				return tmpl(w, data)
			} else if err != nil {
				return err
			}

			if r.FormValue("editNextTrustCorporation") != "" {
				return RedirectError(fmt.Sprintf("/update-trust-corporation?id=%d&caseId=%d&trustCorporationId=%d&replacement=%s", data.DonorId, data.CaseId, data.NextTrustCorporationId, strconv.FormatBool(data.TrustCorporation.IsReplacementAttorney)))
			}

			return RedirectError(fmt.Sprintf("/create-lpa?id=%d&caseId=%d#scroll-to-attorneys-corporation", data.DonorId, data.CaseId))
		}

		return tmpl(w, data)
	}
}

func GetNextTrustCorporationId(id int, isReplacementAttorney bool, trustCorporations []sirius.TrustCorporation) int {
	nextTrustCorporationId := 0
	for _, trustCorporation := range trustCorporations {
		if trustCorporation.ID > id && (nextTrustCorporationId == 0 || trustCorporation.ID < nextTrustCorporationId) && trustCorporation.IsReplacementAttorney == isReplacementAttorney {
			nextTrustCorporationId = trustCorporation.ID
		}
	}
	return nextTrustCorporationId
}
