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
	PageVarsClient
}

type viewDocumentData struct {
	ActionPanelButtons             []ActionPanelButton
	CaseUids                       string
	Document                       sirius.Document
	DonorID                        int
	HasV1PersonsCasesGetPermission bool
	HasV1PersonsGetPermission      bool
	HeaderButtons                  SiriusHeaderButtons
	IsSysAdminUser                 bool
	Pane                           int
	Person                         sirius.Person
	SelectedCaseIds                string
	SelectedCases                  []sirius.Case
	XSRFToken                      string
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
			CaseUids:                       caseUidsStr,
			Document:                       documentData,
			DonorID:                        pageVars.DonorID,
			HasV1PersonsCasesGetPermission: pageVars.HasV1PersonsCasesGetPermission,
			HasV1PersonsGetPermission:      pageVars.HasV1PersonsGetPermission,
			IsSysAdminUser:                 isSysAdminUser,
			Pane:                           pane,
			Person:                         pageVars.Person,
			SelectedCaseIds:                strconv.Itoa(caseId),
			SelectedCases:                  selectedCase,
			XSRFToken:                      ctx.XSRFToken,
		}

		data.ActionPanelButtons = GetActionPanelButtons(data.SelectedCases, data.DonorID, caseUidsStr, pageVars.DraftCount > 0, pageVars.PersonReferences, len(pageVars.Person.Children) > 0, pageVars.TaskIDs, pageVars.UserPermissions)

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
