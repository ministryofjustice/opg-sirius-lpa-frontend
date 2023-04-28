package sirius

import (
	"fmt"
	"net/url"
)

type Address struct {
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
	AddressLine3 string `json:"addressLine3"`
	Town         string `json:"town"`
	Postcode     string `json:"postcode"`
	Description  string `json:"description"`
}

func (c *Client) PostcodeLookup(ctx Context, postcode string) ([]Address, error) {
	var v []Address
	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/postcode-lookup?postcode=%s", url.QueryEscape(postcode)), &v)

	return v, err
}
