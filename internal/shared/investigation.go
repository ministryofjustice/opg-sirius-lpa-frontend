package shared

import "strings"

var InvestigationEventPropertyTranslationMap = map[string]string{
	"INVESTIGATIONTITLE":        "Title",
	"ADDITIONALINFORMATION":     "Information",
	"RISKASSESSMENTDATE":        "Risk assessment date",
	"REPORTAPPROVALDATE":        "Court Application",
	"INVESTIGATIONCLOSUREDATE":  "Investigation closed",
	"REPORTAPPROVALOUTCOME":     "Report approval outcome",
	"TYPE":                      "Type",
	"INVESTIGATIONRECEIVEDDATE": "Investigation received",
}

func TranslateInvestigationEventProperty(property string) string {
	if val, ok := InvestigationEventPropertyTranslationMap[strings.ToUpper(property)]; ok {
		return val
	}
	return property
}
