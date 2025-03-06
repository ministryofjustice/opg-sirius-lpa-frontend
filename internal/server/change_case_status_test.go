package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockChangeCaseStatusClient struct {
	mock.Mock
}

func (m *mockChangeCaseStatusClient) CaseSummary(ctx sirius.Context, uid string, presignImages bool) (sirius.CaseSummary, error) {
	args := m.Called(ctx, uid, presignImages)
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
			LpaStoreData: sirius.LpaStoreData{
				Status: "draft",
			},
		},
	}

	client := &mockChangeCaseStatusClient{}
	client.
		On("CaseSummary", mock.Anything, "M-9876-9876-9876", false).
		Return(caseSummary, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, changeCaseStatusData{
			Entity:    "personal-welfare M-9876-9876-9876",
			CaseUID:   "M-9876-9876-9876",
			OldStatus: "draft",
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

func TestPostChangeCaseStatus(t *testing.T) {
	caseSummary := sirius.CaseSummary{
		DigitalLpa: sirius.DigitalLpa{
			UID: "M-9876-9876-9876",
			SiriusData: sirius.SiriusData{
				ID:      676,
				Subtype: "personal-welfare",
				Status:  "Draft",
			},
			LpaStoreData: sirius.LpaStoreData{
				Status: "draft",
			},
		},
	}

	client := &mockChangeCaseStatusClient{}
	client.
		On("CaseSummary", mock.Anything, "M-9876-9876-9876", false).
		Return(caseSummary, nil)

	client.
		On("EditDigitalLPAStatus", mock.Anything, "M-9876-9876-9876", sirius.CaseStatusData{
			Status: "cannot-register",
		}).
		Return(nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, changeCaseStatusData{
			Success:   true,
			Entity:    "personal-welfare M-9876-9876-9876",
			CaseUID:   "M-9876-9876-9876",
			OldStatus: "cannot-register",
			NewStatus: "cannot-register",
		}).
		Return(nil)

	form := url.Values{
		"status": {"cannot-register"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/change-case-status?uid=M-9876-9876-9876", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := ChangeCaseStatus(client, template.Func)(w, r)
	resp := w.Result()

	assert.Equal(t, RedirectError("/lpa/M-9876-9876-9876"), err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
