package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type TrustCorporationData struct {
	AppointedAs            string
	CaseId                 int
	DonorId                int
	Error                  sirius.ValidationError
	HtmxPost               string
	IsEditing              bool
	IsPartial              bool
	NextTrustCorporationId int
	Title                  string
	TrustCorporation       sirius.TrustCorporation
	XSRFToken              string
}

func TrustCorporation(r *http.Request, title string) (TrustCorporationData, error) {
	ctx := getContext(r)
	donorId, err := strToIntOrStatusError(r.FormValue("id"))
	if err != nil {
		return TrustCorporationData{}, err
	}

	caseId, err := strToIntOrStatusError(r.FormValue("caseId"))
	if err != nil {
		return TrustCorporationData{}, err
	}

	isReplacementAttorney := r.FormValue("replacement") == "true"

	data := TrustCorporationData{
		XSRFToken: ctx.XSRFToken,
		IsPartial: ctx.IsPartial,
		DonorId:   donorId,
		CaseId:    caseId,
		Title:     title,
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
	}

	return data, nil
}
