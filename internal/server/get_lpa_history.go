package server

import (
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type GetLpaHistoryClient interface {
	GetEvents(ctx sirius.Context, donorId int) (sirius.LpaEventsResponse, error)
}

type getLpaHistory struct {
	Events any
}

func GetLpaHistory(client GetLpaHistoryClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		donorId := r.PathValue("donorId")

		ctx := getContext(r)

		donorID, err := strconv.Atoi(donorId)

		eventsData, err := client.GetEvents(ctx, donorID)
		if err != nil {
			return err
		}

		data := getLpaHistory{
			Events: eventsData.Events,
		}

		return tmpl(w, data)
	}
}
