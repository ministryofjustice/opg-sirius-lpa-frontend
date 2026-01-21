package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type DeleteDocumentClient interface {
	DocumentByUUID(ctx sirius.Context, uuid string) (sirius.Document, error)
	GetUserDetails(sirius.Context) (sirius.User, error)
	DeleteDocument(ctx sirius.Context, uuid string) error
}

type deleteDocumentData struct {
	XSRFToken      string
	Document       sirius.Document
	IsSysAdminUser bool
	DonorId        int
}

func DeleteDocument(client DeleteDocumentClient, tmpl template.Template) Handler {
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
		donorId := documentData.CaseItems[0].Donor.ID

		data := deleteDocumentData{
			XSRFToken:      ctx.XSRFToken,
			Document:       documentData,
			IsSysAdminUser: isSysAdminUser,
			DonorId:        donorId,
		}

		if r.Method == http.MethodPost {
			err = client.DeleteDocument(ctx, uuid)
			if err != nil {
				return err
			}

			return RedirectError(fmt.Sprintf("/donor/%d/documents?success=true&documentFriendlyName=%s&documentCreatedTime=%s", donorId, documentData.FriendlyDescription, documentData.CreatedDate))
		}

		return tmpl(w, data)
	}
}
