package sirius

import "errors"

type CaseType string

const (
	CaseTypeLpa = CaseType("lpa")
	CaseTypeEpa = CaseType("epa")
)

func ParseCaseType(s string) (CaseType, error) {
	switch s {
	case "lpa":
		return CaseTypeLpa, nil
	case "epa":
		return CaseTypeEpa, nil
	}

	return CaseType(""), errors.New("could not parse case type")
}
