package server

import (
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type WarningClient interface {
	RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error)
	CreateWarning(ctx sirius.Context, personId int, warningType, warningNote string, caseIDs []int) error
	CasesByDonor(ctx sirius.Context, id int) ([]sirius.Case, error)
}

type warningData struct {
	XSRFToken    string
	WarningTypes []sirius.RefDataItem
	Success      bool
	Error        sirius.ValidationError

	WarningType string
	WarningText string
	Cases       []sirius.Case
}

func Warning(client WarningClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		personId, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}

		ctx := getContext(r)
		warningTypes, err := client.RefDataByCategory(ctx, sirius.WarningTypeCategory)
		if err != nil {
			return err
		}

		cases, err := client.CasesByDonor(ctx, personId)
		if err != nil {
			return err
		}

		data := warningData{
			Success:      false,
			XSRFToken:    ctx.XSRFToken,
			WarningTypes: warningTypes,
			Cases:        cases,
		}

		if r.Method == http.MethodPost {
			warningType := postFormString(r, "warningType")
			warningText := postFormString(r, "warningText")

			var caseIDs = []int{}

			for _, id := range r.PostForm["case-id"] {
				intID, err := strconv.Atoi(id)
				if err != nil {
					return err
				}
				caseIDs = append(caseIDs, intID)
			}

			if len(data.Cases) == 1 {
				caseIDs = []int{data.Cases[0].ID}
			}

			err := client.CreateWarning(ctx, personId, warningType, warningText, caseIDs)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
				data.WarningType = warningType
				data.WarningText = warningText
			} else if err != nil {
				return err
			} else {
				data.Success = true
			}
		}

		return tmpl(w, data)
	}
}
