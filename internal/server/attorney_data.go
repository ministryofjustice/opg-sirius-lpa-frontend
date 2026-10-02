package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type AttorneyData struct {
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

func Attorney(r *http.Request, title string) (AttorneyData, error) {
	ctx := getContext(r)

	donorId, err := strToIntOrStatusError(r.FormValue("id"))
	if err != nil {
		return AttorneyData{}, err
	}

	caseId, err := strToIntOrStatusError(r.FormValue("caseId"))
	if err != nil {
		return AttorneyData{}, err
	}

	caseType := r.FormValue("caseType")

	data := AttorneyData{
		XSRFToken: ctx.XSRFToken,
		IsPartial: ctx.IsPartial,
		DonorId:   donorId,
		CaseId:    caseId,
		CaseType:  caseType,
		Title:     title,
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
	}

	return data, nil
}
