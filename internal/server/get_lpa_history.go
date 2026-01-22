package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type GetLpaHistoryClient interface {
	GetEvents(ctx sirius.Context, donorId string, caseIds []string) (sirius.LpaEventsResponse, error)
}

type getLpaHistory struct {
	Events any
}

func GetLpaHistory(client GetLpaHistoryClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		donorId := r.PathValue("donorId")

		if err := r.ParseForm(); err != nil {
			return err
		}
		caseIDs := r.Form["id[]"]

		ctx := getContext(r)

		eventsData, err := client.GetEvents(ctx, donorId, caseIDs)
		if err != nil {
			return err
		}

		data := getLpaHistory{
			Events: eventsData.Events,
		}

		return tmpl(w, data)
	}
}
