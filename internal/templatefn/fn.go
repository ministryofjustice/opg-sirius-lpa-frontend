package templatefn

import (
	"fmt"
	"html/template"
	"net/url"
	"sort"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

func All(siriusPublicURL, prefix, staticHash string) map[string]interface{} {
	return map[string]interface{}{
		"sirius": func(s string) string {
			return siriusPublicURL + s
		},
		"prefix": func(s string) string {
			return prefix + s
		},
		"prefixAsset": func(s string) string {
			if len(staticHash) >= 11 {
				return prefix + s + "?" + url.QueryEscape(staticHash[3:11])
			} else {
				return prefix + s
			}
		},
		"today": func() string {
			return time.Now().Format("2006-01-02")
		},
		"field":                 field,
		"radios":                radios,
		"item":                  item,
		"fieldID":               fieldID,
		"select":                select_,
		"options":               options,
		"caseTabs":              caseTab,
		"casesWarningAppliedTo": sirius.CasesWarningAppliedTo,
		"fee": func(amount int) string {
			float := float64(amount)
			return fmt.Sprintf("%.2f", float/100)
		},
		"formatDate": func(s sirius.DateString) (string, error) {
			if s != "" {
				return s.ToSirius()
			}
			return "", nil
		},
		"date": func(s sirius.DateString, format string) (string, error) {
			if s == "" {
				return "", nil
			}

			t, err := s.Time()

			if err != nil {
				return "", err
			}

			return t.Format(format), nil
		},
		// s is a date string; layout specifies its structure;
		// if the date is invalid, this returns "invalid date"
		// instead of an error to prevent breaking the page render
		"parseAndFormatDate": func(s string, layout string, format string) string {
			if s == "" {
				return "invalid date"
			}

			t, err := time.Parse(layout, s)
			if err != nil {
				return "invalid date"
			}
			return t.Format(format)
		},
		"translateRefData": func(types []sirius.RefDataItem, tmplHandle string) string {
			for _, refDataType := range types {
				if refDataType.Handle == tmplHandle {
					return refDataType.Label
				}
			}
			return tmplHandle
		},
		"ToLower": strings.ToLower,
		"ToUpper": strings.ToUpper,
		"capitalise": func(text string) string {
			return cases.Title(language.English).String(text)
		},
		"camelcaseToSentence": func(text string) string {
			if text == "" {
				return ""
			}

			if text == "uId" {
				return "UID"
			}

			r, n := utf8.DecodeRuneInString(text)
			text = text[n:]

			s := string(unicode.ToUpper(r))
			for len(text) > 0 {
				r, n := utf8.DecodeRuneInString(text)
				text = text[n:]

				if r >= 'A' && r < 'Z' {
					s += " "
					s += string(unicode.ToLower(r))
				} else if r >= '0' && r < '9' {
					s += " "
					s += string(r)
				} else {
					s += string(r)
				}
			}

			return s
		},
		"contains": func(xs []string, needle string) bool {
			for _, x := range xs {
				if x == needle {
					return true
				}
			}
			return false
		},
		"plusN": func(i int, n int) int {
			return i + n
		},
		"replace": func(s, find, replace string) string {
			return strings.ReplaceAll(s, find, replace)
		},
		"dateYear": func(s sirius.DateString) (string, error) {
			if s != "" {
				return s.GetYear()
			}
			return "", nil
		},
		"filterContent": func(content string) string {
			//Fixes extra newline appearing in text editor due to newline present between the doctype and html tags
			return strings.ReplaceAll(content, "<!DOCTYPE html>\n<html lang=\"en\">", "<!DOCTYPE html><html lang=\"en\">")
		},
		"abs": func(num int) int {
			if num < 0 {
				return -num
			}
			return num
		},
		"attr": func(s string) template.HTMLAttr {
			return template.HTMLAttr(s) //#nosec G203 false positive
		},
		"join": func(s []string, joiner string) string {
			return strings.Join(s, joiner)
		},
		"subtypeShortFormat":     sirius.SubtypeShortFormat,
		"subtypeLongFormat":      subtypeLongFormat,
		"subtypeColour":          subtypeColour,
		"severanceRequiredLabel": severanceRequiredLabel,
		"howAttorneysMakeDecisionsLongForm": func(isSoleAttorney bool, s string) string {
			if isSoleAttorney {
				return "There is only one attorney appointed"
			}

			switch s {
			case "jointly":
				return "Jointly"
			case "jointly-and-severally":
				return "Jointly & severally"
			case "jointly-for-some-severally-for-others":
				return "Jointly for some, severally for others"
			case "":
				return "Not specified"
			default:
				return "howAttorneysMakeDecisions NOT RECOGNISED: " + s
			}
		},
		"howReplacementAttorneysStepInLongForm": func(s string) string {
			switch s {
			case "all-can-no-longer-act":
				return "When all can no longer act"
			case "one-can-no-longer-act":
				return "When one can no longer act"
			case "another-way":
				return "Another way"
			case "":
				return "Not specified"
			default:
				return "howReplacementAttorneysStepIn NOT RECOGNISED: " + s
			}
		},
		"whenTheLpaCanBeUsedLongForm": func(s string) string {
			switch s {
			case "when-has-capacity":
				return "As soon as it's registered"
			case "when-capacity-lost":
				return "When capacity is lost"
			case "":
				return "Not specified"
			default:
				return "whenTheLpaCanBeUsed NOT RECOGNISED: " + s
			}
		},
		"lifeSustainingTreatmentOptionLongForm": func(s string) string {
			switch s {
			case "option-a":
				return "Attorneys can give or refuse consent to LST"
			case "option-b":
				return "Attorneys cannot give or refuse consent to LST"
			case "":
				return "Not specified"
			default:
				return "lifeSustainingTreatmentOption NOT RECOGNISED: " + s
			}
		},
		// translate channel code to long version for Format fields in display
		"channelForFormat": func(s string) string {
			switch s {
			case "paper":
				return "Paper"
			case "online":
				return "Online"
			case "":
				return "Not specified"
			default:
				return "channel NOT RECOGNISED: " + s
			}
		},
		// translate language code to long version for Format fields in display
		"languageForFormat": func(s string) string {
			switch s {
			case "cy":
				return "Welsh"
			case "en":
				return "English"
			case "":
				return "Not specified"
			default:
				return "language NOT RECOGNISED: " + s
			}
		},
		// translate progress indicator context to long version for application progress page
		"progressIndicatorContext": func(s string) string {
			switch s {
			case "FEES":
				return "Fees"
			case "DONOR":
				return "Donor section"
			case "DONOR_ID":
				return "Donor identity confirmation"
			case "CERTIFICATE_PROVIDER_ID":
				return "Certificate provider identity confirmation"
			case "CERTIFICATE_PROVIDER_SIGNATURE":
				return "Certificate provider certificate"
			case "ATTORNEY_SIGNATURES":
				return "Attorney signatures"
			case "PREREGISTRATION_NOTICES":
				return "Pre-registration notices"
			case "REGISTRATION_NOTICES":
				return "Registration notices"
			case "RESTRICTIONS_AND_CONDITIONS":
				return "Restrictions and conditions"
			case "":
				return "Not specified"
			default:
				return "indicator NOT RECOGNISED: " + s
			}
		},
		// translate progress indicator status for application progress page
		"progressIndicatorStatus": func(s string) string {
			switch s {
			case "IN_PROGRESS":
				return "In progress"
			case "COMPLETE":
				return "Complete"
			case "CANNOT_START":
				return "Not started"
			case "":
				return "Not specified"
			default:
				return "status NOT RECOGNISED: " + s
			}
		},
		// translate objection type for confirm objection page
		"objectionType": func(s string) string {
			switch s {
			case "factual":
				return "Factual"
			case "prescribed":
				return "Prescribed"
			case "thirdParty":
				return "Third Party"
			case "":
				return "Not specified"
			default:
				return "objection type NOT RECOGNISED: " + s
			}
		},
		// translate resolution outcome for case summary page
		"resolutionOutcome": func(s string) string {
			switch s {
			case "upheld":
				return "upheld"
			case "notUpheld":
				return "not upheld"
			case "":
				return "Not specified"
			default:
				return "resolution outcome NOT RECOGNISED: " + s
			}
		},
		"compareBoolPointers": func(i *bool, j bool) bool {
			return *i == j
		},
		"inStringArray": func(value string, array []string) bool {
			for _, v := range array {
				if v == value {
					return true
				}
			}
			return false
		},
		"formatEventType": formatEventType,
		"eventTypeColor":  eventTypeColor,
	}
}

type CaseTabData struct {
	CaseSummary       sirius.CaseSummary
	SortedLinkedCases []linkedCase
	TabName           string
}

type linkedCase struct {
	UID         string
	Subtype     string
	Status      shared.CaseStatus
	CreatedDate sirius.DateString
}

// full text for LPA subtype, e.g. "Personal welfare"
func subtypeLongFormat(subtype string) string {
	switch strings.ToLower(subtype) {
	case "personal-welfare":
		return "Personal welfare"
	case "property-and-affairs":
		return "Property and affairs"
	case "hw":
		return "Health and welfare"
	case "pfa":
		return "Property and financial affairs"
	default:
		return ""
	}
}

func subtypeColour(subtype string) string {
	switch strings.ToLower(subtype) {
	case "personal-welfare":
		return "light-green"
	case "property-and-affairs":
		return "turquoise"
	default:
		return ""
	}
}

func severanceRequiredLabel(severanceStatus string) string {
	switch severanceStatus {
	case "REQUIRED":
		return "Yes"
	case "NOT_REQUIRED":
		return "No"
	default:
		return ""
	}
}

func caseTab(caseSummary sirius.CaseSummary, tabName string) CaseTabData {
	lpa := caseSummary.DigitalLpa.SiriusData
	lpaStore := caseSummary.DigitalLpa.LpaStoreData
	status := shared.CaseStatusTypeDraft

	if lpaStore.Status.ReadableString() != "" {
		status = lpaStore.Status
	}

	var linkedCases []linkedCase

	linkedCases = append(linkedCases, linkedCase{lpa.UID, lpa.Subtype, status, lpa.CreatedDate})

	for _, linkedLpa := range lpa.LinkedCases {
		linkedCases = append(linkedCases, linkedCase{linkedLpa.UID, linkedLpa.Subtype, linkedLpa.Status, linkedLpa.CreatedDate})
	}

	sort.Slice(linkedCases, func(i, j int) bool {
		if linkedCases[i].CreatedDate == linkedCases[j].CreatedDate {
			return linkedCases[i].UID < linkedCases[j].UID
		}
		return linkedCases[i].CreatedDate < linkedCases[j].CreatedDate
	})

	return CaseTabData{
		CaseSummary:       caseSummary,
		SortedLinkedCases: linkedCases,
		TabName:           tabName,
	}
}

type fieldData struct {
	Name  string
	Label string
	Value interface{}
	Error map[string]string
	Attrs map[string]interface{}
}

func field(name, label string, value interface{}, error map[string]string, attrs ...interface{}) fieldData {
	return fieldData{
		Name:  name,
		Label: label,
		Value: value,
		Error: error,
		Attrs: collectAttrs(attrs),
	}
}

type radiosData struct {
	Name   string
	Label  string
	Value  interface{}
	Errors map[string]string
	Items  []itemData
}

func radios(name, label string, value interface{}, errors map[string]string, items ...itemData) radiosData {
	return radiosData{
		Name:   name,
		Label:  label,
		Value:  value,
		Errors: errors,
		Items:  items,
	}
}

type itemData struct {
	Value string
	Label string
	Attrs map[string]interface{}
}

func item(value, label string, attrs ...interface{}) itemData {
	return itemData{
		Value: value,
		Label: label,
		Attrs: collectAttrs(attrs),
	}
}

func fieldID(name string, i int) string {
	if i == 0 {
		return name
	}

	return fmt.Sprintf("%s-%d", name, i+1)
}

type selectData struct {
	Name    string
	Label   string
	Value   interface{} // string|int
	Errors  map[string]string
	Options []optionData
	Attrs   map[string]interface{}
}

func select_(name, label string, value interface{}, errors map[string]string, options []optionData, attrs ...interface{}) selectData {
	return selectData{
		Name:    name,
		Label:   label,
		Value:   value,
		Errors:  errors,
		Options: options,
		Attrs:   collectAttrs(attrs),
	}
}

type optionData struct {
	Value interface{} // string|int
	Label string
}

func options(list interface{}, attrs ...interface{}) []optionData {
	attributes := collectAttrs(attrs)
	var datas []optionData

	switch v := list.(type) {
	case []string:
		datas = make([]optionData, len(v))
		for i, item := range v {
			datas[i] = optionData{Value: item, Label: item}
		}

	case []sirius.MiConfigEnum:
		datas = make([]optionData, len(v))
		for i, item := range v {
			datas[i] = optionData{Value: item.Name, Label: item.Description}
		}

	case []sirius.RefDataItem:
		if attributes["filterSelectable"] == true {
			for _, item := range v {
				if item.UserSelectable {
					datas = append(datas, optionData{Value: item.Handle, Label: item.Label})
				}
			}
		} else {
			datas = make([]optionData, len(v))
			for i, item := range v {
				datas[i] = optionData{Value: item.Handle, Label: item.Label}
			}
		}

	case []sirius.Team:
		datas = make([]optionData, len(v))
		for i, item := range v {
			datas[i] = optionData{Value: item.ID, Label: item.DisplayName}
		}
	}

	return datas
}

func collectAttrs(attrs []interface{}) map[string]interface{} {
	attributes := map[string]interface{}{}
	if len(attrs)%2 != 0 {
		panic("must have even number of attrs")
	}

	for i := 0; i < len(attrs); i += 2 {
		attributes[attrs[i].(string)] = attrs[i+1]
	}

	return attributes
}

var eventTypeColors = map[string]string{
	"CORRECTION":                "blue",
	"CHANGE_ATTORNEYS":          "yellow",
	"CERTIFICATE_PROVIDER_SIGN": "green",
	"DONOR_CONFIRM_IDENTITY":    "green",
	"OPG_STATUS_CHANGE":         "purple",
	"STATUTORY_WAITING_PERIOD":  "orange",
}

func eventTypeColor(eventType string) string {
	if color, ok := eventTypeColors[eventType]; ok {
		return color
	}
	return "grey"
}

func formatEventType(eventType string) string {
	// Replace underscores with spaces
	formatted := strings.ReplaceAll(eventType, "_", " ")
	// Title case each word
	return cases.Title(language.English).String(strings.ToLower(formatted))
}
