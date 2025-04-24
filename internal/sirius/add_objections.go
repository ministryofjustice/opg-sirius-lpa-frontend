package sirius

type AddObjection struct {
	LpaUids       []string   `json:"lpaUids"`
	ReceivedDate  DateString `json:"receivedDate"`
	ObjectionType string     `json:"objectionType"`
	Notes         string     `json:"notes"`
}

func (c *Client) AddObjection(ctx Context, objectionDetails AddObjection) error {
	return c.post(ctx, "/lpa-api/v1/objections", objectionDetails, nil)
}
