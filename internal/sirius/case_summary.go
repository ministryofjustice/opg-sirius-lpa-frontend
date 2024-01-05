package sirius

import (
	/*
	"fmt"
	"net/http"
	"net/url"
	*/
	"golang.org/x/sync/errgroup"
)

type CaseSummary struct {
	DigitalLpa  DigitalLpa
	TaskList    []Task
	WarningList []Warning
}

func (c *Client) warningsForCase(ctx Context, caseId int) ([]Warning, error) {
	/*
	path := fmt.Sprintf("/lpa-api/v1/cases/%d/warnings", caseId)

	querystring := url.Values{}
	querystring.Set("limit", "99")
	querystring.Set("filter", "status:Not started,active:true")
	querystring.Set("sort", "duedate:ASC")
	*/

	/*
	req, err := c.newRequestWithQuery(ctx, http.MethodGet, path, querystring, nil)

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
	*/

	/*
	var v taskList
	err = json.NewDecoder(resp.Body).Decode(&v)
	if err != nil {
		return nil, err
	}
	*/

	warnings := []Warning{
		Warning{ID: 2, WarningType: "Warning type", WarningText: "Warning text sits underneath. Warning text sits underneath.", DateAdded: "12th Dec 2023",},
		Warning{ID: 1, WarningType: "Warning type", WarningText: "Warning text sits underneath. Warning text sits underneath.", DateAdded: "11th Aug 2023",},
		Warning{ID: 3, WarningType: "Warning type", WarningText: "Warning text sits underneath. Warning text sits underneath.", DateAdded: "8th Jan 2023",},
		Warning{ID: 4, WarningType: "Warning type", WarningText: "Warning text sits underneath. Warning text sits underneath.", DateAdded: "17th Dec 2022",},
		Warning{ID: 5, WarningType: "Warning type", WarningText: "Warning text sits underneath. Warning text sits underneath.", DateAdded: "1st Sept 2022",},
	}

	return warnings, nil
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
		cs.TaskList, err = c.TasksForCase(ctx.With(groupCtx), cs.DigitalLpa.ID)
		if err != nil {
			return err
		}
		return nil
	})

	group.Go(func() error {
		cs.WarningList, err = c.warningsForCase(ctx.With(groupCtx), cs.DigitalLpa.ID)
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
