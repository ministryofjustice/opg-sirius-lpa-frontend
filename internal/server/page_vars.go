package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type PageVars struct {
	UserPermissions sirius.Permissions
}

func permissionMiddleware(client Client, r *http.Request) (*PageVars, error) {
	ctx := getContext(r)

	userPermissions, err := client.GetUserPermissions(ctx)
	if err != nil {
		return nil, err
	}

	vars := PageVars{
		UserPermissions: userPermissions,
	}

	return &vars, nil
}
