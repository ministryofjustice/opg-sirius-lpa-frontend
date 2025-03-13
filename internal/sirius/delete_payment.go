package sirius

import (
	"fmt"
)

func (c *Client) DeletePayment(ctx Context, paymentID int) error {
	return c.delete(ctx, fmt.Sprintf("/lpa-api/v1/payments/%d", paymentID))
}
