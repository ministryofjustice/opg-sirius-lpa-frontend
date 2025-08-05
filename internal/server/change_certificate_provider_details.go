package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/sync/errgroup"
	"net/http"
	"time"
)

type ChangeCertificateProviderDetailsClient interface {
	CaseSummary(sirius.Context, string) (sirius.CaseSummary, error)
	ChangeCertificateProviderDetails(sirius.Context, string, sirius.ChangeCertificateProviderDetails) error
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
	FirstNames string         `form:"firstNames"`
	LastName   string         `form:"lastName"`
	Address    sirius.Address `form:"address"`
	Channel    string         `form:"channel"`
	Email      string         `form:"email"`
	Phone      string         `form:"phone"`
	SignedAt   dob            `form:"signedAt"`
}

func ChangeCertificateProviderDetails(client ChangeCertificateProviderDetailsClient, tmpl template.Template) Handler {

	return func(w http.ResponseWriter, r *http.Request) error {
		caseUid := r.PathValue("uid")
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
			err := decoder.Decode(&data.Form, r.PostForm)
			if err != nil {
				return err
			}

			certificateProviderDetailsData := sirius.ChangeCertificateProviderDetails{
				FirstNames: data.Form.FirstNames,
				LastName:   data.Form.LastName,
				Address:    data.Form.Address,
				Phone:      data.Form.Phone,
				Email:      data.Form.Email,
				SignedAt:   data.Form.SignedAt.toDateString(),
			}

			err = client.ChangeCertificateProviderDetails(ctx, caseUid, certificateProviderDetailsData)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
			} else if err != nil {
				return err
			} else {
				SetFlash(w, FlashNotification{
					Title: "Update saved aaaa",
				})

				// This sleep is to give LpaStore time to fire off an lpa-updated event to Sirius after it receives the
				// update, which recalculates anomaly before reloading the page
				time.Sleep(1 * time.Second)

				return RedirectError(fmt.Sprintf("/lpa/%s/lpa-details#certificate-provider", caseUid))
			}
		}

		return tmpl(w, data)
	}
}
