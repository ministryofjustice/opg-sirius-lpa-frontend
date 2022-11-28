package sirius

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Document struct {
	ID              int    `json:"id,omitempty"`
	CorrespondentID int    `json:"correspondentID"`
	Type            string `json:"type"`
	SystemType      string `json:"systemType"`
	FileName        string `json:"fileName,omitempty"`
	Content         string `json:"content"`
}

type documentTemplateApiResponse map[string]json.RawMessage
type insertApiResponse map[string]json.RawMessage

type UniversalTemplateData struct {
	Location        string `json:"location"`
	OnScreenSummary string `json:"onScreenSummary"`
}

type DocumentTemplateApiData struct {
	Inserts insertApiResponse `json:"inserts"`
	UniversalTemplateData
}

type Insert struct {
	Key      string // all/pending/imperfect/perfect/withdrawn
	InsertId string
	UniversalTemplateData
}

type DocumentTemplateData struct {
	Inserts    []Insert
	TemplateId string
	UniversalTemplateData
}

func (d documentTemplateApiResponse) toDocumentData() ([]DocumentTemplateData, error) {
	var s []DocumentTemplateData

	for k, v := range d {
		var asDocumentTemplateData DocumentTemplateData
		asDocumentTemplateData.TemplateId = k

		var asDocumentTemplateApiData DocumentTemplateApiData
		if err := json.Unmarshal(v, &asDocumentTemplateApiData); err == nil {
			inserts, err := asDocumentTemplateApiData.Inserts.toInsertData()
			if err != nil {
				return nil, err
			}
			asDocumentTemplateData.Location = asDocumentTemplateApiData.Location
			asDocumentTemplateData.OnScreenSummary = asDocumentTemplateApiData.OnScreenSummary
			asDocumentTemplateData.Inserts = inserts
			s = append(s, asDocumentTemplateData)
			continue
		} else {
			// handle when inserts = []
			var universalTemplateData UniversalTemplateData
			if err := json.Unmarshal(v, &universalTemplateData); err == nil {
				asDocumentTemplateData.Location = universalTemplateData.Location
				asDocumentTemplateData.OnScreenSummary = universalTemplateData.OnScreenSummary
				s = append(s, asDocumentTemplateData)
				continue
			}
		}

		return nil, errors.New("could not format document template data")
	}
	return s, nil
}

func (i insertApiResponse) toInsertData() ([]Insert, error) {
	var inserts []Insert

	for k, v := range i {
		var insert Insert
		insert.Key = k
		var asInsertKeyApiResponse map[string]UniversalTemplateData
		if err := json.Unmarshal(v, &asInsertKeyApiResponse); err == nil {
			for insertId, insertData := range asInsertKeyApiResponse {
				insert.InsertId = insertId
				insert.UniversalTemplateData = insertData
				inserts = append(inserts, insert)
			}
			continue
		}

		return nil, errors.New("could not format insert data")
	}
	return inserts, nil
}

func (c *Client) DocumentTemplates(ctx Context, caseType CaseType) ([]DocumentTemplateData, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/lpa-api/v1/templates/%s", caseType), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var v documentTemplateApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	data, err := v.toDocumentData()
	if err != nil {
		return nil, err
	}
	return data, err
}
