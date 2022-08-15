package sirius

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFeeString(t *testing.T) {
	testcases := map[string]string{
		"2350": "23.50",
		"8200": "82.00",
		"545":  "5.45",
		"0":    "0.00",
	}

	for pence, decimal := range testcases {
		t.Run(pence, func(t *testing.T) {
			var v FeeString
			err := json.Unmarshal([]byte(pence), &v)
			assert.Nil(t, err)
			assert.Equal(t, decimal, string(v))
		})
	}
}

func TestFeeStringEmpty(t *testing.T) {
	var f FeeString
	err := json.Unmarshal([]byte(`""`), &f)
	assert.Nil(t, err)
	assert.Equal(t, "", string(f))
}
