package server

import (
	"errors"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"net/http"
	"strconv"
)

type InvestigationHoldClient interface {
	PlaceInvestigationOnHold(ctx sirius.Context, investigationID int, reason string) error
	TakeInvestigationOffHold(ctx sirius.Context, investigationID int) error
	Investigation(ctx sirius.Context, id int) (sirius.Investigation, error)
}

type investigationHoldData struct {
	XSRFToken     string
	Success       bool
	Error         sirius.ValidationError
	Investigation sirius.Investigation
	Reason        string
}

func InvestigationHold(client InvestigationHoldClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}

		ctx := getContext(r)
		investigation, err := client.Investigation(ctx, id)
		if err != nil {
			return err
		}

		data := investigationHoldData{
			XSRFToken:     ctx.XSRFToken,
			Investigation: investigation,
		}

		var hpID int
		if investigation.IsOnHold {
			for _, hp := range investigation.HoldPeriods {
				if hp.EndDate == "" {
					hpID = hp.ID
					data.Reason = hp.Reason
				}
			}
			if hpID == 0 {
				w.WriteHeader(http.StatusInternalServerError)
				return errors.New("could not find open hold period on investigation")
			}
		}

		if r.Method == http.MethodPost {
			if investigation.IsOnHold {
				err = client.TakeInvestigationOffHold(ctx, hpID)
			} else {
				reason := postFormString(r, "reason")
				data.Reason = reason
				err = client.PlaceInvestigationOnHold(ctx, id, reason)
			}

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
				data.Investigation = investigation
			} else if err != nil {
				return err
			} else {
				data.Success = true
			}

		}

		return tmpl(w, data)
	}
}
