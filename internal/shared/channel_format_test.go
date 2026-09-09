package shared

import (
	"testing"
)

func TestParseChannelFormat(t *testing.T) {
	tests := []struct {
		input    string
		expected ChannelFormat
	}{
		{"paper", ChannelFormatPaper},
		{"online", ChannelFormatOnline},
		{"", ChannelFormatEmpty},
		{"notRecognised", ChannelFormatNotRecognised},
	}

	for _, tt := range tests {
		got := ParseChannelFormat(tt.input)
		if got != tt.expected {
			t.Errorf("ParseChannelFormat(%q) = %v, want %v", tt.input, got, tt.expected)
		}
	}
}

func TestChannelFormatTranslation(t *testing.T) {
	tests := []struct {
		name  string
		input ChannelFormat
		want  string
	}{
		{
			name:  "Paper",
			input: ChannelFormatPaper,
			want:  "Paper",
		},
		{
			name:  "Online",
			input: ChannelFormatOnline,
			want:  "Online",
		},
		{
			name:  "Not specified",
			input: ChannelFormatEmpty,
			want:  "Not specified",
		},
		{
			name:  "Unrecognised value",
			input: ChannelFormatNotRecognised,
			want:  "channel NOT RECOGNISED: Not Recognised",
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

func TestChannelForFormat(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"paper", "Paper"},
		{"online", "Online"},
		{"", "Not specified"},
		{"unknown", "channel NOT RECOGNISED: Not Recognised"},
	}

	for _, tt := range tests {
		got := ChannelForFormat(tt.input)
		if got != tt.want {
			t.Errorf("ChannelForFormat(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
