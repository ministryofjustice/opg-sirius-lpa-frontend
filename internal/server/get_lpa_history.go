package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/sync/errgroup"
)

type GetLpaHistoryClient interface {
	RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error)
	GetEvents(ctx sirius.Context, donorId string, caseIds []string, sourceTypes []string, sortBy string) (sirius.LpaEventsResponse, error)
}

type getLpaHistory struct {
	XSRFToken           string
	DonorID             string
	Events              []sirius.LpaEvent
	EventFilterData     []sirius.SourceType
	Form                FilterLpaEventsForm
	TotalEvents         int
	TotalFilteredEvents int
	IsFiltered          bool
	FeeReductionTypes   []sirius.RefDataItem
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
		group, groupCtx := errgroup.WithContext(ctx.Context)

		data := getLpaHistory{
			XSRFToken: ctx.XSRFToken,
			DonorID:   donorId,
			Form: FilterLpaEventsForm{
				Sort: "desc",
			},
			IsFiltered: false,
		}

		group.Go(func() error {
			eventsData, err := client.GetEvents(ctx.With(groupCtx), donorId, caseIDs, []string{}, "desc")
			if err != nil {
				return err
			}
			data.Events = eventsData.Events
			data.EventFilterData = eventsData.Metadata.SourceTypes
			data.TotalEvents = eventsData.Total
			return nil
		})

		group.Go(func() error {
			feeReductionTypes, err := client.RefDataByCategory(ctx.With(groupCtx), sirius.FeeReductionTypeCategory)
			if err != nil {
				// If the call to get fee reduction types fails, we can just fall-back and display the database values instead of the labels
				data.FeeReductionTypes = []sirius.RefDataItem(nil)
				return nil
			}
			data.FeeReductionTypes = feeReductionTypes
			return nil
		})

		if err := group.Wait(); err != nil {
			return err
		}

		if r.Method == http.MethodPost {
			err := decoder.Decode(&data.Form, r.PostForm)
			if err != nil {
				return err
			}

			eventsData, err := client.GetEvents(ctx, donorId, caseIDs, data.Form.Types, data.Form.Sort)
			if err != nil {
				return err
			}

			data.TotalFilteredEvents = eventsData.Total
			data.Events = eventsData.Events
			data.IsFiltered = true
		}

		return tmpl(w, data)
	}
}
