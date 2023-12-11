package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"golang.org/x/sync/errgroup"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/ministryofjustice/opg-go-common/template"
)

type GetDocumentsClient interface {
	DigitalLpa(ctx sirius.Context, uid string) (sirius.DigitalLpa, error)
	Documents(ctx sirius.Context, caseType sirius.CaseType, caseId int, docTypes []string, notDocTypes []string) ([]sirius.Document, error)
	TasksForCase(ctx sirius.Context, id int) ([]sirius.Task, error)
}

type getDocumentsData struct {
	XSRFToken string

	Lpa          sirius.DigitalLpa
	TaskList     []sirius.Task
	Documents    []sirius.Document
	FlashMessage FlashNotification
}

func GetDocuments(client GetDocumentsClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		var err error

		uid := chi.URLParam(r, "uid")
		ctx := getContext(r)
		group, _ := errgroup.WithContext(ctx.Context)

		data := getDocumentsData{
			XSRFToken: ctx.XSRFToken,
		}

		data.Lpa, err = client.DigitalLpa(ctx, uid)
		if err != nil {
			return err
		}

		group.Go(func() error {
			data.Documents, err = client.Documents(
				ctx,
				"lpa",
				data.Lpa.ID,
				[]string{}, []string{sirius.TypeDraft, sirius.TypePreview})

			if err != nil {
				return err
			}

			return nil
		})

		group.Go(func() error {
			data.TaskList, err = client.TasksForCase(ctx, data.Lpa.ID)

			if err != nil {
				return err
			}

			return nil
		})

		if err = group.Wait(); err != nil {
			return err
		}

		data.FlashMessage, _ = GetFlash(w, r)

		return tmpl(w, data)
	}
}
