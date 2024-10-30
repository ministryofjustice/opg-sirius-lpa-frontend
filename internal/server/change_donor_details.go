package server

import (
	"fmt"
	"github.com/go-playground/form/v4"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/sync/errgroup"
	"net/http"
	"time"
)

type ChangeDonorDetailsClient interface {
	CaseSummary(sirius.Context, string) (sirius.CaseSummary, error)
	ChangeDonorDetails(sirius.Context, string, sirius.ChangeDonorDetails) error
	RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error)
}

type changeDonorDetailsData struct {
	XSRFToken string
	Countries []sirius.RefDataItem
	Success   bool
	Error     sirius.ValidationError
	CaseUID   string
	Form      formDonorDetails
}

type formDonorDetails struct {
	FirstNames        string         `form:"firstNames"`
	LastName          string         `form:"lastName"`
	OtherNamesKnownBy string         `form:"otherNamesKnownBy"`
	DateOfBirth       dob            `form:"dob"`
	Address           sirius.Address `form:"address"`
	PhoneNumber       string         `form:"phoneNumber"`
	Email             string         `form:"email"`
	LpaSignedOn       dob            `form:"lpaSignedOn"`
}

func parseDate(dateString string) dob {
	parsedTime, err := time.Parse("2006-01-02", dateString) // Parses date in "YYYY-MM-DD" format
	if err != nil {
		return dob{}
	}

	return dob{
		Day:   parsedTime.Day(),
		Month: int(parsedTime.Month()),
		Year:  parsedTime.Year(),
	}
}

func parseDateTime(dateTimeString string) (dob, error) {
	parsedTime, err := time.Parse(time.RFC3339, dateTimeString) // Parse ISO 8601 date-time
	if err != nil {
		return dob{}, err
	}

	return dob{
		Day:   parsedTime.Day(),
		Month: int(parsedTime.Month()),
		Year:  parsedTime.Year(),
	}, nil
}

func ChangeDonorDetails(client ChangeDonorDetailsClient, tmpl template.Template) Handler {
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

		lpaStore := cs.DigitalLpa.LpaStoreData
		donorDob := parseDate(lpaStore.Donor.DateOfBirth)

		signedAt, err := parseDateTime(lpaStore.SignedAt)
		if err != nil {
			return err
		}

		data := changeDonorDetailsData{
			XSRFToken: ctx.XSRFToken,
			CaseUID:   caseUID,
			Form: formDonorDetails{
				FirstNames:        lpaStore.Donor.FirstNames,
				LastName:          lpaStore.Donor.LastName,
				OtherNamesKnownBy: lpaStore.Donor.OtherNamesKnownBy,
				DateOfBirth:       donorDob,
				Address: sirius.Address{
					Line1:    lpaStore.Donor.Address.Line1,
					Line2:    lpaStore.Donor.Address.Line2,
					Line3:    lpaStore.Donor.Address.Line3,
					Town:     lpaStore.Donor.Address.Town,
					Postcode: lpaStore.Donor.Address.Postcode,
					Country:  lpaStore.Donor.Address.Country,
				},
				Email:       lpaStore.Donor.Email,
				PhoneNumber: cs.DigitalLpa.SiriusData.Application.PhoneNumber,
				LpaSignedOn: signedAt,
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

			donorDetailsData := sirius.ChangeDonorDetails{
				FirstNames:        data.Form.FirstNames,
				LastName:          data.Form.LastName,
				OtherNamesKnownBy: data.Form.OtherNamesKnownBy,
				DateOfBirth:       data.Form.DateOfBirth.toDateString(),
				Address:           data.Form.Address,
				Phone:             data.Form.PhoneNumber,
				Email:             data.Form.Email,
				LpaSignedOn:       data.Form.LpaSignedOn.toDateString(),
			}

			err = client.ChangeDonorDetails(ctx, caseUID, donorDetailsData)

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
