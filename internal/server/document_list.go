package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type DocumentListClient interface {
	CasesByDonor(ctx sirius.Context, id int) ([]sirius.Case, error)
	GetPersonDocuments(ctx sirius.Context, personID int, caseIDs []string) (sirius.DocumentList, error)
}

type documentListData struct {
	XSRFToken             string
	Entity                string
	Success               bool
	Error                 sirius.ValidationError
	DocumentList          sirius.DocumentList
	SelectedCases         []sirius.Case
	MultipleCasesSelected bool
}

func ToUpper(s string) string {
	return strings.ToUpper(s)
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

		docs, err := client.GetPersonDocuments(ctx, donorID, caseIDs)
		if err != nil {
			return err
		}

		data := documentListData{
			XSRFToken:             ctx.XSRFToken,
			SelectedCases:         selected,
			DocumentList:          docs,
			MultipleCasesSelected: len(caseUIDs) > 1 || (len(caseUIDs) == 0 && len(casesOnDonor) > 1),
		}

		return tmpl(w, data)
	}
}
