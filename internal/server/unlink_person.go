package server

import (
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"net/http"
	"strconv"
)

type UnlinkPersonClient interface {
	Person(sirius.Context, int) (sirius.Person, error)
	UnlinkPerson(sirius.Context, int, int) error
}

type unlinkPersonData struct {
	XSRFToken string
	Success   bool
	Error     sirius.ValidationError

	Person sirius.Person
}

func UnlinkPerson(client UnlinkPersonClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		parentID, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}

		ctx := getContext(r)
		data := unlinkPersonData{XSRFToken: ctx.XSRFToken}

		data.Person, err = client.Person(ctx, parentID)
		if err != nil {
			return err
		}

		var childId int

		if r.Method == http.MethodPost {
			id := r.FormValue("child-id")

			if id == "" {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = sirius.ValidationError{
					Field: sirius.FieldErrors{
						"child": {"reason": "Please select the record to be unlinked"},
					},
				}
			} else {
				childId, err = strconv.Atoi(id)
				if err != nil {
					return err
				}
				err = client.UnlinkPerson(ctx, parentID, childId)
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
