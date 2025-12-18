package shared

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
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
		assert.Equal(t, tc.expected, got)
	}
}

func TestParseDocumentDirectionUnknown(t *testing.T) {
	got := ParseDocumentDirection("invalid")
	assert.Equal(t, DocumentDirectionNotRecognised, got)
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
			assert.Equal(t, tc.want, got)
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
			assert.Equal(t, tc.want, got)
		})
	}

}

func TestDocumentDirectionMarshalJSON(t *testing.T) {
	tests := []struct {
		name  string
		input DocumentDirection
		want  string
	}{
		{"Incoming", DocumentDirectionIn, `"Incoming"`},
		{"Outgoing", DocumentDirectionOut, `"Outgoing"`},
		{"Empty", DocumentDirectionEmpty, `"Empty"`},
		{"NotRecognised", DocumentDirectionNotRecognised, `"Not Recognised"`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b, _ := json.Marshal(tc.input)
			assert.Equal(t, string(b), tc.want)
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
		var d DocumentDirection
		err := json.Unmarshal([]byte(tt.jsonInput), &d)
		assert.Nil(t, err)
		assert.Equal(t, tt.expected, d)
	}
}

func TestDocumentDirectionUnmarshalJSONErrors(t *testing.T) {
	var d DocumentDirection
	err := json.Unmarshal([]byte(`123`), &d)
	assert.Error(t, err)
}
