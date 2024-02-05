package server

import (
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/sync/errgroup"
)

type DeleteRelationshipClient interface {
	PersonReferences(sirius.Context, int) ([]sirius.PersonReference, error)
	DeletePersonReference(sirius.Context, int) error
	Person(sirius.Context, int) (sirius.Person, error)
}

type deleteRelationshipData struct {
	XSRFToken string
	Entity    string
	Success   bool
	Error     sirius.ValidationError

	PersonReferences []sirius.PersonReference
}

func DeleteRelationship(client DeleteRelationshipClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		personID, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}

		ctx := getContext(r)
		data := deleteRelationshipData{XSRFToken: ctx.XSRFToken}

		if r.Method == http.MethodPost {
			referenceID, err := postFormInt(r, "reference-id")
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = sirius.ValidationError{
					Detail: "Select a relationship to delete",
				}
			} else {
				err = client.DeletePersonReference(ctx, referenceID)
				if err != nil {
					return err
				}
				data.Success = true
			}
		}

		group, groupCtx := errgroup.WithContext(ctx.Context)

		group.Go(func() error {
			person, err := client.Person(ctx.With(groupCtx), personID)
			if err != nil {
				return err
			}

			data.Entity = person.Summary()
			return nil
		})

		group.Go(func() error {
			references, err := client.PersonReferences(ctx.With(groupCtx), personID)
			if err != nil {
				return err
			}

			data.PersonReferences = references
			return nil
		})

		if err := group.Wait(); err != nil {
			return err
		}

		return tmpl(w, data)
	}
}
