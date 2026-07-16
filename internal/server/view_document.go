package server

import (
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type ViewDocumentClient interface {
	DocumentByUUID(ctx sirius.Context, uuid string) (sirius.Document, error)
	GetUserDetails(sirius.Context) (sirius.User, error)
	Case(ctx sirius.Context, id int) (sirius.Case, error)
	PageVarsClient
	GetDraftCount(ctx sirius.Context, caseType string, caseId int) (sirius.DocumentDraftCount, error)
	PersonReferences(ctx sirius.Context, id int) ([]sirius.PersonReference, error)
	TasksForCase(ctx sirius.Context, caseId int) ([]sirius.Task, error)
}

type viewDocumentData struct {
	XSRFToken                      string
	Document                       sirius.Document
	IsSysAdminUser                 bool
	Pane                           int
	DonorID                        int
	SelectedCaseIds                string
	Person                         sirius.Person
	CaseUids                       string
	HasV1PersonsGetPermission      bool
	HasV1PersonsCasesGetPermission bool
	SelectedCases                  []sirius.Case
	ActionPanelButtons             []ActionPanelButton
	HeaderButtons                  SiriusHeaderButtons
}

func ViewDocument(client ViewDocumentClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		uuid := r.PathValue("uuid")
		ctx := getContext(r)

		pageVars, err := PageValues(client, r)
		if err != nil {
			return err
		}

		caseId, err := strToIntOrStatusError(r.FormValue("case"))
		if err != nil {
			return err
		}

		documentData, err := client.DocumentByUUID(ctx, uuid)
		if err != nil {
			return err
		}

		user, err := client.GetUserDetails(ctx)
		if err != nil {
			return err
		}
		isSysAdminUser := user.HasRole("System Admin")

		// Extract pane parameter from query string if present
		pane := 1 // Default to pane 1
		if paneStr := r.URL.Query().Get("pane"); paneStr != "" {
			if paneNum, err := strconv.Atoi(paneStr); err == nil && paneNum == 2 {
				pane = 2
			}
		}

		caseUidsStr := ""
		var selectedCase []sirius.Case
		if len(pageVars.SelectedCases) > 0 {
			for _, c := range pageVars.SelectedCases {
				if c.ID == caseId {
					caseUidsStr = "&uid[]=" + c.UID
					selectedCase = []sirius.Case{c}
				}
			}
		}

		data := viewDocumentData{
			XSRFToken:       ctx.XSRFToken,
			Document:        documentData,
			IsSysAdminUser:  isSysAdminUser,
			Pane:            pane,
			DonorID:         pageVars.DonorID,
			SelectedCaseIds: strconv.Itoa(caseId),
			Person:          pageVars.Person,
			CaseUids:        caseUidsStr,
			SelectedCases:   selectedCase,
		}

		data.ActionPanelButtons = GetActionPanelButtons(data.SelectedCases, data.DonorID, uidParams, draftCount > 0, personHasReferences, len(person.Children) > 0, taskIDs, pageVars.UserPermissions)

		data.HeaderButtons = SiriusHeaderButtons{
			BackToTimeline: true,
			CaseInfo:       true,
			PersonInfo:     true,
			Calendar:       true,
		}

		data.HasV1PersonsGetPermission = pageVars.HasV1PersonsGetPermission
		data.HasV1PersonsCasesGetPermission = pageVars.HasV1PersonsCasesGetPermission

		return tmpl(w, data)
	}
}
