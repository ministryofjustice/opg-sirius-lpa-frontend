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
