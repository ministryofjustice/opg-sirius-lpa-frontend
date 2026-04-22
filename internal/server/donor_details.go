package server

import (
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type DonorDetailsClient interface {
	Person(sirius.Context, int) (sirius.Person, error)
}

type DonorDetailsData struct {
	Donor sirius.Person
}

func DonorDetails(client DonorDetailsClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if err := r.ParseForm(); err != nil {
			return err
		}

		donorID, err := strconv.Atoi(r.PathValue("donorId"))
		if err != nil {
			return err
		}

		ctx := getContext(r)

		donorDetails, err := client.Person(ctx, donorID)
		if err != nil {
			return err
		}

		data := DonorDetailsData{
			Donor: donorDetails,
		}

		return tmpl(w, data)
	}
}
