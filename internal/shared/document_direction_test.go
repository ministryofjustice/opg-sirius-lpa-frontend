package shared

import (
	"encoding/json"
	"testing"
)

func TestParseDocumentDirection(t *testing.T) {
	tests := []struct {
		input    string
		expected DocumentDirection
	}{
		{"Incoming", DocumentDirectionIn},
		{"Outgoing", DocumentDirectionOut},
		{"", DocumentDirectionEmpty},
		{"notRecognised", DocumentDirectionNotRecognised},
	}

	for _, tc := range tests {
		got := ParseDocumentDirection(tc.input)
		if got != tc.expected {
			t.Errorf("ParseDocumentDirection(%q) = %v, want %v", tc.input, got, tc.expected)
		}
	}
}

func TestDocumentDirectionTranslation(t *testing.T) {
	tests := []struct {
		name  string
		input DocumentDirection
		want  string
	}{
		{
			name:  "Incoming",
			input: DocumentDirectionIn,
			want:  "In",
		},
		{
			name:  "Outgoing",
			input: DocumentDirectionOut,
			want:  "Out",
		},
		{
			name:  "Not specified",
			input: DocumentDirectionEmpty,
			want:  "Not specified",
		},
		{
			name:  "Unrecognised value",
			input: DocumentDirectionNotRecognised,
			want:  "document direction NOT RECOGNISED: Not Recognised",
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

func TestDocumentDirectionKey(t *testing.T) {
	tests := []struct {
		name  string
		input DocumentDirection
		want  string
	}{
		{
			name:  "Incoming",
			input: DocumentDirectionIn,
			want:  "Incoming",
		},
		{
			name:  "Outgoing",
			input: DocumentDirectionOut,
			want:  "Outgoing",
		},
		{
			name:  "Empty",
			input: DocumentDirectionEmpty,
			want:  "Empty",
		},
		{
			name:  "Not recognised",
			input: DocumentDirectionNotRecognised,
			want:  "Not Recognised",
		},
		{
			name:  "Unrecognised value",
			input: DocumentDirection(4),
			want:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.input.Key()
			if got != tc.want {
				t.Fatalf("Key() = %q, want %q", got, tc.want)
			}
		})
	}

}

func TestDocumentDirectionUnmarshalJSON(t *testing.T) {
	tests := []struct {
		jsonInput string
		expected  DocumentDirection
	}{
		{`"Incoming"`, DocumentDirectionIn},
		{`"Outgoing"`, DocumentDirectionOut},
		{`""`, DocumentDirectionEmpty},
		{`"notRecognised"`, DocumentDirectionNotRecognised},
	}

	for _, tt := range tests {
		var cs DocumentDirection
		err := json.Unmarshal([]byte(tt.jsonInput), &cs)
		if err != nil {
			t.Errorf("UnmarshalJSON(%s) returned error: %v", tt.jsonInput, err)
		}
		if cs != tt.expected {
			t.Errorf("UnmarshalJSON(%s) = %v, want %v", tt.jsonInput, cs, tt.expected)
		}
	}
}
