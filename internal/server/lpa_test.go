package server

import (
	"net/http"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockLpaClient struct {
	mock.Mock
}

func (m *mockLpaClient) CaseSummary(ctx sirius.Context, uid string) (sirius.CaseSummary, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(sirius.CaseSummary), args.Error(1)
}

func TestGetLpaErrorRetrievingCaseSummary(t *testing.T) {
	client := &mockLpaClient{}
	client.
		On("CaseSummary", mock.Anything, "M-AAAA-9876-9876").
		Return(sirius.CaseSummary{}, expectedError)

	template := &mockTemplate{}

	server := newMockServer("/lpa/{uid}", Lpa(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-AAAA-9876-9876", nil)
	_, err := server.serve(req)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetLpa(t *testing.T) {
	caseSummary := sirius.CaseSummary{
		DigitalLpa: sirius.DigitalLpa{
			UID: "M-9876-9876-9876",
			SiriusData: sirius.SiriusData{
				ID:      22,
				Subtype: "hw",
			},
		},
		TaskList: []sirius.Task{},
	}

	client := &mockLpaClient{}
	client.
		On("CaseSummary", mock.Anything, "M-9876-9876-9876").
		Return(caseSummary, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, lpaData{
			CaseSummary: caseSummary,
			DigitalLpa:  caseSummary.DigitalLpa,
		}).
		Return(nil)

	server := newMockServer("/lpa/{uid}", Lpa(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-9876-9876-9876", nil)
	_, err := server.serve(req)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}
