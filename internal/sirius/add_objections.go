package sirius

type AddObjections struct {
	LpaUids       []string   `json:"lpaUids"`
	ReceivedDate  DateString `json:"receivedDate"`
	ObjectionType string     `json:"objectionType"`
	Notes         string     `json:"notes"`
}

func (c *Client) AddObjections(ctx Context, objectionDetails AddObjections) error {
	return c.post(ctx, "/lpa-api/v1/objections", objectionDetails, nil)
}
