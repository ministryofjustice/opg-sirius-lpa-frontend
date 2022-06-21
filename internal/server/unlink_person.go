package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"net/http"
	"strconv"
)

type UnlinkPersonClient interface {
	Person(sirius.Context, int) (sirius.Person, error)
	UnlinkPerson(sirius.Context, int, []int) error
}

type unlinkPersonData struct {
	XSRFToken string
	Success   bool
	Error     sirius.ValidationError

	Entity sirius.Person
}

func UnlinkPerson(client UnlinkPersonClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		parentID, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}

		ctx := getContext(r)
		data := unlinkPersonData{XSRFToken: ctx.XSRFToken}

		data.Entity, err = client.Person(ctx, parentID)
		if err != nil {
			return err
		}

		var childIds []int

		if r.Method == http.MethodPost {
			for i := range data.Entity.Children {
				id := r.FormValue(fmt.Sprintf("child-%d", i))
				if id != "" {
					id, err := strconv.Atoi(id)
					if err != nil {
						return err
					}
					childIds = append(childIds, id)
				}
			}

			if len(childIds) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = sirius.ValidationError{
					Field: sirius.FieldErrors{
						"child": {"reason": "Please select the record to be unlinked"},
					},
				}
			} else {
				err = client.UnlinkPerson(ctx, parentID, childIds)
				if err != nil {
					return err
				} else {
					data.Success = true
				}
			}
		}

		return tmpl(w, data)
	}
}
