package shared

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
