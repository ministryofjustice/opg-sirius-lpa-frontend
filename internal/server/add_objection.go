package server

import (
	"fmt"
	"github.com/go-playground/form/v4"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/sync/errgroup"
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
	CaseUID    string
	LinkedLpas []sirius.SiriusData
	Form       formAddObjections
}

type formAddObjections struct {
	LpaUids       []string `form:"lpaUids"`
	ReceivedDate  dob      `form:"receivedDate"`
	ObjectionType string   `form:"objectionType"`
	Notes         string   `form:"notes"`
}

func isValidStatusForObjection(status string) bool {
	return status == "In progress" || status == "Draft" || status == "Statutory waiting period"
}

func AddObjection(client AddObjectionClient, tmpl template.Template) Handler {
	if decoder == nil {
		decoder = form.NewDecoder()
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		caseUID := r.FormValue("uid")
		ctx := getContext(r)

		var cs sirius.CaseSummary
		var err error

		group, groupCtx := errgroup.WithContext(ctx.Context)

		group.Go(func() error {
			cs, err = client.CaseSummary(ctx.With(groupCtx), caseUID)
			if err != nil {
				return err
			}
			return nil
		})

		if err := group.Wait(); err != nil {
			return err
		}

		var linkedCasesForObjections []sirius.SiriusData

		if isValidStatusForObjection(cs.DigitalLpa.SiriusData.Status) {
			linkedCasesForObjections = append(linkedCasesForObjections, cs.DigitalLpa.SiriusData)
		}

		for _, c := range cs.DigitalLpa.SiriusData.LinkedCases {
			if isValidStatusForObjection(c.Status) {
				linkedCasesForObjections = append(linkedCasesForObjections, c)
			}
		}

		data := addObjectionData{
			XSRFToken:  ctx.XSRFToken,
			CaseUID:    caseUID,
			LinkedLpas: linkedCasesForObjections,
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

				return RedirectError(fmt.Sprintf("/lpa/%s", caseUID))
			}
		}

		return tmpl(w, data)
	}
}
