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

func (m *mockCreateEpaClient) CreateEpa(ctx sirius.Context, donorID int, epa sirius.Epa) (sirius.Epa, error) {
	args := m.Called(ctx, donorID, epa)
	return args.Get(0).(sirius.Epa), args.Error(1)
}

func (m *mockCreateEpaClient) UpdateEpa(ctx sirius.Context, caseId int, epa sirius.Epa) error {
	args := m.Called(ctx, caseId, epa)
	return args.Error(0)
}

func (m *mockCreateEpaClient) Epa(ctx sirius.Context, id int) (sirius.Epa, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Epa), args.Error(1)
}

func TestGetCreateEpa(t *testing.T) {
	client := &mockCreateEpaClient{}

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createEpaData{
			Title: "Create an EPA",
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := CreateEpa(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetCreateEpaEdit(t *testing.T) {
	for _, appointmentType := range []string{"singular", "jointly", "jointly-and-severally"} {
		t.Run(appointmentType, func(t *testing.T) {
			epa := sirius.Epa{
				Case: sirius.Case{
					ReceiptDate:                     sirius.DateString("2022-04-05"),
					CaseAttorneySingular:            boolPtr(appointmentType == "singular"),
					CaseAttorneyJointlyAndSeverally: boolPtr(appointmentType == "jointly-and-severally"),
					CaseAttorneyJointly:             boolPtr(appointmentType == "jointly"),
				},
			}

			client := &mockCreateEpaClient{}

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, createEpaData{
					Title:           "Edit EPA",
					Epa:             epa,
					AppointmentType: appointmentType,
				}).
				Return(nil)

			client.
				On("Epa", mock.Anything, 234).
				Return(epa, nil)

			r, _ := http.NewRequest(http.MethodGet, "/?id=123&caseId=234", nil)
			w := httptest.NewRecorder()

			err := CreateEpa(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestGetCreateEpaBadQuery(t *testing.T) {
	testCases := map[string]string{
		"no-id":       "/",
		"bad-id":      "/?id=test",
		"bad-case-id": "/?id=123&caseId=test",
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

func TestGetCreateEpaEditWhenEpaErrors(t *testing.T) {
	client := &mockCreateEpaClient{}
	client.
		On("Epa", mock.Anything, 234).
		Return(sirius.Epa{}, errExample)

	template := &mockTemplate{}

	r, _ := http.NewRequest(http.MethodGet, "/?id=123&caseId=234", nil)
	w := httptest.NewRecorder()

	err := CreateEpa(client, template.Func)(w, r)

	assert.Equal(t, err, errExample)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateEpa(t *testing.T) {
	truePtr := true
	falsePtr := false
	dateString := "2022-04-05"
	epa := sirius.Epa{
		EpaDonorSignatureDate:   sirius.DateString(dateString),
		EpaDonorNoticeGivenDate: sirius.DateString(dateString),
		DonorHasOtherEpas:       &truePtr,
		OtherEpaInfo:            "More info",
		Case: sirius.Case{
			ReceiptDate:                     sirius.DateString(dateString),
			CaseAttorneySingular:            &truePtr,
			CaseAttorneyJointlyAndSeverally: &falsePtr,
			CaseAttorneyJointly:             &falsePtr,
			PaymentByCheque:                 &falsePtr,
			PaymentExemption:                &truePtr,
			PaymentDate:                     sirius.DateString(dateString),
		},
	}
	client := &mockCreateEpaClient{}
	client.
		On("CreateEpa", mock.Anything, 123, epa).
		Return(sirius.Epa{Case: sirius.Case{ID: 123}}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createEpaData{
			Title:           "Create an EPA",
			Success:         true,
			Epa:             epa,
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

func TestPostCreateEpaEdit(t *testing.T) {
	truePtr := true
	falsePtr := false
	dateString := "2022-04-05"
	epa := sirius.Epa{
		EpaDonorSignatureDate:   sirius.DateString(dateString),
		EpaDonorNoticeGivenDate: sirius.DateString(dateString),
		DonorHasOtherEpas:       &falsePtr,
		OtherEpaInfo:            "More info 1",
		Case: sirius.Case{
			ReceiptDate:                     sirius.DateString(dateString),
			CaseAttorneySingular:            &truePtr,
			CaseAttorneyJointlyAndSeverally: &falsePtr,
			CaseAttorneyJointly:             &falsePtr,
			PaymentByCheque:                 &falsePtr,
			PaymentExemption:                &truePtr,
			PaymentDate:                     sirius.DateString(dateString),
		},
	}
	newEpa := sirius.Epa{
		EpaDonorSignatureDate:   sirius.DateString(dateString),
		EpaDonorNoticeGivenDate: sirius.DateString(dateString),
		DonorHasOtherEpas:       &truePtr,
		OtherEpaInfo:            "More info 2",
		Case: sirius.Case{
			ReceiptDate:                     sirius.DateString(dateString),
			CaseAttorneySingular:            &truePtr,
			CaseAttorneyJointlyAndSeverally: &falsePtr,
			CaseAttorneyJointly:             &falsePtr,
			PaymentByCheque:                 &falsePtr,
			PaymentExemption:                &truePtr,
			PaymentDate:                     sirius.DateString(dateString),
		},
	}
	client := &mockCreateEpaClient{}
	client.
		On("Epa", mock.Anything, 234).
		Return(epa, nil)
	client.
		On("UpdateEpa", mock.Anything, 234, newEpa).
		Return(nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createEpaData{
			Title:           "Edit EPA",
			Success:         true,
			Epa:             newEpa,
			AppointmentType: "singular",
		}).
		Return(nil)

	form := url.Values{
		"epaDonorSignatureDate":   {dateString},
		"epaDonorNoticeGivenDate": {dateString},
		"donorHasOtherEpas":       {"true"},
		"otherEpaInfo":            {"More info 2"},
		"receiptDate":             {dateString},
		"caseAttorney":            {"singular"},
		"paymentByCheque":         {"false"},
		"paymentExemption":        {"true"},
		"paymentDate":             {dateString},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123&caseId=234", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateEpa(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateEpaAddAttorney(t *testing.T) {
	expectedError := RedirectError("/create-attorney?id=123&caseId=456")
	truePtr := true
	falsePtr := false
	dateString := "2022-04-05"
	epa := sirius.Epa{
		EpaDonorSignatureDate:   sirius.DateString(dateString),
		EpaDonorNoticeGivenDate: sirius.DateString(dateString),
		DonorHasOtherEpas:       &truePtr,
		OtherEpaInfo:            "More info",
		Case: sirius.Case{
			ReceiptDate:                     sirius.DateString(dateString),
			CaseAttorneySingular:            &truePtr,
			CaseAttorneyJointlyAndSeverally: &falsePtr,
			CaseAttorneyJointly:             &falsePtr,
			PaymentByCheque:                 &falsePtr,
			PaymentExemption:                &truePtr,
			PaymentDate:                     sirius.DateString(dateString),
		},
	}
	client := &mockCreateEpaClient{}
	client.
		On("CreateEpa", mock.Anything, 123, epa).
		Return(sirius.Epa{Case: sirius.Case{ID: 456}}, nil)

	template := &mockTemplate{}
	
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
		"addAttorney":             {"true"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateEpa(client, template.Func)(w, r)
	resp := w.Result()

	assert.Equal(t, err, expectedError)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateEpaWhenValidationError(t *testing.T) {
	expectedError := sirius.ValidationError{
		Field: sirius.FieldErrors{"field": {"": "problem"}},
	}

	truePtr := true
	falsePtr := false
	dateString := "2022-04-05"
	epa := sirius.Epa{
		EpaDonorSignatureDate:   sirius.DateString(dateString),
		EpaDonorNoticeGivenDate: sirius.DateString(dateString),
		DonorHasOtherEpas:       &truePtr,
		OtherEpaInfo:            "More info",
		Case: sirius.Case{
			CaseAttorneySingular:            &truePtr,
			CaseAttorneyJointlyAndSeverally: &falsePtr,
			CaseAttorneyJointly:             &falsePtr,
			PaymentByCheque:                 &falsePtr,
			PaymentExemption:                &truePtr,
			PaymentDate:                     sirius.DateString(dateString),
		},
	}

	client := &mockCreateEpaClient{}
	client.
		On("CreateEpa", mock.Anything, 123, epa).
		Return(sirius.Epa{}, expectedError)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createEpaData{
			Title:           "Create an EPA",
			Success:         false,
			Error:           expectedError,
			Epa:             epa,
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
