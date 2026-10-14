package sirius

import (
	"fmt"
)

type UpdateDraft struct {
	FirstNames  string     `json:"firstNames"`
	LastName    string     `json:"lastName"`
	DateOfBirth DateString `json:"dateOfBirth"`
	Address     Address    `json:"address"`
	Phone       string     `json:"phoneNumber"`
	Email       string     `json:"email"`
}

func (c *Client) UpdateDraft(ctx Context, caseUID string, draftDetails UpdateDraft) error {
	return c.put(ctx, fmt.Sprintf("/lpa-api/v1/digital-lpas/%s/update-draft", caseUID), draftDetails, nil)
}
