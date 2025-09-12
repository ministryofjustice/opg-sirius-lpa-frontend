package sirius

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type Warning struct {
	ID          int    `json:"id"`
	DateAdded   string `json:"dateAdded"`
	WarningType string `json:"warningType"`
	WarningText string `json:"warningText"`
	CaseItems   []Case `json:"caseItems"`
}

func (c *Client) WarningsForCase(ctx Context, caseId int) ([]Warning, error) {
	var warningList []Warning
	path := fmt.Sprintf("/lpa-api/v1/cases/%d/warnings", caseId)

	err := c.get(ctx, path, &warningList)
	if err != nil {
		return nil, err
	}

	return sortWarningsForCaseSummary(warningList), nil
}

func sortWarningsForCaseSummary(warnings []Warning) []Warning {
	sort.Slice(warnings, func(i, j int) bool {
		if warnings[i].WarningType == "Donor Deceased" {
			return true
		} else if warnings[j].WarningType == "Donor Deceased" {
			return false
		}

		iTime, err := time.Parse("02/01/2006 15:04:05", warnings[i].DateAdded)
		if err != nil {
			return false
		}

		jTime, err := time.Parse("02/01/2006 15:04:05", warnings[j].DateAdded)
		if err != nil {
			return false
		}

		return iTime.After(jTime)
	})

	return warnings
}

// construct string to use in case summary for cases a warning is applied to
func CasesWarningAppliedTo(uid string, cases []Case) string {
	// return value:
	// "" (only this case)
	// or " and <subtype (hw|pw)> <uid>" (one other case)
	// or ", <subtype (hw|pw)> <uid_1>, <subtype (hw|pw)> <uid_2>, ...,
	// <subtype (hw|pw)> <uid_n-1> and <subtype (hw|pw)> <uid_n>" (2 to n other cases)
	if len(cases) == 1 {
		return ""
	}

	var filteredCases []Case
	for _, caseItem := range cases {
		if caseItem.UID != uid {
			filteredCases = append(filteredCases, caseItem)
		}
	}
	numCases := len(filteredCases)

	var b strings.Builder
	for index, caseItem := range filteredCases {
		if index == numCases-1 {
			b.WriteString(" and ")
		} else {
			b.WriteString(", ")
		}
		b.WriteString(subtypeShortFormat(caseItem.SubType))
		b.WriteString(" ")
		b.WriteString(caseItem.UID)
	}

	return b.String()
}

func subtypeShortFormat(subtype string) string {
	switch strings.ToLower(subtype) {
	case "personal-welfare":
		return "PW"
	case "property-and-affairs":
		return "PA"
	case "hw":
		return "HW"
	case "pfa":
		return "PFA"
	default:
		return ""
	}
}
