package server

import (
	"net/http"

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
	CreateDraftLpa(ctx sirius.Context, draft sirius.DraftLpa) (map[string]string, error)
	GetUserDetails(ctx sirius.Context) (sirius.User, error)
}

type createDraftLpaResult struct {
	Uid     string
	Subtype string
}

type createDraftLpaData struct {
	XSRFToken string
	Countries []sirius.RefDataItem
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
		ctx := getContext(r)
		group, _ := errgroup.WithContext(ctx.Context)

		data := createDraftLpaData{
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

			uids, err := client.CreateDraftLpa(ctx, compiledDraft)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
			} else if err != nil {
				return err
			} else {
				data.Success = true

				for subtype, uid := range uids {
					data.Uids = append(data.Uids, createDraftLpaResult{
						Uid:     uid,
						Subtype: subtype,
					})
				}
			}
		}

		return tmpl(w, data)
	}
}
