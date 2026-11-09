package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/sync/errgroup"
)

type GetLpaHistoryClient interface {
	RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error)
	GetEvents(ctx sirius.Context, donorId string, caseIds []string, sourceTypes []string, eventIds []string, sortBy string) (sirius.LpaEventsResponse, error)
	GetUserDetails(sirius.Context) (sirius.User, error)
}

type getLpaHistory struct {
	XSRFToken              string
	DonorID                string
	Events                 []sirius.LpaEvent
	EventFilterData        []sirius.SourceType
	Form                   FilterLpaEventsForm
	IsPartial              bool
	TotalEvents            int
	TotalFilteredEvents    int
	IsFiltered             bool
	FeeReductionTypes      []sirius.RefDataItem
	ComplaintCategories    []sirius.RefDataItem
	ComplaintSubcategories []sirius.RefDataItem
	ComplainantCategories  []sirius.RefDataItem
	ComplaintOrigins       []sirius.RefDataItem
	CompensationTypes      []sirius.RefDataItem
	DecisionTypes          []sirius.RefDataItem
	DonorFieldOrder        []string
	LpaFieldOrder          []string
	EpaFieldOrder          []string
	IsSysAdminUser         bool
}

type FilterLpaEventsForm struct {
	Types []string `form:"type"`
	Sort  string   `form:"sort"`
}

var donorFieldOrder = []string{
	"salutation",
	"firstname",
	"middlenames",
	"surname",
	"otherNames",
	"previousNames",
	"dob",
	"email",
	"correspondenceByPost",
	"correspondenceByPhone",
	"correspondenceByEmail",
	"correspondenceByWelsh",
}

var lpaFieldOrder = []string{
	"applicationType",
	"onlineLpaId",
	"caseAttorneySingular",
	"caseAttorneyJointly",
	"caseAttorneyJointlyAndSeverally",
	"caseAttorneyJointlyAndJointlyAndSeverally",
	"attorneyActDecisions",
	"lifeSustainingTreatment",
	"applicationHasRestrictions",
	"applicationHasGuidance",
	"lpaDonorSignatureDate",
	"certificateProviderSignatureDate",
	"applicantSignatureDate",
	"paymentByDebitCreditCard",
	"paymentByCheque",
	"paymentExemption",
	"paymentRemission",
	"haveAppliedForFeeRemission",
	"anyOtherInfo",
	"additionalInfo",
	"assignee",
	"cancellationDate",
	"registrationDate",
	"dispatchDate",
	"noticeGivenDate",
	"withdrawnDate",
}

var epaFieldOrder = []string{
	"caseAttorneyJointly",
	"caseAttorneySingular",
	"caseAttorneyJointlyAndSeverally",
	"epaDonorSignatureDate",
	"epaDonorNoticeGivenDate",
	"paymentByCheque",
	"paymentExemption",
	"paymentDate",
	"donorHasOtherEpas",
	"otherEpaInfo",
	"assignee",
	"cancellationDate",
	"dispatchDate",
	"dueDate",
	"filingDate",
	"invalidDate",
	"paymentDate",
	"receiptDate",
	"registrationDate",
	"rejectedDate",
	"revokedDate",
	"withdrawnDate",
}

type FieldChange struct {
	OldValue *string
	NewValue *string
}

