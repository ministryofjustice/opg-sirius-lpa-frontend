package sirius

import "fmt"

type Payment struct {
	ID          int           `json:"id,omitempty"`
	Source      PaymentSource `json:"source"`
	Amount      int           `json:"amount"`
	PaymentDate DateString    `json:"paymentDate"`
}

type PaymentSource struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (c *Client) Payments(ctx Context, id int) ([]Payment, error) {
	var p []Payment

	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/cases/%d/payments", id), &p)
	if err != nil {
		return nil, err
	}

	return p, nil
}
