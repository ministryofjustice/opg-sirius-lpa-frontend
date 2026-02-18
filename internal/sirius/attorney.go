package sirius

import "fmt"

type Attorney struct {
	Salutation          string     `json:"salutation,omitempty"`
	FirstName           string     `json:"firstname,omitempty"`
	MiddleNames         string     `json:"middlenames,omitempty"`
	Surname             string     `json:"surname,omitempty"`
	OtherNames          string     `json:"otherNames,omitempty"`
	DOB                 DateString `json:"dob,omitempty"`
	PhoneNumber         string     `json:"phoneNumber,omitempty"`
	Email               string     `json:"email,omitempty"`
	AddressLine1        string     `json:"addressLine1,omitempty"`
	AddressLine2        string     `json:"addressLine2,omitempty"`
	AddressLine3        string     `json:"addressLine3,omitempty"`
	Town                string     `json:"town,omitempty"`
	Postcode            string     `json:"postcode,omitempty"`
	County              string     `json:"county,omitempty"`
	Country             string     `json:"country,omitempty"`
	DateOfDeath         DateString `json:"dateOfDeath,omitempty"`
	CompanyName         string     `json:"companyName,omitempty"`
	CompanyNumber       string     `json:"companyNumber,omitempty"`
	Parent              string     `json:"parent,omitempty"`
	Occupation          string     `json:"occupation,omitempty"`
	SystemStatus        string     `json:"systemStatus,omitempty"`
	SignatureDates      string     `json:"signatureDates,omitempty"`
	ApplicantStatus     string     `json:"applicantStatus,omitempty"`
	RelationshipToDonor string     `json:"relationshipToDonor,omitempty"`
	IsAirmailRequired   bool       `json:"isAirmailRequired,omitempty"`
}

func (c *Client) CreateAttorney(ctx Context, caseId int, attorney Attorney) error {
	return c.post(ctx, fmt.Sprintf("/lpa-api/v1/epas/%d/attorneys", caseId), attorney, nil)
}
