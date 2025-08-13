package templatefn

import (
	"reflect"
	"testing"
	"time"

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

	tests := map[string]map[string]interface{}{
		"Jointly":                          {"isSoleAttorney": false, "value": "jointly", "result": "Jointly"},
		"JointlyAndSeverally":              {"isSoleAttorney": false, "value": "jointly-and-severally", "result": "Jointly & severally"},
		"JointlyForSomeSeverallyForOthers": {"isSoleAttorney": false, "value": "jointly-for-some-severally-for-others", "result": "Jointly for some, severally for others"},
		"Empty":                            {"isSoleAttorney": false, "value": "", "result": "Not specified"},
		"NotValid":                         {"isSoleAttorney": false, "value": "foo", "result": "howAttorneysMakeDecisions NOT RECOGNISED: foo"},
		"IsSoleAttorney":                   {"isSoleAttorney": true, "value": "jointly-for-some-severally-for-others", "result": "There is only one attorney appointed"},
	}

	for _, test := range tests {
		assert.Equal(t, test["result"], fn(test["isSoleAttorney"].(bool), test["value"].(string)))
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

func TestSirius(t *testing.T) {
	fns := All("", "", "")
	fn := fns["sirius"].(func(string) string)

	val := fn("URL")
	assert.Equal(t, "URL", val)
}

func TestPrefix(t *testing.T) {
	fns := All("", "", "")
	fn := fns["prefix"].(func(string) string)

	val := fn("URL")
	assert.Equal(t, "URL", val)
}

func TestPrefixAssetGivenValidCache(t *testing.T) {
	fns := All("", "PREFIX", "248d6a61d20638b8e5c026930c3e6039a33ce45964ff2167f6ecedd419db06c1")
	fn := fns["prefixAsset"].(func(string) string)

	val := fn("")
	assert.Equal(t, "PREFIX?d6a61d20", val)
}

func TestPrefixAssetIgnoresInvalidCache(t *testing.T) {
	fns := All("", "PREFIX", "48d6a61d2")
	fn := fns["prefixAsset"].(func(string) string)

	val := fn("48d6a61d20638b8e5c026930c3e6039a33ce45964ff2167f6ecedd419db06c1")
	assert.Equal(t, "PREFIX48d6a61d20638b8e5c026930c3e6039a33ce45964ff2167f6ecedd419db06c1", val)
}

func TestToday(t *testing.T) {
	fns := All("", "", "")
	fn := fns["today"].(func() string)
	today := time.Now().Format("2006-01-02")
	val := fn()
	assert.Equal(t, today, val)
}

func TestField(t *testing.T) {
	fns := All("", "", "")
	fn := fns["field"].(func(string, string, interface{}, map[string]string, ...interface{}) fieldData)
	expected := fieldData{
		Name:  "name",
		Label: "Name",
		Value: "testing",
		Error: map[string]string{
			"username": "required",
		},
		Attrs: map[string]interface{}{},
	}
	val := fn("name", "Name", "testing", map[string]string{"username": "required"})
	assert.Equal(t, expected, val)
}

func TestRadios(t *testing.T) {
	items := []itemData{{Value: "foo", Label: "Foo"}}
	fns := All("", "", "")
	fn := fns["radios"].(func(string, string, interface{}, map[string]string, ...itemData) radiosData)
	expected := radiosData{
		Name:  "name",
		Label: "Name",
		Value: "testing",
		Errors: map[string]string{
			"username": "required",
		},
		Items: items,
	}

	val := fn("name", "Name", "testing", map[string]string{"username": "required"}, items...)
	assert.Equal(t, expected, val)
}

func TestItem(t *testing.T) {
	fns := All("", "", "")
	fn := fns["item"].(func(string, string, ...interface{}) itemData)
	expected := itemData{
		Value: "testing",
		Label: "Name",
		Attrs: map[string]interface{}{},
	}

	val := fn("testing", "Name")
	assert.Equal(t, expected, val)
}

func TestFieldID(t *testing.T) {
	fns := All("", "", "")
	fn := fns["fieldID"].(func(string, int) string)

	val := fn("testing", 52)
	assert.Equal(t, "testing-53", val)
}

func TestSelect(t *testing.T) {
	fns := All("", "", "")
	fn := fns["select"].(func(string, string, interface{}, map[string]string, []optionData, ...interface{}) selectData)
	expected := selectData{
		Name:  "name",
		Label: "Name",
		Value: "testing",
		Errors: map[string]string{
			"username": "required",
		},
		Options: []optionData{{Value: "foo", Label: "Foo"}},
		Attrs:   map[string]interface{}{},
	}

	val := fn("name", "Name", "testing", map[string]string{"username": "required"}, []optionData{{Value: "foo", Label: "Foo"}})
	assert.Equal(t, expected, val)
}

func TestOptionsStringSlice(t *testing.T) {
	input := []string{"A", "B"}
	expected := []optionData{
		{Value: "A", Label: "A"},
		{Value: "B", Label: "B"},
	}

	result := options(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("got %v, want %v", result, expected)
	}
}

func TestOptionsMiConfigEnumSlice(t *testing.T) {
	input := []sirius.MiConfigEnum{
		{Name: "opt1", Description: "Option 1"},
		{Name: "opt2", Description: "Option 2"},
	}
	expected := []optionData{
		{Value: "opt1", Label: "Option 1"},
		{Value: "opt2", Label: "Option 2"},
	}

	result := options(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("got %v, want %v", result, expected)
	}
}

func TestOptionsRefDataItemSliceWithoutFilterSelectable(t *testing.T) {
	input := []sirius.RefDataItem{
		{Handle: "H1", Label: "L1", UserSelectable: true},
		{Handle: "H2", Label: "L2", UserSelectable: false},
	}
	expected := []optionData{
		{Value: "H1", Label: "L1"},
		{Value: "H2", Label: "L2"},
	}

	result := options(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("got %v, want %v", result, expected)
	}
}

func TestOptionsRefDataItemSliceWithFilterSelectable(t *testing.T) {
	input := []sirius.RefDataItem{
		{Handle: "H1", Label: "L1", UserSelectable: true},
		{Handle: "H2", Label: "L2", UserSelectable: false},
	}
	expected := []optionData{
		{Value: "H1", Label: "L1"},
	}

	result := options(input, "filterSelectable", true)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("got %v, want %v", result, expected)
	}
}

func TestOptionsTeamSlice(t *testing.T) {
	input := []sirius.Team{
		{ID: 1, DisplayName: "Team One"},
		{ID: 2, DisplayName: "Team Two"},
	}
	expected := []optionData{
		{Value: 1, Label: "Team One"},
		{Value: 2, Label: "Team Two"},
	}

	result := options(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("got %v, want %v", result, expected)
	}
}

func TestCaseTab(t *testing.T) {
	caseSummary := sirius.CaseSummary{
		DigitalLpa: sirius.DigitalLpa{
			UID: "123",
			SiriusData: sirius.SiriusData{
				ID:          222,
				UID:         "454654",
				CreatedDate: "2010-01-01",
				Subtype:     "subType",
			},
			LpaStoreData: sirius.LpaStoreData{
				Status: "draft",
			},
		},
	}

	expected := CaseTabData{
		CaseSummary: caseSummary,
		SortedLinkedCases: []linkedCase{
			{
				UID:         "454654",
				Subtype:     "subType",
				Status:      "Draft",
				CreatedDate: "2010-01-01",
			},
		},
		TabName: "TabName",
	}

	val := caseTab(caseSummary, "TabName")
	assert.Equal(t, expected, val)
}

func TestSortWarningsForCaseSummary(t *testing.T) {
	warnings := []sirius.Warning{
		{
			WarningType: "Donor Deceased",
			DateAdded:   "01/01/2020 00:02:03",
		},
		{
			WarningType: "Welsh Language",
			DateAdded:   "11/07/2012 11:02:03",
		},
		{
			WarningType: "Safeguarding",
			DateAdded:   "20/02/2016 00:02:03",
		},
		{
			WarningType: "Violent Warning",
			DateAdded:   "15/09/2011 00:02:03",
		},
		{
			WarningType: "Fee Issue",
			DateAdded:   "11/07/2012 11:02:02",
		},
	}

	expected := []sirius.Warning{
		{
			WarningType: "Donor Deceased",
			DateAdded:   "01/01/2020 00:02:03",
		},
		{
			WarningType: "Safeguarding",
			DateAdded:   "20/02/2016 00:02:03",
		},
		{
			WarningType: "Welsh Language",
			DateAdded:   "11/07/2012 11:02:03",
		},
		{
			WarningType: "Fee Issue",
			DateAdded:   "11/07/2012 11:02:02",
		},
		{
			WarningType: "Violent Warning",
			DateAdded:   "15/09/2011 00:02:03",
		},
	}

	val := sortWarningsForCaseSummary(warnings)
	assert.Equal(t, expected, val)
}

func TestCasesWarningAppliedToOnlyOneCase(t *testing.T) {
	cases := []sirius.Case{
		{
			UID:     "UID123String",
			SubType: "pfa",
		},
	}

	expected := ""

	val := casesWarningAppliedTo("UID123String", cases)
	assert.Equal(t, expected, val)
}

func TestCasesWarningAppliedToCasesWithSameUID(t *testing.T) {
	cases := []sirius.Case{
		{
			UID:     "UID123String",
			SubType: "pfa",
		},
		{
			UID:     "UID123String",
			SubType: "hw",
		},
	}

	expected := ""

	val := casesWarningAppliedTo("UID123String", cases)
	assert.Equal(t, expected, val)
}

func TestCasesWarningAppliedToManyDifferentCases(t *testing.T) {
	cases := []sirius.Case{
		{
			UID:     "UID123String",
			SubType: "pfa",
		},
		{
			UID:     "UID123StringPFANotMatching",
			SubType: "pfa",
		},
		{
			UID:     "UID123StringNtMatching",
			SubType: "hw",
		},
	}

	expected := ", PFA UID123StringPFANotMatching and HW UID123StringNtMatching"

	val := casesWarningAppliedTo("UID123String", cases)
	assert.Equal(t, expected, val)
}

func TestFee(t *testing.T) {
	fns := All("", "", "")
	fn := fns["fee"].(func(int) string)
	expected := "82.00"

	val := fn(8200)
	assert.Equal(t, expected, val)
}

func TestFormatDateForAnEmptyDate(t *testing.T) {
	fns := All("", "", "")
	fn := fns["formatDate"].(func(sirius.DateString) (string, error))
	expected := ""

	val, err := fn("")
	assert.Nil(t, err)
	assert.Equal(t, expected, val)
}

func TestFormatDate(t *testing.T) {
	fns := All("", "", "")
	fn := fns["formatDate"].(func(sirius.DateString) (string, error))
	var expected string
	var err error
	var val string

	expected = "03/02/2024"
	val, err = fn("2024-02-03")
	assert.Nil(t, err)
	assert.Equal(t, expected, val)

	expected = "failed to format non-date"
	val, err = fn("202402/03")
	assert.Equal(t, expected, err.Error())
	assert.Equal(t, "", val)
}

func TestTranslateRefDataForTheLabel(t *testing.T) {
	fns := All("", "", "")
	fn := fns["translateRefData"].(func([]sirius.RefDataItem, string) string)

	types := []sirius.RefDataItem{
		{
			Handle: "REFTYPE",
			Label:  "refType",
		},
	}

	val := fn(types, "REFTYPE")
	assert.Equal(t, "refType", val)
}

func TestTranslateRefDataForTheTmplHandle(t *testing.T) {
	fns := All("", "", "")
	fn := fns["translateRefData"].(func([]sirius.RefDataItem, string) string)

	types := []sirius.RefDataItem{
		{
			Handle: "REFTYPE",
			Label:  "refType",
		},
	}

	val := fn(types, "tmplHandle")
	assert.Equal(t, "tmplHandle", val)
}

func TestCapitalise(t *testing.T) {
	fns := All("", "", "")
	fn := fns["capitalise"].(func(string) string)

	var val string

	val = fn("REFTYPE")
	assert.Equal(t, "Reftype", val)
	val = fn("reftype")
	assert.Equal(t, "Reftype", val)
	val = fn("")
	assert.Equal(t, "", val)
}

func TestContains(t *testing.T) {
	fns := All("", "", "")
	fn := fns["contains"].(func([]string, string) bool)

	var val bool

	val = fn([]string{"a", "b", "c"}, "b")
	assert.Equal(t, true, val)

	val = fn([]string{"a", "b", "c"}, "d")
	assert.Equal(t, false, val)
}

func TestStatusColour(t *testing.T) {
	fns := All("", "", "")
	fn := fns["statusColour"].(func(string) string)

	var val string

	val = fn("RegIStered")
	assert.Equal(t, "green", val)

	val = fn("IN PROGRESS")
	assert.Equal(t, "light-blue", val)

	val = fn("return - unpaid")
	assert.Equal(t, "red", val)

	val = fn("not in list")
	assert.Equal(t, "grey", val)
}

func TestStatusLabel(t *testing.T) {
	var val string

	val = StatusLabelFormat("DRAFT")
	assert.Equal(t, "Draft", val)

	val = StatusLabelFormat("in-progress")
	assert.Equal(t, "In progress", val)

	val = StatusLabelFormat("not in list")
	assert.Equal(t, "draft", val)
}

func TestReplace(t *testing.T) {
	fns := All("", "", "")
	fn := fns["replace"].(func(string, string, string) string)

	val := fn("oink oink oink", "oink", "moo")
	assert.Equal(t, "moo moo moo", val)
}

func TestDateYear(t *testing.T) {
	fns := All("", "", "")
	fn := fns["dateYear"].(func(sirius.DateString) (string, error))

	val, err := fn("")
	assert.Nil(t, err)
	assert.Equal(t, "", val)

	val, _ = fn("2006-01-02")
	assert.Equal(t, "2006", val)

	expected := "failed to format non-date"
	_, err = fn("202402/03")
	assert.Equal(t, expected, err.Error())
}

func TestFilterContent(t *testing.T) {
	fns := All("", "", "")
	fn := fns["filterContent"].(func(string) string)

	val := fn("<>Testing<!DOCTYPE html>\n<html lang=\"en\">!@£$%^&*()<>")
	assert.Equal(t, "<>Testing<!DOCTYPE html><html lang=\"en\">!@£$%^&*()<>", val)
}

func TestAbs(t *testing.T) {
	fns := All("", "", "")
	fn := fns["abs"].(func(int) int)
	var val int

	val = fn(0)
	assert.Equal(t, 0, val)

	val = fn(-1)
	assert.Equal(t, 1, val)

	val = fn(1)
	assert.Equal(t, 1, val)
}

func TestJoin(t *testing.T) {
	fns := All("", "", "")
	fn := fns["join"].(func([]string, string) string)

	val := fn([]string{"a", "b", "c"}, "-")
	assert.Equal(t, "a-b-c", val)
}

func TestSubtypeShortFormat(t *testing.T) {
	var val string
	val = subtypeShortFormat("hw")
	assert.Equal(t, "HW", val)

	val = subtypeShortFormat("not-in-list")
	assert.Equal(t, "", val)
}

func TestSubtypeLongFormat(t *testing.T) {
	var val string
	val = subtypeLongFormat("hw")
	assert.Equal(t, "Health and welfare", val)

	val = subtypeLongFormat("not-in-list")
	assert.Equal(t, "", val)
}

func TestSubtypeColour(t *testing.T) {
	var val string
	val = subtypeColour("personal-welfare")
	assert.Equal(t, "light-green", val)

	val = subtypeColour("not-in-list")
	assert.Equal(t, "", val)
}
