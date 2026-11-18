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

type mockUpdateEpaClient struct {
	mock.Mock
}

func (m *mockUpdateEpaClient) UpdateEpa(ctx sirius.Context, caseId int, epa sirius.Epa) error {
	args := m.Called(ctx, caseId, epa)
	return args.Error(0)
}

func (m *mockUpdateEpaClient) Epa(ctx sirius.Context, id int) (sirius.Epa, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Epa), args.Error(1)
}

func TestGetUpdateEpa(t *testing.T) {
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

			client := &mockUpdateEpaClient{}

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, epaData{
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

			err := UpdateEpa(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestGetUpdateEpaBadQuery(t *testing.T) {
	testCases := map[string]string{
		"no-id":       "/",
		"bad-id":      "/?id=test",
		"bad-case-id": "/?id=123&caseId=test",
	}

	for name, url := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			err := UpdateEpa(nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestGetUpdateEpaWhenEpaErrors(t *testing.T) {
	client := &mockUpdateEpaClient{}
	client.
		On("Epa", mock.Anything, 234).
		Return(sirius.Epa{}, errExample)

	template := &mockTemplate{}

	r, _ := http.NewRequest(http.MethodGet, "/?id=123&caseId=234", nil)
	w := httptest.NewRecorder()

	err := UpdateEpa(client, template.Func)(w, r)

	assert.Equal(t, err, errExample)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostUpdateEpa(t *testing.T) {
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
	newEpa := epa
	newEpa.DonorHasOtherEpas = truePtr
	newEpa.OtherEpaInfo = "More info 2"

	tests := []struct {
		name          string
		expectedError error
	}{
		{
			name:          "Update EPA post form success",
			expectedError: nil,
		},
		{
			name:          "Update EPA returns an API failure",
			expectedError: errExample,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockUpdateEpaClient{}
			client.
				On("Epa", mock.Anything, 234).
				Return(epa, nil)
			client.
				On("UpdateEpa", mock.Anything, 234, newEpa).
				Return(tc.expectedError)

			template := &mockTemplate{}
			if tc.expectedError == nil {
				template.
					On("Func", mock.Anything, epaData{
						DonorId:         123,
						Title:           "Edit EPA",
						Success:         true,
						Epa:             epa,
						AppointmentType: "singular",
						CaseId:          234,
						IsUpdate:        true,
					}).
					Return(nil)
			}

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

			r, _ := http.NewRequest(http.MethodPost, "/update-epa?id=123&caseId=234", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := UpdateEpa(client, template.Func)(w, r)
			resp := w.Result()

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostUpdateEpaRedirects(t *testing.T) {
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
		epa           sirius.Epa
	}{
		"add attorney": {
			query:         "/update-epa?id=123&caseId=456",
			formField:     "addAttorney",
			formValue:     "true",
			expectedError: RedirectError("/create-attorney?id=123&caseId=456&caseType=epa"),
			epa:           sirius.Epa{Case: sirius.Case{ID: 456}},
		},
		"update attorney": {
			query:         "/update-epa?id=123&caseId=456",
			formField:     "updateAttorney",
			formValue:     "789",
			expectedError: RedirectError("/update-attorney?id=123&caseId=456&caseType=epa&attorneyId=789"),
			epa:           sirius.Epa{Case: sirius.Case{ID: 456}},
		},
		"add correspondent": {
			query:         "/update-epa?id=123&caseId=456",
			formField:     "addCorrespondent",
			formValue:     "true",
			expectedError: RedirectError("/create-correspondent?id=123&caseId=456&caseType=epa"),
			epa:           sirius.Epa{Case: sirius.Case{ID: 456}},
		},
		"add correspondent with attorney on case": {
			query:         "/update-epa?id=123&caseId=456",
			formField:     "addCorrespondent",
			formValue:     "true",
			expectedError: RedirectError("/select-or-create-correspondent?id=123&caseId=456&caseType=epa"),
			epa: sirius.Epa{
				Case: sirius.Case{
					ID:        456,
					Attorneys: []sirius.Attorney{{Person: sirius.Person{ID: 1}}},
				},
			},
		},
		"update correspondent": {
			query:         "/update-epa?id=123&caseId=456",
			formField:     "updateCorrespondent",
			formValue:     "1",
			expectedError: RedirectError("/edit-correspondent?id=123&caseId=456&caseType=epa"),
			epa:           sirius.Epa{Case: sirius.Case{ID: 456}},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			client := &mockUpdateEpaClient{}
			template := &mockTemplate{}
			client.
				On("Epa", mock.Anything, 456).
				Return(tc.epa, nil)
			client.
				On("UpdateEpa", mock.Anything, 456, baseEpa).
				Return(nil)

			form := url.Values{
				tc.formField: {tc.formValue},
			}

			r, _ := http.NewRequest(http.MethodPost, tc.query, strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := UpdateEpa(client, template.Func)(w, r)
			resp := w.Result()

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostUpdateEpaReturnsValidationError(t *testing.T) {
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

			client := &mockUpdateEpaClient{}
			client.
				On("Epa", mock.Anything, 456).
				Return(sirius.Epa{}, nil)
			client.
				On("UpdateEpa", mock.Anything, 456, epa).
				Return(returnedError)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, epaData{
					DonorId:         123,
					CaseId:          456,
					Title:           "Edit EPA",
					IsUpdate:        true,
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

			r, _ := http.NewRequest(http.MethodPost, "/update-epa?id=123&caseId=456", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := UpdateEpa(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}
