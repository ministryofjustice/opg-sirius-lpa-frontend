package server

import (
	"fmt"
	"github.com/go-playground/form/v4"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"net/http"
)

type AddObjectionClient interface {
	CaseSummary(sirius.Context, string) (sirius.CaseSummary, error)
	AddObjections(sirius.Context, sirius.AddObjections) error
}

type addObjectionData struct {
	XSRFToken  string
	Success    bool
	Error      sirius.ValidationError
	Case       sirius.CaseSummary
	LinkedLpas []sirius.SiriusData
	Form       formAddObjections
}

type formAddObjections struct {
	LpaUids       []string `form:"lpaUids"`
	ReceivedDate  dob      `form:"receivedDate"`
	ObjectionType string   `form:"objectionType"`
	Notes         string   `form:"notes"`
}

func AddObjection(client AddObjectionClient, tmpl template.Template) Handler {
	if decoder == nil {
		decoder = form.NewDecoder()
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		caseUID := r.FormValue("uid")

		ctx := getContext(r)

		cs, err := client.CaseSummary(ctx, caseUID)
		if err != nil {
			return err
		}

		data := addObjectionData{
			XSRFToken:  ctx.XSRFToken,
			Case:       cs,
			LinkedLpas: cs.DigitalLpa.SiriusData.LinkedCases,
		}

		if r.Method == http.MethodPost {
			err := decoder.Decode(&data.Form, r.PostForm)
			if err != nil {
				return err
			}

			objection := sirius.AddObjections{
				LpaUids:       data.Form.LpaUids,
				ReceivedDate:  data.Form.ReceivedDate.toDateString(),
				ObjectionType: data.Form.ObjectionType,
				Notes:         data.Form.Notes,
			}

			err = client.AddObjections(ctx, objection)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
			} else if err != nil {
				return err
			} else {
				data.Success = true

				SetFlash(w, FlashNotification{
					Title: "Update saved",
				})

				return RedirectError(fmt.Sprintf("/lpa/%s/lpa-details", caseUID))
			}
		}

		return tmpl(w, data)
	}
}
