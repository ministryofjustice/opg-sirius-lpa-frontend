package sirius

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

// DateString is a date in the format "YYYY-MM-DD" that will unmarshal from and
// marshal to the Sirius format of "DD/MM/YYYY"
type DateString string

func (s *DateString) UnmarshalJSON(text []byte) error {
	if bytes.Equal([]byte("null"), text) {
		*s = DateString("")
		return nil
	}

	if text[0] != '"' || text[len(text)-1] != '"' {
		return errors.New("failed to unmarshal non-date")
	}
	text = text[1 : len(text)-1]

	parts := bytes.Split(text, []byte{'/'})

	if len(parts) != 3 {
		return errors.New("failed to unmarshal non-date")
	}

	*s = DateString(fmt.Sprintf("%s-%s-%s", parts[2], parts[1], parts[0]))
	return nil
}

func (s DateString) MarshalJSON() ([]byte, error) {
	if string(s) == "" {
		return []byte(`null`), nil
	}

	parts := strings.Split(string(s), "-")
	if len(parts) != 3 {
		return nil, errors.New("failed to marshal non-date")
	}

	return []byte(fmt.Sprintf(`"%s/%s/%s"`, parts[2], parts[1], parts[0])), nil
}
