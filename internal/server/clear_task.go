package server

import (
	"errors"
	"fmt"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
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

		data := clearTaskData{XSRFToken: ctx.XSRFToken}
		taskID, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}

		task, err := client.Task(ctx, taskID)
		if err != nil {
			return err
		}

		if len(task.CaseItems) > 0 {
			lpa = &task.CaseItems[0]
		}

		data.Uid = lpa.UID
		data.Task = task

		if r.Method == http.MethodPost {
			err = client.ClearTask(ctx, taskID)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
			} else if err != nil {
				return err
			} else if lpa != nil && lpa.CaseType == "DIGITAL_LPA" {
				SetFlash(w, FlashNotification{Title: "Task completed"})
				return RedirectError(fmt.Sprintf("/lpa/%s", lpa.UID))
			} else {
				data.Success = true
			}
		}

		return tmpl(w, data)
	}
}
