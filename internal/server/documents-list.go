package server

import (
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type DocumentsListClient interface {
	Case(ctx sirius.Context, id int) (sirius.Case, error)
	CasesByDonor(ctx sirius.Context, id int) ([]sirius.Case, error)
	GetPersonDocuments(ctx sirius.Context, personID string, caseIDs []string) (sirius.DocumentList, error)
}

type documentListData struct {
	XSRFToken     string
	Entity        string
	Success       bool
	Error         sirius.ValidationError
	DonorId       string
	CaseUIDs      []string
	Complaint     sirius.Complaint
	Documents     sirius.DocumentList
	Cases         []sirius.Case
	SelectedCases []SelectedCaseForDocuments
}

type SelectedCaseForDocuments struct {
	UID     string
	ID      string
	Subtype string
	Type    string
}

func DocumentsList(client DocumentsListClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if err := r.ParseForm(); err != nil {
			return err
		}

		donorID := r.PathValue("id")
		caseUIDs := r.Form["uid[]"]

		ctx := getContext(r)

		donorInt, err := strconv.Atoi(donorID)
		if err != nil {
			return err
		}

		cases, err := client.CasesByDonor(ctx, donorInt)
		if err != nil {
			return err
		}

		casesByUID := make(map[string]sirius.Case, len(cases))
		for _, c := range cases {
			casesByUID[c.UID] = c
		}

		var selected []SelectedCaseForDocuments
		var caseIDs []string
		for _, uid := range caseUIDs {
			if c, ok := casesByUID[uid]; ok {
				selected = append(selected, SelectedCaseForDocuments{
					UID:     c.UID,
					ID:      strconv.Itoa(c.ID),
					Type:    c.CaseType,
					Subtype: c.SubType,
				})
				caseIDs = append(caseIDs, strconv.Itoa(c.ID))
			}
		}

		docs, err := client.GetPersonDocuments(ctx, donorID, caseIDs)
		if err != nil {
			return err
		}

		data := documentListData{
			XSRFToken:     ctx.XSRFToken,
			DonorId:       donorID,
			Cases:         cases,
			SelectedCases: selected,
			Documents:     docs,
		}

		return tmpl(w, data)
	}
}
