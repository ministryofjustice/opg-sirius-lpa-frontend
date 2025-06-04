package server

import (
	"fmt"
	"github.com/go-playground/form/v4"
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
	LpaUids     []string
	Form        formResolveObjection
}

type formResolveObjection struct {
	Resolution      string `form:"resolution"`
	ResolutionNotes string `form:"resolutionNotes"`
}

func ResolveObjection(client ResolveObjectionClient, formTmpl template.Template) Handler {
	if decoder == nil {
		decoder = form.NewDecoder()
	}

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
			LpaUids:     obj.LpaUids,
		}

		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				return err
			}

			uids := r.PostForm["caseUid"]

			results := make([]sirius.ResolutionRequest, len(data.LpaUids))

			var validationErrors sirius.ValidationError
			var hasValidationError bool

			for i := range uids {

				resolution := r.PostForm.Get(fmt.Sprintf("resolution-%d", i))
				notes := r.PostForm.Get(fmt.Sprintf("resolutionNotes-%d", i))

				results[i] = sirius.ResolutionRequest{
					Resolution: resolution,
					Notes:      notes,
				}

				err := client.ResolveObjection(ctx, objectionID, uids[i], results[i])

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
				return formTmpl(w, data)
			}

			data.Success = true
			return RedirectError(fmt.Sprintf("/lpa/%s", caseUID))
		}

		return formTmpl(w, data)
	}
}
