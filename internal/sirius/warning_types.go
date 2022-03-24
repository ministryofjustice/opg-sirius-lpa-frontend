package sirius

type RefDataItem struct {
	Handle string `json:"handle"`
	Label  string `json:"label"`
}

func (c *Client) WarningTypes(ctx Context) ([]RefDataItem, error) {
	var v []RefDataItem
	err := c.get(ctx, "/api/v1/reference-data/warningType", &v)

	return v, err
}
