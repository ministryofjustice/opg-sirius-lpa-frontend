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

func TestGetCreateEpa(t *testing.T) {
	client := &mockCreateEpaClient{}

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, epaData{
			DonorId: 123,
			Title:   "Create an EPA",
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

func TestGetCreateEpaHtmxRequest(t *testing.T) {
	client := &mockCreateEpaClient{}

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, epaData{
			DonorId:   123,
			Title:     "Create an EPA",
			IsPartial: true,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	r.Header.Add("HX-Request", "true")
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
		On("Func", mock.Anything, epaData{
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

	err := CreateEpa(client, template.Func)(w, r)
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
		formField     string
		expectedError error
	}{
		"add attorney": {
			formField:     "addAttorney",
			expectedError: RedirectError("/create-attorney?id=123&caseId=456&caseType=epa"),
		},
		"add correspondent": {
			formField:     "addCorrespondent",
			expectedError: RedirectError("/create-correspondent?id=123&caseId=456&caseType=epa"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			epa := baseEpa
			client := &mockCreateEpaClient{}
			template := &mockTemplate{}

			client.
				On("CreateEpa", mock.Anything, 123, epa).
				Return(sirius.Epa{Case: sirius.Case{ID: 456}}, nil)

			form := url.Values{
				tc.formField: {"true"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/create-epa?id=123", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := CreateEpa(client, template.Func)(w, r)
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
		On("Func", mock.Anything, epaData{
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

	err := CreateEpa(client, template.Func)(w, r)
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
	template.
		On("Func", mock.Anything, epaData{
			DonorId:   123,
			Title:     "Create an EPA",
			Success:   false,
			Error:     expectedError,
			IsPartial: true,
		}).
		Return(nil)

	form := url.Values{}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	r.Header.Add("HX-Request", "true")
	w := httptest.NewRecorder()

	err := CreateEpa(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
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
		On("Func", mock.Anything, epaData{
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

	err := CreateEpa(client, template.Func)(w, r)
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
				On("Func", mock.Anything, epaData{
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

			err := CreateEpa(client, template.Func)(w, r)
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
		On("Func", mock.Anything, epaData{
			DonorId:  123,
			Title:    "Create an EPA",
			IsUpdate: false,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := CreateEpa(client, template.Func)(w, r)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}
