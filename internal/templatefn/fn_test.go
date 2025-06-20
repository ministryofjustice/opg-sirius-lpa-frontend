package templatefn

import (
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"

	"github.com/stretchr/testify/assert"
)

func TestDate(t *testing.T) {
	fns := All("", "", "")
	fn := fns["date"].(func(sirius.DateString, string) (string, error))

	val, err := fn("", "")
	assert.Nil(t, err)
	assert.Equal(t, "", val)

	val, err = fn("aaaa", "2006-01-02")
	assert.NotNil(t, err)
	assert.Equal(t, "", val)

	val, err = fn("2024-06-27", "02 Jan 2006")
	assert.Nil(t, err)
	assert.Equal(t, "27 Jun 2024", val)
}

func TestParseAndFormatDate(t *testing.T) {
	fns := All("", "", "")
	fn := fns["parseAndFormatDate"].(func(string, string, string) string)

	var val string

	val = fn("", "", "")
	assert.Equal(t, "invalid date", val)

	val = fn("2024-13-30", "2006-01-02", "2 January 2006")
	assert.Equal(t, "invalid date", val)

	val = fn("16 April 2024", "2 January 2006", "2006-01-02")
	assert.Equal(t, "2024-04-16", val)

	val = fn("2024-04-11T14:00:39.361141055Z", "2006-01-02T15:04:05Z", "2 January 2006")
	assert.Equal(t, "11 April 2024", val)
}

func TestPlusN(t *testing.T) {
	fns := All("", "", "")
	fn := fns["plusN"].(func(int, int) int)

	val := fn(1, 1)
	assert.Equal(t, 2, val)
}

func testStringMapper(t *testing.T, fnName string, expectations map[string]string) {
	fns := All("", "", "")
	fn := fns[fnName].(func(string) string)

	for input, expected := range expectations {
		assert.Equal(t, expected, fn(input))
	}
}

func TestHowAttorneysMakeDecisionsLongForm(t *testing.T) {
	fns := All("", "", "")
	fn := fns["howAttorneysMakeDecisionsLongForm"].(func(bool, string) string)

	expectations := map[int]map[string]interface{}{
		0: {"soleAttorney": false, "value": "jointly", "result": "Jointly"},
		1: {"soleAttorney": false, "value": "jointly-and-severally", "result": "Jointly & severally"},
		2: {"soleAttorney": false, "value": "jointly-for-some-severally-for-others", "result": "Jointly for some, severally for others"},
		3: {"soleAttorney": false, "value": "", "result": "Not specified"},
		4: {"soleAttorney": false, "value": "foo", "result": "howAttorneysMakeDecisions NOT RECOGNISED: foo"},
		5: {"soleAttorney": true, "value": "jointly-for-some-severally-for-others", "result": "There is only one attorney appointed"},
	}

	for _, expectation := range expectations {
		assert.Equal(t, expectation["result"], fn(expectation["soleAttorney"].(bool), expectation["value"].(string)))
	}
}

func TestHowReplacementAttorneysStepInLongForm(t *testing.T) {
	expectations := map[string]string{
		"all-can-no-longer-act": "When all can no longer act",
		"one-can-no-longer-act": "When one can no longer act",
		"another-way":           "Another way",
		"":                      "Not specified",
		"foo":                   "howReplacementAttorneysStepIn NOT RECOGNISED: foo",
	}

	testStringMapper(t, "howReplacementAttorneysStepInLongForm", expectations)
}

func TestLifeSustainingTreatmentOptionLongForm(t *testing.T) {
	expectations := map[string]string{
		"option-a": "Attorneys can give or refuse consent to LST",
		"option-b": "Attorneys cannot give or refuse consent to LST",
		"":         "Not specified",
		"foo":      "lifeSustainingTreatmentOption NOT RECOGNISED: foo",
	}

	testStringMapper(t, "lifeSustainingTreatmentOptionLongForm", expectations)
}

func TestWhenTheLpaCanBeUsedLongForm(t *testing.T) {
	expectations := map[string]string{
		"when-has-capacity":  "As soon as it's registered",
		"when-capacity-lost": "When capacity is lost",
		"":                   "Not specified",
		"foo":                "whenTheLpaCanBeUsed NOT RECOGNISED: foo",
	}

	testStringMapper(t, "whenTheLpaCanBeUsedLongForm", expectations)
}

