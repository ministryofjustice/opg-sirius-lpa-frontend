package server

import (
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"net/http"
)

type ObjectionOutcomeClient interface {
	GetObjection(sirius.Context, string) (sirius.Objection, error)
}

type objectionOutcomeData struct {
	Objection  sirius.Objection
	Resolution sirius.ObjectionResolution
}

func ObjectionOutcome(client ObjectionOutcomeClient, formTmpl template.Template) Handler {

	return func(w http.ResponseWriter, r *http.Request) error {
		caseUID := r.PathValue("uid")
		objectionID := r.PathValue("id")

		ctx := getContext(r)

		obj, err := client.GetObjection(ctx, objectionID)
		if err != nil {
			return err
		}

		data := objectionOutcomeData{
			Objection: obj,
		}

		for _, resolution := range obj.Resolutions {
			if resolution.Uid == caseUID {
				data.Resolution = resolution
			}
		}

		return formTmpl(w, data)
	}
}
