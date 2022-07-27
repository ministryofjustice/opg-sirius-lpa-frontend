package sirius

func (c *Client) NoteTypes(ctx Context) ([]string, error) {
	var v []string
	err := c.get(ctx, "/lpa-api/v1/note-types/lpa", &v)

	return v, err
}
