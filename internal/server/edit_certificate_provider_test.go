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

type mockEditCertificateProviderClient struct {
	mock.Mock
}

func (m *mockEditCertificateProviderClient) UpdateCertificateProvider(ctx sirius.Context, certificateProviderId int, certificateProvider sirius.Person) error {
	args := m.Called(ctx, certificateProviderId, certificateProvider)
	return args.Error(0)
}

func (m *mockEditCertificateProviderClient) Person(ctx sirius.Context, id int) (sirius.Person, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Person), args.Error(1)
}

var mockCertificateProvider = sirius.Person{
	ID:           3,
	CaseId:       2,
	Salutation:   "Sir",
	Firstname:    "Arthur",
	Middlenames:  "Conan",
	Surname:      "Doyle",
	AddressLine1: "221B",
	AddressLine2: "Baker Street",
	AddressLine3: "Marylebone",
	Town:         "London",
	County:       "Greater London",
	Postcode:     "NW1 6XE",
	Country:      "United Kingdom",
	PersonType:   "Certificate Provider",
}

func TestGetEditCertificateProvidersTest(t *testing.T) {
	tests := []struct {
		name          string
		isHtmxRequest bool
	}{
		{
			name:          "HTMX Request",
			isHtmxRequest: true,
		},
		{
			name:          "Non HTMX Request",
			isHtmxRequest: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockEditCertificateProviderClient{}
			client.
				On("Person", mock.Anything, 3).
				Return(mockCertificateProvider, nil)

			template := &mockTemplate{}
			partialTemplate := &mockTemplate{}
			expectedData := EditCertificateProviderData{
				DonorId:             1,
				CaseId:              2,
				CanAddActor:         false,
				CertificateProvider: mockCertificateProvider,
				Title:               "Edit a certificate provider",
			}

			if tc.isHtmxRequest {
				partialTemplate.
					On("Func", mock.Anything, expectedData).
					Return(nil)
			} else {
				template.
					On("Func", mock.Anything, expectedData).
					Return(nil)
			}

			r, _ := http.NewRequest(http.MethodGet, "/edit-certificate-provider/?id=1&caseId=2&personId=3", nil)
			if tc.isHtmxRequest {
				r.Header.Add("HX-Request", "true")
			}
			w := httptest.NewRecorder()

			err := EditCertificateProvider(client, template.Func, partialTemplate.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template, partialTemplate)
			if tc.isHtmxRequest {
				template.AssertNotCalled(t, "Func")
			} else {
				partialTemplate.AssertNotCalled(t, "Func")
			}
		})
	}
}

