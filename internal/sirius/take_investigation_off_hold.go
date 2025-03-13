package sirius

import (
	"fmt"
)

func (c *Client) TakeInvestigationOffHold(ctx Context, holdPeriodId int) error {
	return c.delete(ctx, fmt.Sprintf("/lpa-api/v1/hold-periods/%d", holdPeriodId))
}
