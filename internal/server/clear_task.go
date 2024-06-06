package server

import (
	"errors"
	"fmt"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
)

type ClearTaskClient interface {
	ClearTask(ctx sirius.Context, taskID int) error
	Task(ctx sirius.Context, id int) (sirius.Task, error)
}

type clearTaskData struct {
	XSRFToken string
	Entities  []string
	Uid       string
	Task      sirius.Task
	Success   bool
	Error     sirius.ValidationError
}

func ClearTask(client ClearTaskClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if err := r.ParseForm(); err != nil {
			return err
		}

		var lpa *sirius.Case
		var taskIDs []int

		for _, id := range r.Form["id"] {
			taskID, err := strconv.Atoi(id)
			if err != nil {
				return err
			}
			taskIDs = append(taskIDs, taskID)
		}
		if len(taskIDs) == 0 {
			return errors.New("no tasks selected")
		}

		ctx := getContext(r)
		group, groupCtx := errgroup.WithContext(ctx.Context)

		data := clearTaskData{XSRFToken: ctx.XSRFToken}
		taskID, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}

		group.Go(func() error {
			task, err := client.Task(ctx.With(groupCtx), taskID)
			if err != nil {
				return err
			}

			if lpa == nil && len(task.CaseItems) > 0 {
				lpa = &task.CaseItems[0]
			}

			data.Uid = lpa.UID
			data.Task = task

			return nil
		})

		if err := group.Wait(); err != nil {
			return err
		}

		if r.Method == http.MethodPost {
			err = client.ClearTask(ctx, taskID)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
			} else if err != nil {
				return err
			}

			data.Success = true
			SetFlash(w, FlashNotification{Title: "Task completed"})
			return RedirectError(fmt.Sprintf("/lpa/%s", lpa.UID))
		}

		return tmpl(w, data)
	}
}
