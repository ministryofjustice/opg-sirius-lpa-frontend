package server

import (
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type mockConfirmAttorneyRemovalClient struct {
	mock.Mock
}

func (m *mockConfirmAttorneyRemovalClient) CaseSummary(ctx sirius.Context, uid string) (sirius.CaseSummary, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(sirius.CaseSummary), args.Error(1)
}

var confirmAttorneyRemovalCaseSummary = sirius.CaseSummary{
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

var attorneyRemoved = sirius.LpaStoreAttorney{
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
}

func TestGetConfirmAttorneyRemoval(t *testing.T) {
	client := &mockConfirmAttorneyRemovalClient{}
	client.
		On("CaseSummary", mock.Anything, "M-1111-2222-3333").
		Return(confirmAttorneyRemovalCaseSummary, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, confirmAttorneyRemovalData{
			CaseSummary: confirmAttorneyRemovalCaseSummary,
			Attorney:    attorneyRemoved,
			Error:       sirius.ValidationError{Field: sirius.FieldErrors{}},
		}).
		Return(nil)

	server := newMockServer("/lpa/{uid}/confirm-attorney-removal/{attorneyUID}", ConfirmAttorneyRemoval(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-1111-2222-3333/confirm-attorney-removal/302b05c7-896c-4290-904e-2005e4f1e81e", nil)
	resp, err := server.serve(req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetConfirmAttorneyRemovalGetCaseSummaryFails(t *testing.T) {
	caseSummary := sirius.CaseSummary{}

	client := &mockConfirmAttorneyRemovalClient{}
	client.
		On("CaseSummary", mock.Anything, "M-1111-2222-3333").
		Return(caseSummary, expectedError)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, confirmAttorneyRemovalData{}).
		Return(nil)

	server := newMockServer("/lpa/{uid}/confirm-attorney-removal/{attorneyUID}", ConfirmAttorneyRemoval(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-1111-2222-3333/confirm-attorney-removal/302b05c7-896c-4290-904e-2005e4f1e81e", nil)
	_, err := server.serve(req)

	assert.Equal(t, expectedError, err)
}

func TestGetConfirmAttorneyRemovalTemplateErrors(t *testing.T) {
	client := &mockConfirmAttorneyRemovalClient{}
	client.
		On("CaseSummary", mock.Anything, "M-1111-2222-3333").
		Return(confirmAttorneyRemovalCaseSummary, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, confirmAttorneyRemovalData{
			CaseSummary: confirmAttorneyRemovalCaseSummary,
			Attorney:    attorneyRemoved,
			Error:       sirius.ValidationError{Field: sirius.FieldErrors{}},
		}).
		Return(expectedError)

	server := newMockServer("/lpa/{uid}/confirm-attorney-removal/{attorneyUID}", ConfirmAttorneyRemoval(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-1111-2222-3333/confirm-attorney-removal/302b05c7-896c-4290-904e-2005e4f1e81e", nil)
	_, err := server.serve(req)

	assert.Equal(t, expectedError, err)
}

func TestPostConfirmAttorneyRemovalInvalidData(t *testing.T) {
	client := &mockConfirmAttorneyRemovalClient{}
	client.
		On("CaseSummary", mock.Anything, "M-1111-2222-3333").
		Return(confirmAttorneyRemovalCaseSummary, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, confirmAttorneyRemovalData{
			SelectedAttorney: "",
			CaseSummary:      confirmAttorneyRemovalCaseSummary,
			Attorney:         attorneyRemoved,
			Error: sirius.ValidationError{Field: sirius.FieldErrors{
				"selectAttorney": {"reason": "Please select an attorney for removal."},
			}},
		}).
		Return(nil)

	server := newMockServer("/lpa/{uid}/confirm-attorney-removal/{attorneyUID}", ConfirmAttorneyRemoval(client, template.Func))

	form := url.Values{}

	req, _ := http.NewRequest(http.MethodPost, "/lpa/M-1111-2222-3333/confirm-attorney-removal/302b05c7-896c-4290-904e-2005e4f1e81e", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", formUrlEncoded)
	resp, err := server.serve(req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostConfirmAttorneyRemovalValidData(t *testing.T) {
	caseSummary := sirius.CaseSummary{
		DigitalLpa: sirius.DigitalLpa{
			UID: "M-1111-2222-3333",
		},
	}

	client := &mockConfirmAttorneyRemovalClient{}
	client.
		On("CaseSummary", mock.Anything, "M-1111-2222-3333").
		Return(caseSummary, nil)

	template := &mockTemplate{}

	server := newMockServer("/lpa/{uid}/confirm-attorney-removal/{attorneyUID}", ConfirmAttorneyRemoval(client, template.Func))

	form := url.Values{
		"selectedAttorney": {"302b05c7-896c-4290-904e-2005e4f1e81e"},
	}

	req, _ := http.NewRequest(http.MethodPost, "/lpa/M-1111-2222-3333/confirm-attorney-removal/302b05c7-896c-4290-904e-2005e4f1e81e", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", formUrlEncoded)
	_, err := server.serve(req)

	assert.Equal(t, RedirectError("/lpa/M-1111-2222-3333/confirm-attorney-removal/302b05c7-896c-4290-904e-2005e4f1e81e"), err)
	mock.AssertExpectationsForObjects(t, client, template)
}
