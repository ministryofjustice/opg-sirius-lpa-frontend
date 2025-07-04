package sirius

import (
	"fmt"
)

type AttorneyDecisions struct {
	UID                      string `json:"uid"`
	CannotMakeJointDecisions bool   `json:"cannotMakeJointDecisions"`
}

type AttorneyDecisionsRequest struct {
	AttorneyDecisions []AttorneyDecisions `json:"attorneyDecisions"`
}

func (c *Client) ManageAttorneyDecisions(ctx Context, caseUID string, attorneyDecisions []AttorneyDecisions) error {
	data := AttorneyDecisionsRequest{AttorneyDecisions: attorneyDecisions}

	return c.put(ctx, fmt.Sprintf("/lpa-api/v1/digital-lpas/%s/attorney-decisions", caseUID), data, nil)
}
