package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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
	data, err := json.Marshal(allocateCasesRequest{Data: allocations})
	if err != nil {
		return err
	}

	caseIDs := strconv.Itoa(allocations[0].ID)
	for _, allocation := range allocations[1:] {
		caseIDs += "+" + strconv.Itoa(allocation.ID)
	}

	req, err := c.newRequest(ctx, http.MethodPut, fmt.Sprintf("/api/v1/users/%d/cases/%s", assigneeID, caseIDs), bytes.NewReader(data))
	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		var v ValidationError
		if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
			return err
		}
		return v
	}

	if resp.StatusCode != http.StatusOK {
		return newStatusError(resp)
	}

	return nil
}
