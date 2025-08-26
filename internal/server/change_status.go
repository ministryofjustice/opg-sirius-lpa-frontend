package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type ChangeStatusClient interface {
	Case(sirius.Context, int) (sirius.Case, error)
	EditCase(sirius.Context, int, sirius.CaseType, sirius.Case) error
	AvailableStatuses(sirius.Context, int, sirius.CaseType) ([]string, error)
}

type changeStatusData struct {
	XSRFToken string
	Entity    string
	Success   bool
	Error     sirius.ValidationError

	AvailableStatuses []string
	NewStatus         string
}

func ChangeStatus(client ChangeStatusClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		caseID, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}

		caseType, err := sirius.ParseCaseType(r.FormValue("case"))
		if err != nil {
			return err
		}

		ctx := getContext(r)

		caseitem, err := client.Case(ctx, caseID)
		if err != nil {
			return err
		}

		availableStatuses, err := client.AvailableStatuses(ctx, caseID, caseType)
		if err != nil {
			return err
		}

		data := changeStatusData{
			XSRFToken:         ctx.XSRFToken,
			Entity:            fmt.Sprintf("%s %s", caseitem.CaseType, caseitem.UID),
			AvailableStatuses: availableStatuses,
			NewStatus:         postFormString(r, "status"),
		}

		if r.Method == http.MethodPost {
			caseDetails := sirius.Case{
				Status: shared.ParseCaseStatusType(data.NewStatus),
			}

			err = client.EditCase(ctx, caseID, sirius.CaseType(caseitem.CaseType), caseDetails)

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
