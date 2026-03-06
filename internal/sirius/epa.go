package sirius

import "fmt"

func (c *Client) CreateEpa(ctx Context, donorID int, epa Case) (int, error) {
	var response Case
	err := c.post(ctx, fmt.Sprintf("/lpa-api/v1/donors/%d/epas", donorID), epa, &response)
	return response.ID, err
}

//func (c *Client) UpdateEpa(ctx Context, caseId int, epa Epa) error {
//	return c.put(ctx, fmt.Sprintf("/lpa-api/v1/epas/%d", caseId), epa, nil)
//}

func (c *Client) UpdateEpa(ctx Context, caseId int, epa Case) error {
	return c.put(ctx, fmt.Sprintf("/lpa-api/v1/epas/%d", caseId), epa, nil)
}
