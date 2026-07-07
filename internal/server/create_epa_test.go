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
			DonorId: 123,
			Title:   "Create an EPA",
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := CreateEpa(client, template.Func, nil)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetCreateEpaHtmxRequest(t *testing.T) {
	client := &mockCreateEpaClient{}

	template := &mockTemplate{}
	partialTemplate := &mockTemplate{}
	partialTemplate.
		On("Func", mock.Anything, createEpaData{
			DonorId: 123,
			Title:   "Create an EPA",
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	r.Header.Add("HX-Request", "true")
	w := httptest.NewRecorder()

	err := CreateEpa(client, template.Func, partialTemplate.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	template.AssertNotCalled(t, "Func", mock.Anything, mock.Anything)
	mock.AssertExpectationsForObjects(t, client, template, partialTemplate)
}

func TestGetCreateEpaEdit(t *testing.T) {
	for _, appointmentType := range []string{"singular", "jointly", "jointly-and-severally"} {
		t.Run(appointmentType, func(t *testing.T) {
			epa := sirius.Epa{
				Case: sirius.Case{
					ReceiptDate:                     sirius.DateString("2022-04-05"),
					CaseAttorneySingular:            shared.BoolPtr(appointmentType == "singular"),
					CaseAttorneyJointlyAndSeverally: shared.BoolPtr(appointmentType == "jointly-and-severally"),
					CaseAttorneyJointly:             shared.BoolPtr(appointmentType == "jointly"),
				},
			}

			client := &mockCreateEpaClient{}

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, createEpaData{
					DonorId:         123,
					Title:           "Edit EPA",
					Epa:             epa,
					AppointmentType: appointmentType,
					CaseId:          234,
					IsUpdate:        true,
				}).
				Return(nil)

			client.
				On("Epa", mock.Anything, 234).
				Return(epa, nil)

			r, _ := http.NewRequest(http.MethodGet, "/?id=123&caseId=234", nil)
			w := httptest.NewRecorder()

			err := CreateEpa(client, template.Func, nil)(w, r)
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

			err := CreateEpa(nil, nil, nil)(w, r)

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

	err := CreateEpa(client, template.Func, nil)(w, r)

	assert.Equal(t, err, errExample)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateEpa(t *testing.T) {
	truePtr := shared.BoolPtr(true)
	falsePtr := shared.BoolPtr(false)
	dateString := "2022-04-05"
	epa := sirius.Epa{
		EpaDonorSignatureDate:   sirius.DateString(dateString),
		EpaDonorNoticeGivenDate: sirius.DateString(dateString),
		DonorHasOtherEpas:       truePtr,
		OtherEpaInfo:            "More info",
		Case: sirius.Case{
			ReceiptDate:                     sirius.DateString(dateString),
			CaseAttorneySingular:            truePtr,
			CaseAttorneyJointlyAndSeverally: falsePtr,
			CaseAttorneyJointly:             falsePtr,
			PaymentByCheque:                 falsePtr,
			PaymentExemption:                truePtr,
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
			DonorId:         123,
			Title:           "Create an EPA",
			Success:         true,
			AppointmentType: "singular",
			CaseId:          123,
			Epa:             sirius.Epa{Case: sirius.Case{ID: 123}},
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

	err := CreateEpa(client, template.Func, nil)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateEpaEdit(t *testing.T) {
	truePtr := shared.BoolPtr(true)
	falsePtr := shared.BoolPtr(false)
	dateString := "2022-04-05"
	epa := sirius.Epa{
		EpaDonorSignatureDate:   sirius.DateString(dateString),
		EpaDonorNoticeGivenDate: sirius.DateString(dateString),
		DonorHasOtherEpas:       falsePtr,
		OtherEpaInfo:            "",
		Case: sirius.Case{
			ReceiptDate:                     sirius.DateString(dateString),
			CaseAttorneySingular:            truePtr,
			CaseAttorneyJointlyAndSeverally: falsePtr,
			CaseAttorneyJointly:             falsePtr,
			PaymentByCheque:                 falsePtr,
			PaymentExemption:                truePtr,
			PaymentDate:                     sirius.DateString(dateString),
		},
	}
	newEpa := sirius.Epa{
		EpaDonorSignatureDate:   sirius.DateString(dateString),
		EpaDonorNoticeGivenDate: sirius.DateString(dateString),
		DonorHasOtherEpas:       truePtr,
		OtherEpaInfo:            "More info 2",
		Case: sirius.Case{
			ReceiptDate:                     sirius.DateString(dateString),
			CaseAttorneySingular:            truePtr,
			CaseAttorneyJointlyAndSeverally: falsePtr,
			CaseAttorneyJointly:             falsePtr,
			PaymentByCheque:                 falsePtr,
			PaymentExemption:                truePtr,
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
			DonorId:         123,
			Title:           "Edit EPA",
			Success:         true,
			Epa:             epa,
			AppointmentType: "singular",
			CaseId:          234,
			IsUpdate:        true,
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

	err := CreateEpa(client, template.Func, nil)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateEpaRedirects(t *testing.T) {
	falsePtr := shared.BoolPtr(false)

	baseEpa := sirius.Epa{
		DonorHasOtherEpas: falsePtr,
		Case: sirius.Case{
			CaseAttorneySingular:            falsePtr,
			CaseAttorneyJointlyAndSeverally: falsePtr,
			CaseAttorneyJointly:             falsePtr,
			PaymentByCheque:                 falsePtr,
			PaymentExemption:                falsePtr,
		},
	}

	testCases := map[string]struct {
		query         string
		formField     string
		formValue     string
		expectedError error
		setupMocks    func(client *mockCreateEpaClient, epa sirius.Epa)
	}{
		"add attorney": {
			query:         "/?id=123",
			formField:     "addAttorney",
			formValue:     "true",
			expectedError: RedirectError("/create-attorney?id=123&caseId=456"),
			setupMocks: func(client *mockCreateEpaClient, epa sirius.Epa) {
				client.
					On("CreateEpa", mock.Anything, 123, epa).
					Return(sirius.Epa{Case: sirius.Case{ID: 456}}, nil)
			},
		},
		"update attorney": {
			query:         "/?id=123",
			formField:     "updateAttorney",
			formValue:     "789",
			expectedError: RedirectError("/create-attorney?id=123&caseId=456&attorneyId=789"),
			setupMocks: func(client *mockCreateEpaClient, epa sirius.Epa) {
				client.
					On("CreateEpa", mock.Anything, 123, epa).
					Return(sirius.Epa{Case: sirius.Case{ID: 456}}, nil)
			},
		},
		"add correspondent": {
			query:         "/?id=123",
			formField:     "addCorrespondent",
			formValue:     "true",
			expectedError: RedirectError("/create-correspondent?id=123&caseId=456"),
			setupMocks: func(client *mockCreateEpaClient, epa sirius.Epa) {
				client.
					On("CreateEpa", mock.Anything, 123, epa).
					Return(sirius.Epa{Case: sirius.Case{ID: 456}}, nil)
			},
		},
		"add correspondent with attorney on case": {
			query:         "/?id=123&caseId=456",
			formField:     "addCorrespondent",
			formValue:     "true",
			expectedError: RedirectError("/select-or-create-correspondent?id=123&caseId=456"),
			setupMocks: func(client *mockCreateEpaClient, epa sirius.Epa) {
				client.
					On("Epa", mock.Anything, 456).
					Return(sirius.Epa{
						Case: sirius.Case{
							ID:        456,
							Attorneys: []sirius.Attorney{{Person: sirius.Person{ID: 1}}},
						},
					}, nil)
				client.
					On("UpdateEpa", mock.Anything, 456, epa).
					Return(nil)
			},
		},
		"update correspondent": {
			query:         "/?id=123&caseId=456",
			formField:     "updateCorrespondent",
			formValue:     "1",
			expectedError: RedirectError("/create-correspondent?id=123&caseId=456"),
			setupMocks: func(client *mockCreateEpaClient, epa sirius.Epa) {
				client.
					On("Epa", mock.Anything, 456).
					Return(sirius.Epa{Case: sirius.Case{ID: 456}}, nil)
				client.
					On("UpdateEpa", mock.Anything, 456, epa).
					Return(nil)
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			epa := baseEpa
			client := &mockCreateEpaClient{}
			template := &mockTemplate{}

			tc.setupMocks(client, epa)

			form := url.Values{
				tc.formField: {tc.formValue},
			}

			r, _ := http.NewRequest(http.MethodPost, tc.query, strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := CreateEpa(client, template.Func, nil)(w, r)
			resp := w.Result()

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostCreateEpaWhenValidationError(t *testing.T) {
	expectedError := sirius.ValidationError{
		Field: sirius.FieldErrors{"field": {"": "problem"}},
	}

	truePtr := shared.BoolPtr(true)
	falsePtr := shared.BoolPtr(false)
	dateString := "2022-04-05"
	epa := sirius.Epa{
		EpaDonorSignatureDate:   sirius.DateString(dateString),
		EpaDonorNoticeGivenDate: sirius.DateString(dateString),
		DonorHasOtherEpas:       truePtr,
		OtherEpaInfo:            "More info",
		Case: sirius.Case{
			CaseAttorneySingular:            truePtr,
			CaseAttorneyJointlyAndSeverally: falsePtr,
			CaseAttorneyJointly:             falsePtr,
			PaymentByCheque:                 falsePtr,
			PaymentExemption:                truePtr,
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
			DonorId:         123,
			Title:           "Create an EPA",
			Success:         false,
			Error:           expectedError,
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

	err := CreateEpa(client, template.Func, nil)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateEpaWhenValidationErrorHtmxRequest(t *testing.T) {
	expectedError := sirius.ValidationError{
		Field: sirius.FieldErrors{"field": {"": "problem"}},
	}

	falsePtr := shared.BoolPtr(false)
	epa := sirius.Epa{
		DonorHasOtherEpas: falsePtr,
		Case: sirius.Case{
			CaseAttorneySingular:            falsePtr,
			CaseAttorneyJointlyAndSeverally: falsePtr,
			CaseAttorneyJointly:             falsePtr,
			PaymentByCheque:                 falsePtr,
			PaymentExemption:                falsePtr,
		},
	}

	client := &mockCreateEpaClient{}
	client.
		On("CreateEpa", mock.Anything, 123, epa).
		Return(sirius.Epa{}, expectedError)

	template := &mockTemplate{}
	partialTemplate := &mockTemplate{}
	partialTemplate.
		On("Func", mock.Anything, createEpaData{
			DonorId: 123,
			Title:   "Create an EPA",
			Success: false,
			Error:   expectedError,
		}).
		Return(nil)

	form := url.Values{}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	r.Header.Add("HX-Request", "true")
	w := httptest.NewRecorder()

	err := CreateEpa(client, template.Func, partialTemplate.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	template.AssertNotCalled(t, "Func", mock.Anything, mock.Anything)
	mock.AssertExpectationsForObjects(t, client, template, partialTemplate)
}

func TestPostCreateEpaWhenValidationErrorOnReceiptDate(t *testing.T) {
	returnedError := sirius.ValidationError{
		Field: sirius.FieldErrors{"receiptDate": {"receiptDate": "problem"}},
	}
	expectedError := sirius.ValidationError{
		Field: sirius.FieldErrors{"receiptDate": {"receiptDate": "Enter or select a receipt date to save and exit"}},
	}

	truePtr := shared.BoolPtr(true)
	falsePtr := shared.BoolPtr(false)
	dateString := "2022-04-05"
	epa := sirius.Epa{
		EpaDonorSignatureDate:   sirius.DateString(dateString),
		EpaDonorNoticeGivenDate: sirius.DateString(dateString),
		DonorHasOtherEpas:       truePtr,
		OtherEpaInfo:            "More info",
		Case: sirius.Case{
			CaseAttorneySingular:            truePtr,
			CaseAttorneyJointlyAndSeverally: falsePtr,
			CaseAttorneyJointly:             falsePtr,
			PaymentByCheque:                 falsePtr,
			PaymentExemption:                truePtr,
			PaymentDate:                     sirius.DateString(dateString),
		},
	}

	client := &mockCreateEpaClient{}
	client.
		On("CreateEpa", mock.Anything, 123, epa).
		Return(sirius.Epa{}, returnedError)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createEpaData{
			DonorId:         123,
			Title:           "Create an EPA",
			Success:         false,
			Error:           expectedError,
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

	err := CreateEpa(client, template.Func, nil)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateEpaAddActorValidationErrorOnReceiptDate(t *testing.T) {
	for _, actorType := range []string{"addAttorney", "addCorrespondent"} {
		t.Run(actorType, func(t *testing.T) {
			returnedError := sirius.ValidationError{
				Field: sirius.FieldErrors{"receiptDate": {"receiptDate": "problem"}},
			}
			expectedError := sirius.ValidationError{
				Field: sirius.FieldErrors{"receiptDate": {"receiptDate": "Enter or select a receipt date to continue to step 3"}},
			}

			truePtr := shared.BoolPtr(true)
			falsePtr := shared.BoolPtr(false)
			dateString := "2022-04-05"
			epa := sirius.Epa{
				EpaDonorSignatureDate:   sirius.DateString(dateString),
				EpaDonorNoticeGivenDate: sirius.DateString(dateString),
				DonorHasOtherEpas:       truePtr,
				OtherEpaInfo:            "More info",
				Case: sirius.Case{
					CaseAttorneySingular:            truePtr,
					CaseAttorneyJointlyAndSeverally: falsePtr,
					CaseAttorneyJointly:             falsePtr,
					PaymentByCheque:                 falsePtr,
					PaymentExemption:                truePtr,
					PaymentDate:                     sirius.DateString(dateString),
				},
			}

			client := &mockCreateEpaClient{}
			client.
				On("CreateEpa", mock.Anything, 123, epa).
				Return(sirius.Epa{}, returnedError)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, createEpaData{
					DonorId:         123,
					Title:           "Create an EPA",
					Success:         false,
					Error:           expectedError,
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
				actorType:                 {"true"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := CreateEpa(client, template.Func, nil)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestGetCreateEpaDoesNotSetIsUpdate(t *testing.T) {
	client := &mockCreateEpaClient{}
	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createEpaData{
			DonorId:  123,
			Title:    "Create an EPA",
			IsUpdate: false,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := CreateEpa(client, template.Func, nil)(w, r)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateEpaEditWhenValidationError(t *testing.T) {
	expectedError := sirius.ValidationError{
		Field: sirius.FieldErrors{"field": {"": "problem"}},
	}

	truePtr := shared.BoolPtr(true)
	falsePtr := shared.BoolPtr(false)
	dateString := "2022-04-05"
	existingEpa := sirius.Epa{
		Case: sirius.Case{ReceiptDate: sirius.DateString(dateString)},
	}
	submittedEpa := sirius.Epa{
		EpaDonorSignatureDate:   sirius.DateString(dateString),
		EpaDonorNoticeGivenDate: sirius.DateString(dateString),
		DonorHasOtherEpas:       truePtr,
		OtherEpaInfo:            "More info",
		Case: sirius.Case{
			ReceiptDate:                     sirius.DateString(dateString),
			CaseAttorneySingular:            truePtr,
			CaseAttorneyJointlyAndSeverally: falsePtr,
			CaseAttorneyJointly:             falsePtr,
			PaymentByCheque:                 falsePtr,
			PaymentExemption:                truePtr,
			PaymentDate:                     sirius.DateString(dateString),
		},
	}

	client := &mockCreateEpaClient{}
	client.
		On("Epa", mock.Anything, 234).
		Return(existingEpa, nil)
	client.
		On("UpdateEpa", mock.Anything, 234, submittedEpa).
		Return(expectedError)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createEpaData{
			DonorId:         123,
			Title:           "Edit EPA",
			IsUpdate:        true,
			Success:         false,
			Epa:             existingEpa,
			Error:           expectedError,
			AppointmentType: "singular",
			CaseId:          234,
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

	r, _ := http.NewRequest(http.MethodPost, "/?id=123&caseId=234", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateEpa(client, template.Func, nil)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}
