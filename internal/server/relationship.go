package server

import (
	"fmt"
	"net/http"
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

	ActionPanelData actionPanelData

	SearchUID  string
	SearchName string
	Reason     string
}

type actionPanelData struct {
	DonorID    int
	CaseUIDs   string
	EntityType string
}

func Relationship(client RelationshipClient, tmpl template.Template, partialTmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		personID, err := strToIntOrStatusError(r.FormValue("id"))
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
			ActionPanelData: actionPanelData{
				DonorID: personID,
			},
		}

		data.ActionPanelData.CaseUIDs = buildUIDQueryString(r.Form["uid[]"])

		if entityType, err := sirius.ParseEntityType(r.FormValue("entity")); err == nil {
			data.ActionPanelData.EntityType = string(entityType)
		}

		if r.Method == http.MethodPost {
			var (
				reason     = postFormString(r, "reason")
				searchUID  string
				searchName string
			)

			parts := strings.SplitN(postFormString(r, "search"), ":", 2)
			if len(parts) == 2 {
				searchUID = parts[0]
				searchName = parts[1]
			}

			err = client.CreatePersonReference(ctx, personID, searchUID, reason)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
				data.Reason = reason
				data.SearchUID = searchUID
				data.SearchName = searchName
			} else if err != nil {
				return err
			} else {
				data.Success = true
			}
		}
		if r.Header.Get("HX-Request") == "true" && partialTmpl != nil {
			return partialTmpl(w, data)
		}

		return tmpl(w, data)
	}
}
