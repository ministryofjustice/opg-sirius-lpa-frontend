package sirius

import (
	"fmt"
	"sort"
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

const (
	TypeDraft   string = "Draft"
	TypePreview string = "Preview"
	TypeSave    string = "Save"
)

func (c *Client) Documents(ctx Context, caseType CaseType, caseId int, docType string) ([]Document, error) {
	var d []Document

	url := fmt.Sprintf("/lpa-api/v1/%s/%d/documents?type[]=%s", caseType+"s", caseId, docType)

	err := c.get(ctx, url, &d)
	if err != nil {
		return nil, err
	}

	sort.Slice(d, func(i, j int) bool {
		return d[i].ID > d[j].ID
	})

	return d, err
}

func (c *Client) DocumentByUUID(ctx Context, uuid string) (Document, error) {
	var d Document
	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/documents/%s", uuid), &d)

	return d, err
}
