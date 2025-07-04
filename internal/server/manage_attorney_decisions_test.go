package server

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockManageAttorneyDecisionsClient struct {
	mock.Mock
}

func (m *mockManageAttorneyDecisionsClient) CaseSummary(ctx sirius.Context, uid string) (sirius.CaseSummary, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(sirius.CaseSummary), args.Error(1)
}

func (m *mockManageAttorneyDecisionsClient) ManageAttorneyDecisions(ctx sirius.Context, caseUID string, attorneyDecisions []sirius.AttorneyDecisions) error {
	args := m.Called(ctx, caseUID, attorneyDecisions)
	return args.Error(0)
}

var manageAttorneyDecisionsSummary = sirius.CaseSummary{
	DigitalLpa: sirius.DigitalLpa{
		UID: "M-6666-6666-6666",
		LpaStoreData: sirius.LpaStoreData{
			Attorneys: []sirius.LpaStoreAttorney{
				{
					LpaStorePerson: sirius.LpaStorePerson{
						Uid:        "302b05c7-896c-4290-904e-2005e4f1e81e",
						FirstNames: "Jack",
						LastName:   "Black",
						Address: sirius.LpaStoreAddress{
							Line1:    "9 Mount Pleasant Drive",
							Town:     "East Harling",
							Postcode: "NR16 2GB",
							Country:  "UK",
						},
					},
					DateOfBirth:     "1990-02-22",
					Status:          shared.ActiveAttorneyStatus.String(),
					AppointmentType: shared.OriginalAppointmentType.String(),
					Email:           "a@example.com",
					Mobile:          "077577575757",
					SignedAt:        "2024-01-12T10:09:09Z",
				},
				{
					LpaStorePerson: sirius.LpaStorePerson{
						Uid:        "987a01b1-456d-4567-813d-2010d3e2d72d",
						FirstNames: "Shelley",
						LastName:   "Jones",
						Address: sirius.LpaStoreAddress{
							Line1:    "29 Broad Road",
							Town:     "Birmingham",
							Postcode: "B29 6BT",
							Country:  "UK",
						},
					},
					DateOfBirth:     "1990-02-27",
					Status:          shared.ActiveAttorneyStatus.String(),
					AppointmentType: shared.ReplacementAppointmentType.String(),
					Email:           "b@example.com",
					Mobile:          "07122121242",
					SignedAt:        "2024-11-28T20:22:11Z",
				},
				{
					LpaStorePerson: sirius.LpaStorePerson{
						Uid:        "123a01b1-456d-5391-813d-2010d3e2d72d",
						FirstNames: "Jack",
						LastName:   "White",
						Address: sirius.LpaStoreAddress{
							Line1:    "29 Grange Road",
							Town:     "Birmingham",
							Postcode: "B29 6BL",
							Country:  "UK",
						},
					},
					DateOfBirth:     "1990-02-22",
					Status:          shared.InactiveAttorneyStatus.String(),
					AppointmentType: shared.ReplacementAppointmentType.String(),
					Email:           "c@example.com",
					Mobile:          "07122121212",
					SignedAt:        "2024-11-28T19:22:11Z",
				},
				{
					LpaStorePerson: sirius.LpaStorePerson{
						Uid:        "123a01b1-456d-5391-813d-2010d3e256f",
						FirstNames: "Jack",
						LastName:   "Green",
						Address: sirius.LpaStoreAddress{
							Line1:    "39 Grange Road",
							Town:     "Birmingham",
							Postcode: "B29 6BL",
							Country:  "UK",
						},
					},
					DateOfBirth:     "1990-02-26",
					Status:          shared.InactiveAttorneyStatus.String(),
					AppointmentType: shared.ReplacementAppointmentType.String(),
					Email:           "d@example.com",
					Mobile:          "07122121232",
					SignedAt:        "2024-11-30T19:22:11Z",
				},
				{
					LpaStorePerson: sirius.LpaStorePerson{
						Uid:        "638f049f-c01f-4ab2-973a-2ea763b3cf7a",
						FirstNames: "Consuelo",
						LastName:   "Swaniawski",
						Address: sirius.LpaStoreAddress{
							Line1:    "14 Meadow Close",
							Town:     "Kutch Court",
							Postcode: "AT28 7WM",
							Country:  "UK",
						},
					},
					DateOfBirth:     "1990-04-15",
					Status:          shared.RemovedAttorneyStatus.String(),
					AppointmentType: shared.OriginalAppointmentType.String(),
					Email:           "Consuelo.Swaniawski@example.com",
					Mobile:          "07004369909",
					SignedAt:        "2024-10-21T13:42:16Z",
				},
			},
		},
	},
}

