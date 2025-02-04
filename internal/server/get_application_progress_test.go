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

type mockApplicationProgressClient struct {
	mock.Mock
}

func (m *mockApplicationProgressClient) CaseSummary(ctx sirius.Context, uid string) (sirius.CaseSummary, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(sirius.CaseSummary), args.Error(1)
}

func (m *mockApplicationProgressClient) ProgressIndicatorsForDigitalLpa(ctx sirius.Context, uid string) ([]sirius.ProgressIndicator, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).([]sirius.ProgressIndicator), args.Error(1)
}

func TestGetApplicationProgressSuccess(t *testing.T) {
	caseSummary := sirius.CaseSummary{
		DigitalLpa: sirius.DigitalLpa{
			UID: "M-9876-9876-9876",
			SiriusData: sirius.SiriusData{
				ID:      22,
				Subtype: "hw",
				Application: sirius.Draft{
					Source: "APPLICANT",
				},
			},
			LpaStoreData: sirius.LpaStoreData{
				Channel: "Digital",
				Attorneys: []sirius.LpaStoreAttorney{
					{
						Status:          shared.InactiveAttorneyStatus.String(),
						AppointmentType: shared.ReplacementAppointmentType.String(),
						LpaStorePerson: sirius.LpaStorePerson{
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
				},
				CertificateProvider: sirius.LpaStoreCertificateProvider{
					LpaStorePerson: sirius.LpaStorePerson{
						FirstNames: "Fake",
						LastName:   "Provider",
					},
					Channel: "Paper",
				},
			},
		},
		TaskList: []sirius.Task{},
	}

	progressIndicators := []sirius.ProgressIndicator{
		sirius.ProgressIndicator{
			Status:    "COMPLETE",
			Indicator: "FEES",
		},
	}

	indicatorView := []IndicatorView{
		{
			UID:                        "M-9876-9876-9876",
			ProgressIndicator:          progressIndicators[0],
			CertificateProviderChannel: "Paper",
			CertificateProviderName:    "Fake Provider",
			ApplicationSource:          "APPLICANT",
		},
	}

	client := &mockApplicationProgressClient{}
	client.
		On("CaseSummary", mock.Anything, "M-9876-9876-9876").
		Return(caseSummary, nil)
	client.
		On("ProgressIndicatorsForDigitalLpa", mock.Anything, "M-9876-9876-9876").
		Return(progressIndicators, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, getApplicationProgressDetails{
			CaseSummary:        caseSummary,
			ProgressIndicators: indicatorView,
		}).
		Return(nil)

	server := newMockServer("/lpa/{uid}", GetApplicationProgressDetails(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-9876-9876-9876", nil)
	_, err := server.serve(req)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetApplicationProgressCaseSummaryFail(t *testing.T) {
	var cs sirius.CaseSummary
	expectedError := errors.New("Case could not be retrieved")

	client := &mockApplicationProgressClient{}
	client.
		On("CaseSummary", mock.Anything, "M-9876-9876-9876").
		Return(cs, expectedError)
	client.
		On("ProgressIndicatorsForDigitalLpa", mock.Anything, "M-9876-9876-9876").
		Return([]sirius.ProgressIndicator{}, nil)

	template := &mockTemplate{}
	template.On("Func", mock.Anything, mock.Anything).Return(nil)

	server := newMockServer("/lpa/{uid}", GetApplicationProgressDetails(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-9876-9876-9876", nil)
	_, err := server.serve(req)

	assert.Equal(t, expectedError, err)
}

func TestGetApplicationProgressProgressIndicatorsFail(t *testing.T) {
	var cs sirius.CaseSummary
	var inds []sirius.ProgressIndicator

	expectedError := errors.New("Progress indicators could not be retrieved")

	client := &mockApplicationProgressClient{}
	client.
		On("CaseSummary", mock.Anything, "M-9876-9876-9876").
		Return(cs, nil)
	client.
		On("ProgressIndicatorsForDigitalLpa", mock.Anything, "M-9876-9876-9876").
		Return(inds, expectedError)

	template := &mockTemplate{}
	template.On("Func", mock.Anything, mock.Anything).Return(nil)

	server := newMockServer("/lpa/{uid}", GetApplicationProgressDetails(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-9876-9876-9876", nil)
	_, err := server.serve(req)

	assert.Equal(t, expectedError, err)
}
