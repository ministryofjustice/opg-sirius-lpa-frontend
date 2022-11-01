package server

import (
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"net/http"
	"strconv"
)

type PlaceInvestigationOnHoldClient interface {
	PlaceInvestigationOnHold(ctx sirius.Context, investigationID int, reason string) error
	Investigation(ctx sirius.Context, id int) (sirius.Investigation, error)
}

type placeInvestigationOnHoldData struct {
	XSRFToken     string
	Success       bool
	Error         sirius.ValidationError
	Investigation sirius.Investigation
	Reason        string
}

func PlaceInvestigationOnHold(client PlaceInvestigationOnHoldClient, tmpl template.Template) Handler {
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

		data := placeInvestigationOnHoldData{
			XSRFToken:     ctx.XSRFToken,
			Investigation: investigation,
		}

		if r.Method == http.MethodPost {
			reason := postFormString(r, "reason")

			err = client.PlaceInvestigationOnHold(ctx, id, reason)

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
