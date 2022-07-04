package server

import (
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
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

		if r.Method == http.MethodPost {
			var childId int
			id := postFormString(r, "child-id")

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

		data.Person, err = client.Person(ctx, parentID)
		if err != nil {
			return err
		}

		return tmpl(w, data)
	}
}
