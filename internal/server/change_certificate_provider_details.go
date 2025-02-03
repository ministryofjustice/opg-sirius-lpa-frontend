package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/sync/errgroup"
	"net/http"
)

type ChangeCertificateProviderDetailsClient interface {
	CaseSummary(sirius.Context, string) (sirius.CaseSummary, error)
	RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error)
}

type changeCertificateProviderDetailsData struct {
	XSRFToken string
	CaseUid   string
	Countries []sirius.RefDataItem
	Error     sirius.ValidationError
	Form      formCertificateProviderDetails
}

type formCertificateProviderDetails struct {
	FirstNames string         `json:"firstNames"`
	LastName   string         `json:"lastName"`
	Address    sirius.Address `json:"address"`
	Channel    string         `json:"channel"`
	Email      string         `json:"email"`
	Phone      string         `json:"phone"`
	SignedAt   dob            `json:"signedAt"`
}

func ChangeCertificateProviderDetails(client ChangeCertificateProviderDetailsClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		caseUid := chi.URLParam(r, "uid")
		ctx := getContext(r)
		caseSummary, err := client.CaseSummary(ctx, caseUid)

		if err != nil {
			return err
		}

		certificateProvider := caseSummary.DigitalLpa.LpaStoreData.CertificateProvider

		data := changeCertificateProviderDetailsData{
			XSRFToken: ctx.XSRFToken,
			CaseUid:   caseUid,
			Form: formCertificateProviderDetails{
				FirstNames: certificateProvider.FirstNames,
				LastName:   certificateProvider.LastName,
				Address: sirius.Address{
					Line1:    certificateProvider.Address.Line1,
					Line2:    certificateProvider.Address.Line2,
					Line3:    certificateProvider.Address.Line3,
					Town:     certificateProvider.Address.Town,
					Postcode: certificateProvider.Address.Postcode,
					Country:  certificateProvider.Address.Country,
				},
				Channel: certificateProvider.Channel,
				Email:   certificateProvider.Email,
				Phone:   certificateProvider.Phone,
			},
		}

		if certificateProvider.SignedAt != "" {
			signedAt, err := parseDateTime(certificateProvider.SignedAt)
			data.Form.SignedAt = signedAt

			if err != nil {
				return err
			}
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

		if r.Method == http.MethodPost {
			SetFlash(w, FlashNotification{
				Title: "Changes confirmed",
			})

			return RedirectError(fmt.Sprintf("/lpa/%s/lpa-details#certificate-provider", caseUid))
		}

		return tmpl(w, data)
	}
}
