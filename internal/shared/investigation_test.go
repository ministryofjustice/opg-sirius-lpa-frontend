package shared

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTranslateInvestigationEventProperty(t *testing.T) {
	tests := map[string]string{
		"ADDITIONALINFORMATION": "Information",
		"123":                   "123",
	}

	for input, expected := range tests {
		result := TranslateInvestigationEventProperty(input)
		assert.Equal(t, expected, result, "TranslateInvestigationEventProperty(%q) should return %q", input, expected)
	}
}
