package sirius

type Investigation struct {
	Title        string     `json:"investigationTitle"`
	Information  string     `json:"additionalInformation"`
	Type         string     `json:"type"`
	DateReceived DateString `json:"investigationReceivedDate"`
}
