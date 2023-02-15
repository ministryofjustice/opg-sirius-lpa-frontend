package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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
	data, err := json.Marshal(addDocumentRequestData{
		CaseId:          caseID,
		CorrespondentID: document.Correspondent.ID,
		Type:            docType,
		FileName:        document.FileName,
		SystemType:      document.SystemType,
		Content:         document.Content,
	})
	if err != nil {
		return Document{}, err
	}

	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/lpa-api/v1/lpas/%d/documents", caseID), bytes.NewReader(data))
	if err != nil {
		return Document{}, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return Document{}, err
	}
	defer resp.Body.Close() //#nosec G307 false positive

	if resp.StatusCode == http.StatusBadRequest {
		var v ValidationError
		if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
			return Document{}, err
		}
		return Document{}, v
	}

	if resp.StatusCode != http.StatusCreated {
		return Document{}, newStatusError(resp)
	}

	var d Document
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return Document{}, err
	}
	return d, nil
}
