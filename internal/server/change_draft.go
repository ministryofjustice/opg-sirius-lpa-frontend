package server

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/sync/errgroup"
)

type ChangeDraftClient interface {
	CaseSummary(sirius.Context, string) (sirius.CaseSummary, error)
	ChangeDraft(sirius.Context, string, sirius.ChangeDraft) error
	RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error)
	ProgressIndicatorsForDigitalLpa(siriusCtx sirius.Context, uid string) ([]sirius.ProgressIndicator, error)
}

type changeDraftData struct {
	XSRFToken                  string
	Countries                  []sirius.RefDataItem
	Success                    bool
	Error                      sirius.ValidationError
	CaseUID                    string
	Form                       formDraftDetails
	DonorIdentityCheckComplete bool
	DonorDobString             string
}

type formDraftDetails struct {
	FirstNames  string         `form:"firstNames"`
	LastName    string         `form:"lastName"`
	DateOfBirth dob            `form:"dob"`
	Address     sirius.Address `form:"address"`
	PhoneNumber string         `form:"phoneNumber"`
	Email       string         `form:"email"`
}

func ChangeDraft(client ChangeDraftClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		caseUID := r.PathValue("uid")
		ctx := getContext(r)

		var countries []sirius.RefDataItem
		var cs sirius.CaseSummary
		var data changeDraftData
		var err error

		group, groupCtx := errgroup.WithContext(ctx.Context)

		group.Go(func() error {
			cs, err = client.CaseSummary(ctx.With(groupCtx), caseUID)
			if err != nil {
				return err
			}
			return nil
		})

		group.Go(func() error {
			countries, err = client.RefDataByCategory(ctx.With(groupCtx), sirius.CountryCategory)
			if err != nil {
				return err
			}
			return nil
		})

		if err := group.Wait(); err != nil {
			return err
		}

		donor := cs.DigitalLpa.SiriusData.Donor

		dob, err := parseDate(string(donor.DateOfBirth))
		if err != nil {
			return err
		}

		donorIdentityCheckComplete := false
		donorDobString := string(donor.DateOfBirth)

		pis, err := client.ProgressIndicatorsForDigitalLpa(ctx, caseUID)
		if err != nil {
			return err
		}

		for _, pi := range pis {
			if pi.Indicator == "DONOR_ID" {
				if pi.Status == "COMPLETE" {
					donorIdentityCheckComplete = true
				}

				break
			}
		}

		data = changeDraftData{
			XSRFToken: ctx.XSRFToken,
			CaseUID:   caseUID,
			Form: formDraftDetails{
				FirstNames:  donor.Firstname,
				LastName:    donor.Surname,
				DateOfBirth: dob,
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
			Countries:                  countries,
			DonorIdentityCheckComplete: donorIdentityCheckComplete,
			DonorDobString:             donorDobString,
		}

		if r.Method == http.MethodPost {
			err := decoder.Decode(&data.Form, r.PostForm)
			if err != nil {
				return err
			}

			draftData := sirius.ChangeDraft{
				FirstNames:  data.Form.FirstNames,
				LastName:    data.Form.LastName,
				DateOfBirth: data.Form.DateOfBirth.toDateString(),
				Address:     data.Form.Address,
				Phone:       data.Form.PhoneNumber,
				Email:       data.Form.Email,
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
					Title: "Update saved",
				})

				return RedirectError(fmt.Sprintf("/lpa/%s/lpa-details", caseUID))
			}
		}

		return tmpl(w, data)
	}
}
