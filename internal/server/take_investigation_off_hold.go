package server

import (
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"net/http"
	"strconv"
)

type TakeInvestigationOffHoldClient interface {
	TakeInvestigationOffHold(ctx sirius.Context, investigationID int) error
	HoldPeriod(ctx sirius.Context, id int) (sirius.HoldPeriod, error)
}

type takeInvestigationOffHoldData struct {
	XSRFToken  string
	Success    bool
	Error      sirius.ValidationError
	HoldPeriod sirius.HoldPeriod
}

func TakeInvestigationOffHold(client TakeInvestigationOffHoldClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}

		ctx := getContext(r)
		holdPeriod, err := client.HoldPeriod(ctx, id)

		if err != nil {
			return err
		}

		data := takeInvestigationOffHoldData{
			XSRFToken:  ctx.XSRFToken,
			HoldPeriod: holdPeriod,
		}

		if r.Method == http.MethodPost {
			err = client.TakeInvestigationOffHold(ctx, id)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
				data.HoldPeriod = holdPeriod
			} else if err != nil {
				return err
			} else {
				data.Success = true
			}
		}

		return tmpl(w, data)
	}
}
