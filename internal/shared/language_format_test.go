package shared

import (
	"encoding/json"
	"testing"
)

func TestParseLanguageFormat(t *testing.T) {
	tests := []struct {
		input    string
		expected LanguageFormat
	}{
		{"cy", LanguageFormatCy},
		{"en", LanguageFormatEn},
		{"", LanguageFormatEmpty},
		{"notRecognised", LanguageFormatNotRecognised},
	}

	for _, tt := range tests {
		got := ParseLanguageFormat(tt.input)
		if got != tt.expected {
			t.Errorf("ParseLanguageFormat(%q) = %v, want %v", tt.input, got, tt.expected)
		}
	}
}

func TestLanguageFormat_Translation(t *testing.T) {
	tests := []struct {
		name  string
		input LanguageFormat
		want  string
	}{
		{
			name:  "Welsh",
			input: LanguageFormatCy,
			want:  "Welsh",
		},
		{
			name:  "English",
			input: LanguageFormatEn,
			want:  "English",
		},
		{
			name:  "Not specified",
			input: LanguageFormatEmpty,
			want:  "Not specified",
		},
		{
			name:  "Unrecognised value",
			input: LanguageFormatNotRecognised,
			want:  "language NOT RECOGNISED: Not Recognised",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.Translation()
			if got != tt.want {
				t.Fatalf("Translation() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLanguageFormatUnmarshalJSON(t *testing.T) {
	tests := []struct {
		jsonInput string
		expected  LanguageFormat
	}{
		{`"cy"`, LanguageFormatCy},
		{`"en"`, LanguageFormatEn},
		{`""`, LanguageFormatEmpty},
		{`"notRecognised"`, LanguageFormatNotRecognised},
	}

	for _, tt := range tests {
		var cs LanguageFormat
		err := json.Unmarshal([]byte(tt.jsonInput), &cs)
		if err != nil {
			t.Errorf("UnmarshalJSON(%s) returned error: %v", tt.jsonInput, err)
		}
		if cs != tt.expected {
			t.Errorf("UnmarshalJSON(%s) = %v, want %v", tt.jsonInput, cs, tt.expected)
		}
	}
}
