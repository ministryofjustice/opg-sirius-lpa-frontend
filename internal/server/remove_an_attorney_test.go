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

type mockRemoveAnAttorneyClient struct {
	mock.Mock
}

func (m *mockRemoveAnAttorneyClient) CaseSummary(ctx sirius.Context, uid string) (sirius.CaseSummary, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(sirius.CaseSummary), args.Error(1)
}

func (m *mockRemoveAnAttorneyClient) ChangeAttorneyStatus(ctx sirius.Context, caseUID string, attorneyUpdatedStatus []sirius.AttorneyUpdatedStatus) error {
	args := m.Called(ctx, caseUID, attorneyUpdatedStatus)
	return args.Error(0)
}

func (m *mockRemoveAnAttorneyClient) RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error) {
	args := m.Called(ctx, category)
	return args.Get(0).([]sirius.RefDataItem), args.Error(1)
}

func (m *mockRemoveAnAttorneyClient) ManageAttorneyDecisions(ctx sirius.Context, caseUID string, attorneyDecisions []sirius.AttorneyDecisions) error {
	args := m.Called(ctx, caseUID, attorneyDecisions)
	return args.Error(0)
}

var ActiveOriginalAttorneyUid = "302b05c7-896c-4290-904e-2005e4f1e81e"
var ActiveReplacementAttorneyUid = "987a01b1-456d-4567-813d-2010d3e2d72d"

