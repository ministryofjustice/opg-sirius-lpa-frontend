package sirius

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDateString(t *testing.T) {
	testcases := map[string]string{
		"normal":  `"03/04/2022"`,
		"escaped": `"03\/04\/2022"`,
	}

	for name, fromSirius := range testcases {
		t.Run(name, func(t *testing.T) {
			var v DateString
			err := json.Unmarshal([]byte(fromSirius), &v)
			assert.Nil(t, err)
			assert.Equal(t, "2022-04-03", string(v))

			data, err := json.Marshal(v)
			assert.Nil(t, err)
			assert.Equal(t, `"03/04/2022"`, string(data))
		})
	}
}

func TestDateStringNull(t *testing.T) {
	fromSirius := `null`

	var v DateString
	err := json.Unmarshal([]byte(fromSirius), &v)
	assert.Nil(t, err)
	assert.Equal(t, "", string(v))

	data, err := json.Marshal(v)
	assert.Nil(t, err)
	assert.Equal(t, fromSirius, string(data))
}

func TestDateStringEmpty(t *testing.T) {
	fromSirius := `""`

	var v DateString
	err := json.Unmarshal([]byte(fromSirius), &v)
	assert.Nil(t, err)
	assert.Equal(t, "", string(v))

	data, err := json.Marshal(v)
	assert.Nil(t, err)
	assert.Equal(t, "null", string(data))
}

func TestDateStringErrors(t *testing.T) {
	fromSirius := `"03/04"`

	var v DateString
	err := json.Unmarshal([]byte(fromSirius), &v)
	assert.NotNil(t, err)

	_, err = json.Marshal(DateString("2022-03"))
	assert.NotNil(t, err)
}
