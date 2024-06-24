package sirius

import (
	"fmt"
)

type Anomaly struct {
	Id            int    `json:"id"`
	Status        string `json:"status"`
	FieldName     string `json:"fieldName"`
	RuleType      string `json:"ruleType"`
	FieldOwnerUid string `json:"fieldOwnerUid"`
}

type Anomalies struct {
	Anomalies []Anomaly `json:"anomalies"`
}

func (c *Client) AnomaliesForDigitalLpa(ctx Context, uid string) ([]Anomaly, error) {
	path := fmt.Sprintf("/lpa-api/v1/digital-lpas/%s/anomalies", uid)

	var receiver Anomalies
	err := c.get(ctx, path, &receiver)

	if err != nil {
		return nil, err
	}

	return receiver.Anomalies, nil
}
