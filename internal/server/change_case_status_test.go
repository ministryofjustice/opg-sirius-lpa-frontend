package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockChangeCaseStatusClient struct {
	mock.Mock
}

func (m *mockChangeCaseStatusClient) CaseSummary(ctx sirius.Context, uid string) (sirius.CaseSummary, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(sirius.CaseSummary), args.Error(1)
}

func (m *mockChangeCaseStatusClient) EditDigitalLPAStatus(ctx sirius.Context, caseUID string, caseStatusData sirius.CaseStatusData) error {
	args := m.Called(ctx, caseUID, caseStatusData)
	return args.Error(0)
}

func TestGetChangeCaseStatus(t *testing.T) {
	caseSummary := sirius.CaseSummary{
		DigitalLpa: sirius.DigitalLpa{
			UID: "M-9876-9876-9876",
			SiriusData: sirius.SiriusData{
				ID:      676,
				Subtype: "personal-welfare",
				Status:  "Draft",
			},
		},
	}

	client := &mockChangeCaseStatusClient{}
	client.
		On("CaseSummary", mock.Anything, "M-9876-9876-9876").
		Return(caseSummary, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, changeCaseStatusData{
			Entity:    "personal-welfare M-9876-9876-9876",
			OldStatus: "Draft",
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/change-case-status?uid=M-9876-9876-9876", nil)
	w := httptest.NewRecorder()

	err := ChangeCaseStatus(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}
