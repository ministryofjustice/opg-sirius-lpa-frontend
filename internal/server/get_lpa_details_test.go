package server

import (
	"encoding/json"
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

func TestGetLpaDetailsSuccess(t *testing.T) {
	caseSummary := sirius.CaseSummary{
		DigitalLpa: sirius.DigitalLpa{
			UID: "M-9876-9876-9876",
			SiriusData: sirius.SiriusData{
				ID:      22,
				Subtype: "hw",
			},
			LpaStoreData: json.RawMessage([]byte(`{"some": "data"}`)),
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
			LpaStoreData: map[string]interface{}{
				"some": "data",
			},
		}).
		Return(nil)

	server := newMockServer("/lpa/{uid}/lpa-details", GetLpaDetails(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-9876-9876-9876/lpa-details", nil)
	_, err := server.serve(req)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}
