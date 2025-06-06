package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/templatefn"
	"golang.org/x/sync/errgroup"
	"net/http"
)

type ChangeCaseStatusClient interface {
	RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error)
	CaseSummary(sirius.Context, string) (sirius.CaseSummary, error)
	EditDigitalLPAStatus(sirius.Context, string, sirius.CaseStatusData) error
}

type statusItem struct {
	Value           string
	Label           string
	ConditionalItem bool
}

type changeCaseStatusData struct {
	XSRFToken string
	Entity    string
	CaseUID   string
	Success   bool
	Error     sirius.ValidationError

	StatusItems             []statusItem
	CaseStatusChangeReasons []sirius.RefDataItem
	OldStatus               string
	NewStatus               string
	StatusChangeReason      string
}

func ChangeCaseStatus(client ChangeCaseStatusClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		caseUID := r.FormValue("uid")

		ctx := getContext(r)

		var cs sirius.CaseSummary
		var caseStatusChangeReasons []sirius.RefDataItem
		var err error

		group, groupCtx := errgroup.WithContext(ctx.Context)

		group.Go(func() error {
			cs, err = client.CaseSummary(ctx.With(groupCtx), caseUID)
			if err != nil {
				return err
			}
			return nil
		})

		group.Go(func() error {
			caseStatusChangeReasons, err = client.RefDataByCategory(ctx.With(groupCtx), sirius.CaseStatusChangeReason)
			if err != nil {
				return err
			}

			return nil
		})

		if err := group.Wait(); err != nil {
			return err
		}

		status := "draft"

		if cs.DigitalLpa.LpaStoreData.Status != "" {
			status = cs.DigitalLpa.LpaStoreData.Status
		}

		data := changeCaseStatusData{
			XSRFToken:               ctx.XSRFToken,
			Error:                   sirius.ValidationError{Field: sirius.FieldErrors{}},
			Entity:                  fmt.Sprintf("%s %s", cs.DigitalLpa.SiriusData.Subtype, caseUID),
			CaseUID:                 caseUID,
			OldStatus:               status,
			NewStatus:               postFormString(r, "status"),
			StatusChangeReason:      postFormString(r, "statusReason"),
			CaseStatusChangeReasons: caseStatusChangeReasons,
		}

		data.StatusItems = []statusItem{
			{Value: "draft", Label: "Draft", ConditionalItem: false},
			{Value: "in-progress", Label: "In progress", ConditionalItem: false},
			{Value: "statutory-waiting-period", Label: "Statutory waiting period", ConditionalItem: false},
			{Value: "registered", Label: "Registered", ConditionalItem: false},
			{Value: "suspended", Label: "Suspended", ConditionalItem: false},
			{Value: "do-not-register", Label: "Do not register", ConditionalItem: false},
			{Value: "expired", Label: "Expired", ConditionalItem: false},
			{Value: "cannot-register", Label: "Cannot register", ConditionalItem: true},
			{Value: "cancelled", Label: "Cancelled", ConditionalItem: true},
			{Value: "de-registered", Label: "De-registered", ConditionalItem: false},
		}

		if r.Method == http.MethodPost {
			if (data.NewStatus == "cannot-register" || data.NewStatus == "cancelled") && data.StatusChangeReason == "" {
				data.OldStatus = data.NewStatus
				w.WriteHeader(http.StatusBadRequest)
				data.Error.Field["changeReason"] = map[string]string{
					"reason": "Please select a reason",
				}
			}

			if !data.Error.Any() {
				caseStatusData := sirius.CaseStatusData{
					Status:           data.NewStatus,
					CaseChangeReason: data.StatusChangeReason,
				}

				err = client.EditDigitalLPAStatus(ctx, caseUID, caseStatusData)

				if ve, ok := err.(sirius.ValidationError); ok {
					w.WriteHeader(http.StatusBadRequest)
					data.Error = ve
				} else if err != nil {
					return err
				} else {
					data.Success = true
					data.OldStatus = data.NewStatus

					SetFlash(w, FlashNotification{
						Title: fmt.Sprintf("Status changed to %s", templatefn.StatusLabelFormat(data.NewStatus)),
					})
					return RedirectError(fmt.Sprintf("/lpa/%s", data.CaseUID))
				}
			}
		}

		return tmpl(w, data)
	}
}
