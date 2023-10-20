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

func TestAddFeeDecisionInvalidPostData(t *testing.T) {
	testCases := []struct{
		name string
		postData url.Values
		expectedTemplateData addFeeDecisionData
		expectedValidationError sirius.ValidationError
	}{
		{
			name: "Invalid decision type",
			postData: url.Values{
				"decisionReason": {"Invalid evidence"},
				"decisionDate":   {"2023-10-10"},
			},
			expectedTemplateData: addFeeDecisionData{
				DecisionReason: "Invalid evidence",
				DecisionType: "",
				DecisionDate: "2023-10-10",
			},
			expectedValidationError: sirius.ValidationError{
				Field: sirius.FieldErrors{
					"decisionType": {"reason": "Value is required and can't be empty"},
				},
			},
		},
		{
			name: "Invalid decision date",
			postData: url.Values{
				"decisionReason": {"Invalid evidence"},
				"decisionType":   {"DECLINED_REMISSION"},
			},
			expectedTemplateData: addFeeDecisionData{
				DecisionReason: "Invalid evidence",
				DecisionType: "DECLINED_REMISSION",
				DecisionDate: "",
			},
			expectedValidationError: sirius.ValidationError{
				Field: sirius.FieldErrors{
					"decisionDate": {"reason": "Value is required and can't be empty"},
				},
			},
		},
	}

	caseItem := sirius.Case{CaseType: "DIGITAL_LPA", UID: "M-AAAA-BBBB-DDDD"}
	client := &mockAddFeeDecisionClient{}
	template := &mockTemplate{}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			client.
				On("Case", mock.Anything, 22222).
				Return(caseItem, nil)
			client.
				On("RefDataByCategory", mock.Anything, sirius.FeeDecisionTypeCategory).
				Return(feeDecisionTypes, nil)
			client.
				On("AddFeeDecision",
					mock.Anything,
					22222,
					test.expectedTemplateData.DecisionType,
					test.expectedTemplateData.DecisionReason,
					sirius.DateString(test.expectedTemplateData.DecisionDate)).
				Return(test.expectedValidationError)

			test.expectedTemplateData.Case = caseItem
			test.expectedTemplateData.DecisionTypes = feeDecisionTypes
			test.expectedTemplateData.ReturnUrl = "/lpa/M-AAAA-BBBB-DDDD/payments"
			test.expectedTemplateData.Error = test.expectedValidationError

			template.
				On("Func", mock.Anything, test.expectedTemplateData).
				Return(nil)

			r, _ := http.NewRequest(http.MethodPost, "/?id=22222", strings.NewReader(test.postData.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := AddFeeDecision(client, template.Func)(w, r)

			resp := w.Result()
			mock.AssertExpectationsForObjects(t, client, template)
			assert.Nil(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
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

func TestAddFeeDecisionNonValidationClientError(t *testing.T) {
	caseItem := sirius.Case{CaseType: "DIGITAL_LPA", UID: "M-AAAA-BBBB-ZZEE"}

	client := &mockAddFeeDecisionClient{}
	client.
		On("Case", mock.Anything, 765).
		Return(caseItem, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.FeeDecisionTypeCategory).
		Return(feeDecisionTypes, nil)

	// client gets a 500 error back when adding the fee decision via the API
	serverError := sirius.StatusError{
		Code: http.StatusInternalServerError,
		URL: "https://not.real/",
		Method: http.MethodPost,
		CorrelationId: "Uncorrelated",
	}
	client.
		On("AddFeeDecision", mock.Anything, 765, "DECLINED_REMISSION", "Invalid evidence", sirius.DateString("2023-10-10")).
		Return(serverError)

	template := &mockTemplate{}

	form := url.Values{
		"decisionReason": {"Invalid evidence"},
		"decisionType":   {"DECLINED_REMISSION"},
		"decisionDate":   {"2023-10-10"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=765", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := AddFeeDecision(client, template.Func)(w, r)

	assert.Equal(t, serverError, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestAddFeeDecisionToDigitalLpaPostSuccess(t *testing.T) {
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
