package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type PageVars struct {
	DonorID                        int
	SelectedCaseIds                string
	Person                         sirius.Person
	CaseUidsCollection             []string
	HasV1PersonsGetPermission      bool
	HasV1PersonsCasesGetPermission bool
	ActionPanelButtons             []ActionPanelButton
	SelectedCases                  []sirius.Case
	HeaderButtons                  SiriusHeaderButtons
	PersonReferences               bool
	CasesOnDonor                   []sirius.Case
	UserPermissions                sirius.Permissions
	DraftCount                     int
	CaseIDs                        []string
	DocumentList                   sirius.DocumentList
}

type PageVarsClient interface {
	CasesByDonor(ctx sirius.Context, id int) ([]sirius.Case, error)
	GetUserPermissions(ctx sirius.Context) (sirius.Permissions, error)
	Person(ctx sirius.Context, id int) (sirius.Person, error)
	PersonReferences(ctx sirius.Context, id int) ([]sirius.PersonReference, error)
	GetDraftCount(ctx sirius.Context, caseType string, caseId int) (sirius.DocumentDraftCount, error)
	GetPersonDocuments(ctx sirius.Context, personID int, caseIDs []string) (sirius.DocumentList, error)
}

func PageValues(client PageVarsClient, r *http.Request) (PageVars, error) {
	ctx := getContext(r)

	donorID, err := strToIntOrStatusError(r.PathValue("id"))
	if err != nil {
		return PageVars{}, err
	}

	caseUIDs := r.Form["uid[]"]

	if len(caseUIDs) == 0 {
		caseUID := r.PathValue("caseUid")
		if caseUID != "" {
			caseUIDs = []string{caseUID}
		}
	}

	userPermissions, err := client.GetUserPermissions(ctx)
	if err != nil {
		return PageVars{}, err
	}

	casesOnDonor, err := client.CasesByDonor(ctx, donorID)
	if err != nil {
		return PageVars{}, err
	}

	person, err := client.Person(ctx, donorID)
	if err != nil {
		return PageVars{}, err
	}

	personReferences, err := client.PersonReferences(ctx, donorID)
	if err != nil {
		return PageVars{}, err
	}
	personHasReferences := len(personReferences) > 0

	var selected []sirius.Case
	var caseIDs []string

	if len(caseUIDs) > 0 {
		casesByUID := make(map[string]sirius.Case, len(casesOnDonor))
		for _, c := range casesOnDonor {
			casesByUID[c.UID] = c
		}

		for _, uid := range caseUIDs {
			if c, ok := casesByUID[uid]; ok {
				selected = append(selected, c)
				caseIDs = append(caseIDs, strconv.Itoa(c.ID))
			}
		}
	} else {
		selected = casesOnDonor
	}

	var draftCount int
	if len(selected) == 1 {
		documentDraftCount, err := client.GetDraftCount(ctx, strings.ToLower(selected[0].CaseType), selected[0].ID)
		if err != nil {
			return PageVars{}, err
		}
		draftCount = documentDraftCount.DraftCount
	}

	docs, err := client.GetPersonDocuments(ctx, donorID, caseIDs)
	if err != nil {
		return PageVars{}, err
	}

	vars := PageVars{
		DonorID:            donorID,
		Person:             person,
		CaseUidsCollection: caseUIDs,
		CaseIDs:            caseIDs,
		CasesOnDonor:       casesOnDonor,
		PersonReferences:   personHasReferences,
		UserPermissions:    userPermissions,
		DraftCount:         draftCount,
		SelectedCases:      selected,
		DocumentList:       docs,
	}

	vars.HasV1PersonsGetPermission = userPermissions.Includes("v1-persons", "GET")
	vars.HasV1PersonsCasesGetPermission = userPermissions.Includes("v1-persons-cases", "GET")

	return vars, nil
}
