package templatefn

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAndFormatDate(t *testing.T) {
	fns := All("", "", "")

	fn := fns["parseAndFormatDate"].(func(string, string, string)(string, error))

	var val string
	var err error

	val, err = fn("", "", "")
	assert.Equal(t, "", val)
	assert.Equal(t, errors.New("Not a date"), err)

	val, err = fn("2024-13-30", "2006-01-02", "2 January 2006")
	assert.Equal(t, "", val)
	assert.Equal(t, "parsing time \"2024-13-30\": month out of range", err.Error())

	val, err = fn("16 April 2024", "2 January 2006", "2006-01-02")
	assert.Equal(t, "2024-04-16", val)
	assert.Equal(t, nil, err)
}
