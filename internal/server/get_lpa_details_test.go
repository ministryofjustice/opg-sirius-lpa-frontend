package server

import (
	"errors"
	"net/http"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockGetLpaDetailsClient struct {
	mock.Mock
}

func (m *mockGetLpaDetailsClient) CaseSummary(ctx sirius.Context, uid string, presignImages bool) (sirius.CaseSummary, error) {
	args := m.Called(ctx, uid, presignImages)
	return args.Get(0).(sirius.CaseSummary), args.Error(1)
}

func (m *mockGetLpaDetailsClient) AnomaliesForDigitalLpa(ctx sirius.Context, uid string) ([]sirius.Anomaly, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).([]sirius.Anomaly), args.Error(1)
}

func TestGetLpaDetailsCaseSummaryFail(t *testing.T) {
	expectedError := errors.New("network error")

	client := &mockGetLpaDetailsClient{}
	client.
		On("CaseSummary", mock.Anything, "M-EEEE-9876-9876", true).
		Return(sirius.CaseSummary{}, expectedError)
	client.
		On("AnomaliesForDigitalLpa", mock.Anything, "M-EEEE-9876-9876").
		Return([]sirius.Anomaly{}, nil)

	template := &mockTemplate{}

	server := newMockServer("/lpa/{uid}/lpa-details", GetLpaDetails(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-EEEE-9876-9876/lpa-details", nil)
	_, err := server.serve(req)

	assert.Error(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetLpaDetailsSuccess(t *testing.T) {
	caseSummary := sirius.CaseSummary{
		DigitalLpa: sirius.DigitalLpa{
			UID: "M-9876-9876-9876",
			SiriusData: sirius.SiriusData{
				ID:      22,
				Subtype: "hw",
			},
			LpaStoreData: sirius.LpaStoreData{
				Attorneys: []sirius.LpaStoreAttorney{
					{
						Status:          shared.InactiveAttorneyStatus.String(),
						AppointmentType: shared.ReplacementAppointmentType.String(),
						LpaStorePerson: sirius.LpaStorePerson{
							Uid:   "1",
							Email: "first@does.not.exist",
						},
					},
					{
						Status:          shared.InactiveAttorneyStatus.String(),
						AppointmentType: shared.ReplacementAppointmentType.String(),
						LpaStorePerson: sirius.LpaStorePerson{
							Email: "second@does.not.exist",
						},
					},
					{
						Status:          shared.ActiveAttorneyStatus.String(),
						AppointmentType: shared.OriginalAppointmentType.String(),
						LpaStorePerson: sirius.LpaStorePerson{
							Email: "third@does.not.exist",
						},
					},
					{
						Status:          shared.ActiveAttorneyStatus.String(),
						AppointmentType: shared.OriginalAppointmentType.String(),
						LpaStorePerson: sirius.LpaStorePerson{
							Email: "fourth@does.not.exist",
						},
					},
					{
						Status:          shared.RemovedAttorneyStatus.String(),
						AppointmentType: shared.OriginalAppointmentType.String(),
						LpaStorePerson: sirius.LpaStorePerson{
							Email: "fifth@does.not.exist",
						},
					},
					{
						Status:          shared.ActiveAttorneyStatus.String(),
						AppointmentType: shared.ReplacementAppointmentType.String(),
						LpaStorePerson: sirius.LpaStorePerson{
							Email: "sixth@does.not.exist",
						},
					},
				},
			},
		},
		TaskList: []sirius.Task{},
	}

	anomalies := []sirius.Anomaly{
		{
			Id:            999,
			Status:        sirius.AnomalyFatal,
			FieldName:     sirius.ObjectFieldName("lastName"),
			FieldOwnerUid: sirius.ObjectUid("1"),
		},
	}

	expectedAnomalyDisplay := sirius.AnomalyDisplay{
		AnomaliesBySection: map[sirius.AnomalyDisplaySection]sirius.AnomaliesForSection{
			sirius.ReplacementAttorneysSection: {
				Section: sirius.ReplacementAttorneysSection,
				Objects: map[sirius.ObjectUid]sirius.AnomaliesForObject{
					sirius.ObjectUid("1"): {
						Uid: sirius.ObjectUid("1"),
						Anomalies: map[sirius.ObjectFieldName][]sirius.Anomaly{
							sirius.ObjectFieldName("lastName"): {
								{
									Id:            999,
									Status:        sirius.AnomalyFatal,
									FieldName:     sirius.ObjectFieldName("lastName"),
									FieldOwnerUid: sirius.ObjectUid("1"),
								},
							},
						},
					},
				},
			},
		},
	}

	client := &mockGetLpaDetailsClient{}
	client.
		On("CaseSummary", mock.Anything, "M-9876-9876-9876", true).
		Return(caseSummary, nil)
	client.
		On("AnomaliesForDigitalLpa", mock.Anything, "M-9876-9876-9876").
		Return(anomalies, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, getLpaDetails{
			CaseSummary:    caseSummary,
			DigitalLpa:     caseSummary.DigitalLpa,
			AnomalyDisplay: &expectedAnomalyDisplay,
			ReplacementAttorneys: []sirius.LpaStoreAttorney{
				{
					Status:          shared.InactiveAttorneyStatus.String(),
					AppointmentType: shared.ReplacementAppointmentType.String(),
					LpaStorePerson: sirius.LpaStorePerson{
						Email: "first@does.not.exist",
						Uid:   "1",
					},
				},
				{
					Status:          shared.InactiveAttorneyStatus.String(),
					AppointmentType: shared.ReplacementAppointmentType.String(),
					LpaStorePerson: sirius.LpaStorePerson{
						Email: "second@does.not.exist",
					},
				},
			},
			NonReplacementAttorneys: []sirius.LpaStoreAttorney{
				{
					Status:          shared.ActiveAttorneyStatus.String(),
					AppointmentType: shared.OriginalAppointmentType.String(),
					LpaStorePerson: sirius.LpaStorePerson{
						Email: "third@does.not.exist",
					},
				},
				{
					Status:          shared.ActiveAttorneyStatus.String(),
					AppointmentType: shared.OriginalAppointmentType.String(),
					LpaStorePerson: sirius.LpaStorePerson{
						Email: "fourth@does.not.exist",
					},
				},
				{
					Status:          shared.ActiveAttorneyStatus.String(),
					AppointmentType: shared.ReplacementAppointmentType.String(),
					LpaStorePerson: sirius.LpaStorePerson{
						Email: "sixth@does.not.exist",
					},
				},
			},
		}).
		Return(nil)

	server := newMockServer("/lpa/{uid}/lpa-details", GetLpaDetails(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-9876-9876-9876/lpa-details", nil)
	_, err := server.serve(req)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}
