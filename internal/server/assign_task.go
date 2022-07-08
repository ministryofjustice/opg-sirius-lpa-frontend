package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/sync/errgroup"
)

type AssignTaskClient interface {
	AssignTask(ctx sirius.Context, assigneeID, taskID int) error
	Task(ctx sirius.Context, id int) (sirius.Task, error)
	Teams(ctx sirius.Context) ([]sirius.Team, error)
}

type assignTaskData struct {
	XSRFToken string
	Entity    string
	Success   bool
	Error     sirius.ValidationError

	Teams            []sirius.Team
	AssignTo         string
	AssigneeUserName string
}

func AssignTask(client AssignTaskClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		taskID, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}

		ctx := getContext(r)
		data := assignTaskData{XSRFToken: ctx.XSRFToken}

		group, groupCtx := errgroup.WithContext(ctx.Context)

		group.Go(func() error {
			teams, err := client.Teams(ctx.With(groupCtx))
			if err != nil {
				return err
			}

			data.Teams = teams
			return nil
		})

		group.Go(func() error {
			task, err := client.Task(ctx.With(groupCtx), taskID)
			if err != nil {
				return err
			}

			data.Entity = task.Summary()
			return nil
		})

		if err := group.Wait(); err != nil {
			return err
		}

		if r.Method == http.MethodPost {
			var assigneeID int
			assignTo := postFormString(r, "assignTo")

			switch assignTo {
			case "user":
				parts := strings.SplitN(postFormString(r, "assigneeUser"), ":", 2)
				if len(parts) == 2 {
					assigneeID, _ = strconv.Atoi(parts[0])
					data.AssigneeUserName = parts[1]
				}
			case "team":
				assigneeID, _ = postFormInt(r, "assigneeTeam")
			}

			err := client.AssignTask(ctx, assigneeID, taskID)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
				data.AssignTo = assignTo

				switch data.AssignTo {
				case "user":
					data.Error.Field["assigneeUser"] = data.Error.Field["assigneeId"]
				case "team":
					data.Error.Field["assigneeTeam"] = data.Error.Field["assigneeId"]
				default:
					data.Error.Field["assignTo"] = map[string]string{"": "Assignee not set"}
				}
				delete(data.Error.Field, "assigneeId")
			} else if err != nil {
				return err
			} else {
				data.Success = true
			}
		}

		return tmpl(w, data)
	}
}
