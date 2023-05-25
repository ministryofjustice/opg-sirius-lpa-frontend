package server

import (
	"encoding/json"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type PostcodeLookupClient interface {
	PostcodeLookup(ctx sirius.Context, postcode string) ([]sirius.PostcodeLookupAddress, error)
}

func SearchPostcode(client PostcodeLookupClient) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		addresses, err := client.PostcodeLookup(ctx, r.FormValue("postcode"))
		if err != nil {
			return err
		}

		return json.NewEncoder(w).Encode(addresses)
	}
}
