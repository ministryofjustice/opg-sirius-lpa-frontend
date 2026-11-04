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
	flowCases := []struct {
		name      string
		flowQuery string
	}{
		{name: "Without flow"},
		{name: "Create flow", flowQuery: "&flow=create"},
	}

	for _, flowCase := range flowCases {
		t.Run(flowCase.name, func(t *testing.T) {
			for _, isPartial := range []bool{true, false} {
				t.Run("Is partial: "+strconv.FormatBool(isPartial), func(t *testing.T) {
					template := &mockTemplate{}
					template.
						On("Func", mock.Anything, CertificateProviderData{
							FlowQuery: flowCase.flowQuery,
							DonorId:   1,
							CaseId:    2,
							Title:     "Add a certificate provider",
							PostURL:   "/create-certificate-provider?id=1&caseId=2" + flowCase.flowQuery,
							IsPartial: isPartial,
						}).
						Return(nil)

					r, _ := http.NewRequest(http.MethodGet, "/create-certificate-provider?id=1&caseId=2"+flowCase.flowQuery, nil)
					if isPartial {
						r.Header.Add("HX-Request", "true")
					}
					w := httptest.NewRecorder()

					err := CreateCertificateProvider(nil, template.Func)(w, r)
					resp := w.Result()

					assert.Nil(t, err)
					assert.Equal(t, http.StatusOK, resp.StatusCode)
					mock.AssertExpectationsForObjects(t, template)
				})
			}
		})
	}
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
	flowCases := []struct {
		name      string
		flowQuery string
	}{
		{name: "Without flow"},
		{name: "Create flow", flowQuery: "&flow=create"},
	}

	for _, flowCase := range flowCases {
		t.Run(flowCase.name, func(t *testing.T) {
			for _, isHtmx := range []bool{false, true} {
				t.Run("Is Htmx: "+strconv.FormatBool(isHtmx), func(t *testing.T) {
					client := &mockCreateCertificateProviderClient{}
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
					if isHtmx {
						template.
							On("Func", mock.Anything, CertificateProviderData{
								FlowQuery:    flowCase.flowQuery,
								DonorId:      1,
								CaseId:       2,
								HtmxRedirect: "/create-lpa?id=1&caseId=2" + flowCase.flowQuery + "#accordion-create-lpa-heading-3",
								HtmxSwap:     "innerHTML show:#accordion-create-lpa-heading-3:top",
								Title:        "Add a certificate provider",
								PostURL:      "/create-certificate-provider?id=1&caseId=2" + flowCase.flowQuery,
								IsPartial:    true,
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
					}

					r, _ := http.NewRequest(http.MethodPost, "/create-certificate-provider?id=1&caseId=2"+flowCase.flowQuery, strings.NewReader(form.Encode()))
					r.Header.Add("Content-Type", formUrlEncoded)
					if isHtmx {
						r.Header.Add("HX-Request", "true")
					}
					w := httptest.NewRecorder()

					err := CreateCertificateProvider(client, template.Func)(w, r)
					resp := w.Result()
					if isHtmx {
						assert.Nil(t, err)
					} else {
						assert.Equal(t, RedirectError("/create-lpa?id=1&caseId=2"+flowCase.flowQuery+"#accordion-create-lpa-heading-3"), err)
					}
					assert.Equal(t, http.StatusOK, resp.StatusCode)
					mock.AssertExpectationsForObjects(t, client, template)
				})
			}
		})
	}
}

func TestPostCreateCertificateProviderWhenAPIFails(t *testing.T) {
	client := &mockCreateCertificateProviderClient{}
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
			CaseId:  2,
			DonorId: 1,
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
