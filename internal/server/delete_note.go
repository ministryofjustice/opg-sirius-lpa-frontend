package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type DeleteNoteClient interface {
	GetEvents(ctx sirius.Context, donorId string, caseIds []string, sourceTypes []string, eventIds []string, sortBy string) (sirius.LpaEventsResponse, error)
	DeleteNote(ctx sirius.Context, id int) error
}

type deleteNoteData struct {
	XSRFToken string
	DonorId   string
	Event     sirius.LpaEvent
}

func DeleteNote(client DeleteNoteClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		donorId := r.FormValue("donorId")
		eventId := r.FormValue("eventId")

		eventsResp, err := client.GetEvents(ctx, donorId, []string{}, []string{}, []string{eventId}, "desc")
		if err != nil {
			return err
		}

		noteId := eventsResp.Events[0].SourceNote.ID

		data := deleteNoteData{
			XSRFToken: ctx.XSRFToken,
			DonorId:   donorId,
			Event:     eventsResp.Events[0],
		}

		if r.Method == http.MethodPost {
			err = client.DeleteNote(ctx, noteId)
			if err != nil {
				return err
			}

			return RedirectError(fmt.Sprintf("/donor/%s/history", donorId))
		}

		return tmpl(w, data)
	}
}
