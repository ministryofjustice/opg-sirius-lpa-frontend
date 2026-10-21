package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type NotifyErrorClient interface {
	DocumentByUUID(ctx sirius.Context, uuid string) (sirius.Document, error)
}

type notifyErrorData struct {
	XSRFToken string
	IsPartial bool
	Success   bool
	Error     sirius.ValidationError
	Case      sirius.Case

	CreatedDocument      sirius.Document
	NotifyWarning        string
	DocumentFriendlyName string
	DonorID              int
	CaseUUID             string
	SuccessMessage       string
	CaseType             string
}

func NotifyError(client NotifyErrorClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		donorId, err := strToIntOrStatusError(r.PathValue("donorId"))
		if err != nil {
			return err
		}
		docUUID := r.PathValue("docUUID")
		caseType := r.PathValue("caseType")
		caseUUID := r.PathValue("caseUUID")

		data := notifyErrorData{
			XSRFToken: ctx.XSRFToken,
			IsPartial: ctx.IsPartial,
			DonorID:   donorId,
			CaseUUID:  caseUUID,
			CaseType:  caseType,
		}

		data.NotifyWarning = "The published document has more than 10 pages and will not be sent automatically via GOV.UK,Notify. It must be manually uploaded to Xerox for sending."

		document, err := client.DocumentByUUID(ctx, docUUID)
		if err != nil {
			return err
		}
		data.SuccessMessage = "Document has been published named " + document.FriendlyDescription

		return tmpl(w, data)
	}
}
