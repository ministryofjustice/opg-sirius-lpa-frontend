package server

import (
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type CompareDocWithDocListClient interface {
	DocumentByUUID(ctx sirius.Context, uuid string) (sirius.Document, error)
	GetPersonDocuments(ctx sirius.Context, personID int, caseIDs []string) (sirius.DocumentList, error)
}

func CompareDocWithDocList(client CompareDocWithDocListClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		donorID, err := strToIntOrStatusError(r.PathValue("id"))
		if err != nil {
			return err
		}

		docUUIDs := r.URL.Query()["uid[]"]
		ctx := getContext(r)

		documentData, err := client.DocumentByUUID(ctx, docUUIDs[0])
		if err != nil {
			return err
		}

		caseItems := []sirius.Case{}
		if len(documentData.CaseItems) > 0 {
			caseItems = documentData.CaseItems
		}

		caseIDs := []string{}
		for _, caseItem := range caseItems {
			caseIDs = append(caseIDs, strconv.Itoa(caseItem.ID))
		}

		docs, err := client.GetPersonDocuments(ctx, donorID, caseIDs)
		if err != nil {
			return err
		}

		data := documentPageData{
			XSRFToken:     ctx.XSRFToken,
			DocumentList:  docs,
			Document:      documentData,
			SelectedCases: caseItems,
			Comparing:     true,
		}

		return tmpl(w, data)
	}
}
