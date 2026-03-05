package shared

import "strings"

var ComplaintPropertyTranslationMap = map[string]string{
	"CATEGORY":             "Category",
	"SUBCATEGORY":          "Subcategory",
	"COMPLAINANTCATEGORY":  "Complainant category",
	"ORIGIN":               "Origin",
	"SEVERITY":             "Severity",
	"INVESTIGATINGOFFICER": "Investigating officer",
	"COMPLAINANTNAME":      "Complainant name",
	"SUMMARY":              "Title",
	"DESCRIPTION":          "Description",
	"COMPENSATIONAMOUNT":   "Compensation amount",
	"COMPENSATIONTYPE":     "Compensation type",
	"RECEIVEDDATE":         "Received date",
	"RESOLUTIONDATE":       "Resolution date",
	"RESOLUTION":           "Resolution",
	"RESOLUTIONINFO":       "Resolution notes",
}

func TranslateComplaintProperty(property string) string {
	if val, ok := ComplaintPropertyTranslationMap[strings.ToUpper(property)]; ok {
		return val
	}
	return property
}
