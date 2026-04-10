package shared

import (
	"encoding/json"
)

type ComplaintSeverity int

const (
	ComplaintSeverityMinor ComplaintSeverity = iota
	ComplaintSeverityMajor
	ComplaintSeveritySecurityBreach
	ComplaintSeverityNotRecognised
)

var complaintSeverityMap = map[string]ComplaintSeverity{
	"Minor":           ComplaintSeverityMinor,
	"Major":           ComplaintSeverityMajor,
	"Security Breach": ComplaintSeveritySecurityBreach,
	"NotRecognised":   ComplaintSeverityNotRecognised,
}

func (d ComplaintSeverity) Translation() string {
	switch d {
	case ComplaintSeverityMinor:
		return "Minor"
	case ComplaintSeverityMajor:
		return "Major"
	case ComplaintSeveritySecurityBreach:
		return "Security Breach"
	default:
		return "complaint severity NOT RECOGNISED"
	}
}

func ParseComplaintSeverity(s string) ComplaintSeverity {
	value, ok := complaintSeverityMap[s]
	if !ok {
		return ComplaintSeverityNotRecognised
	}
	return value
}

func (d ComplaintSeverity) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Translation())
}

func (d *ComplaintSeverity) UnmarshalJSON(data []byte) (err error) {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*d = ParseComplaintSeverity(s)
	return nil
}
