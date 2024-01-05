package sirius

type CaseSummary struct {
	DigitalLpa  DigitalLpa
	TaskList    []Task
	WarningList []Warning
}

/**
 * Get data for the case summary area (digital LPA record, tasks, and warnings for that LPA)
 */
func (c *Client) CaseSummary(ctx Context, uid string) (CaseSummary, error) {
	var cs CaseSummary
	var err error

	cs.DigitalLpa, err = c.DigitalLpa(ctx, uid)
	if err != nil {
		return cs, err
	}

	cs.TaskList, err = c.TasksForCase(ctx, cs.DigitalLpa.SiriusData.ID)
	if err != nil {
		return cs, err
	}

	return cs, nil
}
