package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockChangeCaseStatusClient struct {
	mock.Mock
}

func (m *mockChangeCaseStatusClient) CaseSummary(ctx sirius.Context, uid string) (sirius.CaseSummary, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(sirius.CaseSummary), args.Error(1)
}

func (m *mockChangeCaseStatusClient) EditDigitalLPAStatus(ctx sirius.Context, caseUID string, caseStatusData sirius.CaseStatusData) error {
	args := m.Called(ctx, caseUID, caseStatusData)
	return args.Error(0)
}

func (m *mockChangeCaseStatusClient) RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error) {
	args := m.Called(ctx, category)
	if args.Get(0) != nil {
		return args.Get(0).([]sirius.RefDataItem), args.Error(1)
	}
	return nil, args.Error(1)
}

var statusItems = []statusItem{
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

func TestGetChangeCaseStatus(t *testing.T) {
	caseSummary := sirius.CaseSummary{
		DigitalLpa: sirius.DigitalLpa{
			UID: "M-9876-9876-9876",
			SiriusData: sirius.SiriusData{
				ID:      676,
				Subtype: "personal-welfare",
				Status:  "Draft",
			},
			LpaStoreData: sirius.LpaStoreData{
				Status: "draft",
			},
		},
	}

	statusChangeReasons := []sirius.RefDataItem{
		{
			Handle:        "LPA_DOES_NOT_WORK",
			Label:         "The LPA does not work and cannot be changed",
			ParentSources: []string{"cannot-register"},
		},
	}

	client := &mockChangeCaseStatusClient{}
	client.
		On("CaseSummary", mock.Anything, "M-9876-9876-9876").
		Return(caseSummary, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.CaseStatusChangeReason).
		Return(statusChangeReasons, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, changeCaseStatusData{
			Entity:                  "personal-welfare M-9876-9876-9876",
			CaseUID:                 "M-9876-9876-9876",
			OldStatus:               "draft",
			StatusItems:             statusItems,
			CaseStatusChangeReasons: statusChangeReasons,
			Error:                   sirius.ValidationError{Field: sirius.FieldErrors{}},
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/change-case-status?uid=M-9876-9876-9876", nil)
	w := httptest.NewRecorder()

	err := ChangeCaseStatus(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostChangeCaseStatus(t *testing.T) {
	caseSummary := sirius.CaseSummary{
		DigitalLpa: sirius.DigitalLpa{
			UID: "M-9876-9876-9876",
			SiriusData: sirius.SiriusData{
				ID:      676,
				Subtype: "personal-welfare",
				Status:  "Draft",
			},
			LpaStoreData: sirius.LpaStoreData{
				Status: "draft",
			},
		},
	}

	statusChangeReasons := []sirius.RefDataItem{
		{
			Handle:        "LPA_DOES_NOT_WORK",
			Label:         "The LPA does not work and cannot be changed",
			ParentSources: []string{"cannot-register"},
		},
		{
			Handle:        "CANCELLED_BY_COP",
			Label:         "Cancelled by the Court of Protection",
			ParentSources: []string{"cancelled"},
		},
	}

	client := &mockChangeCaseStatusClient{}
	client.
		On("CaseSummary", mock.Anything, "M-9876-9876-9876").
		Return(caseSummary, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.CaseStatusChangeReason).
		Return(statusChangeReasons, nil)

	client.
		On("EditDigitalLPAStatus", mock.Anything, "M-9876-9876-9876", sirius.CaseStatusData{
			Status: "expired",
		}).
		Return(nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, changeCaseStatusData{
			Success:                 true,
			Entity:                  "personal-welfare M-9876-9876-9876",
			CaseUID:                 "M-9876-9876-9876",
			OldStatus:               "in-progress",
			NewStatus:               "expired",
			StatusItems:             statusItems,
			CaseStatusChangeReasons: statusChangeReasons,
			Error:                   sirius.ValidationError{Field: sirius.FieldErrors{}},
		}).
		Return(nil)

	form := url.Values{
		"status": {"expired"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/change-case-status?uid=M-9876-9876-9876", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := ChangeCaseStatus(client, template.Func)(w, r)
	resp := w.Result()

	assert.Equal(t, RedirectError("/lpa/M-9876-9876-9876"), err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostChangeCaseStatusWithReason(t *testing.T) {
	caseSummary := sirius.CaseSummary{
		DigitalLpa: sirius.DigitalLpa{
			UID: "M-9876-9876-9876",
			SiriusData: sirius.SiriusData{
				ID:      676,
				Subtype: "personal-welfare",
				Status:  "Draft",
			},
			LpaStoreData: sirius.LpaStoreData{
				Status: "draft",
			},
		},
	}

	statusChangeReasons := []sirius.RefDataItem{
		{
			Handle:        "LPA_DOES_NOT_WORK",
			Label:         "The LPA does not work and cannot be changed",
			ParentSources: []string{"cannot-register"},
		},
		{
			Handle:        "CANCELLED_BY_COP",
			Label:         "Cancelled by the Court of Protection",
			ParentSources: []string{"cancelled"},
		},
	}

	client := &mockChangeCaseStatusClient{}
	client.
		On("CaseSummary", mock.Anything, "M-9876-9876-9876").
		Return(caseSummary, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.CaseStatusChangeReason).
		Return(statusChangeReasons, nil)

	client.
		On("EditDigitalLPAStatus", mock.Anything, "M-9876-9876-9876", sirius.CaseStatusData{
			Status:           "cannot-register",
			CaseChangeReason: "LPA_DOES_NOT_WORK",
		}).
		Return(nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, changeCaseStatusData{
			Success:                 true,
			Entity:                  "personal-welfare M-9876-9876-9876",
			CaseUID:                 "M-9876-9876-9876",
			OldStatus:               "in-progress",
			NewStatus:               "cannot-register",
			StatusItems:             statusItems,
			CaseStatusChangeReasons: statusChangeReasons,
			StatusChangeReason:      "LPA_DOES_NOT_WORK",
			Error:                   sirius.ValidationError{Field: sirius.FieldErrors{}},
		}).
		Return(nil)

	form := url.Values{
		"status":       {"cannot-register"},
		"statusReason": {"LPA_DOES_NOT_WORK"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/change-case-status?uid=M-9876-9876-9876", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := ChangeCaseStatus(client, template.Func)(w, r)
	resp := w.Result()

	assert.Equal(t, RedirectError("/lpa/M-9876-9876-9876"), err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
