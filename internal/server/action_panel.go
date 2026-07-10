package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/sync/errgroup"
)

type ActionPanelClient interface {
	CasesByDonor(ctx sirius.Context, id int) ([]sirius.Case, error)
	GetDraftCount(ctx sirius.Context, caseType string, caseId int) (sirius.DocumentDraftCount, error)
	Person(ctx sirius.Context, id int) (sirius.Person, error)
	PersonReferences(ctx sirius.Context, id int) ([]sirius.PersonReference, error)
	GetUserPermissions(ctx sirius.Context) (sirius.Permissions, error)
}

type ActionPanelData struct {
	XSRFToken          string
	ActionPanelButtons []ActionPanelButton
}

func ActionPanel(client ActionPanelClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		err := r.ParseForm()
		if err != nil {
			return err
		}

		ctx := getContext(r)
		data := ActionPanelData{XSRFToken: ctx.XSRFToken}

		// Support both query-string usage (?donorId=123) and path-based params
		donorIDStr := r.URL.Query().Get("donorId")
		if donorIDStr == "" {
			donorIDStr = r.URL.Query().Get("id")
		}

		var donorId int
		if donorIDStr != "" {
			if donorId, err = strconv.Atoi(donorIDStr); err != nil {
				return err
			}
		}

		caseUIDs := r.Form["uid[]"]
		caseUidsString := buildUIDQueryString(caseUIDs)

		group, groupCtx := errgroup.WithContext(ctx.Context)

		var draftCount int
		var personHasReferences bool
		var selectedCases []sirius.Case
		var userPermissions sirius.Permissions
		group.Go(func() error {
			if donorId > 0 {
				cases, err := client.CasesByDonor(ctx.With(groupCtx), donorId)
				if err != nil {
					return err
				}

				// Filter cases by uid[] parameter if provided
				if len(caseUIDs) > 0 {
					casesByUID := make(map[string]sirius.Case, len(cases))
					for _, c := range cases {
						casesByUID[c.UID] = c
					}

					var filtered []sirius.Case
					for _, uid := range caseUIDs {
						if c, ok := casesByUID[uid]; ok {
							filtered = append(filtered, c)
						}
					}
					selectedCases = filtered
				} else {
					selectedCases = cases
				}
			}

			if len(selectedCases) == 1 {
				documentDraftCount, err := client.GetDraftCount(ctx.With(groupCtx), strings.ToLower(selectedCases[0].CaseType), selectedCases[0].ID)
				if err != nil {
					return err
				}
				draftCount = documentDraftCount.DraftCount
			}
			return nil
		})

		group.Go(func() error {
			userPermissions, err = client.GetUserPermissions(ctx)
			if err != nil {
				return err
			}

			return nil
		})

		group.Go(func() error {
			if donorId > 0 {
				personReferences, err := client.PersonReferences(ctx.With(groupCtx), donorId)
				if err != nil {
					return err
				}
				personHasReferences = len(personReferences) > 0
			}
			return nil
		})

		var personHasLinks bool
		group.Go(func() error {
			if donorId > 0 {
				person, err := client.Person(ctx.With(groupCtx), donorId)
				if err != nil {
					return err
				}
				personHasLinks = len(person.Children) > 0
			}
			return nil
		})

		if err := group.Wait(); err != nil {
			return err
		}

		data.ActionPanelButtons = GetActionPanelButtons(selectedCases, donorId, caseUidsString, draftCount > 0, personHasReferences, personHasLinks, userPermissions)

		return tmpl(w, data)
	}
}

type ActionPanelButton struct {
	Label    string
	URL      string
	IconName string
	Disabled bool
	Hidden   bool
}

