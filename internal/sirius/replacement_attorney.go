package sirius

import "fmt"

type createReplacementAttorneyRequest struct {
	Attorney
	PersonType string `json:"personType"`
	CaseID     int    `json:"caseId"`
}

func (c *Client) CreateReplacementAttorney(ctx Context, caseId int, attorney Attorney) (Attorney, error) {
	req := []createReplacementAttorneyRequest{
		{
			Attorney:   attorney,
			PersonType: "ReplacementAttorney",
			CaseID:     caseId,
		},
	}

	var resp []Attorney
	if err := c.post(ctx, "/lpa-api/v1/persons", req, &resp); err != nil {
		return Attorney{}, err
	}

	if len(resp) == 0 {
		return Attorney{}, nil
	}

	return resp[0], nil
}

func (c *Client) UpdateReplacementAttorney(ctx Context, attorneyId int, attorney Attorney) error {
	return c.put(ctx, fmt.Sprintf("/lpa-api/v1/replacement-attorneys/%d", attorneyId), attorney, nil)
}
