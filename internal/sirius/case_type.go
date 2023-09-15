package sirius

import (
	"errors"
	"strings"
)

type CaseType string

const (
	CaseTypeLpa        = CaseType("lpa")
	CaseTypeEpa        = CaseType("epa")
	CaseTypeDigitalLpa = CaseType("digital_lpa")
)

func ParseCaseType(s string) (CaseType, error) {
	switch strings.ToLower(s) {
	case "lpa":
		return CaseTypeLpa, nil
	case "epa":
		return CaseTypeEpa, nil
	case "digital_lpa":
		return CaseTypeDigitalLpa, nil
	}

	return CaseType(""), errors.New("could not parse case type")
}
