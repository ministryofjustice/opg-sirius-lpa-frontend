package server

import (
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type mockManageRestrictionsClient struct {
	mock.Mock
}

func (m *mockManageRestrictionsClient) CaseSummary(ctx sirius.Context, uid string) (sirius.CaseSummary, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(sirius.CaseSummary), args.Error(1)
}

func (m *mockManageRestrictionsClient) ClearTask(ctx sirius.Context, taskID int) error {
	args := m.Called(ctx, taskID)
	return args.Error(0)
}

func (m *mockManageRestrictionsClient) UpdateSeveranceStatus(ctx sirius.Context, caseUID string, severanceStatusData sirius.SeveranceStatusData) error {
	args := m.Called(ctx, caseUID, severanceStatusData)
	return args.Error(0)
}

func (m *mockManageRestrictionsClient) EditSeveranceApplication(ctx sirius.Context, caseUID string, severanceApplicationDetails sirius.SeveranceApplication) error {
	args := m.Called(ctx, caseUID, severanceApplicationDetails)
	return args.Error(0)
}

func boolPointer(b bool) *bool {
	return &b
}

var restrictionsCaseSummary = sirius.CaseSummary{
	DigitalLpa: sirius.DigitalLpa{
		UID: "M-1111-2222-3333",
	},
	TaskList: []sirius.Task{
		{
			ID:     1,
			Name:   "Review restrictions and conditions",
			Status: "Not started",
		},
	},
}

var restrictionsCaseSummaryWithSeveranceRequired = sirius.CaseSummary{
	DigitalLpa: sirius.DigitalLpa{
		UID: "M-1111-2222-3333",
		SiriusData: sirius.SiriusData{
			Application: sirius.Draft{
				SeveranceStatus: "REQUIRED",
			},
		},
	},
	TaskList: []sirius.Task{
		{
			ID:     1,
			Name:   "Review restrictions and conditions",
			Status: "Not started",
		},
	},
}

var restrictionsCaseSummaryWithDonorConsentGiven = sirius.CaseSummary{
	DigitalLpa: sirius.DigitalLpa{
		UID: "M-1111-2222-3333",
		SiriusData: sirius.SiriusData{
			Application: sirius.Draft{
				SeveranceStatus: "REQUIRED",
				SeveranceApplication: &sirius.SeveranceApplication{
					HasDonorConsented: boolPointer(true),
				},
			},
		},
	},
	TaskList: []sirius.Task{
		{
			ID:     1,
			Name:   "Review restrictions and conditions",
			Status: "Not started",
		},
	},
}

func TestGetManageRestrictionsCases(t *testing.T) {
	tests := []struct {
		name          string
		caseSummary   sirius.CaseSummary
		action        string
		templateError error
		expectedError error
	}{
		{
			name:          "Get manage restrictions request succeeds",
			caseSummary:   restrictionsCaseSummary,
			action:        "",
			templateError: nil,
			expectedError: nil,
		},
		{
			name:          "Get manage restrictions with severance required request succeeds",
			caseSummary:   restrictionsCaseSummaryWithSeveranceRequired,
			action:        "donor-consent",
			templateError: nil,
			expectedError: nil,
		},
		{
			name:          "Get manage restrictions with donor consent given succeeds",
			caseSummary:   restrictionsCaseSummaryWithDonorConsentGiven,
			action:        "court-order",
			templateError: nil,
			expectedError: nil,
		},
		{
			name:          "Get case summary errors",
			caseSummary:   sirius.CaseSummary{},
			action:        "",
			templateError: nil,
			expectedError: errExample,
		},
		{
			name:          "Template errors",
			caseSummary:   restrictionsCaseSummary,
			action:        "",
			templateError: errExample,
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockManageRestrictionsClient{}
			client.
				On("CaseSummary", mock.Anything, "M-1111-2222-3333").
				Return(tc.caseSummary, tc.expectedError)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, manageRestrictionsData{
					CaseSummary: tc.caseSummary,
					CaseUID:     "M-1111-2222-3333",
					FormAction:  tc.action,
					Error:       sirius.ValidationError{Field: sirius.FieldErrors{}},
				}).
				Return(tc.templateError)

			server := newMockServer("/lpa/{uid}/manage-restrictions", ManageRestrictions(client, template.Func))

			req, _ := http.NewRequest(http.MethodGet, "/lpa/M-1111-2222-3333/manage-restrictions", nil)
			resp, err := server.serve(req)

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
			} else if tc.templateError != nil {
				assert.Equal(t, tc.templateError, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, resp.Code)
				mock.AssertExpectationsForObjects(t, client, template)
			}
		})
	}
}

