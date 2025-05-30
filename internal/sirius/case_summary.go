package sirius

import (
	"golang.org/x/sync/errgroup"
)

type CaseSummary struct {
	DigitalLpa  DigitalLpa
	TaskList    []Task
	WarningList []Warning
	Objections  []ObjectionForCase
	Resolution  ObjectionResolution
}

/**
 * Get data for the case summary area (digital LPA record, tasks, and warnings for that LPA)
 */
func (c *Client) CaseSummary(ctx Context, uid string) (CaseSummary, error) {
	return c.getCaseSummary(ctx, uid, false)
}

/**
 * Get data for the case summary area (digital LPA record, tasks, and warnings
 * for that LPA) including presigned URLs for images
 */
func (c *Client) CaseSummaryWithImages(ctx Context, uid string) (CaseSummary, error) {
	return c.getCaseSummary(ctx, uid, true)
}

func (c *Client) getCaseSummary(ctx Context, uid string, presignImages bool) (CaseSummary, error) {
	var cs CaseSummary
	var err error
	group, groupCtx := errgroup.WithContext(ctx.Context)

	cs.DigitalLpa, err = c.DigitalLpa(ctx, uid, presignImages)
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

	group.Go(func() error {
		cs.Objections, err = c.ObjectionsForCase(ctx.With(groupCtx), uid)
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
