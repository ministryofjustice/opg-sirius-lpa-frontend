package sirius

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDateString(t *testing.T) {
	fromSirius := `"03/04/2022"`

	var v DateString
	err := json.Unmarshal([]byte(fromSirius), &v)
	assert.Nil(t, err)
	assert.Equal(t, "2022-04-03", string(v))

	data, err := json.Marshal(v)
	assert.Nil(t, err)
	assert.Equal(t, fromSirius, string(data))
}

func TestDateStringErrors(t *testing.T) {
	fromSirius := `"03/04"`

	var v DateString
	err := json.Unmarshal([]byte(fromSirius), &v)
	assert.NotNil(t, err)

	_, err = json.Marshal(DateString("2022-03"))
	assert.NotNil(t, err)
}
