package sirius

import (
	"fmt"
)

type UpdateAttorneyDetails struct {
	FirstNames  string     `json:"firstNames"`
	LastName    string     `json:"lastName"`
	DateOfBirth DateString `json:"dateOfBirth"`
	Address     Address    `json:"address"`
	Phone       string     `json:"phoneNumber"`
	Email       string     `json:"email"`
	SignedAt    DateString `json:"signedAt"`
}

func (c *Client) UpdateAttorneyDetails(ctx Context, caseUID string, attorneyUID string, attorneyDetails UpdateAttorneyDetails) error {
	path := fmt.Sprintf("/lpa-api/v1/digital-lpas/%s/attorney/%s/change-details", caseUID, attorneyUID)

	return c.put(ctx, path, attorneyDetails, nil)
}
