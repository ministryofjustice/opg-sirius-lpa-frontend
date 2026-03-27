package shared

import (
	"encoding/json"
)

type WhenLpaCanBeUsed int

const (
	WhenLpaCanBeUsedUnknown WhenLpaCanBeUsed = iota
	WhenLpaCanBeUsedHasCapacity
	WhenLpaCanBeUsedCapacityLost
)

type whenLpaCanBeUsedMeta struct {
	Readable string
	API      string
}

var whenLpaCanBeUsedMetadata = map[WhenLpaCanBeUsed]whenLpaCanBeUsedMeta{
	WhenLpaCanBeUsedHasCapacity:  {"As soon as it's registered", "when-has-capacity"},
	WhenLpaCanBeUsedCapacityLost: {"When capacity is lost", "when-capacity-lost"},
}

func (w WhenLpaCanBeUsed) ReadableString() string {
	if meta, ok := whenLpaCanBeUsedMetadata[w]; ok {
		return meta.Readable
	}
	return "Not specified"
}

func (w WhenLpaCanBeUsed) StringForApi() string {
	if meta, ok := whenLpaCanBeUsedMetadata[w]; ok {
		return meta.API
	}
	return ""
}

func ParseWhenLpaCanBeUsed(s string) WhenLpaCanBeUsed {
	for w, meta := range whenLpaCanBeUsedMetadata {
		if (meta.Readable != "" && meta.Readable == s) ||
			(meta.API != "" && meta.API == s) {
			return w
		}
	}
	return WhenLpaCanBeUsedUnknown
}

func (w WhenLpaCanBeUsed) MarshalJSON() ([]byte, error) {
	return json.Marshal(w.StringForApi())
}

func (w *WhenLpaCanBeUsed) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*w = ParseWhenLpaCanBeUsed(s)
	return nil
}

