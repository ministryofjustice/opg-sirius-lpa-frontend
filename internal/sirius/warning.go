package sirius

import (
	"fmt"
)

type Warning struct {
	ID          int    `json:"id"`
	DateAdded   string `json:"dateAdded"`
	WarningType string `json:"warningType"`
	WarningText string `json:"warningText"`
	CaseItems   []Case `json:"caseItems"`
}

func (c *Client) WarningsForCase(ctx Context, caseId int) ([]Warning, error) {
	var warningList []Warning
	path := fmt.Sprintf("/lpa-api/v1/cases/%d/warnings", caseId)

	err := c.get(ctx, path, &warningList)
	if err != nil {
		return nil, err
	}

	return warningList, nil
}
