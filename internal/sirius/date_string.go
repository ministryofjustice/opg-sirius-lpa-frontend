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

func (s *DateString) UnmarshalText(text []byte) error {
	parts := bytes.Split(text, []byte{'/'})

	if len(parts) != 3 {
		return errors.New("failed to unmarshal non-date")
	}

	*s = DateString(fmt.Sprintf("%s-%s-%s", parts[2], parts[1], parts[0]))
	return nil
}

func (s DateString) MarshalText() ([]byte, error) {
	parts := strings.Split(string(s), "-")
	if len(parts) != 3 {
		return nil, errors.New("failed to marshal non-date")
	}

	return []byte(fmt.Sprintf("%s/%s/%s", parts[2], parts[1], parts[0])), nil
}
