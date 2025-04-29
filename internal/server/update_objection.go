package server

import (
	"fmt"
	"github.com/go-playground/form/v4"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/sync/errgroup"
	"net/http"
)

type UpdateObjectionClient interface {
	CaseSummary(sirius.Context, string) (sirius.CaseSummary, error)
	GetObjection(sirius.Context, string) (sirius.Objection, error)
	UpdateObjection(sirius.Context, string, sirius.AddObjection) error
}

type updateObjectionData struct {
	XSRFToken  string
	Success    bool
	Error      sirius.ValidationError
	CaseUID    string
	LinkedLpas []sirius.SiriusData
	Form       formUpdateObjection
}

type formUpdateObjection struct {
	LpaUids       []string `form:"lpaUids"`
	ReceivedDate  dob      `form:"receivedDate"`
	ObjectionType string   `form:"objectionType"`
	Notes         string   `form:"notes"`
}

func UpdateObjection(client UpdateObjectionClient, tmpl template.Template) Handler {
	if decoder == nil {
		decoder = form.NewDecoder()
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		caseUID := r.PathValue("uid")
		objectionID := r.PathValue("id")

		ctx := getContext(r)

		var cs sirius.CaseSummary
		var obj sirius.Objection
		var err error

		group, groupCtx := errgroup.WithContext(ctx.Context)

		group.Go(func() error {
			cs, err = client.CaseSummary(ctx.With(groupCtx), caseUID)
			if err != nil {
				return err
			}
			return nil
		})

		group.Go(func() error {
			obj, err = client.GetObjection(ctx.With(groupCtx), objectionID)
			if err != nil {
				return err
			}
			return nil
		})

		if err := group.Wait(); err != nil {
			return err
		}

		receivedDate, err := parseDate(obj.ReceivedDate)
		if err != nil {
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

		data := updateObjectionData{
			XSRFToken:  ctx.XSRFToken,
			CaseUID:    caseUID,
			LinkedLpas: linkedCasesForObjections,
			Form: formUpdateObjection{
				LpaUids:       obj.LpaUids,
				ReceivedDate:  receivedDate,
				ObjectionType: obj.ObjectionType,
				Notes:         obj.Notes,
			},
		}

		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				return err
			}

			data.Form = formUpdateObjection{}

			err := decoder.Decode(&data.Form, r.PostForm)
			if err != nil {
				return err
			}

			objection := sirius.AddObjection{
				LpaUids:       data.Form.LpaUids,
				ReceivedDate:  data.Form.ReceivedDate.toDateString(),
				ObjectionType: data.Form.ObjectionType,
				Notes:         data.Form.Notes,
			}

			err = client.UpdateObjection(ctx, objectionID, objection)

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
