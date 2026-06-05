package server

import (
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type ViewDocumentClient interface {
	DocumentByUUID(ctx sirius.Context, uuid string) (sirius.Document, error)
	GetUserDetails(sirius.Context) (sirius.User, error)
}

type viewDocumentData struct {
	XSRFToken      string
	Document       sirius.Document
	IsSysAdminUser bool
	Pane           int
}

func ViewDocument(client ViewDocumentClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		uuid := r.PathValue("uuid")
		ctx := getContext(r)

		documentData, err := client.DocumentByUUID(ctx, uuid)
		if err != nil {
			return err
		}

		user, err := client.GetUserDetails(ctx)
		if err != nil {
			return err
		}
		isSysAdminUser := user.HasRole("System Admin")

		// Extract pane parameter from query string if present
		pane := 1 // Default to pane 1
		if paneStr := r.URL.Query().Get("pane"); paneStr != "" {
			if paneNum, err := strconv.Atoi(paneStr); err == nil && paneNum == 2 {
				pane = 2
			}
		}

		data := viewDocumentData{
			XSRFToken:      ctx.XSRFToken,
			Document:       documentData,
			IsSysAdminUser: isSysAdminUser,
			Pane:           pane,
		}

		return tmpl(w, data)
	}
}
