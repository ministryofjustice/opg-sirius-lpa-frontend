package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) CreateWarning(ctx Context, personId int, warningType string, warningNote string) error {
	postData, err := json.Marshal(struct {
		WarningType string `json:"warningType"`
		WarningText string `json:"warningText"`
	}{
		WarningType: warningType,
		WarningText: warningNote,
	})

	if err != nil {
		return err
	}

	req, err := c.newRequest(
		ctx,
		http.MethodPost,
		fmt.Sprintf("/lpa-api/v1/persons/%d/warnings", personId),
		bytes.NewReader(postData),
	)

	if err != nil {
		return err
	}

	res, err := c.http.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

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
