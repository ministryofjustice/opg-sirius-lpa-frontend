package server

import (
	"github.com/go-playground/form/v4"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/sync/errgroup"
	"net/http"
)

type UpdateObjectionClient interface {
	CaseSummary(sirius.Context, string) (sirius.CaseSummary, error)
	GetObjection(sirius.Context, string) (sirius.Objection, error)
	//UpdateObjection(sirius.Context, string, sirius.ChangeDonorDetails) error
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

		return tmpl(w, data)
	}
}
