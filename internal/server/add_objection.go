package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type AddObjectionClient interface {
	CaseSummary(sirius.Context, string) (sirius.CaseSummary, error)
	AddObjection(sirius.Context, sirius.ObjectionRequest) error
}

type addObjectionData struct {
	XSRFToken  string
	Title      string
	Success    bool
	Error      sirius.ValidationError
	CaseUID    string
	LinkedLpas []sirius.SiriusData
	Form       formObjection
}

type formObjection struct {
	LpaUids       []string `form:"lpaUids"`
	ReceivedDate  dob      `form:"receivedDate"`
	ObjectionType string   `form:"objectionType"`
	Notes         string   `form:"notes"`
}

func AddObjection(client AddObjectionClient, tmpl template.Template) Handler {

	return func(w http.ResponseWriter, r *http.Request) error {
		caseUID := r.FormValue("uid")
		ctx := getContext(r)

		cs, err := client.CaseSummary(ctx, caseUID)
		if err != nil {
			return err
		}

		var linkedCasesForObjections []sirius.SiriusData

		cs.DigitalLpa.SiriusData.Status.IsValidStatusForObjection() {
			linkedCasesForObjections = append(linkedCasesForObjections, cs.DigitalLpa.SiriusData)
		}

		for _, c := range cs.DigitalLpa.SiriusData.LinkedCases {
			if c.Status.IsValidStatusForObjection() {
				linkedCasesForObjections = append(linkedCasesForObjections, c)
			}
		}

		data := addObjectionData{
			XSRFToken:  ctx.XSRFToken,
			Title:      "Add Objection",
			CaseUID:    caseUID,
			LinkedLpas: linkedCasesForObjections,
		}

		if r.Method == http.MethodPost {
			err := decoder.Decode(&data.Form, r.PostForm)
			if err != nil {
				return err
			}

			objection := sirius.ObjectionRequest{
				LpaUids:       data.Form.LpaUids,
				ReceivedDate:  data.Form.ReceivedDate.toDateString(),
				ObjectionType: data.Form.ObjectionType,
				Notes:         data.Form.Notes,
			}

			err = client.AddObjection(ctx, objection)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
			} else if err != nil {
				return err
			} else {
				data.Success = true

				return RedirectError(fmt.Sprintf("/lpa/%s", caseUID))
			}
		}

		return tmpl(w, data)
	}
}
