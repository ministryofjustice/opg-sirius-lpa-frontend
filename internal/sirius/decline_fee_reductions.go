package sirius

import (
	"fmt"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
)

type DeclineFeeReductions struct {
	ID             int                 `json:"id,omitempty"`
	DecisionDate   DateString          `json:"decisionDate,omitempty"`
	DecisionType   shared.DecisionType `json:"decisionType,omitempty"`
	DecisionReason string              `json:"decisionReason,omitempty"`
}

func (c *Client) DeclineFeeReductions(ctx Context, id int) ([]DeclineFeeReductions, error) {
	var declineFeeReductions []DeclineFeeReductions

	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/cases/%d/fee-decisions", id), &declineFeeReductions)
	if err != nil {
		return nil, err
	}

	return declineFeeReductions, nil
}
