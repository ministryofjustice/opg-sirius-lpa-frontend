package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/sync/errgroup"
)

type CaseSummary struct {
	DigitalLpa  DigitalLpa
	TaskList    []Task
	WarningList []Warning
}

func (c *Client) warningsForCase(ctx Context, caseId int) ([]Warning, error) {
	path := fmt.Sprintf("/lpa-api/v1/cases/%d/warnings", caseId)

	req, err := c.newRequest(ctx, http.MethodGet, path, nil)

	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close() //#nosec G307 false positive

	if resp.StatusCode != http.StatusOK {
		return nil, newStatusError(resp)
	}

	var warningList []Warning
	err = json.NewDecoder(resp.Body).Decode(&warningList)
	if err != nil {
		return nil, err
	}

	return warningList, nil
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
		cs.WarningList, err = c.warningsForCase(ctx.With(groupCtx), cs.DigitalLpa.SiriusData.ID)
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
