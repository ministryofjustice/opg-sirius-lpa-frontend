package sirius

type CaseSummary struct {
	Lpa DigitalLpa
	TaskList []Task
}

/**
 * Get data for the case summary area (digital LPA record and tasks for that LPA).
 */
func (c *Client) CaseSummary(ctx Context, uid string) (CaseSummary, error) {
	var cs CaseSummary
	var err error

	cs.Lpa, err = c.DigitalLpa(ctx, uid)
	if err != nil {
		return cs, err
	}

	cs.TaskList, err = c.TasksForCase(ctx, cs.Lpa.ID)
	if err != nil {
		return cs, err
	}

	return cs, nil
}
