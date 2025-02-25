package server

import (
	"fmt"
	"github.com/go-playground/form/v4"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/sync/errgroup"
	"net/http"
)

type ChangeDraftClient interface {
	CaseSummary(sirius.Context, string) (sirius.CaseSummary, error)
	ChangeDraft(sirius.Context, string, sirius.ChangeDraft) error
	RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error)
}

type changeDraftData struct {
	XSRFToken string
	Countries []sirius.RefDataItem
	Success   bool
	Error     sirius.ValidationError
	CaseUID   string
	Form      formDraftDetails
}

type formDraftDetails struct {
	FirstNames  string            `form:"firstNames"`
	LastName    string            `form:"lastName"`
	DateOfBirth sirius.DateString `form:"dob"`
	Address     sirius.Address    `form:"address"`
	PhoneNumber string            `form:"phoneNumber"`
	Email       string            `form:"email"`
}

func ChangeDraft(client ChangeDraftClient, tmpl template.Template) Handler {
	if decoder == nil {
		decoder = form.NewDecoder()
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		caseUID := r.FormValue("uid")

		ctx := getContext(r)

		cs, err := client.CaseSummary(ctx, caseUID)
		if err != nil {
			return err
		}

		donor := cs.DigitalLpa.SiriusData.Donor

		data := changeDraftData{
			XSRFToken: ctx.XSRFToken,
			CaseUID:   caseUID,
			Form: formDraftDetails{
				FirstNames:  donor.Firstname,
				LastName:    donor.Surname,
				DateOfBirth: donor.DateOfBirth,
				Address: sirius.Address{
					Line1:    donor.AddressLine1,
					Line2:    donor.AddressLine2,
					Line3:    donor.AddressLine3,
					Town:     donor.Town,
					Postcode: donor.Postcode,
					Country:  donor.Country,
				},
				PhoneNumber: donor.Phone,
				Email:       donor.Email,
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

		if r.Method == http.MethodPost {
			err := decoder.Decode(&data.Form, r.PostForm)
			if err != nil {
				return err
			}

			draftData := sirius.ChangeDraft{
				FirstNames: data.Form.FirstNames,
				LastName:   data.Form.LastName,
				//DateOfBirth: data.Form.DateOfBirth.toDateString(),
				Address: data.Form.Address,
				Phone:   data.Form.PhoneNumber,
				Email:   data.Form.Email,
			}

			err = client.ChangeDraft(ctx, caseUID, draftData)

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
