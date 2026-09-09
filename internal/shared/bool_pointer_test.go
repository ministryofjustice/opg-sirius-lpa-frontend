package shared

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoolPointer(t *testing.T) {
	truePtr := true
	falsePtr := false
	testCases := []struct {
		input    bool
		expected *bool
	}{
		{input: true, expected: &truePtr},
		{input: false, expected: &falsePtr},
	}

	for _, tc := range testCases {
		result := BoolPtr(tc.input)
		assert.Equal(t, tc.expected, result)
	}
}
