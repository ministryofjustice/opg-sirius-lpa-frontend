package server

import (
	"net/http"
	"strconv"

	"github.com/go-playground/form/v4"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/sync/errgroup"
)

type formAdditionalDraft struct {
	SubTypes                  []string       `form:"subtype"`
	Recipient                 string         `form:"recipient"`
	CorrespondentFirstname    string         `form:"correspondentFirstname"`
	CorrespondentSurname      string         `form:"correspondentSurname"`
	AlternativeAddress        sirius.Address `form:"alternativeAddress"`
	CorrespondentAddress      sirius.Address `form:"correspondentAddress"`
	CorrespondenceByWelsh     bool           `form:"correspondenceByWelsh"`
	CorrespondenceLargeFormat bool           `form:"correspondenceLargeFormat"`
}

type CreateAdditionalDraftClient interface {
	CreateAdditionalDraft(ctx sirius.Context, donorID int, draft sirius.AdditionalDraft) (map[string]string, error)
	GetUserDetails(ctx sirius.Context) (sirius.User, error)
	Person(ctx sirius.Context, personID int) (sirius.Person, error)
}

type createAdditionalDraftResult struct {
	Uid     string
	Subtype string
}

type createAdditionalDraftData struct {
	XSRFToken string
	Countries []sirius.RefDataItem
	Donor     sirius.Person
	Form      formAdditionalDraft
	Error     sirius.ValidationError
	Success   bool
	Uids      []createAdditionalDraftResult
}

func CreateAdditionalDraft(client CreateAdditionalDraftClient, tmpl template.Template) Handler {
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

		data := createAdditionalDraftData{
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

			compiledDraft := sirius.AdditionalDraft{
				CaseType:                  data.Form.SubTypes,
				CorrespondenceByWelsh:     data.Form.CorrespondenceByWelsh,
				CorrespondenceLargeFormat: data.Form.CorrespondenceLargeFormat,
				Source:                    "PHONE",
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

			uids, err := client.CreateAdditionalDraft(ctx, id, compiledDraft)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
			} else if err != nil {
				return err
			} else {
				data.Success = true

				for subtype, uid := range uids {
					data.Uids = append(data.Uids, createAdditionalDraftResult{
						Subtype: subtype,
						Uid:     uid,
					})
				}
			}
		}

		return tmpl(w, data)
	}
}
