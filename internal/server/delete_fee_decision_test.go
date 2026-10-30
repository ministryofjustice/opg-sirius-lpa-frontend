package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockDeleteFeeDecisionClient struct {
	mock.Mock
}

func (m *mockDeleteFeeDecisionClient) DeleteFeeDecision(ctx sirius.Context, feeDecisionID int) error {
	return m.Called(ctx, feeDecisionID).Error(0)
}

func (m *mockDeleteFeeDecisionClient) GetFeeDecision(ctx sirius.Context, feeDecisionID int) (sirius.FeeDecision, error) {
	args := m.Called(ctx, feeDecisionID)
	return args.Get(0).(sirius.FeeDecision), args.Error(1)
}

func TestGetDeleteFeeDecision(t *testing.T) {
	feeDecision := sirius.FeeDecision{
		ID:             8,
		DecisionType:   shared.DecisionTypeExemption,
		DecisionReason: "test",
		DecisionDate:   sirius.DateString("2026-01-01"),
	}
	client := &mockDeleteFeeDecisionClient{}
	client.
		On("GetFeeDecision", mock.Anything, 8).
		Return(feeDecision, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, deleteFeeDecisionData{
			CaseID:          3,
			FeeDecision:     feeDecision,
			ReturnUrl:       "/payments/3",
			FeeDecisionType: "declined exemption",
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/delete-decline-fee-decision?caseId=3&feeDecisionId=8", nil)
	w := httptest.NewRecorder()

	err := DeleteFeeDecision(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestDeleteFeeDecisionBadID(t *testing.T) {
	testCases := map[string]string{
		"no-id":       "/",
		"bad-case-id": "/?caseId=test",
		"bad-fee-id":  "/?caseId=1&feeDecisionId=test",
	}

	for name, testUrl := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, testUrl, nil)
			w := httptest.NewRecorder()

			err := DeleteFeeDecision(nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestGetDeleteFeeDecisionWhenGetFeeDecisionErrors(t *testing.T) {
	client := &mockDeleteFeeDecisionClient{}
	client.
		On("GetFeeDecision", mock.Anything, 8).
		Return(sirius.FeeDecision{}, errExample)

	r, _ := http.NewRequest(http.MethodGet, "/delete-decline-fee-decision?caseId=3&feeDecisionId=8", nil)
	w := httptest.NewRecorder()

	err := DeleteFeeDecision(client, nil)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestPostDeleteFeeDecisionHtmx(t *testing.T) {
	feeDecision := sirius.FeeDecision{
		ID:             8,
		DecisionType:   shared.DecisionTypeRemission,
		DecisionReason: "test",
		DecisionDate:   sirius.DateString("2026-01-01"),
	}
	client := &mockDeleteFeeDecisionClient{}
	client.
		On("GetFeeDecision", mock.Anything, 8).
		Return(feeDecision, nil)
	client.
		On("DeleteFeeDecision", mock.Anything, 8).
		Return(nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, deleteFeeDecisionData{
			CaseID:          3,
			FeeDecision:     feeDecision,
			FeeDecisionType: "declined remission",
			ReturnUrl:       "/payments/3",
			HtmxRedirect:    "/payments/3",
			IsPartial:       true,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodPost, "/delete-decline-fee-decision?caseId=3&feeDecisionId=8", nil)
	r.Header.Add("HX-Request", "true")

	w := httptest.NewRecorder()

	err := DeleteFeeDecision(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostDeleteFeeDecision(t *testing.T) {
	cases := []struct {
		name                 string
		deleteFeeDecisionErr error
		expectedErr          error
	}{
		{
			name:                 "Success - OK",
			deleteFeeDecisionErr: nil,
			expectedErr:          RedirectError("/payments/3"),
		},
		{
			name:                 "API error",
			deleteFeeDecisionErr: errExample,
			expectedErr:          errExample,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			feeDecision := sirius.FeeDecision{
				ID:             8,
				DecisionType:   shared.DecisionTypeNotRecognised,
				DecisionReason: "test",
				DecisionDate:   sirius.DateString("2026-01-01"),
			}
			client := &mockDeleteFeeDecisionClient{}
			client.
				On("GetFeeDecision", mock.Anything, 8).
				Return(feeDecision, nil)
			client.
				On("DeleteFeeDecision", mock.Anything, 8).
				Return(tc.deleteFeeDecisionErr)

			template := &mockTemplate{}

			r, _ := http.NewRequest(http.MethodPost, "/delete-decline-fee-decision?caseId=3&feeDecisionId=8", nil)
			w := httptest.NewRecorder()

			err := DeleteFeeDecision(client, template.Func)(w, r)
			resp := w.Result()

			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}
