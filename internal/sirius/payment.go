package sirius

import "fmt"

type Payment struct {
	ID          int           `json:"id,omitempty"`
	CaseID      int           `json:"case_id,omitempty"`
	Source      PaymentSource `json:"source"`
	Amount      FeeString     `json:"amount"`
	PaymentDate DateString    `json:"paymentdate"`
	Type        TypeOfPayment `json:"type"`
	CreatedDate DateString    `json:"createddate"`
	Locked      bool          `json:"locked,omitempty"`
	CreatedByID int           `json:"createdby_id"`
}

type PaymentSource struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type TypeOfPayment struct {
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
