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

func (m *mockGetLpaDetailsClient) CaseSummaryWithImages(ctx sirius.Context, uid string) (sirius.CaseSummary, error) {
	args := m.Called(ctx, uid)
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
		On("CaseSummaryWithImages", mock.Anything, "M-EEEE-9876-9876").
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
	tests := []struct {
		name               string
		taskList           []sirius.Task
		reviewRestrictions bool
	}{

		{
			name:               "Empty Task List",
			taskList:           []sirius.Task{},
			reviewRestrictions: false,
		},
		{
			name: "A closed review restrictions and conditions task exists",
			taskList: []sirius.Task{
				{
					ID:     1,
					Name:   "Review restrictions and conditions",
					Status: "Completed",
				},
			},
			reviewRestrictions: false,
		},
		{
			name: "An Open Review Restrictions and Conditions Task Exists",
			taskList: []sirius.Task{
				{
					ID:     1,
					Name:   "Review restrictions and conditions",
					Status: "Not started",
				},
			},
			reviewRestrictions: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
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
								Decisions:       false,
								Status:          shared.InactiveAttorneyStatus.String(),
								AppointmentType: shared.ReplacementAppointmentType.String(),
								LpaStorePerson: sirius.LpaStorePerson{
									Uid:   "1",
									Email: "first@does.not.exist",
								},
							},
							{
								Decisions:       false,
								Status:          shared.InactiveAttorneyStatus.String(),
								AppointmentType: shared.ReplacementAppointmentType.String(),
								LpaStorePerson: sirius.LpaStorePerson{
									Email: "second@does.not.exist",
								},
							},
							{
								Decisions:       true,
								Status:          shared.ActiveAttorneyStatus.String(),
								AppointmentType: shared.OriginalAppointmentType.String(),
								LpaStorePerson: sirius.LpaStorePerson{
									Email: "third@does.not.exist",
								},
							},
							{
								Decisions:       false,
								Status:          shared.ActiveAttorneyStatus.String(),
								AppointmentType: shared.OriginalAppointmentType.String(),
								LpaStorePerson: sirius.LpaStorePerson{
									Email: "fourth@does.not.exist",
								},
							},
							{
								Decisions:       false,
								Status:          shared.RemovedAttorneyStatus.String(),
								AppointmentType: shared.OriginalAppointmentType.String(),
								LpaStorePerson: sirius.LpaStorePerson{
									Email: "fifth@does.not.exist",
								},
							},
							{
								Decisions:       true,
								Status:          shared.ActiveAttorneyStatus.String(),
								AppointmentType: shared.ReplacementAppointmentType.String(),
								LpaStorePerson: sirius.LpaStorePerson{
									Email: "sixth@does.not.exist",
								},
							},
						},
						TrustCorporations: []sirius.LpaStoreTrustCorporation{
							{
								Name:          "Trust Me Once",
								CompanyNumber: "123456789",
								Signatories: []sirius.Signatory{
									{
										FirstNames: "First",
									},
								},
								LpaStoreAttorney: sirius.LpaStoreAttorney{
									Decisions:       true,
									Status:          shared.ActiveAttorneyStatus.String(),
									AppointmentType: shared.OriginalAppointmentType.String(),
									LpaStorePerson: sirius.LpaStorePerson{
										Email: "trust.me.once@does.not.exist",
									},
								},
							},
							{
								Name:          "Trust Me Twice",
								CompanyNumber: "987654321",
								Signatories: []sirius.Signatory{
									{
										FirstNames: "Second",
									},
								},
								LpaStoreAttorney: sirius.LpaStoreAttorney{
									Decisions:       false,
									Status:          shared.InactiveAttorneyStatus.String(),
									AppointmentType: shared.ReplacementAppointmentType.String(),
									LpaStorePerson: sirius.LpaStorePerson{
										Email: "trust.me.twice@does.not.exist",
									},
								},
							},
						},
					},
				},
				TaskList: tc.taskList,
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
				On("CaseSummaryWithImages", mock.Anything, "M-9876-9876-9876").
				Return(caseSummary, nil)
			client.
				On("AnomaliesForDigitalLpa", mock.Anything, "M-9876-9876-9876").
				Return(anomalies, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, getLpaDetails{
					CaseSummary:        caseSummary,
					DigitalLpa:         caseSummary.DigitalLpa,
					AnomalyDisplay:     &expectedAnomalyDisplay,
					ReviewRestrictions: tc.reviewRestrictions,
					ReplacementAttorneys: []sirius.LpaStoreAttorney{
						{
							Decisions:       false,
							Status:          shared.InactiveAttorneyStatus.String(),
							AppointmentType: shared.ReplacementAppointmentType.String(),
							LpaStorePerson: sirius.LpaStorePerson{
								Email: "first@does.not.exist",
								Uid:   "1",
							},
						},
						{
							Decisions:       false,
							Status:          shared.InactiveAttorneyStatus.String(),
							AppointmentType: shared.ReplacementAppointmentType.String(),
							LpaStorePerson: sirius.LpaStorePerson{
								Email: "second@does.not.exist",
							},
						},
					},
					NonReplacementAttorneys: []sirius.LpaStoreAttorney{
						{
							Decisions:       true,
							Status:          shared.ActiveAttorneyStatus.String(),
							AppointmentType: shared.OriginalAppointmentType.String(),
							LpaStorePerson: sirius.LpaStorePerson{
								Email: "third@does.not.exist",
							},
						},
						{
							Decisions:       false,
							Status:          shared.ActiveAttorneyStatus.String(),
							AppointmentType: shared.OriginalAppointmentType.String(),
							LpaStorePerson: sirius.LpaStorePerson{
								Email: "fourth@does.not.exist",
							},
						},
						{
							Decisions:       true,
							Status:          shared.ActiveAttorneyStatus.String(),
							AppointmentType: shared.ReplacementAppointmentType.String(),
							LpaStorePerson: sirius.LpaStorePerson{
								Email: "sixth@does.not.exist",
							},
						},
					},
					RemovedAttorneys: []sirius.LpaStoreAttorney{
						{
							Decisions:       false,
							Status:          shared.RemovedAttorneyStatus.String(),
							AppointmentType: shared.OriginalAppointmentType.String(),
							LpaStorePerson: sirius.LpaStorePerson{
								Email: "fifth@does.not.exist",
							},
						},
					},
					DecisionAttorneys: []sirius.LpaStoreAttorney{
						{
							Decisions:       true,
							Status:          shared.ActiveAttorneyStatus.String(),
							AppointmentType: shared.OriginalAppointmentType.String(),
							LpaStorePerson: sirius.LpaStorePerson{
								Email: "third@does.not.exist",
							},
						},
						{
							Decisions:       true,
							Status:          shared.ActiveAttorneyStatus.String(),
							AppointmentType: shared.ReplacementAppointmentType.String(),
							LpaStorePerson: sirius.LpaStorePerson{
								Email: "sixth@does.not.exist",
							},
						},
					},
					ReplacementTrustCorporations: []sirius.LpaStoreTrustCorporation{
						{
							Name:          "Trust Me Twice",
							CompanyNumber: "987654321",
							Signatories: []sirius.Signatory{
								{
									FirstNames: "Second",
								},
							},
							LpaStoreAttorney: sirius.LpaStoreAttorney{
								Decisions:       false,
								Status:          shared.InactiveAttorneyStatus.String(),
								AppointmentType: shared.ReplacementAppointmentType.String(),
								LpaStorePerson: sirius.LpaStorePerson{
									Email: "trust.me.twice@does.not.exist",
								},
							},
						},
					},
					NonReplacementTrustCorporations: []sirius.LpaStoreTrustCorporation{
						{
							Name:          "Trust Me Once",
							CompanyNumber: "123456789",
							Signatories: []sirius.Signatory{
								{
									FirstNames: "First",
								},
							},
							LpaStoreAttorney: sirius.LpaStoreAttorney{
								Decisions:       true,
								Status:          shared.ActiveAttorneyStatus.String(),
								AppointmentType: shared.OriginalAppointmentType.String(),
								LpaStorePerson: sirius.LpaStorePerson{
									Email: "trust.me.once@does.not.exist",
								},
							},
						},
					},
					DecisionTrustCorporations: []sirius.LpaStoreTrustCorporation{
						{
							Name:          "Trust Me Once",
							CompanyNumber: "123456789",
							Signatories: []sirius.Signatory{
								{
									FirstNames: "First",
								},
							},
							LpaStoreAttorney: sirius.LpaStoreAttorney{
								Decisions:       true,
								Status:          shared.ActiveAttorneyStatus.String(),
								AppointmentType: shared.OriginalAppointmentType.String(),
								LpaStorePerson: sirius.LpaStorePerson{
									Email: "trust.me.once@does.not.exist",
								},
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
		})
	}
}
