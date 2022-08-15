package sirius

import (
	"bytes"
	"fmt"
	"strconv"
)

type FeeString string

func (f *FeeString) UnmarshalJSON(amount []byte) error {
	if bytes.Equal([]byte(`""`), amount) {
		*f = ""
		return nil
	}
	float, err := strconv.ParseFloat(string(amount), 32)
	if err != nil {
		return err
	}
	*f = FeeString(fmt.Sprintf("%.2f", float/100))
	return nil
}
