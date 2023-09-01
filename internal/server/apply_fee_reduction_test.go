package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockApplyFeeReductionClient struct {
	mock.Mock
}

func (m *mockApplyFeeReductionClient) ApplyFeeReduction(ctx sirius.Context, caseID int, feeReductionType string, paymentEvidence string, paymentDate sirius.DateString) error {
	return m.Called(ctx, caseID, feeReductionType, paymentEvidence, paymentDate).Error(0)
}

func (m *mockApplyFeeReductionClient) Case(ctx sirius.Context, id int) (sirius.Case, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Case), args.Error(1)
}

func (m *mockApplyFeeReductionClient) RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error) {
	args := m.Called(ctx, category)
	if args.Get(0) != nil {
		return args.Get(0).([]sirius.RefDataItem), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestGetApplyFeeReduction(t *testing.T) {
	caseItem := sirius.Case{
		UID:     "7000-0000-0021",
		SubType: "pfa",
	}

	feeReductionTypes := []sirius.RefDataItem{
		{
			Handle: "REMISSION",
			Label:  "Remission",
		},
	}

	client := &mockApplyFeeReductionClient{}
	client.
		On("Case", mock.Anything, 4).
		Return(caseItem, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.FeeReductionTypeCategory).
		Return(feeReductionTypes, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, applyFeeReductionData{
			Case:              caseItem,
			FeeReductionTypes: feeReductionTypes,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=4", nil)
	w := httptest.NewRecorder()

	err := ApplyFeeReduction(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestApplyFeeReductionNoID(t *testing.T) {
	testCases := map[string]string{
		"no-id":  "/",
		"bad-id": "/?id=test",
	}

	for name, testUrl := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, testUrl, nil)
			w := httptest.NewRecorder()

			err := ApplyFeeReduction(nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestApplyFeeReductionWhenFailureOnGetCase(t *testing.T) {
	feeReductionTypes := []sirius.RefDataItem{
		{
			Handle: "REMISSION",
			Label:  "Remission",
		},
	}
	client := &mockApplyFeeReductionClient{}
	client.
		On("Case", mock.Anything, 4).
		Return(sirius.Case{}, expectedError)
	client.
		On("RefDataByCategory", mock.Anything, sirius.FeeReductionTypeCategory).
		Return(feeReductionTypes, nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=4", nil)
	w := httptest.NewRecorder()

	err := ApplyFeeReduction(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestApplyFeeReductionWhenTemplateErrors(t *testing.T) {
	caseItem := sirius.Case{
		UID:     "7000-0000-0021",
		SubType: "pfa",
	}

	feeReductionTypes := []sirius.RefDataItem{
		{
			Handle: "REMISSION",
			Label:  "Remission",
		},
	}

	client := &mockApplyFeeReductionClient{}
	client.
		On("Case", mock.Anything, 123).
		Return(caseItem, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.FeeReductionTypeCategory).
		Return(feeReductionTypes, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, applyFeeReductionData{
			Case:              caseItem,
			FeeReductionTypes: feeReductionTypes,
		}).
		Return(expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := ApplyFeeReduction(client, template.Func)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestApplyFeeReductionWhenFailureOnGetFeeReductionTypesRefData(t *testing.T) {
	caseItem := sirius.Case{
		UID:     "7000-0000-0021",
		SubType: "pfa",
	}

	client := &mockApplyFeeReductionClient{}
	client.
		On("Case", mock.Anything, 123).
		Return(caseItem, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.FeeReductionTypeCategory).
		Return([]sirius.RefDataItem{}, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := ApplyFeeReduction(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestPostFeeReduction(t *testing.T) {
	testCases := []struct {
		name string
		id int
		caseItem sirius.Case
		paymentDate string
		expectedRedirect string
	}{
		{
			name: "Non-digital LPA",
			id: 123,
			caseItem: sirius.Case{CaseType: "LPA", UID: "700700"},
			paymentDate: "2022-01-23",
			expectedRedirect: "/payments/123",
		},
		{
			name: "Digital LPA",
			id: 456,
			caseItem: sirius.Case{CaseType: "DIGITAL_LPA", UID: "M-AAA-BBB-CCC"},
			paymentDate: "2023-09-01",
			expectedRedirect: "/lpa/M-AAA-BBB-CCC/payments",
		},
	}

	feeReductionTypes := []sirius.RefDataItem{
		{
			Handle: "REMISSION",
			Label:  "Remission",
		},
	}

	for _, tc := range(testCases) {
		client := &mockApplyFeeReductionClient{}
		client.
			On("ApplyFeeReduction", mock.Anything, tc.id, "REMISSION", "Test evidence", sirius.DateString(tc.paymentDate)).
			Return(nil)
		client.
			On("Case", mock.Anything, tc.id).
			Return(tc.caseItem, nil)
		client.
			On("RefDataByCategory", mock.Anything, sirius.FeeReductionTypeCategory).
			Return(feeReductionTypes, nil)

		template := &mockTemplate{}

		form := url.Values{
			"feeReductionType": {"REMISSION"},
			"paymentDate":      {tc.paymentDate},
			"paymentEvidence":  {"Test evidence"},
		}

		r, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/?id=%d", tc.id), strings.NewReader(form.Encode()))
		r.Header.Add("Content-Type", formUrlEncoded)
		w := httptest.NewRecorder()

		err := ApplyFeeReduction(client, template.Func)(w, r)
		resp := w.Result()

		assert.Equal(t, RedirectError(tc.expectedRedirect), err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mock.AssertExpectationsForObjects(t, client, template)
	}
}