func GetLpaHistory(client GetLpaHistoryClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		donorId := r.PathValue("donorId")
		caseIDs := r.URL.Query()["id[]"]

		ctx := getContext(r)
		group, groupCtx := errgroup.WithContext(ctx.Context)

		data := getLpaHistory{
			XSRFToken: ctx.XSRFToken,
			DonorID:   donorId,
			Form: FilterLpaEventsForm{
				Sort: "desc",
			},
			IsFiltered:      false,
			IsPartial:       ctx.IsPartial,
			DonorFieldOrder: donorFieldOrder,
			LpaFieldOrder:   lpaFieldOrder,
			EpaFieldOrder:   epaFieldOrder,
		}

		group.Go(func() error {
			user, err := client.GetUserDetails(ctx)
			if err != nil {
				return err
			}
			data.IsSysAdminUser = user.HasRole("System Admin")
			return nil
		})

		group.Go(func() error {
			eventsData, err := client.GetEvents(ctx.With(groupCtx), donorId, caseIDs, []string{}, []string{}, "desc")
			if err != nil {
				return err
			}

			normaliseComplaintTitleChanges(eventsData.Events)
			normaliseFeeDecisionChanges(eventsData.Events)

			data.Events = eventsData.Events
			data.EventFilterData = eventsData.Metadata.SourceTypes
			data.TotalEvents = eventsData.Total
			return nil
		})

		group.Go(func() error {
			feeReductionTypes, err := client.RefDataByCategory(ctx.With(groupCtx), sirius.FeeReductionTypeCategory)
			if err != nil {
				// If the call to get fee reduction types fails, we can just fall-back and display the database values instead of the labels
				data.FeeReductionTypes = []sirius.RefDataItem(nil)
				return nil
			}
			data.FeeReductionTypes = feeReductionTypes
			return nil
		})

		group.Go(func() error {
			complaintCategories, err := client.RefDataByCategory(ctx.With(groupCtx), sirius.ComplaintCategory)
			if err != nil {
				data.ComplaintCategories = []sirius.RefDataItem(nil)
				data.ComplaintSubcategories = []sirius.RefDataItem(nil)
				return nil
			}
			data.ComplaintCategories = complaintCategories

			for _, category := range complaintCategories {
				if len(category.Subcategories) != 0 {
					data.ComplaintSubcategories = append(data.ComplaintSubcategories, category.Subcategories...)
				}
			}
			return nil
		})

		group.Go(func() error {
			complainantCategories, err := client.RefDataByCategory(ctx.With(groupCtx), sirius.ComplainantCategory)
			if err != nil {
				data.ComplainantCategories = []sirius.RefDataItem(nil)
				return nil
			}
			data.ComplainantCategories = complainantCategories
			return nil
		})

		group.Go(func() error {
			complaintOrigins, err := client.RefDataByCategory(ctx.With(groupCtx), sirius.ComplaintOrigin)
			if err != nil {
				data.ComplaintOrigins = []sirius.RefDataItem(nil)
				return nil
			}
			data.ComplaintOrigins = complaintOrigins
			return nil
		})

		group.Go(func() error {
			compensationType, err := client.RefDataByCategory(ctx.With(groupCtx), sirius.CompensationType)
			if err != nil {
				data.CompensationTypes = []sirius.RefDataItem(nil)
				return nil
			}
			data.CompensationTypes = compensationType
			return nil
		})

		group.Go(func() error {
			decisionTypes, err := client.RefDataByCategory(ctx.With(groupCtx), sirius.FeeDecisionTypeCategory)
			if err != nil {
				data.DecisionTypes = []sirius.RefDataItem(nil)
				return nil
			}
			data.DecisionTypes = decisionTypes
			return nil
		})

		if err := group.Wait(); err != nil {
			return err
		}

		if r.Method == http.MethodPost {
			err := decoder.Decode(&data.Form, r.PostForm)
			if err != nil {
				return err
			}

			eventsData, err := client.GetEvents(ctx, donorId, caseIDs, data.Form.Types, []string{}, data.Form.Sort)
			if err != nil {
				return err
			}

			normaliseComplaintTitleChanges(eventsData.Events)
			normaliseFeeDecisionChanges(eventsData.Events)

			data.TotalFilteredEvents = eventsData.Total
			data.Events = eventsData.Events
			data.IsFiltered = true
		}

		return tmpl(w, data)
	}
}

func normaliseComplaintTitleChanges(events []sirius.LpaEvent) {
	for i, event := range events {
		if event.SourceType == shared.LpaEventSourceTypeComplaint {
			changes, isMap := event.Changes.(map[string]interface{})
			if isMap {
				title, hasTitle := changes["title"]
				if hasTitle {
					changes["title"] = normaliseChange(title)
				}
				events[i].Changes = changes
			}
		}
	}
}

func normaliseChange(v interface{}) FieldChange {
	var change FieldChange

	val, isList := v.([]interface{})
	if isList {
		if len(val) >= 1 {
			s := fmt.Sprintf("%v", val[0])
			change.OldValue = &s
		}
		if len(val) >= 2 {
			s := fmt.Sprintf("%v", val[1])
			change.NewValue = &s
		}
	}

	val2, isMap := v.(map[string]interface{})
	if isMap {
		newValue, hasNewValue := val2["1"]
		if hasNewValue {
			s := fmt.Sprintf("%v", newValue)
			change.NewValue = &s
		}
	}

	return change
}

var feeDecisionChangeFields = []string{"decisionType", "decisionDate", "decisionReason"}

func normaliseFeeDecisionChanges(events []sirius.LpaEvent) {
	for i, event := range events {
		if event.SourceType != shared.LpaEventSourceTypeFeeDecision || event.Type != "UPD" {
			continue
		}

		changes, isMap := event.Changes.(map[string]interface{})
		if !isMap {
			continue
		}

		for _, field := range feeDecisionChangeFields {
			value, hasField := changes[field]
			if hasField {
				changes[field] = normaliseFeeDecisionChange(value)
			}
		}

		events[i].Changes = changes
	}
}

func normaliseFeeDecisionChange(v interface{}) FieldChange {
	var change FieldChange

	list, isList := v.([]interface{})
	if isList {
		if len(list) >= 1 {
			change.OldValue = formatFeeDecisionValue(list[0])
		}
		if len(list) >= 2 {
			change.NewValue = formatFeeDecisionValue(list[1])
		}
	}

	keyed, isMap := v.(map[string]interface{})
	if isMap {
		oldValue, hasOldValue := keyed["0"]
		if hasOldValue {
			change.OldValue = formatFeeDecisionValue(oldValue)
		}
		newValue, hasNewValue := keyed["1"]
		if hasNewValue {
			change.NewValue = formatFeeDecisionValue(newValue)
		}
	}

	return change
}

func formatFeeDecisionValue(v interface{}) *string {
	switch value := v.(type) {
	case nil:
		return nil
	case string:
		if value == "" {
			return nil
		}
		return &value
	case map[string]interface{}:
		date, isDate := value["date"].(string)
		if isDate {
			parsed, err := time.Parse("2006-01-02 15:04:05.000000", date)
			if err != nil {
				return &date
			}
			formatted := parsed.Format("02/01/2006")
			return &formatted
		}
	}

	s := fmt.Sprintf("%v", v)
	return &s
}
