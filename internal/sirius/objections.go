package sirius

import (
	"fmt"
)

type ObjectionForCase struct {
	ID            int                   `json:"id"`
	Notes         string                `json:"notes"`
	ObjectionType string                `json:"objectionType"`
	ReceivedDate  string                `json:"receivedDate"`
	LpaUids       []string              `json:"lpaUids"`
	Resolutions   []ObjectionResolution `json:"objectionLpas"`
}

type Objection struct {
	ObjectionForCase
	Resolutions []ObjectionResolution `json:"resolutions"`
}

type ObjectionResolution struct {
	Uid             string `json:"uid"`
	Resolution      string `json:"resolution"`
	ResolutionNotes string `json:"resolutionNotes"`
	ResolutionDate  string `json:"resolutionDate,omitempty"`
}

type ObjectionRequest struct {
	LpaUids       []string   `json:"lpaUids"`
	ReceivedDate  DateString `json:"receivedDate"`
	ObjectionType string     `json:"objectionType"`
	Notes         string     `json:"notes"`
}

type ResolutionRequest struct {
	Resolution string `json:"resolution"`
	Notes      string `json:"resolutionNotes"`
}

func (c *Client) AddObjection(ctx Context, objectionDetails ObjectionRequest) error {
	return c.post(ctx, "/lpa-api/v1/objections", objectionDetails, nil)
}

func (c *Client) UpdateObjection(ctx Context, objectionId string, objectionDetails ObjectionRequest) error {
	return c.put(ctx, fmt.Sprintf("/lpa-api/v1/objections/%s", objectionId), objectionDetails, nil)
}

func (c *Client) ObjectionsForCase(ctx Context, caseUID string) ([]ObjectionForCase, error) {
	var objectionList []ObjectionForCase
	path := fmt.Sprintf("/lpa-api/v1/digital-lpas/%s/objections", caseUID)

	err := c.get(ctx, path, &objectionList)

	return objectionList, err
}

func (c *Client) GetObjection(ctx Context, Id string) (Objection, error) {
	var objection Objection
	path := fmt.Sprintf("/lpa-api/v1/objections/%s", Id)

	err := c.get(ctx, path, &objection)

	return objection, err
}

func (c *Client) ResolveObjection(ctx Context, objectionId string, lpaUid string, resolutionDetails ResolutionRequest) error {
	return c.put(ctx, fmt.Sprintf("/lpa-api/v1/objections/%s/resolution/%s", objectionId, lpaUid), resolutionDetails, nil)
}
