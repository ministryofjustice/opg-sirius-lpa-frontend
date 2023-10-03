package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"net/http"
)

type LpaClient interface {
	DigitalLpa(ctx sirius.Context, uid string) (sirius.DigitalLpa, error)
	TasksForCase(ctx sirius.Context, id int) ([]sirius.Task, error)
}

type lpaData struct {
	Lpa      sirius.DigitalLpa
	TaskList []sirius.Task
}

func Lpa(client LpaClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		uid := chi.URLParam(r, "uid")
		ctx := getContext(r)

		lpa, err := client.DigitalLpa(ctx, uid)

		if err != nil {
			return err
		}

		tasks, err := client.TasksForCase(ctx, lpa.ID)

		if err != nil {
			return err
		}

		data := lpaData{
			Lpa:      lpa,
			TaskList: tasks,
		}

		return tmpl(w, data)
	}
}
