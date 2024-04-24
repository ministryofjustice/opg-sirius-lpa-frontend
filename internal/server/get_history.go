package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
)

type GetHistoryClient interface {
	GetEvents(ctx sirius.Context, donorId int, caseId int) (any, error)
	CaseSummary(ctx sirius.Context, uid string) (sirius.CaseSummary, error)
}

type getHistory struct {
	CaseSummary sirius.CaseSummary
	EventData   any
}

func GetHistory(client GetHistoryClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		uid := chi.URLParam(r, "uid")
		ctx := getContext(r)

		caseSummary, err := client.CaseSummary(ctx, uid)
		if err != nil {
			return err
		}

		caseId := caseSummary.DigitalLpa.SiriusData.ID

		eventDetails, err := client.GetEvents(ctx, caseSummary.DigitalLpa.SiriusData.Donor.ID, caseId)
		if err != nil {
			return err
		}

		data := getHistory{
			CaseSummary: caseSummary,
			EventData:   eventDetails,
		}

		return tmpl(w, data)
	}
}
