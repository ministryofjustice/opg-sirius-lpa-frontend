package sirius

import (
	"fmt"
)

type SeveranceApplicationDetails struct {
	HasDonorConsented      bool       `json:"hasDonorConsented"`
	SeveranceOrdered       bool       `json:"severanceOrdered"`
	CourtOrderDecisionMade DateString `json:"courtOrderDecisionMade"`
	CourtOrderReceived     DateString `json:"courtOrderReceived"`
}

func (c *Client) EditSeveranceApplication(ctx Context, caseUID string, severanceApplicationDetails SeveranceApplicationDetails) error {
	return c.put(ctx, fmt.Sprintf("/lpa-api/v1/digital-lpas/%s/severance", caseUID), severanceApplicationDetails, nil)
}
