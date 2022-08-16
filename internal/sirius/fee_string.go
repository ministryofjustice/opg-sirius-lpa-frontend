package sirius

import (
	"bytes"
	"fmt"
	"strconv"
)

type FeeString string

func (f *FeeString) UnmarshalJSON(amount []byte) error {
	if bytes.Equal([]byte("null"), amount) || bytes.Equal([]byte(`""`), amount) {
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

func (f FeeString) MarshalJSON() ([]byte, error) {
	if string(f) == "" {
		return []byte(`null`), nil
	}

	amount, err := f.ToPence()
	if err != nil {
		return nil, err
	}

	return []byte(`"` + amount + `"`), nil
}

func (f FeeString) ToPence() (string, error) {
	float, err := strconv.ParseFloat(string(f), 64)
	if err != nil {
		return "", err
	}

	pence := float * 100

	return fmt.Sprintf("%.0f", pence), nil
}
