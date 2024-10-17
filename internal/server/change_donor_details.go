package server

import (
	"fmt"
	"github.com/go-playground/form/v4"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/sync/errgroup"
	"net/http"
)

type ChangeDonorDetailsClient interface {
	CaseSummary(sirius.Context, string) (sirius.CaseSummary, error)
	ChangeDonorDetails(sirius.Context, string, sirius.ChangeDonorDetailsData) error
	GetUserDetails(ctx sirius.Context) (sirius.User, error)
	RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error)
}

type ChangeDonorDetailsData struct {
	XSRFToken string
	Countries []sirius.RefDataItem
	Entity    string
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

		group, groupCtx := errgroup.WithContext(ctx.Context)

		data := ChangeDonorDetailsData{
			XSRFToken: ctx.XSRFToken,
			Entity:    fmt.Sprintf("%s %s", cs.DigitalLpa.SiriusData.Subtype, caseUID),
			CaseUID:   caseUID,
		}

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

			donorDetailsData := sirius.ChangeDonorDetailsData{
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

				return RedirectError(fmt.Sprintf("/lpa/%s/lpa-details", caseUID))
			}
		}

		return tmpl(w, data)
	}
}
