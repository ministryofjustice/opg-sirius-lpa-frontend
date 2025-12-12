package shared

import (
	"encoding/json"
)

type CaseLabel int

const (
	CaseLabelHW CaseLabel = iota
	CaseLabelPFA
	CaseLabelEmpty
	CaseLabelNotRecognised
)

var caseLabelMap = map[string]CaseLabel{
	"hw":            CaseLabelHW,
	"pfa":           CaseLabelPFA,
	"":              CaseLabelEmpty,
	"notRecognised": CaseLabelNotRecognised,
}

func (c CaseLabel) String() string {
	return c.Key()
}

func (c CaseLabel) Translation() string {
	switch c {
	case CaseLabelHW:
		return "colour-govuk-grass-green"
	case CaseLabelPFA:
		return "colour-govuk-turquoise"
	case CaseLabelEmpty:
		return "Not specified"
	default:
		return "case label NOT RECOGNISED: " + c.String()
	}
}

func (c CaseLabel) Key() string {
	switch c {
	case CaseLabelHW:
		return "hw"
	case CaseLabelPFA:
		return "pfa"
	case CaseLabelEmpty:
		return "Empty"
	case CaseLabelNotRecognised:
		return "Not Recognised"
	default:
		return ""
	}
}

func ParseCaseLabel(s string) CaseLabel {
	value, ok := caseLabelMap[s]
	if !ok {
		return CaseLabelNotRecognised
	}
	return value
}

func (c CaseLabel) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Key())
}

func (c *CaseLabel) UnmarshalJSON(data []byte) (err error) {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*c = ParseCaseLabel(s)
	return nil
}
