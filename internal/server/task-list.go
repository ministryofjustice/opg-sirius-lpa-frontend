package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type GetTaskListClient interface {
	TasksForDonor(ctx sirius.Context, donorId int) ([]sirius.Task, error)
}

type taskListData struct {
	XSRFToken  string
	Error      sirius.ValidationError
	CasesTasks []sirius.Task
}

func GetTaskList(client GetTaskListClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if err := r.ParseForm(); err != nil {
			return err
		}
		donorId, err := strconv.Atoi(r.PathValue("donorId"))
		if err != nil {
			return err
		}
		ctx := getContext(r)

		donorTasks, err := client.TasksForDonor(ctx, donorId)
		if err != nil {
			fmt.Println("Error getting tasks for donor: ", err)
			return err
		}

		data := taskListData{
			XSRFToken:  ctx.XSRFToken,
			CasesTasks: donorTasks,
		}

		return tmpl(w, data)
	}
}
