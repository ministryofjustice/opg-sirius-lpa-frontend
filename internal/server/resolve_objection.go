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
	Error       sirius.ValidationError
	CaseUID     string
	ObjectionId string
	Objection   sirius.Objection
	Form        formResolveObjection
}

type formResolveObjectionItem struct {
	UID            string `form:"uid"`
	Resolution     string `form:"resolution"`
	ResolutionNote string `form:"resolutionNotes"`
}

type formResolveObjection struct {
	Resolutions []formResolveObjectionItem `form:"resolutions"`
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

		initialResolutions := make([]formResolveObjectionItem, len(obj.LpaUids))
		for i, uid := range obj.LpaUids {
			initialResolutions[i] = formResolveObjectionItem{
				UID: uid,
			}
		}

		data := resolveObjectionData{
			XSRFToken:   ctx.XSRFToken,
			CaseUID:     caseUID,
			ObjectionId: objectionID,
			Objection:   obj,
			Form: formResolveObjection{
				Resolutions: initialResolutions,
			},
		}

		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				return err
			}

			var formData formResolveObjection

			if err := decoder.Decode(&formData, r.PostForm); err != nil {
				return err
			}

			var validationErrors sirius.ValidationError
			var hasValidationError bool

			for _, item := range formData.Resolutions {
				err := client.ResolveObjection(ctx, objectionID, item.UID, sirius.ResolutionRequest{
					Resolution: item.Resolution,
					Notes:      item.ResolutionNote,
				})

				if ve, ok := err.(sirius.ValidationError); ok {
					hasValidationError = true
					if validationErrors.Field == nil {
						validationErrors = ve
					}
				} else if err != nil {
					return err
				}
			}

			if hasValidationError {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = validationErrors
				data.Form = formData
				return formTmpl(w, data)
			}

			data.Success = true
			return RedirectError(fmt.Sprintf("/lpa/%s", caseUID))
		}

		return formTmpl(w, data)
	}
}
