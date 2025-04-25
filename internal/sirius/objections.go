package sirius

import (
	"fmt"
)

type ObjectionsForCase struct {
	UID           string      `json:"uid"`
	ObjectionList []Objection `json:"objections"`
}

type Objection struct {
	ID            int    `json:"id"`
	Notes         string `json:"notes"`
	ObjectionType string `json:"objectionType"`
	ReceivedDate  string `json:"receivedDate"`
}

func (c *Client) ObjectionsForCase(ctx Context, caseUID string) ([]Objection, error) {
	var caseObjections ObjectionsForCase
	path := fmt.Sprintf("/lpa-api/v1/digital-lpas/%s/objections", caseUID)

	err := c.get(ctx, path, &caseObjections)

	return caseObjections.ObjectionList, err
}
