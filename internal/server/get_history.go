package server

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"log"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
)

type GetHistoryClient interface {
	GetEvents(ctx sirius.Context, id int, caseId string) (any, error)
	CaseSummary(ctx sirius.Context, uid string) (sirius.CaseSummary, error)
}

type getHistory struct {
	CaseSummary  sirius.CaseSummary
	DigitalLpa   sirius.DigitalLpa
	LpaStoreData map[string]interface{}
	EventData    any
}

func GetHistory(client GetHistoryClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		uid := chi.URLParam(r, "uid")
		ctx := getContext(r)

		caseSummary, err := client.CaseSummary(ctx, uid)
		if err != nil {
			return err
		}

		var lpaStoreData map[string]interface{}
		err = json.Unmarshal(caseSummary.DigitalLpa.LpaStoreData, &lpaStoreData)
		if err != nil {
			return err
		}

		eventDetails, err := client.GetEvents(ctx, caseSummary.DigitalLpa.SiriusData.Donor.ID, uid)

		var eventDump map[string]interface{}
		//err = json.Unmarshal(eventDetails, &eventDump)
		if err != nil {
			return err
		}

		data := getHistory{
			CaseSummary:  caseSummary,
			DigitalLpa:   caseSummary.DigitalLpa,
			LpaStoreData: lpaStoreData,
			EventData:    &eventDetails,
		}

		log.Print("MISH")
		log.Print(eventDump)

		return tmpl(w, data)
	}
}