var removeAnAttorneyCaseSummary = sirius.CaseSummary{
	DigitalLpa: sirius.DigitalLpa{
		UID: "M-1111-2222-3333",
		SiriusData: sirius.SiriusData{
			Subtype: "personal-welfare",
		},
		LpaStoreData: sirius.LpaStoreData{
			HowAttorneysMakeDecisions:            "jointly",
			HowReplacementAttorneysMakeDecisions: "jointly",
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

var activeAttorneys = []sirius.LpaStoreAttorney{
	{
		LpaStorePerson: sirius.LpaStorePerson{
			Uid:        ActiveOriginalAttorneyUid,
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
			Uid:        ActiveReplacementAttorneyUid,
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
}

var inactiveAttorneys = []sirius.LpaStoreAttorney{
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
			Uid:        InactiveReplacementAttorneyUID,
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
}

var decisionAttorneys = []sirius.LpaStoreAttorney{
	{
		LpaStorePerson: sirius.LpaStorePerson{
			Uid:        ActiveReplacementAttorneyUid,
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
			Uid:        InactiveReplacementAttorneyUID,
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
}

var removeAttorneyReasons = []sirius.RefDataItem{
	{
		Handle:        "BANKRUPT",
		Label:         "Bankrupt",
		ValidSubTypes: []string{"property-and-affairs"},
	},
	{
		Handle:        "DECEASED",
		Label:         "Deceased",
		ValidSubTypes: []string{"property-and-affairs", "personal-welfare"},
	},
}

func TestGetRemoveAnAttorney(t *testing.T) {
	client := &mockRemoveAnAttorneyClient{}
	client.
		On("CaseSummary", mock.Anything, "M-1111-2222-3333").
		Return(removeAnAttorneyCaseSummary, nil)

	client.
		On("RefDataByCategory", mock.Anything, sirius.AttorneyRemovedReasonCategory).
		Return(removeAttorneyReasons, nil)

	removeTemplate := &mockTemplate{}
	removeTemplate.
		On("Func", mock.Anything, removeAnAttorneyData{
			FormName:          "remove",
			Decisions:         "jointly",
			CaseSummary:       removeAnAttorneyCaseSummary,
			ActiveAttorneys:   activeAttorneys,
			InactiveAttorneys: inactiveAttorneys,
			RemovedReasons:    removeAttorneyReasons[1:2], // only second reason is valid for "personal-welfare"
			Error:             sirius.ValidationError{Field: sirius.FieldErrors{}},
		}).
		Return(nil)

	confirmTemplate := &mockTemplate{}
	decisionsTemplate := &mockTemplate{}

	server := newMockServer("/lpa/{uid}/remove-an-attorney", RemoveAnAttorney(client, removeTemplate.Func, confirmTemplate.Func, decisionsTemplate.Func))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-1111-2222-3333/remove-an-attorney", nil)
	resp, err := server.serve(req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	mock.AssertExpectationsForObjects(t, client, removeTemplate)
}

func TestGetRemoveAnAttorneyGetCaseSummaryFails(t *testing.T) {
	caseSummary := sirius.CaseSummary{}

	client := &mockRemoveAnAttorneyClient{}
	client.
		On("CaseSummary", mock.Anything, "M-1111-2222-3333").
		Return(caseSummary, errExample)
	client.
		On("RefDataByCategory", mock.Anything, sirius.AttorneyRemovedReasonCategory).
		Return(removeAttorneyReasons, nil)

	removeTemplate := &mockTemplate{}
	removeTemplate.
		On("Func", mock.Anything, removeAnAttorneyData{}).
		Return(nil)

	confirmTemplate := &mockTemplate{}

	decisionsTemplate := &mockTemplate{}

	server := newMockServer("/lpa/{uid}/remove-an-attorney", RemoveAnAttorney(client, removeTemplate.Func, confirmTemplate.Func, decisionsTemplate.Func))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-1111-2222-3333/remove-an-attorney", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
}

func TestGetRemoveAnAttorneyTemplateErrors(t *testing.T) {
	client := &mockRemoveAnAttorneyClient{}
	client.
		On("CaseSummary", mock.Anything, "M-1111-2222-3333").
		Return(removeAnAttorneyCaseSummary, nil)

	client.
		On("RefDataByCategory", mock.Anything, sirius.AttorneyRemovedReasonCategory).
		Return(removeAttorneyReasons, nil)

	removeTemplate := &mockTemplate{}
	removeTemplate.
		On("Func", mock.Anything, removeAnAttorneyData{
			FormName:          "remove",
			Decisions:         "jointly",
			CaseSummary:       removeAnAttorneyCaseSummary,
			ActiveAttorneys:   activeAttorneys,
			InactiveAttorneys: inactiveAttorneys,
			RemovedReasons:    removeAttorneyReasons[1:2], // only second reason is valid for "personal-welfare"
			Error:             sirius.ValidationError{Field: sirius.FieldErrors{}},
		}).
		Return(errExample)

	confirmTemplate := &mockTemplate{}
	decisionsTemplate := &mockTemplate{}

	server := newMockServer("/lpa/{uid}/remove-an-attorney", RemoveAnAttorney(client, removeTemplate.Func, confirmTemplate.Func, decisionsTemplate.Func))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-1111-2222-3333/remove-an-attorney", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
}

func TestPostRemoveAnAttorneyValidationErrorsRemoveTemplate(t *testing.T) {
	tests := []struct {
		name         string
		form         url.Values
		expectedErr  sirius.ValidationError
		expectedForm formRemoveAttorney
	}{
		{
			name: "Validation error when form fields are empty",
			form: url.Values{},
			expectedErr: sirius.ValidationError{Field: sirius.FieldErrors{
				"removeAttorney": {"reason": "Please select an attorney for removal"},
				"enableAttorney": {"reason": "Please select either the attorneys that can be enabled or skip the replacement of the attorneys"},
				"removedReason":  {"reason": "Please select a reason for removal"},
			}},
			expectedForm: formRemoveAttorney{},
		},
		{
			name: "Validation errors when skip and enable attorneys selected",
			form: url.Values{
				"removedAttorney":    {ActiveOriginalAttorneyUid},
				"enabledAttorney":    {InactiveReplacementAttorneyUID},
				"removedReason":      {"DECEASED"},
				"skipEnableAttorney": {"yes"},
			},
			expectedErr: sirius.ValidationError{Field: sirius.FieldErrors{
				"enableAttorney": {"reason": "Please do not select both a replacement attorney and the option to skip"},
			}},
			expectedForm: formRemoveAttorney{
				RemovedAttorneyUid:  ActiveOriginalAttorneyUid,
				EnabledAttorneyUids: []string{InactiveReplacementAttorneyUID},
				SkipEnableAttorney:  "yes",
				RemovedReason:       "DECEASED",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockRemoveAnAttorneyClient{}
			client.On("CaseSummary", mock.Anything, "M-1111-2222-3333").
				Return(removeAnAttorneyCaseSummary, nil)
			client.On("RefDataByCategory", mock.Anything, sirius.AttorneyRemovedReasonCategory).
				Return(removeAttorneyReasons, nil)

			removeTemplate := &mockTemplate{}
			removeTemplate.
				On("Func", mock.Anything, removeAnAttorneyData{
					FormName:          "remove",
					Decisions:         "jointly",
					CaseSummary:       removeAnAttorneyCaseSummary,
					ActiveAttorneys:   activeAttorneys,
					InactiveAttorneys: inactiveAttorneys,
					RemovedReasons:    removeAttorneyReasons[1:2], // only second reason is valid for "personal-welfare"
					Form:              tc.expectedForm,
					Error:             tc.expectedErr,
				}).
				Return(nil)

			confirmTemplate := &mockTemplate{}
			decisionsTemplate := &mockTemplate{}

			server := newMockServer("/lpa/{uid}/remove-an-attorney", RemoveAnAttorney(client, removeTemplate.Func, confirmTemplate.Func, decisionsTemplate.Func))

			req, _ := http.NewRequest(http.MethodPost, "/lpa/M-1111-2222-3333/remove-an-attorney", strings.NewReader(tc.form.Encode()))
			req.Header.Add("Content-Type", formUrlEncoded)

			resp, err := server.serve(req)

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.Code)
			mock.AssertExpectationsForObjects(t, client, removeTemplate)
		})
	}
}

func TestPostDecisionAttorneyValidationError(t *testing.T) {
	client := &mockRemoveAnAttorneyClient{}
	client.
		On("CaseSummary", mock.Anything, "M-6666-6666-6666").
		Return(manageAttorneyDecisionsSummary, nil)

	client.
		On("RefDataByCategory", mock.Anything, sirius.AttorneyRemovedReasonCategory).
		Return(removeAttorneyReasons, nil)

	removeTemplate := &mockTemplate{}
	decisionsTemplate := &mockTemplate{}
	decisionsTemplate.
		On("Func", mock.Anything, removeAnAttorneyData{
			FormName:          "remove",
			Decisions:         "jointly-for-some-severally-for-others",
			CaseSummary:       manageAttorneyDecisionsSummary,
			ActiveAttorneys:   activeAttorneys,
			InactiveAttorneys: inactiveAttorneys,
			DecisionAttorneys: decisionAttorneys,
			Form: formRemoveAttorney{
				RemovedAttorneyUid:    ActiveOriginalAttorneyUid,
				EnabledAttorneyUids:   []string{InactiveReplacementAttorneyUID},
				RemovedReason:         "DECEASED",
				DecisionAttorneysUids: []string{ActiveReplacementAttorneyUid},
				SkipDecisionAttorney:  "yes",
			},
			Error: sirius.ValidationError{Field: sirius.FieldErrors{
				"decisionAttorney": {"reason": "Select who cannot make joint decisions, or select 'Joint decisions can be made by all attorneys'"},
			}},
		}).
		Return(nil)

	confirmTemplate := &mockTemplate{}

	server := newMockServer("/lpa/{uid}/remove-an-attorney", RemoveAnAttorney(client, removeTemplate.Func, confirmTemplate.Func, decisionsTemplate.Func))

	form := url.Values{
		"removedAttorney":      {ActiveOriginalAttorneyUid},
		"enabledAttorney":      {InactiveReplacementAttorneyUID},
		"decisionAttorney":     {ActiveReplacementAttorneyUid},
		"removedReason":        {"DECEASED"},
		"skipDecisionAttorney": {"yes"},
		"step":                 {"decision"},
	}

	req, _ := http.NewRequest(http.MethodPost, "/lpa/M-6666-6666-6666/remove-an-attorney", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", formUrlEncoded)
	resp, err := server.serve(req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	mock.AssertExpectationsForObjects(t, client, decisionsTemplate)
}

func TestPostRemoveAnAttorneyWithoutDecisionsConfirmTemplate(t *testing.T) {
	client := &mockRemoveAnAttorneyClient{}
	client.
		On("CaseSummary", mock.Anything, "M-1111-2222-3333").
		Return(removeAnAttorneyCaseSummary, nil)

	client.
		On("RefDataByCategory", mock.Anything, sirius.AttorneyRemovedReasonCategory).
		Return(removeAttorneyReasons, nil)

	removeTemplate := &mockTemplate{}
	confirmTemplate := &mockTemplate{}
	confirmTemplate.
		On("Func", mock.Anything, removeAnAttorneyData{
			FormName:          "remove",
			Decisions:         "jointly",
			CaseSummary:       removeAnAttorneyCaseSummary,
			ActiveAttorneys:   activeAttorneys,
			InactiveAttorneys: inactiveAttorneys,
			RemovedReasons:    removeAttorneyReasons[1:2], // only second reason is valid for "personal-welfare"
			Form: formRemoveAttorney{
				RemovedAttorneyUid:  ActiveOriginalAttorneyUid,
				EnabledAttorneyUids: []string{InactiveReplacementAttorneyUID},
				SkipEnableAttorney:  "",
				RemovedReason:       "DECEASED",
			},
			RemovedAttorneysDetails: SelectedAttorneyDetails{
				SelectedAttorneyName: "Jack Black",
				SelectedAttorneyDob:  "1990-02-22",
			},
			EnabledAttorneysDetails: []SelectedAttorneyDetails{
				{
					SelectedAttorneyName: "Jack Green",
					SelectedAttorneyDob:  "1990-02-26",
				},
			},
			RemovedReason: removeAttorneyReasons[1],
			Error:         sirius.ValidationError{Field: sirius.FieldErrors{}},
		}).
		Return(nil)
	decisionsTemplate := &mockTemplate{}

	server := newMockServer("/lpa/{uid}/remove-an-attorney", RemoveAnAttorney(client, removeTemplate.Func, confirmTemplate.Func, decisionsTemplate.Func))

	form := url.Values{
		"removedAttorney": {ActiveOriginalAttorneyUid},
		"enabledAttorney": {InactiveReplacementAttorneyUID},
		"removedReason":   {"DECEASED"},
		"step":            {"remove"},
	}

	req, _ := http.NewRequest(http.MethodPost, "/lpa/M-1111-2222-3333/remove-an-attorney", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", formUrlEncoded)
	resp, err := server.serve(req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	mock.AssertExpectationsForObjects(t, client, confirmTemplate)
}

func TestPostConfirmAttorneyRemovalWithoutDecisionsRedirects(t *testing.T) {
	client := &mockRemoveAnAttorneyClient{}
	client.
		On("CaseSummary", mock.Anything, "M-1111-2222-3333").
		Return(removeAnAttorneyCaseSummary, nil)

	client.
		On("RefDataByCategory", mock.Anything, sirius.AttorneyRemovedReasonCategory).
		Return(removeAttorneyReasons, nil)

	client.
		On("ChangeAttorneyStatus", mock.Anything, "M-1111-2222-3333", []sirius.AttorneyUpdatedStatus{
			{UID: ActiveOriginalAttorneyUid, Status: "removed", RemovedReason: "DECEASED"},
			{UID: InactiveReplacementAttorneyUID, Status: "active"},
		}).
		Return(nil)

	removeTemplate := &mockTemplate{}
	confirmTemplate := &mockTemplate{}
	decisionsTemplate := &mockTemplate{}

	server := newMockServer("/lpa/{uid}/remove-an-attorney", RemoveAnAttorney(client, removeTemplate.Func, confirmTemplate.Func, decisionsTemplate.Func))

	form := url.Values{
		"removedAttorney": {ActiveOriginalAttorneyUid},
		"enabledAttorney": {InactiveReplacementAttorneyUID},
		"removedReason":   {"DECEASED"},
		"step":            {"confirm"},
	}

	req, _ := http.NewRequest(http.MethodPost, "/lpa/M-1111-2222-3333/remove-an-attorney", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", formUrlEncoded)
	resp, err := server.serve(req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, RedirectError("/lpa/M-1111-2222-3333"), err)
}

func TestPostRemoveAttorneyWithDecisionsOnDecisionsTemplate(t *testing.T) {
	client := &mockRemoveAnAttorneyClient{}
	client.
		On("CaseSummary", mock.Anything, "M-6666-6666-6666").
		Return(manageAttorneyDecisionsSummary, nil)

	client.
		On("RefDataByCategory", mock.Anything, sirius.AttorneyRemovedReasonCategory).
		Return(removeAttorneyReasons, nil)

	removeTemplate := &mockTemplate{}
	decisionsTemplate := &mockTemplate{}
	decisionsTemplate.
		On("Func", mock.Anything, removeAnAttorneyData{
			FormName:          "remove",
			Decisions:         "jointly-for-some-severally-for-others",
			CaseSummary:       manageAttorneyDecisionsSummary,
			ActiveAttorneys:   activeAttorneys,
			InactiveAttorneys: inactiveAttorneys,
			DecisionAttorneys: decisionAttorneys,
			Form: formRemoveAttorney{
				RemovedAttorneyUid:  ActiveOriginalAttorneyUid,
				EnabledAttorneyUids: []string{InactiveReplacementAttorneyUID},
				SkipEnableAttorney:  "",
				RemovedReason:       "DECEASED",
			},
			RemovedAttorneysDetails: SelectedAttorneyDetails{
				SelectedAttorneyName: "Jack Black",
				SelectedAttorneyDob:  "1990-02-22",
			},
			EnabledAttorneysDetails: []SelectedAttorneyDetails{
				{
					SelectedAttorneyName: "Jack Green",
					SelectedAttorneyDob:  "1990-02-26",
				},
			},
			RemovedReason: removeAttorneyReasons[1],
			Error:         sirius.ValidationError{Field: sirius.FieldErrors{}},
		}).
		Return(nil)
	confirmTemplate := &mockTemplate{}

	server := newMockServer("/lpa/{uid}/remove-an-attorney", RemoveAnAttorney(client, removeTemplate.Func, confirmTemplate.Func, decisionsTemplate.Func))

	form := url.Values{
		"removedAttorney": {ActiveOriginalAttorneyUid},
		"enabledAttorney": {InactiveReplacementAttorneyUID},
		"removedReason":   {"DECEASED"},
		"step":            {"remove"},
	}

	req, _ := http.NewRequest(http.MethodPost, "/lpa/M-6666-6666-6666/remove-an-attorney", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", formUrlEncoded)
	resp, err := server.serve(req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	mock.AssertExpectationsForObjects(t, client, decisionsTemplate)
}

func TestPostRemoveAttorneyWithDecisionsOnConfirmTemplate(t *testing.T) {
	client := &mockRemoveAnAttorneyClient{}
	client.
		On("CaseSummary", mock.Anything, "M-6666-6666-6666").
		Return(manageAttorneyDecisionsSummary, nil)

	client.
		On("RefDataByCategory", mock.Anything, sirius.AttorneyRemovedReasonCategory).
		Return(removeAttorneyReasons, nil)

	removeTemplate := &mockTemplate{}
	decisionsTemplate := &mockTemplate{}
	confirmTemplate := &mockTemplate{}
	confirmTemplate.
		On("Func", mock.Anything, removeAnAttorneyData{
			FormName:          "remove",
			Decisions:         "jointly-for-some-severally-for-others",
			CaseSummary:       manageAttorneyDecisionsSummary,
			ActiveAttorneys:   activeAttorneys,
			InactiveAttorneys: inactiveAttorneys,
			Form: formRemoveAttorney{
				RemovedAttorneyUid:    ActiveOriginalAttorneyUid,
				EnabledAttorneyUids:   []string{InactiveReplacementAttorneyUID},
				RemovedReason:         "DECEASED",
				DecisionAttorneysUids: []string{ActiveReplacementAttorneyUid},
			},
			RemovedAttorneysDetails: SelectedAttorneyDetails{
				SelectedAttorneyName: "Jack Black",
				SelectedAttorneyDob:  "1990-02-22",
			},
			EnabledAttorneysDetails: []SelectedAttorneyDetails{
				{
					SelectedAttorneyName: "Jack Green",
					SelectedAttorneyDob:  "1990-02-26",
				},
			},
			DecisionAttorneysDetails: []AttorneyDetails{
				{
					AttorneyName:    "Shelley Jones",
					AttorneyDob:     "1990-02-27",
					AppointmentType: shared.ReplacementAppointmentType.String(),
				},
			},
			RemovedReason: removeAttorneyReasons[1],
			Error:         sirius.ValidationError{Field: sirius.FieldErrors{}},
		}).
		Return(nil)

	server := newMockServer("/lpa/{uid}/remove-an-attorney", RemoveAnAttorney(client, removeTemplate.Func, confirmTemplate.Func, decisionsTemplate.Func))

	form := url.Values{
		"removedAttorney":  {ActiveOriginalAttorneyUid},
		"enabledAttorney":  {InactiveReplacementAttorneyUID},
		"decisionAttorney": {ActiveReplacementAttorneyUid},
		"removedReason":    {"DECEASED"},
		"step":             {"decision"},
	}

	req, _ := http.NewRequest(http.MethodPost, "/lpa/M-6666-6666-6666/remove-an-attorney", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", formUrlEncoded)
	resp, err := server.serve(req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	mock.AssertExpectationsForObjects(t, client, confirmTemplate)
}

func TestPostConfirmAttorneyRemovalWithDecisionAttorneys(t *testing.T) {
	tests := []struct {
		name              string
		formKey           string
		formValue         string
		expectedDecisions []sirius.AttorneyDecisions
	}{
		{
			name:      "User selects specific attorneys who cannot make Decisions",
			formKey:   "decisionAttorney",
			formValue: ActiveReplacementAttorneyUid,
			expectedDecisions: []sirius.AttorneyDecisions{
				{UID: ActiveReplacementAttorneyUid, CannotMakeJointDecisions: true},
				{UID: InactiveReplacementAttorneyUID, CannotMakeJointDecisions: false},
				{UID: ActiveOriginalAttorneyUid, CannotMakeJointDecisions: false},
			},
		},
		{
			name:      "User allows all attorneys to make decisions",
			formKey:   "skipDecisionAttorney",
			formValue: "yes",
			expectedDecisions: []sirius.AttorneyDecisions{
				{UID: ActiveOriginalAttorneyUid, CannotMakeJointDecisions: false},
				{UID: ActiveReplacementAttorneyUid, CannotMakeJointDecisions: false},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockRemoveAnAttorneyClient{}
			client.
				On("CaseSummary", mock.Anything, "M-6666-6666-6666").
				Return(manageAttorneyDecisionsSummary, nil)

			client.
				On("RefDataByCategory", mock.Anything, sirius.AttorneyRemovedReasonCategory).
				Return(removeAttorneyReasons, nil)

			client.
				On("ChangeAttorneyStatus", mock.Anything, "M-6666-6666-6666", []sirius.AttorneyUpdatedStatus{
					{UID: ActiveOriginalAttorneyUid, Status: "removed", RemovedReason: "DECEASED"},
					{UID: InactiveReplacementAttorneyUID, Status: "active"},
				}).
				Return(nil)

			client.
				On("ManageAttorneyDecisions", mock.Anything, "M-6666-6666-6666", tc.expectedDecisions).
				Return(nil)

			removeTemplate := &mockTemplate{}
			confirmTemplate := &mockTemplate{}
			decisionsTemplate := &mockTemplate{}

			server := newMockServer("/lpa/{uid}/remove-an-attorney", RemoveAnAttorney(client, removeTemplate.Func, confirmTemplate.Func, decisionsTemplate.Func))

			formData := url.Values{
				"removedAttorney": {ActiveOriginalAttorneyUid},
				"enabledAttorney": {InactiveReplacementAttorneyUID},
				"removedReason":   {"DECEASED"},
				"step":            {"confirm"},
				tc.formKey:        {tc.formValue},
			}

			req, _ := http.NewRequest(http.MethodPost, "/lpa/M-6666-6666-6666/remove-an-attorney", strings.NewReader(formData.Encode()))
			req.Header.Add("Content-Type", formUrlEncoded)

			resp, err := server.serve(req)

			assert.Equal(t, http.StatusOK, resp.Code)
			assert.Equal(t, RedirectError("/lpa/M-6666-6666-6666"), err)
		})
	}
}
