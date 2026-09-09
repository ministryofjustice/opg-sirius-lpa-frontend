package shared
import (
"encoding/json"
"testing"
)
func TestParseLifeSustainingTreatmentOption(t *testing.T) {
	tests := []struct {
		input    string
		expected LifeSustainingTreatmentOption
	}{
		{"option-a", LifeSustainingTreatmentOptionA},
		{"option-b", LifeSustainingTreatmentOptionB},
		{"", LifeSustainingTreatmentOptionEmpty},
		{"unknown-value", LifeSustainingTreatmentOptionNotRecognised},
	}
	for _, tt := range tests {
		got := ParseLifeSustainingTreatmentOption(tt.input)
		if got != tt.expected {
			t.Errorf("ParseLifeSustainingTreatmentOption(%q) = %v, want %v", tt.input, got, tt.expected)
		}
	}
}
func TestLifeSustainingTreatmentOptionLongForm(t *testing.T) {
	tests := []struct {
		input    LifeSustainingTreatmentOption
		expected string
	}{
		{LifeSustainingTreatmentOptionA, "Attorneys can give or refuse consent to LST"},
		{LifeSustainingTreatmentOptionB, "Attorneys cannot give or refuse consent to LST"},
		{LifeSustainingTreatmentOptionEmpty, "Not specified"},
		{LifeSustainingTreatmentOptionNotRecognised, "lifeSustainingTreatmentOption NOT RECOGNISED: notRecognised"},
	}
	for _, tt := range tests {
		got := tt.input.LongForm()
		if got != tt.expected {
			t.Errorf("LongForm() = %q, want %q", got, tt.expected)
		}
	}
}
func TestLifeSustainingTreatmentOptionKey(t *testing.T) {
	tests := []struct {
		input    LifeSustainingTreatmentOption
		expected string
	}{
		{LifeSustainingTreatmentOptionA, "option-a"},
		{LifeSustainingTreatmentOptionB, "option-b"},
		{LifeSustainingTreatmentOptionEmpty, ""},
		{LifeSustainingTreatmentOptionNotRecognised, "notRecognised"},
	}
	for _, tt := range tests {
		got := tt.input.Key()
		if got != tt.expected {
			t.Errorf("Key() = %q, want %q", got, tt.expected)
		}
	}
}
func TestLifeSustainingTreatmentOptionString(t *testing.T) {
	tests := []struct {
		input    LifeSustainingTreatmentOption
		expected string
	}{
		{LifeSustainingTreatmentOptionA, "option-a"},
		{LifeSustainingTreatmentOptionB, "option-b"},
		{LifeSustainingTreatmentOptionEmpty, ""},
		{LifeSustainingTreatmentOptionNotRecognised, "notRecognised"},
	}
	for _, tt := range tests {
		got := tt.input.String()
		if got != tt.expected {
			t.Errorf("String() = %q, want %q", got, tt.expected)
		}
	}
}
func TestLifeSustainingTreatmentOptionUnmarshalJSON(t *testing.T) {
	tests := []struct {
		jsonInput string
		expected  LifeSustainingTreatmentOption
	}{
		{`"option-a"`, LifeSustainingTreatmentOptionA},
		{`"option-b"`, LifeSustainingTreatmentOptionB},
		{`""`, LifeSustainingTreatmentOptionEmpty},
		{`"unknown"`, LifeSustainingTreatmentOptionNotRecognised},
	}
	for _, tt := range tests {
		var l LifeSustainingTreatmentOption
		err := json.Unmarshal([]byte(tt.jsonInput), &l)
		if err != nil {
			t.Errorf("UnmarshalJSON(%s) returned error: %v", tt.jsonInput, err)
		}
		if l != tt.expected {
			t.Errorf("UnmarshalJSON(%s) = %v, want %v", tt.jsonInput, l, tt.expected)
		}
	}
}
func TestLifeSustainingTreatmentOptionMarshalJSON(t *testing.T) {
	tests := []struct {
		input    LifeSustainingTreatmentOption
		expected string
	}{
		{LifeSustainingTreatmentOptionA, `"option-a"`},
		{LifeSustainingTreatmentOptionB, `"option-b"`},
		{LifeSustainingTreatmentOptionEmpty, `""`},
		{LifeSustainingTreatmentOptionNotRecognised, `"notRecognised"`},
	}
	for _, tt := range tests {
		got, err := json.Marshal(tt.input)
		if err != nil {
			t.Errorf("MarshalJSON() returned error: %v", err)
		}
		if string(got) != tt.expected {
			t.Errorf("MarshalJSON() = %s, want %s", string(got), tt.expected)
		}
	}
}
