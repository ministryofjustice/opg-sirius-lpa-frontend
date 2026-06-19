package sirius

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func getCountries() interface{} {
	return []interface{}{
		map[string]string{
			"Handle": "UK",
			"Label":  "United Kingdom",
		},
	}
}

func TestCacheWhenEmpty(t *testing.T) {
	val, ok := getCached("not set")

	assert.Nil(t, val)
	assert.Equal(t, ok, false)
}

func TestCacheWhenSet(t *testing.T) {
	setCached("countries", getCountries())
	val, ok := getCached("countries")

	assert.Equal(t, val, getCountries())
	assert.Equal(t, ok, true)
}

func TestCacheWhenMiss(t *testing.T) {
	setCached("countries", getCountries())
	val, ok := getCached("cities")

	assert.Nil(t, val)
	assert.Equal(t, ok, false)
}
