package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type EditComplaintClient interface {
	Case(ctx sirius.Context, id int) (sirius.Case, error)
	Complaint(ctx sirius.Context, id int) (sirius.Complaint, error)
	EditComplaint(ctx sirius.Context, id int, complaint sirius.Complaint) error
	RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error)
}

type editComplaintData struct {
	XSRFToken string
	Success   bool
	Error     sirius.ValidationError

	Categories            map[string]complaintCategory
	ComplainantCategories []sirius.RefDataItem
	Origins               []sirius.RefDataItem
	CompensationTypes     []sirius.RefDataItem

	Complaint sirius.Complaint
}

func EditComplaint(client EditComplaintClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}

		ctx := getContext(r)

		data := editComplaintData{
			XSRFToken:  ctx.XSRFToken,
			Categories: complaintCategories,
		}

		if err != nil {
			return err
		}

		data.CompensationTypes, err = client.RefDataByCategory(ctx, sirius.CompensationType)
		if err != nil {
			return err
		}

		data.Complaint, err = client.Complaint(ctx, id)
		if err != nil {
			return err
		}

		if r.Method == http.MethodPost {
			complaint := sirius.Complaint{
				Category:             postFormString(r, "category"),
				Description:          postFormString(r, "description"),
				ReceivedDate:         postFormDateString(r, "receivedDate"),
				Severity:             postFormString(r, "severity"),
				InvestigatingOfficer: postFormString(r, "investigatingOfficer"),
				SubCategory:          getValidSubcategory(postFormString(r, "category"), r.PostForm["subCategory"]),
				ComplainantCategory:  postFormString(r, "complainantCategory"),
				Origin:               postFormString(r, "origin"),
				CompensationType:     postFormString(r, "compensationType"),
				Summary:              postFormString(r, "summary"),
				Resolution:           postFormString(r, "resolution"),
				ResolutionInfo:       postFormString(r, "resolutionInfo"),
				ResolutionDate:       postFormDateString(r, "resolutionDate"),
			}

			if complaint.CompensationType != "NOT_APPLICABLE" {
				complaint.CompensationAmount = postFormString(r, fmt.Sprintf("compensationAmount%s", complaint.CompensationType))
			}

			err = client.EditComplaint(ctx, id, complaint)

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
