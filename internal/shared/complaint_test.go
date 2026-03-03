package shared

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTranslateComplaintProperty(t *testing.T) {
	tests := map[string]string{
		"CATEGORY": "Category",
		"123":      "123",
	}

	for input, expected := range tests {
		result := TranslateComplaintProperty(input)
		assert.Equal(t, expected, result, "TranslateComplaintProperty(%q) should return %q", input, expected)
	}
}

func TestTranslateComplaintValue(t *testing.T) {
	tests := map[string]map[string]string{
		"CATEGORY": {
			"01":  "Correspondence",
			"123": "123",
		},
		"234": {
			"06": "06",
		},
	}

	for property, test := range tests {
		for value, expected := range test {
			result := TranslateComplaintValue(property, value)
			assert.Equal(t, expected, result, "TranslateComplaintValue(%q, %q) should return %q", property, value, expected)
		}
	}
}
