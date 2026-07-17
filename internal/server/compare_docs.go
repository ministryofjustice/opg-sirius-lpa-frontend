package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type CompareDocsClient interface {
	DocumentByUUID(ctx sirius.Context, uuid string) (sirius.Document, error)
	GetPersonDocuments(ctx sirius.Context, personID int, caseIDs []string) (sirius.DocumentList, error)
	Person(ctx sirius.Context, id int) (sirius.Person, error)
	GetUserPermissions(ctx sirius.Context) (sirius.Permissions, error)
	Case(ctx sirius.Context, id int) (sirius.Case, error)
	GetDraftCount(ctx sirius.Context, caseType string, caseId int) (sirius.DocumentDraftCount, error)
	PersonReferences(ctx sirius.Context, id int) ([]sirius.PersonReference, error)
	PageVarsClient
	TasksForCase(ctx sirius.Context, caseId int) ([]sirius.Task, error)
}

type compareDocsData struct {
	DocListPane1Data               documentPageData
	DocListPane2Data               documentPageData
	Pane1                          string
	Pane2                          string
	View1                          *viewingDocumentData
	View2                          *viewingDocumentData
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

type viewingDocumentData struct {
	Document sirius.Document
	Pane     int
	BackURL  string
	CloseURL string
}

func CompareDocs(client CompareDocsClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		pageVars, err := PageValues(client, r)
		if err != nil {
			return err
		}

		id, _ := strconv.Atoi(pageVars.CaseIDs[0])
		var selectedCase []sirius.Case
		if caseData, err := client.Case(ctx, id); err == nil {
			selectedCase = []sirius.Case{caseData}
		}
		var draftCount int
		var taskIDs []int
		if len(selectedCase) > 0 {
			documentDraftCount, err := client.GetDraftCount(ctx, strings.ToLower(selectedCase[0].CaseType), selectedCase[0].ID)
			if err != nil {
				return err
			}
			draftCount = documentDraftCount.DraftCount

			tasks, err := client.TasksForCase(ctx, selectedCase[0].ID)
			if err != nil {
				return err
			}
			for _, task := range tasks {
				taskIDs = append(taskIDs, task.ID)
			}
		}

		baseURL := fmt.Sprintf("/compare/%d/%s", pageVars.DonorID, pageVars.CaseIDs[0])

		data := compareDocsData{
			Pane1: "list",
			Pane2: "list",
			DocListPane1Data: documentPageData{
				XSRFToken:     ctx.XSRFToken,
				DocumentList:  pageVars.DocumentList,
				SelectedCases: selectedCase,
				Comparing:     true,
				DonorID:       pageVars.DonorID,
			},
			DocListPane2Data: documentPageData{
				XSRFToken:     ctx.XSRFToken,
				DocumentList:  pageVars.DocumentList,
				SelectedCases: selectedCase,
				Comparing:     true,
				DonorID:       pageVars.DonorID,
			},
			DonorID:         pageVars.DonorID,
			SelectedCaseIds: pageVars.CaseIDs[0],
			Person:          pageVars.Person,
			CaseUids:        "&uid[]=" + selectedCase[0].UID,
			SelectedCases:   selectedCase,
		}

		pane1UUID := r.URL.Query().Get("pane1")
		pane2UUID := r.URL.Query().Get("pane2")

		data.DocListPane1Data.CompareURLs = make(map[string]string)
		data.DocListPane2Data.CompareURLs = make(map[string]string)

		for _, doc := range pageVars.DocumentList.Documents {
			panel1Url := baseURL + "?pane1=" + doc.UUID
			panel2Url := baseURL + "?pane2=" + doc.UUID
			if pane1UUID != "" {
				panel2Url += "&pane1=" + pane1UUID
			}
			if pane2UUID != "" {
				panel1Url += "&pane2=" + pane2UUID
			}
			data.DocListPane1Data.CompareURLs[doc.UUID] = panel1Url
			data.DocListPane2Data.CompareURLs[doc.UUID] = panel2Url
		}

		if pane1UUID != "" {
			doc, err := client.DocumentByUUID(ctx, pane1UUID)
			if err != nil {
				return err
			}

			backURL := baseURL
			if pane2UUID != "" {
				backURL = baseURL + "?pane2=" + pane2UUID
			}

			data.Pane1 = "doc"
			data.View1 = &viewingDocumentData{
				Document: doc,
				Pane:     1,
				BackURL:  backURL,
			}
		}

		if pane2UUID != "" {
			doc, err := client.DocumentByUUID(ctx, pane2UUID)
			if err != nil {
				return err
			}

			backURL := baseURL
			if pane1UUID != "" {
				backURL = baseURL + "?pane1=" + pane1UUID
			}

			data.Pane2 = "doc"
			data.View2 = &viewingDocumentData{
				Document: doc,
				Pane:     2,
				BackURL:  backURL,
			}
		}

		data.ActionPanelButtons = GetActionPanelButtons(pageVars.SelectedCases, data.DonorID, data.CaseUids, draftCount > 0, pageVars.PersonReferences, len(pageVars.Person.Children) > 0, taskIDs, ctx.Permissions)

		data.HeaderButtons = SiriusHeaderButtons{
			BackToTimeline: true,
			CaseInfo:       true,
			PersonInfo:     true,
			Calendar:       true,
		}

		data.HasV1PersonsGetPermission = pageVars.HasV1PersonsGetPermission
		data.HasV1PersonsCasesGetPermission = pageVars.HasV1PersonsCasesGetPermission

		viewingADocumentAndList := data.Pane1 == "doc" && data.Pane2 == "list"
		if viewingADocumentAndList {
			data.DocListPane2Data.CloseURL = fmt.Sprintf("/view-document/%s/%d?case=%d&pane=1", data.View1.Document.UUID, pageVars.DonorID, selectedCase[0].ID)
			data.View1.CloseURL = fmt.Sprintf("/donor/%d/documents?uid[]=%s", pageVars.DonorID, data.DocListPane2Data.DocumentList.Documents[0].CaseItems[0].UID)
		}

		viewingAListAndDocument := data.Pane1 == "list" && data.Pane2 == "doc"
		if viewingAListAndDocument {
			data.DocListPane1Data.CloseURL = fmt.Sprintf("/view-document/%s/%d?case=%d&pane=2", data.View2.Document.UUID, pageVars.DonorID, selectedCase[0].ID)
			data.View2.CloseURL = fmt.Sprintf("/donor/%d/documents?uid[]=%s", pageVars.DonorID, data.DocListPane1Data.DocumentList.Documents[0].CaseItems[0].UID)
		}

		bothSidesAreDocuments := data.Pane1 == "doc" && data.Pane2 == "doc"
		if bothSidesAreDocuments {
			data.View1.CloseURL = fmt.Sprintf("/view-document/%s/%d?case=%d&pane=2", data.View2.Document.UUID, pageVars.DonorID, selectedCase[0].ID)
			data.View2.CloseURL = fmt.Sprintf("/view-document/%s/%d?case=%d&pane=1", data.View1.Document.UUID, pageVars.DonorID, selectedCase[0].ID)
		}

		bothSidesAreLists := data.Pane1 == "list" && data.Pane2 == "list"
		if bothSidesAreLists {
			data.DocListPane1Data.CloseURL = fmt.Sprintf("/donor/%d/documents?uid[]=%s", pageVars.DonorID, data.DocListPane2Data.DocumentList.Documents[0].CaseItems[0].UID)
			data.DocListPane2Data.CloseURL = fmt.Sprintf("/donor/%d/documents?uid[]=%s", pageVars.DonorID, data.DocListPane1Data.DocumentList.Documents[0].CaseItems[0].UID)
		}

		return tmpl(w, data)
	}
}
