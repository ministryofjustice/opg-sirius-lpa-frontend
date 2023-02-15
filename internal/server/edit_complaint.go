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

type addOrEditComplaintData struct {
	XSRFToken string
	Entity    string
	Success   bool
	Error     sirius.ValidationError

	Categories            map[string]complaintCategory
	ComplainantCategories []sirius.RefDataItem
	Origins               []sirius.RefDataItem
	CompensationTypes     []sirius.RefDataItem

	Complaint sirius.Complaint
}

var fieldsToBeValidated = []string{"category", "severity", "investigatingOfficer", "complainantName", "origin", "compensationType", "summary", "resolutionDate", "receivedDate", "complainantCategory"}

func EditComplaint(client EditComplaintClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}

		ctx := getContext(r)

		data := addOrEditComplaintData{
			XSRFToken:  ctx.XSRFToken,
			Categories: complaintCategories,
		}

		data.ComplainantCategories, err = client.RefDataByCategory(ctx, sirius.ComplainantCategory)
		if err != nil {
			return err
		}

		data.Origins, err = client.RefDataByCategory(ctx, sirius.ComplaintOrigin)
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
			resolution := postFormString(r, "resolution")

			if resolution != "" {
				fieldErrors := sirius.FieldErrors{}

				for _, field := range fieldsToBeValidated {
					if !isFieldPopulated(field, r) {
						fieldErrors[field] = map[string]string{
							"reason": "Value is required and can't be empty",
						}
					}
				}

				if len(fieldErrors) > 0 {
					w.WriteHeader(http.StatusBadRequest)
					data.Error.Field = fieldErrors
					data.Complaint = populateComplaint(r)
					return tmpl(w, data)
				}
			}

			category := postFormString(r, "category")
			if category != "" {
				if getValidSubcategory(postFormString(r, "category"), r.PostForm["subCategory"]) == "" {
					return getSubcategoryValidationError(w, r, tmpl, data)
				}
			}

			complaint := populateComplaint(r)

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

func populateComplaint(r *http.Request) sirius.Complaint {
	return sirius.Complaint{
		Category:             postFormString(r, "category"),
		Description:          postFormString(r, "description"),
		ReceivedDate:         postFormDateString(r, "receivedDate"),
		Severity:             postFormString(r, "severity"),
		InvestigatingOfficer: postFormString(r, "investigatingOfficer"),
		ComplainantName:      postFormString(r, "complainantName"),
		SubCategory:          getValidSubcategory(postFormString(r, "category"), r.PostForm["subCategory"]),
		ComplainantCategory:  postFormString(r, "complainantCategory"),
		Origin:               postFormString(r, "origin"),
		CompensationType:     postFormString(r, "compensationType"),
		Summary:              postFormString(r, "summary"),
		Resolution:           postFormString(r, "resolution"),
		ResolutionInfo:       postFormString(r, "resolutionInfo"),
		ResolutionDate:       postFormDateString(r, "resolutionDate"),
	}
}

func isFieldPopulated(field string, r *http.Request) bool {
	if field == "receivedDate" || field == "resolutionDate" {
		if postFormDateString(r, field) == "" {
			return false
		}
	} else {
		if postFormString(r, field) == "" {
			return false
		}
	}
	return true
}

func getSubcategoryValidationError(w http.ResponseWriter, r *http.Request, tmpl template.Template, data addOrEditComplaintData) error {
	w.WriteHeader(http.StatusBadRequest)
	data.Error = sirius.ValidationError{
		Field: sirius.FieldErrors{
			"subCategory": {"reason": "Please select a subcategory"},
		},
	}
	data.Complaint = populateComplaint(r)
	return tmpl(w, data)
}
