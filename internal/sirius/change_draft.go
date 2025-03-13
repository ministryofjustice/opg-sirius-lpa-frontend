package sirius

import (
	"fmt"
)

type ChangeDraft struct {
	FirstNames  string     `json:"firstNames"`
	LastName    string     `json:"lastName"`
	DateOfBirth DateString `json:"dateOfBirth"`
	Address     Address    `json:"address"`
	Phone       string     `json:"phoneNumber"`
	Email       string     `json:"email"`
}

func (c *Client) ChangeDraft(ctx Context, caseUID string, draftDetails ChangeDraft) error {
	return c.put(ctx, fmt.Sprintf("/lpa-api/v1/digital-lpas/%s/change-draft", caseUID), draftDetails, nil)
}