func TestGetManageAttorneyDecisions(t *testing.T) {
	client := &mockManageAttorneyDecisionsClient{}
	client.
		On("CaseSummary", mock.Anything, "M-6666-6666-6666").
		Return(manageAttorneyDecisionsSummary, nil)

	formTemplate := &mockTemplate{}
	formTemplate.
		On("Func", mock.Anything, manageAttorneyDecisionsData{
			CaseSummary:     manageAttorneyDecisionsSummary,
			ActiveAttorneys: activeAttorneys,
			Error:           sirius.ValidationError{Field: sirius.FieldErrors{}},
		}).
		Return(nil)

	confirmTemplate := &mockTemplate{}

	server := newMockServer("/lpa/{uid}/manage-attorney-decisions", AttorneyDecisions(client, formTemplate.Func, confirmTemplate.Func))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-6666-6666-6666/manage-attorney-decisions", nil)
	resp, err := server.serve(req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	mock.AssertExpectationsForObjects(t, client, formTemplate)
}

func TestGetManageAttorneyDecisionsCaseSummaryFails(t *testing.T) {
	client := &mockManageAttorneyDecisionsClient{}
	client.
		On("CaseSummary", mock.Anything, "M-6666-6666-6666").
		Return(sirius.CaseSummary{}, errExample)

	formTemplate := &mockTemplate{}
	formTemplate.
		On("Func", mock.Anything, manageAttorneyDecisionsData{}).
		Return(nil)

	confirmTemplate := &mockTemplate{}

	server := newMockServer("/lpa/{uid}/manage-attorney-decisions", AttorneyDecisions(client, formTemplate.Func, confirmTemplate.Func))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-6666-6666-6666/manage-attorney-decisions", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
}

func TestGetManageAttorneyDecisionsTemplateErrors(t *testing.T) {
	client := &mockManageAttorneyDecisionsClient{}
	client.
		On("CaseSummary", mock.Anything, "M-6666-6666-6666").
		Return(manageAttorneyDecisionsSummary, nil)

	formTemplate := &mockTemplate{}
	formTemplate.
		On("Func", mock.Anything, manageAttorneyDecisionsData{
			CaseSummary:     manageAttorneyDecisionsSummary,
			ActiveAttorneys: activeAttorneys,
			Error:           sirius.ValidationError{Field: sirius.FieldErrors{}},
		}).
		Return(errExample)

	confirmTemplate := &mockTemplate{}

	server := newMockServer("/lpa/{uid}/manage-attorney-decisions", AttorneyDecisions(client, formTemplate.Func, confirmTemplate.Func))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-6666-6666-6666/manage-attorney-decisions", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
}

func TestPostManageAttorneyDecisionsInvalidData(t *testing.T) {
	tests := []struct {
		name     string
		form     url.Values
		formData formManageAttorneyDecisions
	}{
		{
			name: "neither an attorney nor all attorneys to make join decisions are selected",
			form: url.Values{},
			formData: formManageAttorneyDecisions{
				DecisionAttorneysUids: nil,
				SkipDecisionAttorney:  "",
			},
		},
		{
			name: "both an attorney and all attorneys to make join decisions are selected",
			form: url.Values{
				"decisionAttorney":     {"302b05c7-896c-4290-904e-2005e4f1e81e"},
				"skipDecisionAttorney": {"yes"},
			},
			formData: formManageAttorneyDecisions{
				DecisionAttorneysUids: []string{"302b05c7-896c-4290-904e-2005e4f1e81e"},
				SkipDecisionAttorney:  "yes",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockManageAttorneyDecisionsClient{}
			client.
				On("CaseSummary", mock.Anything, "M-6666-6666-6666").
				Return(manageAttorneyDecisionsSummary, nil)

			formTemplate := &mockTemplate{}
			formTemplate.
				On("Func", mock.Anything, manageAttorneyDecisionsData{
					CaseSummary:     manageAttorneyDecisionsSummary,
					ActiveAttorneys: activeAttorneys,
					Form:            tc.formData,
					Error: sirius.ValidationError{Field: sirius.FieldErrors{
						"decisionAttorney": {"reason": "Select who cannot make joint decisions, or select 'Joint decisions can be made by all attorneys'"},
					}},
				}).
				Return(nil)

			confirmTemplate := &mockTemplate{}

			server := newMockServer("/lpa/{uid}/manage-attorney-decisions", AttorneyDecisions(client, formTemplate.Func, confirmTemplate.Func))

			req, _ := http.NewRequest(http.MethodPost, "/lpa/M-6666-6666-6666/manage-attorney-decisions", strings.NewReader(tc.form.Encode()))
			req.Header.Add("Content-Type", formUrlEncoded)

			resp, err := server.serve(req)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.Code)
			mock.AssertExpectationsForObjects(t, client, formTemplate)
		})
	}
}

