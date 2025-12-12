package shared

import (
	"encoding/json"
	"testing"
)

func TestParseCaseLabel(t *testing.T) {
	tests := []struct {
		input    string
		expected CaseLabel
	}{
		{"hw", CaseLabelHW},
		{"pfa", CaseLabelPFA},
		{"", CaseLabelEmpty},
		{"notRecognised", CaseLabelNotRecognised},
	}

	for _, tc := range tests {
		got := ParseCaseLabel(tc.input)
		if got != tc.expected {
			t.Errorf("ParseCaseLabel(%q) = %v, want %v", tc.input, got, tc.expected)
		}
	}
}

func TestCaseLabelTranslation(t *testing.T) {
	tests := []struct {
		name  string
		input CaseLabel
		want  string
	}{
		{
			name:  "hw",
			input: CaseLabelHW,
			want:  "colour-govuk-grass-green",
		},
		{
			name:  "pfa",
			input: CaseLabelPFA,
			want:  "colour-govuk-turquoise",
		},
		{
			name:  "Not specified",
			input: CaseLabelEmpty,
			want:  "Not specified",
		},
		{
			name:  "Unrecognised value",
			input: CaseLabelNotRecognised,
			want:  "case label NOT RECOGNISED: Not Recognised",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.input.Translation()
			if got != tc.want {
				t.Fatalf("Translation() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestCaseLabelUnmarshalJSON(t *testing.T) {
	tests := []struct {
		jsonInput string
		expected  CaseLabel
	}{
		{`"hw"`, CaseLabelHW},
		{`"pfa"`, CaseLabelPFA},
		{`""`, CaseLabelEmpty},
		{`"notRecognised"`, CaseLabelNotRecognised},
	}

	for _, tt := range tests {
		var cs CaseLabel
		err := json.Unmarshal([]byte(tt.jsonInput), &cs)
		if err != nil {
			t.Errorf("UnmarshalJSON(%s) returned error: %v", tt.jsonInput, err)
		}
		if cs != tt.expected {
			t.Errorf("UnmarshalJSON(%s) = %v, want %v", tt.jsonInput, cs, tt.expected)
		}
	}
}
