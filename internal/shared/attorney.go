package shared

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
