package shared

import (
	"encoding/json"
)

type AttorneyStatus string

const (
	ActiveAttorneyStatus   AttorneyStatus = "active"
	InactiveAttorneyStatus AttorneyStatus = "inactive"
	RemovedAttorneyStatus  AttorneyStatus = "removed"
)

func (attorneyStatus AttorneyStatus) String() string {
	return string(attorneyStatus)
}

type AppointmentType string

const (
	OriginalAppointmentType    AppointmentType = "original"
	ReplacementAppointmentType AppointmentType = "replacement"
)

func (appointmentType AppointmentType) String() string {
	return string(appointmentType)
}

type HowAttorneysMakeDecisions int

const (
	HowAttorneysMakeDecisionsEmpty HowAttorneysMakeDecisions = iota
	HowAttorneysMakeDecisionsJointly
	HowAttorneysMakeDecisionsJointlyAndSeverally
	HowAttorneysMakeDecisionsJointlyForSomeSeverallyForOthers
	HowAttorneysMakeDecisionsNotRecognised
)

var howAttorneysMakeDecisionsMap = map[string]HowAttorneysMakeDecisions{
	"jointly":                              HowAttorneysMakeDecisionsJointly,
	"jointly-and-severally":                HowAttorneysMakeDecisionsJointlyAndSeverally,
	"jointly-for-some-severally-for-others": HowAttorneysMakeDecisionsJointlyForSomeSeverallyForOthers,
	"":                                     HowAttorneysMakeDecisionsEmpty,
	"notRecognised":                        HowAttorneysMakeDecisionsNotRecognised,
}

func (h HowAttorneysMakeDecisions) String() string {
	return h.Key()
}

func (h HowAttorneysMakeDecisions) Translation(isSoleAttorney bool) string {
	if isSoleAttorney {
		return "There is only one attorney appointed"
	}

	switch h {
	case HowAttorneysMakeDecisionsJointly:
		return "Jointly"
	case HowAttorneysMakeDecisionsJointlyAndSeverally:
		return "Jointly & severally"
	case HowAttorneysMakeDecisionsJointlyForSomeSeverallyForOthers:
		return "Jointly for some, severally for others"
	case HowAttorneysMakeDecisionsEmpty:
		return "Not specified"
	default:
		return "howAttorneysMakeDecisions NOT RECOGNISED: " + h.String()
	}
}

func (h HowAttorneysMakeDecisions) Key() string {
	switch h {
	case HowAttorneysMakeDecisionsJointly:
		return "jointly"
	case HowAttorneysMakeDecisionsJointlyAndSeverally:
		return "jointly-and-severally"
	case HowAttorneysMakeDecisionsJointlyForSomeSeverallyForOthers:
		return "jointly-for-some-severally-for-others"
	case HowAttorneysMakeDecisionsEmpty:
		return ""
	case HowAttorneysMakeDecisionsNotRecognised:
		return "notRecognised"
	default:
		return ""
	}
}

func ParseHowAttorneysMakeDecisions(s string) HowAttorneysMakeDecisions {
	value, ok := howAttorneysMakeDecisionsMap[s]
	if !ok {
		return HowAttorneysMakeDecisionsNotRecognised
	}
	return value
}

func (h HowAttorneysMakeDecisions) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.Key())
}

func (h *HowAttorneysMakeDecisions) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*h = ParseHowAttorneysMakeDecisions(s)
	return nil
}

type HowReplacementAttorneysStepIn int

const (
	HowReplacementAttorneysStepInEmpty HowReplacementAttorneysStepIn = iota
	HowReplacementAttorneysStepInAllCanNoLongerAct
	HowReplacementAttorneysStepInOneCanNoLongerAct
	HowReplacementAttorneysStepInAnotherWay
	HowReplacementAttorneysStepInNotRecognised
)

var howReplacementAttorneysStepInMap = map[string]HowReplacementAttorneysStepIn{
	"all-can-no-longer-act": HowReplacementAttorneysStepInAllCanNoLongerAct,
	"one-can-no-longer-act": HowReplacementAttorneysStepInOneCanNoLongerAct,
	"another-way":           HowReplacementAttorneysStepInAnotherWay,
	"":                      HowReplacementAttorneysStepInEmpty,
	"notRecognised":         HowReplacementAttorneysStepInNotRecognised,
}

func (h HowReplacementAttorneysStepIn) String() string {
	return h.Key()
}

func (h HowReplacementAttorneysStepIn) Translation() string {
	switch h {
	case HowReplacementAttorneysStepInAllCanNoLongerAct:
		return "When all can no longer act"
	case HowReplacementAttorneysStepInOneCanNoLongerAct:
		return "When one can no longer act"
	case HowReplacementAttorneysStepInAnotherWay:
		return "Another way"
	case HowReplacementAttorneysStepInEmpty:
		return "Not specified"
	default:
		return "howReplacementAttorneysStepIn NOT RECOGNISED: " + h.String()
	}
}

func (h HowReplacementAttorneysStepIn) Key() string {
	switch h {
	case HowReplacementAttorneysStepInAllCanNoLongerAct:
		return "all-can-no-longer-act"
	case HowReplacementAttorneysStepInOneCanNoLongerAct:
		return "one-can-no-longer-act"
	case HowReplacementAttorneysStepInAnotherWay:
		return "another-way"
	case HowReplacementAttorneysStepInEmpty:
		return ""
	case HowReplacementAttorneysStepInNotRecognised:
		return "notRecognised"
	default:
		return ""
	}
}

func ParseHowReplacementAttorneysStepIn(s string) HowReplacementAttorneysStepIn {
	value, ok := howReplacementAttorneysStepInMap[s]
	if !ok {
		return HowReplacementAttorneysStepInNotRecognised
	}
	return value
}

func (h HowReplacementAttorneysStepIn) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.Key())
}

func (h *HowReplacementAttorneysStepIn) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*h = ParseHowReplacementAttorneysStepIn(s)
	return nil
}

