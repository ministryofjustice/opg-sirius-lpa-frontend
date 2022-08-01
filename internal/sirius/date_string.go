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
	if bytes.Equal([]byte("null"), text) || bytes.Equal([]byte(`""`), text) {
		*s = DateString("")
		return nil
	}

	if text[0] != '"' || text[len(text)-1] != '"' {
		return errors.New("failed to unmarshal non-date")
	}
	text = text[1 : len(text)-1]

	// Sirius gives dates as "03\/04\/2022", which pointlessly escapes the forward
	// slashes, we can safely remove them
	text = bytes.ReplaceAll(text, []byte{'\\'}, []byte{})

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

	date, err := s.ToSirius()
	if err != nil {
		return nil, err
	}

	return []byte(`"` + date + `"`), nil
}

func (s DateString) ToSirius() (string, error) {
	parts := strings.Split(string(s), "-")
	if len(parts) != 3 {
		return "", errors.New("failed to format non-date")
	}

	return fmt.Sprintf(`%s/%s/%s`, parts[2], parts[1], parts[0]), nil
}
