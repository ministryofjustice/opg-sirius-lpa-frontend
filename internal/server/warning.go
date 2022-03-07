package server

import (
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type WarningClient interface {
	WarningTypes(ctx sirius.Context) ([]sirius.RefDataItem, error)
	CreateWarning(ctx sirius.Context, personId int, warningType, warningNote string) error
}

type warningData struct {
	WasWarningCreated bool
	Error             error
	ValidationErr     sirius.ValidationError
	XSRFToken         string
	WarningTypes      []sirius.RefDataItem
}

func Warning(client WarningClient, t Template) Handler {
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
			WasWarningCreated: false,
			Error:             err,
			XSRFToken:         ctx.XSRFToken,
			WarningTypes:      warningTypes,
		}

		if r.Method == http.MethodPost {
			warningType := r.FormValue("warning-type")
			warningNotes := r.FormValue("warning-notes")

			err := client.CreateWarning(ctx, personId, warningType, warningNotes)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.ValidationErr = ve
			} else if err != nil {
				return err
			} else {
				data.WasWarningCreated = true
			}
		}

		return t.ExecuteTemplate(w, "page", data)
	}
}
