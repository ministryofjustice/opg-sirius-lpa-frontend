package shared

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAttorneyStatus(t *testing.T) {
	active := ActiveAttorneyStatus
	inactive := InactiveAttorneyStatus
	removed := RemovedAttorneyStatus

	assert.Equal(t, active.String(), "active")
	assert.Equal(t, inactive.String(), "inactive")
	assert.Equal(t, removed.String(), "removed")
}

func TestAppointmentType(t *testing.T) {
	original := OriginalAppointmentType
	replacement := ReplacementAppointmentType

	assert.Equal(t, original.String(), "original")
	assert.Equal(t, replacement.String(), "replacement")
}

func TestParseHowAttorneysMakeDecisions(t *testing.T) {
	tests := []struct {
		input    string
		expected HowAttorneysMakeDecisions
	}{
		{"jointly", HowAttorneysMakeDecisionsJointly},
		{"jointly-and-severally", HowAttorneysMakeDecisionsJointlyAndSeverally},
		{"jointly-for-some-severally-for-others", HowAttorneysMakeDecisionsJointlyForSomeSeverallyForOthers},
		{"", HowAttorneysMakeDecisionsEmpty},
		{"notRecognised", HowAttorneysMakeDecisionsNotRecognised},
	}

	for _, tc := range tests {
		got := ParseHowAttorneysMakeDecisions(tc.input)
		assert.Equal(t, tc.expected, got)
	}
}

func TestParseHowAttorneysMakeDecisionsUnknown(t *testing.T) {
	got := ParseHowAttorneysMakeDecisions("invalid")
	assert.Equal(t, HowAttorneysMakeDecisionsNotRecognised, got)
}

