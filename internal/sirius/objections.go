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

type GetObjection struct {
	Objection Objection `json:"objection"`
	LpaUids   []string  `json:"lpaUids"`
}

func (c *Client) ObjectionsForCase(ctx Context, caseUID string) ([]Objection, error) {
	var caseObjections ObjectionsForCase
	path := fmt.Sprintf("/lpa-api/v1/digital-lpas/%s/objections", caseUID)

	err := c.get(ctx, path, &caseObjections)

	return caseObjections.ObjectionList, err
}

func (c *Client) GetObjectionByID(ctx Context, objectionId int) (GetObjection, error) {
	var objection GetObjection
	path := fmt.Sprintf("/api/v1/objections/%d", objectionId)

	err := c.get(ctx, path, &objection)

	return objection, err
}
