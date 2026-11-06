package sirius

import (
	"fmt"
	"strings"
)

func (c *Client) EditCase(ctx Context, caseID int, caseType CaseType, caseDetails Case) error {
	return c.put(ctx, fmt.Sprintf("/lpa-api/v1/%ss/%d", strings.ToLower(string(caseType)), caseID), caseDetails, nil)
}
