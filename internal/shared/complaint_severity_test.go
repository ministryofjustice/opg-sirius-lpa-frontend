package shared

import (
	"encoding/json"
	"testing"
)

func TestParseComplaintSeverity(t *testing.T) {
	tests := []struct {
		input    string
		expected ComplaintSeverity
	}{
		{"Minor", ComplaintSeverityMinor},
		{"Major", ComplaintSeverityMajor},
		{"Security Breach", ComplaintSeveritySecurityBreach},
		{"NotRecognised", ComplaintSeverityNotRecognised},
		{"unknown-severity", ComplaintSeverityNotRecognised},
	}

	for _, tt := range tests {
		got := ParseComplaintSeverity(tt.input)
		if got != tt.expected {
			t.Errorf("ParseComplaintSeverity(%q) = %v, want %v", tt.input, got, tt.expected)
		}
	}
}

func TestComplaintSeverityUnmarshalJSON(t *testing.T) {
	tests := []struct {
		jsonInput string
		expected  ComplaintSeverity
	}{
		{`"Minor"`, ComplaintSeverityMinor},
		{`"Major"`, ComplaintSeverityMajor},
		{`"Security Breach"`, ComplaintSeveritySecurityBreach},
		{`"NotRecognised"`, ComplaintSeverityNotRecognised},
		{`"Invalid"`, ComplaintSeverityNotRecognised},
	}

	for _, tt := range tests {
		var cs ComplaintSeverity
		err := json.Unmarshal([]byte(tt.jsonInput), &cs)
		if err != nil {
			t.Errorf("UnmarshalJSON(%s) returned error: %v", tt.jsonInput, err)
		}
		if cs != tt.expected {
			t.Errorf("UnmarshalJSON(%s) = %v, want %v", tt.jsonInput, cs, tt.expected)
		}
	}
}

func TestComplaintSeverityMarshalJSON(t *testing.T) {
	tests := []struct {
		input    ComplaintSeverity
		expected string
	}{
		{ComplaintSeverityMinor, `"Minor"`},
		{ComplaintSeverityMajor, `"Major"`},
		{ComplaintSeveritySecurityBreach, `"Security Breach"`},
		{ComplaintSeverityNotRecognised, `"complaint severity NOT RECOGNISED"`},
	}

	for _, tt := range tests {
		got, err := json.Marshal(tt.input)
		if err != nil {
			t.Errorf("MarshalJSON(%v) returned error: %v", tt.input, err)
		}
		if string(got) != tt.expected {
			t.Errorf("MarshalJSON(%v) = %s, want %s", tt.input, string(got), tt.expected)
		}
	}
}

func TestComplaintSeverityTranslation(t *testing.T) {
	tests := []struct {
		severity ComplaintSeverity
		expected string
	}{
		{ComplaintSeverityMinor, "Minor"},
		{ComplaintSeverityMajor, "Major"},
		{ComplaintSeveritySecurityBreach, "Security Breach"},
		{ComplaintSeverityNotRecognised, "complaint severity NOT RECOGNISED"},
	}

	for _, tt := range tests {
		got := tt.severity.Translation()
		if got != tt.expected {
			t.Errorf("Translation(%v) = %s, want %s", tt.severity, got, tt.expected)
		}
	}
}

