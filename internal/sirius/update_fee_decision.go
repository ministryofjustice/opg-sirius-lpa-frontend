package sirius

import "fmt"

func (c *Client) UpdateFeeDecision(ctx Context, caseID int, feeDecisionID int, decisionType string, decisionReason string, decisionDate DateString) error {
	data := struct {
		DecisionType   string     `json:"decisionType"`
		DecisionReason string     `json:"decisionReason"`
		DecisionDate   DateString `json:"decisionDate"`
	}{
		DecisionType:   decisionType,
		DecisionReason: decisionReason,
		DecisionDate:   decisionDate,
	}

	path := fmt.Sprintf("/lpa-api/v1/cases/%d/fee-decisions/%d", caseID, feeDecisionID)

	return c.put(ctx, path, data, nil)
}
