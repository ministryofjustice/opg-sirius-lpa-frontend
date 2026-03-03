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

var ComplaintValueTranslationMap = map[string]map[string]string{
	"CATEGORY": {
		"01": "Correspondence",
		"02": "OPG Decisions",
		"03": "Non OPG",
		"04": "Customer Service",
		"05": "Policy",
	},
	"SUBCATEGORY": {
		"06": "General Query",
		"07": "Chase up",
		"08": "Typo / Grammar",
		"09": "Quality of Documents",
		"10": "Third Parties",
		"11": "Refund Request",
		"12": "Digital Tool",
		"13": "Finance",
		"14": "Customer Service",
		"15": "POA Decisions",
		"16": "Supervision Decisions",
		"17": "Investigation Outcomes",
		"18": "Fee Decision",
		"19": "Safeguarding Decisions",
		"20": "Other",
		"21": "Banks / Utilities",
		"22": "COP / Judicial",
		"23": "DX / Royal Mail",
		"24": "Health / Social Care",
		"25": "Solicitors",
		"26": "Deputy / Attorney",
		"27": "Other",
		"28": "Letter Content",
		"29": "Delays",
		"30": "Contact with OPG",
		"31": "Quality",
		"32": "Incorrect / Confusing Advice",
		"33": "Failure to Follow Procedure",
		"34": "Lost Documents",
		"35": "Security Breach",
		"36": "Other",
		"37": "Mental Capacity Act",
		"38": "Fee Policy",
		"39": "Donor Deceased Policy",
		"40": "Refund Policy",
		"41": "Forms / Guidance",
		"42": "Digital Product",
		"43": "Safeguarding Policy",
		"44": "Jurisdiction",
		"45": "Other",
	},
	"COMPLAINANTCATEGORY": {
		"LPA_DONOR":    "LPA donor",
		"LPA_ATTORNEY": "LPA attorney",
		"EPA_DONOR":    "EPA donor",
		"EPA_ATTORNEY": "EPA attorney",
		"SOLICITOR":    "Solicitor",
		"OTHER":        "Other",
	},
	"ORIGIN": {
		"CONTACT_CENTRE": "Contact centre",
		"LETTER":         "Letter",
		"EMAIL":          "E-mail",
		"PHONE":          "Phone call",
		"OTHER":          "Other",
	},
	"COMPENSATIONTYPE": {
		"REFUND":         "Refund",
		"FEE_WAIVER":     "Fee waiver",
		"EX_GRATIA":      "Ex gratia",
		"COMPENSATORY":   "Compensatory",
		"NOT_APPLICABLE": "N/A",
	},
}

func TranslateComplaintProperty(property string) string {
	if val, ok := ComplaintPropertyTranslationMap[strings.ToUpper(property)]; ok {
		return val
	}
	return property
}

func TranslateComplaintValue(property, value string) string {
	if val, ok := ComplaintValueTranslationMap[strings.ToUpper(property)][strings.ToUpper(value)]; ok {
		return val
	}
	return value
}
