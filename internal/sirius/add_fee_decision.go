package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) AddFeeDecision(ctx Context, caseID int, decisionType string, decisionReason string, decisionDate DateString) error {
	postData, err := json.Marshal(struct {
		DecisionType   string     `json:"decisionType"`
		DecisionReason string     `json:"decisionReason"`
		DecisionDate   DateString `json:"decisionDate"`
	}{
		DecisionType:   decisionType,
		DecisionReason: decisionReason,
		DecisionDate:   decisionDate,
	})

	if err != nil {
		return err
	}

	req, err := c.newRequest(
		ctx,
		http.MethodPost,
		fmt.Sprintf("/lpa-api/v1/cases/%d/fee-decisions", caseID),
		bytes.NewReader(postData),
	)

	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //#nosec G307 false positive

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
