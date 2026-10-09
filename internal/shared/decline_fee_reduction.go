package shared

import (
	"encoding/json"
)

type DecisionType int

const (
	DecisionTypeExemption DecisionType = iota
	DecisionTypeRemission
	DecisionTypeNotRecognised
)

var decisionTypeMap = map[string]DecisionType{
	"DECLINED_EXEMPTION": DecisionTypeExemption,
	"DECLINED_REMISSION": DecisionTypeRemission,
	"notRecognised":      DecisionTypeNotRecognised,
}

func (d DecisionType) String() string {
	return d.Key()
}

func (d DecisionType) Translation() string {
	switch d {
	case DecisionTypeExemption:
		return "Exemption declined"
	case DecisionTypeRemission:
		return "Remission declined"
	case DecisionTypeNotRecognised:
		return "Not recognised"
	default:
		return "decision type NOT RECOGNISED: " + d.String()
	}
}

func (d DecisionType) Key() string {
	switch d {
	case DecisionTypeExemption:
		return "DECLINED_EXEMPTION"
	case DecisionTypeRemission:
		return "DECLINED_REMISSION"
	case DecisionTypeNotRecognised:
		return "DECLINED_NOT_RECOGNISED"
	default:
		return ""
	}
}

func ParseDecisionType(s string) DecisionType {
	value, ok := decisionTypeMap[s]
	if !ok {
		return DecisionTypeNotRecognised
	}
	return value
}

func (d DecisionType) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Key())
}

func (d *DecisionType) UnmarshalJSON(data []byte) (err error) {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*d = ParseDecisionType(s)
	return nil
}
