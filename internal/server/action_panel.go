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

type ActionPanelClient interface {
	CasesByDonor(ctx sirius.Context, id int) ([]sirius.Case, error)
	GetDraftCount(ctx sirius.Context, caseType string, caseId int) (sirius.DocumentDraftCount, error)
}

type ActionPanelData struct {
	XSRFToken          string
	DonorID            int
	SelectedCases      []sirius.Case
	CaseUids           string
	CaseType           string
	ActionPanelButtons []ActionPanelButton
}

func ActionPanel(client ActionPanelClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if err := r.ParseForm(); err != nil {
			return err
		}

		ctx := getContext(r)
		data := ActionPanelData{XSRFToken: ctx.XSRFToken}

		// Support both query-string usage (?donorId=123) and path-based params
		donorIDStr := r.URL.Query().Get("donorId")
		if donorIDStr == "" {
			donorIDStr = r.URL.Query().Get("id")
		}
		if donorIDStr != "" {
			if donorID, err := strconv.Atoi(donorIDStr); err == nil {
				data.DonorID = donorID
			}
		}

		caseUIDs := r.Form["uid[]"]
		data.CaseUids = buildUIDQueryString(caseUIDs)

		entityType, err := sirius.ParseEntityType(r.FormValue("entity"))
		if err != nil {
			return err
		}
		data.CaseType = string(entityType)

		group, groupCtx := errgroup.WithContext(ctx.Context)

		var draftCount int
		group.Go(func() error {
			if data.DonorID > 0 {
				cases, err := client.CasesByDonor(ctx.With(groupCtx), data.DonorID)
				if err != nil {
					return err
				}

				// Filter cases by uid[] parameter if provided
				if len(caseUIDs) > 0 {
					casesByUID := make(map[string]sirius.Case, len(cases))
					for _, c := range cases {
						casesByUID[c.UID] = c
					}

					var filtered []sirius.Case
					for _, uid := range caseUIDs {
						if c, ok := casesByUID[uid]; ok {
							filtered = append(filtered, c)
						}
					}
					data.SelectedCases = filtered
				} else {
					data.SelectedCases = cases
				}
			}

			if len(data.SelectedCases) == 1 {
				documentDraftCount, err := client.GetDraftCount(ctx.With(groupCtx), strings.ToLower(data.SelectedCases[0].CaseType), data.SelectedCases[0].ID)
				if err != nil {
					return err
				}
				draftCount = documentDraftCount.DraftCount
			}
			return nil
		})

		if err := group.Wait(); err != nil {
			return err
		}

		data.ActionPanelButtons = GetActionPanelButtons(data.SelectedCases, data.DonorID, data.CaseUids, draftCount > 0)

		return tmpl(w, data)
	}
}

type ActionPanelButton struct {
	Label    string
	URL      string
	IconName string
	Disabled bool
}

func GetActionPanelButtons(selectedCases []sirius.Case, donorId int, caseUids string, hasDrafts bool) []ActionPanelButton {
	warningUrl := fmt.Sprintf("/create-warning?id=%d&entity=person%s", donorId, caseUids)
	eventUrl := fmt.Sprintf("/create-event?id=%d&entity=person%s", donorId, caseUids)
	createDonorUrl := fmt.Sprintf("/create-donor?id=%d&entity=person%s", donorId, caseUids)
	editDonorUrl := fmt.Sprintf("/edit-donor?id=%d&entity=person%s", donorId, caseUids)
	complaintUrl := ""
	createDocumentUrl := ""
	editDocumentUrl := ""
	changeStatusUrl := ""
	paymentsUrl := ""
	newTaskUrl := ""
	editDatesUrl := ""
	allocateCasesUrl := ""

	if len(selectedCases) == 1 {
		selectedCase := selectedCases[0]
		caseType := strings.ToLower(selectedCase.CaseType)
		id := selectedCase.ID

		warningUrl = fmt.Sprintf("/create-warning?id=%d&entity=%s%s", donorId, caseType, caseUids)
		complaintUrl = fmt.Sprintf("/add-complaint?id=%d&case=%s", id, caseType)
		createDocumentUrl = fmt.Sprintf("/create-document?id=%d&case=%s", id, caseType)
		changeStatusUrl = fmt.Sprintf("/change-status?id=%d&case=%s&donorId=%d%s", id, caseType, donorId, caseUids)
		paymentsUrl = fmt.Sprintf("/payments/%d", id)
		newTaskUrl = fmt.Sprintf("/create-task?id=%d&entity=%s", id, caseType)
		editDatesUrl = fmt.Sprintf("/edit-dates?id=%d&case=%s", id, caseType)
		if hasDrafts {
			editDocumentUrl = fmt.Sprintf("/edit-document?id=%d&case=%s", id, caseType)
		}
	}
	if len(selectedCases) > 1 {
		idQuery := ""
		caseType := strings.ToLower(selectedCases[0].CaseType)
		for i, c := range selectedCases {
			if i == 0 {
				idQuery += fmt.Sprintf("id=%d", c.ID)
			} else {
				idQuery += fmt.Sprintf("&id=%d", c.ID)
			}
		}
		allocateCasesUrl = fmt.Sprintf("/allocate-cases?%s&entity=%s", idQuery, caseType)
	}

	return []ActionPanelButton{
		{
			Label:    "Create warning",
			URL:      warningUrl,
			IconName: "aw-create-warning",
			Disabled: false,
		},
		{
			Label:    "Create event",
			URL:      eventUrl,
			IconName: "aw-new-event",
			Disabled: false,
		},
		{
			Label:    "Add complaint",
			URL:      complaintUrl,
			IconName: "aw-log-complaint",
			Disabled: len(selectedCases) != 1,
		},
		{
			Label:    "Create document",
			URL:      createDocumentUrl,
			IconName: "aw-new-template",
			Disabled: len(selectedCases) != 1,
		},
		{
			Label:    "Retrieve draft",
			URL:      editDocumentUrl,
			IconName: "aw-new-template",
			Disabled: len(selectedCases) != 1 || !hasDrafts,
		},
		{
			Label:    "Change status",
			URL:      changeStatusUrl,
			IconName: "aw-change-status",
			Disabled: len(selectedCases) != 1,
		},
		{
			Label:    "Fees",
			URL:      paymentsUrl,
			IconName: "aw-fees",
			Disabled: len(selectedCases) != 1,
		},
		{
			Label:    "New task",
			URL:      newTaskUrl,
			IconName: "aw-new-task",
			Disabled: len(selectedCases) != 1,
		},
		{
			Label:    "Create donor",
			URL:      createDonorUrl,
			IconName: "aw-create-person",
			Disabled: false,
		},
		{
			Label:    "Edit donor",
			URL:      editDonorUrl,
			IconName: "aw-edit-person",
			Disabled: false,
		},
		{
			Label:    "Edit dates",
			URL:      editDatesUrl,
			IconName: "calendar-open",
			Disabled: len(selectedCases) != 1,
		},
		{
			Label:    "Allocate Case",
			URL:      allocateCasesUrl,
			IconName: "aw-allocate-case",
			Disabled: len(selectedCases) == 0,
		},
	}
}
