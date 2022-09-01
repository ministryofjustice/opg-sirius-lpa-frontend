package sirius

import "fmt"

type Payment struct {
	ID          int        `json:"id,omitempty"`
	Source      string     `json:"source"`
	Amount      int        `json:"amount"`
	PaymentDate DateString `json:"paymentDate"`
}

func (c *Client) Payments(ctx Context, id int) ([]Payment, error) {
	var p []Payment

	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/cases/%d/payments", id), &p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (c *Client) PaymentByID(ctx Context, id int) (Payment, error) {
	var p Payment
	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/payments/%d", id), &p)

	return p, err
}

func PoundsToPence(pounds float64) int {
	return int(pounds * 100)
}

func PenceToPounds(pence int) float64 {
	return float64(pence) / 100
}
