package sirius

import (
	"fmt"
)

type SeveranceApplication struct {
	HasDonorConsented      *bool      `json:"hasDonorConsented,omitempty"`
	SeveranceOrdered       *bool      `json:"severanceOrdered,omitempty"`
	CourtOrderDecisionMade DateString `json:"courtOrderDecisionMade,omitempty"`
	CourtOrderReceived     DateString `json:"courtOrderReceived,omitempty"`
}

func (c *Client) EditSeveranceApplication(ctx Context, caseUID string, severanceApplication SeveranceApplication) error {
	return c.put(ctx, fmt.Sprintf("/lpa-api/v1/digital-lpas/%s/severance", caseUID), severanceApplication, nil)
}
