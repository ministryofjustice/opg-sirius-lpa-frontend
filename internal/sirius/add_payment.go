package sirius

import (
	"fmt"
)

func (c *Client) AddPayment(ctx Context, caseID int, amount int, source string, paymentDate DateString) error {
	data := struct {
		Amount      int        `json:"amount"`
		Source      string     `json:"source"`
		PaymentDate DateString `json:"paymentDate"`
	}{
		Amount:      amount,
		Source:      source,
		PaymentDate: paymentDate,
	}

	return c.post(ctx, fmt.Sprintf("/lpa-api/v1/cases/%d/payments", caseID), data, nil)
}
