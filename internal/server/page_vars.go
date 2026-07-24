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
	TaskIDs                        []int
	UserPermissions                sirius.Permissions
}

type PageVarsClient interface {
	CasesByDonor(ctx sirius.Context, id int) ([]sirius.Case, error)
	GetUserPermissions(ctx sirius.Context) (sirius.Permissions, error)
	Person(ctx sirius.Context, id int) (sirius.Person, error)
	PersonReferences(ctx sirius.Context, id int) ([]sirius.PersonReference, error)
	GetDraftCount(ctx sirius.Context, caseType string, caseId int) (sirius.DocumentDraftCount, error)
	GetPersonDocuments(ctx sirius.Context, personID int, caseIDs []string) (sirius.DocumentList, error)
	TasksForCase(ctx sirius.Context, id int) ([]sirius.Task, error)
}

func PageValues(client PageVarsClient, r *http.Request) (PageVars, error) {
	ctx := getContext(r)

	donorID, _ := strToIntOrStatusError(r.PathValue("id"))
	if donorID == 0 {
		donorID, _ = strToIntOrStatusError(r.URL.Query().Get("id"))
	}
	if donorID == 0 {
		donorID, _ = strToIntOrStatusError(r.URL.Query().Get("donorId"))
	}

	caseUIDs := r.Form["uid[]"]

	if len(caseUIDs) == 0 {
		caseUID := r.PathValue("caseUid")
		if caseUID == "" {
			caseUID = r.PathValue("caseId")
		}
		if caseUID != "" {
			caseUIDs = []string{caseUID}
		}
	}

	var userPermissions sirius.Permissions
	var casesOnDonor []sirius.Case
	var person sirius.Person
	var personReferences []sirius.PersonReference
	var personHasReferences bool

	if donorID > 0 {
		group, _ := errgroup.WithContext(ctx.Context)

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

		personHasReferences = len(personReferences) > 0
	}

	var selected []sirius.Case
	var caseIDs []string

	if len(caseUIDs) > 0 {
		casesByUID := make(map[string]sirius.Case, len(casesOnDonor))
		casesByID := make(map[string]sirius.Case, len(casesOnDonor))
		for _, c := range casesOnDonor {
			casesByUID[c.UID] = c
			casesByID[strconv.Itoa(c.ID)] = c
		}

		for _, uid := range caseUIDs {
			// Try to match by UID first
			if c, ok := casesByUID[uid]; ok {
				selected = append(selected, c)
				caseIDs = append(caseIDs, strconv.Itoa(c.ID))
			} else if c, ok := casesByID[uid]; ok {
				// If not found by UID, try by ID (for path parameters like caseId)
				selected = append(selected, c)
				caseIDs = append(caseIDs, uid)
			}
		}
	} else {
		selected = casesOnDonor
	}

	var draftCount int
	var docs sirius.DocumentList
	var taskIDs []int

	group, _ := errgroup.WithContext(ctx.Context)

	if len(selected) == 1 {
		group.Go(func() error {
			documentDraftCount, err := client.GetDraftCount(ctx, strings.ToLower(selected[0].CaseType), selected[0].ID)
			if err != nil {
				return err
			}
			draftCount = documentDraftCount.DraftCount
			return nil
		})

		group.Go(func() error {
			tasks, err := client.TasksForCase(ctx, selected[0].ID)
			if err != nil {
				return err
			}
			for _, task := range tasks {
				taskIDs = append(taskIDs, task.ID)
			}
			return nil
		})
	}

	if donorID != 0 || len(caseIDs) != 0 {
		group.Go(func() error {
			var err error
			docs, err = client.GetPersonDocuments(ctx, donorID, caseIDs)
			return err
		})

	}

	group.Go(func() error {
		var err error
		userPermissions, err = client.GetUserPermissions(ctx)
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
		TaskIDs:            taskIDs,
		UserPermissions:    userPermissions,
	}

	vars.HasV1PersonsGetPermission = userPermissions.Includes("v1-persons", "GET")
	vars.HasV1PersonsCasesGetPermission = userPermissions.Includes("v1-persons-cases", "GET")

	return vars, nil
}
