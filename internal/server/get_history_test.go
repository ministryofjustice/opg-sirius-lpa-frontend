package server

import (
	"encoding/json"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
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
	return args.Get(0).(any), args.Error(1)
}

func TestGetHistorySuccess(t *testing.T) {
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
			LpaStoreData: json.RawMessage(`{"some": "data"}`),
		},
		TaskList: []sirius.Task{},
	}

	client := &mockGetHistoryClient{}
	client.
		On("CaseSummary", mock.Anything, "M-9876-9876-9999").
		Return(caseSummary, nil)
	client.
		On("GetEvents", mock.Anything, 8, 12).
		Return("blah", nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, getHistory{
			CaseSummary: caseSummary,
			EventData:   "blah",
		}).
		Return(nil)

	server := newMockServer("/lpa/{uid}/history", GetHistory(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-9876-9876-9999/history", nil)
	_, err := server.serve(req)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}
