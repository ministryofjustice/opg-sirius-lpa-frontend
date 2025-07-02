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

var removeAnAttorneyCaseSummary = sirius.CaseSummary{
	DigitalLpa: sirius.DigitalLpa{
		UID: "M-1111-2222-3333",
		SiriusData: sirius.SiriusData{
			Subtype: "personal-welfare",
		},
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

var activeAttorneys = []sirius.LpaStoreAttorney{
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
			CaseSummary:       removeAnAttorneyCaseSummary,
			ActiveAttorneys:   activeAttorneys,
			InactiveAttorneys: inactiveAttorneys,
			RemovedReasons:    removeAttorneyReasons[1:2], // only second reason is valid for "personal-welfare"
			Error:             sirius.ValidationError{Field: sirius.FieldErrors{}},
		}).
		Return(nil)

	confirmTemplate := &mockTemplate{}

	server := newMockServer("/lpa/{uid}/remove-an-attorney", RemoveAnAttorney(client, removeTemplate.Func, confirmTemplate.Func))

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

	removeTemplate := &mockTemplate{}
	removeTemplate.
		On("Func", mock.Anything, removeAnAttorneyData{}).
		Return(nil)

	confirmTemplate := &mockTemplate{}

	server := newMockServer("/lpa/{uid}/remove-an-attorney", RemoveAnAttorney(client, removeTemplate.Func, confirmTemplate.Func))

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
			CaseSummary:       removeAnAttorneyCaseSummary,
			ActiveAttorneys:   activeAttorneys,
			InactiveAttorneys: inactiveAttorneys,
			RemovedReasons:    removeAttorneyReasons[1:2], // only second reason is valid for "personal-welfare"
			Error:             sirius.ValidationError{Field: sirius.FieldErrors{}},
		}).
		Return(errExample)

	confirmTemplate := &mockTemplate{}

	server := newMockServer("/lpa/{uid}/remove-an-attorney", RemoveAnAttorney(client, removeTemplate.Func, confirmTemplate.Func))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-1111-2222-3333/remove-an-attorney", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
}

func TestPostRemoveAnAttorneyInvalidData(t *testing.T) {
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
			CaseSummary:       removeAnAttorneyCaseSummary,
			ActiveAttorneys:   activeAttorneys,
			InactiveAttorneys: inactiveAttorneys,
			RemovedReasons:    removeAttorneyReasons[1:2], // only second reason is valid for "personal-welfare"
			Form: formRemoveAttorney{
				RemovedAttorneyUid:  "",
				EnabledAttorneyUids: nil,
				SkipEnableAttorney:  "",
				RemovedReason:       "",
			},
			Error: sirius.ValidationError{Field: sirius.FieldErrors{
				"removeAttorney": {"reason": "Please select an attorney for removal"},
				"enableAttorney": {"reason": "Please select either the attorneys that can be enabled or skip the replacement of the attorneys"},
				"removedReason":  {"reason": "Please select a reason for removal"},
			}},
		}).
		Return(nil)

	confirmTemplate := &mockTemplate{}

	server := newMockServer("/lpa/{uid}/remove-an-attorney", RemoveAnAttorney(client, removeTemplate.Func, confirmTemplate.Func))

	form := url.Values{}

	req, _ := http.NewRequest(http.MethodPost, "/lpa/M-1111-2222-3333/remove-an-attorney", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", formUrlEncoded)
	resp, err := server.serve(req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	mock.AssertExpectationsForObjects(t, client, removeTemplate)
}

func TestPostRemoveAnAttorneyValidData(t *testing.T) {
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
			CaseSummary:       removeAnAttorneyCaseSummary,
			ActiveAttorneys:   activeAttorneys,
			InactiveAttorneys: inactiveAttorneys,
			RemovedReasons:    removeAttorneyReasons[1:2], // only second reason is valid for "personal-welfare"
			Form: formRemoveAttorney{
				RemovedAttorneyUid:  "302b05c7-896c-4290-904e-2005e4f1e81e",
				EnabledAttorneyUids: []string{"123a01b1-456d-5391-813d-2010d3e256f"},
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

	server := newMockServer("/lpa/{uid}/remove-an-attorney", RemoveAnAttorney(client, removeTemplate.Func, confirmTemplate.Func))

	form := url.Values{
		"removedAttorney": {"302b05c7-896c-4290-904e-2005e4f1e81e"},
		"enabledAttorney": {"123a01b1-456d-5391-813d-2010d3e256f"},
		"removedReason":   {"DECEASED"},
	}

	req, _ := http.NewRequest(http.MethodPost, "/lpa/M-1111-2222-3333/remove-an-attorney", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", formUrlEncoded)
	resp, err := server.serve(req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	mock.AssertExpectationsForObjects(t, client, confirmTemplate)
}

func TestPostConfirmAttorneyRemovalValidData(t *testing.T) {
	client := &mockRemoveAnAttorneyClient{}
	client.
		On("CaseSummary", mock.Anything, "M-1111-2222-3333").
		Return(removeAnAttorneyCaseSummary, nil)

	client.
		On("RefDataByCategory", mock.Anything, sirius.AttorneyRemovedReasonCategory).
		Return(removeAttorneyReasons, nil)

	client.
		On("ChangeAttorneyStatus", mock.Anything, "M-1111-2222-3333", []sirius.AttorneyUpdatedStatus{
			{UID: "302b05c7-896c-4290-904e-2005e4f1e81e", Status: "removed", RemovedReason: "DECEASED"},
			{UID: "123a01b1-456d-5391-813d-2010d3e256f", Status: "active"},
		}).
		Return(nil)

	removeTemplate := &mockTemplate{}
	confirmTemplate := &mockTemplate{}

	server := newMockServer("/lpa/{uid}/remove-an-attorney", RemoveAnAttorney(client, removeTemplate.Func, confirmTemplate.Func))

	form := url.Values{
		"removedAttorney": {"302b05c7-896c-4290-904e-2005e4f1e81e"},
		"enabledAttorney": {"123a01b1-456d-5391-813d-2010d3e256f"},
		"removedReason":   {"DECEASED"},
		"confirmRemoval":  {""},
	}

	req, _ := http.NewRequest(http.MethodPost, "/lpa/M-1111-2222-3333/remove-an-attorney", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", formUrlEncoded)
	resp, err := server.serve(req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, RedirectError("/lpa/M-1111-2222-3333"), err)
}
