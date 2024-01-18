package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/sync/errgroup"
)

type AssignTaskClient interface {
	AssignTasks(ctx sirius.Context, assigneeID int, taskIDs []int) error
	Task(ctx sirius.Context, id int) (sirius.Task, error)
	Teams(ctx sirius.Context) ([]sirius.Team, error)
	GetUserDetails(ctx sirius.Context) (sirius.User, error)
}

type assignTaskData struct {
	XSRFToken string
	Entities  []string
	Success   bool
	Error     sirius.ValidationError

	Teams            []sirius.Team
	AssignTo         string
	AssigneeUserName string
}

func AssignTask(client AssignTaskClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if err := r.ParseForm(); err != nil {
			return err
		}

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

		var tasksMu sync.Mutex
		var lpa *sirius.Case

		for _, taskID := range taskIDs {
			taskID := taskID

			group.Go(func() error {
				task, err := client.Task(ctx.With(groupCtx), taskID)
				if err != nil {
					return err
				}

				if lpa == nil && len(task.CaseItems) > 0 {
					lpa = &task.CaseItems[0]
				}

				tasksMu.Lock()
				data.Entities = append(data.Entities, task.Summary())
				tasksMu.Unlock()
				return nil
			})
		}

		if err := group.Wait(); err != nil {
			return err
		}

		if r.Method == http.MethodPost {
			var assigneeID int
			assignTo := postFormString(r, "assignTo")

			switch assignTo {
			case "me":
				user, err := client.GetUserDetails(ctx)
				if (err != nil){
					return err
				} else {
					assigneeID = user.ID
					data.AssigneeUserName = user.DisplayName
				}
			case "user":
				parts := strings.SplitN(postFormString(r, "assigneeUser"), ":", 2)
				if len(parts) == 2 {
					assigneeID, _ = strconv.Atoi(parts[0])
					data.AssigneeUserName = parts[1]
				}
			case "team":
				assigneeID, _ = postFormInt(r, "assigneeTeam")
			}

			err := client.AssignTasks(ctx, assigneeID, taskIDs)

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
			} else if lpa != nil && lpa.CaseType == "DIGITAL_LPA" {
				// redirect
				SetFlash(w, FlashNotification{
					Title: "Task assigned",
				})
				return RedirectError(fmt.Sprintf("/lpa/%s", lpa.UID))
			} else {
				data.Success = true
			}
		}

		return tmpl(w, data)
	}
}
