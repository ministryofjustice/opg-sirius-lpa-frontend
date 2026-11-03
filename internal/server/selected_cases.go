package server

import "fmt"

var selectedCaseUids []string

func GetSelectedCaseUids() []string {
	fmt.Println("Getting case uids: ", selectedCaseUids)
	return selectedCaseUids
}

func SetSelectedCaseUids(caseUids []string) {
	selectedCaseUids = caseUids
	fmt.Println("Setting case uids: ", selectedCaseUids)
}

func AddSelectedCaseUid(caseUid string) {
	selectedCaseUids = append(selectedCaseUids, caseUid)
}

func ClearSelectedCaseUids() {
	selectedCaseUids = []string{}
}
