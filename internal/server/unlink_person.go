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
	CaseUids  string

	Person sirius.Person
}

func UnlinkPerson(client UnlinkPersonClient, tmpl template.Template, partialTmpl template.Template) Handler {
	return func(pageVars PageVars, w http.ResponseWriter, r *http.Request) error {
		parentID, err := strToIntOrStatusError(r.FormValue("id"))
		if err != nil {
			return err
		}

		ctx := getContext(r)
		data := unlinkPersonData{
			XSRFToken: ctx.XSRFToken,
			CaseUids:  buildUIDQueryString(r.Form["uid[]"]),
		}

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

		if r.Header.Get("HX-Request") == "true" {
			return partialTmpl(w, data)
		}

		return tmpl(w, data)
	}
}
