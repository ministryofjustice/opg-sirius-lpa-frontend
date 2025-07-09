package sirius

import (
	"fmt"
	"sort"
	"strings"
)

type Document struct {
	ID                  int    `json:"id,omitempty"`
	UUID                string `json:"uuid,omitempty"`
	Type                string `json:"type"`
	FriendlyDescription string `json:"friendlyDescription"`
	CreatedDate         string `json:"createdDate"`
	ReceivedDateTime    string `json:"receivedDateTime"`
	Direction           string `json:"direction"`
	MimeType            string `json:"mimeType"`
	SystemType          string `json:"systemType"`
	SubType             string `json:"subType"`
	FileName            string `json:"fileName,omitempty"`
	Content             string `json:"content,omitempty"`
	Correspondent       Person `json:"correspondent"`
	ChildCount          int    `json:"childCount"`
	CaseItems           []Case `json:"caseItems,omitempty"`
}

func (d *Document) IsViewable() bool {
	if d.SubType == "Reduced fee request evidence" && d.Direction == "Incoming" {
		return d.ReceivedDateTime != ""
	}
	return true
}

const (
	TypeDraft   string = "Draft"
	TypePreview string = "Preview"
	TypeSave    string = "Save"
)

func (c *Client) Documents(ctx Context, caseType CaseType, caseId int, docTypes []string, notDocTypes []string) ([]Document, error) {
	var d []Document

	if caseType == CaseTypeDigitalLpa {
		caseType = CaseTypeLpa
	}

	query := ""

	for _, docType := range docTypes {
		query = query + "&type[]=" + docType
	}

	for _, docType := range notDocTypes {
		query = query + "&type[-][]=" + docType
	}

	url := fmt.Sprintf("/lpa-api/v1/%s/%d/documents?%s", caseType+"s", caseId, strings.Trim(query, "&"))

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
