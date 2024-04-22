package templatefn

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAndFormatDate(t *testing.T) {
	fns := All("", "", "")
	fn := fns["parseAndFormatDate"].(func(string, string, string)(string, error))

	var val string
	var err error

	val, err = fn("", "", "")
	assert.Equal(t, "", val)
	assert.Equal(t, errors.New("Not a date"), err)

	val, err = fn("2024-13-30", "2006-01-02", "2 January 2006")
	assert.Equal(t, "", val)
	assert.Equal(t, "parsing time \"2024-13-30\": month out of range", err.Error())

	val, err = fn("16 April 2024", "2 January 2006", "2006-01-02")
	assert.Equal(t, "2024-04-16", val)
	assert.Equal(t, nil, err)
}

func TestPlusN(t *testing.T) {
	fns := All("", "", "")
	fn := fns["plusN"].(func(int, int)int)

	val := fn(1, 1)
	assert.Equal(t, 2, val)
}

func testStringMapper(t *testing.T, fnName string, expectations map[string]string) {
	fns := All("", "", "")
	fn := fns[fnName].(func(string)string)

	for input, expected := range expectations {
		assert.Equal(t, expected, fn(input))
	}
}

func TestHowAttorneysMakeDecisionsLongForm(t *testing.T) {
	expectations := map[string]string{
		"jointly": "Jointly",
		"jointly-and-severally": "Jointly & severally",
		"jointly-for-some-severally-for-others": "Jointly for some, severally for others",
		"": "Not specified",
		"foo": "howAttorneysMakeDecisions NOT RECOGNISED: foo",
	}

	testStringMapper(t, "howAttorneysMakeDecisionsLongForm", expectations)
}

func TestHowReplacementAttorneysStepInLongForm(t *testing.T) {
	expectations := map[string]string{
		"all-can-no-longer-act": "When all can no longer act",
		"one-can-no-longer-act": "When one can no longer act",
		"another-way": "Another way",
		"": "Not specified",
		"foo": "howReplacementAttorneysStepIn NOT RECOGNISED: foo",
	}

	testStringMapper(t, "howReplacementAttorneysStepInLongForm", expectations)
}

func TestWhenTheLpaCanBeUsedLongForm(t *testing.T) {
	expectations := map[string]string{
		"when-has-capacity": "As soon as it's registered",
		"when-capacity-lost": "When capacity is lost",
		"": "Not specified",
		"foo": "whenTheLpaCanBeUsed NOT RECOGNISED: foo",
	}

	testStringMapper(t, "whenTheLpaCanBeUsedLongForm", expectations)
}

func TestChannelForFormat(t *testing.T) {
	expectations := map[string]string{
		"paper": "Paper",
		"online": "Online",
		"": "Not specified",
		"foo": "channel NOT RECOGNISED: foo",
	}

	testStringMapper(t, "channelForFormat", expectations)
}

func TestLanguageForFormat(t *testing.T) {
	expectations := map[string]string{
		"cy": "Welsh",
		"en": "English",
		"": "Not specified",
		"foo": "language NOT RECOGNISED: foo",
	}

	testStringMapper(t, "languageForFormat", expectations)
}