func TestPostManageRestrictions(t *testing.T) {
	tests := []struct {
		name            string
		form            url.Values
		expectedData    manageRestrictionsData
		severanceAction string
		error           sirius.ValidationError
	}{
		{
			name:            "No option selected",
			form:            url.Values{},
			severanceAction: "",
			error: sirius.ValidationError{Field: sirius.FieldErrors{
				"severanceAction": {"reason": "Please select an option"},
			}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockManageRestrictionsClient{}
			client.
				On("CaseSummary", mock.Anything, "M-1111-2222-3333").
				Return(restrictionsCaseSummary, nil)

			restrictionsData := manageRestrictionsData{
				SeveranceAction: tc.severanceAction,
				CaseSummary:     restrictionsCaseSummary,
				CaseUID:         "M-1111-2222-3333",
				FormAction:      "",
				Error:           tc.error,
			}

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, restrictionsData).
				Return(nil)

			server := newMockServer("/lpa/{uid}/manage-restrictions", ManageRestrictions(client, template.Func))

			req, _ := http.NewRequest(http.MethodPost, "/lpa/M-1111-2222-3333/manage-restrictions", strings.NewReader(tc.form.Encode()))
			req.Header.Add("Content-Type", formUrlEncoded)
			resp, err := server.serve(req)

			assert.Nil(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.Code)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostManageRestrictionsRedirects(t *testing.T) {
	tests := []struct {
		name            string
		severanceAction string
		severanceStatus *sirius.SeveranceStatusData
	}{
		{
			name:            "Severance application not required",
			severanceAction: "severance-application-not-required",
			severanceStatus: &sirius.SeveranceStatusData{SeveranceStatus: "NOT_REQUIRED"},
		},
		{
			name:            "Severance application required",
			severanceAction: "severance-application-required",
			severanceStatus: &sirius.SeveranceStatusData{SeveranceStatus: "REQUIRED"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockManageRestrictionsClient{}
			client.
				On("CaseSummary", mock.Anything, "M-1111-2222-3333").
				Return(restrictionsCaseSummary, nil)

			if tc.severanceAction == "severance-application-not-required" {
				client.On("ClearTask", mock.Anything, 1).Return(nil)
			}

			client.
				On("UpdateSeveranceStatus", mock.Anything, "M-1111-2222-3333", *tc.severanceStatus).
				Return(nil)

			template := &mockTemplate{}
			server := newMockServer("/lpa/{uid}/manage-restrictions", ManageRestrictions(client, template.Func))

			form := url.Values{
				"severanceAction": {tc.severanceAction},
			}

			req, _ := http.NewRequest(http.MethodPost, "/lpa/M-1111-2222-3333/manage-restrictions", strings.NewReader(form.Encode()))
			req.Header.Add("Content-Type", formUrlEncoded)
			_, err := server.serve(req)

			assert.Equal(t, RedirectError("/lpa/M-1111-2222-3333/lpa-details"), err)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostManageRestrictionsWithSeveranceRequiredRedirects(t *testing.T) {
	tests := []struct {
		name               string
		donorConsentAction string
		severanceDetails   *sirius.SeveranceApplication
	}{
		{
			name:               "Donor consent given",
			donorConsentAction: "donor-consent-given",
			severanceDetails:   &sirius.SeveranceApplication{HasDonorConsented: boolPointer(true)},
		},
		{
			name:               "Donor refused severance",
			donorConsentAction: "donor-consent-not-given",
			severanceDetails:   &sirius.SeveranceApplication{HasDonorConsented: boolPointer(false)},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockManageRestrictionsClient{}
			client.
				On("CaseSummary", mock.Anything, "M-1111-2222-3333").
				Return(restrictionsCaseSummaryWithSeveranceRequired, nil)

			client.
				On("EditSeveranceApplication", mock.Anything, "M-1111-2222-3333", *tc.severanceDetails).
				Return(nil)

			template := &mockTemplate{}
			server := newMockServer("/lpa/{uid}/manage-restrictions", ManageRestrictions(client, template.Func))

			form := url.Values{
				"donorConsentGiven": {tc.donorConsentAction},
				"action":            {"donor-consent"},
			}

			req, _ := http.NewRequest(http.MethodPost, "/lpa/M-1111-2222-3333/manage-restrictions", strings.NewReader(form.Encode()))
			req.Header.Add("Content-Type", formUrlEncoded)
			_, err := server.serve(req)

			assert.Equal(t, RedirectError("/lpa/M-1111-2222-3333/lpa-details"), err)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostManageRestrictionsWithDonorConsentGivenRedirects(t *testing.T) {
	tests := []struct {
		name                   string
		courtOrderDecisionMade string
		courtOrderReceived     string
		severanceOrderedAction string
		severanceDetails       *sirius.SeveranceApplication
	}{
		{
			name:                   "Court order decision made date given",
			courtOrderDecisionMade: "2025-04-05",
			courtOrderReceived:     "",
			severanceOrderedAction: "",
			severanceDetails: &sirius.SeveranceApplication{
				CourtOrderDecisionMade: "2025-04-05",
			},
		},
		{
			name:                   "Court order received date and severance ordered given",
			courtOrderDecisionMade: "2025-04-05",
			courtOrderReceived:     "2025-04-10",
			severanceOrderedAction: "severance-ordered",
			severanceDetails: &sirius.SeveranceApplication{
				CourtOrderDecisionMade: "2025-04-05",
				CourtOrderReceived:     "2025-04-10",
				SeveranceOrdered:       boolPointer(true),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockManageRestrictionsClient{}
			client.
				On("CaseSummary", mock.Anything, "M-1111-2222-3333").
				Return(restrictionsCaseSummaryWithDonorConsentGiven, nil)

			client.
				On("EditSeveranceApplication", mock.Anything, "M-1111-2222-3333", *tc.severanceDetails).
				Return(nil)

			template := &mockTemplate{}
			server := newMockServer("/lpa/{uid}/manage-restrictions", ManageRestrictions(client, template.Func))

			form := url.Values{
				"courtOrderDecisionMade": {tc.courtOrderDecisionMade},
				"courtOrderReceived":     {tc.courtOrderReceived},
				"severanceOrdered":       {tc.severanceOrderedAction},
				"action":                 {"court-order"},
			}

			req, _ := http.NewRequest(http.MethodPost, "/lpa/M-1111-2222-3333/manage-restrictions", strings.NewReader(form.Encode()))
			req.Header.Add("Content-Type", formUrlEncoded)
			_, err := server.serve(req)

			assert.Equal(t, RedirectError("/lpa/M-1111-2222-3333/lpa-details"), err)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}
