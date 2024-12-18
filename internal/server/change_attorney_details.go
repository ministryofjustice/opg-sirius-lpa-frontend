package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/form/v4"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/sync/errgroup"
	"net/http"
)

type ChangeAttorneyDetailsClient interface {
	CaseSummary(sirius.Context, string) (sirius.CaseSummary, error)
	ChangeAttorneyDetails(sirius.Context, string, string, sirius.ChangeAttorneyDetails) error
	RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error)
}

type changeAttorneyDetailsData struct {
	XSRFToken      string
	Countries      []sirius.RefDataItem
	Success        bool
	Error          sirius.ValidationError
	CaseUID        string
	Form           formAttorneyDetails
	AttorneyStatus string
}

type formAttorneyDetails struct {
	FirstNames  string         `form:"firstNames"`
	LastName    string         `form:"lastName"`
	DateOfBirth dob            `form:"dob"`
	Address     sirius.Address `form:"address"`
	PhoneNumber string         `form:"phoneNumber"`
	Email       string         `form:"email"`
	SignedAt    dob            `form:"signedAt"`
}

func ChangeAttorneyDetails(client ChangeAttorneyDetailsClient, tmpl template.Template) Handler {
	if decoder == nil {
		decoder = form.NewDecoder()
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		caseUID := chi.URLParam(r, "uid")
		attorneyUID := chi.URLParam(r, "attorneyUID")

		ctx := getContext(r)

		cs, err := client.CaseSummary(ctx, caseUID)
		if err != nil {
			return err
		}

		attorneys := cs.DigitalLpa.LpaStoreData.Attorneys

		var attorney sirius.LpaStoreAttorney
		for i, att := range attorneys {
			if att.Uid == attorneyUID {
				attorney = attorneys[i]
			}
		}

		attorneyDob, err := parseDate(attorney.DateOfBirth)
		if err != nil {
			return err
		}

		data := changeAttorneyDetailsData{
			XSRFToken: ctx.XSRFToken,
			CaseUID:   caseUID,
			Form: formAttorneyDetails{
				FirstNames:  attorney.FirstNames,
				LastName:    attorney.LastName,
				DateOfBirth: attorneyDob,
				Address: sirius.Address{
					Line1:    attorney.Address.Line1,
					Line2:    attorney.Address.Line2,
					Line3:    attorney.Address.Line3,
					Town:     attorney.Address.Town,
					Postcode: attorney.Address.Postcode,
					Country:  attorney.Address.Country,
				},
				PhoneNumber: attorney.Mobile,
				Email:       attorney.Email,
			},
			AttorneyStatus: attorney.Status,
		}

		if attorney.SignedAt != "" {
			signedAt, err := parseDateTime(attorney.SignedAt)
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

			attorneyDetailsData := sirius.ChangeAttorneyDetails{
				FirstNames:  data.Form.FirstNames,
				LastName:    data.Form.LastName,
				DateOfBirth: data.Form.DateOfBirth.toDateString(),
				Address:     data.Form.Address,
				Phone:       data.Form.PhoneNumber,
				Email:       data.Form.Email,
				SignedAt:    data.Form.SignedAt.toDateString(),
			}

			err = client.ChangeAttorneyDetails(ctx, caseUID, attorneyUID, attorneyDetailsData)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
			} else if err != nil {
				return err
			} else {
				data.Success = true

				SetFlash(w, FlashNotification{
					Title: "Changes confirmed",
				})

				return RedirectError(fmt.Sprintf("/lpa/%s/lpa-details", caseUID))
			}
		}

		return tmpl(w, data)
	}
}
