package server

import (
	"net/http"
	"strings"
	"time"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type dob struct {
	Day   int
	Month int
	Year  int
}

type draft struct {
	SubTypes                []string
	DonorFirstname          string
	DonorMiddlename         string
	DonorSurname            string
	DonorAddress            sirius.Address
	Recipient               string
	CorrespondentFirstname  string
	CorrespondentMiddlename string
	CorrespondentSurname    string
	AlternativeAddress      sirius.Address
	CorrespondentAddress    sirius.Address
	Dob                     dob
	Email                   string
	Phone                   string
}

type CreateDraftClient interface {
	CreateDraft(ctx sirius.Context, draft sirius.Draft) (map[string]string, error)
}

type createDraftData struct {
	XSRFToken string
	Draft     draft
	Error     sirius.ValidationError
	Success   bool
	Uids      map[string]string
}

func buildName(parts ...string) string {
	nonEmptyParts := []string{}

	for _, part := range parts {
		if part != "" {
			nonEmptyParts = append(nonEmptyParts, part)
		}
	}

	return strings.Join(nonEmptyParts, " ")
}

func fallback(val string, fallback string) string {
	if val == "" {
		return fallback
	} else {
		return val
	}
}

func CreateDraft(client CreateDraftClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		data := createDraftData{
			XSRFToken: ctx.XSRFToken,
		}

		if r.Method == "POST" {
			dobDay, _ := postFormInt(r, "dobDay")
			dobMonth, _ := postFormInt(r, "dobMonth")
			dobYear, _ := postFormInt(r, "dobYear")

			data.Draft = draft{
				SubTypes:        r.PostForm["subtype"],
				DonorFirstname:  postFormString(r, "donorFirstname"),
				DonorMiddlename: postFormString(r, "donorMiddlename"),
				DonorSurname:    postFormString(r, "donorSurname"),
				Dob: dob{
					Day:   dobDay,
					Month: dobMonth,
					Year:  dobYear,
				},
				DonorAddress: sirius.Address{
					Line1:    postFormString(r, "donorAddressLine1"),
					Line2:    postFormString(r, "donorAddressLine2"),
					Line3:    postFormString(r, "donorAddressLine3"),
					Town:     postFormString(r, "donorTown"),
					Postcode: postFormString(r, "donorPostcode"),
					Country:  postFormString(r, "donorCountry"),
				},
				CorrespondentFirstname:  postFormString(r, "correspondentFirstname"),
				CorrespondentMiddlename: postFormString(r, "correspondentMiddlename"),
				CorrespondentSurname:    postFormString(r, "correspondentSurname"),
				AlternativeAddress: sirius.Address{
					Line1:    postFormString(r, "alternativeAddressLine1"),
					Line2:    postFormString(r, "alternativeAddressLine2"),
					Line3:    postFormString(r, "alternativeAddressLine3"),
					Town:     postFormString(r, "alternativeTown"),
					Postcode: postFormString(r, "alternativePostcode"),
					Country:  postFormString(r, "alternativeCountry"),
				},
				CorrespondentAddress: sirius.Address{
					Line1:    postFormString(r, "correspondentAddressLine1"),
					Line2:    postFormString(r, "correspondentAddressLine2"),
					Line3:    postFormString(r, "correspondentAddressLine3"),
					Town:     postFormString(r, "correspondentTown"),
					Postcode: postFormString(r, "correspondentPostcode"),
					Country:  postFormString(r, "correspondentCountry"),
				},
				Recipient: postFormString(r, "recipient"),
				Email:     postFormString(r, "donorEmail"),
				Phone:     postFormString(r, "donorPhone"),
			}

			compiledDraft := sirius.Draft{
				CaseType:  data.Draft.SubTypes,
				Source:    "PHONE",
				DonorName: buildName(data.Draft.DonorFirstname, data.Draft.DonorMiddlename, data.Draft.DonorSurname),
				DonorAddress: sirius.Address{
					Line1:    data.Draft.DonorAddress.Line1,
					Line2:    data.Draft.DonorAddress.Line2,
					Line3:    data.Draft.DonorAddress.Line3,
					Town:     data.Draft.DonorAddress.Town,
					Postcode: data.Draft.DonorAddress.Postcode,
					Country:  fallback(data.Draft.DonorAddress.Country, "GB"),
				},
				Email:       data.Draft.Email,
				PhoneNumber: data.Draft.Phone,
			}

			if data.Draft.Recipient == "donor-other-address" {
				compiledDraft.CorrespondentAddress = &sirius.Address{
					Line1:    data.Draft.AlternativeAddress.Line1,
					Line2:    data.Draft.AlternativeAddress.Line2,
					Line3:    data.Draft.AlternativeAddress.Line3,
					Town:     data.Draft.AlternativeAddress.Town,
					Postcode: data.Draft.AlternativeAddress.Postcode,
					Country:  fallback(data.Draft.AlternativeAddress.Country, "GB"),
				}
			} else if data.Draft.Recipient == "other" {
				compiledDraft.CorrespondentName = buildName(data.Draft.CorrespondentFirstname, data.Draft.CorrespondentMiddlename, data.Draft.CorrespondentSurname)
				compiledDraft.CorrespondentAddress = &sirius.Address{
					Line1:    data.Draft.CorrespondentAddress.Line1,
					Line2:    data.Draft.CorrespondentAddress.Line2,
					Line3:    data.Draft.CorrespondentAddress.Line3,
					Town:     data.Draft.CorrespondentAddress.Town,
					Postcode: data.Draft.CorrespondentAddress.Postcode,
					Country:  fallback(data.Draft.CorrespondentAddress.Country, "GB"),
				}
			}

			dob := time.Date(data.Draft.Dob.Year, time.Month(data.Draft.Dob.Month), data.Draft.Dob.Day, 0, 0, 0, 0, time.UTC)
			date1900, _ := time.Parse("2006-01-02", "1900-01-01")
			if dob.After(date1900) {
				compiledDraft.DonorDob = sirius.DateString(dob.Format("2006-01-02"))
			}
			uids, err := client.CreateDraft(ctx, compiledDraft)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
			} else if err != nil {
				return err
			} else {
				data.Success = true
				data.Uids = uids
			}
		}

		return tmpl(w, data)
	}
}
