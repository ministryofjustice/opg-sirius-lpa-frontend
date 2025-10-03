package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"

	"github.com/ministryofjustice/opg-go-common/template"
)

type GetHistoryClient interface {
	GetEvents(ctx sirius.Context, donorId int, caseId int) (any, error)
	GetCombinedEvents(ctx sirius.Context, uid string) (any, error)
	CaseSummary(ctx sirius.Context, uid string) (sirius.CaseSummary, error)
}

type getHistory struct {
	CaseSummary sirius.CaseSummary
	EventData   any
}

func GetHistory(client GetHistoryClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		uid := r.PathValue("uid")
		ctx := getContext(r)

		caseSummary, err := client.CaseSummary(ctx, uid)
		if err != nil {
			return err
		}

		var eventDetails any

		if caseSummary.DigitalLpa.LpaStoreData.Status != "" {
			// Digital LPA - use combined events
			eventDetails, err = client.GetCombinedEvents(ctx, uid)
		} else {
			// Traditional LPA - use Sirius events only
			donorId := caseSummary.DigitalLpa.SiriusData.Donor.ID
			caseId := caseSummary.DigitalLpa.SiriusData.ID
			eventDetails, err = client.GetEvents(ctx, donorId, caseId)
		}

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
