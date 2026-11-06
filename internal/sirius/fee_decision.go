package sirius

import (
	"fmt"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
)

type FeeDecision struct {
	ID             int                 `json:"id,omitempty"`
	DecisionDate   DateString          `json:"decisionDate,omitempty"`
	DecisionType   shared.DecisionType `json:"decisionType,omitempty"`
	DecisionReason string              `json:"decisionReason,omitempty"`
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

func (c *Client) FeeDecisionsByCase(ctx Context, id int) ([]FeeDecision, error) {
	var feeDecisions []FeeDecision

	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/cases/%d/fee-decisions", id), &feeDecisions)
	if err != nil {
		return nil, err
	}

	return feeDecisions, nil
}
