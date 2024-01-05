package sirius

type Warning struct {
	ID          int
	WarningType string
	WarningText string
	DateAdded   DateString
	CaseItems   []Case
}
