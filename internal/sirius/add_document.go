package sirius

import (
	"fmt"
)

type addDocumentRequestData struct {
	CaseId          int    `json:"caseId"`
	CorrespondentID int    `json:"correspondentId"`
	Type            string `json:"type"`
	FileName        string `json:"filename"`
	SystemType      string `json:"systemType"`
	Content         string `json:"content"`
}

func (c *Client) AddDocument(ctx Context, caseID int, document Document, docType string) (Document, error) {
	data := addDocumentRequestData{
		CaseId:          caseID,
		CorrespondentID: document.Correspondent.ID,
		Type:            docType,
		FileName:        document.FileName,
		SystemType:      document.SystemType,
		Content:         document.Content,
	}

	var d Document

	err := c.post(ctx, fmt.Sprintf("/lpa-api/v1/lpas/%d/documents", caseID), data, &d)
	if err != nil {
		return Document{}, err
	}

	return d, nil
}
