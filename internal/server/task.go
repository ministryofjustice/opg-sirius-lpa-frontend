package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/sync/errgroup"
)

type TaskClient interface {
	CreateTask(ctx sirius.Context, caseID int, task sirius.TaskRequest) error
	TaskTypes(ctx sirius.Context) ([]string, error)
	Teams(ctx sirius.Context) ([]sirius.Team, error)
	Case(ctx sirius.Context, id int) (sirius.Case, error)
}

type taskData struct {
	XSRFToken string
	Entity    string
	Success   bool
	Error     sirius.ValidationError

	Case             sirius.Case
	TaskTypes        []string
	Teams            []sirius.Team
	Task             sirius.TaskRequest
	AssignTo         string
	AssigneeUserName string
}

func Task(client TaskClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		caseID, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}

		ctx := getContext(r)
		data := taskData{
			XSRFToken: ctx.XSRFToken,
		}

		group, groupCtx := errgroup.WithContext(ctx.Context)

		group.Go(func() error {
			taskTypes, err := client.TaskTypes(ctx.With(groupCtx))
			if err != nil {
				return err
			}

			data.TaskTypes = taskTypes
			return nil
		})

		group.Go(func() error {
			teams, err := client.Teams(ctx.With(groupCtx))
			if err != nil {
				return err
			}

			data.Teams = teams
			return nil
		})

		group.Go(func() error {
			caseitem, err := client.Case(ctx.With(groupCtx), caseID)
			if err != nil {
				return err
			}
			data.Entity = caseitem.Summary()
			data.Case = caseitem
			return nil
		})

		if err := group.Wait(); err != nil {
			return err
		}

		if r.Method == http.MethodPost {
			task := sirius.TaskRequest{
				Type:        postFormString(r, "taskType"),
				DueDate:     postFormDateString(r, "dueDate"),
				Name:        postFormString(r, "name"),
				Description: postFormString(r, "description"),
			}
			assignTo := postFormString(r, "assignTo")

			switch assignTo {
			case "user":
				parts := strings.SplitN(postFormString(r, "assigneeUser"), ":", 2)
				if len(parts) == 2 {
					assigneeID, _ := strconv.Atoi(parts[0])
					task.AssigneeID = assigneeID
					data.AssigneeUserName = parts[1]
				}
			case "team":
				assigneeID, _ := postFormInt(r, "assigneeTeam")
				task.AssigneeID = assigneeID
			}

			err = client.CreateTask(ctx, caseID, task)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Task = task
				data.AssignTo = assignTo
				data.Error = ve

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

				SetFlash(w, FlashNotification{
					Title: "Task created",
				})

				if data.Case.CaseType == "DIGITAL_LPA" {
					return RedirectError(fmt.Sprintf("/lpa/%s", data.Case.UID))
				}
			}
		}

		return tmpl(w, data)
	}
}
