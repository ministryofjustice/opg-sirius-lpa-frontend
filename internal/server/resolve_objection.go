package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"net/http"
)

type ResolveObjectionClient interface {
	GetObjection(sirius.Context, string) (sirius.Objection, error)
	ResolveObjection(sirius.Context, string, string, sirius.ResolutionRequest) error
}

type resolveObjectionData struct {
	XSRFToken   string
	Success     bool
	Error        sirius.ValidationError
	ErrorsByCase map[string]sirius.ValidationError
	CaseUID      string
	ObjectionId string
	Objection   sirius.Objection
	Form        []formResolveObjection
}

type formResolveObjection struct {
	UID             string `form:"uid"`
	Resolution      string `form:"resolution"`
	ResolutionNotes string `form:"resolutionNotes"`
}

func ResolveObjection(client ResolveObjectionClient, formTmpl template.Template) Handler {

	return func(w http.ResponseWriter, r *http.Request) error {
		caseUID := r.PathValue("uid")
		objectionID := r.PathValue("id")

		ctx := getContext(r)

		obj, err := client.GetObjection(ctx, objectionID)
		if err != nil {
			return err
		}

		data := resolveObjectionData{
			XSRFToken:   ctx.XSRFToken,
			CaseUID:     caseUID,
			ObjectionId: objectionID,
			Objection:   obj,
			Form:        make([]formResolveObjection, len(obj.LpaUids)),
		}

		for i, uid := range obj.LpaUids {
			data.Form[i].UID = uid
		}

		Errors := make(map[string]sirius.ValidationError)

		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				return err
			}

			var validationErrors sirius.ValidationError

			for i, uid := range obj.LpaUids {
				resolution := r.PostForm.Get(fmt.Sprintf("resolution-%s", uid))
				notes := r.PostForm.Get(fmt.Sprintf("resolutionNotes-%s", uid))

				data.Form[i].Resolution = resolution
				data.Form[i].ResolutionNotes = notes

				err := client.ResolveObjection(ctx, objectionID, uid, sirius.ResolutionRequest{
					Resolution: resolution,
					Notes:      notes,
				})

				if ve, ok := err.(sirius.ValidationError); ok {
					Errors[uid] = ve
					validationErrors = ve
				} else if err != nil {
					return err
				}
			}

			if len(Errors) > 0 {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = validationErrors
				data.ErrorsByCase = Errors
				return formTmpl(w, data)
			}

			data.Success = true
			return RedirectError(fmt.Sprintf("/lpa/%s", caseUID))
		}

		return formTmpl(w, data)
	}
}
