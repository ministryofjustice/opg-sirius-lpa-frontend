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
	CasesByDonor(ctx sirius.Context, id int) ([]sirius.Case, error)
}

type ActionPanelData struct {
	XSRFToken     string
	NoteTypes     []string
	Entity        string
	IsDigitalLpa  bool
	CaseUID       string
	Success       bool
	Error         sirius.ValidationError
	DonorID       int
	SelectedCases []sirius.Case
	CaseUids      string
	CaseType      string
}

func ActionPanel(client ActionPanelClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if err := r.ParseForm(); err != nil {
			return err
		}

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

		caseUIDs := r.Form["uid[]"]
		data.CaseUids = buildUIDQueryString(caseUIDs)

		entityType, err := sirius.ParseEntityType(r.FormValue("entity"))
		if err != nil {
			return err
		}
		data.CaseType = string(entityType)

		group, groupCtx := errgroup.WithContext(ctx.Context)

		group.Go(func() error {
			noteTypes, err := client.NoteTypes(ctx.With(groupCtx))
			if err != nil {
				return err
			}

			data.NoteTypes = noteTypes
			return nil
		})

		group.Go(func() error {
			if data.DonorID > 0 {
				cases, err := client.CasesByDonor(ctx.With(groupCtx), data.DonorID)
				if err != nil {
					return err
				}

				// Filter cases by uid[] parameter if provided
				if len(caseUIDs) > 0 {
					casesByUID := make(map[string]sirius.Case, len(cases))
					for _, c := range cases {
						casesByUID[c.UID] = c
					}

					var filtered []sirius.Case
					for _, uid := range caseUIDs {
						if c, ok := casesByUID[uid]; ok {
							filtered = append(filtered, c)
						}
					}
					data.SelectedCases = filtered
				} else {
					data.SelectedCases = cases
				}
			}
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
