package server

import (
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/sync/errgroup"
)

type AllocateCasesClient interface {
	AllocateCases(ctx sirius.Context, assigneeID int, allocations []sirius.CaseAllocation) error
	Case(ctx sirius.Context, id int) (sirius.Case, error)
	Teams(ctx sirius.Context) ([]sirius.Team, error)
}

type allocateCasesData struct {
	XSRFToken        string
	Entities         []string
	Success          bool
	Error            sirius.ValidationError
	CaseID           int
	CaseIDs          []int
	CaseUIDs         string
	DonorID          int
	EntityType       string
	CaseUids         string
	Teams            []sirius.Team
	AssignTo         string
	AssigneeUserName string
}

func AllocateCases(client AllocateCasesClient, tmpl template.Template, partialTmpl template.Template) Handler {
	return func(pageVars PageVars, w http.ResponseWriter, r *http.Request) error {
		if err := r.ParseForm(); err != nil {
			return err
		}

		errNotFound := sirius.StatusError{Code: http.StatusNotFound}

		var caseIDs []int
		for _, id := range r.Form["id"] {
			if id == "" {
				return errNotFound
			}

			caseID, err := strconv.Atoi(id)
			if err != nil {
				return errNotFound
			}
			caseIDs = append(caseIDs, caseID)
		}

		if len(caseIDs) == 0 {
			return errNotFound
		}

		ctx := getContext(r)
		data := allocateCasesData{XSRFToken: ctx.XSRFToken, CaseID: caseIDs[0], CaseIDs: caseIDs}

		group, groupCtx := errgroup.WithContext(ctx.Context)

		group.Go(func() error {
			teams, err := client.Teams(ctx.With(groupCtx))
			if err != nil {
				return err
			}

			data.Teams = teams
			return nil
		})

		var casesMu sync.Mutex
		var allocations []sirius.CaseAllocation

		for _, caseID := range caseIDs {
			caseID := caseID

			group.Go(func() error {
				caseitem, err := client.Case(ctx.With(groupCtx), caseID)
				if err != nil {
					return err
				}

				casesMu.Lock()
				data.Entities = append(data.Entities, caseitem.Summary())
				allocations = append(allocations, sirius.CaseAllocation{
					ID:       caseID,
					CaseType: caseitem.CaseType,
				})
				if caseitem.Donor != nil {
					data.DonorID = caseitem.Donor.ID
				}
				casesMu.Unlock()
				return nil
			})
		}

		if err := group.Wait(); err != nil {
			return err
		}

		data.CaseUIDs = buildUIDQueryString(r.Form["uid[]"])

		if entityType, err := sirius.ParseEntityType(r.FormValue("entity")); err == nil {
			data.EntityType = string(entityType)
		}

		if r.Method == http.MethodPost {
			var assigneeID int
			assignTo := postFormString(r, "assignTo")

			switch assignTo {
			case "user":
				parts := strings.SplitN(postFormString(r, "assigneeUser"), ":", 2)
				if len(parts) == 2 {
					assigneeID, _ = strconv.Atoi(parts[0])
					data.AssigneeUserName = parts[1]
				}
			case "team":
				assigneeID, _ = postFormInt(r, "assigneeTeam")
			}

			err := client.AllocateCases(ctx, assigneeID, allocations)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
				data.AssignTo = assignTo

				switch data.AssignTo {
				case "user":
					data.Error.Field["assigneeUser"] = data.Error.Field["assigneeId"]
				case "team":
					data.Error.Field["assigneeTeam"] = data.Error.Field["assigneeId"]
				default:
					data.Error.Field["assignTo"] = map[string]string{"": "Assignee not set"}
				}
				delete(data.Error.Field, "assigneeId")
			} else if err != nil {
				return err
			} else {
				data.Success = true
			}
		}
		if r.Header.Get("HX-Request") == "true" {
			return partialTmpl(w, data)
		}

		return tmpl(w, data)
	}
}
