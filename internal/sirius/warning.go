package sirius

type Warning struct {
	ID          int
	Name        string
	WarningType string
	WarningText string
	DateAdded   DateString
	CaseItems   []Case
}
