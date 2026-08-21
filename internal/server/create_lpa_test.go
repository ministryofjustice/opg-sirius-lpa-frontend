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

func (m *mockCreateLpaClient) UpdateLpa(ctx sirius.Context, caseID int, lpa sirius.Lpa) error {
	args := m.Called(ctx, caseID, lpa)
	return args.Error(0)
}

func (m *mockCreateLpaClient) GetUserPermissions(ctx sirius.Context) (sirius.Permissions, error) {
	args := m.Called(ctx)
	return args.Get(0).(sirius.Permissions), args.Error(1)
}

func (m *mockCreateLpaClient) UpdateAttorney(ctx sirius.Context, attorneyId int, attorney sirius.Attorney) error {
	args := m.Called(ctx, attorneyId, attorney)
	return args.Error(0)
}

func (m *mockCreateLpaClient) UpdateReplacementAttorney(ctx sirius.Context, attorneyId int, attorney sirius.Attorney) error {
	args := m.Called(ctx, attorneyId, attorney)
	return args.Error(0)
}

func (m *mockCreateLpaClient) UpdateTrustCorporation(ctx sirius.Context, trustCorporationId int, trustCorporation sirius.TrustCorporation) error {
	args := m.Called(ctx, trustCorporationId, trustCorporation)
	return args.Error(0)
}

