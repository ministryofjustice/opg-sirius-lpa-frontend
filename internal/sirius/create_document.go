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

func (c *Client) CreateDocument(ctx Context, caseID, correspondentID int, templateID string, inserts []string) (*DocumentData, error) {
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
		return nil, err
	}

	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/lpa-api/v1/lpas/%d/documents/draft", caseID), bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		var v ValidationError
		if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
			return nil, err
		}
		return nil, v
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, newStatusError(resp)
	}

	var v DocumentData
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	return &v, nil
}
