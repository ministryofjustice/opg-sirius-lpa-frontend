package sirius

import "fmt"

func (c *Client) UpdateObjection(ctx Context, objectionId string, objectionDetails AddObjection) error {
	return c.put(ctx, fmt.Sprintf("/lpa-api/v1/objections/%s", objectionId), objectionDetails, nil)
}
