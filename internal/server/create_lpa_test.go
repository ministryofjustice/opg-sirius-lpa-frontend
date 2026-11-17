package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCreateLpaClient struct {
	mock.Mock
}

func (m *mockCreateLpaClient) Person(ctx sirius.Context, id int) (sirius.Person, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Person), args.Error(1)
}

func (m *mockCreateLpaClient) Lpa(ctx sirius.Context, id int) (sirius.Lpa, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Lpa), args.Error(1)
}

func (m *mockCreateLpaClient) CreateLpa(ctx sirius.Context, donorID int, lpa sirius.Lpa) (sirius.Lpa, error) {
	args := m.Called(ctx, donorID, lpa)
	return args.Get(0).(sirius.Lpa), args.Error(1)
}

func TestGetCreateLpa(t *testing.T) {
	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{Firstname: "Firstname", Surname: "Surname"}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, lpaData{
			AllowNewNotifiedPerson: true,
			DonorId:                123,
			DonorName:              "Firstname Surname",
			FlowQuery:              "&flow=create",
			IsCreating:             true,
			Title:                  "Create an LPA",
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := CreateLpa(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetCreateLpaHtmxRequest(t *testing.T) {
	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{Firstname: "Firstname", Surname: "Surname"}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, lpaData{
			AllowNewNotifiedPerson: true,
			DonorId:                123,
			DonorName:              "Firstname Surname",
			FlowQuery:              "&flow=create",
			IsCreating:             true,
			IsPartial:              true,
			Title:                  "Create an LPA",
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	r.Header.Add("HX-Request", "true")
	w := httptest.NewRecorder()

	err := CreateLpa(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetCreateLpaDoesNotSetIsUpdate(t *testing.T) {
	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{Firstname: "Firstname", Surname: "Surname"}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, lpaData{
			AllowNewNotifiedPerson: true,
			DonorId:                123,
			DonorName:              "Firstname Surname",
			FlowQuery:              "&flow=create",
			IsCreating:             true,
			IsUpdate:               false,
			Title:                  "Create an LPA",
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := CreateLpa(client, template.Func)(w, r)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetCreateLpaBadQuery(t *testing.T) {
	testCases := map[string]string{
		"no-id":  "/",
		"bad-id": "/?id=test",
	}

	for name, url := range testCases {
		t.Run(name, func(t *testing.T) {
			client := &mockCreateLpaClient{}
			client.
				On("Person", mock.Anything, 123).
				Return(sirius.Person{}, nil)

			r, _ := http.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			err := CreateLpa(client, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestCreateLpaWhenPersonErrors(t *testing.T) {
	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, errExample)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := CreateLpa(client, nil)(w, r)

	assert.Equal(t, err, errExample)
	mock.AssertExpectationsForObjects(t, client)
}

func TestPostCreateLpa(t *testing.T) {
	dateString := "2022-04-05"
	lpa := sirius.Lpa{
		OnlineLpaId:                      "A12345678901",
		AttorneyActDecisions:             "When Registered",
		ApplicantType:                    "donor",
		ApplicantIds:                     []int{123},
		AnyOtherInfo:                     shared.BoolPtr(true),
		AdditionalInfo:                   "Some extra info",
		ApplicationHasContinuationSheets: shared.BoolPtr(true),
		ApplicationHasGuidance:           shared.BoolPtr(true),
		ApplicationHasRestrictions:       shared.BoolPtr(false),
		PaymentByDebitCreditCard:         shared.BoolPtr(true),
		PaymentRemission:                 shared.BoolPtr(false),
		RepeatApplication:                shared.BoolPtr(false),
		CardPaymentContact:               "01234 567890",
		CertificateProviderSignatureDate: sirius.DateString(dateString),
		Case: sirius.Case{
			SubType:                         "pfa",
			ApplicationType:                 "Online",
			ReceiptDate:                     sirius.DateString(dateString),
			LpaDonorSignatureDate:           sirius.DateString(dateString),
			CaseAttorneySingular:            shared.BoolPtr(true),
			CaseAttorneyJointly:             shared.BoolPtr(false),
			CaseAttorneyJointlyAndSeverally: shared.BoolPtr(false),
			CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(false),
			PaymentByCheque: shared.BoolPtr(false),
		},
	}

	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{Firstname: "Firstname", Surname: "Surname"}, nil)
	client.
		On("CreateLpa", mock.Anything, 123, lpa).
		Return(sirius.Lpa{Case: sirius.Case{ID: 456}}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, lpaData{
			AllowNewNotifiedPerson: true,
			AttorneyApplicants:     nil,
			CaseId:                 456,
			DonorId:                123,
			DonorName:              "Firstname Surname",
			IsCreating:             true,
			FlowQuery:              "&flow=create",
			Lpa:                    sirius.Lpa{Case: sirius.Case{ID: 456}},
			Success:                true,
			SuccessMessage:         "You have successfully created an LPA.",
			Title:                  "Create an LPA",
		}).
		Return(nil)

	form := url.Values{
		"additionalInfo":                   {"Some extra info"},
		"anyOtherInfo":                     {"true"},
		"applicantType":                    {"donor"},
		"applicationFee":                   {"card"},
		"applicationType":                  {"Online"},
		"attorneyActDecisions":             {"When Registered"},
		"cardPaymentContact":               {"01234 567890"},
		"caseAttorney":                     {"singular"},
		"caseSubtype":                      {"pfa"},
		"certificateProviderSignatureDate": {dateString},
		"lpaDonorSignatureDate":            {dateString},
		"onlineLpaId":                      {"A12345678901"},
		"preferencesAndInstructions":       {"guidance"},
		"continuationSheets":               {"continuationSheets"},
		"receiptDate":                      {dateString},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateLpa(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateLpaClearsMismatchedSubtypeOnlyFields(t *testing.T) {
	lpa := sirius.Lpa{
		ApplicationHasContinuationSheets: shared.BoolPtr(false),
		ApplicationHasGuidance:           shared.BoolPtr(false),
		ApplicationHasRestrictions:       shared.BoolPtr(false),
		PaymentByDebitCreditCard:         shared.BoolPtr(false),
		PaymentRemission:                 shared.BoolPtr(false),
		RepeatApplication:                shared.BoolPtr(false),
		AnyOtherInfo:                     shared.BoolPtr(false),
		Case: sirius.Case{
			SubType:                                   "hw",
			LifeSustainingTreatment:                   "Option A",
			CaseAttorneySingular:                      shared.BoolPtr(true),
			CaseAttorneyJointly:                       shared.BoolPtr(false),
			CaseAttorneyJointlyAndSeverally:           shared.BoolPtr(false),
			CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(false),
			PaymentByCheque:                           shared.BoolPtr(false),
		},
	}

	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("CreateLpa", mock.Anything, 123, lpa).
		Return(sirius.Lpa{}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.Anything).
		Return(nil)

	form := url.Values{
		"attorneyActDecisions":    {"When Registered"},
		"caseAttorney":            {"singular"},
		"caseSubtype":             {"hw"},
		"lifeSustainingTreatment": {"Option A"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateLpa(client, template.Func)(w, r)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateLpaPreferencesNoneClearsOtherSelections(t *testing.T) {
	lpa := sirius.Lpa{
		ApplicationHasContinuationSheets: shared.BoolPtr(false),
		ApplicationHasGuidance:           shared.BoolPtr(false),
		ApplicationHasRestrictions:       shared.BoolPtr(false),
		PaymentByDebitCreditCard:         shared.BoolPtr(false),
		PaymentRemission:                 shared.BoolPtr(false),
		RepeatApplication:                shared.BoolPtr(false),
		AnyOtherInfo:                     shared.BoolPtr(false),
		Case: sirius.Case{
			CaseAttorneySingular:                      shared.BoolPtr(true),
			CaseAttorneyJointly:                       shared.BoolPtr(false),
			CaseAttorneyJointlyAndSeverally:           shared.BoolPtr(false),
			CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(false),
			PaymentByCheque:                           shared.BoolPtr(false),
		},
	}

	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("CreateLpa", mock.Anything, 123, lpa).
		Return(sirius.Lpa{}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.Anything).
		Return(nil)

	form := url.Values{
		"caseAttorney":               {"singular"},
		"preferencesAndInstructions": {"guidance", "restrictions", "none"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateLpa(client, template.Func)(w, r)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateLpaDropsOnlineLpaIdWhenNotOnline(t *testing.T) {
	lpa := sirius.Lpa{
		ApplicationHasContinuationSheets: shared.BoolPtr(false),
		ApplicationHasGuidance:           shared.BoolPtr(false),
		ApplicationHasRestrictions:       shared.BoolPtr(false),
		PaymentByDebitCreditCard:         shared.BoolPtr(false),
		PaymentRemission:                 shared.BoolPtr(false),
		RepeatApplication:                shared.BoolPtr(false),
		AnyOtherInfo:                     shared.BoolPtr(false),
		Case: sirius.Case{
			SubType:                                   "pfa",
			ApplicationType:                           "Classic",
			CaseAttorneySingular:                      shared.BoolPtr(true),
			CaseAttorneyJointly:                       shared.BoolPtr(false),
			CaseAttorneyJointlyAndSeverally:           shared.BoolPtr(false),
			CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(false),
			PaymentByCheque:                           shared.BoolPtr(false),
		},
	}

	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("CreateLpa", mock.Anything, 123, lpa).
		Return(sirius.Lpa{}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.Anything).
		Return(nil)

	form := url.Values{
		"applicationType": {"Classic"},
		"caseAttorney":    {"singular"},
		"caseSubtype":     {"pfa"},
		"onlineLpaId":     {"A12345678901"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateLpa(client, template.Func)(w, r)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateLpaDropsCardPaymentContactWhenCardNotSelected(t *testing.T) {
	lpa := sirius.Lpa{
		ApplicationHasContinuationSheets: shared.BoolPtr(false),
		ApplicationHasGuidance:           shared.BoolPtr(false),
		ApplicationHasRestrictions:       shared.BoolPtr(false),
		PaymentByDebitCreditCard:         shared.BoolPtr(false),
		PaymentRemission:                 shared.BoolPtr(false),
		RepeatApplication:                shared.BoolPtr(false),
		AnyOtherInfo:                     shared.BoolPtr(false),
		Case: sirius.Case{
			CaseAttorneySingular:                      shared.BoolPtr(true),
			CaseAttorneyJointly:                       shared.BoolPtr(false),
			CaseAttorneyJointlyAndSeverally:           shared.BoolPtr(false),
			CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(false),
			PaymentByCheque:                           shared.BoolPtr(true),
		},
	}

	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("CreateLpa", mock.Anything, 123, lpa).
		Return(sirius.Lpa{}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.Anything).
		Return(nil)

	form := url.Values{
		"applicationFee":     {"cheque"}, // not "card"
		"caseAttorney":       {"singular"},
		"cardPaymentContact": {"01234 567890"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateLpa(client, template.Func)(w, r)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateLpaDropsAdditionalInfoWhenAnyOtherInfoNotSelected(t *testing.T) {
	lpa := sirius.Lpa{
		ApplicationHasContinuationSheets: shared.BoolPtr(false),
		ApplicationHasGuidance:           shared.BoolPtr(false),
		ApplicationHasRestrictions:       shared.BoolPtr(false),
		PaymentByDebitCreditCard:         shared.BoolPtr(false),
		PaymentRemission:                 shared.BoolPtr(false),
		RepeatApplication:                shared.BoolPtr(false),
		AnyOtherInfo:                     shared.BoolPtr(false),
		Case: sirius.Case{
			CaseAttorneySingular:                      shared.BoolPtr(true),
			CaseAttorneyJointly:                       shared.BoolPtr(false),
			CaseAttorneyJointlyAndSeverally:           shared.BoolPtr(false),
			CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(false),
			PaymentByCheque:                           shared.BoolPtr(false),
		},
	}

	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("CreateLpa", mock.Anything, 123, lpa).
		Return(sirius.Lpa{}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.Anything).
		Return(nil)

	form := url.Values{
		"additionalInfo": {"Some info the user typed then unchecked yes"},
		"caseAttorney":   {"singular"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateLpa(client, template.Func)(w, r)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateLpaWhenValidationError(t *testing.T) {
	expectedError := sirius.ValidationError{
		Field: sirius.FieldErrors{"receiptDate": {"receiptDate": "Select the date of receipt"}},
	}

	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("CreateLpa", mock.Anything, 123, mock.Anything).
		Return(sirius.Lpa{}, expectedError)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.MatchedBy(func(d lpaData) bool {
			return d.Error.Error() == expectedError.Error() && !d.Success
		})).
		Return(nil)

	form := url.Values{"caseAttorney": {"singular"}}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateLpa(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateLpaWhenValidationErrorHtmxRequest(t *testing.T) {
	expectedError := sirius.ValidationError{
		Field: sirius.FieldErrors{"receiptDate": {"receiptDate": "Select the date of receipt"}},
	}

	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("CreateLpa", mock.Anything, 123, mock.Anything).
		Return(sirius.Lpa{}, expectedError)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.MatchedBy(func(d lpaData) bool {
			return d.Error.Error() == expectedError.Error() && !d.Success
		})).
		Return(nil)

	form := url.Values{"caseAttorney": {"singular"}}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	r.Header.Add("HX-Request", "true")
	w := httptest.NewRecorder()

	err := CreateLpa(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateLpaWhenGenericError(t *testing.T) {
	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("CreateLpa", mock.Anything, 123, mock.Anything).
		Return(sirius.Lpa{}, errExample)

	form := url.Values{"caseAttorney": {"singular"}}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateLpa(client, nil)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestPostCreateLpaAddReplacementAttorney(t *testing.T) {
	for _, isHtmx := range []bool{false, true} {
		t.Run("Is Htmx: "+strconv.FormatBool(isHtmx), func(t *testing.T) {
			dateString := "2022-04-05"
			lpa := sirius.Lpa{
				OnlineLpaId:                      "A12345678901",
				AttorneyActDecisions:             "When Registered",
				ApplicantType:                    "donor",
				ApplicantIds:                     []int{123},
				AnyOtherInfo:                     shared.BoolPtr(true),
				AdditionalInfo:                   "Some extra info",
				ApplicationHasContinuationSheets: shared.BoolPtr(false),
				ApplicationHasGuidance:           shared.BoolPtr(true),
				ApplicationHasRestrictions:       shared.BoolPtr(false),
				PaymentByDebitCreditCard:         shared.BoolPtr(true),
				PaymentRemission:                 shared.BoolPtr(false),
				RepeatApplication:                shared.BoolPtr(false),
				CardPaymentContact:               "01234 567890",
				Case: sirius.Case{
					SubType:                         "pfa",
					ApplicationType:                 "Online",
					ReceiptDate:                     sirius.DateString(dateString),
					LpaDonorSignatureDate:           sirius.DateString(dateString),
					CaseAttorneySingular:            shared.BoolPtr(true),
					CaseAttorneyJointly:             shared.BoolPtr(false),
					CaseAttorneyJointlyAndSeverally: shared.BoolPtr(false),
					CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(false),
					PaymentByCheque: shared.BoolPtr(false),
				},
			}

			client := &mockCreateLpaClient{}
			client.
				On("Person", mock.Anything, 123).
				Return(sirius.Person{Firstname: "Firstname", Surname: "Surname"}, nil)
			client.
				On("CreateLpa", mock.Anything, 123, lpa).
				Return(sirius.Lpa{Case: sirius.Case{ID: 456}}, nil)

			template := &mockTemplate{}

			form := url.Values{
				"additionalInfo":             {"Some extra info"},
				"addReplacementAttorney":     {"Add replacement attorney"},
				"anyOtherInfo":               {"true"},
				"applicantType":              {"donor"},
				"applicationFee":             {"card"},
				"applicationType":            {"Online"},
				"attorneyActDecisions":       {"When Registered"},
				"cardPaymentContact":         {"01234 567890"},
				"caseAttorney":               {"singular"},
				"caseSubtype":                {"pfa"},
				"lpaDonorSignatureDate":      {dateString},
				"onlineLpaId":                {"A12345678901"},
				"preferencesAndInstructions": {"guidance"},
				"receiptDate":                {dateString},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			if isHtmx {
				r.Header.Add("HX-Request", "true")
			}
			w := httptest.NewRecorder()

			err := CreateLpa(client, template.Func)(w, r)
			resp := w.Result()

			expectedRedirect := RedirectError("/create-replacement-attorney?id=123&caseId=456&flow=create")
			assert.Equal(t, expectedRedirect, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostCreateLpaRedirects(t *testing.T) {
	tests := []struct {
		name        string
		formKey     string
		formValue   string
		expectedErr error
	}{
		{
			name:        "Add attorney redirects",
			formKey:     "addAttorney",
			formValue:   "true",
			expectedErr: RedirectError("/create-attorney?id=1&caseId=2&caseType=lpa"),
		},
		{
			name:        "Add certificate provider redirects",
			formKey:     "addCertificateProvider",
			formValue:   "true",
			expectedErr: RedirectError("/create-certificate-provider?id=1&caseId=2"),
		},
		{
			name:        "Add correspondent",
			formKey:     "addCorrespondent",
			formValue:   "true",
			expectedErr: RedirectError("/select-or-create-correspondent?id=1&caseId=2&caseType=lpa"),
		},
		{
			name:        "Add notified person redirects",
			formKey:     "addNotifiedPerson",
			formValue:   "true",
			expectedErr: RedirectError("/create-notified-person?id=1&caseId=2"),
		},
	}

	flowCases := []struct {
		name      string
		query     string
		flowQuery string
	}{
		{name: "New LPA", flowQuery: "&flow=create"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for _, flowCase := range flowCases {
				t.Run(flowCase.name, func(t *testing.T) {
					client := &mockCreateLpaClient{}
					client.
						On("Person", mock.Anything, 1).
						Return(sirius.Person{Firstname: "John", Surname: "Doe"}, nil)

					client.
						On("CreateLpa", mock.Anything, 1, mock.Anything).
						Return(sirius.Lpa{Case: sirius.Case{ID: 2}}, nil)

					template := &mockTemplate{}

					form := url.Values{
						"caseSubtype": {"pfa"},
						tc.formKey:    {tc.formValue},
					}

					r, _ := http.NewRequest(http.MethodPost, "/create-lpa?id=1"+flowCase.query, strings.NewReader(form.Encode()))
					r.Header.Add("Content-Type", formUrlEncoded)
					w := httptest.NewRecorder()

					err := CreateLpa(client, template.Func)(w, r)
					resp := w.Result()

					expectedErr := tc.expectedErr
					if redirectErr, ok := expectedErr.(RedirectError); ok {
						expectedErr = RedirectError(string(redirectErr) + flowCase.flowQuery)
					}
					assert.Equal(t, expectedErr, err)
					assert.Equal(t, http.StatusOK, resp.StatusCode)
					mock.AssertExpectationsForObjects(t, client, template)
				})
			}
		})
	}
}
