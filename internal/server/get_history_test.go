package server

import (
	"net/http"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockGetHistoryClient struct {
	mock.Mock
}

func (m *mockGetHistoryClient) CaseSummary(ctx sirius.Context, uid string) (sirius.CaseSummary, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(sirius.CaseSummary), args.Error(1)
}

func (m *mockGetHistoryClient) GetEvents(ctx sirius.Context, donorId int, caseId int) (any, error) {
	args := m.Called(ctx, donorId, caseId)
	return args.Get(0), args.Error(1)
}

func (m *mockGetHistoryClient) GetCombinedEvents(ctx sirius.Context, uid string, donorId int, caseId int) (any, error) {
	args := m.Called(ctx, uid, donorId, caseId)
	return args.Get(0), args.Error(1)
}

func TestGetHistorySuccessForDigitalLpa(t *testing.T) {
	caseSummary := sirius.CaseSummary{
		DigitalLpa: sirius.DigitalLpa{
			UID: "M-9876-9876-9999",
			SiriusData: sirius.SiriusData{
				ID:      12,
				Subtype: "hw",
				Donor: sirius.Donor{
					ID: 8,
				},
			},
			LpaStoreData: sirius.LpaStoreData{
				Status: "processing", // Non-empty status indicates digital LPA
			},
		},
	}

	client := &mockGetHistoryClient{}
	client.
		On("CaseSummary", mock.Anything, "M-9876-9876-9999").
		Return(caseSummary, nil)
	client.
		On("GetCombinedEvents", mock.Anything, "M-9876-9876-9999", 8, 12).
		Return(map[string]string{"event": "combined event details"}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, getHistory{
			CaseSummary: caseSummary,
			EventData:   map[string]string{"event": "combined event details"},
		}).
		Return(nil)

	server := newMockServer("/lpa/{uid}/history", GetHistory(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-9876-9876-9999/history", nil)
	_, err := server.serve(req)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetHistorySuccessForTraditionalLpa(t *testing.T) {
	caseSummary := sirius.CaseSummary{
		DigitalLpa: sirius.DigitalLpa{
			UID: "M-9876-9876-9999",
			SiriusData: sirius.SiriusData{
				ID:      12,
				Subtype: "hw",
				Donor: sirius.Donor{
					ID: 8,
				},
			},
			LpaStoreData: sirius.LpaStoreData{
				Status: "", // Empty status indicates traditional LPA
			},
		},
	}

	client := &mockGetHistoryClient{}
	client.
		On("CaseSummary", mock.Anything, "M-9876-9876-9999").
		Return(caseSummary, nil)
	client.
		On("GetEvents", mock.Anything, 8, 12).
		Return(map[string]string{"event": "sirius only event details"}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, getHistory{
			CaseSummary: caseSummary,
			EventData:   map[string]string{"event": "sirius only event details"},
		}).
		Return(nil)

	server := newMockServer("/lpa/{uid}/history", GetHistory(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-9876-9876-9999/history", nil)
	_, err := server.serve(req)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetHistoryWhenFailureOnGetCombinedEventsForDigitalLpa(t *testing.T) {
	caseSummary := sirius.CaseSummary{
		DigitalLpa: sirius.DigitalLpa{
			UID: "M-9876-9876-9999",
			SiriusData: sirius.SiriusData{
				ID:      12,
				Subtype: "hw",
				Donor: sirius.Donor{
					ID: 8,
				},
			},
			LpaStoreData: sirius.LpaStoreData{
				Status: "processing",
			},
		},
	}

	client := &mockGetHistoryClient{}
	client.
		On("CaseSummary", mock.Anything, "M-9876-9876-9999").
		Return(caseSummary, nil)
	client.
		On("GetCombinedEvents", mock.Anything, "M-9876-9876-9999", 8, 12).
		Return(nil, errExample)

	server := newMockServer("/lpa/{uid}/history", GetHistory(client, nil))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-9876-9876-9999/history", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetHistoryWhenFailureOnGetEventsForTraditionalLpa(t *testing.T) {
	caseSummary := sirius.CaseSummary{
		DigitalLpa: sirius.DigitalLpa{
			UID: "M-9876-9876-9999",
			SiriusData: sirius.SiriusData{
				ID:      12,
				Subtype: "hw",
				Donor: sirius.Donor{
					ID: 8,
				},
			},
			LpaStoreData: sirius.LpaStoreData{
				Status: "",
			},
		},
	}

	client := &mockGetHistoryClient{}
	client.
		On("CaseSummary", mock.Anything, "M-9876-9876-9999").
		Return(caseSummary, nil)
	client.
		On("GetEvents", mock.Anything, 8, 12).
		Return(nil, errExample)

	server := newMockServer("/lpa/{uid}/history", GetHistory(client, nil))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-9876-9876-9999/history", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetHistoryWhenFailureOnGetCaseSummary(t *testing.T) {
	client := &mockGetHistoryClient{}
	client.
		On("CaseSummary", mock.Anything, "M-9876-9876-9999").
		Return(sirius.CaseSummary{}, errExample)

	template := &mockTemplate{}

	server := newMockServer("/lpa/{uid}/history", GetHistory(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-9876-9876-9999/history", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}
