package sirius

import (
	"fmt"
)

type Document struct {
	ID                  int    `json:"id,omitempty"`
	UUID                string `json:"uuid,omitempty"`
	Type                string `json:"type"`
	FriendlyDescription string `json:"friendlyDescription"`
	CreatedDate         string `json:"createdDate"`
	Direction           string `json:"direction"`
	MimeType            string `json:"mimeType"`
	SystemType          string `json:"systemType"`
	FileName            string `json:"fileName,omitempty"`
	Content             string `json:"content,omitempty"`
	Correspondent       Person `json:"correspondent"`
	ChildCount          int    `json:"childCount"`
	CaseItems           []Case `json:"caseItems,omitempty"`
}

func (c *Client) Documents(ctx Context, caseType CaseType, caseId int) ([]Document, error) {
	var d []Document

	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/%s/%d/documents", caseType+"s", caseId), &d)
	if err != nil {
		return nil, err
	}

	return d, err
}

func (c *Client) DocumentByUUID(ctx Context, uuid string) (Document, error) {
	var d Document
	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/documents/%s", uuid), &d)

	return d, err
}