func TestChannelForFormat(t *testing.T) {
	expectations := map[string]string{
		"paper":  "Paper",
		"online": "Online",
		"":       "Not specified",
		"foo":    "channel NOT RECOGNISED: foo",
	}

	testStringMapper(t, "channelForFormat", expectations)
}

func TestLanguageForFormat(t *testing.T) {
	expectations := map[string]string{
		"cy":  "Welsh",
		"en":  "English",
		"":    "Not specified",
		"foo": "language NOT RECOGNISED: foo",
	}

	testStringMapper(t, "languageForFormat", expectations)
}

func TestProgressIndicatorContext(t *testing.T) {
	expectations := map[string]string{
		"FEES":                           "Fees",
		"DONOR":                          "Donor section",
		"DONOR_ID":                       "Donor identity confirmation",
		"CERTIFICATE_PROVIDER_ID":        "Certificate provider identity confirmation",
		"CERTIFICATE_PROVIDER_SIGNATURE": "Certificate provider certificate",
		"ATTORNEY_SIGNATURES":            "Attorney signatures",
		"PREREGISTRATION_NOTICES":        "Pre-registration notices",
		"REGISTRATION_NOTICES":           "Registration notices",
		"":                               "Not specified",
		"foo":                            "indicator NOT RECOGNISED: foo",
	}

	testStringMapper(t, "progressIndicatorContext", expectations)
}

func TestProgressIndicatorStatus(t *testing.T) {
	expectations := map[string]string{
		"IN_PROGRESS":  "In progress",
		"CANNOT_START": "Not started",
		"COMPLETE":     "Complete",
		"":             "Not specified",
		"foo":          "status NOT RECOGNISED: foo",
	}

	testStringMapper(t, "progressIndicatorStatus", expectations)
}

func TestObjectionType(t *testing.T) {
	expectations := map[string]string{
		"factual":    "Factual",
		"prescribed": "Prescribed",
		"thirdParty": "Third Party",
		"":           "Not specified",
		"foo":        "objection type NOT RECOGNISED: foo",
	}

	testStringMapper(t, "objectionType", expectations)
}

func TestResolutionOutcome(t *testing.T) {
	expectations := map[string]string{
		"upheld":    "upheld",
		"notUpheld": "not upheld",
		"":          "Not specified",
		"foo":       "resolution outcome NOT RECOGNISED: foo",
	}

	testStringMapper(t, "resolutionOutcome", expectations)
}

func TestCamelcaseToSentence(t *testing.T) {
	expectations := map[string]string{
		"uId":   "UID",
		"abc":   "Abc",
		"aBc":   "A bc",
		"aBCd":  "A b cd",
		"aBcDe": "A bc de",
		"a2B1":  "A 2 b 1",
		"":      "",
	}

	testStringMapper(t, "camelcaseToSentence", expectations)
}

func TestSeveranceRequiredLabel(t *testing.T) {
	expectations := map[string]string{
		"REQUIRED":     "Yes",
		"NOT_REQUIRED": "No",
		"":             "",
	}

	testStringMapper(t, "severanceRequiredLabel", expectations)
}

// Helper function to easily get a pointer to a bool
func boolPtr(b bool) *bool {
	return &b
}

func TestCompareBoolPointer(t *testing.T) {
	fns := All("", "", "")
	fn := fns["compareBoolPointers"].(func(*bool, bool) bool)

	var val bool

	val = fn(boolPtr(true), true)
	assert.Equal(t, true, val)

	val = fn(boolPtr(true), false)
	assert.Equal(t, false, val)

	val = fn(boolPtr(false), false)
	assert.Equal(t, true, val)
}

func TestInStringArray(t *testing.T) {
	fns := All("", "", "")
	fn := fns["inStringArray"].(func(string, []string) bool)

	var val bool

	val = fn("in-progress", []string{"draft", "in-progress"})
	assert.Equal(t, true, val)

	val = fn("cancelled", []string{"draft", "in-progress"})
	assert.Equal(t, false, val)
}
