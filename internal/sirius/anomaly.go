package sirius

import (
	"fmt"
)

type ObjectUid string
type ObjectFieldName string

type AnomalyStatus string

const (
	AnomalyAccepted = AnomalyStatus("accepted")
	AnomalyDetected = AnomalyStatus("detected")
	AnomalyFatal    = AnomalyStatus("fatal")
	AnomalyResolved = AnomalyStatus("resolved")
)

type Anomaly struct {
	Id            int             `json:"id"`
	Status        AnomalyStatus   `json:"status"`
	FieldName     ObjectFieldName `json:"fieldName"`
	RuleType      AnomalyRuleType `json:"ruleType"`
	FieldOwnerUid ObjectUid       `json:"fieldOwnerUid"`
}

func (c *Client) AnomaliesForDigitalLpa(ctx Context, uid string) ([]Anomaly, error) {
	path := fmt.Sprintf("/lpa-api/v1/digital-lpas/%s/anomalies", uid)

	var receiver []Anomaly
	err := c.get(ctx, path, &receiver)

	if err != nil {
		return nil, err
	}

	return receiver, nil
}
