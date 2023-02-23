package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/sync/errgroup"
)

type AddComplaintClient interface {
	AddComplaint(ctx sirius.Context, caseID int, caseType sirius.CaseType, complaint sirius.Complaint) error
	Case(ctx sirius.Context, id int) (sirius.Case, error)
	RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error)
}

type complaintCategory struct {
	Label    string
	Children map[string]string
}

var complaintCategories = map[string]complaintCategory{
	"01": {
		Label: "Correspondence",
		Children: map[string]string{
			"06": "General Query",
			"07": "Chase up",
			"08": "Typo / Grammar",
			"09": "Quality of Documents",
			"10": "Third Parties",
			"11": "Refund Request",
			"12": "Digital Tool",
			"13": "Finance",
			"14": "Customer Service",
		},
	},
	"02": {
		Label: "OPG Decisions",
		Children: map[string]string{
			"15": "POA Decisions",
			"16": "Supervision Decisions",
			"17": "Investigation Outcomes",
			"18": "Fee Decision",
			"19": "Safeguarding Decisions",
			"20": "Other",
		},
	},
	"03": {
		Label: "Non OPG",
		Children: map[string]string{
			"21": "Banks / Utilities",
			"22": "COP / Judicial",
			"23": "DX / Royal Mail",
			"24": "Health / Social Care",
			"25": "Solicitors",
			"26": "Deputy / Attorney",
			"27": "Other",
		},
	},
	"04": {
		Label: "Customer Service",
		Children: map[string]string{
			"28": "Letter Content",
			"29": "Delays",
			"30": "Contact with OPG",
			"31": "Quality",
			"32": "Incorrect / Confusing Advice",
			"33": "Failure to Follow Procedure",
			"34": "Lost Documents",
			"35": "Security Breach",
			"36": "Other",
		},
	},
	"05": {
		Label: "Policy",
		Children: map[string]string{
			"37": "Mental Capacity Act",
			"38": "Fee Policy",
			"39": "Donor Deceased Policy",
			"40": "Refund Policy",
			"41": "Forms / Guidance",
			"42": "Digital Product",
			"43": "Safeguarding Policy",
			"44": "Jurisdiction",
			"45": "Other",
		},
	},
}

func getValidSubcategory(category string, subCategories []string) string {
	if category, ok := complaintCategories[category]; ok {
		for _, subCategory := range subCategories {
			s := strings.TrimSpace(subCategory)
			if _, ok := category.Children[s]; ok {
				return s
			}
		}
	}

	return ""
}

type addComplaintData struct {
	XSRFToken string
	Entity    string
	Success   bool
	Error     sirius.ValidationError

	Categories            map[string]complaintCategory
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
			XSRFToken:  ctx.XSRFToken,
			Categories: complaintCategories,
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
				SubCategory:          getValidSubcategory(postFormString(r, "category"), r.PostForm["subCategory"]),
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
