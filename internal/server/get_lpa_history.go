package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type GetLpaHistoryClient interface {
	GetEvents(ctx sirius.Context, donorId string, caseIds []string, sourceTypes []string, sortBy string) (sirius.LpaEventsResponse, error)
}

type getLpaHistory struct {
	XSRFToken           string
	DonorID             string
	Events              []LpaEventWithContext
	EventFilterData     []sirius.SourceType
	Form                FilterLpaEventsForm
	TotalEvents         int
	TotalFilteredEvents int
	IsFiltered          bool
}

type LpaEventWithContext struct {
	sirius.LpaEvent
	DonorID string
}

type FilterLpaEventsForm struct {
	Types []string `form:"type"`
	Sort  string   `form:"sort"`
}

func GetLpaHistory(client GetLpaHistoryClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		donorId := r.PathValue("donorId")
		caseIDs := r.URL.Query()["id[]"]

		ctx := getContext(r)

		eventsData, err := client.GetEvents(ctx, donorId, caseIDs, []string{}, "desc")
		if err != nil {
			return err
		}
		eventsWithContext := make([]LpaEventWithContext, len(eventsData.Events))
		for i, event := range eventsData.Events {
			eventsWithContext[i] = LpaEventWithContext{
				LpaEvent: event,
				DonorID:  donorId,
			}
		}

		data := getLpaHistory{
			XSRFToken:       ctx.XSRFToken,
			DonorID:         donorId,
			Events:          eventsWithContext,
			EventFilterData: eventsData.Metadata.SourceTypes,
			TotalEvents:     eventsData.Total,
			IsFiltered:      false,
			Form: FilterLpaEventsForm{
				Sort: "desc",
			},
		}

		if r.Method == http.MethodPost {
			err := decoder.Decode(&data.Form, r.PostForm)
			if err != nil {
				return err
			}

			eventsData, err = client.GetEvents(ctx, donorId, caseIDs, data.Form.Types, data.Form.Sort)
			if err != nil {
				return err
			}

			data.TotalFilteredEvents = eventsData.Total
			eventsWithContext = make([]LpaEventWithContext, len(eventsData.Events))
			for i, event := range eventsData.Events {
				eventsWithContext[i] = LpaEventWithContext{
					LpaEvent: event,
					DonorID:  donorId,
				}
			}
			data.Events = eventsWithContext
			data.IsFiltered = true
		}

		return tmpl(w, data)
	}
}
