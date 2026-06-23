package server

import (
	"encoding/json"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type SiriusHeaderCalendarClient interface {
	BankHolidays(ctx sirius.Context) (sirius.BankHolidays, error)
}

type siriusHeaderCalendarData struct {
	BankHolidaysJSON string
}

func SiriusHeaderCalendars(client SiriusHeaderCalendarClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		bankHolidays, err := client.BankHolidays(ctx)
		if err != nil {
			bankHolidays = sirius.BankHolidays{}
		}

		bankHolidaysJSON, err := json.Marshal(bankHolidays)
		if err != nil {
			bankHolidaysJSON = []byte("{}")
		}

		data := siriusHeaderCalendarData{
			BankHolidaysJSON: string(bankHolidaysJSON),
		}

		return tmpl(w, data)
	}
}
