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

var removeAnAttorneyCaseSummary = sirius.CaseSummary{
	DigitalLpa: sirius.DigitalLpa{
		UID: "M-1111-2222-3333",
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
					Email:           "b@example.com",
					Mobile:          "07122121212",
					SignedAt:        "2024-11-28T19:22:11Z",
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

func TestGetRemoveAnAttorney(t *testing.T) {
	client := &mockRemoveAnAttorneyClient{}
	client.
		On("CaseSummary", mock.Anything, "M-1111-2222-3333").
		Return(removeAnAttorneyCaseSummary, nil)

	removeTemplate := &mockTemplate{}
	removeTemplate.
		On("Func", mock.Anything, removeAnAttorneyData{
			CaseSummary:     removeAnAttorneyCaseSummary,
			ActiveAttorneys: activeAttorneys,
			Error:           sirius.ValidationError{Field: sirius.FieldErrors{}},
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
		Return(caseSummary, expectedError)

	removeTemplate := &mockTemplate{}
	removeTemplate.
		On("Func", mock.Anything, removeAnAttorneyData{}).
		Return(nil)

	confirmTemplate := &mockTemplate{}

	server := newMockServer("/lpa/{uid}/remove-an-attorney", RemoveAnAttorney(client, removeTemplate.Func, confirmTemplate.Func))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-1111-2222-3333/remove-an-attorney", nil)
	_, err := server.serve(req)

	assert.Equal(t, expectedError, err)
}

func TestGetRemoveAnAttorneyTemplateErrors(t *testing.T) {
	client := &mockRemoveAnAttorneyClient{}
	client.
		On("CaseSummary", mock.Anything, "M-1111-2222-3333").
		Return(removeAnAttorneyCaseSummary, nil)

	removeTemplate := &mockTemplate{}
	removeTemplate.
		On("Func", mock.Anything, removeAnAttorneyData{
			CaseSummary:     removeAnAttorneyCaseSummary,
			ActiveAttorneys: activeAttorneys,
			Error:           sirius.ValidationError{Field: sirius.FieldErrors{}},
		}).
		Return(expectedError)

	confirmTemplate := &mockTemplate{}

	server := newMockServer("/lpa/{uid}/remove-an-attorney", RemoveAnAttorney(client, removeTemplate.Func, confirmTemplate.Func))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-1111-2222-3333/remove-an-attorney", nil)
	_, err := server.serve(req)

	assert.Equal(t, expectedError, err)
}

func TestPostRemoveAnAttorneyInvalidData(t *testing.T) {
	client := &mockRemoveAnAttorneyClient{}
	client.
		On("CaseSummary", mock.Anything, "M-1111-2222-3333").
		Return(removeAnAttorneyCaseSummary, nil)

	removeTemplate := &mockTemplate{}
	removeTemplate.
		On("Func", mock.Anything, removeAnAttorneyData{
			SelectedAttorneyUid: "",
			CaseSummary:         removeAnAttorneyCaseSummary,
			ActiveAttorneys:     activeAttorneys,
			Error: sirius.ValidationError{Field: sirius.FieldErrors{
				"selectAttorney": {"reason": "Please select an attorney for removal."},
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
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	mock.AssertExpectationsForObjects(t, client, removeTemplate)
}

func TestPostRemoveAnAttorneyValidData(t *testing.T) {
	client := &mockRemoveAnAttorneyClient{}
	client.
		On("CaseSummary", mock.Anything, "M-1111-2222-3333").
		Return(removeAnAttorneyCaseSummary, nil)

	removeTemplate := &mockTemplate{}
	confirmTemplate := &mockTemplate{}
	confirmTemplate.
		On("Func", mock.Anything, removeAnAttorneyData{
			CaseSummary:          removeAnAttorneyCaseSummary,
			ActiveAttorneys:      activeAttorneys,
			SelectedAttorneyUid:  "302b05c7-896c-4290-904e-2005e4f1e81e",
			SelectedAttorneyName: "Jack Black",
			SelectedAttorneyDob:  "1990-02-22",
			Error:                sirius.ValidationError{Field: sirius.FieldErrors{}},
		}).
		Return(nil)

	server := newMockServer("/lpa/{uid}/remove-an-attorney", RemoveAnAttorney(client, removeTemplate.Func, confirmTemplate.Func))

	form := url.Values{
		"selectedAttorney": {"302b05c7-896c-4290-904e-2005e4f1e81e"},
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

	removeTemplate := &mockTemplate{}
	confirmTemplate := &mockTemplate{}

	server := newMockServer("/lpa/{uid}/remove-an-attorney", RemoveAnAttorney(client, removeTemplate.Func, confirmTemplate.Func))

	form := url.Values{
		"selectedAttorney": {"302b05c7-896c-4290-904e-2005e4f1e81e"},
		"confirmRemoval":   {""},
	}

	req, _ := http.NewRequest(http.MethodPost, "/lpa/M-1111-2222-3333/remove-an-attorney", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", formUrlEncoded)
	_, err := server.serve(req)

	assert.Equal(t, RedirectError("/lpa/M-1111-2222-3333"), err)
	mock.AssertExpectationsForObjects(t, client, removeTemplate)
}
