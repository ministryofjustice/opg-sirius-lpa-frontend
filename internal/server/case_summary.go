package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type GetCaseSummaryClient interface {
	Case(ctx sirius.Context, id int) (sirius.Case, error)
	WarningsForCase(ctx sirius.Context, caseId int) ([]sirius.Warning, error)
	//TasksForCase(ctx sirius.Context, caseId int) ([]sirius.Task, error)
	TasksForDonor(ctx sirius.Context, donorId int) ([]sirius.Task, error)
	Complaints(ctx sirius.Context, caseType string, caseId int) ([]sirius.Complaint, error)
	Investigations(ctx sirius.Context, caseType string, caseId int) ([]sirius.Investigation, error)
}

type caseSummaryData struct {
	XSRFToken           string
	Error               sirius.ValidationError
	CasesWarnings       []sirius.Warning
	CasesTasks          []sirius.Task
	CasesComplaints     []sirius.Complaint
	CasesInvestigations []sirius.Investigation
}

func GetCaseSummary(client GetCaseSummaryClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if err := r.ParseForm(); err != nil {
			return err
		}
		donorId, err := strconv.Atoi(r.PathValue("donorId"))
		if err != nil {
			return err
		}
		caseUIDs := r.Form["id[]"]
		ctx := getContext(r)

		var casesWarnings []sirius.Warning
		//var casesTasks []sirius.Task
		var casesComplaints []sirius.Complaint
		var casesInvestigations []sirius.Investigation

		for _, id := range caseUIDs {
			caseidToInt, _ := strconv.Atoi(id)
			caseitem, err := client.Case(ctx, caseidToInt)
			if err != nil {
				return err
			}

			caseWarning, err := client.WarningsForCase(ctx, caseidToInt)
			if err != nil {
				return err
			}

			//caseTasks, err := client.TasksForCase(ctx, caseidToInt)
			//if err != nil {
			//	return err
			//}

			caseComplaints, err := client.Complaints(ctx, strings.ToLower(caseitem.CaseType), caseidToInt)
			if err != nil {
				return err
			}

			caseInvestigations, err := client.Investigations(ctx, strings.ToLower(caseitem.CaseType), caseidToInt)
			if err != nil {
				return err
			}

			casesWarnings = append(casesWarnings, caseWarning...)
			//casesTasks = append(casesTasks, caseTasks...)
			casesComplaints = append(casesComplaints, caseComplaints...)
			casesInvestigations = append(casesInvestigations, caseInvestigations...)
		}

		donorTasks, err := client.TasksForDonor(ctx, donorId)
		if err != nil {
			fmt.Println("Error getting tasks for donor: ", err)
			return err
		}

		data := caseSummaryData{
			XSRFToken:           ctx.XSRFToken,
			CasesWarnings:       casesWarnings,
			CasesTasks:          donorTasks,
			CasesComplaints:     casesComplaints,
			CasesInvestigations: casesInvestigations,
		}

		return tmpl(w, data)
	}
}
