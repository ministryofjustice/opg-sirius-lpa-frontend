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
	CasesByDonor(ctx sirius.Context, id int) ([]sirius.Case, error)
	GetPersonDocuments(ctx sirius.Context, personID int, caseIDs []string) (sirius.DocumentList, error)
	DownloadMultiple(ctx sirius.Context, docIDs []string) (*http.Response, error)
}

type documentListData struct {
	XSRFToken             string
	Entity                string
	Success               bool
	SuccessMessage        string
	Error                 sirius.ValidationError
	DocumentList          sirius.DocumentList
	SelectedCases         []sirius.Case
	MultipleCasesSelected bool
}

func DocumentList(client DocumentListClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if err := r.ParseForm(); err != nil {
			return err
		}

		donorID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			return err
		}

		caseUIDs := r.Form["uid[]"]

		ctx := getContext(r)

		casesOnDonor, err := client.CasesByDonor(ctx, donorID)
		if err != nil {
			return err
		}

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

		docs, err := client.GetPersonDocuments(ctx, donorID, caseIDs)
		if err != nil {
			return err
		}

		var validationErr sirius.ValidationError
		if r.Method == http.MethodPost && len(selectedDocUUIDs) == 0 && r.FormValue("actionDownload") == "true" {
			validationErr.Detail = "Select one or more documents and try again."
		}

		successMessage := ""
		isSuccess := r.URL.Query().Get("success") == "true" && r.FormValue("dismissNotification") != "true"

		if isSuccess {
			successMessage = successMessageFormatter(r.URL.Query().Get("documentFriendlyName"), r.URL.Query().Get("documentCreatedTime"), "02/01/2006 15:04:05", "02/01/2006")
		}

		data := documentListData{
			XSRFToken:             ctx.XSRFToken,
			SelectedCases:         selected,
			DocumentList:          docs,
			MultipleCasesSelected: len(caseUIDs) > 1 || (len(caseUIDs) == 0 && len(casesOnDonor) > 1),
			Error:                 validationErr,
			Success:               isSuccess,
			SuccessMessage:        successMessage,
		}

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