func TestHowAttorneysMakeDecisionsTranslation(t *testing.T) {
	tests := []struct {
		name           string
		input          HowAttorneysMakeDecisions
		isSoleAttorney bool
		want           string
	}{
		{
			name:           "Jointly",
			input:          HowAttorneysMakeDecisionsJointly,
			isSoleAttorney: false,
			want:           "Jointly",
		},
		{
			name:           "Jointly and severally",
			input:          HowAttorneysMakeDecisionsJointlyAndSeverally,
			isSoleAttorney: false,
			want:           "Jointly & severally",
		},
		{
			name:           "Jointly for some severally for others",
			input:          HowAttorneysMakeDecisionsJointlyForSomeSeverallyForOthers,
			isSoleAttorney: false,
			want:           "Jointly for some, severally for others",
		},
		{
			name:           "Empty",
			input:          HowAttorneysMakeDecisionsEmpty,
			isSoleAttorney: false,
			want:           "Not specified",
		},
		{
			name:           "Not recognised",
			input:          HowAttorneysMakeDecisionsNotRecognised,
			isSoleAttorney: false,
			want:           "howAttorneysMakeDecisions NOT RECOGNISED: notRecognised",
		},
		{
			name:           "Sole attorney overrides value",
			input:          HowAttorneysMakeDecisionsJointlyForSomeSeverallyForOthers,
			isSoleAttorney: true,
			want:           "There is only one attorney appointed",
		},
		{
			name:           "Sole attorney with empty value",
			input:          HowAttorneysMakeDecisionsEmpty,
			isSoleAttorney: true,
			want:           "There is only one attorney appointed",
		},
		{
			name:           "Unknown int value",
			input:          HowAttorneysMakeDecisions(99),
			isSoleAttorney: false,
			want:           "howAttorneysMakeDecisions NOT RECOGNISED: ",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.input.Translation(tc.isSoleAttorney)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestHowAttorneysMakeDecisionsKey(t *testing.T) {
	tests := []struct {
		name  string
		input HowAttorneysMakeDecisions
		want  string
	}{
		{"Jointly", HowAttorneysMakeDecisionsJointly, "jointly"},
		{"Jointly and severally", HowAttorneysMakeDecisionsJointlyAndSeverally, "jointly-and-severally"},
		{"Jointly for some severally for others", HowAttorneysMakeDecisionsJointlyForSomeSeverallyForOthers, "jointly-for-some-severally-for-others"},
		{"Empty", HowAttorneysMakeDecisionsEmpty, ""},
		{"Not recognised", HowAttorneysMakeDecisionsNotRecognised, "notRecognised"},
		{"Unknown int value", HowAttorneysMakeDecisions(99), ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.input.Key()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestHowAttorneysMakeDecisionsString(t *testing.T) {
	assert.Equal(t, "jointly", HowAttorneysMakeDecisionsJointly.String())
	assert.Equal(t, "", HowAttorneysMakeDecisionsEmpty.String())
}

func TestHowAttorneysMakeDecisionsMarshalJSON(t *testing.T) {
	tests := []struct {
		name  string
		input HowAttorneysMakeDecisions
		want  string
	}{
		{"Jointly", HowAttorneysMakeDecisionsJointly, `"jointly"`},
		{"Jointly and severally", HowAttorneysMakeDecisionsJointlyAndSeverally, `"jointly-and-severally"`},
		{"Jointly for some severally for others", HowAttorneysMakeDecisionsJointlyForSomeSeverallyForOthers, `"jointly-for-some-severally-for-others"`},
		{"Empty", HowAttorneysMakeDecisionsEmpty, `""`},
		{"Not recognised", HowAttorneysMakeDecisionsNotRecognised, `"notRecognised"`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b, err := json.Marshal(tc.input)
			assert.Nil(t, err)
			assert.Equal(t, tc.want, string(b))
		})
	}
}

func TestHowAttorneysMakeDecisionsUnmarshalJSON(t *testing.T) {
	tests := []struct {
		jsonInput string
		expected  HowAttorneysMakeDecisions
	}{
		{`"jointly"`, HowAttorneysMakeDecisionsJointly},
		{`"jointly-and-severally"`, HowAttorneysMakeDecisionsJointlyAndSeverally},
		{`"jointly-for-some-severally-for-others"`, HowAttorneysMakeDecisionsJointlyForSomeSeverallyForOthers},
		{`""`, HowAttorneysMakeDecisionsEmpty},
		{`"notRecognised"`, HowAttorneysMakeDecisionsNotRecognised},
		{`"unknown-value"`, HowAttorneysMakeDecisionsNotRecognised},
	}

	for _, tc := range tests {
		var h HowAttorneysMakeDecisions
		err := json.Unmarshal([]byte(tc.jsonInput), &h)
		assert.Nil(t, err)
		assert.Equal(t, tc.expected, h)
	}
}

func TestHowAttorneysMakeDecisionsUnmarshalJSONErrors(t *testing.T) {
	var h HowAttorneysMakeDecisions
	err := json.Unmarshal([]byte(`123`), &h)
	assert.Error(t, err)
}

func TestParseHowReplacementAttorneysStepIn(t *testing.T) {
	tests := []struct {
		input    string
		expected HowReplacementAttorneysStepIn
	}{
		{"all-can-no-longer-act", HowReplacementAttorneysStepInAllCanNoLongerAct},
		{"one-can-no-longer-act", HowReplacementAttorneysStepInOneCanNoLongerAct},
		{"another-way", HowReplacementAttorneysStepInAnotherWay},
		{"", HowReplacementAttorneysStepInEmpty},
		{"notRecognised", HowReplacementAttorneysStepInNotRecognised},
	}

	for _, tc := range tests {
		got := ParseHowReplacementAttorneysStepIn(tc.input)
		assert.Equal(t, tc.expected, got)
	}
}

func TestParseHowReplacementAttorneysStepInUnknown(t *testing.T) {
	got := ParseHowReplacementAttorneysStepIn("invalid")
	assert.Equal(t, HowReplacementAttorneysStepInNotRecognised, got)
}

func TestHowReplacementAttorneysStepInTranslation(t *testing.T) {
	tests := []struct {
		name  string
		input HowReplacementAttorneysStepIn
		want  string
	}{
		{
			name:  "All can no longer act",
			input: HowReplacementAttorneysStepInAllCanNoLongerAct,
			want:  "When all can no longer act",
		},
		{
			name:  "One can no longer act",
			input: HowReplacementAttorneysStepInOneCanNoLongerAct,
			want:  "When one can no longer act",
		},
		{
			name:  "Another way",
			input: HowReplacementAttorneysStepInAnotherWay,
			want:  "Another way",
		},
		{
			name:  "Empty",
			input: HowReplacementAttorneysStepInEmpty,
			want:  "Not specified",
		},
		{
			name:  "Not recognised",
			input: HowReplacementAttorneysStepInNotRecognised,
			want:  "howReplacementAttorneysStepIn NOT RECOGNISED: notRecognised",
		},
		{
			name:  "Unknown int value",
			input: HowReplacementAttorneysStepIn(99),
			want:  "howReplacementAttorneysStepIn NOT RECOGNISED: ",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.input.Translation()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestHowReplacementAttorneysStepInKey(t *testing.T) {
	tests := []struct {
		name  string
		input HowReplacementAttorneysStepIn
		want  string
	}{
		{"All can no longer act", HowReplacementAttorneysStepInAllCanNoLongerAct, "all-can-no-longer-act"},
		{"One can no longer act", HowReplacementAttorneysStepInOneCanNoLongerAct, "one-can-no-longer-act"},
		{"Another way", HowReplacementAttorneysStepInAnotherWay, "another-way"},
		{"Empty", HowReplacementAttorneysStepInEmpty, ""},
		{"Not recognised", HowReplacementAttorneysStepInNotRecognised, "notRecognised"},
		{"Unknown int value", HowReplacementAttorneysStepIn(99), ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.input.Key()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestHowReplacementAttorneysStepInString(t *testing.T) {
	assert.Equal(t, "all-can-no-longer-act", HowReplacementAttorneysStepInAllCanNoLongerAct.String())
	assert.Equal(t, "", HowReplacementAttorneysStepInEmpty.String())
}

func TestHowReplacementAttorneysStepInMarshalJSON(t *testing.T) {
	tests := []struct {
		name  string
		input HowReplacementAttorneysStepIn
		want  string
	}{
		{"All can no longer act", HowReplacementAttorneysStepInAllCanNoLongerAct, `"all-can-no-longer-act"`},
		{"One can no longer act", HowReplacementAttorneysStepInOneCanNoLongerAct, `"one-can-no-longer-act"`},
		{"Another way", HowReplacementAttorneysStepInAnotherWay, `"another-way"`},
		{"Empty", HowReplacementAttorneysStepInEmpty, `""`},
		{"Not recognised", HowReplacementAttorneysStepInNotRecognised, `"notRecognised"`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b, err := json.Marshal(tc.input)
			assert.Nil(t, err)
			assert.Equal(t, tc.want, string(b))
		})
	}
}

func TestHowReplacementAttorneysStepInUnmarshalJSON(t *testing.T) {
	tests := []struct {
		jsonInput string
		expected  HowReplacementAttorneysStepIn
	}{
		{`"all-can-no-longer-act"`, HowReplacementAttorneysStepInAllCanNoLongerAct},
		{`"one-can-no-longer-act"`, HowReplacementAttorneysStepInOneCanNoLongerAct},
		{`"another-way"`, HowReplacementAttorneysStepInAnotherWay},
		{`""`, HowReplacementAttorneysStepInEmpty},
		{`"notRecognised"`, HowReplacementAttorneysStepInNotRecognised},
		{`"unknown-value"`, HowReplacementAttorneysStepInNotRecognised},
	}

	for _, tc := range tests {
		var h HowReplacementAttorneysStepIn
		err := json.Unmarshal([]byte(tc.jsonInput), &h)
		assert.Nil(t, err)
		assert.Equal(t, tc.expected, h)
	}
}

func TestHowReplacementAttorneysStepInUnmarshalJSONErrors(t *testing.T) {
	var h HowReplacementAttorneysStepIn
	err := json.Unmarshal([]byte(`123`), &h)
	assert.Error(t, err)
}
