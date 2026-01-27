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
	Events []sirius.LpaEvent
}

func GetLpaHistory(client GetLpaHistoryClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		donorId := r.PathValue("donorId")
		caseIDs := r.URL.Query()["id[]"]

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
