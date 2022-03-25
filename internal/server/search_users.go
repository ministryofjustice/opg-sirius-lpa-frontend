package server

import (
	"encoding/json"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type SearchUsersClient interface {
	SearchUsers(ctx sirius.Context, term string) ([]sirius.User, error)
}

func SearchUsers(client SearchUsersClient) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		users, err := client.SearchUsers(ctx, r.FormValue("q"))
		if err != nil {
			return err
		}

		return json.NewEncoder(w).Encode(users)
	}
}
