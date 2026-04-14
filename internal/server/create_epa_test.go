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

type mockCreateEpaClient struct {
	mock.Mock
}

func (m *mockCreateEpaClient) CreateEpa(ctx sirius.Context, donorID int, epa sirius.Case) error {
	args := m.Called(ctx, donorID, epa)
	return args.Error(0)
}

func TestGetCreateEpa(t *testing.T) {
	client := &mockCreateEpaClient{}

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createEpaData{}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := CreateEpa(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetCreateEpaBadQuery(t *testing.T) {
	testCases := map[string]string{
		"no-id":  "/",
		"bad-id": "/?id=test",
	}

	for name, url := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			err := CreateEpa(nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestPostCreateEpa(t *testing.T) {
	truePtr := true
	falsePtr := false
	dateString := "2022-04-05"
	epa := sirius.Case{
		EpaDonorSignatureDate:           sirius.DateString(dateString),
		EpaDonorNoticeGivenDate:         sirius.DateString(dateString),
		DonorHasOtherEpas:               &truePtr,
		OtherEpaInfo:                    "More info",
		ReceiptDate:                     sirius.DateString(dateString),
		CaseAttorneySingular:            &truePtr,
		CaseAttorneyJointlyAndSeverally: &falsePtr,
		CaseAttorneyJointly:             &falsePtr,
		PaymentByCheque:                 &falsePtr,
		PaymentExemption:                &truePtr,
		PaymentDate:                     sirius.DateString(dateString),
	}
	client := &mockCreateEpaClient{}
	client.
		On("CreateEpa", mock.Anything, 123, epa).
		Return(nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createEpaData{
			Success:         true,
			Case:            epa,
			AppointmentType: "singular",
		}).
		Return(nil)

	form := url.Values{
		"epaDonorSignatureDate":   {dateString},
		"epaDonorNoticeGivenDate": {dateString},
		"donorHasOtherEpas":       {"true"},
		"otherEpaInfo":            {"More info"},
		"receiptDate":             {dateString},
		"caseAttorney":            {"singular"},
		"paymentByCheque":         {"false"},
		"paymentExemption":        {"true"},
		"paymentDate":             {dateString},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateEpa(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostAddComplaintWhenCreateEpaValidationError(t *testing.T) {
	expectedError := sirius.ValidationError{
		Field: sirius.FieldErrors{"field": {"": "problem"}},
	}

	truePtr := true
	falsePtr := false
	dateString := "2022-04-05"
	epa := sirius.Case{
		EpaDonorSignatureDate:           sirius.DateString(dateString),
		EpaDonorNoticeGivenDate:         sirius.DateString(dateString),
		DonorHasOtherEpas:               &truePtr,
		OtherEpaInfo:                    "More info",
		CaseAttorneySingular:            &truePtr,
		CaseAttorneyJointlyAndSeverally: &falsePtr,
		CaseAttorneyJointly:             &falsePtr,
		PaymentByCheque:                 &falsePtr,
		PaymentExemption:                &truePtr,
		PaymentDate:                     sirius.DateString(dateString),
	}

	client := &mockCreateEpaClient{}
	client.
		On("CreateEpa", mock.Anything, 123, epa).
		Return(expectedError)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createEpaData{
			Success:         false,
			Error:           expectedError,
			Case:            epa,
			AppointmentType: "singular",
		}).
		Return(nil)

	form := url.Values{
		"epaDonorSignatureDate":   {dateString},
		"epaDonorNoticeGivenDate": {dateString},
		"donorHasOtherEpas":       {"true"},
		"otherEpaInfo":            {"More info"},
		"caseAttorney":            {"singular"},
		"paymentByCheque":         {"false"},
		"paymentExemption":        {"true"},
		"paymentDate":             {dateString},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateEpa(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}
