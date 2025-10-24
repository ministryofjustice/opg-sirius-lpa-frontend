package server

import (
	"fmt"
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
	ProgressIndicatorsForDigitalLpa(siriusCtx sirius.Context, uid string) ([]sirius.ProgressIndicator, error)
}

type changeDonorDetailsData struct {
	XSRFToken                  string
	Countries                  []sirius.RefDataItem
	Success                    bool
	Error                      sirius.ValidationError
	CaseUID                    string
	Form                       formDonorDetails
	DonorIdentityCheckComplete bool
	DonorDobString             string
	SignedByWitnessTwoLabel    string
}

type formDonorDetails struct {
	FirstNames                string         `form:"firstNames"`
	LastName                  string         `form:"lastName"`
	OtherNamesKnownBy         string         `form:"otherNamesKnownBy"`
	DateOfBirth               dob            `form:"dob"`
	Address                   sirius.Address `form:"address"`
	PhoneNumber               string         `form:"phoneNumber"`
	Email                     string         `form:"email"`
	LpaSignedOn               dob            `form:"lpaSignedOn"`
	AuthorisedSignatory       string         `form:"authorisedSignatory"`
	SignedByWitnessOne        string         `form:"signedByWitnessOne"`
	SignedByWitnessTwo        string         `form:"signedByWitnessTwo"`
	IndependentWitnessName    string         `form:"independentWitnessName"`
	IndependentWitnessAddress sirius.Address `form:"independentWitnessAddress"`
}

func parseDate(dateString string) (dob, error) {
	parsedTime, err := time.Parse("2006-01-02", dateString) // Parses date in "YYYY-MM-DD" format
	if err != nil {
		return dob{}, err
	}

	return dob{
		Day:   parsedTime.Day(),
		Month: int(parsedTime.Month()),
		Year:  parsedTime.Year(),
	}, nil
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
	return func(w http.ResponseWriter, r *http.Request) error {
		caseUID := r.FormValue("uid")

		ctx := getContext(r)

		cs, err := client.CaseSummary(ctx, caseUID)
		if err != nil {
			return err
		}

		lpaStore := cs.DigitalLpa.LpaStoreData
		donorDob, err := parseDate(lpaStore.Donor.DateOfBirth)
		if err != nil {
			return err
		}

		signedAt, err := parseDateTime(lpaStore.SignedAt)
		if err != nil {
			return err
		}

		donorIdentityCheckComplete := false
		donorDobString := lpaStore.Donor.DateOfBirth

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

		signedByWitnessOne := "No"
		signedByWitnessTwo := "No"
		lpaStoreWitnessedByCertificateProviderAt := time.Time{}

		if lpaStore.WitnessedByCertificateProviderAt != "" {
			lpaStoreWitnessedByCertificateProviderAt, err = time.Parse(time.RFC3339, lpaStore.WitnessedByCertificateProviderAt)

			if err != nil {
				return err
			}
		}

		if !lpaStoreWitnessedByCertificateProviderAt.IsZero() {
			signedByWitnessOne = "Yes"
		}

		if lpaStore.WitnessedByIndependentWitnessAt != "" {
			signedByWitnessTwo = "Yes"
		}

		signedByWitnessTwoLabel := "Signed by witness 2"
		independentWitnessName := ""

		if lpaStore.IndependentWitness != nil {
			independentWitnessName = lpaStore.IndependentWitness.FirstNames + " " + lpaStore.IndependentWitness.LastName

			if lpaStore.IndependentWitness.FirstNames != "" || lpaStore.IndependentWitness.LastName != "" {
				signedByWitnessTwoLabel += " - " + independentWitnessName
			}
		}

		independentWitnessAddress := sirius.Address{}

		if lpaStore.IndependentWitness != nil {
			independentWitnessAddress.Line1 = lpaStore.IndependentWitness.Address.Line1
			independentWitnessAddress.Line2 = lpaStore.IndependentWitness.Address.Line2
			independentWitnessAddress.Line3 = lpaStore.IndependentWitness.Address.Line3
			independentWitnessAddress.Town = lpaStore.IndependentWitness.Address.Town
			independentWitnessAddress.Postcode = lpaStore.IndependentWitness.Address.Postcode
			independentWitnessAddress.Country = lpaStore.IndependentWitness.Address.Country
		}

		authorisedSignatoryName := ""

		if lpaStore.AuthorisedSignatory != nil {
			authorisedSignatoryName = lpaStore.AuthorisedSignatory.FirstNames + " " + lpaStore.AuthorisedSignatory.LastName
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
				Email:                     lpaStore.Donor.Email,
				PhoneNumber:               cs.DigitalLpa.SiriusData.Application.PhoneNumber,
				LpaSignedOn:               signedAt,
				AuthorisedSignatory:       authorisedSignatoryName,
				SignedByWitnessOne:        signedByWitnessOne,
				SignedByWitnessTwo:        signedByWitnessTwo,
				IndependentWitnessName:    independentWitnessName,
				IndependentWitnessAddress: independentWitnessAddress,
			},
			DonorIdentityCheckComplete: donorIdentityCheckComplete,
			DonorDobString:             donorDobString,
			SignedByWitnessTwoLabel:    signedByWitnessTwoLabel,
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
				Address:           data.Form.Address,
				Phone:             data.Form.PhoneNumber,
				Email:             data.Form.Email,
				LpaSignedOn:       data.Form.LpaSignedOn.toDateString(),
			}

			witnessedByCertificateProviderAt := time.Time{}
			var witnessedByIndependentWitnessAt *time.Time
			signedAtTime, err := time.Parse(time.RFC3339, lpaStore.SignedAt)

			if err != nil {
				return err
			}

			if data.Form.SignedByWitnessOne == "Yes" {
				if lpaStoreWitnessedByCertificateProviderAt.IsZero() {
					witnessedByCertificateProviderAt = signedAtTime
				} else {
					witnessedByCertificateProviderAt = lpaStoreWitnessedByCertificateProviderAt
				}
			}

			if data.Form.SignedByWitnessTwo == "Yes" {
				if lpaStore.WitnessedByIndependentWitnessAt == "" {
					witnessedByIndependentWitnessAt = &signedAtTime
				} else {
					lpaStoreWitnessedByIndependentWitnessAt, err := time.Parse(time.RFC3339, lpaStore.WitnessedByIndependentWitnessAt)

					if err != nil {
						return err
					}

					witnessedByIndependentWitnessAt = &lpaStoreWitnessedByIndependentWitnessAt
				}
			}

			donorDetailsData.AuthorisedSignatory = data.Form.AuthorisedSignatory
			donorDetailsData.WitnessedByCertificateProviderAt = witnessedByCertificateProviderAt
			donorDetailsData.WitnessedByIndependentWitnessAt = witnessedByIndependentWitnessAt
			donorDetailsData.IndependentWitnessName = data.Form.IndependentWitnessName
			donorDetailsData.IndependentWitnessAddress = data.Form.IndependentWitnessAddress

			if donorIdentityCheckComplete {
				donorDetailsData.DateOfBirth = donorDob.toDateString()
			} else {
				donorDetailsData.DateOfBirth = data.Form.DateOfBirth.toDateString()
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
					Title: "Update saved",
				})

				return RedirectError(fmt.Sprintf("/lpa/%s/lpa-details", caseUID))
			}
		}

		return tmpl(w, data)
	}
}
