package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type ActionPanelClient interface {
	PageVarsClient
	TasksForCase(ctx sirius.Context, caseId int) ([]sirius.Task, error)
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

		pageVars, err := PageValues(client, r)
		if err != nil {
			return err
		}

		caseUidsString := buildUIDQueryString(pageVars.CaseUidsCollection)

		var personHasLinks bool
		if pageVars.DonorID > 0 {
			personHasLinks = len(pageVars.Person.Children) > 0
		}

		data.ActionPanelButtons = GetActionPanelButtons(pageVars.SelectedCases, pageVars.DonorID, caseUidsString, pageVars.DraftCount > 0, pageVars.PersonReferences, personHasLinks, pageVars.UserPermissions)

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

func GetActionPanelButtons(selectedCases []sirius.Case, donorId int, caseUids string, hasDrafts bool, hasReferences bool, hasLinks bool, taskIDs []int, userPermissions sirius.Permissions) []ActionPanelButton {
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
	assignTaskUrl := ""

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

		if len(taskIDs) > 0 {
			idQuery := ""
			for i, taskID := range taskIDs {
				if i == 0 {
					idQuery += fmt.Sprintf("id=%d", taskID)
				} else {
					idQuery += fmt.Sprintf("&id=%d", taskID)
				}
			}
			assignTaskUrl = fmt.Sprintf("/assign-task?%s&donorId=%d%s", idQuery, donorId, caseUids)
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
			Label:    "Assign task",
			URL:      assignTaskUrl,
			IconName: "aw-assign-task",
			Disabled: len(selectedCases) != 1 || len(taskIDs) == 0,
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
