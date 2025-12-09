package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type DocumentsListClient interface {
	CasesByDonor(ctx sirius.Context, id int) ([]sirius.Case, error)
	GetPersonDocuments(ctx sirius.Context, personID int, caseIDs []string) (sirius.DocumentList, error)
}

type documentListData struct {
	XSRFToken             string
	Entity                string
	Success               bool
	Error                 sirius.ValidationError
	DocumentList          sirius.DocumentList
	SelectedCases         []SelectedCaseForDocuments
	MultipleCasesSelected bool
}

type SelectedCaseForDocuments struct {
	UID      string
	SubType  string
	CaseType string
}

func ToUpper(s string) string {
	return strings.ToUpper(s)
}

func DocumentsList(client DocumentsListClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		donorID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			return err
		}

		caseUIDs := r.Form["uid[]"]
		multipleCasesForSelection := false

		ctx := getContext(r)

		casesOnDonor, err := client.CasesByDonor(ctx, donorID)
		if err != nil {
			return err
		}

		var selected []SelectedCaseForDocuments
		var caseIDs []string

		if len(caseUIDs) > 0 {
			if len(caseUIDs) > 1 {
				multipleCasesForSelection = true
			}

			casesByUID := make(map[string]sirius.Case, len(casesOnDonor))
			for _, c := range casesOnDonor {
				casesByUID[c.UID] = c
			}

			for _, uid := range caseUIDs {
				if c, ok := casesByUID[uid]; ok {
					selected = append(selected, SelectedCaseForDocuments{
						UID:      c.UID,
						CaseType: c.CaseType,
						SubType:  c.SubType,
					})
					caseIDs = append(caseIDs, strconv.Itoa(c.ID))
				}
			}

		} else {
			if len(casesOnDonor) > 1 {
				multipleCasesForSelection = true
			}

			for _, c := range casesOnDonor {
				selected = append(selected, SelectedCaseForDocuments{
					UID:      c.UID,
					CaseType: c.CaseType,
					SubType:  c.SubType,
				})
			}
		}

		docs, err := client.GetPersonDocuments(ctx, donorID, caseIDs)
		if err != nil {
			return err
		}

		data := documentListData{
			XSRFToken:             ctx.XSRFToken,
			SelectedCases:         selected,
			DocumentList:          docs,
			MultipleCasesSelected: multipleCasesForSelection,
		}

		return tmpl(w, data)
	}
}
