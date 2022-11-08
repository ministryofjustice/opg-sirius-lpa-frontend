package sirius

import "fmt"

type Investigation struct {
	ID                       int          `json:"id,omitempty"`
	Title                    string       `json:"investigationTitle"`
	Information              string       `json:"additionalInformation"`
	Type                     string       `json:"type"`
	DateReceived             DateString   `json:"investigationReceivedDate"`
	HoldPeriods              []HoldPeriod `json:"holdPeriods,omitempty"`
	IsOnHold                 bool         `json:"isOnHold,omitempty"`
	ApprovalDate             DateString   `json:"approvalDate,omitempty"`
	ApprovalOutcome          string       `json:"approvalOutcome,omitempty"`
	InvestigationClosureDate DateString   `json:"investigationClosureDate,omitempty"`
}

func (c *Client) Investigation(ctx Context, id int) (Investigation, error) {
	var v Investigation
	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/investigations/%d", id), &v)

	return v, err
}
