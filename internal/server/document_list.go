package server

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type DocumentListClient interface {
	PageVarsClient
	DownloadMultiple(ctx sirius.Context, docIDs []string) (*http.Response, error)
	GetUserPermissions(ctx sirius.Context) (sirius.Permissions, error)
	GetDraftCount(ctx sirius.Context, caseType string, caseId int) (sirius.DocumentDraftCount, error)
	PersonReferences(ctx sirius.Context, id int) ([]sirius.PersonReference, error)
	TasksForCase(ctx sirius.Context, caseId int) ([]sirius.Task, error)
}

type documentPageData struct {
	XSRFToken                      string
	Entity                         string
	Success                        bool
	SuccessMessage                 string
	Error                          sirius.ValidationError
	DocumentList                   sirius.DocumentList
	Document                       sirius.Document
	MultipleCasesSelected          bool
	Comparing                      bool
	CompareURLs                    map[string]string
	CloseURL                       string
	DonorID                        int
	SelectedCaseIds                string
	Person                         sirius.Person
	CaseUids                       string
	HasV1PersonsGetPermission      bool
	HasV1PersonsCasesGetPermission bool
	ActionPanelButtons             []ActionPanelButton
	SelectedCases                  []sirius.Case
	HeaderButtons                  SiriusHeaderButtons
}

func DocumentList(client DocumentListClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)
		if err := r.ParseForm(); err != nil {
			return err
		}

		pageVars, err := PageValues(client, r)
		if err != nil {
			return err
		}

		selectedDocUUIDs := r.Form["document"]

		if r.Method == http.MethodPost && len(selectedDocUUIDs) > 0 && r.FormValue("actionDownload") == "true" {
			downloadResp, err := client.DownloadMultiple(ctx, selectedDocUUIDs)
			if err != nil {
				return err
			}
			defer downloadResp.Body.Close() //nolint:errcheck // no need to check error when closing body

			for key, values := range downloadResp.Header {
				for _, value := range values {
					w.Header().Add(key, value)
				}
			}

			w.WriteHeader(downloadResp.StatusCode)
			if _, err := io.Copy(w, downloadResp.Body); err != nil {
				return err
			}

			return nil
		}

		compareView := r.FormValue("comparing") == "true"
		var validationErr sirius.ValidationError
		if r.Method == http.MethodPost && len(selectedDocUUIDs) == 0 && r.FormValue("actionDownload") == "true" {
			if compareView {
				w.WriteHeader(http.StatusNoContent)
				return nil
			}
			validationErr.Detail = "Select one or more documents and try again."
		}

		successMessage := ""
		isSuccess := r.URL.Query().Get("success") == "true" && r.FormValue("dismissNotification") != "true"

		if isSuccess {
			successMessage = successMessageFormatter(r.URL.Query().Get("documentFriendlyName"), r.URL.Query().Get("documentCreatedTime"), "02/01/2006 15:04:05", "02/01/2006")
		}

		data := documentPageData{
			XSRFToken:             ctx.XSRFToken,
			SelectedCases:         pageVars.SelectedCases,
			Person:                pageVars.Person,
			DocumentList:          pageVars.DocumentList,
			MultipleCasesSelected: len(pageVars.CaseUidsCollection) > 1 || (len(pageVars.CaseUidsCollection) == 0 && len(pageVars.CasesOnDonor) > 1),
			Error:                 validationErr,
			Success:               isSuccess,
			SuccessMessage:        successMessage,
			Comparing:             compareView,
			DonorID:               pageVars.DonorID,
		}

		uidParams := buildUIDQueryString(pageVars.CaseUidsCollection)

		data.CaseUids = uidParams

		for index, selectedCase := range data.SelectedCases {
			if index != 0 {
				data.SelectedCaseIds += "+"
			}
			data.SelectedCaseIds += strconv.Itoa(selectedCase.ID)
		}

		data.ActionPanelButtons = GetActionPanelButtons(data.SelectedCases, data.DonorID, uidParams, draftCount > 0, personHasReferences, len(person.Children) > 0, taskIDs, pageVars.UserPermissions)

		data.HeaderButtons = SiriusHeaderButtons{
			BackToTimeline: true,
			Calendar:       true,
		}

		data.HasV1PersonsGetPermission = pageVars.HasV1PersonsGetPermission
		data.HasV1PersonsCasesGetPermission = pageVars.HasV1PersonsCasesGetPermission

		return tmpl(w, data)
	}
}

func successMessageFormatter(docFriendlyName string, docCreatedTime, layout string, format string) string {
	t, err := time.Parse(layout, docCreatedTime)
	if err != nil {
		return "invalid date"
	}

	return t.Format(format) + " " + docFriendlyName
}

func buildUIDQueryString(uids []string) string {
	var result string
	for _, uid := range uids {
		result += "&uid[]=" + uid
	}
	return result
}
