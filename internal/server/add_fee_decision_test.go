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

type mockAddFeeDecisionClient struct {
	mock.Mock
}

func (m *mockAddFeeDecisionClient) AddFeeDecision(ctx sirius.Context, caseID int, decisionType string, decisionReason string, decisionDate sirius.DateString) error {
	return m.Called(ctx, caseID, decisionType, decisionReason, decisionDate).Error(0)
}

func (m *mockAddFeeDecisionClient) Case(ctx sirius.Context, id int) (sirius.Case, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Case), args.Error(1)
}

func (m *mockAddFeeDecisionClient) RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error) {
	args := m.Called(ctx, category)
	if args.Get(0) != nil {
		return args.Get(0).([]sirius.RefDataItem), args.Error(1)
	}
	return nil, args.Error(1)
}

var feeDecisionTypes = []sirius.RefDataItem{
	{
		Handle: "DECLINED_EXEMPTION",
		Label:  "Declined exemption",
	},
}

func TestGetAddFeeDecision(t *testing.T) {
	caseItem := sirius.Case{
		UID: "7000-0000-0021",
	}

	client := &mockAddFeeDecisionClient{}
	client.
		On("Case", mock.Anything, 4).
		Return(caseItem, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.FeeDecisionTypeCategory).
		Return(feeDecisionTypes, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, addFeeDecisionData{
			Case: caseItem,
			DecisionTypes: feeDecisionTypes,
			ReturnUrl: "/payments/4",
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=4", nil)
	w := httptest.NewRecorder()

	err := AddFeeDecision(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestAddFeeDecisionNoID(t *testing.T) {
	testCases := map[string]string{
		"no-id":  "/",
		"bad-id": "/?id=test",
	}

	for name, testUrl := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, testUrl, nil)
			w := httptest.NewRecorder()

			err := AddFeeDecision(nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestAddFeeDecisionWhenFailureOnGetCase(t *testing.T) {
	client := &mockAddFeeDecisionClient{}
	client.
		On("Case", mock.Anything, 75757).
		Return(sirius.Case{}, expectedError)
	client.
		On("RefDataByCategory", mock.Anything, sirius.FeeDecisionTypeCategory).
		Return(feeDecisionTypes, nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=75757", nil)
	w := httptest.NewRecorder()

	err := AddFeeDecision(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestAddFeeDecisionWhenTemplateErrors(t *testing.T) {
	caseItem := sirius.Case{
		UID:     "7000-0000-0021",
		SubType: "pfa",
	}

	client := &mockAddFeeDecisionClient{}
	client.
		On("Case", mock.Anything, 111).
		Return(caseItem, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.FeeDecisionTypeCategory).
		Return(feeDecisionTypes, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, addFeeDecisionData{
			Case: caseItem,
			DecisionTypes: feeDecisionTypes,
			ReturnUrl: "/payments/111",
		}).
		Return(expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=111", nil)
	w := httptest.NewRecorder()

	err := AddFeeDecision(client, template.Func)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestAddFeeDecisionWhenFailureOnGetRefData(t *testing.T) {
	caseItem := sirius.Case{
		UID:     "7000-0000-0021",
		SubType: "pfa",
	}

	client := &mockAddFeeDecisionClient{}
	client.
		On("Case", mock.Anything, 232).
		Return(caseItem, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.FeeDecisionTypeCategory).
		Return([]sirius.RefDataItem{}, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=232", nil)
	w := httptest.NewRecorder()

	err := AddFeeDecision(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestPostAddFeeDecisionToDigitalLpa(t *testing.T) {
	caseItem := sirius.Case{CaseType: "DIGITAL_LPA", UID: "M-AAAA-BBBB-DDDD"}

	client := &mockAddFeeDecisionClient{}
	client.
		On("AddFeeDecision", mock.Anything, 454, "DECLINED_REMISSION", "Invalid evidence", sirius.DateString("2023-10-10")).
		Return(nil)
	client.
		On("Case", mock.Anything, 454).
		Return(caseItem, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.FeeDecisionTypeCategory).
		Return(feeDecisionTypes, nil)

	template := &mockTemplate{}

	form := url.Values{
		"decisionReason": {"Invalid evidence"},
		"decisionType":   {"DECLINED_REMISSION"},
		"decisionDate":   {"2023-10-10"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=454", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := AddFeeDecision(client, template.Func)(w, r)
	resp := w.Result()

	assert.Equal(t, RedirectError("/lpa/M-AAAA-BBBB-DDDD/payments"), err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}
