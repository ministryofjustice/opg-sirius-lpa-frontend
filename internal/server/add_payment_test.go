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

type mockAddPaymentClient struct {
	mock.Mock
}

func (m *mockAddPaymentClient) AddPayment(ctx sirius.Context, caseID int, amount int, source string, paymentDate sirius.DateString) error {
	return m.Called(ctx, caseID, amount, source, paymentDate).Error(0)
}

func (m *mockAddPaymentClient) Case(ctx sirius.Context, id int) (sirius.Case, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Case), args.Error(1)
}

func (m *mockAddPaymentClient) RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error) {
	args := m.Called(ctx, category)
	if args.Get(0) != nil {
		return args.Get(0).([]sirius.RefDataItem), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestGetAddPayment(t *testing.T) {
	caseItem := sirius.Case{
		UID:     "7000-0000-0021",
		SubType: "pfa",
	}

	paymentSources := []sirius.RefDataItem{
		{
			Handle:         "PHONE",
			Label:          "Paid over the phone",
			UserSelectable: true,
		},
	}

	client := &mockAddPaymentClient{}
	client.
		On("Case", mock.Anything, 4).
		Return(caseItem, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.PaymentSourceCategory).
		Return(paymentSources, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, addPaymentData{
			Case:           caseItem,
			PaymentSources: paymentSources,
			ReturnUrl: "/payments/4",
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=4", nil)
	w := httptest.NewRecorder()

	err := AddPayment(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestAddPaymentNoID(t *testing.T) {
	testCases := map[string]string{
		"no-id":  "/",
		"bad-id": "/?id=test",
	}

	for name, testUrl := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, testUrl, nil)
			w := httptest.NewRecorder()

			err := AddPayment(nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestAddPaymentWhenFailureOnGetCase(t *testing.T) {
	paymentSources := []sirius.RefDataItem{
		{
			Handle:         "PHONE",
			Label:          "Paid over the phone",
			UserSelectable: true,
		},
	}

	client := &mockAddPaymentClient{}
	client.
		On("Case", mock.Anything, 4).
		Return(sirius.Case{}, expectedError)
	client.
		On("RefDataByCategory", mock.Anything, sirius.PaymentSourceCategory).
		Return(paymentSources, nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=4", nil)
	w := httptest.NewRecorder()

	err := AddPayment(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestAddPaymentWhenTemplateErrors(t *testing.T) {
	caseItem := sirius.Case{
		UID:     "7000-0000-0021",
		SubType: "pfa",
	}

	paymentSources := []sirius.RefDataItem{
		{
			Handle:         "PHONE",
			Label:          "Paid over the phone",
			UserSelectable: true,
		},
	}

	client := &mockAddPaymentClient{}
	client.
		On("Case", mock.Anything, 123).
		Return(caseItem, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.PaymentSourceCategory).
		Return(paymentSources, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, addPaymentData{
			Case:           caseItem,
			PaymentSources: paymentSources,
			ReturnUrl: "/payments/123",
		}).
		Return(expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := AddPayment(client, template.Func)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestAddPaymentWhenFailureOnGetPaymentSourceRefData(t *testing.T) {
	caseItem := sirius.Case{
		UID:     "7000-0000-0021",
		SubType: "pfa",
	}

	client := &mockAddPaymentClient{}
	client.
		On("Case", mock.Anything, 123).
		Return(caseItem, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.PaymentSourceCategory).
		Return([]sirius.RefDataItem{}, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := AddPayment(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestPostAddPayment(t *testing.T) {
	caseitem := sirius.Case{CaseType: "lpa", UID: "700700"}

	paymentSources := []sirius.RefDataItem{
		{
			Handle:         "PHONE",
			Label:          "Paid over the phone",
			UserSelectable: true,
		},
	}

	client := &mockAddPaymentClient{}
	client.
		On("AddPayment", mock.Anything, 123, 4100, "MAKE", sirius.DateString("2022-01-23")).
		Return(nil)
	client.
		On("Case", mock.Anything, 123).
		Return(caseitem, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.PaymentSourceCategory).
		Return(paymentSources, nil)

	template := &mockTemplate{}

	form := url.Values{
		"amount":      {"41.00"},
		"source":      {"MAKE"},
		"paymentDate": {"2022-01-23"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := AddPayment(client, template.Func)(w, r)
	resp := w.Result()

	assert.Equal(t, RedirectError("/payments/123"), err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostAddPaymentAmountIncorrectFormat(t *testing.T) {
	for _, amount := range []string{"41", "41.5", "41.555", ".45"} {
		t.Run(amount, func(t *testing.T) {
			caseitem := sirius.Case{CaseType: "lpa", UID: "700700"}

			paymentSources := []sirius.RefDataItem{
				{
					Handle:         "PHONE",
					Label:          "Paid over the phone",
					UserSelectable: true,
				},
			}

			client := &mockAddPaymentClient{}
			client.
				On("Case", mock.Anything, 123).
				Return(caseitem, nil)
			client.
				On("RefDataByCategory", mock.Anything, sirius.PaymentSourceCategory).
				Return(paymentSources, nil)

			validationError := sirius.ValidationError{
				Field: sirius.FieldErrors{
					"amount": {"reason": "Please enter the amount to 2 decimal places"},
				},
			}

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, addPaymentData{
					Case:           caseitem,
					Amount:         amount,
					Source:         "MAKE",
					PaymentDate:    sirius.DateString("2022-01-23"),
					Error:          validationError,
					PaymentSources: paymentSources,
					ReturnUrl:      "/payments/123",
				}).
				Return(nil)

			form := url.Values{
				"amount":      {amount},
				"source":      {"MAKE"},
				"paymentDate": {"2022-01-23"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := AddPayment(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostAddPaymentToDigitalLpa(t *testing.T) {
	caseitem := sirius.Case{CaseType: "DIGITAL_LPA", UID: "M-AAAA-BBBB-CCCC"}

	paymentSources := []sirius.RefDataItem{
		{
			Handle:         "PHONE",
			Label:          "Paid over the phone",
			UserSelectable: true,
		},
	}

	client := &mockAddPaymentClient{}
	client.
		On("AddPayment", mock.Anything, 444, 5200, "PHONE", sirius.DateString("2023-08-31")).
		Return(nil)
	client.
		On("Case", mock.Anything, 444).
		Return(caseitem, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.PaymentSourceCategory).
		Return(paymentSources, nil)

	template := &mockTemplate{}

	form := url.Values{
		"amount":      {"52.00"},
		"source":      {"PHONE"},
		"paymentDate": {"2023-08-31"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=444", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := AddPayment(client, template.Func)(w, r)
	resp := w.Result()

	assert.Equal(t, RedirectError("/lpa/M-AAAA-BBBB-CCCC/payments"), err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}
