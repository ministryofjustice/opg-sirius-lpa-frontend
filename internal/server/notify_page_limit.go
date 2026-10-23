package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type NotifyPageLimitClient interface {
	DocumentByUUID(ctx sirius.Context, uuid string) (sirius.Document, error)
}

type notifyPageLimitData struct {
	CaseID              int
	CaseType            string
	CaseUUID            string
	CreatedDocumentUUID string
	DefaultDocumentUUID string
	DonorID             int
	Error               sirius.ValidationError
	IsPartial           bool
	NotifyWarning       string
	XSRFToken           string
}

func NotifyPageLimit(client NotifyPageLimitClient, tmpl template.Template) Handler {
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
			XSRFToken:           ctx.XSRFToken,
			IsPartial:           ctx.IsPartial,
			DonorID:             donorId,
			CaseID:              caseID,
			CaseType:            caseType,
			DefaultDocumentUUID: defaultDocumentUUID,
			CaseUUID:            caseUUID,
		}

		document, err := client.DocumentByUUID(ctx, createdDocumentUUID)
		if err != nil {
			return err
		}
		data.CreatedDocumentUUID = document.UUID
		data.NotifyWarning = "Document has been published named " + document.FriendlyDescription +
			". The published document has more than 10 pages and will not be sent automatically via GOV.UK,Notify. It must be manually uploaded to Xerox for sending."

		return tmpl(w, data)
	}
}
