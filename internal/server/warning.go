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
	Person(ctx sirius.Context, personID int) (sirius.Person, error)
}

type warningData struct {
	XSRFToken    string
	WarningTypes []sirius.RefDataItem
	Success      bool
	Error        sirius.ValidationError

	WarningType string
	WarningText string
	Donor       sirius.Person
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

		donor, err := client.Person(ctx, personId)
		if err != nil {
			return err
		}

		data := warningData{
			Success:      false,
			XSRFToken:    ctx.XSRFToken,
			WarningTypes: warningTypes,
			Donor:        donor,
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

			if len(donor.Cases) == 1 {
				caseIDs = []int{donor.Cases[0].ID}
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
