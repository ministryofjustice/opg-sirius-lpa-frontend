package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type RelationshipClient interface {
	CreatePersonReference(ctx sirius.Context, personID int, referencedUID, reason string) error
	Person(ctx sirius.Context, id int) (sirius.Person, error)
}

type relationshipData struct {
	XSRFToken string
	Entity    string
	Success   bool
	Error     sirius.ValidationError

	SearchUID  string
	SearchName string
	Reason     string
}

func Relationship(client RelationshipClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		personID, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}

		ctx := getContext(r)

		person, err := client.Person(ctx, personID)
		if err != nil {
			return err
		}

		data := relationshipData{
			XSRFToken: ctx.XSRFToken,
			Entity:    fmt.Sprintf("%s %s", person.Firstname, person.Surname),
		}

		if r.Method == http.MethodPost {
			data.Reason = r.FormValue("reason")

			parts := strings.SplitN(r.FormValue("search"), ":", 2)
			if len(parts) == 2 {
				data.SearchUID = parts[0]
				data.SearchName = parts[1]
			}

			err = client.CreatePersonReference(ctx, personID, data.SearchUID, data.Reason)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
			} else if err != nil {
				return err
			} else {
				data.Success = true
			}
		}

		return tmpl(w, data)
	}
}
