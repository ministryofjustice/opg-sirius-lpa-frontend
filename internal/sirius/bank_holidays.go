package sirius

type BankHolidays map[string]map[string]string

func (c *Client) BankHolidays(ctx Context) (BankHolidays, error) {
	var b BankHolidays

	if cached, ok := getCached("bank-holidays"); ok {
		return cached.(BankHolidays), nil
	}

	err := c.get(ctx, "/lpa-api/v1/dates/bank-holidays", &b)

	setCached("bank-holidays", b)

	return b, err
}
