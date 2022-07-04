package server

import (
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type LinkPersonClient interface {
	LinkPeople(sirius.Context, int, int) error
	Person(sirius.Context, int) (sirius.Person, error)
	PersonByUid(sirius.Context, string) (sirius.Person, error)
}

type linkPersonData struct {
	XSRFToken        string
	Entity           sirius.Person
	OtherPerson      sirius.Person
	PrimaryId        int
	CanChangePrimary bool
	Error            sirius.ValidationError
	Success          bool
}

func LinkPerson(client LinkPersonClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		person1ID, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}

		ctx := getContext(r)
		data := linkPersonData{XSRFToken: ctx.XSRFToken}

		data.Entity, err = client.Person(ctx, person1ID)
		if err != nil {
			return err
		}

		if r.Method == http.MethodPost {
			data.OtherPerson, err = client.PersonByUid(ctx, postFormString(r, "uid"))
			if ve, ok := err.(sirius.StatusError); ok && ve.Code == http.StatusNotFound {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = sirius.ValidationError{
					Field: sirius.FieldErrors{
						"uid": map[string]string{
							"notFound": "A record matching the supplied uId cannot be found.",
						},
					},
				}

				return tmpl(w, data)
			} else if err != nil {
				return err
			}

			if len(data.Entity.Children) == 0 && len(data.OtherPerson.Children) > 0 {
				data.PrimaryId = data.OtherPerson.ID
				data.CanChangePrimary = false
			} else if len(data.Entity.Children) > 0 && len(data.OtherPerson.Children) == 0 {
				data.PrimaryId = data.Entity.ID
				data.CanChangePrimary = false
			} else {
				data.CanChangePrimary = true
			}

			if postFormString(r, "primary-id") != "" {
				if data.CanChangePrimary {
					data.PrimaryId, err = postFormInt(r, "primary-id")
					if err != nil {
						return err
					}
				}

				if data.PrimaryId == data.Entity.ID {
					err = client.LinkPeople(ctx, data.Entity.ID, data.OtherPerson.ID)
				} else if data.PrimaryId == data.OtherPerson.ID {
					err = client.LinkPeople(ctx, data.OtherPerson.ID, data.Entity.ID)
				} else {
					return tmpl(w, data)
				}

				if ve, ok := err.(sirius.ValidationError); ok {
					w.WriteHeader(http.StatusBadRequest)
					data.Error = ve
				} else if err != nil {
					return err
				} else {
					data.Success = true
				}
			}
		}

		return tmpl(w, data)
	}
}
