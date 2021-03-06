package server

import (
	"net/http"
	"strconv"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type EditComplaintClient interface {
	EditComplaint(ctx sirius.Context, id int, complaint sirius.Complaint) error
	Complaint(ctx sirius.Context, id int) (sirius.Complaint, error)
	Case(ctx sirius.Context, id int) (sirius.Case, error)
}

type editComplaintData struct {
	XSRFToken string
	Success   bool
	Error     sirius.ValidationError

	Categories map[string]complaintCategory
	Complaint  sirius.Complaint
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

		if r.Method == http.MethodPost {
			complaint := sirius.Complaint{
				Category:       postFormString(r, "category"),
				Description:    postFormString(r, "description"),
				ReceivedDate:   postFormDateString(r, "receivedDate"),
				Severity:       postFormString(r, "severity"),
				SubCategory:    getValidSubcategory(postFormString(r, "category"), r.PostForm["subCategory"]),
				Summary:        postFormString(r, "summary"),
				Resolution:     postFormString(r, "resolution"),
				ResolutionInfo: postFormString(r, "resolutionInfo"),
				ResolutionDate: postFormDateString(r, "resolutionDate"),
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

		complaint, err := client.Complaint(ctx, id)
		if err != nil {
			return err
		}
		data.Complaint = complaint

		return tmpl(w, data)
	}
}
