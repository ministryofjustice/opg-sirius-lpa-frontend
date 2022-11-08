package server

import (
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"net/http"
	"strconv"
)

type EditInvestigationClient interface {
	EditInvestigation(ctx sirius.Context, investigationID int, investigation sirius.Investigation) error
	Investigation(ctx sirius.Context, id int) (sirius.Investigation, error)
}

type editInvestigationData struct {
	XSRFToken            string
	Success              bool
	Error                sirius.ValidationError
	Investigation        sirius.Investigation
	ApprovalOutcomeTypes []string
}

var approvalOutcomeTypes = []string{"Court Application", "Further Action", "No Further Action", "Instrument no longer Valid"}

func EditInvestigation(client EditInvestigationClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		investigationID, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}

		ctx := getContext(r)

		data := editInvestigationData{
			XSRFToken:            ctx.XSRFToken,
			ApprovalOutcomeTypes: approvalOutcomeTypes,
		}

		if r.Method == http.MethodPost {
			investigation := sirius.Investigation{
				Title:                    postFormString(r, "title"),
				Information:              postFormString(r, "information"),
				Type:                     postFormString(r, "type"),
				DateReceived:             postFormDateString(r, "dateReceived"),
				ApprovalDate:             postFormDateString(r, "approvalDate"),
				ApprovalOutcome:          postFormString(r, "approvalOutcome"),
				InvestigationClosureDate: postFormDateString(r, "investigationClosureDate"),
			}

			err = client.EditInvestigation(ctx, investigationID, investigation)

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

		investigation, err := client.Investigation(ctx, investigationID)
		if err != nil {
			return err
		}
		data.Investigation = investigation

		return tmpl(w, data)
	}
}
