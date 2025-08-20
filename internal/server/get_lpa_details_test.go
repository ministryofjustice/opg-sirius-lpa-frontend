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

func TestLpaDetailsSignedOnBehalf(t *testing.T) {
	caseSummary := sirius.CaseSummary{
		DigitalLpa: sirius.DigitalLpa{
			UID: "M-1234-5678-ABCD",
			SiriusData: sirius.SiriusData{
				ID:      123,
				Subtype: "hw",
			},
			LpaStoreData: sirius.LpaStoreData{
				AuthorisedSignatory: &sirius.LpaStoreAuthorisedSignatory{
					FirstNames: "John",
					LastName:   "Smith",
					SignedAt:   "2024-01-15T10:30:00Z",
				},
				WitnessedByCertificateProviderAt: "2024-01-15T10:31:00Z",
				WitnessedByIndependentWitnessAt:  "2024-01-15T10:32:00Z",
				IndependentWitness: &sirius.LpaStoreIndependentWitness{
					FirstNames: "Jane",
					LastName:   "Doe",
					Address: sirius.LpaStoreAddress{
						Line1:    "123 Test Street",
						Town:     "Test Town",
						Postcode: "T3ST 1NG",
						Country:  "GB",
					},
				},
				Attorneys: []sirius.LpaStoreAttorney{},
			},
		},
		TaskList: []sirius.Task{},
	}

	var anomalies []sirius.Anomaly

	client := &mockGetLpaDetailsClient{}
	client.
		On("CaseSummaryWithImages", mock.Anything, "M-1234-5678-ABCD").
		Return(caseSummary, nil)
	client.
		On("AnomaliesForDigitalLpa", mock.Anything, "M-1234-5678-ABCD").
		Return(anomalies, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.MatchedBy(func(data getLpaDetails) bool {
			lpa := data.CaseSummary.DigitalLpa

			signatoryCheck := lpa.WasSignedOnBehalfOfDonor() == true
			signatoryNameCheck := lpa.GetAuthorisedSignatoryFullName() == "John Smith"
			witness1Check := lpa.WasWitnessedByCertificateProvider() == true
			witness2Check := lpa.WasWitnessedByIndependentWitness() == true
			witness2NameCheck := lpa.GetIndependentWitnessFullName() == "Jane Doe"

			address := lpa.GetIndependentWitnessAddress()
			addressCheck := address.Line1 == "123 Test Street" && address.Town == "Test Town"

			dataStructureCheck := data.CaseSummary.DigitalLpa.UID == "M-1234-5678-ABCD" &&
				data.AnomalyDisplay != nil &&
				data.ReviewRestrictions == false

			return signatoryCheck && signatoryNameCheck && witness1Check &&
				witness2Check && witness2NameCheck && addressCheck && dataStructureCheck
		})).
		Return(nil)

	server := newMockServer("/lpa/{uid}/lpa-details", GetLpaDetails(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-1234-5678-ABCD/lpa-details", nil)
	_, err := server.serve(req)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestLpaDetailsSignedDirectly(t *testing.T) {
	caseSummary := sirius.CaseSummary{
		DigitalLpa: sirius.DigitalLpa{
			UID: "M-5678-1234-DCBA",
			SiriusData: sirius.SiriusData{
				ID:      456,
				Subtype: "hw",
			},
			LpaStoreData: sirius.LpaStoreData{
				// No AuthorisedSignatory - donor signed directly
				SignedAt:                         "2024-01-15T10:30:00Z",
				WitnessedByCertificateProviderAt: "2024-01-15T10:31:00Z",
				// No IndependentWitness or WitnessedByIndependentWitnessAt
				Attorneys: []sirius.LpaStoreAttorney{},
			},
		},
		TaskList: []sirius.Task{},
	}

	var anomalies []sirius.Anomaly

	client := &mockGetLpaDetailsClient{}
	client.
		On("CaseSummaryWithImages", mock.Anything, "M-5678-1234-DCBA").
		Return(caseSummary, nil)
	client.
		On("AnomaliesForDigitalLpa", mock.Anything, "M-5678-1234-DCBA").
		Return(anomalies, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.MatchedBy(func(data getLpaDetails) bool {
			lpa := data.CaseSummary.DigitalLpa

			signatoryCheck := lpa.WasSignedOnBehalfOfDonor() == false
			signatoryNameCheck := lpa.GetAuthorisedSignatoryFullName() == ""
			witness1Check := lpa.WasWitnessedByCertificateProvider() == true
			witness2Check := lpa.WasWitnessedByIndependentWitness() == false
			witness2NameCheck := lpa.GetIndependentWitnessFullName() == ""

			address := lpa.GetIndependentWitnessAddress()
			addressCheck := address.Line1 == "" && address.Town == ""

			return signatoryCheck && signatoryNameCheck && witness1Check &&
				witness2Check && witness2NameCheck && addressCheck
		})).
		Return(nil)

	server := newMockServer("/lpa/{uid}/lpa-details", GetLpaDetails(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-5678-1234-DCBA/lpa-details", nil)
	_, err := server.serve(req)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}
