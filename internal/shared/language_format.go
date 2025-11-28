package shared

import (
	"encoding/json"
)

type LanguageFormat int

const (
	LanguageFormatCy LanguageFormat = iota
	LanguageFormatEn
	LanguageFormatEmpty
	LanguageFormatNotRecognised
)

var languageFormatMap = map[string]LanguageFormat{
	"cy":            LanguageFormatCy,
	"en":            LanguageFormatEn,
	"":              LanguageFormatEmpty,
	"notRecognised": LanguageFormatNotRecognised,
}

func (l LanguageFormat) String() string {
	return l.Key()
}

func (l LanguageFormat) Translation() string {
	switch l {
	case LanguageFormatCy:
		return "Welsh"
	case LanguageFormatEn:
		return "English"
	case LanguageFormatEmpty:
		return "Not specified"
	default:
		return "language NOT RECOGNISED: " + l.String()
	}
}

func (l LanguageFormat) Key() string {
	switch l {
	case LanguageFormatCy:
		return "Cy"
	case LanguageFormatEn:
		return "En"
	case LanguageFormatEmpty:
		return "Empty"
	case LanguageFormatNotRecognised:
		return "Not Recognised"
	default:
		return ""
	}
}

func ParseLanguageFormat(s string) LanguageFormat {
	value, ok := languageFormatMap[s]
	if !ok {
		return LanguageFormatNotRecognised
	}
	return value
}

func (l LanguageFormat) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.Key())
}

func (l *LanguageFormat) UnmarshalJSON(data []byte) (err error) {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*l = ParseLanguageFormat(s)
	return nil
}
