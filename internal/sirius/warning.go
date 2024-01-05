package sirius

type Warning struct {
	ID          int        `json:"id"`
	WarningType string     `json:"warningType"`
	WarningText string     `json:"warningText"`
	DateAdded   DateString `json:"dateAdded"`
	CaseItems   []Case     `json:"caseItems"`
}
