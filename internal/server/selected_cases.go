package server

var selectedCaseUids []string

func GetSelectedCaseUids() []string {
	return selectedCaseUids
}

func SetSelectedCaseUids(caseUids []string) {
	selectedCaseUids = caseUids
}

func ClearSelectedCaseUids() {
	selectedCaseUids = []string{}
}
