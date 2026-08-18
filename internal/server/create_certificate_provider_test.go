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

type mockCreateCertificateProviderClient struct {
	mock.Mock
}

func (m *mockCreateCertificateProviderClient) CreateCertificateProvider(ctx sirius.Context, caseId int, certificateProvider sirius.Person) error {
	args := m.Called(ctx, caseId, certificateProvider)
	return args.Error(0)
}

func (m *mockCreateCertificateProviderClient) Lpa(ctx sirius.Context, id int) (sirius.Lpa, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Lpa), args.Error(1)
}

func TestGetCreateCertificateProviders(t *testing.T) {
	tests := []struct {
		name                 string
		certificateProviders []sirius.Person
		canAddActor          bool
		isPartial            bool
	}{
		{
			name:                 "Can add another certificate provider",
			certificateProviders: nil,
			canAddActor:          true,
		},
		{
			name: "Cannot add certificate provider",
			certificateProviders: []sirius.Person{
				{ID: 1},
			},
			canAddActor: false,
		},
		{
			name:                 "Can add certificate provider on htmx",
			certificateProviders: nil,
			canAddActor:          true,
			isPartial:            true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lpa := sirius.Lpa{
				Case: sirius.Case{
					ID:                   2,
					CertificateProviders: tc.certificateProviders,
				},
			}

			client := &mockCreateCertificateProviderClient{}
			client.
				On("Lpa", mock.Anything, 2).
				Return(lpa, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, CertificateProviderData{
					DonorId:     1,
					CaseId:      2,
					CanAddActor: tc.canAddActor,
					Title:       "Add a certificate provider",
					PostURL:     "/create-certificate-provider?id=1&caseId=2",
					IsPartial:   tc.isPartial,
				}).
				Return(nil)

			r, _ := http.NewRequest(http.MethodGet, "/create-certificate-provider?id=1&caseId=2", nil)
			if tc.isPartial {
				r.Header.Add("HX-Request", "true")
			}
			w := httptest.NewRecorder()

			err := CreateCertificateProvider(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestGetCreateCertificateProviderLpaFail(t *testing.T) {
	lpa := sirius.Lpa{Case: sirius.Case{ID: 2}}

	client := &mockCreateCertificateProviderClient{}
	client.
		On("Lpa", mock.Anything, 2).
		Return(lpa, errExample)

	template := &mockTemplate{}

	r, _ := http.NewRequest(http.MethodGet, "/create-certificate-provider?id=1&caseId=2", nil)
	w := httptest.NewRecorder()

	err := CreateCertificateProvider(client, template.Func)(w, r)
	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetCreateCertificateProviderBadQuery(t *testing.T) {
	testCases := map[string]string{
		"no-id":       "/",
		"bad-id":      "/?id=test",
		"bad-case-id": "/?id=123&caseId=test",
	}

	for name, query := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, query, nil)
			w := httptest.NewRecorder()

			err := CreateCertificateProvider(nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestPostCreateCertificateProvider(t *testing.T) {
	tests := []struct {
		name        string
		addActor    string
		redirectURL string
		htmxRequest bool
		swap        string
	}{
		{
			name:        "Submit",
			addActor:    "",
			redirectURL: "/create-lpa?id=1&caseId=2",
			htmxRequest: false,
		},
		{
			name:        "Submit htmx request",
			addActor:    "",
			redirectURL: "/create-lpa?id=1&caseId=2",
			htmxRequest: true,
			swap:        "innerHTML show:#accordion-create-lpa-heading-3:top",
		},
		{
			name:        "Submit and add another certificate provider",
			addActor:    "true",
			redirectURL: "/create-certificate-provider?id=1&caseId=2",
			htmxRequest: false,
		},
		{
			name:        "Submit and add another certificate provider htmx request",
			addActor:    "true",
			redirectURL: "/create-certificate-provider?id=1&caseId=2",
			htmxRequest: true,
			swap:        "innerHTML scroll:.action-panel__content:top",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockCreateCertificateProviderClient{}
			client.
				On("Lpa", mock.Anything, 2).
				Return(sirius.Lpa{Case: sirius.Case{ID: 2}}, nil)
			client.
				On("CreateCertificateProvider", mock.Anything, 2, sirius.Person{
					Salutation:   "Sir",
					Firstname:    "Arthur",
					Middlenames:  "Conan",
					Surname:      "Doyle",
					AddressLine1: "221B",
					AddressLine2: "Baker Street",
					AddressLine3: "Marylebone",
					Town:         "London",
					Postcode:     "NW1 6XE",
					County:       "Greater London",
					Country:      "United Kingdom",
				}).
				Return(nil)

			template := &mockTemplate{}
			partialTemplate := &mockTemplate{}
			if tc.htmxRequest {
				partialTemplate.
					On("Func", mock.Anything, CertificateProviderData{
						DonorId:      1,
						CaseId:       2,
						CanAddActor:  true,
						HtmxRedirect: tc.redirectURL,
						HtmxSwap:     tc.swap,
						Title:        "Add a certificate provider",
						PostURL:      "/create-certificate-provider?id=1&caseId=2",
					}).
					Return(nil)
			}

			form := url.Values{
				"salutation":   {"Sir"},
				"firstname":    {"Arthur"},
				"middlenames":  {"Conan"},
				"surname":      {"Doyle"},
				"addressLine1": {"221B"},
				"addressLine2": {"Baker Street"},
				"addressLine3": {"Marylebone"},
				"town":         {"London"},
				"county":       {"Greater London"},
				"postcode":     {"NW1 6XE"},
				"country":      {"United Kingdom"},
				"add-another":  {tc.addActor},
			}

			r, _ := http.NewRequest(http.MethodPost, "/create-certificate-provider?id=1&caseId=2", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			if tc.htmxRequest {
				r.Header.Add("HX-Request", "true")
			}
			w := httptest.NewRecorder()

			err := CreateCertificateProvider(client, template.Func)(w, r)
			resp := w.Result()
			assert.Equal(t, RedirectError(tc.redirectURL), err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostCreateCertificateProviderWhenAPIFails(t *testing.T) {
	client := &mockCreateCertificateProviderClient{}
	client.
		On("Lpa", mock.Anything, 2).
		Return(sirius.Lpa{Case: sirius.Case{ID: 2}}, nil)
	client.
		On("CreateCertificateProvider", mock.Anything, 2, sirius.Person{
			Salutation:   "Sir",
			Firstname:    "Arthur",
			Middlenames:  "Conan",
			Surname:      "Doyle",
			AddressLine1: "221B",
			AddressLine2: "Baker Street",
			AddressLine3: "Marylebone",
			Town:         "London",
			Postcode:     "NW1 6XE",
			County:       "Greater London",
			Country:      "United Kingdom",
		}).
		Return(errExample)

	template := &mockTemplate{}

	form := url.Values{
		"salutation":   {"Sir"},
		"firstname":    {"Arthur"},
		"middlenames":  {"Conan"},
		"surname":      {"Doyle"},
		"addressLine1": {"221B"},
		"addressLine2": {"Baker Street"},
		"addressLine3": {"Marylebone"},
		"town":         {"London"},
		"county":       {"Greater London"},
		"postcode":     {"NW1 6XE"},
		"country":      {"United Kingdom"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/create-certificate-provider?id=1&caseId=2", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateCertificateProvider(client, template.Func)(w, r)
	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateCertificateProviderValidationError(t *testing.T) {
	client := &mockCreateCertificateProviderClient{}
	client.
		On("Lpa", mock.Anything, 2).
		Return(sirius.Lpa{Case: sirius.Case{ID: 2}}, nil)
	client.
		On("CreateCertificateProvider", mock.Anything, 2, sirius.Person{
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
		On("Func", mock.Anything, CertificateProviderData{
			CanAddActor: true,
			CaseId:      2,
			DonorId:     1,
			Error: sirius.ValidationError{
				Field: sirius.FieldErrors{
					"firstname": {"required": "This field is required"},
				},
			},
			Title:   "Add a certificate provider",
			PostURL: "/create-certificate-provider?id=1&caseId=2",
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

	r, _ := http.NewRequest(http.MethodPost, "/create-certificate-provider?id=1&caseId=2", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateCertificateProvider(client, template.Func)(w, r)
	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}
