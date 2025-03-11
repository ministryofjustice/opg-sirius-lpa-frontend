package sirius

import (
	"fmt"
)

func (c *Client) AddFeeDecision(ctx Context, caseID int, decisionType string, decisionReason string, decisionDate DateString) error {
	data := struct {
		DecisionType   string     `json:"decisionType"`
		DecisionReason string     `json:"decisionReason"`
		DecisionDate   DateString `json:"decisionDate"`
	}{
		DecisionType:   decisionType,
		DecisionReason: decisionReason,
		DecisionDate:   decisionDate,
	}

	return c.post(ctx, fmt.Sprintf("/lpa-api/v1/cases/%d/fee-decisions", caseID), data, nil)
}
