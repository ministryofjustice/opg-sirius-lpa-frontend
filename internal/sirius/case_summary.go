package sirius

import (
	"golang.org/x/sync/errgroup"
)

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
	group, groupCtx := errgroup.WithContext(ctx.Context)

	cs.DigitalLpa, err = c.DigitalLpa(ctx, uid)
	if err != nil {
		return cs, err
	}

	group.Go(func() error {
		cs.TaskList, err = c.TasksForCase(ctx.With(groupCtx), cs.DigitalLpa.SiriusData.ID)
		if err != nil {
			return err
		}
		return nil
	})

	group.Go(func() error {
		cs.WarningList, err = c.WarningsForCase(ctx.With(groupCtx), cs.DigitalLpa.SiriusData.ID)
		if err != nil {
			return err
		}
		return nil
	})

	if err := group.Wait(); err != nil {
		return cs, err
	}

	return cs, nil
}
