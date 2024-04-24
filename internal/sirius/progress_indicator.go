package sirius

import (
	"fmt"
)

type ProgressIndicator struct {
	Indicator string `json:"indicator"`
	Status    string `json:"status"`
}

type ProgressIndicators struct {
	ProgressIndicators []ProgressIndicator `json:"progressIndicators"`
}

func (c *Client) ProgressIndicatorsForDigitalLpa(ctx Context, uid string) ([]ProgressIndicator, error) {
	path := fmt.Sprintf("/lpa-api/v1/digital-lpas/%s/progress-indicators", uid)

	var receiver ProgressIndicators
	err := c.get(ctx, path, &receiver)

	if err != nil {
		return nil, err
	}

	return receiver.ProgressIndicators, nil
}
