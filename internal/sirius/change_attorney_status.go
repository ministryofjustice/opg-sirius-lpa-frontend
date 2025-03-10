package sirius

import (
	"fmt"
)

type AttorneyUpdatedStatus struct {
	UID    string `json:"uid"`
	Status string `json:"status"`
}

type attorneyStatusRequest struct {
	AttorneyStatuses []AttorneyUpdatedStatus `json:"attorneyStatuses"`
}

func (c *Client) ChangeAttorneyStatus(ctx Context, caseUID string, attorneyUpdatedStatus []AttorneyUpdatedStatus) error {
	data := attorneyStatusRequest{AttorneyStatuses: attorneyUpdatedStatus}

	return c.put(ctx, fmt.Sprintf("/lpa-api/v1/digital-lpas/%s/attorney-status", caseUID), data, nil)
}
