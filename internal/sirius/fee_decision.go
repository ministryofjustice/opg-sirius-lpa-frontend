package sirius

import (
	"fmt"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
)

type FeeDecision struct {
	ID             int                 `json:"id"`
	DecisionType   shared.DecisionType `json:"decisionType"`
	DecisionReason string              `json:"decisionReason"`
	DecisionDate   DateString          `json:"decisionDate"`
}

func (c *Client) CreateFeeDecision(ctx Context, caseID int, decisionType string, decisionReason string, decisionDate DateString) error {
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

func (c *Client) GetFeeDecision(ctx Context, feeDecisionID int) (FeeDecision, error) {
	var v FeeDecision
	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/fee-decisions/%d", feeDecisionID), &v)

	return v, err
}

func (c *Client) DeleteFeeDecision(ctx Context, feeDecisionID int) error {
	return c.delete(ctx, fmt.Sprintf("/lpa-api/v1/fee-decisions/%d", feeDecisionID))
}