func TestPostManageAttorneyDecisionsValidData(t *testing.T) {
	client := &mockManageAttorneyDecisionsClient{}
	client.
		On("CaseSummary", mock.Anything, "M-6666-6666-6666").
		Return(manageAttorneyDecisionsSummary, nil)

	client.
		On("ManageAttorneyDecisions", mock.Anything, "M-6666-6666-6666", []sirius.AttorneyDecisions{
			{UID: "302b05c7-896c-4290-904e-2005e4f1e81e", CannotMakeJointDecisions: true},
			{UID: "987a01b1-456d-4567-813d-2010d3e2d72d", CannotMakeJointDecisions: false},
		}).
		Return(nil)

	formTemplate := &mockTemplate{}

	confirmTemplate := &mockTemplate{}
	confirmTemplate.
		On("Func", mock.Anything, manageAttorneyDecisionsData{
			CaseSummary:     manageAttorneyDecisionsSummary,
			ActiveAttorneys: activeAttorneys,
			Form: formManageAttorneyDecisions{
				DecisionAttorneysUids: []string{"302b05c7-896c-4290-904e-2005e4f1e81e"},
			},
			DecisionAttorneysDetails: []AttorneyDetails{
				{
					AttorneyName:    "Jack Black",
					AttorneyDob:     "1990-02-22",
					AppointmentType: "original",
				},
			},
			Error:   sirius.ValidationError{Field: sirius.FieldErrors{}},
			Success: false,
		}).
		Return(nil)

	server := newMockServer("/lpa/{uid}/manage-attorney-decisions", AttorneyDecisions(client, formTemplate.Func, confirmTemplate.Func))

	form := url.Values{
		"decisionAttorney": {"302b05c7-896c-4290-904e-2005e4f1e81e"},
	}

	req, _ := http.NewRequest(http.MethodPost, "/lpa/M-6666-6666-6666/manage-attorney-decisions", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", formUrlEncoded)
	resp, err := server.serve(req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Nil(t, err)
}

func TestPostConfirmManageAttorneyDecisionsValidData(t *testing.T) {
	tests := []struct {
		name              string
		form              url.Values
		expectedDecisions []sirius.AttorneyDecisions
	}{
		{
			name: "single attorney cannot make joint decisions",
			form: url.Values{
				"decisionAttorney": {"302b05c7-896c-4290-904e-2005e4f1e81e"},
				"confirmDecisions": {""},
			},
			expectedDecisions: []sirius.AttorneyDecisions{
				{UID: "302b05c7-896c-4290-904e-2005e4f1e81e", CannotMakeJointDecisions: true},
				{UID: "987a01b1-456d-4567-813d-2010d3e2d72d", CannotMakeJointDecisions: false},
			},
		},
		{
			name: "all attorneys can make joint decisions",
			form: url.Values{
				"skipDecisionAttorney": {"yes"},
				"confirmDecisions":     {""},
			},
			expectedDecisions: []sirius.AttorneyDecisions{
				{UID: "302b05c7-896c-4290-904e-2005e4f1e81e", CannotMakeJointDecisions: false},
				{UID: "987a01b1-456d-4567-813d-2010d3e2d72d", CannotMakeJointDecisions: false},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockManageAttorneyDecisionsClient{}
			client.
				On("CaseSummary", mock.Anything, "M-6666-6666-6666").
				Return(manageAttorneyDecisionsSummary, nil)

			client.
				On("ManageAttorneyDecisions", mock.Anything, "M-6666-6666-6666", tc.expectedDecisions).
				Return(nil)

			formTemplate := &mockTemplate{}
			confirmTemplate := &mockTemplate{}

			server := newMockServer("/lpa/{uid}/manage-attorney-decisions", AttorneyDecisions(client, formTemplate.Func, confirmTemplate.Func))

			req, _ := http.NewRequest(http.MethodPost, "/lpa/M-6666-6666-6666/manage-attorney-decisions", strings.NewReader(tc.form.Encode()))
			req.Header.Add("Content-Type", formUrlEncoded)

			resp, err := server.serve(req)

			assert.Equal(t, http.StatusOK, resp.Code)
			assert.Equal(t, RedirectError("/lpa/M-6666-6666-6666"), err)

			mock.AssertExpectationsForObjects(t, client)
		})
	}
}
