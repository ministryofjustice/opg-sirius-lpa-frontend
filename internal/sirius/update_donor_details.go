package sirius

import (
	"fmt"
	"time"
)

type UpdateDonorDetails struct {
	FirstNames                       string     `json:"firstNames"`
	LastName                         string     `json:"lastName"`
	OtherNamesKnownBy                string     `json:"otherNamesKnownBy"`
	DateOfBirth                      DateString `json:"dateOfBirth"`
	Address                          Address    `json:"address"`
	Phone                            string     `json:"phoneNumber"`
	Email                            string     `json:"email"`
	LpaSignedOn                      DateString `json:"lpaSignedOn"`
	AuthorisedSignatory              string     `json:"authorisedSignatory"`
	WitnessedByCertificateProviderAt *time.Time `json:"witnessedByCertificateProviderAt"`
	WitnessedByIndependentWitnessAt  *time.Time `json:"witnessedByIndependentWitnessAt"`
	IndependentWitnessName           string     `json:"independentWitnessName"`
	IndependentWitnessAddress        Address    `json:"independentWitnessAddress"`
}

func (c *Client) UpdateDonorDetails(ctx Context, caseUID string, donorDetails UpdateDonorDetails) error {
	return c.put(ctx, fmt.Sprintf("/lpa-api/v1/digital-lpas/%s/update-donor-details", caseUID), donorDetails, nil)
}
