package server

import (
	"net/http"
	"strconv"

	"github.com/go-playground/form/v4"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/sync/errgroup"
)

type formDraftLpa struct {
	SubTypes                  []string       `form:"subtype"`
	Recipient                 string         `form:"recipient"`
	CorrespondentFirstname    string         `form:"correspondentFirstname"`
	CorrespondentSurname      string         `form:"correspondentSurname"`
	AlternativeAddress        sirius.Address `form:"alternativeAddress"`
	CorrespondentAddress      sirius.Address `form:"correspondentAddress"`
	CorrespondenceByWelsh     bool           `form:"correspondenceByWelsh"`
	CorrespondenceLargeFormat bool           `form:"correspondenceLargeFormat"`
}

type CreateDraftLpaClient interface {
	CreateDraftLpa(ctx sirius.Context, donorID int, draft sirius.DraftLpa) (map[string]string, error)
	GetUserDetails(ctx sirius.Context) (sirius.User, error)
	Person(ctx sirius.Context, personID int) (sirius.Person, error)
}

type createDraftLpaResult struct {
	Uid     string
	Subtype string
}

type createDraftLpaData struct {
	XSRFToken string
	Countries []sirius.RefDataItem
	Donor     sirius.Person
	Form      formDraftLpa
	Error     sirius.ValidationError
	Success   bool
	Uids      []createDraftLpaResult
}

func CreateDraftLpa(client CreateDraftLpaClient, tmpl template.Template) Handler {
	if decoder == nil {
		decoder = form.NewDecoder()
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}

		ctx := getContext(r)
		donor, err := client.Person(ctx, id)
		if err != nil {
			return err
		}

		group, _ := errgroup.WithContext(ctx.Context)

		data := createDraftLpaData{
			XSRFToken: ctx.XSRFToken,
			Donor:     donor,
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

		if err := group.Wait(); err != nil {
			return err
		}

		if r.Method == "POST" {
			err := decoder.Decode(&data.Form, r.PostForm)
			if err != nil {
				return err
			}

			compiledDraft := sirius.DraftLpa{
				CaseType:                  data.Form.SubTypes,
				CorrespondenceByWelsh:     data.Form.CorrespondenceByWelsh,
				CorrespondenceLargeFormat: data.Form.CorrespondenceLargeFormat,
			}

			if data.Form.Recipient == "donor-other-address" {
				correspondentAddress := addDefaultCountry(data.Form.AlternativeAddress)
				compiledDraft.CorrespondentAddress = &correspondentAddress
			} else if data.Form.Recipient == "other" {
				correspondentAddress := addDefaultCountry(data.Form.CorrespondentAddress)
				compiledDraft.CorrespondentAddress = &correspondentAddress
				compiledDraft.CorrespondentFirstNames = data.Form.CorrespondentFirstname
				compiledDraft.CorrespondentLastName = data.Form.CorrespondentSurname
			}

			uids, err := client.CreateDraftLpa(ctx, id, compiledDraft)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
			} else if err != nil {
				return err
			} else {
				data.Success = true

				for subtype, uid := range uids {
					data.Uids = append(data.Uids, createDraftLpaResult{
						Subtype: subtype,
						Uid:     uid,
					})
				}
			}
		}

		return tmpl(w, data)
	}
}
