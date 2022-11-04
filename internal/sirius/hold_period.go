package sirius

type HoldPeriod struct {
	ID        int        `json:"id,omitempty"`
	StartDate DateString `json:"startDate"`
	EndDate   DateString `json:"endDate"`
	Reason    string     `json:"reason"`
}
