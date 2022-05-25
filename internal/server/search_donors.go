package server

import (
	"encoding/json"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type SearchDonorsClient interface {
	SearchDonors(ctx sirius.Context, term string) ([]sirius.Person, error)
}

func SearchDonors(client SearchDonorsClient) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		persons, err := client.SearchDonors(ctx, r.FormValue("q"))
		if err != nil {
			return err
		}

		return json.NewEncoder(w).Encode(persons)
	}
}
