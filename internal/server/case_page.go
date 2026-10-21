package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type CasePageClient interface {
	PageVarsClient
	FeatureToggles(ctx sirius.Context) (sirius.FeatureToggles, error)
}

type casePageData struct {
	ActionPanelButtons             []ActionPanelButton
	CaseUids                       string
	Donor                          sirius.Person
	DonorID                        int
	HasV1PersonsCasesGetPermission bool
	HasV1PersonsGetPermission      bool
	HeaderButtons                  SiriusHeaderButtons
	Person                         sirius.Person
	SelectedCases                  []sirius.Case
	XSRFToken                      string
}

func CasePage(client CasePageClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		featureToggles, err := client.FeatureToggles(ctx)
		if err != nil || !featureToggles.Enabled("poasNewCasePage") {
			return sirius.StatusError{Code: http.StatusForbidden}
		}

		if err := r.ParseForm(); err != nil {
			return err
		}

		pageVars, err := GetPageValues(client, r)
		if err != nil {
			return err
		}

		if pageVars.DonorID == 0 {
			return sirius.StatusError{Code: http.StatusNotFound}
		}

		uidParams := buildUIDQueryString(pageVars.CaseUidsCollection)

		data := casePageData{
			XSRFToken:                      ctx.XSRFToken,
			CaseUids:                       uidParams,
			Donor:                          pageVars.Person,
			DonorID:                        pageVars.DonorID,
			HasV1PersonsCasesGetPermission: pageVars.HasV1PersonsCasesGetPermission,
			HasV1PersonsGetPermission:      pageVars.HasV1PersonsGetPermission,
			Person:                         pageVars.Person,
			SelectedCases:                  pageVars.SelectedCases,
		}

		data.ActionPanelButtons = GetActionPanelButtons(data.SelectedCases, data.DonorID, uidParams, pageVars.DraftCount > 0, pageVars.PersonReferences, len(pageVars.Person.Children) > 0, pageVars.TaskIDs, pageVars.UserPermissions)

		data.HeaderButtons = SiriusHeaderButtons{
			Calendar:  true,
			Documents: true,
		}

		return tmpl(w, data)
	}
}
