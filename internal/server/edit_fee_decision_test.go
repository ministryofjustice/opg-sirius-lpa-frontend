package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockEditFeeDecisionClient struct {
	mock.Mock
}

func (m *mockEditFeeDecisionClient) Case(ctx sirius.Context, id int) (sirius.Case, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Case), args.Error(1)
}

func (m *mockEditFeeDecisionClient) DeclineFeeReductions(ctx sirius.Context, id int) ([]sirius.DeclineFeeReductions, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).([]sirius.DeclineFeeReductions), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockEditFeeDecisionClient) EditFeeDecision(ctx sirius.Context, caseID int, feeDecisionID int, decisionType string, decisionReason string, decisionDate sirius.DateString) error {
	return m.Called(ctx, caseID, feeDecisionID, decisionType, decisionReason, decisionDate).Error(0)
}

func (m *mockEditFeeDecisionClient) RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error) {
	args := m.Called(ctx, category)
	if args.Get(0) != nil {
		return args.Get(0).([]sirius.RefDataItem), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestGetEditFeeDecision(t *testing.T) {
	caseItem := sirius.Case{ID: 4, UID: "7000-0000-0021"}
	decisions := []sirius.DeclineFeeReductions{
		{
			ID:             123,
			DecisionDate:   "2022-08-23",
			DecisionType:   shared.DecisionTypeExemption,
			DecisionReason: "Insufficient evidence",
		},
	}
	decisionTypes := []sirius.RefDataItem{{Handle: "DECLINED_EXEMPTION", Label: "Declined exemption"}}

	client := &mockEditFeeDecisionClient{}
	client.
		On("Case", mock.Anything, 4).
		Return(caseItem, nil).
		On("DeclineFeeReductions", mock.Anything, 4).
		Return(decisions, nil).
		On("RefDataByCategory", mock.Anything, sirius.FeeDecisionTypeCategory).
		Return(decisionTypes, nil)

	tmpl := &mockTemplate{}
	tmpl.
		On("Func", mock.Anything, editFeeDecisionData{
			Case:           caseItem,
			DecisionTypes:  decisionTypes,
			DecisionType:   "DECLINED_EXEMPTION",
			DecisionReason: "Insufficient evidence",
			DecisionDate:   "2022-08-23",
			FeeDecisionID:  123,
			IsEditing:      true,
			ReturnUrl:      "/payments/4",
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?caseId=4&feeDecisionId=123", nil)
	w := httptest.NewRecorder()

	err := EditFeeDecision(client, tmpl.Func)(w, r)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	mock.AssertExpectationsForObjects(t, client, tmpl)
}

func TestEditFeeDecisionOnPostSuccess(t *testing.T) {
	caseItem := sirius.Case{ID: 4, UID: "M-AAAA-BBBB-DDDD", CaseType: "DIGITAL_LPA"}
	decisions := []sirius.DeclineFeeReductions{{ID: 123, DecisionType: shared.DecisionTypeExemption}}

	client := &mockEditFeeDecisionClient{}
	client.
		On("Case", mock.Anything, 4).
		Return(caseItem, nil).
		On("DeclineFeeReductions", mock.Anything, 4).
		Return(decisions, nil).
		On("RefDataByCategory", mock.Anything, sirius.FeeDecisionTypeCategory).
		Return(feeDecisionTypes, nil).
		On("EditFeeDecision", mock.Anything, 4, 123, "DECLINED_REMISSION", "Updated reason", sirius.DateString("2023-10-10")).
		Return(nil)

	form := url.Values{
		"decisionType":   {"DECLINED_REMISSION"},
		"decisionReason": {"Updated reason"},
		"decisionDate":   {"2023-10-10"},
	}
	r, _ := http.NewRequest(http.MethodPost, "/?caseId=4&feeDecisionId=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := EditFeeDecision(client, nil)(w, r)

	assert.Equal(t, RedirectError("/lpa/M-AAAA-BBBB-DDDD/payments"), err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestEditFeeDecisionOnValidationError(t *testing.T) {
	caseItem := sirius.Case{ID: 4, UID: "7000-0000-0021"}
	decisions := []sirius.DeclineFeeReductions{{ID: 123, DecisionType: shared.DecisionTypeExemption}}
	validationError := sirius.ValidationError{
		Field: sirius.FieldErrors{
			"decisionReason": {"isInvalid": "Reason is invalid"},
		},
	}

	client := &mockEditFeeDecisionClient{}
	client.
		On("Case", mock.Anything, 4).
		Return(caseItem, nil).
		On("DeclineFeeReductions", mock.Anything, 4).
		Return(decisions, nil).
		On("RefDataByCategory", mock.Anything, sirius.FeeDecisionTypeCategory).
		Return(feeDecisionTypes, nil).
		On("EditFeeDecision", mock.Anything, 4, 123, "DECLINED_REMISSION", "Updated reason", sirius.DateString("2023-10-10")).
		Return(validationError)

	tmpl := &mockTemplate{}
	tmpl.
		On("Func", mock.Anything, mock.MatchedBy(func(data editFeeDecisionData) bool {
			return data.Case.ID == caseItem.ID &&
				data.ReturnUrl == "/payments/4" &&
				data.DecisionType == "DECLINED_REMISSION" &&
				data.DecisionReason == "Updated reason" &&
				data.DecisionDate == "2023-10-10" &&
				data.FeeDecisionID == 123 &&
				data.Error.Detail == validationError.Detail &&
				data.Error.Field["decisionReason"]["isInvalid"] == "Reason is invalid"
		})).
		Return(nil)

	form := url.Values{
		"decisionType":   {"DECLINED_REMISSION"},
		"decisionReason": {"Updated reason"},
		"decisionDate":   {"2023-10-10"},
	}
	r, _ := http.NewRequest(http.MethodPost, "/?caseId=4&feeDecisionId=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := EditFeeDecision(client, tmpl.Func)(w, r)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	mock.AssertExpectationsForObjects(t, client, tmpl)
}

func TestPostEditFeeDecisionHTMX(t *testing.T) {
	caseItem := sirius.Case{ID: 4, UID: "M-AAAA-BBBB-DDDD", CaseType: "DIGITAL_LPA"}
	decisions := []sirius.DeclineFeeReductions{{ID: 123, DecisionType: shared.DecisionTypeExemption}}

	client := &mockEditFeeDecisionClient{}
	client.
		On("Case", mock.Anything, 4).
		Return(caseItem, nil).
		On("DeclineFeeReductions", mock.Anything, 4).
		Return(decisions, nil).
		On("RefDataByCategory", mock.Anything, sirius.FeeDecisionTypeCategory).
		Return(feeDecisionTypes, nil).
		On("EditFeeDecision", mock.Anything, 4, 123, "DECLINED_REMISSION", "Updated reason", sirius.DateString("2023-10-10")).
		Return(nil)

	tmpl := &mockTemplate{}
	tmpl.
		On("Func", mock.Anything, editFeeDecisionData{
			Case:           caseItem,
			DecisionTypes:  feeDecisionTypes,
			DecisionType:   "DECLINED_REMISSION",
			DecisionReason: "Updated reason",
			DecisionDate:   "2023-10-10",
			FeeDecisionID:  123,
			IsEditing:      true,
			IsPartial:      true,
			ReturnUrl:      "/lpa/M-AAAA-BBBB-DDDD/payments",
			HtmxRedirect:   "/lpa/M-AAAA-BBBB-DDDD/payments",
		}).
		Return(nil)

	form := url.Values{
		"decisionType":   {"DECLINED_REMISSION"},
		"decisionReason": {"Updated reason"},
		"decisionDate":   {"2023-10-10"},
	}
	r, _ := http.NewRequest(http.MethodPost, "/?caseId=4&feeDecisionId=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	r.Header.Add("HX-Request", "true")
	w := httptest.NewRecorder()

	err := EditFeeDecision(client, tmpl.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, tmpl)
}

func TestEditFeeDecisionNotFound(t *testing.T) {
	client := &mockEditFeeDecisionClient{}
	client.
		On("Case", mock.Anything, 4).
		Return(sirius.Case{ID: 4}, nil).
		On("DeclineFeeReductions", mock.Anything, 4).
		Return([]sirius.DeclineFeeReductions{}, nil).
		On("RefDataByCategory", mock.Anything, sirius.FeeDecisionTypeCategory).
		Return(feeDecisionTypes, nil)

	r, _ := http.NewRequest(http.MethodGet, "/?caseId=4&feeDecisionId=123", nil)
	w := httptest.NewRecorder()

	err := EditFeeDecision(client, nil)(w, r)

	assert.Equal(t, sirius.StatusError{Code: http.StatusNotFound}, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestEditFeeDecisionWhenRequestFails(t *testing.T) {
	for _, failedRequest := range []string{"Case", "DeclineFeeReductions", "RefDataByCategory"} {
		t.Run(failedRequest, func(t *testing.T) {
			client := &mockEditFeeDecisionClient{}
			client.
				On("Case", mock.Anything, 4).
				Return(sirius.Case{ID: 4}, errorForRequest(failedRequest, "Case")).
				On("DeclineFeeReductions", mock.Anything, 4).
				Return([]sirius.DeclineFeeReductions{}, errorForRequest(failedRequest, "DeclineFeeReductions")).
				On("RefDataByCategory", mock.Anything, sirius.FeeDecisionTypeCategory).
				Return(feeDecisionTypes, errorForRequest(failedRequest, "RefDataByCategory"))

			r, _ := http.NewRequest(http.MethodGet, "/?caseId=4&feeDecisionId=123", nil)
			w := httptest.NewRecorder()

			err := EditFeeDecision(client, nil)(w, r)

			assert.Equal(t, errExample, err)
			mock.AssertExpectationsForObjects(t, client)
		})
	}
}

func errorForRequest(failedRequest string, request string) error {
	if failedRequest == request {
		return errExample
	}

	return nil
}
