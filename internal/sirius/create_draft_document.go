package sirius

import (
	"fmt"
)

func (c *Client) CreateDocument(ctx Context, caseID, correspondentID int, templateID string, inserts []string) (Document, error) {
	var d Document

	data := struct {
		TemplateID      string   `json:"templateId"`
		Inserts         []string `json:"inserts"`
		CorrespondentID int      `json:"correspondentId"`
	}{
		TemplateID:      templateID,
		Inserts:         inserts,
		CorrespondentID: correspondentID,
	}

	err := c.post(ctx, fmt.Sprintf("/lpa-api/v1/lpas/%d/documents/draft", caseID), data, &d)
	if err != nil {
		return Document{}, nil
	}

	return d, nil
}
