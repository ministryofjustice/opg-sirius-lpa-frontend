package sirius

type BankHolidays map[string]map[string]string

func (c *Client) BankHolidays(ctx Context) (BankHolidays, error) {
	var b BankHolidays
	err := c.get(ctx, "/lpa-api/v1/dates/bank-holidays", &b)

	return b, err
}
