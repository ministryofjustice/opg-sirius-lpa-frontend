package server

import (
	"errors"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type mockApplyFeeReductionClient struct {
	mock.Mock
}

func (m *mockApplyFeeReductionClient) ApplyFeeReduction(ctx sirius.Context, caseID int, source string, feeReductionType string, paymentEvidence string, appliedDate sirius.DateString) error {
	return m.Called(ctx, caseID, source, feeReductionType, paymentEvidence, appliedDate).Error(0)
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
	expectedError := errors.New("err")

	client := &mockApplyFeeReductionClient{}
	client.
		On("Case", mock.Anything, 4).
		Return(sirius.Case{}, expectedError)

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

	expectedError := errors.New("err")

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

	expectedError := errors.New("err")

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
	caseitem := sirius.Case{CaseType: "lpa", UID: "700700"}

	feeReductionTypes := []sirius.RefDataItem{
		{
			Handle: "REMISSION",
			Label:  "Remission",
		},
	}

	client := &mockApplyFeeReductionClient{}
	client.
		On("ApplyFeeReduction", mock.Anything, 123, "FEE_REDUCTION", "REMISSION", "Test evidence", sirius.DateString("2022-01-23")).
		Return(nil)
	client.
		On("Case", mock.Anything, 123).
		Return(caseitem, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.FeeReductionTypeCategory).
		Return(feeReductionTypes, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, applyFeeReductionData{
			Success:           true,
			Case:              caseitem,
			FeeReductionType:  "REMISSION",
			Source:            "FEE_REDUCTION",
			AppliedDate:       "2022-01-23",
			PaymentEvidence:   "Test evidence",
			FeeReductionTypes: feeReductionTypes,
		}).
		Return(nil)

	form := url.Values{
		"feeReductionType": {"REMISSION"},
		"source":           {"FEE_REDUCTION"},
		"appliedDate":      {"2022-01-23"},
		"paymentEvidence":  {"Test evidence"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := ApplyFeeReduction(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}
