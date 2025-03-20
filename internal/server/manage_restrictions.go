package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/form/v4"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/sync/errgroup"
	"net/http"
)

type ManageRestrictionsClient interface {
	CaseSummary(sirius.Context, string) (sirius.CaseSummary, error)
	ClearTask(sirius.Context, int) error
	UpdateSeveranceStatus(sirius.Context, string, sirius.SeveranceStatusData) error
}

type manageRestrictionsData struct {
	XSRFToken       string
	Error           sirius.ValidationError
	CaseUID         string
	CaseSummary     sirius.CaseSummary
	SeveranceAction string
	Success         bool
}

func ManageRestrictions(client ManageRestrictionsClient, tmpl template.Template) Handler {
	if decoder == nil {
		decoder = form.NewDecoder()
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		caseUID := chi.URLParam(r, "uid")
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

		data := manageRestrictionsData{
			CaseSummary:     cs,
			SeveranceAction: postFormString(r, "severanceAction"),
			XSRFToken:       ctx.XSRFToken,
			Error:           sirius.ValidationError{Field: sirius.FieldErrors{}},
			CaseUID:         caseUID,
		}

		if r.Method == http.MethodPost {
			taskList := data.CaseSummary.TaskList
			var taskID int
			for _, task := range taskList {
				if task.Name == "Review restrictions and conditions" && task.Status != "Completed" {
					taskID = task.ID
				}
			}

			switch data.SeveranceAction {
			case "severance-application-not-required":
				err := client.ClearTask(ctx, taskID)

				if ve, ok := err.(sirius.ValidationError); ok {
					w.WriteHeader(http.StatusBadRequest)
					data.Error = ve
				} else if err != nil {
					return err
				} else {
					data.Success = true

					SetFlash(w, FlashNotification{
						Title: "Changes confirmed",
					})

					return RedirectError(fmt.Sprintf("/lpa/%s/lpa-details", caseUID))
				}

			case "severance-application-required":
				err := client.UpdateSeveranceStatus(ctx, caseUID, sirius.SeveranceStatusData{
					SeveranceStatus: "REQUIRED",
				})

				if ve, ok := err.(sirius.ValidationError); ok {
					w.WriteHeader(http.StatusBadRequest)
					data.Error = ve
				} else if err != nil {
					return err
				} else {
					data.Success = true

					SetFlash(w, FlashNotification{
						Title: "Changes confirmed",
					})

					return RedirectError(fmt.Sprintf("/lpa/%s/lpa-details", caseUID))
				}

			default:
				w.WriteHeader(http.StatusBadRequest)
				data.Error.Field["severanceAction"] = map[string]string{
					"reason": "Please select an option",
				}
			}
		}
		return tmpl(w, data)
	}
}
