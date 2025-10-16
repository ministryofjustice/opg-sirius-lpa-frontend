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
	AuthorisedSignatory              string     `json:"authorisedSignatory"`
	WitnessedByCertificateProviderAt time.Time  `json:"witnessedByCertificateProviderAt"`
	WitnessedByIndependentWitnessAt  *time.Time `json:"witnessedByIndependentWitnessAt"`
	IndependentWitnessName           string     `json:"independentWitnessName"`
	IndependentWitnessAddress        Address    `json:"independentWitnessAddress"`
}

func (c *Client) ChangeDonorDetails(ctx Context, caseUID string, donorDetails ChangeDonorDetails) error {
	return c.put(ctx, fmt.Sprintf("/lpa-api/v1/digital-lpas/%s/change-donor-details", caseUID), donorDetails, nil)
}
