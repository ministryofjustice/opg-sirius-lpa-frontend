package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type ViewDocumentClient interface {
	DocumentByUUID(ctx sirius.Context, uuid string) (sirius.Document, error)
	GetUserDetails(sirius.Context) (sirius.User, error)
	Person(ctx sirius.Context, id int) (sirius.Person, error)
	GetUserPermissions(ctx sirius.Context) (sirius.Permissions, error)
	Case(ctx sirius.Context, id int) (sirius.Case, error)
	GetDraftCount(ctx sirius.Context, caseType string, caseId int) (sirius.DocumentDraftCount, error)
	PersonReferences(ctx sirius.Context, id int) ([]sirius.PersonReference, error)
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
	ID                             int
	ActionPanelButtons             []ActionPanelButton
	HeaderButtons                  SiriusHeaderButtons
}

func ViewDocument(client ViewDocumentClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		uuid := r.PathValue("uuid")
		ctx := getContext(r)

		donorID, err := strconv.Atoi(r.PathValue("donorId"))
		if err != nil {
			return err
		}

		person, err := client.Person(ctx, donorID)
		if err != nil {
			return err
		}

		caseId := r.FormValue("case")

		documentData, err := client.DocumentByUUID(ctx, uuid)
		if err != nil {
			return err
		}

		user, err := client.GetUserDetails(ctx)
		if err != nil {
			return err
		}
		isSysAdminUser := user.HasRole("System Admin")

		personReferences, err := client.PersonReferences(ctx, donorID)
		if err != nil {
			return err
		}
		personHasReferences := len(personReferences) > 0

		// Extract pane parameter from query string if present
		pane := 1 // Default to pane 1
		if paneStr := r.URL.Query().Get("pane"); paneStr != "" {
			if paneNum, err := strconv.Atoi(paneStr); err == nil && paneNum == 2 {
				pane = 2
			}
		}

		id, _ := strconv.Atoi(caseId)
		var selectedCase []sirius.Case
		if caseData, err := client.Case(ctx, id); err == nil {
			selectedCase = []sirius.Case{caseData}
		}

		var draftCount int
		if len(selectedCase) > 0 {
			documentDraftCount, err := client.GetDraftCount(ctx, strings.ToLower(selectedCase[0].CaseType), selectedCase[0].ID)
			if err != nil {
				return err
			}
			draftCount = documentDraftCount.DraftCount
		}

		caseUidsStr := ""
		uidParams := ""
		if len(selectedCase) > 0 {
			caseUidsStr = "&uid[]=" + selectedCase[0].UID
			uidParams = caseUidsStr
		}

		data := viewDocumentData{
			XSRFToken:       ctx.XSRFToken,
			Document:        documentData,
			IsSysAdminUser:  isSysAdminUser,
			Pane:            pane,
			DonorID:         donorID,
			SelectedCaseIds: caseId,
			Person:          person,
			CaseUids:        caseUidsStr,
			SelectedCases:   selectedCase,
		}

		userPermissions, err := client.GetUserPermissions(ctx)
		if err != nil {
			return err
		}

		data.ActionPanelButtons = GetActionPanelButtons(data.SelectedCases, data.DonorID, uidParams, draftCount > 0, personHasReferences, len(person.Children) > 0, userPermissions)
		data.HeaderButtons = SiriusHeaderButtons{
			BackToTimeline: true,
			CaseInfo:       true,
			PersonInfo:     true,
			Calendar:       true,
		}

		data.HasV1PersonsGetPermission = userPermissions.Includes("v1-persons", "GET")
		data.HasV1PersonsCasesGetPermission = userPermissions.Includes("v1-persons-cases", "GET")

		return tmpl(w, data)
	}
}
