package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
)

type notifyPageLimitData struct {
	CaseID              int
	CaseType            string
	CaseUUID            string
	CreatedDocumentUUID string
	DefaultDocumentUUID string
	DonorID             int
	IsPartial           bool
	NotifyWarning       string
}

func NotifyPageLimit(tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		donorId, err := strToIntOrStatusError(r.PathValue("donorId"))
		if err != nil {
			return err
		}
		createdDocumentUUID := r.PathValue("createdDocumentUUID")
		caseType := r.PathValue("caseType")
		caseUUID := r.PathValue("caseUUID")
		caseID, err := strToIntOrStatusError(r.PathValue("caseID"))
		if err != nil {
			return err
		}
		defaultDocumentUUID := r.URL.Query().Get("defaultDocumentUUID")

		data := notifyPageLimitData{
			CaseID:              caseID,
			CaseType:            caseType,
			CaseUUID:            caseUUID,
			CreatedDocumentUUID: createdDocumentUUID,
			DefaultDocumentUUID: defaultDocumentUUID,
			DonorID:             donorId,
			IsPartial:           ctx.IsPartial,
		}

		if err != nil {
			return err
		}
		data.NotifyWarning = "Document has been published but has not been posted via GOV.UK,Notify.\n" +
			"The published document has more than 10 pages and will not be posted automatically via GOV.UK,Notify. It must be manually uploaded to Xerox for posting."

		return tmpl(w, data)
	}
}
