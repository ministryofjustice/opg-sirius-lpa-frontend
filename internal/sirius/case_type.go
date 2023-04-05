package sirius

import (
	"errors"
	"strings"
)

type CaseType string

const (
	CaseTypeLpa = CaseType("lpa")
	CaseTypeEpa = CaseType("epa")
)

func ParseCaseType(s string) (CaseType, error) {
	switch strings.ToLower(s) {
	case "lpa":
		return CaseTypeLpa, nil
	case "epa":
		return CaseTypeEpa, nil
	}

	return CaseType(""), errors.New("could not parse case type")
}
