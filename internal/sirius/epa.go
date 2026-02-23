package sirius

import "fmt"

type Epa struct {
	Id                              int        `json:"id,omitempty"`
	EpaDonorSignatureDate           DateString `json:"epaDonorSignatureDate,omitempty"`
	EpaDonorNoticeGivenDate         DateString `json:"epaDonorNoticeGivenDate,omitempty"`
	DonorHasOtherEpas               bool       `json:"donorHasOtherEpas,omitempty"`
	ReceiptDate                     DateString `json:"receiptDate,omitempty"`
	RegistrationDate                DateString `json:"registrationDate,omitempty"`
	DispatchDate                    DateString `json:"dispatchDate,omitempty"`
	CaseAttorneySingular            bool       `json:"caseAttorneySingular,omitempty"`
	CaseAttorneyJointlyAndSeverally bool       `json:"caseAttorneyJointlyAndSeverally,omitempty"`
	CaseAttorneyJointly             bool       `json:"caseAttorneyJointly,omitempty"`
	AttorneyRelationshipToDonor     string     `json:"attorneyRelationshipToDonor,omitempty"`
	AppointmentType                 string     `json:"appointmentType,omitempty"`
	AttorneyTitle                   string     `json:"attorneyTitle,omitempty"`
	AttorneyFirstName               string     `json:"attorneyFirstName,omitempty"`
	AttorneyMiddleName              string     `json:"attorneyMiddleName,omitempty"`
	AttorneyLastName                string     `json:"attorneyLastName,omitempty"`
	AttorneyDob                     DateString `json:"attorneyDob,omitempty"`
	AttorneyCompanyName             string     `json:"attorneyCompanyName,omitempty"`
	AttorneyAddressLineOne          string     `json:"attorneyAddressLineOne,omitempty"`
	AttorneyAddressLineTwo          string     `json:"attorneyAddressLineTwo,omitempty"`
	AttorneyAddressLineThree        string     `json:"attorneyAddressLineThree,omitempty"`
	AttorneyTown                    string     `json:"attorneyTown,omitempty"`
	AttorneyCountry                 string     `json:"attorneyCountry,omitempty"`
	AttorneyPostcode                string     `json:"attorneyPostcode,omitempty"`
	AttorneyAirmail                 string     `json:"attorneyAirmail,omitempty"`
	AttorneyApplyingToRegEpa        string     `json:"attorneyApplyingToRegEpa,omitempty"`
	AttorneyIsActive                string     `json:"attorneyIsActive,omitempty"`
}

func (c *Client) CreateEpa(ctx Context, donorID int, epa Epa) (int, error) {
	var response Epa
	err := c.post(ctx, fmt.Sprintf("/lpa-api/v1/donors/%d/epas", donorID), epa, &response)
	return response.Id, err
}

func (c *Client) UpdateEpa(ctx Context, caseId int, epa Epa) error {
	return c.post(ctx, fmt.Sprintf("/lpa-api/v1/epas/%d", caseId), epa, nil)
}
