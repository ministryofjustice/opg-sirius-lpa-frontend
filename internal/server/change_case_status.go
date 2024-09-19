package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type ChangeCaseStatusClient interface {
	Case(sirius.Context, int) (sirius.Case, error)
	EditCase(sirius.Context, int, sirius.CaseType, sirius.Case) error
}

type changeCaseStatusData struct {
	XSRFToken string
	Entity    string
	Success   bool
	Error     sirius.ValidationError

	NewStatus string
}

func ChangeCaseStatus(client ChangeCaseStatusClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		caseID, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}

		ctx := getContext(r)

		caseitem, err := client.Case(ctx, caseID)
		if err != nil {
			return err
		}

		data := changeCaseStatusData{
			XSRFToken: ctx.XSRFToken,
			Entity:    fmt.Sprintf("%s %s", caseitem.CaseType, caseitem.UID),
			NewStatus: postFormString(r, "status"),
		}

		if r.Method == http.MethodPost {
			caseDetails := sirius.Case{
				Status: data.NewStatus,
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
