package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type DeleteNoteClient interface {
	DeleteNote(ctx sirius.Context, id int) error
}

type deleteNoteData struct {
	XSRFToken string
	DonorId   int
}

func DeleteNote(client DeleteNoteClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		donorId, err := strToIntOrStatusError(r.FormValue("donorId"))
		if err != nil {
			return err
		}

		noteId, err := strToIntOrStatusError(r.FormValue("noteId"))
		if err != nil {
			return err
		}

		data := deleteNoteData{
			XSRFToken: ctx.XSRFToken,
			DonorId:   donorId,
		}

		if r.Method == http.MethodPost {
			err = client.DeleteNote(ctx, noteId)
			if err != nil {
				return err
			}

			return RedirectError(fmt.Sprintf("/donor/%d/history", donorId))
		}

		return tmpl(w, data)
	}
}
