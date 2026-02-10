package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type CompareDocListWithDocListClient interface {
	DocumentByUUID(ctx sirius.Context, uuid string) (sirius.Document, error)
	GetPersonDocuments(ctx sirius.Context, personID int, caseIDs []string) (sirius.DocumentList, error)
}

func CompareDocListWithDocList(client CompareDocListWithDocListClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		donorID, err := strToIntOrStatusError(r.PathValue("id"))
		if err != nil {
			return err
		}

		caseID := r.PathValue("caseId")
		ctx := getContext(r)

		docs, err := client.GetPersonDocuments(ctx, donorID, []string{caseID})
		if err != nil {
			return err
		}

		selected := docs.Documents[0].CaseItems

		data := documentPageData{
			XSRFToken:     ctx.XSRFToken,
			DocumentList:  docs,
			SelectedCases: selected,
			Comparing:     true,
		}

		return tmpl(w, data)
	}
}
