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

	for pence, pounds := range testcases {
		t.Run(pence, func(t *testing.T) {
			var v FeeString
			err := json.Unmarshal([]byte(pence), &v)
			assert.Nil(t, err)
			assert.Equal(t, pounds, string(v))

			s, err := v.ToPence()
			assert.Nil(t, err)
			assert.Equal(t, pence, s)

			data, err := json.Marshal(v)

			assert.Equal(t, `"`+pence+`"`, string(data))
			assert.Nil(t, err)
		})
	}
}

func TestFeeStringEmpty(t *testing.T) {
	var v FeeString
	err := json.Unmarshal([]byte(`""`), &v)
	assert.Nil(t, err)
	assert.Equal(t, "", string(v))
}

func TestFeeStringNull(t *testing.T) {
	var v FeeString
	err := json.Unmarshal([]byte(`null`), &v)
	assert.Nil(t, err)
	assert.Equal(t, "", string(v))
}

func TestFeeStringErrors(t *testing.T) {
	var v FeeString
	err := json.Unmarshal([]byte(`"hello"`), &v)
	assert.NotNil(t, err)

	_, err = json.Marshal(FeeString("hello"))
	assert.NotNil(t, err)
}
