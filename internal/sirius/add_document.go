package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) AddDocument(ctx Context, caseID int, document Document, docType string) error {
	data, err := json.Marshal(struct {
		CaseId          int    `json:"caseId"`
		CorrespondentID int    `json:"correspondentId"`
		Type            string `json:"type"`
		FileName        string `json:"filename"`
		SystemType      string `json:"systemType"`
		Content         string `json:"content"`
	}{
		CaseId:          caseID,
		CorrespondentID: document.Correspondent.ID,
		Type:            docType,
		FileName:        document.FileName,
		SystemType:      document.SystemType,
		Content:         document.Content,
	})
	if err != nil {
		return err
	}

	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/lpa-api/v1/lpas/%d/documents", caseID), bytes.NewReader(data))
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
