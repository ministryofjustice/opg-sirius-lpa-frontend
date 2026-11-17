package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUpdateLpaClient struct {
	mock.Mock
}

func (m *mockUpdateLpaClient) Person(ctx sirius.Context, id int) (sirius.Person, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Person), args.Error(1)
}

func (m *mockUpdateLpaClient) Lpa(ctx sirius.Context, id int) (sirius.Lpa, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Lpa), args.Error(1)
}

func (m *mockUpdateLpaClient) UpdateLpa(ctx sirius.Context, caseID int, lpa sirius.Lpa) error {
	args := m.Called(ctx, caseID, lpa)
	return args.Error(0)
}

func (m *mockUpdateLpaClient) UpdateAttorney(ctx sirius.Context, attorneyId int, attorney sirius.Attorney) error {
	args := m.Called(ctx, attorneyId, attorney)
	return args.Error(0)
}

func (m *mockUpdateLpaClient) UpdateReplacementAttorney(ctx sirius.Context, attorneyId int, attorney sirius.Attorney) error {
	args := m.Called(ctx, attorneyId, attorney)
	return args.Error(0)
}

func (m *mockUpdateLpaClient) UpdateTrustCorporation(ctx sirius.Context, trustCorporationId int, trustCorporation sirius.TrustCorporation) error {
	args := m.Called(ctx, trustCorporationId, trustCorporation)
	return args.Error(0)
}

