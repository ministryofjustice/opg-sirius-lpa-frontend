package sirius

import (
	"fmt"
	"time"
)

type ChangeDonorDetails struct {
	FirstNames                       string     `json:"firstNames"`
	LastName                         string     `json:"lastName"`
	OtherNamesKnownBy                string     `json:"otherNamesKnownBy"`
	DateOfBirth                      DateString `json:"dateOfBirth"`
	Address                          Address    `json:"address"`
	Phone                            string     `json:"phoneNumber"`
	Email                            string     `json:"email"`
	LpaSignedOn                      DateString `json:"lpaSignedOn"`
	AuthorisedSignatory              string     `form:"authorisedSignatory"`
	WitnessedByCertificateProviderAt time.Time  `form:"witnessedByCertificateProviderAt"`
	WitnessedByIndependentWitnessAt  *time.Time `form:"witnessedByIndependentWitnessAt"`
	IndependentWitnessName           string     `form:"independentWitnessName"`
	IndependentWitnessAddress        Address    `form:"independentWitnessAddress"`
}

func (c *Client) ChangeDonorDetails(ctx Context, caseUID string, donorDetails ChangeDonorDetails) error {
	return c.put(ctx, fmt.Sprintf("/lpa-api/v1/digital-lpas/%s/change-donor-details", caseUID), donorDetails, nil)
}
