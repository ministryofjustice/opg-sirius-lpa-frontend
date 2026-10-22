package shared

import (
	"encoding/json"
	"testing"
)

func TestParseDecisionType(t *testing.T) {
	tests := []struct {
		input    string
		expected DecisionType
	}{
		{"DECLINED_EXEMPTION", DecisionTypeExemption},
		{"DECLINED_REMISSION", DecisionTypeRemission},
		{"notRecognised", DecisionTypeNotRecognised},
		{"unknown-decision", DecisionTypeNotRecognised},
		{"", DecisionTypeNotRecognised},
	}

	for _, tt := range tests {
		got := ParseDecisionType(tt.input)
		if got != tt.expected {
			t.Errorf("ParseDecisionType(%q) = %v, want %v", tt.input, got, tt.expected)
		}
	}
}

func TestDecisionTypeString(t *testing.T) {
	tests := []struct {
		decisionType DecisionType
		expected     string
	}{
		{DecisionTypeExemption, "DECLINED_EXEMPTION"},
		{DecisionTypeRemission, "DECLINED_REMISSION"},
		{DecisionTypeNotRecognised, "DECLINED_NOT_RECOGNISED"},
	}

	for _, tt := range tests {
		got := tt.decisionType.String()
		if got != tt.expected {
			t.Errorf("String(%v) = %s, want %s", tt.decisionType, got, tt.expected)
		}
	}
}

func TestDecisionTypeKey(t *testing.T) {
	tests := []struct {
		decisionType DecisionType
		expected     string
	}{
		{DecisionTypeExemption, "DECLINED_EXEMPTION"},
		{DecisionTypeRemission, "DECLINED_REMISSION"},
		{DecisionTypeNotRecognised, "DECLINED_NOT_RECOGNISED"},
	}

	for _, tt := range tests {
		got := tt.decisionType.Key()
		if got != tt.expected {
			t.Errorf("Key(%v) = %s, want %s", tt.decisionType, got, tt.expected)
		}
	}
}

func TestDecisionTypeTranslation(t *testing.T) {
	tests := []struct {
		decisionType DecisionType
		expected     string
	}{
		{DecisionTypeExemption, "Exemption declined"},
		{DecisionTypeRemission, "Remission declined"},
		{DecisionTypeNotRecognised, "Not recognised"},
	}

	for _, tt := range tests {
		got := tt.decisionType.Translation()
		if got != tt.expected {
			t.Errorf("Translation(%v) = %s, want %s", tt.decisionType, got, tt.expected)
		}
	}
}

func TestDecisionTypeMarshalJSON(t *testing.T) {
	tests := []struct {
		decisionType DecisionType
		expected     string
	}{
		{DecisionTypeExemption, `"DECLINED_EXEMPTION"`},
		{DecisionTypeRemission, `"DECLINED_REMISSION"`},
		{DecisionTypeNotRecognised, `"DECLINED_NOT_RECOGNISED"`},
	}

	for _, tt := range tests {
		got, err := tt.decisionType.MarshalJSON()
		if err != nil {
			t.Errorf("MarshalJSON(%v) returned error: %v", tt.decisionType, err)
		}
		if string(got) != tt.expected {
			t.Errorf("MarshalJSON(%v) = %s, want %s", tt.decisionType, string(got), tt.expected)
		}
	}
}

func TestDecisionTypeUnmarshalJSON(t *testing.T) {
	tests := []struct {
		jsonInput string
		expected  DecisionType
	}{
		{`"DECLINED_EXEMPTION"`, DecisionTypeExemption},
		{`"DECLINED_REMISSION"`, DecisionTypeRemission},
		{`"notRecognised"`, DecisionTypeNotRecognised},
		{`"unknown"`, DecisionTypeNotRecognised},
	}

	for _, tt := range tests {
		var dt DecisionType
		err := json.Unmarshal([]byte(tt.jsonInput), &dt)
		if err != nil {
			t.Errorf("UnmarshalJSON(%s) returned error: %v", tt.jsonInput, err)
		}
		if dt != tt.expected {
			t.Errorf("UnmarshalJSON(%s) = %v, want %v", tt.jsonInput, dt, tt.expected)
		}
	}
}

func TestDecisionTypeJSONRoundTrip(t *testing.T) {
	tests := []DecisionType{
		DecisionTypeExemption,
		DecisionTypeRemission,
		DecisionTypeNotRecognised,
	}

	for _, tt := range tests {
		// Marshal to JSON
		marshalled, err := json.Marshal(tt)
		if err != nil {
			t.Errorf("Marshal(%v) returned error: %v", tt, err)
		}

		// Unmarshal back
		var dt DecisionType
		err = json.Unmarshal(marshalled, &dt)
		if err != nil {
			t.Errorf("Unmarshal(%s) returned error: %v", string(marshalled), err)
		}

		if dt != tt {
			t.Errorf("Round trip failed for %v: got %v", tt, dt)
		}
	}
}
