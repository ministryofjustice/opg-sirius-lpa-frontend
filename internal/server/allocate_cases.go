package server

import (
	"errors"
	"fmt"
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
	XSRFToken string
	Entities  []string
	Success   bool
	Error     sirius.ValidationError

	Teams            []sirius.Team
	AssignTo         string
	AssigneeUserName string
}

func AllocateCases(client AllocateCasesClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if err := r.ParseForm(); err != nil {
			return err
		}

		var caseIDs []int
		for _, id := range r.Form["id"] {
			caseID, err := strconv.Atoi(id)
			if err != nil {
				return err
			}
			caseIDs = append(caseIDs, caseID)
		}

		if len(caseIDs) == 0 {
			return errors.New("no cases selected")
		}

		ctx := getContext(r)
		data := allocateCasesData{XSRFToken: ctx.XSRFToken}

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
				data.Entities = append(data.Entities, fmt.Sprintf("%s %s", caseitem.CaseType, caseitem.UID))
				allocations = append(allocations, sirius.CaseAllocation{
					ID:       caseID,
					CaseType: caseitem.CaseType,
				})
				casesMu.Unlock()
				return nil
			})
		}

		if err := group.Wait(); err != nil {
			return err
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

		return tmpl(w, data)
	}
}