func TestGetCreateLpa(t *testing.T) {
	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{Firstname: "Firstname", Surname: "Surname"}, nil)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(sirius.Permissions{}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createLpaData{
			DonorId:                123,
			DonorName:              "Firstname Surname",
			Title:                  "Create an LPA",
			AllowNewNotifiedPerson: true,
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
	client.
		On("GetUserPermissions", mock.Anything).
		Return(sirius.Permissions{}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createLpaData{
			DonorId:                123,
			DonorName:              "Firstname Surname",
			Title:                  "Create an LPA",
			AllowNewNotifiedPerson: true,
			IsPartial:              true,
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
	client.
		On("GetUserPermissions", mock.Anything).
		Return(sirius.Permissions{}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createLpaData{
			DonorId:                123,
			DonorName:              "Firstname Surname",
			Title:                  "Create an LPA",
			IsUpdate:               false,
			AllowNewNotifiedPerson: true,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := CreateLpa(client, template.Func)(w, r)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetCreateLpaCanEditReceiptDate(t *testing.T) {
	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{Firstname: "Firstname", Surname: "Surname"}, nil)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(sirius.Permissions{"v1-lpas-edit-dates": sirius.PermissionType{Permissions: []string{"PUT"}}}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createLpaData{
			DonorId:                123,
			DonorName:              "Firstname Surname",
			Title:                  "Create an LPA",
			CanEditReceiptDate:     true,
			AllowNewNotifiedPerson: true,
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

func TestGetCreateLpaEdit(t *testing.T) {
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
			client := &mockCreateLpaClient{}
			client.
				On("Person", mock.Anything, 123).
				Return(sirius.Person{Firstname: "Firstname", Surname: "Surname"}, nil)
			client.
				On("GetUserPermissions", mock.Anything).
				Return(sirius.Permissions{}, nil)
			client.
				On("Lpa", mock.Anything, 456).
				Return(tc.lpa, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, createLpaData{
					DonorId:                123,
					DonorName:              "Firstname Surname",
					Title:                  "Edit LPA",
					CaseId:                 456,
					Lpa:                    tc.lpa,
					AppointmentType:        tc.formValue,
					IsUpdate:               true,
					AllowNewNotifiedPerson: true,
				}).
				Return(nil)

			r, _ := http.NewRequest(http.MethodGet, "/?id=123&caseId=456", nil)
			w := httptest.NewRecorder()

			err := CreateLpa(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestGetCreateLpaEditWithTrustCorporations(t *testing.T) {
	for _, isReplacementAttorney := range []string{"true", "false"} {
		t.Run("Is Replacement Attorney"+isReplacementAttorney, func(t *testing.T) {
			lpa := sirius.Lpa{Case: sirius.Case{TrustCorporations: []sirius.TrustCorporation{{IsReplacementAttorney: isReplacementAttorney == "true"}}}}
			client := &mockCreateLpaClient{}
			client.
				On("Person", mock.Anything, 123).
				Return(sirius.Person{Firstname: "Firstname", Surname: "Surname"}, nil)
			client.
				On("GetUserPermissions", mock.Anything).
				Return(sirius.Permissions{}, nil)
			client.
				On("Lpa", mock.Anything, 456).
				Return(lpa, nil)

			data := createLpaData{
				DonorId:                123,
				DonorName:              "Firstname Surname",
				Title:                  "Edit LPA",
				CaseId:                 456,
				Lpa:                    lpa,
				IsUpdate:               true,
				AllowNewNotifiedPerson: true,
			}

			if isReplacementAttorney == "true" {
				data.ReplacementAttorneyTrustCorporations = []sirius.TrustCorporation{{IsReplacementAttorney: true}}
			} else {
				data.AttorneyTrustCorporations = []sirius.TrustCorporation{{IsReplacementAttorney: false}}
			}

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, data).
				Return(nil)

			r, _ := http.NewRequest(http.MethodGet, "/?id=123&caseId=456", nil)
			w := httptest.NewRecorder()

			err := CreateLpa(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestGetCreateLpaBadQuery(t *testing.T) {
	testCases := map[string]string{
		"no-id":       "/",
		"bad-id":      "/?id=test",
		"bad-case-id": "/?id=123&caseId=test",
	}

	for name, url := range testCases {
		t.Run(name, func(t *testing.T) {
			client := &mockCreateLpaClient{}
			client.
				On("Person", mock.Anything, 123).
				Return(sirius.Person{}, nil)
			client.
				On("GetUserPermissions", mock.Anything).
				Return(sirius.Permissions{}, nil)

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

func TestCreateLpaWhenPermissionsError(t *testing.T) {
	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(sirius.Permissions{}, errExample)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := CreateLpa(client, nil)(w, r)

	assert.Equal(t, err, errExample)
	mock.AssertExpectationsForObjects(t, client)
}

func TestCreateLpaWhenLpaErrors(t *testing.T) {
	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(sirius.Permissions{}, nil)
	client.
		On("Lpa", mock.Anything, 456).
		Return(sirius.Lpa{}, errExample)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123&caseId=456", nil)
	w := httptest.NewRecorder()

	err := CreateLpa(client, nil)(w, r)

	assert.Equal(t, err, errExample)
	mock.AssertExpectationsForObjects(t, client)
}

func TestPostCreateLpa(t *testing.T) {
	dateString := "2022-04-05"
	lpa := sirius.Lpa{
		OnlineLpaId:                "A12345678901",
		AttorneyActDecisions:       "When Registered",
		ApplicantType:              "donor",
		AnyOtherInfo:               shared.BoolPtr(true),
		AdditionalInfo:             "Some extra info",
		ApplicationHasGuidance:     shared.BoolPtr(true),
		ApplicationHasRestrictions: shared.BoolPtr(false),
		PaymentByDebitCreditCard:   shared.BoolPtr(true),
		PaymentRemission:           shared.BoolPtr(false),
		RepeatApplication:          shared.BoolPtr(false),
		CardPaymentContact:         "01234 567890",
		Case: sirius.Case{
			SubType:                         "pfa",
			ApplicationType:                 "Online",
			ReceiptDate:                     sirius.DateString(dateString),
			LpaDonorSignatureDate:           sirius.DateString(dateString),
			CaseAttorneySingular:            shared.BoolPtr(true),
			CaseAttorneyJointly:             shared.BoolPtr(false),
			CaseAttorneyJointlyAndSeverally: shared.BoolPtr(false),
			CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(false),
			PaymentByCheque:  shared.BoolPtr(false),
			PaymentExemption: shared.BoolPtr(false),
		},
	}

	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{Firstname: "Firstname", Surname: "Surname"}, nil)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(sirius.Permissions{"v1-lpas-edit-dates": sirius.PermissionType{Permissions: []string{"PUT"}}}, nil)
	client.
		On("CreateLpa", mock.Anything, 123, lpa).
		Return(sirius.Lpa{Case: sirius.Case{ID: 456}}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createLpaData{
			DonorId:                123,
			DonorName:              "Firstname Surname",
			Title:                  "Create an LPA",
			Success:                true,
			SuccessMessage:         "You have successfully created an LPA.",
			CanEditReceiptDate:     true,
			AppointmentType:        "singular",
			CaseId:                 456,
			Lpa:                    sirius.Lpa{Case: sirius.Case{ID: 456}},
			AllowNewNotifiedPerson: true,
		}).
		Return(nil)

	form := url.Values{
		"caseSubtype":                {"pfa"},
		"applicationType":            {"Online"},
		"onlineLpaId":                {"A12345678901"},
		"receiptDate":                {dateString},
		"lpaDonorSignatureDate":      {dateString},
		"caseAttorney":               {"singular"},
		"attorneyActDecisions":       {"When Registered"},
		"preferencesAndInstructions": {"guidance"},
		"applicantType":              {"donor"},
		"applicationFee":             {"card"},
		"cardPaymentContact":         {"01234 567890"},
		"anyOtherInfo":               {"true"},
		"additionalInfo":             {"Some extra info"},
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
		ApplicationHasGuidance:                    shared.BoolPtr(false),
		ApplicationHasRestrictions:                shared.BoolPtr(false),
		PaymentByDebitCreditCard:                  shared.BoolPtr(false),
		PaymentRemission:                          shared.BoolPtr(false),
		RepeatApplication:                         shared.BoolPtr(false),
		AnyOtherInfo:                              shared.BoolPtr(false),
		LifeSustainingTreatmentSignedAndWitnessed: shared.BoolPtr(false),
		Case: sirius.Case{
			SubType:                                   "hw",
			LifeSustainingTreatment:                   "Option A",
			CaseAttorneySingular:                      shared.BoolPtr(true),
			CaseAttorneyJointly:                       shared.BoolPtr(false),
			CaseAttorneyJointlyAndSeverally:           shared.BoolPtr(false),
			CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(false),
			PaymentByCheque:                           shared.BoolPtr(false),
			PaymentExemption:                          shared.BoolPtr(false),
		},
	}

	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(sirius.Permissions{}, nil)
	client.
		On("CreateLpa", mock.Anything, 123, lpa).
		Return(sirius.Lpa{}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.Anything).
		Return(nil)

	form := url.Values{
		"caseSubtype":             {"hw"},
		"lifeSustainingTreatment": {"Option A"},
		"attorneyActDecisions":    {"When Registered"},
		"caseAttorney":            {"singular"},
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
		ApplicationHasGuidance:     shared.BoolPtr(false),
		ApplicationHasRestrictions: shared.BoolPtr(false),
		PaymentByDebitCreditCard:   shared.BoolPtr(false),
		PaymentRemission:           shared.BoolPtr(false),
		RepeatApplication:          shared.BoolPtr(false),
		AnyOtherInfo:               shared.BoolPtr(false),
		Case: sirius.Case{
			CaseAttorneySingular:                      shared.BoolPtr(true),
			CaseAttorneyJointly:                       shared.BoolPtr(false),
			CaseAttorneyJointlyAndSeverally:           shared.BoolPtr(false),
			CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(false),
			PaymentByCheque:                           shared.BoolPtr(false),
			PaymentExemption:                          shared.BoolPtr(false),
		},
	}

	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(sirius.Permissions{}, nil)
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
		ApplicationHasGuidance:     shared.BoolPtr(false),
		ApplicationHasRestrictions: shared.BoolPtr(false),
		PaymentByDebitCreditCard:   shared.BoolPtr(false),
		PaymentRemission:           shared.BoolPtr(false),
		RepeatApplication:          shared.BoolPtr(false),
		AnyOtherInfo:               shared.BoolPtr(false),
		Case: sirius.Case{
			SubType:                                   "pfa",
			ApplicationType:                           "Classic",
			CaseAttorneySingular:                      shared.BoolPtr(true),
			CaseAttorneyJointly:                       shared.BoolPtr(false),
			CaseAttorneyJointlyAndSeverally:           shared.BoolPtr(false),
			CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(false),
			PaymentByCheque:                           shared.BoolPtr(false),
			PaymentExemption:                          shared.BoolPtr(false),
		},
	}

	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(sirius.Permissions{}, nil)
	client.
		On("CreateLpa", mock.Anything, 123, lpa).
		Return(sirius.Lpa{}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.Anything).
		Return(nil)

	form := url.Values{
		"caseSubtype":     {"pfa"},
		"applicationType": {"Classic"},
		"onlineLpaId":     {"A12345678901"},
		"caseAttorney":    {"singular"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateLpa(client, template.Func)(w, r)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateLpaIgnoresReceiptDateWithoutPermission(t *testing.T) {
	existingLpa := sirius.Lpa{Case: sirius.Case{ID: 456, ReceiptDate: sirius.DateString("2022-01-01")}}
	submittedLpa := sirius.Lpa{
		ApplicationHasGuidance:                    shared.BoolPtr(false),
		ApplicationHasRestrictions:                shared.BoolPtr(false),
		PaymentByDebitCreditCard:                  shared.BoolPtr(false),
		PaymentRemission:                          shared.BoolPtr(false),
		RepeatApplication:                         shared.BoolPtr(false),
		AnyOtherInfo:                              shared.BoolPtr(false),
		LifeSustainingTreatmentSignedAndWitnessed: shared.BoolPtr(false),
		Case: sirius.Case{
			SubType:                         "hw",
			ReceiptDate:                     sirius.DateString("2022-01-01"),
			CaseAttorneySingular:            shared.BoolPtr(false),
			CaseAttorneyJointly:             shared.BoolPtr(true),
			CaseAttorneyJointlyAndSeverally: shared.BoolPtr(false),
			CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(false),
			PaymentByCheque:  shared.BoolPtr(false),
			PaymentExemption: shared.BoolPtr(false),
		},
	}

	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(sirius.Permissions{}, nil)
	client.
		On("Lpa", mock.Anything, 456).
		Return(existingLpa, nil).
		Once()
	client.
		On("UpdateLpa", mock.Anything, 456, submittedLpa).
		Return(nil)
	client.
		On("Lpa", mock.Anything, 456).
		Return(sirius.Lpa{Case: sirius.Case{ID: 456}}, nil).
		Once()

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.Anything).
		Return(nil)

	form := url.Values{
		"caseSubtype":  {"hw"},
		"receiptDate":  {"2099-01-01"},
		"caseAttorney": {"jointly"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123&caseId=456", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateLpa(client, template.Func)(w, r)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateLpaDropsCardPaymentContactWhenCardNotSelected(t *testing.T) {
	lpa := sirius.Lpa{
		ApplicationHasGuidance:     shared.BoolPtr(false),
		ApplicationHasRestrictions: shared.BoolPtr(false),
		PaymentByDebitCreditCard:   shared.BoolPtr(false),
		PaymentRemission:           shared.BoolPtr(false),
		RepeatApplication:          shared.BoolPtr(false),
		AnyOtherInfo:               shared.BoolPtr(false),
		Case: sirius.Case{
			CaseAttorneySingular:                      shared.BoolPtr(true),
			CaseAttorneyJointly:                       shared.BoolPtr(false),
			CaseAttorneyJointlyAndSeverally:           shared.BoolPtr(false),
			CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(false),
			PaymentByCheque:                           shared.BoolPtr(true),
			PaymentExemption:                          shared.BoolPtr(false),
		},
	}

	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(sirius.Permissions{}, nil)
	client.
		On("CreateLpa", mock.Anything, 123, lpa).
		Return(sirius.Lpa{}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.Anything).
		Return(nil)

	form := url.Values{
		"caseAttorney":       {"singular"},
		"applicationFee":     {"cheque"}, // not "card"
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
		ApplicationHasGuidance:     shared.BoolPtr(false),
		ApplicationHasRestrictions: shared.BoolPtr(false),
		PaymentByDebitCreditCard:   shared.BoolPtr(false),
		PaymentRemission:           shared.BoolPtr(false),
		RepeatApplication:          shared.BoolPtr(false),
		AnyOtherInfo:               shared.BoolPtr(false),
		Case: sirius.Case{
			CaseAttorneySingular:                      shared.BoolPtr(true),
			CaseAttorneyJointly:                       shared.BoolPtr(false),
			CaseAttorneyJointlyAndSeverally:           shared.BoolPtr(false),
			CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(false),
			PaymentByCheque:                           shared.BoolPtr(false),
			PaymentExemption:                          shared.BoolPtr(false),
		},
	}

	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(sirius.Permissions{}, nil)
	client.
		On("CreateLpa", mock.Anything, 123, lpa).
		Return(sirius.Lpa{}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.Anything).
		Return(nil)

	form := url.Values{
		"caseAttorney":   {"singular"},
		"additionalInfo": {"Some info the user typed then unchecked yes"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateLpa(client, template.Func)(w, r)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateLpaApplicationFeeReducedFeeExemption(t *testing.T) {
	lpa := sirius.Lpa{
		ApplicationHasGuidance:     shared.BoolPtr(false),
		ApplicationHasRestrictions: shared.BoolPtr(false),
		PaymentByDebitCreditCard:   shared.BoolPtr(false),
		PaymentRemission:           shared.BoolPtr(false),
		RepeatApplication:          shared.BoolPtr(true),
		AnyOtherInfo:               shared.BoolPtr(false),
		Case: sirius.Case{
			CaseAttorneySingular:                      shared.BoolPtr(true),
			CaseAttorneyJointly:                       shared.BoolPtr(false),
			CaseAttorneyJointlyAndSeverally:           shared.BoolPtr(false),
			CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(false),
			PaymentByCheque:                           shared.BoolPtr(true),
			PaymentExemption:                          shared.BoolPtr(true),
		},
	}

	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(sirius.Permissions{}, nil)
	client.
		On("CreateLpa", mock.Anything, 123, lpa).
		Return(sirius.Lpa{}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.Anything).
		Return(nil)

	form := url.Values{
		"caseAttorney":   {"singular"},
		"applicationFee": {"cheque", "reducedFee", "repeatApplication"},
		"reducedFeeType": {"exemption"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateLpa(client, template.Func)(w, r)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateLpaApplicationFeeReducedFeeRemission(t *testing.T) {
	lpa := sirius.Lpa{
		ApplicationHasGuidance:     shared.BoolPtr(false),
		ApplicationHasRestrictions: shared.BoolPtr(false),
		PaymentByDebitCreditCard:   shared.BoolPtr(false),
		PaymentRemission:           shared.BoolPtr(true),
		RepeatApplication:          shared.BoolPtr(false),
		AnyOtherInfo:               shared.BoolPtr(false),
		Case: sirius.Case{
			CaseAttorneySingular:                      shared.BoolPtr(true),
			CaseAttorneyJointly:                       shared.BoolPtr(false),
			CaseAttorneyJointlyAndSeverally:           shared.BoolPtr(false),
			CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(false),
			PaymentByCheque:                           shared.BoolPtr(false),
			PaymentExemption:                          shared.BoolPtr(false),
		},
	}

	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(sirius.Permissions{}, nil)
	client.
		On("CreateLpa", mock.Anything, 123, lpa).
		Return(sirius.Lpa{}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.Anything).
		Return(nil)

	form := url.Values{
		"caseAttorney":   {"singular"},
		"applicationFee": {"reducedFee"},
		"reducedFeeType": {"remission"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateLpa(client, template.Func)(w, r)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateLpaApplicationFeeReducedFeeTypeIgnoredWhenNotSelected(t *testing.T) {
	lpa := sirius.Lpa{
		ApplicationHasGuidance:     shared.BoolPtr(false),
		ApplicationHasRestrictions: shared.BoolPtr(false),
		PaymentByDebitCreditCard:   shared.BoolPtr(false),
		PaymentRemission:           shared.BoolPtr(false),
		RepeatApplication:          shared.BoolPtr(false),
		AnyOtherInfo:               shared.BoolPtr(false),
		Case: sirius.Case{
			CaseAttorneySingular:                      shared.BoolPtr(true),
			CaseAttorneyJointly:                       shared.BoolPtr(false),
			CaseAttorneyJointlyAndSeverally:           shared.BoolPtr(false),
			CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(false),
			PaymentByCheque:                           shared.BoolPtr(false),
			PaymentExemption:                          shared.BoolPtr(false),
		},
	}

	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(sirius.Permissions{}, nil)
	client.
		On("CreateLpa", mock.Anything, 123, lpa).
		Return(sirius.Lpa{}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.Anything).
		Return(nil)

	form := url.Values{
		"caseAttorney":   {"singular"},
		"reducedFeeType": {"exemption"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateLpa(client, template.Func)(w, r)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateLpaApplicantAndLifeSustainingTreatmentFields(t *testing.T) {
	lpa := sirius.Lpa{
		ApplicationHasGuidance:                    shared.BoolPtr(false),
		ApplicationHasRestrictions:                shared.BoolPtr(false),
		PaymentByDebitCreditCard:                  shared.BoolPtr(false),
		PaymentRemission:                          shared.BoolPtr(false),
		RepeatApplication:                         shared.BoolPtr(false),
		AnyOtherInfo:                              shared.BoolPtr(false),
		LifeSustainingTreatmentSignedAndWitnessed: shared.BoolPtr(true),
		ApplicantType:                             "attorney",
		ApplicantSignatureDate:                    "2022-04-06",
		ApplicantIds:                              []int{1, 2},
		Case: sirius.Case{
			SubType:                                   "hw",
			LifeSustainingTreatment:                   "None",
			LifeSustainingTreatmentSignatureDateA:     "2022-04-05",
			CaseAttorneySingular:                      shared.BoolPtr(false),
			CaseAttorneyJointly:                       shared.BoolPtr(true),
			CaseAttorneyJointlyAndSeverally:           shared.BoolPtr(false),
			CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(false),
			PaymentByCheque:                           shared.BoolPtr(false),
			PaymentExemption:                          shared.BoolPtr(false),
		},
	}

	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(sirius.Permissions{}, nil)
	client.
		On("CreateLpa", mock.Anything, 123, lpa).
		Return(sirius.Lpa{}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.Anything).
		Return(nil)

	form := url.Values{
		"caseSubtype":                               {"hw"},
		"caseAttorney":                              {"jointly"},
		"lifeSustainingTreatment":                   {"None"},
		"lifeSustainingTreatmentSignatureDate":      {"2022-04-05"},
		"lifeSustainingTreatmentSignedAndWitnessed": {"true"},
		"applicantType":                             {"attorney"},
		"applicantSignatureDate":                    {"2022-04-06"},
		"applicantIds":                              {"1", "2"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateLpa(client, template.Func)(w, r)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateLpaEditAttorneySignatureDates(t *testing.T) {
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
		ApplicationHasGuidance:                    shared.BoolPtr(false),
		ApplicationHasRestrictions:                shared.BoolPtr(false),
		PaymentByDebitCreditCard:                  shared.BoolPtr(false),
		PaymentRemission:                          shared.BoolPtr(false),
		RepeatApplication:                         shared.BoolPtr(false),
		AnyOtherInfo:                              shared.BoolPtr(false),
		LifeSustainingTreatmentSignedAndWitnessed: shared.BoolPtr(false),
		Case: sirius.Case{
			SubType:                         "hw",
			ReceiptDate:                     sirius.DateString("2022-01-01"),
			CaseAttorneySingular:            shared.BoolPtr(false),
			CaseAttorneyJointly:             shared.BoolPtr(true),
			CaseAttorneyJointlyAndSeverally: shared.BoolPtr(false),
			CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(false),
			PaymentByCheque:  shared.BoolPtr(false),
			PaymentExemption: shared.BoolPtr(false),
		},
	}

	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(sirius.Permissions{}, nil)
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
		"caseSubtype":               {"hw"},
		"receiptDate":               {"2099-01-01"},
		"caseAttorney":              {"jointly"},
		"lpaPartCSignatureDate-876": {"2022-01-02"},
		"lpaPartCSignatureDate-987": {"2022-01-02"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123&caseId=456", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateLpa(client, template.Func)(w, r)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateLpaEditReplacementAttorneySignatureDates(t *testing.T) {
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
		ApplicationHasGuidance:                    shared.BoolPtr(false),
		ApplicationHasRestrictions:                shared.BoolPtr(false),
		PaymentByDebitCreditCard:                  shared.BoolPtr(false),
		PaymentRemission:                          shared.BoolPtr(false),
		RepeatApplication:                         shared.BoolPtr(false),
		AnyOtherInfo:                              shared.BoolPtr(false),
		LifeSustainingTreatmentSignedAndWitnessed: shared.BoolPtr(false),
		Case: sirius.Case{
			SubType:                         "hw",
			ReceiptDate:                     sirius.DateString("2022-01-01"),
			CaseAttorneySingular:            shared.BoolPtr(false),
			CaseAttorneyJointly:             shared.BoolPtr(true),
			CaseAttorneyJointlyAndSeverally: shared.BoolPtr(false),
			CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(false),
			PaymentByCheque:  shared.BoolPtr(false),
			PaymentExemption: shared.BoolPtr(false),
		},
	}

	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(sirius.Permissions{}, nil)
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
		"caseSubtype":               {"hw"},
		"receiptDate":               {"2099-01-01"},
		"caseAttorney":              {"jointly"},
		"lpaPartCSignatureDate-876": {"2022-01-02"},
		"lpaPartCSignatureDate-987": {"2022-01-02"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123&caseId=456", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateLpa(client, template.Func)(w, r)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateLpaEditTrustCorporationSignatureDates(t *testing.T) {
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
		ApplicationHasGuidance:                    shared.BoolPtr(false),
		ApplicationHasRestrictions:                shared.BoolPtr(false),
		PaymentByDebitCreditCard:                  shared.BoolPtr(false),
		PaymentRemission:                          shared.BoolPtr(false),
		RepeatApplication:                         shared.BoolPtr(false),
		AnyOtherInfo:                              shared.BoolPtr(false),
		LifeSustainingTreatmentSignedAndWitnessed: shared.BoolPtr(false),
		Case: sirius.Case{
			SubType:                         "hw",
			ReceiptDate:                     sirius.DateString("2022-01-01"),
			CaseAttorneySingular:            shared.BoolPtr(false),
			CaseAttorneyJointly:             shared.BoolPtr(true),
			CaseAttorneyJointlyAndSeverally: shared.BoolPtr(false),
			CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(false),
			PaymentByCheque:  shared.BoolPtr(false),
			PaymentExemption: shared.BoolPtr(false),
		},
	}

	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(sirius.Permissions{}, nil)
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
		"caseSubtype":               {"hw"},
		"receiptDate":               {"2099-01-01"},
		"caseAttorney":              {"jointly"},
		"lpaPartCSignatureDate-876": {"2022-01-02"},
		"lpaPartCSignatureDate-987": {"2022-01-02"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123&caseId=456", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateLpa(client, template.Func)(w, r)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateLpaEditAttorneySignatureDatesError(t *testing.T) {
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
		ApplicationHasGuidance:                    shared.BoolPtr(false),
		ApplicationHasRestrictions:                shared.BoolPtr(false),
		PaymentByDebitCreditCard:                  shared.BoolPtr(false),
		PaymentRemission:                          shared.BoolPtr(false),
		RepeatApplication:                         shared.BoolPtr(false),
		AnyOtherInfo:                              shared.BoolPtr(false),
		LifeSustainingTreatmentSignedAndWitnessed: shared.BoolPtr(false),
		Case: sirius.Case{
			SubType:                         "hw",
			ReceiptDate:                     sirius.DateString("2022-01-01"),
			CaseAttorneySingular:            shared.BoolPtr(false),
			CaseAttorneyJointly:             shared.BoolPtr(true),
			CaseAttorneyJointlyAndSeverally: shared.BoolPtr(false),
			CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(false),
			PaymentByCheque:  shared.BoolPtr(false),
			PaymentExemption: shared.BoolPtr(false),
		},
	}

	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(sirius.Permissions{}, nil)
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
		"caseSubtype":               {"hw"},
		"receiptDate":               {"2099-01-01"},
		"caseAttorney":              {"jointly"},
		"lpaPartCSignatureDate-876": {"2022-01-02"},
		"lpaPartCSignatureDate-987": {"2022-01-02"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123&caseId=456", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateLpa(client, nil)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestPostCreateLpaEditReplacementAttorneySignatureDatesError(t *testing.T) {
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
		ApplicationHasGuidance:                    shared.BoolPtr(false),
		ApplicationHasRestrictions:                shared.BoolPtr(false),
		PaymentByDebitCreditCard:                  shared.BoolPtr(false),
		PaymentRemission:                          shared.BoolPtr(false),
		RepeatApplication:                         shared.BoolPtr(false),
		AnyOtherInfo:                              shared.BoolPtr(false),
		LifeSustainingTreatmentSignedAndWitnessed: shared.BoolPtr(false),
		Case: sirius.Case{
			SubType:                         "hw",
			ReceiptDate:                     sirius.DateString("2022-01-01"),
			CaseAttorneySingular:            shared.BoolPtr(false),
			CaseAttorneyJointly:             shared.BoolPtr(true),
			CaseAttorneyJointlyAndSeverally: shared.BoolPtr(false),
			CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(false),
			PaymentByCheque:  shared.BoolPtr(false),
			PaymentExemption: shared.BoolPtr(false),
		},
	}

	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(sirius.Permissions{}, nil)
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
		"caseSubtype":               {"hw"},
		"receiptDate":               {"2099-01-01"},
		"caseAttorney":              {"jointly"},
		"lpaPartCSignatureDate-876": {"2022-01-02"},
		"lpaPartCSignatureDate-987": {"2022-01-02"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123&caseId=456", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateLpa(client, nil)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestPostCreateLpaEditTrustCorporationSignatureDatesError(t *testing.T) {
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
		ApplicationHasGuidance:                    shared.BoolPtr(false),
		ApplicationHasRestrictions:                shared.BoolPtr(false),
		PaymentByDebitCreditCard:                  shared.BoolPtr(false),
		PaymentRemission:                          shared.BoolPtr(false),
		RepeatApplication:                         shared.BoolPtr(false),
		AnyOtherInfo:                              shared.BoolPtr(false),
		LifeSustainingTreatmentSignedAndWitnessed: shared.BoolPtr(false),
		Case: sirius.Case{
			SubType:                         "hw",
			ReceiptDate:                     sirius.DateString("2022-01-01"),
			CaseAttorneySingular:            shared.BoolPtr(false),
			CaseAttorneyJointly:             shared.BoolPtr(true),
			CaseAttorneyJointlyAndSeverally: shared.BoolPtr(false),
			CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(false),
			PaymentByCheque:  shared.BoolPtr(false),
			PaymentExemption: shared.BoolPtr(false),
		},
	}

	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(sirius.Permissions{}, nil)
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
		"caseSubtype":               {"hw"},
		"receiptDate":               {"2099-01-01"},
		"caseAttorney":              {"jointly"},
		"lpaPartCSignatureDate-876": {"2022-01-02"},
		"lpaPartCSignatureDate-987": {"2022-01-02"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123&caseId=456", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateLpa(client, nil)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
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
		On("GetUserPermissions", mock.Anything).
		Return(sirius.Permissions{"v1-lpas-edit-dates": sirius.PermissionType{Permissions: []string{"PUT"}}}, nil)
	client.
		On("CreateLpa", mock.Anything, 123, mock.Anything).
		Return(sirius.Lpa{}, expectedError)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.MatchedBy(func(d createLpaData) bool {
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
		On("GetUserPermissions", mock.Anything).
		Return(sirius.Permissions{"v1-lpas-edit-dates": sirius.PermissionType{Permissions: []string{"PUT"}}}, nil)
	client.
		On("CreateLpa", mock.Anything, 123, mock.Anything).
		Return(sirius.Lpa{}, expectedError)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.MatchedBy(func(d createLpaData) bool {
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

func TestPostCreateLpaEditWhenValidationError(t *testing.T) {
	expectedError := sirius.ValidationError{
		Field: sirius.FieldErrors{"receiptDate": {"receiptDate": "Select the date of receipt"}},
	}

	existingLpa := sirius.Lpa{Case: sirius.Case{ID: 456, ReceiptDate: sirius.DateString("2022-01-01")}}

	submittedLpa := sirius.Lpa{
		ApplicationHasGuidance:     shared.BoolPtr(false),
		ApplicationHasRestrictions: shared.BoolPtr(false),
		PaymentByDebitCreditCard:   shared.BoolPtr(false),
		PaymentRemission:           shared.BoolPtr(false),
		RepeatApplication:          shared.BoolPtr(false),
		AnyOtherInfo:               shared.BoolPtr(false),
		Case: sirius.Case{
			CaseAttorneySingular:                      shared.BoolPtr(true),
			CaseAttorneyJointly:                       shared.BoolPtr(false),
			CaseAttorneyJointlyAndSeverally:           shared.BoolPtr(false),
			CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(false),
			PaymentByCheque:                           shared.BoolPtr(false),
			PaymentExemption:                          shared.BoolPtr(false),
		},
	}

	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(sirius.Permissions{"v1-lpas-edit-dates": sirius.PermissionType{Permissions: []string{"PUT"}}}, nil)
	client.
		On("Lpa", mock.Anything, 456).
		Return(existingLpa, nil).
		Once()
	client.
		On("UpdateLpa", mock.Anything, 456, submittedLpa).
		Return(expectedError)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.MatchedBy(func(d createLpaData) bool {
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
		On("GetUserPermissions", mock.Anything).
		Return(sirius.Permissions{}, nil)
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
				OnlineLpaId:                "A12345678901",
				AttorneyActDecisions:       "When Registered",
				ApplicantType:              "donor",
				AnyOtherInfo:               shared.BoolPtr(true),
				AdditionalInfo:             "Some extra info",
				ApplicationHasGuidance:     shared.BoolPtr(true),
				ApplicationHasRestrictions: shared.BoolPtr(false),
				PaymentByDebitCreditCard:   shared.BoolPtr(true),
				PaymentRemission:           shared.BoolPtr(false),
				RepeatApplication:          shared.BoolPtr(false),
				CardPaymentContact:         "01234 567890",
				Case: sirius.Case{
					SubType:                         "pfa",
					ApplicationType:                 "Online",
					ReceiptDate:                     sirius.DateString(dateString),
					LpaDonorSignatureDate:           sirius.DateString(dateString),
					CaseAttorneySingular:            shared.BoolPtr(true),
					CaseAttorneyJointly:             shared.BoolPtr(false),
					CaseAttorneyJointlyAndSeverally: shared.BoolPtr(false),
					CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(false),
					PaymentByCheque:  shared.BoolPtr(false),
					PaymentExemption: shared.BoolPtr(false),
				},
			}

			client := &mockCreateLpaClient{}
			client.
				On("Person", mock.Anything, 123).
				Return(sirius.Person{Firstname: "Firstname", Surname: "Surname"}, nil)
			client.
				On("GetUserPermissions", mock.Anything).
				Return(sirius.Permissions{"v1-lpas-edit-dates": sirius.PermissionType{Permissions: []string{"PUT"}}}, nil)
			client.
				On("CreateLpa", mock.Anything, 123, lpa).
				Return(sirius.Lpa{Case: sirius.Case{ID: 456}}, nil)

			template := &mockTemplate{}

			expectedData := createLpaData{
				AllowNewNotifiedPerson: true,
				DonorId:                123,
				DonorName:              "Firstname Surname",
				Title:                  "Create an LPA",
				Success:                true,
				SuccessMessage:         "You have successfully created an LPA.",
				CanEditReceiptDate:     true,
				AppointmentType:        "singular",
				CaseId:                 456,
				Lpa:                    sirius.Lpa{Case: sirius.Case{ID: 456}},
				HtmxRedirect:           "/create-replacement-attorney?id=123&caseId=456",
				HtmxSwap:               "innerHTML",
				IsPartial:              isHtmx,
			}

			if isHtmx {
				template.
					On("Func", mock.Anything, expectedData).
					Return(nil)
			}

			form := url.Values{
				"caseSubtype":                {"pfa"},
				"applicationType":            {"Online"},
				"onlineLpaId":                {"A12345678901"},
				"receiptDate":                {dateString},
				"lpaDonorSignatureDate":      {dateString},
				"caseAttorney":               {"singular"},
				"attorneyActDecisions":       {"When Registered"},
				"preferencesAndInstructions": {"guidance"},
				"applicantType":              {"donor"},
				"applicationFee":             {"card"},
				"cardPaymentContact":         {"01234 567890"},
				"anyOtherInfo":               {"true"},
				"additionalInfo":             {"Some extra info"},
				"addReplacementAttorney":     {"Add replacement attorney"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			if isHtmx {
				r.Header.Add("HX-Request", "true")
			}
			w := httptest.NewRecorder()

			err := CreateLpa(client, template.Func)(w, r)
			resp := w.Result()

			if !isHtmx {
				expectedRedirect := RedirectError("/create-replacement-attorney?id=123&caseId=456")
				assert.Equal(t, expectedRedirect, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostCreateLpaUpdateAttorney(t *testing.T) {
	for _, isHtmx := range []bool{false, true} {
		t.Run("Is Htmx: "+strconv.FormatBool(isHtmx), func(t *testing.T) {
			existingLpa := sirius.Lpa{
				Case: sirius.Case{
					ID:          456,
					ReceiptDate: sirius.DateString("2022-01-01"),
					Attorneys: []sirius.Attorney{
						{Person: sirius.Person{ID: 999, Firstname: "Rudolph", Surname: "Stotesbury"}},
					},
				},
			}

			submittedLpa := sirius.Lpa{
				ApplicationHasGuidance:                    shared.BoolPtr(false),
				ApplicationHasRestrictions:                shared.BoolPtr(false),
				PaymentByDebitCreditCard:                  shared.BoolPtr(false),
				PaymentRemission:                          shared.BoolPtr(false),
				RepeatApplication:                         shared.BoolPtr(false),
				AnyOtherInfo:                              shared.BoolPtr(false),
				LifeSustainingTreatmentSignedAndWitnessed: shared.BoolPtr(false),
				Case: sirius.Case{
					SubType:                         "hw",
					ReceiptDate:                     sirius.DateString("2022-01-01"),
					CaseAttorneySingular:            shared.BoolPtr(false),
					CaseAttorneyJointly:             shared.BoolPtr(true),
					CaseAttorneyJointlyAndSeverally: shared.BoolPtr(false),
					CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(false),
					PaymentByCheque:  shared.BoolPtr(false),
					PaymentExemption: shared.BoolPtr(false),
				},
			}

			client := &mockCreateLpaClient{}
			client.
				On("Person", mock.Anything, 123).
				Return(sirius.Person{Firstname: "Firstname", Surname: "Surname"}, nil)
			client.
				On("GetUserPermissions", mock.Anything).
				Return(sirius.Permissions{}, nil)
			client.
				On("Lpa", mock.Anything, 456).
				Return(existingLpa, nil)
			client.
				On("UpdateLpa", mock.Anything, 456, submittedLpa).
				Return(nil)

			template := &mockTemplate{}

			expectedData := createLpaData{
				AllowNewNotifiedPerson: true,
				DonorId:                123,
				DonorName:              "Firstname Surname",
				Title:                  "Edit LPA",
				IsUpdate:               true,
				Success:                true,
				SuccessMessage:         "You have successfully updated an LPA.",
				AppointmentType:        "jointly",
				CaseId:                 456,
				Lpa:                    existingLpa,
				HtmxRedirect:           "/create-attorney?id=123&caseId=456&caseType=lpa&attorneyId=999",
				HtmxSwap:               "innerHTML",
				IsPartial:              isHtmx,
			}

			if isHtmx {
				template.
					On("Func", mock.Anything, expectedData).
					Return(nil)
			}

			form := url.Values{
				"caseSubtype":    {"hw"},
				"caseAttorney":   {"jointly"},
				"updateAttorney": {"999"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=123&caseId=456", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			if isHtmx {
				r.Header.Add("HX-Request", "true")
			}
			w := httptest.NewRecorder()

			err := CreateLpa(client, template.Func)(w, r)
			resp := w.Result()

			if !isHtmx {
				expectedRedirect := RedirectError("/create-attorney?id=123&caseId=456&caseType=lpa&attorneyId=999")
				assert.Equal(t, expectedRedirect, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostCreateLpaUpdateAttorneyBadId(t *testing.T) {
	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(sirius.Permissions{}, nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123&updateAttorney=not-a-number", nil)
	w := httptest.NewRecorder()

	err := CreateLpa(client, nil)(w, r)

	assert.NotNil(t, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestPostCreateLpaUpdateTrustCorporationBadId(t *testing.T) {
	for _, isReplacementAttorney := range []bool{false, true} {
		t.Run("Is Replacement Attorney: "+strconv.FormatBool(isReplacementAttorney), func(t *testing.T) {
			client := &mockCreateLpaClient{}
			client.
				On("Person", mock.Anything, 1).
				Return(sirius.Person{Firstname: "John", Surname: "Doe"}, nil)
			client.
				On("GetUserPermissions", mock.Anything).
				Return(sirius.Permissions{}, nil)
			client.
				On("CreateLpa", mock.Anything, 1, mock.Anything).
				Return(sirius.Lpa{Case: sirius.Case{ID: 2}}, nil)

			form := url.Values{
				"caseSubtype": {"pfa"},
			}

			if isReplacementAttorney {
				form["updateTrustCorporationReplacementAttorney"] = []string{"bad-id"}
			} else {
				form["updateTrustCorporationAttorney"] = []string{"bad-id"}
			}

			r, _ := http.NewRequest(http.MethodPost, "/create-lpa?id=1", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := CreateLpa(client, nil)(w, r)

			assert.NotNil(t, err)
			mock.AssertExpectationsForObjects(t, client)
		})
	}
}

func TestPostCreateLpaUpdateReplacementAttorney(t *testing.T) {
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
				ApplicationHasGuidance:                    shared.BoolPtr(false),
				ApplicationHasRestrictions:                shared.BoolPtr(false),
				PaymentByDebitCreditCard:                  shared.BoolPtr(false),
				PaymentRemission:                          shared.BoolPtr(false),
				RepeatApplication:                         shared.BoolPtr(false),
				AnyOtherInfo:                              shared.BoolPtr(false),
				LifeSustainingTreatmentSignedAndWitnessed: shared.BoolPtr(false),
				Case: sirius.Case{
					SubType:                         "hw",
					ReceiptDate:                     sirius.DateString("2022-01-01"),
					CaseAttorneySingular:            shared.BoolPtr(false),
					CaseAttorneyJointly:             shared.BoolPtr(true),
					CaseAttorneyJointlyAndSeverally: shared.BoolPtr(false),
					CaseAttorneyJointlyAndJointlyAndSeverally: shared.BoolPtr(false),
					PaymentByCheque:  shared.BoolPtr(false),
					PaymentExemption: shared.BoolPtr(false),
				},
			}

			client := &mockCreateLpaClient{}
			client.
				On("Person", mock.Anything, 123).
				Return(sirius.Person{Firstname: "Firstname", Surname: "Surname"}, nil)
			client.
				On("GetUserPermissions", mock.Anything).
				Return(sirius.Permissions{}, nil)
			client.
				On("Lpa", mock.Anything, 456).
				Return(existingLpa, nil)
			client.
				On("UpdateLpa", mock.Anything, 456, submittedLpa).
				Return(nil)

			template := &mockTemplate{}

			expectedData := createLpaData{
				AllowNewNotifiedPerson: true,
				DonorId:                123,
				DonorName:              "Firstname Surname",
				Title:                  "Edit LPA",
				IsUpdate:               true,
				Success:                true,
				SuccessMessage:         "You have successfully updated an LPA.",
				AppointmentType:        "jointly",
				CaseId:                 456,
				Lpa:                    existingLpa,
				HtmxRedirect:           "/create-replacement-attorney?id=123&caseId=456&attorneyId=999",
				HtmxSwap:               "innerHTML",
				IsPartial:              isHtmx,
			}

			if isHtmx {
				template.
					On("Func", mock.Anything, expectedData).
					Return(nil)
			}

			form := url.Values{
				"caseSubtype":               {"hw"},
				"caseAttorney":              {"jointly"},
				"updateReplacementAttorney": {"999"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=123&caseId=456", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			if isHtmx {
				r.Header.Add("HX-Request", "true")
			}
			w := httptest.NewRecorder()

			err := CreateLpa(client, template.Func)(w, r)
			resp := w.Result()

			if !isHtmx {
				expectedRedirect := RedirectError("/create-replacement-attorney?id=123&caseId=456&attorneyId=999")
				assert.Equal(t, expectedRedirect, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostCreateLpaUpdateReplacementAttorneyBadId(t *testing.T) {
	client := &mockCreateLpaClient{}
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(sirius.Permissions{}, nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123&updateReplacementAttorney=not-a-number", nil)
	w := httptest.NewRecorder()

	err := CreateLpa(client, nil)(w, r)

	assert.NotNil(t, err)
	mock.AssertExpectationsForObjects(t, client)
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
			name:        "Update correspondent",
			formKey:     "updateCorrespondent",
			formValue:   "true",
			expectedErr: RedirectError("/create-correspondent?id=1&caseId=2&caseType=lpa"),
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
			expectedErr: RedirectError("/create-notified-person?id=1&caseId=2&notifiedPersonId=111"),
		},
		{
			name:        "Update certificate provider redirects",
			formKey:     "updateCertificateProvider",
			formValue:   "3",
			expectedErr: RedirectError("/edit-certificate-provider?id=1&caseId=2&personId=3"),
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
			expectedErr: RedirectError("/create-replacement-attorney?id=1&caseId=2&attorneyId=999"),
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
			expectedErr: RedirectError("/create-attorney?id=1&caseId=2&caseType=lpa&attorneyId=999"),
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
			expectedErr: RedirectError("/create-trust-corporation?id=1&caseId=2&trustCorporationId=111&replacement=false"),
		},
		{
			name:        "Update replacement attorney trust corporation redirects",
			formKey:     "updateTrustCorporationReplacementAttorney",
			formValue:   "111",
			expectedErr: RedirectError("/create-trust-corporation?id=1&caseId=2&trustCorporationId=111&replacement=true"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockCreateLpaClient{}
			client.
				On("Person", mock.Anything, 1).
				Return(sirius.Person{Firstname: "John", Surname: "Doe"}, nil)
			client.
				On("GetUserPermissions", mock.Anything).
				Return(sirius.Permissions{}, nil)
			client.
				On("CreateLpa", mock.Anything, 1, mock.Anything).
				Return(sirius.Lpa{Case: sirius.Case{ID: 2}}, nil)

			template := &mockTemplate{}

			form := url.Values{
				"caseSubtype": {"pfa"},
				tc.formKey:    {tc.formValue},
			}

			r, _ := http.NewRequest(http.MethodPost, "/create-lpa?id=1", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := CreateLpa(client, template.Func)(w, r)
			resp := w.Result()

			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}
