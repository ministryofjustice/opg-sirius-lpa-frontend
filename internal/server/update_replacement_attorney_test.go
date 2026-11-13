package server

import (
	"fmt"
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

type mockUpdateReplacementAttorneyClient struct {
	mock.Mock
}

func (m *mockUpdateReplacementAttorneyClient) Lpa(ctx sirius.Context, id int) (sirius.Lpa, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Lpa), args.Error(1)
}

func (m *mockUpdateReplacementAttorneyClient) UpdateReplacementAttorney(ctx sirius.Context, attorneyId int, attorney sirius.Attorney) error {
	args := m.Called(ctx, attorneyId, attorney)
	return args.Error(0)
}

func TestGetUpdateReplacementAttorney(t *testing.T) {
	existingAttorney := sirius.Attorney{
		Person: sirius.Person{
			ID:        4,
			Firstname: "Rudolph",
			Surname:   "Stotesbury",
		},
	}

	client := &mockUpdateReplacementAttorneyClient{}
	client.
		On("Lpa", mock.Anything, 2).
		Return(sirius.Lpa{
			Case: sirius.Case{
				ReceiptDate:          sirius.DateString("2026-08-01"),
				ReplacementAttorneys: []sirius.Attorney{existingAttorney},
			},
		}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, replacementAttorneyData{
			DonorId:   1,
			CaseId:    2,
			Attorney:  existingAttorney,
			IsEditing: true,
			Title:     "Update replacement attorney details",
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=1&caseId=2&attorneyId=4", nil)
	w := httptest.NewRecorder()

	err := UpdateReplacementAttorney(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetUpdateReplacementAttorneyLPAFail(t *testing.T) {
	client := &mockUpdateReplacementAttorneyClient{}
	client.
		On("Lpa", mock.Anything, 2).
		Return(sirius.Lpa{}, errExample)

	template := &mockTemplate{}

	r, _ := http.NewRequest(http.MethodGet, "/update-replacement-attorney?id=1&caseId=2&attorneyId=3", nil)
	w := httptest.NewRecorder()

	err := UpdateReplacementAttorney(client, template.Func)(w, r)
	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostUpdateReplacementAttorneyTransitions(t *testing.T) {
	tests := []struct {
		name           string
		formKey        string
		redirect       string
		fragment       string
		htmxSwap       string
		nextAttorneyID int
		lpa            sirius.Lpa
	}{
		{
			name:     "Save and return to LPA",
			redirect: "/create-lpa?id=1&caseId=2",
			fragment: "#scroll-to-replacement-attorneys",
			htmxSwap: "innerHTML show:#scroll-to-replacement-attorneys:top",
			lpa: sirius.Lpa{Case: sirius.Case{
				ReceiptDate:          sirius.DateString("2026-08-01"),
				ReplacementAttorneys: []sirius.Attorney{{Person: sirius.Person{ID: 4}}},
			}},
		},
		{
			name:           "Update next attorney",
			formKey:        "update-next-attorney",
			redirect:       "/update-replacement-attorney?id=1&caseId=2&attorneyId=5",
			htmxSwap:       "innerHTML scroll:.action-panel__content:top",
			nextAttorneyID: 5,
			lpa: sirius.Lpa{
				Case: sirius.Case{
					ReceiptDate: sirius.DateString("2026-08-01"),
					ReplacementAttorneys: []sirius.Attorney{
						{Person: sirius.Person{ID: 4}},
						{Person: sirius.Person{ID: 5, PersonType: "Replacement Attorney"}, SystemStatus: shared.BoolPtr(false)},
					},
				},
			},
		},
		{
			name:           "Update next trust corporation",
			formKey:        "update-next-attorney",
			redirect:       "/update-trust-corporation?id=1&caseId=2&trustCorporationId=5&replacement=true",
			htmxSwap:       "innerHTML scroll:.action-panel__content:top",
			nextAttorneyID: 5,
			lpa: sirius.Lpa{
				Case: sirius.Case{
					ReceiptDate:          sirius.DateString("2026-08-01"),
					ReplacementAttorneys: []sirius.Attorney{{Person: sirius.Person{ID: 4}}},
					TrustCorporations: []sirius.TrustCorporation{
						{Attorney: sirius.Attorney{Person: sirius.Person{ID: 5, PersonType: "Trust Corporation"}}, IsReplacementAttorney: true},
					},
				},
			},
		},
	}

	flows := []struct{ name, query string }{
		{"no flow", ""},
		{"create flow", "&flow=create"},
	}

	for _, tc := range tests {
		for _, flow := range flows {
			for _, isHtmx := range []bool{false, true} {
				name := fmt.Sprintf("%s, %s, htmx=%t", tc.name, flow.name, isHtmx)

				t.Run(name, func(t *testing.T) {
					updatedAttorney := sirius.Attorney{
						Person: sirius.Person{
							Firstname: "Rudolph",
							Surname:   "Stotesbury",
						},
					}

					client := &mockUpdateReplacementAttorneyClient{}
					client.
						On("Lpa", mock.Anything, 2).
						Return(tc.lpa, nil).
						On("UpdateReplacementAttorney", mock.Anything, 4, updatedAttorney).
						Return(nil)

					template := &mockTemplate{}

					if isHtmx {
						template.
							On("Func", mock.Anything, replacementAttorneyData{
								FlowQuery:      flow.query,
								DonorId:        1,
								CaseId:         2,
								Attorney:       updatedAttorney,
								IsEditing:      true,
								Title:          "Update replacement attorney details",
								HtmxRedirect:   tc.redirect + flow.query,
								HtmxSwap:       tc.htmxSwap,
								IsPartial:      true,
								NextAttorneyId: tc.nextAttorneyID,
							}).
							Return(nil)
					}

					form := url.Values{
						"firstname": {"Rudolph"},
						"surname":   {"Stotesbury"},
					}
					if tc.formKey != "" {
						form.Set(tc.formKey, "true")
					}

					r, _ := http.NewRequest(http.MethodPost, "/update-replacement-attorney?id=1&caseId=2&attorneyId=4"+flow.query, strings.NewReader(form.Encode()))
					r.Header.Add("Content-Type", formUrlEncoded)
					if isHtmx {
						r.Header.Add("HX-Request", "true")
					}
					w := httptest.NewRecorder()

					err := UpdateReplacementAttorney(client, template.Func)(w, r)
					resp := w.Result()

					if !isHtmx {
						expectedError := RedirectError(tc.redirect + flow.query + tc.fragment)
						assert.Equal(t, err, expectedError)
					}
					assert.Equal(t, http.StatusOK, resp.StatusCode)
					mock.AssertExpectationsForObjects(t, client, template)
				})
			}
		}
	}
}

func TestPostUpdateReplacementAttorneyBadQuery(t *testing.T) {
	testCases := map[string]string{
		"no-id":           "/",
		"bad-id":          "/?id=test",
		"bad-case-id":     "/?id=123&caseId=test",
		"bad-attorney-id": "/?id=123&caseId=123&attorneyId=test",
	}

	for name, query := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, query, nil)
			w := httptest.NewRecorder()

			err := UpdateReplacementAttorney(nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestPostUpdateReplacementAttorneyValidationError(t *testing.T) {
	lpa := sirius.Lpa{Case: sirius.Case{
		ID:                   3,
		ReplacementAttorneys: []sirius.Attorney{{Person: sirius.Person{ID: 4}}},
	}}

	client := &mockUpdateReplacementAttorneyClient{}

	client.
		On("Lpa", mock.Anything, 3).
		Return(lpa, nil)
	client.
		On("UpdateReplacementAttorney", mock.Anything, 4, sirius.Attorney{
			Person: sirius.Person{
				Salutation:   "Sir",
				Surname:      "Doyle",
				AddressLine1: "221B",
				Town:         "London",
				Postcode:     "NW1 6XE",
				Country:      "United Kingdom",
			},
		}).
		Return(sirius.ValidationError{Field: sirius.FieldErrors{
			"firstname": {"required": "This field is required"},
		}})

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, replacementAttorneyData{
			CaseId:  3,
			DonorId: 1,
			Attorney: sirius.Attorney{
				Person: sirius.Person{
					Salutation:   "Sir",
					Surname:      "Doyle",
					AddressLine1: "221B",
					Town:         "London",
					Postcode:     "NW1 6XE",
					Country:      "United Kingdom",
				},
			},
			Error: sirius.ValidationError{
				Field: sirius.FieldErrors{
					"firstname": {"required": "This field is required"},
				},
			},
			Title:     "Update replacement attorney details",
			IsEditing: true,
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

	r, _ := http.NewRequest(http.MethodPost, "/update-replacement-attorney?id=1&caseId=3&attorneyId=4", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := UpdateReplacementAttorney(client, template.Func)(w, r)
	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostUpdateReplacementAttorneyWhenUpdateErrors(t *testing.T) {
	lpa := sirius.Lpa{Case: sirius.Case{
		ID:                   3,
		ReplacementAttorneys: []sirius.Attorney{{Person: sirius.Person{ID: 4}}},
	}}

	client := &mockUpdateReplacementAttorneyClient{}
	client.
		On("Lpa", mock.Anything, 3).
		Return(lpa, nil)
	client.
		On("UpdateReplacementAttorney", mock.Anything, 4, sirius.Attorney{
			Person: sirius.Person{Surname: "Doyle"},
		}).
		Return(errExample)

	template := &mockTemplate{}

	form := url.Values{"surname": {"Doyle"}}

	r, _ := http.NewRequest(http.MethodPost, "/update-replacement-attorney?id=1&caseId=3&attorneyId=4", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := UpdateReplacementAttorney(client, template.Func)(w, r)

	assert.Equal(t, errExample, err)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}
