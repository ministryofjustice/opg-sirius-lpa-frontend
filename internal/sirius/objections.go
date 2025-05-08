package sirius

import (
	"fmt"
)

type Objection struct {
	ID            int      `json:"id"`
	Notes         string   `json:"notes"`
	ObjectionType string   `json:"objectionType"`
	ReceivedDate  string   `json:"receivedDate"`
	LpaUids       []string `json:"lpaUids"`
}

func (c *Client) ObjectionsForCase(ctx Context, caseUID string) ([]Objection, error) {
	var objectionList []Objection
	path := fmt.Sprintf("/lpa-api/v1/digital-lpas/%s/objections", caseUID)

	err := c.get(ctx, path, &objectionList)

	return objectionList, err
}
