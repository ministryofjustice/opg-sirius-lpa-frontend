package sirius

import "fmt"

const (
	PaymentSourceCategory string = "paymentSource"
	WarningTypeCategory   string = "warningType"
	FeeReductionTypeCategory string = "feeReductionType"
)

type RefDataItem struct {
	Handle         string `json:"handle"`
	Label          string `json:"label"`
	UserSelectable bool   `json:"userSelectable"`
}

func (c *Client) RefDataByCategory(ctx Context, category string) ([]RefDataItem, error) {
	var v []RefDataItem
	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/reference-data/%s", category), &v)

	return v, err
}
