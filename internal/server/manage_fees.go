package server

import (
    "fmt"
    "golang.org/x/sync/errgroup"
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

        ctx := getContext(r)
        group, groupCtx := errgroup.WithContext(ctx.Context)

        data := manageFeesData{}

        group.Go(func() error {
            data.Case, err = client.Case(ctx.With(groupCtx), caseID)
            if err != nil {
                return err
            }

            return nil
        })

        if err := group.Wait(); err != nil {
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
