package server

import (
	"log"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type LpaClient interface {
	DigitalLpa(ctx sirius.Context, uid string) (sirius.DigitalLpa, error)
}

type lpaData struct {
	Lpa sirius.DigitalLpa
}

func Lpa(client LpaClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		uid := r.FormValue("uid")
		ctx := getContext(r)

		lpa, err := client.DigitalLpa(ctx, uid)

		if err != nil {
			return err
		}

		data := lpaData{
			Lpa: lpa,
		}

		log.Print(data)

		return tmpl(w, data)
	}
}
