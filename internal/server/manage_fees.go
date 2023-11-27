package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type ManageFeesClient interface {
	Case(sirius.Context, int) (sirius.Case, error)
}

type manageFeesData struct {
	Case      sirius.Case
	ReturnUrl string
}

func ManageFees(client ManageFeesClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		caseID, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}

		data := manageFeesData{}

		data.Case, err = client.Case(getContext(r), caseID)
		if err != nil {
			return err
		}

		if data.Case.CaseType == "DIGITAL_LPA" {
			data.ReturnUrl = fmt.Sprintf("/lpa/%s/payments", data.Case.UID)
		} else {
			data.ReturnUrl = fmt.Sprintf("/payments/%d", caseID)
		}

		return tmpl(w, data)
	}
}
