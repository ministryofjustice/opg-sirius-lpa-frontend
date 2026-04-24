package server

import (
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/sync/errgroup"
)

type ActionPanelClient interface {
	NoteTypes(ctx sirius.Context) ([]string, error)
}

type ActionPanelData struct {
	XSRFToken    string
	NoteTypes    []string
	Entity       string
	IsDigitalLpa bool
	CaseUID      string
	Success      bool
	Error        sirius.ValidationError
	DonorID      int

	Type        string
	Name        string
	Description string
}

func ActionPanel(client ActionPanelClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)
		data := ActionPanelData{XSRFToken: ctx.XSRFToken}

		// Support both query-string usage (?donorId=123) and path-based params
		donorIDStr := r.URL.Query().Get("donorId")
		if donorIDStr == "" {
			donorIDStr = r.URL.Query().Get("id")
		}
		if donorIDStr != "" {
			if donorID, err := strconv.Atoi(donorIDStr); err == nil {
				data.DonorID = donorID
			}
		}

		group, groupCtx := errgroup.WithContext(ctx.Context)

		group.Go(func() error {
			noteTypes, err := client.NoteTypes(ctx.With(groupCtx))
			if err != nil {
				return err
			}

			data.NoteTypes = noteTypes
			return nil
		})

		if err := group.Wait(); err != nil {
			return err
		}

		//if r.Header.Get("HX-Request") == "true" && partialTmpl != nil {
		//	return partialTmpl(w, data)
		//}

		return tmpl(w, data)
	}
}
