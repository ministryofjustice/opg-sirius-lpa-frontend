package sirius

import "sort"

func (c *Client) TaskTypes(ctx Context) ([]string, error) {
	var v struct {
		TaskTypes map[string]struct{} `json:"task_types"`
	}
	if err := c.get(ctx, "/api/v1/tasktypes/lpa", &v); err != nil {
		return nil, err
	}

	var types []string
	for k := range v.TaskTypes {
		types = append(types, k)
	}
	sort.Strings(types)

	return types, nil
}
