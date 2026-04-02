package shared

import (
	"encoding/json"
	"testing"
)

func TestParseWhenLpaCanBeUsed(t *testing.T) {
	tests := []struct {
		input    string
		expected WhenLpaCanBeUsed
	}{
		{"As soon as it's registered", WhenLpaCanBeUsedHasCapacity},
		{"when-has-capacity", WhenLpaCanBeUsedHasCapacity},
		{"When capacity is lost", WhenLpaCanBeUsedCapacityLost},
		{"when-capacity-lost", WhenLpaCanBeUsedCapacityLost},
		{"unknown-value", WhenLpaCanBeUsedUnknown},
	}

	for _, tt := range tests {
		got := ParseWhenLpaCanBeUsed(tt.input)
		if got != tt.expected {
			t.Errorf("ParseWhenLpaCanBeUsed(%q) = %v, want %v", tt.input, got, tt.expected)
		}
	}
}

func TestWhenLpaCanBeUsedReadableString(t *testing.T) {
	tests := []struct {
		input    WhenLpaCanBeUsed
		expected string
	}{
		{WhenLpaCanBeUsedHasCapacity, "As soon as it's registered"},
		{WhenLpaCanBeUsedCapacityLost, "When capacity is lost"},
		{WhenLpaCanBeUsedUnknown, "Not specified"},
	}

	for _, tt := range tests {
		got := tt.input.ReadableString()
		if got != tt.expected {
			t.Errorf("ReadableString() = %q, want %q", got, tt.expected)
		}
	}
}

func TestWhenLpaCanBeUsedStringForApi(t *testing.T) {
	tests := []struct {
		input    WhenLpaCanBeUsed
		expected string
	}{
		{WhenLpaCanBeUsedHasCapacity, "when-has-capacity"},
		{WhenLpaCanBeUsedCapacityLost, "when-capacity-lost"},
		{WhenLpaCanBeUsedUnknown, ""},
	}

	for _, tt := range tests {
		got := tt.input.StringForApi()
		if got != tt.expected {
			t.Errorf("StringForApi() = %q, want %q", got, tt.expected)
		}
	}
}

func TestWhenLpaCanBeUsedUnmarshalJSON(t *testing.T) {
	tests := []struct {
		jsonInput string
		expected  WhenLpaCanBeUsed
	}{
		{`"when-has-capacity"`, WhenLpaCanBeUsedHasCapacity},
		{`"when-capacity-lost"`, WhenLpaCanBeUsedCapacityLost},
		{`""`, WhenLpaCanBeUsedUnknown},
		{`"unknown"`, WhenLpaCanBeUsedUnknown},
	}

	for _, tt := range tests {
		var w WhenLpaCanBeUsed
		err := json.Unmarshal([]byte(tt.jsonInput), &w)
		if err != nil {
			t.Errorf("UnmarshalJSON(%s) returned error: %v", tt.jsonInput, err)
		}
		if w != tt.expected {
			t.Errorf("UnmarshalJSON(%s) = %v, want %v", tt.jsonInput, w, tt.expected)
		}
	}
}

func TestWhenLpaCanBeUsedMarshalJSON(t *testing.T) {
	tests := []struct {
		input    WhenLpaCanBeUsed
		expected string
	}{
		{WhenLpaCanBeUsedHasCapacity, `"when-has-capacity"`},
		{WhenLpaCanBeUsedCapacityLost, `"when-capacity-lost"`},
		{WhenLpaCanBeUsedUnknown, `""`},
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

