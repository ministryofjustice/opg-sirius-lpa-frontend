package shared

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatMonetaryValue(t *testing.T) {
	expected := "82.00"

	val := FormatMonetaryValue(8200)
	assert.Equal(t, expected, val)
}

func TestFormatMonetaryFloat(t *testing.T) {
	expected := "82.00"

	val := FormatMonetaryFloat(float64(8200))
	assert.Equal(t, expected, val)
}
