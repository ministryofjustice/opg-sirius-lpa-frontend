package server

import (
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

type mockGetLpaDetailsClient struct {
	mock.Mock
}

func (m *mockGetLpaDetailsClient) CaseSummary(ctx sirius.Context, uid string) (sirius.CaseSummary, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(sirius.CaseSummary), args.Error(1)
}

func (m *mockGetLpaDetailsClient) AnomaliesForDigitalLpa(ctx sirius.Context, uid string) ([]sirius.Anomaly, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).([]sirius.Anomaly), args.Error(1)
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
						Status: "replacement",
						LpaStorePerson: sirius.LpaStorePerson{
							Email: "first@does.not.exist",
						},
					},
					{
						Status: "replacement",
						LpaStorePerson: sirius.LpaStorePerson{
							Email: "second@does.not.exist",
						},
					},
					{
						Status: "active",
						LpaStorePerson: sirius.LpaStorePerson{
							Email: "third@does.not.exist",
						},
					},
					{
						Status: "active",
						LpaStorePerson: sirius.LpaStorePerson{
							Email: "fourth@does.not.exist",
						},
					},
					{
						Status: "removed",
						LpaStorePerson: sirius.LpaStorePerson{
							Email: "fifth@does.not.exist",
						},
					},
				},
			},
		},
		TaskList: []sirius.Task{},
	}

	client := &mockGetLpaDetailsClient{}
	client.
		On("CaseSummary", mock.Anything, "M-9876-9876-9876").
		Return(caseSummary, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, getLpaDetails{
			CaseSummary: caseSummary,
			DigitalLpa:  caseSummary.DigitalLpa,
			ReplacementAttorneys: []sirius.LpaStoreAttorney{
				{
					Status: "replacement",
					LpaStorePerson: sirius.LpaStorePerson{
						Email: "first@does.not.exist",
					},
				},
				{
					Status: "replacement",
					LpaStorePerson: sirius.LpaStorePerson{
						Email: "second@does.not.exist",
					},
				},
			},
			NonReplacementAttorneys: []sirius.LpaStoreAttorney{
				{
					Status: "active",
					LpaStorePerson: sirius.LpaStorePerson{
						Email: "third@does.not.exist",
					},
				},
				{
					Status: "active",
					LpaStorePerson: sirius.LpaStorePerson{
						Email: "fourth@does.not.exist",
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
