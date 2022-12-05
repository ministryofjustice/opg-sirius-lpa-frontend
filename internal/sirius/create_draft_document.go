package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) CreateDraftDocument(ctx Context, caseID, correspondentID int, templateID string, inserts []string) error {
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
		return err
	}

	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/lpa-api/v1/lpas/%d/documents/draft", caseID), bytes.NewReader(data))
	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		var v ValidationError
		if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
			return err
		}
		return v
	}

	if resp.StatusCode != http.StatusCreated {
		return newStatusError(resp)
	}

	return nil
}
