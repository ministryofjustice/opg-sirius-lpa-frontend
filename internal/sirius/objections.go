package sirius

import (
	"fmt"
)

type ObjectionsForCase struct {
	UID           string      `json:"uid"`
	ObjectionList []Objection `json:"objections"`
}

type Objection struct {
	ID            int      `json:"id"`
	Notes         string   `json:"notes"`
	ObjectionType string   `json:"objectionType"`
	ReceivedDate  string   `json:"receivedDate"`
	LpaUids       []string `json:"lpaUids"`
}

type ObjectionRequest struct {
	LpaUids       []string   `json:"lpaUids"`
	ReceivedDate  DateString `json:"receivedDate"`
	ObjectionType string     `json:"objectionType"`
	Notes         string     `json:"notes"`
}

func (c *Client) AddObjection(ctx Context, objectionDetails ObjectionRequest) error {
	return c.post(ctx, "/lpa-api/v1/objections", objectionDetails, nil)
}

func (c *Client) UpdateObjection(ctx Context, objectionId string, objectionDetails ObjectionRequest) error {
	return c.put(ctx, fmt.Sprintf("/lpa-api/v1/objections/%s", objectionId), objectionDetails, nil)
}

func (c *Client) ObjectionsForCase(ctx Context, caseUID string) ([]Objection, error) {
	var caseObjections ObjectionsForCase
	path := fmt.Sprintf("/lpa-api/v1/digital-lpas/%s/objections", caseUID)

	err := c.get(ctx, path, &caseObjections)

	return caseObjections.ObjectionList, err
}

func (c *Client) GetObjection(ctx Context, Id string) (Objection, error) {
	var objection Objection
	path := fmt.Sprintf("/lpa-api/v1/objections/%s", Id)

	err := c.get(ctx, path, &objection)

	return objection, err
}
