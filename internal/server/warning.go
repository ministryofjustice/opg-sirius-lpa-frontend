package server

import (
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type WarningClient interface {
	WarningTypes(ctx sirius.Context) ([]sirius.RefDataItem, error)
	CreateWarning(ctx sirius.Context, personId int, warningType, warningNote string) error
}

type warningData struct {
	XSRFToken    string
	WarningTypes []sirius.RefDataItem
	Success      bool
	Error        sirius.ValidationError
}

func Warning(client WarningClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		personId, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}

		ctx := getContext(r)

		warningTypes, err := client.WarningTypes(ctx)
		if err != nil {
			return err
		}

		data := warningData{
			Success:      false,
			XSRFToken:    ctx.XSRFToken,
			WarningTypes: warningTypes,
		}

		if r.Method == http.MethodPost {
			warningType := postFormString(r, "warning-type")
			warningNotes := postFormString(r, "warning-notes")

			err := client.CreateWarning(ctx, personId, warningType, warningNotes)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
			} else if err != nil {
				return err
			} else {
				data.Success = true
			}
		}

		return tmpl(w, data)
	}
}
