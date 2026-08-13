package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type CreateTrustCorporationClient interface {
	CreateTrustCorporation(ctx sirius.Context, caseId int, trustCorporation sirius.TrustCorporation) error
	Lpa(ctx sirius.Context, caseId int) (sirius.Lpa, error)
	UpdateTrustCorporation(ctx sirius.Context, trustCorporationId int, trustCorporation sirius.TrustCorporation) error
}

type createTrustCorporationData struct {
	XSRFToken              string
	IsPartial              bool
	TrustCorporation       sirius.TrustCorporation
	Error                  sirius.ValidationError
	DonorId                int
	CaseId                 int
	IsEditing              bool
	Title                  string
	NextTrustCorporationId int
	HtmxRedirect           string
	HtmxSwap               string
	HtmxPost               string
	AppointedAs            string
}

func CreateTrustCorporation(client CreateTrustCorporationClient, tmpl template.Template) Handler {
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

		isReplacementAttorney := r.FormValue("replacement") == "true"

		data := createTrustCorporationData{
			XSRFToken: ctx.XSRFToken,
			IsPartial: r.Header.Get("HX-Request") == "true",
			DonorId:   donorId,
			CaseId:    caseId,
			Title:     "Add a trust corporation",
			HtmxPost:  fmt.Sprintf("/create-trust-corporation?id=%d&caseId=%d&replacement=%s", donorId, caseId, strconv.FormatBool(isReplacementAttorney)),
			TrustCorporation: sirius.TrustCorporation{
				IsReplacementAttorney: isReplacementAttorney,
				Attorney:              sirius.Attorney{SystemStatus: shared.BoolPtr(true)},
			},
		}

		if isReplacementAttorney {
			data.AppointedAs = "Replacement attorney"
		} else {
			data.AppointedAs = "Attorney"
		}

		lpa, err := client.Lpa(ctx, caseId)
		if err != nil {
			return err
		}

		var trustCorporationId int
		trustCorporationIdStr := r.FormValue("trustCorporationId")
		isEditing := trustCorporationIdStr != ""
		if isEditing {
			trustCorporationId, err = strToIntOrStatusError(trustCorporationIdStr)
			if err != nil {
				return err
			}

			for _, trustCorporation := range lpa.TrustCorporations {
				if trustCorporation.ID == trustCorporationId {
					data.TrustCorporation = trustCorporation
					break
				}
			}

			data.Title = "Update trust corporation details"
			data.IsEditing = true
			data.NextTrustCorporationId = GetNextTrustCorporationId(trustCorporationId, lpa.TrustCorporations)
			data.HtmxPost = fmt.Sprintf("/create-trust-corporation?id=%d&caseId=%d&trustCorporationId=%d&isReplacementAttorney=%s", donorId, caseId, trustCorporationId, strconv.FormatBool(isReplacementAttorney))

			if data.TrustCorporation.IsReplacementAttorney {
				data.AppointedAs = "Replacement attorney"
			} else {
				data.AppointedAs = "Attorney"
			}
		}

		if r.Method == http.MethodPost {
			trustCorporation := sirius.TrustCorporation{
				Attorney: sirius.Attorney{
					Person: sirius.Person{
						CompanyName:       postFormString(r, "companyName"),
						CompanyReference:  postFormString(r, "companyReference"),
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
				},
				IsReplacementAttorney: postFormString(r, "isReplacementAttorney") == "true",
			}

			if trustCorporation.IsReplacementAttorney {
				trustCorporation.TrustCorporationAppointedAs = "Replacement Attorney"
				trustCorporation.SystemStatus = shared.BoolPtr(false)
			} else {
				trustCorporation.TrustCorporationAppointedAs = "Attorney"
				trustCorporation.SystemStatus = shared.BoolPtr(postFormString(r, "isTrustCorporationActive") == "true")
			}

			data.TrustCorporation = trustCorporation

			if isEditing {
				err = client.UpdateTrustCorporation(ctx, trustCorporationId, trustCorporation)
			} else {
				err = client.CreateTrustCorporation(ctx, caseId, trustCorporation)
			}

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve

				return tmpl(w, data)
			} else if err != nil {
				return err
			}

			if r.FormValue("add-another-trust-corporation") != "" {
				if data.IsPartial {
					data.HtmxRedirect = fmt.Sprintf("/create-trust-corportation?id=%d&caseId=%d&replacement=%s", donorId, caseId, strconv.FormatBool(trustCorporation.IsReplacementAttorney))
					data.HtmxSwap = "innerHTML scroll:.action-panel__content:top"
					return tmpl(w, data)
				}
				return RedirectError(fmt.Sprintf("/create-trust-corporation?id=%d&caseId=%d&replacement=%s", donorId, caseId, strconv.FormatBool(trustCorporation.IsReplacementAttorney)))
			}

			if r.FormValue("edit-next-trust-corporation") != "" {
				if data.IsPartial {
					data.HtmxRedirect = fmt.Sprintf("/create-trust-corporation?id=%d&caseId=%d&trustCorporationId=%d&replacement=%s", donorId, caseId, data.NextTrustCorporationId, strconv.FormatBool(trustCorporation.IsReplacementAttorney))
					data.HtmxSwap = "innerHTML scroll:.action-panel__content:top"
					return tmpl(w, data)
				}
				return RedirectError(fmt.Sprintf("/create-trust-corporation?id=%d&caseId=%d&trustCorporationId=%d&replacement=%s", donorId, caseId, data.NextTrustCorporationId, strconv.FormatBool(trustCorporation.IsReplacementAttorney)))
			}

			if data.IsPartial {
				data.HtmxRedirect = fmt.Sprintf("/create-lpa?id=%d&caseId=%d", donorId, caseId)
				data.HtmxSwap = "innerHTML show:#accordion-create-lpa-heading-1:top"
				return tmpl(w, data)
			}

			return RedirectError(fmt.Sprintf("/create-lpa?id=%d&caseId=%d#accordion-create-lpa-heading-1", donorId, caseId))
		}

		return tmpl(w, data)
	}
}

func GetNextTrustCorporationId(id int, trustCorporations []sirius.TrustCorporation) int {
	nextTrustCorporationId := 0
	for _, trustCorporation := range trustCorporations {
		if trustCorporation.ID > id && (nextTrustCorporationId == 0 || trustCorporation.ID < nextTrustCorporationId) {
			nextTrustCorporationId = trustCorporation.ID
		}
	}
	return nextTrustCorporationId
}