func TestGetUpdateLpa(t *testing.T) {
	for _, tc := range []struct {
		formValue string
		lpa       sirius.Lpa
	}{
		{"singular", sirius.Lpa{Case: sirius.Case{CaseAttorneySingular: shared.BoolPtr(true)}}},
		{"jointly", sirius.Lpa{Case: sirius.Case{CaseAttorneyJointly: shared.BoolPtr(true)}}},
		{"jointly-and-severally", sirius.Lpa{Case: sirius.Case{CaseAttorneyJointlyAndSeverally: shared.BoolPtr(true)}}},
		{"jointly-and-jointly-and-severally", sirius.Lpa{Case: sirius.Case{CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(true)}}},
	} {
		t.Run(tc.formValue, func(t *testing.T) {
			client := &mockUpdateLpaClient{}
			client.
				On("Person", mock.Anything, 123).
				Return(sirius.Person{Firstname: "Firstname", Surname: "Surname"}, nil)
			client.
				On("Lpa", mock.Anything, 456).
				Return(tc.lpa, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, lpaData{
					AllowNewNotifiedPerson: true,
					AppointmentType:        tc.formValue,
					AttorneyApplicants:     []sirius.Attorney{},
					CaseId:                 456,
					DonorId:                123,
					DonorName:              "Firstname Surname",
					IsUpdate:               true,
					Lpa:                    tc.lpa,
					ReplacementAttorneys:   []sirius.Attorney{},
					Title:                  "Edit LPA",
				}).
				Return(nil)

			r, _ := http.NewRequest(http.MethodGet, "/?id=123&caseId=456", nil)
			w := httptest.NewRecorder()

			err := UpdateLpa(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestGetUpdateLpaWithTrustCorporations(t *testing.T) {
	tests := []struct {
		name                  string
		isReplacementAttorney bool
		AttorneyApplicants    []sirius.Attorney
	}{
		{
			name:                  "Trust corp is a replacement attorney",
			isReplacementAttorney: true,
			AttorneyApplicants:    []sirius.Attorney{},
		},
		{
			name:                  "Trust corp is an attorney",
			isReplacementAttorney: false,
			AttorneyApplicants:    []sirius.Attorney{{Person: sirius.Person{ID: 1}}},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			trustCorporation := sirius.TrustCorporation{
				IsReplacementAttorney: tc.isReplacementAttorney,
				Attorney:              sirius.Attorney{Person: sirius.Person{ID: 1}},
			}

			lpa := sirius.Lpa{Case: sirius.Case{TrustCorporations: []sirius.TrustCorporation{trustCorporation}}}
			client := &mockUpdateLpaClient{}
			client.
				On("Person", mock.Anything, 123).
				Return(sirius.Person{Firstname: "Firstname", Surname: "Surname"}, nil)
			client.
				On("Lpa", mock.Anything, 456).
				Return(lpa, nil)

			data := lpaData{
				AllowNewNotifiedPerson: true,
				AttorneyApplicants:     tc.AttorneyApplicants,
				CaseId:                 456,
				DonorId:                123,
				DonorName:              "Firstname Surname",
				IsUpdate:               true,
				Lpa:                    lpa,
				Title:                  "Edit LPA",
			}

			if tc.isReplacementAttorney {
				data.ReplacementAttorneyTrustCorporations = []sirius.TrustCorporation{trustCorporation}
				data.ReplacementAttorneys = []sirius.Attorney{trustCorporation.Attorney}
			} else {
				data.AttorneyTrustCorporations = []sirius.TrustCorporation{trustCorporation}
				data.ReplacementAttorneys = []sirius.Attorney{}
			}

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, data).
				Return(nil)

			r, _ := http.NewRequest(http.MethodGet, "/?id=123&caseId=456", nil)
			w := httptest.NewRecorder()

			err := UpdateLpa(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestUpdateLpaWhenLpaErrors(t *testing.T) {
	client := &mockUpdateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("Lpa", mock.Anything, 456).
		Return(sirius.Lpa{}, errExample)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123&caseId=456", nil)
	w := httptest.NewRecorder()

	err := UpdateLpa(client, nil)(w, r)

	assert.Equal(t, err, errExample)
	mock.AssertExpectationsForObjects(t, client)
}

func TestPostUpdateLpaApplicantAndLifeSustainingTreatmentFields(t *testing.T) {
	lpa := sirius.Lpa{
		ApplicationHasContinuationSheets: shared.BoolPtr(false),
		ApplicationHasGuidance:           shared.BoolPtr(false),
		ApplicationHasRestrictions:       shared.BoolPtr(false),
		PaymentByDebitCreditCard:         shared.BoolPtr(false),
		PaymentRemission:                 shared.BoolPtr(false),
		RepeatApplication:                shared.BoolPtr(false),
		AnyOtherInfo:                     shared.BoolPtr(false),
		ApplicantType:                    "attorney",
		ApplicantSignatureDate:           "2022-04-06",
		ApplicantIds:                     []int{1, 2},
		Case: sirius.Case{
			SubType:                                   "hw",
			LifeSustainingTreatment:                   "None",
			LifeSustainingTreatmentSignatureDateA:     "2022-04-05",
			CaseAttorneySingular:                      shared.BoolPtr(false),
			CaseAttorneyJointly:                       shared.BoolPtr(true),
			CaseAttorneyJointlyAndSeverally:           shared.BoolPtr(false),
			CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(false),
			PaymentByCheque:                           shared.BoolPtr(false),
		},
	}

	client := &mockUpdateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("Lpa", mock.Anything, 456).
		Return(sirius.Lpa{}, nil).
		Twice()
	client.
		On("UpdateLpa", mock.Anything, 456, lpa).
		Return(nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.Anything).
		Return(nil)

	form := url.Values{
		"caseSubtype":                          {"hw"},
		"caseAttorney":                         {"jointly"},
		"lifeSustainingTreatment":              {"None"},
		"lifeSustainingTreatmentSignatureDate": {"2022-04-05"},
		"applicantType":                        {"attorney"},
		"applicantSignatureDate":               {"2022-04-06"},
		"applicantIds":                         {"1", "2"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123&caseId=456", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := UpdateLpa(client, template.Func)(w, r)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostUpdateLpaAttorneySignatureDates(t *testing.T) {
	existingLpa := sirius.Lpa{
		Case: sirius.Case{
			ID:          456,
			ReceiptDate: sirius.DateString("2022-01-01"),
			Attorneys: []sirius.Attorney{
				{Person: sirius.Person{ID: 876}, LpaPartCSignatureDate: sirius.DateString("2022-01-01")},
				{Person: sirius.Person{ID: 987}, LpaPartCSignatureDate: sirius.DateString("2022-01-02")},
				{Person: sirius.Person{ID: 786}},
			},
		},
	}
	submittedLpa := sirius.Lpa{
		ApplicationHasContinuationSheets: shared.BoolPtr(false),
		ApplicationHasGuidance:           shared.BoolPtr(false),
		ApplicationHasRestrictions:       shared.BoolPtr(false),
		PaymentByDebitCreditCard:         shared.BoolPtr(false),
		PaymentRemission:                 shared.BoolPtr(false),
		RepeatApplication:                shared.BoolPtr(false),
		AnyOtherInfo:                     shared.BoolPtr(false),
		Case: sirius.Case{
			SubType:                                   "hw",
			CaseAttorneySingular:                      shared.BoolPtr(false),
			CaseAttorneyJointly:                       shared.BoolPtr(true),
			CaseAttorneyJointlyAndSeverally:           shared.BoolPtr(false),
			CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(false),
			PaymentByCheque:                           shared.BoolPtr(false),
		},
	}

	client := &mockUpdateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("Lpa", mock.Anything, 456).
		Return(existingLpa, nil)
	client.
		On("UpdateLpa", mock.Anything, 456, submittedLpa).
		Return(nil)
	client.
		On("UpdateAttorney", mock.Anything, 876, sirius.Attorney{Person: sirius.Person{ID: 876}, LpaPartCSignatureDate: sirius.DateString("2022-01-02")}).
		Return(nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.Anything).
		Return(nil)

	form := url.Values{
		"caseAttorney":              {"jointly"},
		"caseSubtype":               {"hw"},
		"lpaPartCSignatureDate-876": {"2022-01-02"},
		"lpaPartCSignatureDate-987": {"2022-01-02"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123&caseId=456", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := UpdateLpa(client, template.Func)(w, r)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostUpdateLpaReplacementAttorneySignatureDates(t *testing.T) {
	existingLpa := sirius.Lpa{
		Case: sirius.Case{
			ID:          456,
			ReceiptDate: sirius.DateString("2022-01-01"),
			ReplacementAttorneys: []sirius.Attorney{
				{Person: sirius.Person{ID: 876}, LpaPartCSignatureDate: sirius.DateString("2022-01-01")},
				{Person: sirius.Person{ID: 987}, LpaPartCSignatureDate: sirius.DateString("2022-01-02")},
				{Person: sirius.Person{ID: 786}},
			},
		},
	}
	submittedLpa := sirius.Lpa{
		ApplicationHasContinuationSheets: shared.BoolPtr(false),
		ApplicationHasGuidance:           shared.BoolPtr(false),
		ApplicationHasRestrictions:       shared.BoolPtr(false),
		PaymentByDebitCreditCard:         shared.BoolPtr(false),
		PaymentRemission:                 shared.BoolPtr(false),
		RepeatApplication:                shared.BoolPtr(false),
		AnyOtherInfo:                     shared.BoolPtr(false),
		Case: sirius.Case{
			SubType:                                   "hw",
			CaseAttorneySingular:                      shared.BoolPtr(false),
			CaseAttorneyJointly:                       shared.BoolPtr(true),
			CaseAttorneyJointlyAndSeverally:           shared.BoolPtr(false),
			CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(false),
			PaymentByCheque:                           shared.BoolPtr(false),
		},
	}

	client := &mockUpdateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("Lpa", mock.Anything, 456).
		Return(existingLpa, nil)
	client.
		On("UpdateLpa", mock.Anything, 456, submittedLpa).
		Return(nil)
	client.
		On("UpdateReplacementAttorney", mock.Anything, 876, sirius.Attorney{Person: sirius.Person{ID: 876}, LpaPartCSignatureDate: sirius.DateString("2022-01-02")}).
		Return(nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.Anything).
		Return(nil)

	form := url.Values{
		"caseAttorney":              {"jointly"},
		"caseSubtype":               {"hw"},
		"lpaPartCSignatureDate-876": {"2022-01-02"},
		"lpaPartCSignatureDate-987": {"2022-01-02"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123&caseId=456", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := UpdateLpa(client, template.Func)(w, r)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostUpdateLpaTrustCorporationSignatureDates(t *testing.T) {
	existingLpa := sirius.Lpa{
		Case: sirius.Case{
			ID:          456,
			ReceiptDate: sirius.DateString("2022-01-01"),
			TrustCorporations: []sirius.TrustCorporation{
				{Attorney: sirius.Attorney{Person: sirius.Person{ID: 876}, LpaPartCSignatureDate: sirius.DateString("2022-01-01")}},
				{Attorney: sirius.Attorney{Person: sirius.Person{ID: 987}, LpaPartCSignatureDate: sirius.DateString("2022-01-02")}},
				{Attorney: sirius.Attorney{Person: sirius.Person{ID: 786}}, TrustCorporationAppointedAs: "Replacement Attorney"},
			},
		},
	}
	submittedLpa := sirius.Lpa{
		ApplicationHasContinuationSheets: shared.BoolPtr(false),
		ApplicationHasGuidance:           shared.BoolPtr(false),
		ApplicationHasRestrictions:       shared.BoolPtr(false),
		PaymentByDebitCreditCard:         shared.BoolPtr(false),
		PaymentRemission:                 shared.BoolPtr(false),
		RepeatApplication:                shared.BoolPtr(false),
		AnyOtherInfo:                     shared.BoolPtr(false),
		Case: sirius.Case{
			SubType:                                   "hw",
			CaseAttorneySingular:                      shared.BoolPtr(false),
			CaseAttorneyJointly:                       shared.BoolPtr(true),
			CaseAttorneyJointlyAndSeverally:           shared.BoolPtr(false),
			CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(false),
			PaymentByCheque:                           shared.BoolPtr(false),
		},
	}

	client := &mockUpdateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("Lpa", mock.Anything, 456).
		Return(existingLpa, nil)
	client.
		On("UpdateLpa", mock.Anything, 456, submittedLpa).
		Return(nil)
	client.
		On("UpdateTrustCorporation", mock.Anything, 876, sirius.TrustCorporation{Attorney: sirius.Attorney{Person: sirius.Person{ID: 876}, LpaPartCSignatureDate: sirius.DateString("2022-01-02")}}).
		Return(nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.Anything).
		Return(nil)

	form := url.Values{
		"caseAttorney":              {"jointly"},
		"caseSubtype":               {"hw"},
		"lpaPartCSignatureDate-876": {"2022-01-02"},
		"lpaPartCSignatureDate-987": {"2022-01-02"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123&caseId=456", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := UpdateLpa(client, template.Func)(w, r)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostUpdateLpaAttorneySignatureDatesError(t *testing.T) {
	existingLpa := sirius.Lpa{
		Case: sirius.Case{
			ID:          456,
			ReceiptDate: sirius.DateString("2022-01-01"),
			Attorneys: []sirius.Attorney{
				{Person: sirius.Person{ID: 876}, LpaPartCSignatureDate: sirius.DateString("2022-01-01")},
				{Person: sirius.Person{ID: 987}, LpaPartCSignatureDate: sirius.DateString("2022-01-02")},
				{Person: sirius.Person{ID: 786}},
			},
		},
	}
	submittedLpa := sirius.Lpa{
		ApplicationHasContinuationSheets: shared.BoolPtr(false),
		ApplicationHasGuidance:           shared.BoolPtr(false),
		ApplicationHasRestrictions:       shared.BoolPtr(false),
		PaymentByDebitCreditCard:         shared.BoolPtr(false),
		PaymentRemission:                 shared.BoolPtr(false),
		RepeatApplication:                shared.BoolPtr(false),
		AnyOtherInfo:                     shared.BoolPtr(false),
		Case: sirius.Case{
			SubType:                                   "hw",
			CaseAttorneySingular:                      shared.BoolPtr(false),
			CaseAttorneyJointly:                       shared.BoolPtr(true),
			CaseAttorneyJointlyAndSeverally:           shared.BoolPtr(false),
			CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(false),
			PaymentByCheque:                           shared.BoolPtr(false),
		},
	}

	client := &mockUpdateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("Lpa", mock.Anything, 456).
		Return(existingLpa, nil)
	client.
		On("UpdateLpa", mock.Anything, 456, submittedLpa).
		Return(nil)
	client.
		On("UpdateAttorney", mock.Anything, 876, sirius.Attorney{Person: sirius.Person{ID: 876}, LpaPartCSignatureDate: sirius.DateString("2022-01-02")}).
		Return(errExample)

	form := url.Values{
		"caseAttorney":              {"jointly"},
		"caseSubtype":               {"hw"},
		"lpaPartCSignatureDate-876": {"2022-01-02"},
		"lpaPartCSignatureDate-987": {"2022-01-02"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123&caseId=456", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := UpdateLpa(client, nil)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestPostUpdateLpaReplacementAttorneySignatureDatesError(t *testing.T) {
	existingLpa := sirius.Lpa{
		Case: sirius.Case{
			ID:          456,
			ReceiptDate: sirius.DateString("2022-01-01"),
			ReplacementAttorneys: []sirius.Attorney{
				{Person: sirius.Person{ID: 876}, LpaPartCSignatureDate: sirius.DateString("2022-01-01")},
				{Person: sirius.Person{ID: 987}, LpaPartCSignatureDate: sirius.DateString("2022-01-02")},
				{Person: sirius.Person{ID: 786}},
			},
		},
	}
	submittedLpa := sirius.Lpa{
		ApplicationHasContinuationSheets: shared.BoolPtr(false),
		ApplicationHasGuidance:           shared.BoolPtr(false),
		ApplicationHasRestrictions:       shared.BoolPtr(false),
		PaymentByDebitCreditCard:         shared.BoolPtr(false),
		PaymentRemission:                 shared.BoolPtr(false),
		RepeatApplication:                shared.BoolPtr(false),
		AnyOtherInfo:                     shared.BoolPtr(false),
		Case: sirius.Case{
			SubType:                                   "hw",
			CaseAttorneySingular:                      shared.BoolPtr(false),
			CaseAttorneyJointly:                       shared.BoolPtr(true),
			CaseAttorneyJointlyAndSeverally:           shared.BoolPtr(false),
			CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(false),
			PaymentByCheque:                           shared.BoolPtr(false),
		},
	}

	client := &mockUpdateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("Lpa", mock.Anything, 456).
		Return(existingLpa, nil)
	client.
		On("UpdateLpa", mock.Anything, 456, submittedLpa).
		Return(nil)
	client.
		On("UpdateReplacementAttorney", mock.Anything, 876, sirius.Attorney{Person: sirius.Person{ID: 876}, LpaPartCSignatureDate: sirius.DateString("2022-01-02")}).
		Return(errExample)

	form := url.Values{
		"caseAttorney":              {"jointly"},
		"caseSubtype":               {"hw"},
		"lpaPartCSignatureDate-876": {"2022-01-02"},
		"lpaPartCSignatureDate-987": {"2022-01-02"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123&caseId=456", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := UpdateLpa(client, nil)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestPostUpdateLpaTrustCorporationSignatureDatesError(t *testing.T) {
	existingLpa := sirius.Lpa{
		Case: sirius.Case{
			ID:          456,
			ReceiptDate: sirius.DateString("2022-01-01"),
			TrustCorporations: []sirius.TrustCorporation{
				{Attorney: sirius.Attorney{Person: sirius.Person{ID: 876}, LpaPartCSignatureDate: sirius.DateString("2022-01-01")}},
				{Attorney: sirius.Attorney{Person: sirius.Person{ID: 987}, LpaPartCSignatureDate: sirius.DateString("2022-01-02")}},
				{Attorney: sirius.Attorney{Person: sirius.Person{ID: 786}}},
			},
		},
	}
	submittedLpa := sirius.Lpa{
		ApplicationHasContinuationSheets: shared.BoolPtr(false),
		ApplicationHasGuidance:           shared.BoolPtr(false),
		ApplicationHasRestrictions:       shared.BoolPtr(false),
		PaymentByDebitCreditCard:         shared.BoolPtr(false),
		PaymentRemission:                 shared.BoolPtr(false),
		RepeatApplication:                shared.BoolPtr(false),
		AnyOtherInfo:                     shared.BoolPtr(false),
		Case: sirius.Case{
			SubType:                                   "hw",
			CaseAttorneySingular:                      shared.BoolPtr(false),
			CaseAttorneyJointly:                       shared.BoolPtr(true),
			CaseAttorneyJointlyAndSeverally:           shared.BoolPtr(false),
			CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(false),
			PaymentByCheque:                           shared.BoolPtr(false),
		},
	}

	client := &mockUpdateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("Lpa", mock.Anything, 456).
		Return(existingLpa, nil)
	client.
		On("UpdateLpa", mock.Anything, 456, submittedLpa).
		Return(nil)
	client.
		On("UpdateTrustCorporation", mock.Anything, 876, sirius.TrustCorporation{Attorney: sirius.Attorney{Person: sirius.Person{ID: 876}, LpaPartCSignatureDate: sirius.DateString("2022-01-02")}}).
		Return(errExample)

	form := url.Values{
		"caseAttorney":              {"jointly"},
		"caseSubtype":               {"hw"},
		"lpaPartCSignatureDate-876": {"2022-01-02"},
		"lpaPartCSignatureDate-987": {"2022-01-02"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123&caseId=456", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := UpdateLpa(client, nil)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestPostUpdateLpaWhenValidationError(t *testing.T) {
	expectedError := sirius.ValidationError{
		Field: sirius.FieldErrors{"receiptDate": {"receiptDate": "Select the date of receipt"}},
	}

	existingLpa := sirius.Lpa{Case: sirius.Case{ID: 456, ReceiptDate: sirius.DateString("2022-01-01")}}

	submittedLpa := sirius.Lpa{
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

	client := &mockUpdateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("Lpa", mock.Anything, 456).
		Return(existingLpa, nil).
		Once()
	client.
		On("UpdateLpa", mock.Anything, 456, submittedLpa).
		Return(expectedError)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.MatchedBy(func(d lpaData) bool {
			return d.Error.Error() == expectedError.Error() &&
				!d.Success &&
				d.IsUpdate &&
				d.CaseId == 456 &&
				reflect.DeepEqual(d.Lpa, submittedLpa)
		})).
		Return(nil)

	form := url.Values{"caseAttorney": {"singular"}}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123&caseId=456", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := UpdateLpa(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostUpdateLpaUpdateAttorney(t *testing.T) {
	for _, isHtmx := range []bool{false, true} {
		t.Run("Is Htmx: "+strconv.FormatBool(isHtmx), func(t *testing.T) {
			attorney := sirius.Attorney{Person: sirius.Person{ID: 999, Firstname: "Rudolph", Surname: "Stotesbury"}}
			existingLpa := sirius.Lpa{
				Case: sirius.Case{
					ID:                   456,
					ReceiptDate:          sirius.DateString("2022-01-01"),
					Attorneys:            []sirius.Attorney{attorney},
					ReplacementAttorneys: []sirius.Attorney{},
				},
			}

			submittedLpa := sirius.Lpa{
				ApplicationHasContinuationSheets: shared.BoolPtr(false),
				ApplicationHasGuidance:           shared.BoolPtr(false),
				ApplicationHasRestrictions:       shared.BoolPtr(false),
				PaymentByDebitCreditCard:         shared.BoolPtr(false),
				PaymentRemission:                 shared.BoolPtr(false),
				RepeatApplication:                shared.BoolPtr(false),
				AnyOtherInfo:                     shared.BoolPtr(false),
				Case: sirius.Case{
					SubType:                                   "hw",
					CaseAttorneySingular:                      shared.BoolPtr(false),
					CaseAttorneyJointly:                       shared.BoolPtr(true),
					CaseAttorneyJointlyAndSeverally:           shared.BoolPtr(false),
					CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(false),
					PaymentByCheque:                           shared.BoolPtr(false),
				},
			}

			client := &mockUpdateLpaClient{}
			client.
				On("Person", mock.Anything, 123).
				Return(sirius.Person{Firstname: "Firstname", Surname: "Surname"}, nil)
			client.
				On("Lpa", mock.Anything, 456).
				Return(existingLpa, nil)
			client.
				On("UpdateLpa", mock.Anything, 456, submittedLpa).
				Return(nil)

			template := &mockTemplate{}

			form := url.Values{
				"caseAttorney":   {"jointly"},
				"caseSubtype":    {"hw"},
				"updateAttorney": {"999"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=123&caseId=456", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			if isHtmx {
				r.Header.Add("HX-Request", "true")
			}
			w := httptest.NewRecorder()

			err := UpdateLpa(client, template.Func)(w, r)
			resp := w.Result()

			expectedRedirect := RedirectError("/update-attorney?id=123&caseId=456&caseType=lpa&attorneyId=999")
			assert.Equal(t, expectedRedirect, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostUpdateLpaUpdateAttorneyBadId(t *testing.T) {
	client := &mockUpdateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("Lpa", mock.Anything, 456).
		Return(sirius.Lpa{}, nil).
		Twice()
	client.
		On("UpdateLpa", mock.Anything, 456, mock.Anything).
		Return(nil)

	form := url.Values{"updateAttorney": {"not-a-number"}}
	r, _ := http.NewRequest(http.MethodPost, "/?id=123&caseId=456", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := UpdateLpa(client, nil)(w, r)

	assert.NotNil(t, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestPostUpdateLpaUpdateTrustCorporationBadId(t *testing.T) {
	for _, isReplacementAttorney := range []bool{false, true} {
		t.Run("Is Replacement Attorney: "+strconv.FormatBool(isReplacementAttorney), func(t *testing.T) {
			client := &mockUpdateLpaClient{}
			client.
				On("Person", mock.Anything, 1).
				Return(sirius.Person{Firstname: "John", Surname: "Doe"}, nil)
			client.
				On("Lpa", mock.Anything, 2).
				Return(sirius.Lpa{Case: sirius.Case{ID: 2}}, nil).
				Twice()
			client.
				On("UpdateLpa", mock.Anything, 2, mock.Anything).
				Return(nil)

			form := url.Values{
				"caseSubtype": {"pfa"},
			}

			if isReplacementAttorney {
				form["updateTrustCorporationReplacementAttorney"] = []string{"bad-id"}
			} else {
				form["updateTrustCorporationAttorney"] = []string{"bad-id"}
			}

			r, _ := http.NewRequest(http.MethodPost, "/update-lpa?id=1&caseId=2", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := UpdateLpa(client, nil)(w, r)

			assert.NotNil(t, err)
			mock.AssertExpectationsForObjects(t, client)
		})
	}
}

func TestPostUpdateLpaUpdateReplacementAttorney(t *testing.T) {
	for _, isHtmx := range []bool{false, true} {
		t.Run("Is Htmx: "+strconv.FormatBool(isHtmx), func(t *testing.T) {
			existingLpa := sirius.Lpa{
				Case: sirius.Case{
					ID:          456,
					ReceiptDate: sirius.DateString("2022-01-01"),
					ReplacementAttorneys: []sirius.Attorney{
						{Person: sirius.Person{ID: 999, Firstname: "Rudolph", Surname: "Stotesbury"}},
					},
				},
			}

			submittedLpa := sirius.Lpa{
				ApplicationHasContinuationSheets: shared.BoolPtr(false),
				ApplicationHasGuidance:           shared.BoolPtr(false),
				ApplicationHasRestrictions:       shared.BoolPtr(false),
				PaymentByDebitCreditCard:         shared.BoolPtr(false),
				PaymentRemission:                 shared.BoolPtr(false),
				RepeatApplication:                shared.BoolPtr(false),
				AnyOtherInfo:                     shared.BoolPtr(false),
				Case: sirius.Case{
					SubType:                                   "hw",
					CaseAttorneySingular:                      shared.BoolPtr(false),
					CaseAttorneyJointly:                       shared.BoolPtr(true),
					CaseAttorneyJointlyAndSeverally:           shared.BoolPtr(false),
					CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(false),
					PaymentByCheque:                           shared.BoolPtr(false),
				},
			}

			client := &mockUpdateLpaClient{}
			client.
				On("Person", mock.Anything, 123).
				Return(sirius.Person{Firstname: "Firstname", Surname: "Surname"}, nil)
			client.
				On("Lpa", mock.Anything, 456).
				Return(existingLpa, nil)
			client.
				On("UpdateLpa", mock.Anything, 456, submittedLpa).
				Return(nil)

			template := &mockTemplate{}

			form := url.Values{
				"caseAttorney":              {"jointly"},
				"caseSubtype":               {"hw"},
				"updateReplacementAttorney": {"999"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=123&caseId=456", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			if isHtmx {
				r.Header.Add("HX-Request", "true")
			}
			w := httptest.NewRecorder()

			err := UpdateLpa(client, template.Func)(w, r)
			resp := w.Result()

			expectedRedirect := RedirectError("/update-replacement-attorney?id=123&caseId=456&attorneyId=999")
			assert.Equal(t, expectedRedirect, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostUpdateLpaUpdateReplacementAttorneyBadId(t *testing.T) {
	client := &mockUpdateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("Lpa", mock.Anything, 456).
		Return(sirius.Lpa{}, nil).
		Twice()
	client.
		On("UpdateLpa", mock.Anything, 456, mock.Anything).
		Return(nil)

	form := url.Values{"updateReplacementAttorney": {"not-a-number"}}
	r, _ := http.NewRequest(http.MethodPost, "/?id=123&caseId=456", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := UpdateLpa(client, nil)(w, r)

	assert.NotNil(t, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestPostUpdateLpaErrorWhenAttorneyRadioSelected(t *testing.T) {
	dateString := "2022-04-05"

	client := &mockUpdateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{Firstname: "Firstname", Surname: "Surname"}, nil)
	client.
		On("Lpa", mock.Anything, 456).
		Return(sirius.Lpa{}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, lpaData{
			AllowNewNotifiedPerson: true,
			AppointmentType:        "singular",
			AttorneyApplicants:     []sirius.Attorney{},
			CaseId:                 456,
			DonorId:                123,
			DonorName:              "Firstname Surname",
			Error: sirius.ValidationError{
				Field: sirius.FieldErrors{"applicantType": {"": "Select an applicant"}},
			},
			FlowQuery:            "&flow=create",
			IsCreating:           true,
			IsUpdate:             true,
			ReplacementAttorneys: []sirius.Attorney{},
			Lpa: sirius.Lpa{
				Case: sirius.Case{
					ApplicationType:                           "Online",
					CaseAttorneyJointly:                       shared.BoolPtr(false),
					CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(false),
					CaseAttorneyJointlyAndSeverally:           shared.BoolPtr(false),
					CaseAttorneySingular:                      shared.BoolPtr(true),
					LpaDonorSignatureDate:                     sirius.DateString(dateString),
					PaymentByCheque:                           shared.BoolPtr(false),
					SubType:                                   "pfa",
				},
				AdditionalInfo:                   "Some extra info",
				AnyOtherInfo:                     shared.BoolPtr(true),
				ApplicantType:                    "attorney",
				ApplicationHasContinuationSheets: shared.BoolPtr(false),
				ApplicationHasGuidance:           shared.BoolPtr(true),
				ApplicationHasRestrictions:       shared.BoolPtr(false),
				AttorneyActDecisions:             "When Registered",
				CardPaymentContact:               "01234 567890",
				CertificateProviderSignatureDate: sirius.DateString(dateString),
				OnlineLpaId:                      "A12345678901",
				PaymentByDebitCreditCard:         shared.BoolPtr(true),
				PaymentRemission:                 shared.BoolPtr(false),
				RepeatApplication:                shared.BoolPtr(false),
			},
			Success:        false,
			SuccessMessage: "",
			Title:          "Create an LPA",
		}).
		Return(nil)

	form := url.Values{
		"additionalInfo":                   {"Some extra info"},
		"anyOtherInfo":                     {"true"},
		"applicantIds":                     {},
		"applicantType":                    {"attorney"},
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
		"receiptDate":                      {dateString},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123&caseId=456&flow=create", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := UpdateLpa(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}
func TestGetUpdateLpaBadQuery(t *testing.T) {
	testCases := map[string]string{
		"bad-case-id": "/?id=123&caseId=test",
	}

	for name, url := range testCases {
		t.Run(name, func(t *testing.T) {
			client := &mockUpdateLpaClient{}
			client.
				On("Person", mock.Anything, 123).
				Return(sirius.Person{}, nil)

			r, _ := http.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			err := UpdateLpa(client, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestPostUpdateLpaRedirects(t *testing.T) {
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
			name:        "Update correspondent",
			formKey:     "updateCorrespondent",
			formValue:   "true",
			expectedErr: RedirectError("/edit-correspondent?id=1&caseId=2&caseType=lpa"),
		},
		{
			name:        "Add notified person redirects",
			formKey:     "addNotifiedPerson",
			formValue:   "true",
			expectedErr: RedirectError("/create-notified-person?id=1&caseId=2"),
		},
		{
			name:        "Update notified person redirects",
			formKey:     "updateNotifiedPerson",
			formValue:   "111",
			expectedErr: RedirectError("/update-notified-person?id=1&caseId=2&notifiedPersonId=111"),
		},
		{
			name:        "Update certificate provider redirects",
			formKey:     "updateCertificateProvider",
			formValue:   "3",
			expectedErr: RedirectError("/update-certificate-provider?id=1&caseId=2&personId=3"),
		},
		{
			name:        "Update notified person with invalid ID errors",
			formKey:     "updateNotifiedPerson",
			formValue:   "not-a-number",
			expectedErr: sirius.StatusError{Code: http.StatusBadRequest},
		},
		{
			name:        "Update certificate provider with invalid ID errors",
			formKey:     "updateCertificateProvider",
			formValue:   "not-a-number",
			expectedErr: sirius.StatusError{Code: http.StatusBadRequest},
		},
		{
			name:        "Update replacement attorney redirects",
			formKey:     "updateReplacementAttorney",
			formValue:   "999",
			expectedErr: RedirectError("/update-replacement-attorney?id=1&caseId=2&attorneyId=999"),
		},
		{
			name:        "Update replacement attorney with invalid ID errors",
			formKey:     "updateReplacementAttorney",
			formValue:   "not-a-number",
			expectedErr: sirius.StatusError{Code: http.StatusBadRequest},
		},
		{
			name:        "Update attorney redirects",
			formKey:     "updateAttorney",
			formValue:   "999",
			expectedErr: RedirectError("/update-attorney?id=1&caseId=2&caseType=lpa&attorneyId=999"),
		},
		{
			name:        "Update attorney with invalid ID errors",
			formKey:     "updateAttorney",
			formValue:   "not-a-number",
			expectedErr: sirius.StatusError{Code: http.StatusBadRequest},
		},
		{
			name:        "Update attorney trust corporation redirects",
			formKey:     "updateTrustCorporationAttorney",
			formValue:   "111",
			expectedErr: RedirectError("/update-trust-corporation?id=1&caseId=2&trustCorporationId=111&replacement=false"),
		},
		{
			name:        "Update replacement attorney trust corporation redirects",
			formKey:     "updateTrustCorporationReplacementAttorney",
			formValue:   "111",
			expectedErr: RedirectError("/update-trust-corporation?id=1&caseId=2&trustCorporationId=111&replacement=true"),
		},
	}

	flowCases := []struct {
		name      string
		query     string
		flowQuery string
	}{
		{name: "LPA creation flow", query: "&caseId=2&flow=create", flowQuery: "&flow=create"},
		{name: "Edit LPA", query: "&caseId=2"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for _, flowCase := range flowCases {
				t.Run(flowCase.name, func(t *testing.T) {
					client := &mockUpdateLpaClient{}
					client.
						On("Person", mock.Anything, 1).
						Return(sirius.Person{Firstname: "John", Surname: "Doe"}, nil)

					client.
						On("Lpa", mock.Anything, 2).
						Return(sirius.Lpa{Case: sirius.Case{ID: 2}}, nil).
						Twice()
					client.
						On("UpdateLpa", mock.Anything, 2, mock.Anything).
						Return(nil)

					template := &mockTemplate{}

					form := url.Values{
						"caseSubtype": {"pfa"},
						tc.formKey:    {tc.formValue},
					}

					r, _ := http.NewRequest(http.MethodPost, "/update-lpa?id=1"+flowCase.query, strings.NewReader(form.Encode()))
					r.Header.Add("Content-Type", formUrlEncoded)
					w := httptest.NewRecorder()

					err := UpdateLpa(client, template.Func)(w, r)
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