func TestGetEditCertificateProviders(t *testing.T) {
	client := &mockEditCertificateProviderClient{}
	client.
		On("Person", mock.Anything, 3).
		Return(mockCertificateProvider, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, EditCertificateProviderData{
			DonorId:             1,
			CaseId:              2,
			CanAddActor:         false,
			CertificateProvider: mockCertificateProvider,
			Title:               "Edit a certificate provider",
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/edit-certificate-provider/?id=1&caseId=2&personId=3", nil)
	w := httptest.NewRecorder()

	err := EditCertificateProvider(client, template.Func, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)

}

func TestGetEditCertificateProviderPersonFail(t *testing.T) {
	client := &mockEditCertificateProviderClient{}
	client.
		On("Person", mock.Anything, 3).
		Return(mockCertificateProvider, errExample)

	template := &mockTemplate{}

	r, _ := http.NewRequest(http.MethodGet, "/edit-certificate-provider/?id=1&caseId=2&personId=3", nil)
	w := httptest.NewRecorder()

	err := EditCertificateProvider(client, template.Func, template.Func)(w, r)
	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetEditCertificateProviderBadQuery(t *testing.T) {
	testCases := map[string]string{
		"no-id":         "/",
		"bad-id":        "/?id=test",
		"bad-case-id":   "/?id=123&caseId=test",
		"bad-person-id": "/?id=123&caseId=123&personId=test",
	}

	for name, query := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, query, nil)
			w := httptest.NewRecorder()

			err := EditCertificateProvider(nil, nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestPostEditCertificateProvider(t *testing.T) {
	tests := []struct {
		name        string
		htmxRequest bool
		swap        string
		error       error
		expectedErr error
	}{
		{
			name:        "Submit",
			htmxRequest: false,
			error:       nil,
			expectedErr: RedirectError("/create-lpa?id=1&caseId=2"),
		},
		{
			name:        "Submit API Failure",
			htmxRequest: false,
			error:       errExample,
			expectedErr: errExample,
		},
		{
			name:        "Submit htmx request",
			htmxRequest: true,
			swap:        "innerHTML show:#accordion-create-lpa-heading-3:top",
			error:       nil,
			expectedErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockEditCertificateProviderClient{}
			client.
				On("Person", mock.Anything, 3).
				Return(mockCertificateProvider, nil)
			client.
				On("UpdateCertificateProvider", mock.Anything, 3, sirius.Person{
					Salutation:   "Dr",
					Firstname:    "John",
					Middlenames:  "",
					Surname:      "Watson",
					AddressLine1: "221B",
					AddressLine2: "Baker Street",
					AddressLine3: "Marylebone",
					Town:         "London",
					Postcode:     "NW1 6XE",
					County:       "Greater London",
					Country:      "United Kingdom",
				}).
				Return(tc.error)

			template := &mockTemplate{}
			partialTemplate := &mockTemplate{}
			if tc.htmxRequest {
				partialTemplate.
					On("Func", mock.Anything, EditCertificateProviderData{
						DonorId:             1,
						CaseId:              2,
						CanAddActor:         false,
						CertificateProvider: mockCertificateProvider,
						HtmxRedirect:        "/create-lpa?id=1&caseId=2",
						HtmxSwap:            tc.swap,
						Title:               "Edit a certificate provider",
					}).
					Return(nil)
			}

			form := url.Values{
				"salutation":   {"Dr"},
				"firstname":    {"John"},
				"middlenames":  {""},
				"surname":      {"Watson"},
				"addressLine1": {"221B"},
				"addressLine2": {"Baker Street"},
				"addressLine3": {"Marylebone"},
				"town":         {"London"},
				"county":       {"Greater London"},
				"postcode":     {"NW1 6XE"},
				"country":      {"United Kingdom"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/edit-certificate-provider/?id=1&caseId=2&personId=3", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			if tc.htmxRequest {
				r.Header.Add("HX-Request", "true")
			}
			w := httptest.NewRecorder()

			err := EditCertificateProvider(client, template.Func, partialTemplate.Func)(w, r)
			resp := w.Result()

			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template, partialTemplate)
		})
	}
}

func TestPostEditCertificateProviderValidationError(t *testing.T) {
	client := &mockEditCertificateProviderClient{}
	client.
		On("Person", mock.Anything, 3).
		Return(mockCertificateProvider, nil)
	client.
		On("UpdateCertificateProvider", mock.Anything, 3, sirius.Person{
			Salutation:   "Sir",
			Surname:      "Doyle",
			AddressLine1: "221B",
			Town:         "London",
			Postcode:     "NW1 6XE",
			Country:      "United Kingdom",
		}).
		Return(sirius.ValidationError{Field: sirius.FieldErrors{
			"firstname": {"required": "This field is required"},
		}})

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, EditCertificateProviderData{
			CanAddActor:         false,
			CaseId:              2,
			CertificateProvider: mockCertificateProvider,
			DonorId:             1,
			Error: sirius.ValidationError{
				Field: sirius.FieldErrors{
					"firstname": {"required": "This field is required"},
				},
			},
			Title: "Edit a certificate provider",
		}).
		Return(nil)

	form := url.Values{
		"salutation":   {"Sir"},
		"surname":      {"Doyle"},
		"addressLine1": {"221B"},
		"town":         {"London"},
		"postcode":     {"NW1 6XE"},
		"country":      {"United Kingdom"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/edit-certificate-provider/?id=1&caseId=2&personId=3", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := EditCertificateProvider(client, template.Func, template.Func)(w, r)
	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}
