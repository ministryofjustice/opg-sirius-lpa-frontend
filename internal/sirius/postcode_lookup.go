package sirius

import (
	"fmt"
	"net/url"
)

type PostcodeLookupAddress struct {
	Line1       string `json:"addressLine1"`
	Line2       string `json:"addressLine2"`
	Line3       string `json:"addressLine3"`
	Town        string `json:"town"`
	Postcode    string `json:"postcode"`
	Country     string `json:"country"`
	Description string `json:"description"`
}

func (c *Client) PostcodeLookup(ctx Context, postcode string) ([]PostcodeLookupAddress, error) {
	var v []PostcodeLookupAddress
	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/postcode-lookup?postcode=%s", url.QueryEscape(postcode)), &v)

	return v, err
}
