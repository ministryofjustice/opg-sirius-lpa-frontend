package shared

import (
	"encoding/json"
)

type LifeSustainingTreatmentOption int

const (
	LifeSustainingTreatmentOptionEmpty LifeSustainingTreatmentOption = iota
	LifeSustainingTreatmentOptionA
	LifeSustainingTreatmentOptionB
	LifeSustainingTreatmentOptionNotRecognised
)

var lifeSustainingTreatmentOptionMap = map[string]LifeSustainingTreatmentOption{
	"":         LifeSustainingTreatmentOptionEmpty,
	"option-a": LifeSustainingTreatmentOptionA,
	"option-b": LifeSustainingTreatmentOptionB,
}

func (l LifeSustainingTreatmentOption) String() string {
	return l.Key()
}

func (l LifeSustainingTreatmentOption) LongForm() string {
	switch l {
	case LifeSustainingTreatmentOptionA:
		return "Attorneys can give or refuse consent to LST"
	case LifeSustainingTreatmentOptionB:
		return "Attorneys cannot give or refuse consent to LST"
	case LifeSustainingTreatmentOptionEmpty:
		return "Not specified"
	default:
		return "lifeSustainingTreatmentOption NOT RECOGNISED: " + l.Key()
	}
}

func (l LifeSustainingTreatmentOption) Key() string {
	switch l {
	case LifeSustainingTreatmentOptionA:
		return "option-a"
	case LifeSustainingTreatmentOptionB:
		return "option-b"
	case LifeSustainingTreatmentOptionEmpty:
		return ""
	case LifeSustainingTreatmentOptionNotRecognised:
		return "notRecognised"
	default:
		return ""
	}
}

func ParseLifeSustainingTreatmentOption(s string) LifeSustainingTreatmentOption {
	value, ok := lifeSustainingTreatmentOptionMap[s]
	if !ok {
		return LifeSustainingTreatmentOptionNotRecognised
	}
	return value
}

func (l LifeSustainingTreatmentOption) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.Key())
}

func (l *LifeSustainingTreatmentOption) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*l = ParseLifeSustainingTreatmentOption(s)
	return nil
}

