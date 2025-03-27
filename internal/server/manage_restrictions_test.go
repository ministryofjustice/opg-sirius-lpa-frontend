package server

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestGetManageRestrictionsCases(t *testing.T) {
	tests := []struct {
		name          string
		caseSummary   sirius.CaseSummary
		templateError error
		expectedError error
	}{
		{
			name:          "Get manage restrictions request succeeds",
			caseSummary:   restrictionsCaseSummary,
			templateError: nil,
			expectedError: nil,
		},
		{
			name:          "Get case summary errors",
			caseSummary:   sirius.CaseSummary{},
			templateError: nil,
			expectedError: errExample,
		},
		{
			name:          "Template errors",
			caseSummary:   restrictionsCaseSummary,
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
