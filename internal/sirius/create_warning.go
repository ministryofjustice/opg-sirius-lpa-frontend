package sirius

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func (c *Client) CreateWarning(ctx Context, personId int, warningType string, warningNote string) error {
	postData, err := json.Marshal(struct {
		PersonID    int    `json:"personId"`
		WarningType string `json:"warningType"`
		WarningText string `json:"warningText"`
	}{
		PersonID:    personId,
		WarningType: warningType,
		WarningText: warningNote,
	})

	if err != nil {
		return err
	}

	req, err := c.newRequest(
		ctx,
		http.MethodPost,
		"/lpa-api/v1/warnings",
		bytes.NewReader(postData),
	)

	if err != nil {
		return err
	}

	res, err := c.http.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close() //#nosec G307 false positive

	if res.StatusCode == http.StatusBadRequest {
		var v ValidationError
		if err := json.NewDecoder(res.Body).Decode(&v); err != nil {
			return err
		}
		return v
	}

	if res.StatusCode != http.StatusCreated {
		return newStatusError(res)
	}
	return nil
}