func GetActionPanelButtons(selectedCases []sirius.Case, donorId int, caseUids string, hasDrafts bool, hasReferences bool, hasLinks bool, userPermissions sirius.Permissions) []ActionPanelButton {
	warningUrl := fmt.Sprintf("/create-warning?id=%d&entity=person%s", donorId, caseUids)
	eventUrl := fmt.Sprintf("/create-event?id=%d&entity=person%s", donorId, caseUids)
	createDonorUrl := fmt.Sprintf("/create-donor?id=%d&entity=person%s", donorId, caseUids)
	editDonorUrl := fmt.Sprintf("/edit-donor?id=%d&entity=person%s", donorId, caseUids)
	miReportingUrl := fmt.Sprintf("/mi-reporting?donorId=%d%s", donorId, caseUids)
	linkPersonUrl := fmt.Sprintf("/link-person?id=%d%s", donorId, caseUids)
	unlinkPersonUrl := fmt.Sprintf("/unlink-person?id=%d%s", donorId, caseUids)
	deleteRelationshipUrl := fmt.Sprintf("/delete-relationship?id=%d%s", donorId, caseUids)
	createRelationshipUrl := fmt.Sprintf("/create-relationship?id=%d&entity=person%s", donorId, caseUids)
	createEpaUrl := fmt.Sprintf("/create-epa?id=%d", donorId)
	editEpaUrl := ""
	complaintUrl := ""
	createDocumentUrl := ""
	createInvestigationUrl := ""
	editDocumentUrl := ""
	changeStatusUrl := ""
	paymentsUrl := ""
	newTaskUrl := ""
	editDatesUrl := ""
	allocateCasesUrl := ""

	if len(selectedCases) == 1 {
		selectedCase := selectedCases[0]
		caseType := strings.ToLower(selectedCase.CaseType)
		id := selectedCase.ID

		warningUrl = fmt.Sprintf("/create-warning?id=%d&entity=%s%s", donorId, caseType, caseUids)
		complaintUrl = fmt.Sprintf("/add-complaint?id=%d&case=%s", id, caseType)
		createDocumentUrl = fmt.Sprintf("/create-document?id=%d&case=%s", id, caseType)
		changeStatusUrl = fmt.Sprintf("/change-status?id=%d&case=%s&donorId=%d%s", id, caseType, donorId, caseUids)
		paymentsUrl = fmt.Sprintf("/payments/%d", id)
		newTaskUrl = fmt.Sprintf("/create-task?id=%d&entity=%s%s", id, caseType, caseUids)
		editDatesUrl = fmt.Sprintf("/edit-dates?id=%d&case=%s", id, caseType)
		allocateCasesUrl = fmt.Sprintf("/allocate-cases?id=%d&entity=%s%s", id, caseType, caseUids)
		createInvestigationUrl = fmt.Sprintf("/create-investigation?id=%d&case=%s%s", id, caseType, caseUids)

		if strings.ToLower(selectedCase.CaseType) == "epa" {
			editEpaUrl = fmt.Sprintf("/create-epa?id=%d&caseId=%d", donorId, selectedCases[0].ID)
		}

		if hasDrafts {
			editDocumentUrl = fmt.Sprintf("/edit-document?id=%d&case=%s", id, caseType)
		}
	}
	if len(selectedCases) > 1 {
		idQuery := ""
		caseType := strings.ToLower(selectedCases[0].CaseType)
		for i, c := range selectedCases {
			if i == 0 {
				idQuery += fmt.Sprintf("id=%d", c.ID)
			} else {
				idQuery += fmt.Sprintf("&id=%d", c.ID)
			}
		}
		allocateCasesUrl = fmt.Sprintf("/allocate-cases?%s&entity=%s%s", idQuery, caseType, caseUids)
	}

	return []ActionPanelButton{
		{
			Label:    "Create warning",
			URL:      warningUrl,
			IconName: "aw-create-warning",
			Disabled: false,
			Hidden:   !userPermissions.Includes("v1-warnings", "POST"),
		},
		{
			Label:    "Create event",
			URL:      eventUrl,
			IconName: "aw-new-event",
			Disabled: false,
			Hidden:   !userPermissions.Includes("v1-notes", "POST"),
		},
		{
			Label:    "Add complaint",
			URL:      complaintUrl,
			IconName: "aw-log-complaint",
			Disabled: len(selectedCases) != 1,
			Hidden:   !userPermissions.Includes("v1-warnings", "POST"),
		},
		{
			Label:    "Create document",
			URL:      createDocumentUrl,
			IconName: "aw-new-template",
			Disabled: len(selectedCases) != 1,
			Hidden:   !userPermissions.Includes("v1-lpas-documents-draft", "POST"),
		},
		{
			Label:    "Retrieve draft",
			URL:      editDocumentUrl,
			IconName: "aw-new-template",
			Disabled: len(selectedCases) != 1 || !hasDrafts,
			Hidden:   !userPermissions.Includes("v1-lpas-documents-draft", "POST"),
		},
		{
			Label:    "Change status",
			URL:      changeStatusUrl,
			IconName: "aw-change-status",
			Disabled: len(selectedCases) != 1,
			Hidden:   !userPermissions.Includes("v1-lpas", "PUT"),
		},
		{
			Label:    "Fees",
			URL:      paymentsUrl,
			IconName: "aw-fees",
			Disabled: len(selectedCases) != 1,
			Hidden:   !userPermissions.Includes("v1-payments", "GET"),
		},
		{
			Label:    "New task",
			URL:      newTaskUrl,
			IconName: "aw-new-task",
			Disabled: len(selectedCases) != 1,
			Hidden:   !userPermissions.Includes("v1-cases-tasks-post", "POST"),
		},
		{
			Label:    "Create donor",
			URL:      createDonorUrl,
			IconName: "aw-create-person",
			Disabled: false,
			Hidden:   !userPermissions.Includes("v1-donors", "POST"),
		},
		{
			Label:    "Edit donor",
			URL:      editDonorUrl,
			IconName: "aw-edit-person",
			Disabled: false,
			Hidden:   !userPermissions.Includes("v1-donors", "PUT"),
		},
		{
			Label:    "Edit dates",
			URL:      editDatesUrl,
			IconName: "calendar-open",
			Disabled: len(selectedCases) != 1,
			Hidden:   !userPermissions.Includes("v1-lpas", "PUT"),
		},
		{
			Label:    "MI reporting",
			URL:      miReportingUrl,
			IconName: "aw-mi",
			Disabled: false,
			Hidden:   !userPermissions.Includes("reporting", "GET"),
		},
		{
			Label:    "Allocate Case",
			URL:      allocateCasesUrl,
			IconName: "aw-allocate-case",
			Disabled: len(selectedCases) == 0,
			Hidden:   !userPermissions.Includes("v1-users-updateusercases", "PUT"),
		},
		{
			Label:    "Link record",
			URL:      linkPersonUrl,
			IconName: "aw-link",
			Disabled: donorId == 0,
			Hidden:   !userPermissions.Includes("v1-person-links", "POST"),
		},
		{
			Label:    "Unlink record",
			URL:      unlinkPersonUrl,
			IconName: "aw-unlink",
			Disabled: donorId == 0 || !hasLinks,
			Hidden:   !userPermissions.Includes("v1-person-links", "PATCH"),
		},
		{
			Label:    "Delete relationship",
			URL:      deleteRelationshipUrl,
			IconName: "icon-minus",
			Disabled: !hasReferences,
			Hidden:   !userPermissions.Includes("v1-person-references", "DELETE"),
		},
		{
			Label:    "Create relationship",
			URL:      createRelationshipUrl,
			IconName: "aw-relationship",
			Disabled: false,
			Hidden:   !userPermissions.Includes("v1-persons-references", "POST"),
		},
		{
			Label:    "Create epa case",
			URL:      createEpaUrl,
			IconName: "aw-create-case",
			Disabled: caseUids != "",
			Hidden:   !userPermissions.Includes("v1-donors-epas", "POST"),
		},
		{
			Label:    "Edit epa case",
			URL:      editEpaUrl,
			IconName: "aw-edit-case",
			Disabled: len(selectedCases) != 1 || strings.ToLower(selectedCases[0].CaseType) != "epa",
			Hidden:   !userPermissions.Includes("v1-lpas", "PUT"),
		},
		{
			Label:    "Add investigation",
			URL:      createInvestigationUrl,
			IconName: "icon-investigation",
			Disabled: len(selectedCases) != 1,
			Hidden:   !userPermissions.Includes("v1-lpas-investigations", "POST"),
		},
	}
}
