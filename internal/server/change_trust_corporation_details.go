package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/sync/errgroup"
)

type ChangeTrustCorporationDetailsClient interface {
	CaseSummary(sirius.Context, string) (sirius.CaseSummary, error)
	RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error)
}

type changeTrustCorporationDetailsData struct {
	XSRFToken string
	Countries []sirius.RefDataItem
	Success   bool
	Error     sirius.ValidationError
	CaseUID   string
	Form      formTrustCorporationDetails
}

type formTrustCorporationDetails struct {
	Name          string         `form:"name"`
	Address       sirius.Address `form:"address"`
	Email         string         `form:"email"`
	PhoneNumber   string         `form:"phoneNumber"`
	CompanyNumber string         `form:"companyNumber"`
}

func ChangeTrustCorporationDetails(client ChangeTrustCorporationDetailsClient, tmpl template.Template) Handler {

	return func(w http.ResponseWriter, r *http.Request) error {
		caseUID := r.PathValue("uid")
		trustCorpUID := r.PathValue("trustCorporationUID")

		ctx := getContext(r)

		cs, err := client.CaseSummary(ctx, caseUID)
		if err != nil {
			return err
		}

		trustCorps := cs.DigitalLpa.LpaStoreData.TrustCorporations

		var trustCorporation sirius.LpaStoreTrustCorporation
		for i, trustCorp := range trustCorps {
			if trustCorp.Uid == trustCorpUID {
				trustCorporation = trustCorps[i]
			}
		}

		data := changeTrustCorporationDetailsData{
			XSRFToken: ctx.XSRFToken,
			CaseUID:   caseUID,
			Form: formTrustCorporationDetails{
				Name: trustCorporation.Name,
				Address: sirius.Address{
					Line1:    trustCorporation.Address.Line1,
					Line2:    trustCorporation.Address.Line2,
					Line3:    trustCorporation.Address.Line3,
					Town:     trustCorporation.Address.Town,
					Postcode: trustCorporation.Address.Postcode,
					Country:  trustCorporation.Address.Country,
				},
				Email:         trustCorporation.Email,
				PhoneNumber:   trustCorporation.Mobile,
				CompanyNumber: trustCorporation.CompanyNumber,
			},
		}

		group, groupCtx := errgroup.WithContext(ctx.Context)

		group.Go(func() error {
			var err error
			data.Countries, err = client.RefDataByCategory(ctx.With(groupCtx), sirius.CountryCategory)
			if err != nil {
				return err
			}

			return nil
		})

		if err := group.Wait(); err != nil {
			return err
		}

		return tmpl(w, data)
	}
}
