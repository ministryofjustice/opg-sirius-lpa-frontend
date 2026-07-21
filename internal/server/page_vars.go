package server

import (
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/sync/errgroup"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type PageVars struct {
	ActionPanelButtons             []ActionPanelButton
	CaseIDs                        []string
	CaseUidsCollection             []string
	CasesOnDonor                   []sirius.Case
	DocumentList                   sirius.DocumentList
	DonorID                        int
	DraftCount                     int
	HasV1PersonsCasesGetPermission bool
	HasV1PersonsGetPermission      bool
	HeaderButtons                  SiriusHeaderButtons
	Person                         sirius.Person
	PersonReferences               bool
	SelectedCaseIds                string
	SelectedCases                  []sirius.Case
	UserPermissions                sirius.Permissions
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
	if err != nil || donorID == 0 {
		donorID, err = strToIntOrStatusError(r.URL.Query().Get("id"))
	}

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

	var userPermissions sirius.Permissions
	var casesOnDonor []sirius.Case
	var person sirius.Person
	var personReferences []sirius.PersonReference

	group, _ := errgroup.WithContext(ctx.Context)

	group.Go(func() error {
		var err error
		userPermissions, err = client.GetUserPermissions(ctx)
		return err
	})

	group.Go(func() error {
		var err error
		casesOnDonor, err = client.CasesByDonor(ctx, donorID)
		return err
	})

	group.Go(func() error {
		var err error
		person, err = client.Person(ctx, donorID)
		return err
	})

	group.Go(func() error {
		var err error
		personReferences, err = client.PersonReferences(ctx, donorID)
		return err
	})

	if err := group.Wait(); err != nil {
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
	var docs sirius.DocumentList

	group, _ = errgroup.WithContext(ctx.Context)

	if len(selected) == 1 {
		group.Go(func() error {
			documentDraftCount, err := client.GetDraftCount(ctx, strings.ToLower(selected[0].CaseType), selected[0].ID)
			if err != nil {
				return err
			}
			draftCount = documentDraftCount.DraftCount
			return nil
		})
	}

	group.Go(func() error {
		var err error
		docs, err = client.GetPersonDocuments(ctx, donorID, caseIDs)
		return err
	})

	if err := group.Wait(); err != nil {
		return PageVars{}, err
	}

	vars := PageVars{
		CaseIDs:            caseIDs,
		CaseUidsCollection: caseUIDs,
		CasesOnDonor:       casesOnDonor,
		DocumentList:       docs,
		DonorID:            donorID,
		DraftCount:         draftCount,
		Person:             person,
		PersonReferences:   personHasReferences,
		SelectedCases:      selected,
		UserPermissions:    userPermissions,
	}

	vars.HasV1PersonsGetPermission = userPermissions.Includes("v1-persons", "GET")
	vars.HasV1PersonsCasesGetPermission = userPermissions.Includes("v1-persons-cases", "GET")

	return vars, nil
}
