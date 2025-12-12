package shared

import (
	"encoding/json"
)

type DocumentDirection int

const (
	DocumentDirectionIn DocumentDirection = iota
	DocumentDirectionOut
	DocumentDirectionEmpty
	DocumentDirectionNotRecognised
)

var documentDirectionMap = map[string]DocumentDirection{
	"Incoming":      DocumentDirectionIn,
	"Outgoing":      DocumentDirectionOut,
	"":              DocumentDirectionEmpty,
	"notRecognised": DocumentDirectionNotRecognised,
}

func (d DocumentDirection) String() string {
	return d.Key()
}

func (d DocumentDirection) Translation() string {
	switch d {
	case DocumentDirectionIn:
		return "In"
	case DocumentDirectionOut:
		return "Out"
	case DocumentDirectionEmpty:
		return "Not specified"
	default:
		return "document direction NOT RECOGNISED: " + d.String()
	}
}

func (d DocumentDirection) Key() string {
	switch d {
	case DocumentDirectionIn:
		return "Incoming"
	case DocumentDirectionOut:
		return "Outgoing"
	case DocumentDirectionEmpty:
		return "Empty"
	case DocumentDirectionNotRecognised:
		return "Not Recognised"
	default:
		return ""
	}
}

func ParseDocumentDirection(s string) DocumentDirection {
	value, ok := documentDirectionMap[s]
	if !ok {
		return DocumentDirectionNotRecognised
	}
	return value
}

func (d DocumentDirection) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Key())
}

func (d *DocumentDirection) UnmarshalJSON(data []byte) (err error) {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*d = ParseDocumentDirection(s)
	return nil
}
