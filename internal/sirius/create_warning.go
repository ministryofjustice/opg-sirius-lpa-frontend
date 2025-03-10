package sirius

func (c *Client) CreateWarning(ctx Context, personId int, warningType string, warningNote string, caseIDs []int) error {
	type WarningRequest struct {
		PersonID    int    `json:"personId,omitempty"`
		CaseIDs     []int  `json:"caseIds,omitempty"`
		WarningType string `json:"warningType"`
		WarningText string `json:"warningText"`
	}

	data := WarningRequest{
		PersonID:    personId,
		CaseIDs:     caseIDs,
		WarningType: warningType,
		WarningText: warningNote,
	}

	return c.post(ctx, "/lpa-api/v1/warnings", data, nil)
}
