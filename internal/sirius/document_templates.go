package sirius

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
)

type documentTemplateApiResponse map[string]json.RawMessage
type insertApiResponse map[string]json.RawMessage

type UniversalTemplateData struct {
	Label      string `json:"label"`
	Order      int    `json:"order"`
	UsesNotify bool   `json:"govukNotify"`
}

type documentTemplateApiData struct {
	Inserts    insertApiResponse `json:"inserts"`
	Label      string            `json:"label"`
	UsesNotify bool              `json:"govukNotify"`
}

type Insert struct {
	Key      string
	InsertId string
	Label    string `json:"label"`
	Order    int    `json:"order"`
}

type DocumentTemplateData struct {
	Inserts    []Insert
	TemplateId string
	Label      string `json:"label"`
	UsesNotify bool
}

func (d documentTemplateApiResponse) toDocumentData() ([]DocumentTemplateData, error) {
	var s []DocumentTemplateData

	for k, v := range d {
		var asDocumentTemplateData DocumentTemplateData
		asDocumentTemplateData.TemplateId = k

		var asDocumentTemplateApiData documentTemplateApiData
		if err := json.Unmarshal(v, &asDocumentTemplateApiData); err == nil {
			inserts, err := asDocumentTemplateApiData.Inserts.toInsertData()
			if err != nil {
				return nil, err
			}
			asDocumentTemplateData.Label = asDocumentTemplateApiData.Label
			asDocumentTemplateData.UsesNotify = asDocumentTemplateApiData.UsesNotify
			asDocumentTemplateData.Inserts = inserts
			s = append(s, asDocumentTemplateData)
			continue
		} else {
			// handle when inserts = []
			var universalTemplateData UniversalTemplateData
			if err := json.Unmarshal(v, &universalTemplateData); err == nil {
				asDocumentTemplateData.Label = universalTemplateData.Label
				asDocumentTemplateData.UsesNotify = universalTemplateData.UsesNotify
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
				insert.Label = insertData.Label
				insert.Order = insertData.Order
				inserts = append(inserts, insert)
			}
			continue
		}

		return nil, errors.New("could not format insert data")
	}

	sort.Slice(inserts, func(i, j int) bool {
		if inserts[i].Key != inserts[j].Key {
			return inserts[i].Key < inserts[j].Key
		}

		return inserts[i].Order < inserts[j].Order
	})

	return inserts, nil
}

func (c *Client) DocumentTemplates(ctx Context, caseType CaseType) ([]DocumentTemplateData, error) {
	templateSet := caseType
	if caseType == CaseTypeDigitalLpa {
		templateSet = "digitallpa"
	}

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/lpa-api/v1/templates/%s", templateSet), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //#nosec G307 false positive

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
