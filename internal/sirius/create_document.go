package sirius

import (
	"fmt"
)

type createDocumentRequestData struct {
	CaseId            int      `json:"caseId"`
	CorrespondentID   int      `json:"correspondentId"`
	Type              string   `json:"type"`
	FileName          string   `json:"filename"`
	SystemType        string   `json:"systemType"`
	Content           string   `json:"content"`
	AppendLpaSections []string `json:"appendLpaSections"`
	SendViaNotify     bool     `json:"sendViaNotify"`
}

func (c *Client) CreateDocument(ctx Context, caseID int, document Document, docType string, blankSections []string, sendViaNotify bool) (Document, error) {
	data := createDocumentRequestData{
		CaseId:            caseID,
		CorrespondentID:   document.Correspondent.ID,
		Type:              docType,
		FileName:          document.FileName,
		SystemType:        document.SystemType,
		Content:           document.Content,
		AppendLpaSections: blankSections,
		SendViaNotify:     sendViaNotify,
	}

	var d Document

	err := c.post(ctx, fmt.Sprintf("/lpa-api/v1/lpas/%d/documents", caseID), data, &d)
	if err != nil {
		return Document{}, err
	}

	return d, nil
}
