package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
	AppointedAs      string
	CaseId           int
	CaseType         string
	DonorId          int
	Error            sirius.ValidationError
	HtmxPost         string
	IsEditing        bool
	IsPartial        bool
	NextPersonId     int
	NextPersonType   string
	Title            string
	TrustCorporation sirius.TrustCorporation
	XSRFToken        string
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

		lpa, err := client.Lpa(ctx, caseId)
		if err != nil {
			return err
		}

		isReplacementAttorney := r.FormValue("replacement") == "true"

		data := createTrustCorporationData{
			XSRFToken: ctx.XSRFToken,
			IsPartial: r.Header.Get("HX-Request") == "true",
			DonorId:   donorId,
			CaseId:    caseId,
			CaseType:  strings.ToLower(lpa.CaseType),
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
			data.HtmxPost = fmt.Sprintf("/create-trust-corporation?id=%d&caseId=%d&trustCorporationId=%d&replacement=%s", donorId, caseId, trustCorporationId, strconv.FormatBool(isReplacementAttorney))

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
					CompanyNumber: postFormString(r, "companyNumber"),
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

			if r.FormValue("add-another") != "" {
				if trustCorporation.IsReplacementAttorney {
					return RedirectError(fmt.Sprintf("/create-replacement-attorney?id=%d&caseId=%d", donorId, caseId))
				} else {
					return RedirectError(fmt.Sprintf("/create-attorney?id=%d&caseId=%d&caseType=%s", donorId, caseId, data.CaseType))
				}
			}

			if r.FormValue("next-attorney") != "" {
				if trustCorporation.IsReplacementAttorney {
					data.NextPersonId, data.NextPersonType = GetIdForNextAttorney(trustCorporationId, true, lpa.TrustCorporations, lpa.ReplacementAttorneys)
					if data.NextPersonType == "Trust Corporation" {
						return RedirectError(fmt.Sprintf("/create-trust-corporation?id=%d&caseId=%d&trustCorporationId=%d&replacement=%s", donorId, caseId, data.NextPersonId, strconv.FormatBool(trustCorporation.IsReplacementAttorney)))
					} else {
						return RedirectError(fmt.Sprintf("/create-replacement-attorney?id=%d&caseId=%d&attorneyId=%d", donorId, caseId, data.NextPersonId))
					}
				} else {
					data.NextPersonId, data.NextPersonType = GetIdForNextAttorney(trustCorporationId, false, lpa.TrustCorporations, lpa.Attorneys)
					if data.NextPersonType == "Trust Corporation" {
						return RedirectError(fmt.Sprintf("/create-trust-corporation?id=%d&caseId=%d&trustCorporationId=%d&replacement=%s", donorId, caseId, data.NextPersonId, strconv.FormatBool(trustCorporation.IsReplacementAttorney)))
					}
					return RedirectError(fmt.Sprintf("/create-attorney?id=%d&caseId=%d&caseType=%s&attorneyId=%d", donorId, caseId, data.CaseType, data.NextPersonId))
				}
			}

			return RedirectError(fmt.Sprintf("/create-lpa?id=%d&caseId=%d#accordion-create-lpa-heading-1", donorId, caseId))
		}

		return tmpl(w, data)
	}
}

func GetIdForNextAttorney(id int, isReplacementAttorney bool, trustCorporations []sirius.TrustCorporation, attorneys []sirius.Attorney) (int, string) {
	nextID := 0
	personType := ""

	for _, trustCorporation := range trustCorporations {
		if trustCorporation.ID > id && trustCorporation.IsReplacementAttorney == isReplacementAttorney && (nextID == 0 || trustCorporation.ID < nextID) {
			nextID = trustCorporation.ID
			personType = trustCorporation.PersonType
		}
	}

	for _, attorney := range attorneys {
		if attorney.ID > id && (nextID == 0 || attorney.ID < nextID) {
			nextID = attorney.ID
			personType = attorney.PersonType
		}
	}
	return nextID, personType
}
