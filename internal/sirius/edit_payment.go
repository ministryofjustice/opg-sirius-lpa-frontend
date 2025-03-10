package sirius

import (
	"fmt"
)

func (c *Client) EditPayment(ctx Context, paymentID int, payment Payment) error {
	return c.put(ctx, fmt.Sprintf("/lpa-api/v1/payments/%d", paymentID), payment, nil)
}
