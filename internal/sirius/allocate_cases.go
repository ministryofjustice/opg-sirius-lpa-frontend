package sirius

import (
	"fmt"
	"strconv"
)

type CaseAllocation struct {
	ID       int    `json:"id"`
	CaseType string `json:"caseType"`
}

type allocateCasesRequest struct {
	Data []CaseAllocation `json:"data"`
}

func (c *Client) AllocateCases(ctx Context, assigneeID int, allocations []CaseAllocation) error {
	data := allocateCasesRequest{Data: allocations}

	caseIDs := strconv.Itoa(allocations[0].ID)
	for _, allocation := range allocations[1:] {
		caseIDs += "+" + strconv.Itoa(allocation.ID)
	}

	return c.put(ctx, fmt.Sprintf("/lpa-api/v1/users/%d/cases/%s", assigneeID, caseIDs), data, nil)
}
