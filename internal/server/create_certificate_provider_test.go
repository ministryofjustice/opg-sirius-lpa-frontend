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

	flowCases := []struct {
		name      string
		flowQuery string
	}{
		{name: "Without flow"},
		{name: "Create flow", flowQuery: "&flow=create"},
	}

	for _, flowCase := range flowCases {
		t.Run(flowCase.name, func(t *testing.T) {
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
							FlowQuery:   flowCase.flowQuery,
							DonorId:     1,
							CaseId:      2,
							CanAddActor: tc.canAddActor,
							Title:       "Add a certificate provider",
							PostURL:     "/create-certificate-provider?id=1&caseId=2" + flowCase.flowQuery,
							IsPartial:   tc.isPartial,
						}).
						Return(nil)

					r, _ := http.NewRequest(http.MethodGet, "/create-certificate-provider?id=1&caseId=2"+flowCase.flowQuery, nil)
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
	transitions := []struct {
		name       string
		addAnother bool
		redirect   string
		fragment   string
		htmxSwap   string
	}{
		{
			name:     "Save and return to LPA",
			redirect: "/create-lpa?id=1&caseId=2",
			fragment: "#accordion-create-lpa-heading-3",
			htmxSwap: "innerHTML show:#accordion-create-lpa-heading-3:top",
		},
		{
			name:       "Save and add another",
			addAnother: true,
			redirect:   "/create-certificate-provider?id=1&caseId=2",
			htmxSwap:   "innerHTML scroll:.action-panel__content:top",
		},
	}
	flowCases := []struct {
		name      string
		flowQuery string
	}{
		{name: "Without flow"},
		{name: "Create flow", flowQuery: "&flow=create"},
	}

	for _, transition := range transitions {
		t.Run(transition.name, func(t *testing.T) {
			for _, flowCase := range flowCases {
				t.Run(flowCase.name, func(t *testing.T) {
					for _, isHtmx := range []bool{false, true} {
						t.Run("Is Htmx: "+strconv.FormatBool(isHtmx), func(t *testing.T) {
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
							if isHtmx {
								template.
									On("Func", mock.Anything, CertificateProviderData{
										FlowQuery:    flowCase.flowQuery,
										DonorId:      1,
										CaseId:       2,
										CanAddActor:  true,
										HtmxRedirect: transition.redirect + flowCase.flowQuery + transition.fragment,
										HtmxSwap:     transition.htmxSwap,
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
							if transition.addAnother {
								form.Set("add-another", "true")
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
								assert.Equal(t, RedirectError(transition.redirect+flowCase.flowQuery+transition.fragment), err)
							}
							assert.Equal(t, http.StatusOK, resp.StatusCode)
							mock.AssertExpectationsForObjects(t, client, template)
						})
					}
				})
			}
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
