package sirius

func (c *Client) LinkPeople(ctx Context, parentId int, childId int) error {
	postData := struct {
		ParentId int `json:"parentId"`
		ChildId  int `json:"childId"`
	}{
		ParentId: parentId,
		ChildId:  childId,
	}

	return c.post(ctx, "/lpa-api/v1/person-links", postData, nil)
}
