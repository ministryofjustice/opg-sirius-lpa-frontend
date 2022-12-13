package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type DocumentData struct {
	DocumentID int `json:"id"`
}

func (c *Client) CreateDocument(ctx Context, caseID, correspondentID int, templateID string, inserts []string) (DocumentData, error) {
	var d DocumentData

	data, err := json.Marshal(struct {
		TemplateID      string   `json:"templateId"`
		Inserts         []string `json:"inserts"`
		CorrespondentID int      `json:"correspondentId"`
	}{
		TemplateID:      templateID,
		Inserts:         inserts,
		CorrespondentID: correspondentID,
	})

	if err != nil {
		return d, err
	}

	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/lpa-api/v1/lpas/%d/documents/draft", caseID), bytes.NewReader(data))
	if err != nil {
		return d, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return d, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		var v ValidationError
		if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
			return d, err
		}
		return d, v
	}

	if resp.StatusCode != http.StatusCreated {
		return d, newStatusError(resp)
	}

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return d, err
	}

	return d, nil
}
