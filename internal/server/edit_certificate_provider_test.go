package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
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

func TestGetEditCertificateProviders(t *testing.T) {
	flowCases := []struct {
		name      string
		flowQuery string
	}{
		{name: "Without flow"},
		{name: "Create flow", flowQuery: "&flow=create"},
	}
	for _, flowCase := range flowCases {
		t.Run(flowCase.name, func(t *testing.T) {
			for _, isHtmxRequest := range []bool{false, true} {
				t.Run("Is Htmx: "+strconv.FormatBool(isHtmxRequest), func(t *testing.T) {
					client := &mockEditCertificateProviderClient{}
					client.
						On("Person", mock.Anything, 3).
						Return(mockCertificateProvider, nil)

					template := &mockTemplate{}
					template.
						On("Func", mock.Anything, CertificateProviderData{
							FlowQuery:           flowCase.flowQuery,
							DonorId:             1,
							CaseId:              2,
							CanAddActor:         false,
							CertificateProvider: mockCertificateProvider,
							Title:               "Edit a certificate provider",
							PostURL:             "/edit-certificate-provider?id=1&caseId=2&personId=3" + flowCase.flowQuery,
							IsPartial:           isHtmxRequest,
						}).
						Return(nil)

					r, _ := http.NewRequest(http.MethodGet, "/edit-certificate-provider/?id=1&caseId=2&personId=3"+flowCase.flowQuery, nil)
					if isHtmxRequest {
						r.Header.Add("HX-Request", "true")
					}
					w := httptest.NewRecorder()

					err := EditCertificateProvider(client, template.Func)(w, r)
					resp := w.Result()

					assert.Nil(t, err)
					assert.Equal(t, http.StatusOK, resp.StatusCode)
					mock.AssertExpectationsForObjects(t, client, template)
				})
			}
		})
	}
}

func TestGetEditCertificateProviderPersonFail(t *testing.T) {
	client := &mockEditCertificateProviderClient{}
	client.
		On("Person", mock.Anything, 3).
		Return(mockCertificateProvider, errExample)

	template := &mockTemplate{}

	r, _ := http.NewRequest(http.MethodGet, "/edit-certificate-provider?id=1&caseId=2&personId=3", nil)
	w := httptest.NewRecorder()

	err := EditCertificateProvider(client, template.Func)(w, r)
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

			err := EditCertificateProvider(nil, nil)(w, r)

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
		redirect    string
		fragment    string
	}{
		{
			name:        "Submit",
			htmxRequest: false,
			redirect:    "/create-lpa?id=1&caseId=2",
			fragment:    "#accordion-create-lpa-heading-3",
		},
		{
			name:        "Submit API Failure",
			htmxRequest: false,
			error:       errExample,
		},
		{
			name:        "Submit htmx request",
			htmxRequest: true,
			swap:        "innerHTML show:#accordion-create-lpa-heading-3:top",
			redirect:    "/create-lpa?id=1&caseId=2",
			fragment:    "#accordion-create-lpa-heading-3",
		},
	}
	flowCases := []struct {
		name      string
		flowQuery string
	}{
		{name: "Without flow"},
		{name: "Create flow", flowQuery: "&flow=create"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for _, flowCase := range flowCases {
				t.Run(flowCase.name, func(t *testing.T) {
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
					if tc.htmxRequest {
						template.
							On("Func", mock.Anything, CertificateProviderData{
								FlowQuery:           flowCase.flowQuery,
								DonorId:             1,
								CaseId:              2,
								CanAddActor:         false,
								CertificateProvider: mockCertificateProvider,
								HtmxRedirect:        tc.redirect + flowCase.flowQuery + tc.fragment,
								HtmxSwap:            tc.swap,
								Title:               "Edit a certificate provider",
								PostURL:             "/edit-certificate-provider?id=1&caseId=2&personId=3" + flowCase.flowQuery,
								IsPartial:           true,
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

					r, _ := http.NewRequest(http.MethodPost, "/edit-certificate-provider?id=1&caseId=2&personId=3"+flowCase.flowQuery, strings.NewReader(form.Encode()))
					r.Header.Add("Content-Type", formUrlEncoded)
					if tc.htmxRequest {
						r.Header.Add("HX-Request", "true")
					}
					w := httptest.NewRecorder()

					err := EditCertificateProvider(client, template.Func)(w, r)
					resp := w.Result()

					expectedErr := tc.error
					if expectedErr == nil && !tc.htmxRequest {
						expectedErr = RedirectError(tc.redirect + flowCase.flowQuery + tc.fragment)
					}
					assert.Equal(t, expectedErr, err)
					assert.Equal(t, http.StatusOK, resp.StatusCode)
					mock.AssertExpectationsForObjects(t, client, template)
				})
			}
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
		On("Func", mock.Anything, CertificateProviderData{
			CanAddActor:         false,
			CaseId:              2,
			CertificateProvider: mockCertificateProvider,
			DonorId:             1,
			Error: sirius.ValidationError{
				Field: sirius.FieldErrors{
					"firstname": {"required": "This field is required"},
				},
			},
			Title:   "Edit a certificate provider",
			PostURL: "/edit-certificate-provider?id=1&caseId=2&personId=3",
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

	r, _ := http.NewRequest(http.MethodPost, "/edit-certificate-provider?id=1&caseId=2&personId=3", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := EditCertificateProvider(client, template.Func)(w, r)
	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}
