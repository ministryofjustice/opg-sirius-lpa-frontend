package sirius

import (
	"fmt"
)

type ChangeAttorneyDetails struct {
	FirstNames  string     `json:"firstNames"`
	LastName    string     `json:"lastName"`
	DateOfBirth DateString `json:"dateOfBirth"`
	Address     Address    `json:"address"`
	Phone       string     `json:"phoneNumber"`
	Email       string     `json:"email"`
	SignedAt    DateString `json:"signedAt"`
}

func (c *Client) ChangeAttorneyDetails(ctx Context, caseUID string, attorneyUID string, attorneyDetails ChangeAttorneyDetails) error {
	return c.put(ctx, fmt.Sprintf("/lpa-api/v1/digital-lpas/%s/attorney/%s/change-details", caseUID, attorneyUID), attorneyDetails, nil)
}
