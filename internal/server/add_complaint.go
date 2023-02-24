package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
)

type AddComplaintClient interface {
	AddComplaint(ctx sirius.Context, caseID int, caseType sirius.CaseType, complaint sirius.Complaint) error
	Case(ctx sirius.Context, id int) (sirius.Case, error)
	RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error)
}

type addComplaintData struct {
	XSRFToken string
	Entity    string
	Success   bool
	Error     sirius.ValidationError

	Categories            []sirius.RefDataItem
	ComplainantCategories []sirius.RefDataItem
	Origins               []sirius.RefDataItem

	Complaint sirius.Complaint
}

func AddComplaint(client AddComplaintClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		caseID, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}

		caseType, err := sirius.ParseCaseType(r.FormValue("case"))
		if err != nil {
			return err
		}

		ctx := getContext(r)
		group, groupCtx := errgroup.WithContext(ctx.Context)

		data := addComplaintData{
			XSRFToken: ctx.XSRFToken,
		}

		group.Go(func() error {
			caseitem, err := client.Case(ctx.With(groupCtx), caseID)
			if err != nil {
				return err
			}

			data.Entity = fmt.Sprintf("%s %s", caseitem.CaseType, caseitem.UID)
			return nil
		})

		group.Go(func() error {
			data.ComplainantCategories, err = client.RefDataByCategory(ctx.With(groupCtx), sirius.ComplainantCategory)
			if err != nil {
				return err
			}

			return nil
		})

		group.Go(func() error {
			data.Origins, err = client.RefDataByCategory(ctx.With(groupCtx), sirius.ComplaintOrigin)
			if err != nil {
				return err
			}

			return nil
		})

		group.Go(func() error {
			data.Categories, err = client.RefDataByCategory(ctx.With(groupCtx), sirius.ComplaintCategory)
			if err != nil {
				return err
			}

			return nil
		})

		if err := group.Wait(); err != nil {
			return err
		}

		if r.Method == http.MethodPost {
			complaint := sirius.Complaint{
				Category:             postFormString(r, "category"),
				Description:          postFormString(r, "description"),
				ReceivedDate:         postFormDateString(r, "receivedDate"),
				Severity:             postFormString(r, "severity"),
				InvestigatingOfficer: postFormString(r, "investigatingOfficer"),
				ComplainantName:      postFormString(r, "complainantName"),
				SubCategory:          postFormString(r, "subCategory"),
				ComplainantCategory:  postFormString(r, "complainantCategory"),
				Origin:               postFormString(r, "origin"),
				Summary:              postFormString(r, "summary"),
			}

			err = client.AddComplaint(ctx, caseID, caseType, complaint)

			if ve, ok := err.(sirius.ValidationError); ok {
				w.WriteHeader(http.StatusBadRequest)
				data.Error = ve
				data.Complaint = complaint
			} else if err != nil {
				return err
			} else {
				data.Success = true
			}
		}

		return tmpl(w, data)
	}
}
