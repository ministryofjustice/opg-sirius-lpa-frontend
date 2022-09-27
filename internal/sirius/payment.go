package sirius

import (
	"fmt"
	"regexp"
)

type PaymentReference struct {
	Reference string `json:"reference"`
	Type      string `json:"type"`
}

type Payment struct {
	ID          int                `json:"id,omitempty"`
	Source      string             `json:"source,omitempty"`
	Amount      int                `json:"amount,omitempty"`
	PaymentDate DateString         `json:"paymentDate,omitempty"`
	Case        *Case              `json:"case,omitempty"`
	Locked      bool               `json:"locked,omitempty"`
	References  []PaymentReference `json:"references,omitempty"`
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

func IsAmountValid(amount string) bool {
	m, err := regexp.Match(`^\d+\.\d{2}$`, []byte(amount))
	if err != nil {
		return false
	}
	return m
}
