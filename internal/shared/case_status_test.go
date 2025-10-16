package shared

import (
	"encoding/json"
	"testing"
)

func TestParseCaseStatusType(t *testing.T) {
	tests := []struct {
		input    string
		expected CaseStatus
	}{
		{"Draft", CaseStatusTypeDraft},
		{"draft", CaseStatusTypeDraft},
		{"in-progress", CaseStatusTypeInProgress},
		{"unknown-status", CaseStatusTypeUnknown},
	}

	for _, tt := range tests {
		got := ParseCaseStatusType(tt.input)
		if got != tt.expected {
			t.Errorf("ParseCaseStatusType(%q) = %v, want %v", tt.input, got, tt.expected)
		}
	}
}

func TestUnmarshalJSON(t *testing.T) {
	tests := []struct {
		jsonInput string
		expected  CaseStatus
	}{
		{`"Draft"`, CaseStatusTypeDraft},
		{`"draft"`, CaseStatusTypeDraft},
		{`"Registered"`, CaseStatusTypeRegistered},
		{`"Invalid"`, CaseStatusTypeInvalid},
	}

	for _, tt := range tests {
		var cs CaseStatus
		err := json.Unmarshal([]byte(tt.jsonInput), &cs)
		if err != nil {
			t.Errorf("UnmarshalJSON(%s) returned error: %v", tt.jsonInput, err)
		}
		if cs != tt.expected {
			t.Errorf("UnmarshalJSON(%s) = %v, want %v", tt.jsonInput, cs, tt.expected)
		}
	}
}

func TestColour(t *testing.T) {
	tests := []struct {
		status   CaseStatus
		expected string
	}{
		{CaseStatusTypeRegistered, "green"},
		{CaseStatusTypePerfect, "turquoise"},
		{CaseStatusTypeStatutoryWaitingPeriod, "yellow"},
		{CaseStatusTypeInProgress, "light-blue"},
		{CaseStatusTypeDraft, "purple"},
		{CaseStatusTypeCancelled, "red"},
		{CaseStatusTypeUnknown, "grey"},
	}

	for _, tt := range tests {
		got := tt.status.Colour()
		if got != tt.expected {
			t.Errorf("Colour(%v) = %s, want %s", tt.status, got, tt.expected)
		}
	}
}

func TestIsValidStatusForObjection(t *testing.T) {
	tests := []struct {
		status   CaseStatus
		expected bool
	}{
		{CaseStatusTypeInProgress, true},
		{CaseStatusTypeDraft, true},
		{CaseStatusTypeStatutoryWaitingPeriod, true},
		{CaseStatusTypeRegistered, false},
		{CaseStatusTypeCancelled, false},
		{CaseStatusTypeUnknown, false},
	}

	for _, tt := range tests {
		got := tt.status.IsValidStatusForObjection()
		if got != tt.expected {
			t.Errorf("IsValidStatusForObjection(%v) = %v, want %v", tt.status, got, tt.expected)
		}
	}
}

func TestIsDraft(t *testing.T) {
	tests := []struct {
		status   CaseStatus
		expected bool
	}{
		{CaseStatusTypeDraft, true},
		{CaseStatusTypeInProgress, false},
		{CaseStatusTypeCancelled, false},
		{CaseStatusTypeUnknown, false},
	}

	for _, tt := range tests {
		got := tt.status.IsDraft()
		if got != tt.expected {
			t.Errorf("IsDraft(%v) = %v, want %v", tt.status, got, tt.expected)
		}
	}
}

func TestString(t *testing.T) {
	if CaseStatusTypeDraft.ReadableString() != "Draft" {
		t.Errorf("CaseStatusTypeDraft.ReadableString()() = %q, want %q", CaseStatusTypeDraft.ReadableString(), "Draft")
	}
	if CaseStatusTypeRegistered.ReadableString() != "Registered" {
		t.Errorf("CaseStatusTypeRegistered.Key() = %q, want %q", CaseStatusTypeRegistered.ReadableString(), "Registered")
	}
}
