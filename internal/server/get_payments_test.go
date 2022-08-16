package server

import (
	"errors"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockGetPayments struct {
	mock.Mock
}

func (m *mockGetPayments) Payments(ctx sirius.Context, id int) ([]sirius.Payment, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]sirius.Payment), args.Error(1)
}

func (m *mockGetPayments) Case(ctx sirius.Context, id int) (sirius.Case, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Case), args.Error(1)
}

func TestGetPayments(t *testing.T) {
	payments := []sirius.Payment{
		{
			ID:     2,
			CaseID: 4,
			Amount: "41.00",
			Locked: false,
		},
		{
			ID:     3,
			CaseID: 4,
			Amount: "14.38",
			Locked: false,
		},
	}

	caseItem := sirius.Case{
		UID:     "7000-0000-0021",
		SubType: "pfa",
	}

	client := &mockGetPayments{}
	client.
		On("Payments", mock.Anything, 4).
		Return(payments, nil)

	client.
		On("Case", mock.Anything, 4).
		Return(caseItem, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, getPaymentsData{
			Payments:  payments,
			Case:      caseItem,
			TotalPaid: "55.38",
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=4", nil)
	w := httptest.NewRecorder()

	err := GetPayments(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetPaymentsNoID(t *testing.T) {
	testCases := map[string]string{
		"no-id":  "/",
		"bad-id": "/?id=test",
	}

	for name, testUrl := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, testUrl, nil)
			w := httptest.NewRecorder()

			err := GetPayments(nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestGetPaymentsWhenFailure(t *testing.T) {
	expectedError := errors.New("err")

	client := &mockGetPayments{}
	client.
		On("Case", mock.Anything, 4).
		Return(sirius.Case{}, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=4", nil)
	w := httptest.NewRecorder()

	err := GetPayments(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetPaymentsWhenTemplateErrors(t *testing.T) {
	payments := []sirius.Payment{
		{
			ID:     2,
			CaseID: 4,
			Source: sirius.PaymentSource{
				Name:  "Phone",
				Value: "PHONE",
			},
			Amount:      "41.00",
			PaymentDate: sirius.DateString("2022-08-23T14:55:20+00:00"),
			Type: sirius.TypeOfPayment{
				Name:  "Card",
				Value: "CARD",
			},
			CreatedDate: sirius.DateString("2022-08-24T14:55:20+00:00"),
			Locked:      false,
		},
	}

	caseItem := sirius.Case{
		UID:     "7000-0000-0021",
		SubType: "pfa",
	}

	client := &mockGetPayments{}
	client.
		On("Payments", mock.Anything, 4).
		Return(payments, nil)

	client.
		On("Case", mock.Anything, 4).
		Return(caseItem, nil)

	expectedError := errors.New("err")

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, getPaymentsData{
			Payments:  payments,
			Case:      caseItem,
			TotalPaid: "41.00",
		}).
		Return(expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=4", nil)
	w := httptest.NewRecorder()

	err := GetPayments(client, template.Func)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client, template)
}
