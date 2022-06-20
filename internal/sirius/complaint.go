package sirius

import "fmt"

type Complaint struct {
	Category       string     `json:"category"`
	Description    string     `json:"description"`
	ReceivedDate   DateString `json:"receivedDate"`
	Severity       string     `json:"severity"`
	SubCategory    string     `json:"subCategory"`
	Summary        string     `json:"summary"`
	Resolution     string     `json:"resolution,omitempty"`
	ResolutionInfo string     `json:"resolutionInfo,omitempty"`
	ResolutionDate DateString `json:"resolutionDate,omitempty"`
}

func (c *Client) Complaint(ctx Context, id int) (Complaint, error) {
	var v Complaint
	err := c.get(ctx, fmt.Sprintf("/api/v1/complaints/%d", id), &v)

	return v, err
}
