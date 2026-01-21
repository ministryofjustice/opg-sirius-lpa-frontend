package server

import (
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type GetLpaHistoryClient interface {
	Case(ctx sirius.Context, id int) (sirius.Case, error)
	GetEvents(ctx sirius.Context, donorId int, caseId int, sourceTypes []string, sortBy string) (sirius.LpaEvents, error)
}

type EventFilter struct {
	Type  string
	Count int
}

type getLpaHistory struct {
	XSRFToken       string
	Form            LpaHistoryForm
	Case            sirius.Case
	EventData       any
	EventFilterData []EventFilter
}

type LpaHistoryForm struct {
	Types []string `form:"type"`
	Sort  string   `form:"sort"`
}

func GetLpaHistory(client GetLpaHistoryClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		donorId := r.PathValue("donorId")
		caseId := r.PathValue("caseId")

		ctx := getContext(r)

		caseID, err := strconv.Atoi(caseId)
		donorID, err := strconv.Atoi(donorId)

		caseInfo, err := client.Case(ctx, caseID)

		allEvents, err := client.GetEvents(ctx, donorID, caseID, []string{}, "desc")
		if err != nil {
			return err
		}

		countEvents := make(map[string]int)
		for _, event := range allEvents {
			countEvents[event.SourceType]++
		}

		filters := make([]EventFilter, 0, len(countEvents))
		for sourceType, count := range countEvents {
			filters = append(filters, EventFilter{
				Type:  sourceType,
				Count: count,
			})
		}

		eventDetails := allEvents
		var selectedTypes []string

		data := getLpaHistory{
			Case:            caseInfo,
			XSRFToken:       ctx.XSRFToken,
			EventData:       eventDetails,
			EventFilterData: filters,
			Form: LpaHistoryForm{
				Sort: "desc",
			},
		}

		if r.Method == http.MethodPost {

			err := decoder.Decode(&data.Form, r.PostForm)
			if err != nil {
				return err
			}

			selectedTypes = data.Form.Types
			sortBy := data.Form.Sort

			eventDetails, err = client.GetEvents(ctx, donorID, caseID, selectedTypes, sortBy)
			if err != nil {
				return err
			}

			data.EventData = eventDetails
		}

		return tmpl(w, data)
	}
}
