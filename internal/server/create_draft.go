package server

import (
	"fmt"
	"net/http"

	"github.com/go-playground/form/v4"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/sync/errgroup"
)

type dob struct {
	Day   int `form:"day"`
	Month int `form:"month"`
	Year  int `form:"year"`
}

func (d *dob) toDateString() sirius.DateString {
	if d.Day == 0 && d.Month == 0 && d.Year == 0 {
		return ""
	}

	monthSpacing := ""
	daySpacing := ""

	if d.Month < 10 {
		monthSpacing = "0"
	}
	if d.Day < 10 {
		daySpacing = "0"
	}

	return sirius.DateString(fmt.Sprintf("%d-%s%d-%s%d", d.Year, monthSpacing, d.Month, daySpacing, d.Day))
}

type formDraft struct {
	SubTypes                  []string       `form:"subtype"`
	DonorFirstname            string         `form:"donorFirstname"`
	DonorSurname              string         `form:"donorSurname"`
	DonorAddress              sirius.Address `form:"donorAddress"`
	Recipient                 string         `form:"recipient"`
	CorrespondentFirstname    string         `form:"correspondentFirstname"`
	CorrespondentSurname      string         `form:"correspondentSurname"`
	AlternativeAddress        sirius.Address `form:"alternativeAddress"`
	CorrespondentAddress      sirius.Address `form:"correspondentAddress"`
	Dob                       dob            `form:"dob"`
	Email                     string         `form:"donorEmail"`
	Phone                     string         `form:"donorPhone"`
	CorrespondenceByWelsh     bool           `form:"correspondenceByWelsh"`
	CorrespondenceLargeFormat bool           `form:"correspondenceLargeFormat"`
}

type CreateDraftClient interface {
	CreateDraft(ctx sirius.Context, draft sirius.Draft) (map[string]string, error)
	GetUserDetails(ctx sirius.Context) (sirius.User, error)
	RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error)
}

type createDraftResult struct {
	Uid     string
	Subtype string
}

type createDraftData struct {
	XSRFToken string
	Countries []sirius.RefDataItem
	Form      formDraft
	Error     sirius.ValidationError
	Success   bool
	Uids      []createDraftResult
}

func CreateDraft(client CreateDraftClient, tmpl template.Template) Handler {
	if decoder == nil {
		decoder = form.NewDecoder()
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)
		group, groupCtx := errgroup.WithContext(ctx.Context)

		data := createDraftData{
			XSRFToken: ctx.XSRFToken,
		}

		group.Go(func() error {
			user, err := client.GetUserDetails(ctx)
			if err != nil {
				return err
			}

			if !user.HasRole("private-mlpa") {
				return sirius.StatusError{Code: http.StatusForbidden}
			}

			return nil
		})

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

		if r.Method == "POST" {
			err := decoder.Decode(&data.Form, r.PostForm)
			if err != nil {
				return err
			}

			compiledDraft := sirius.Draft{
				CaseType:                  data.Form.SubTypes,
				Source:                    "PHONE",
				DonorFirstNames:           data.Form.DonorFirstname,
				DonorLastName:             data.Form.DonorSurname,
				DonorDob:                  data.Form.Dob.toDateString(),
				DonorAddress:              data.Form.DonorAddress,
				Email:                     data.Form.Email,
				PhoneNumber:               data.Form.Phone,
				CorrespondenceByWelsh:     data.Form.CorrespondenceByWelsh,
				CorrespondenceLargeFormat: data.Form.CorrespondenceLargeFormat,
			}

			if data.Form.Recipient == "donor-other-address" {
				correspondentAddress := data.Form.AlternativeAddress
				compiledDraft.CorrespondentAddress = &correspondentAddress
			} else if data.Form.Recipient == "other" {
				correspondentAddress := data.Form.CorrespondentAddress
				compiledDraft.CorrespondentAddress = &correspondentAddress
				compiledDraft.CorrespondentFirstNames = data.Form.CorrespondentFirstname
				compiledDraft.CorrespondentLastName = data.Form.CorrespondentSurname
			}

			uids, err := client.CreateDraft(ctx, compiledDraft)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
			} else if err != nil {
				return err
			} else {
				data.Success = true

				for subtype, uid := range uids {
					data.Uids = append(data.Uids, createDraftResult{
						Subtype: subtype,
						Uid:     uid,
					})
				}
			}
		}

		return tmpl(w, data)
	}
}
