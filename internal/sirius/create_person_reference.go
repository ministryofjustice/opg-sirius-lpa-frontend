package sirius

import (
	"fmt"
)

type createPersonReferenceRequest struct {
	Reason        string `json:"reason"`
	ReferencedUID string `json:"referencedUid"`
}

func (c *Client) CreatePersonReference(ctx Context, personID int, referencedUID, reason string) error {
	data := createPersonReferenceRequest{
		Reason:        reason,
		ReferencedUID: referencedUID,
	}

	return c.post(ctx, fmt.Sprintf("/lpa-api/v1/persons/%d/references", personID), data, nil)
}
