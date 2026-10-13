package server

import (
	"fmt"
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

type mockCreateTrustCorporationClient struct {
	mock.Mock
}

func (m *mockCreateTrustCorporationClient) Lpa(ctx sirius.Context, id int) (sirius.Lpa, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Lpa), args.Error(1)
}

func (m *mockCreateTrustCorporationClient) CreateTrustCorporation(ctx sirius.Context, caseId int, trustCorporation sirius.TrustCorporation) error {
	args := m.Called(ctx, caseId, trustCorporation)
	return args.Error(0)
}

func (m *mockCreateTrustCorporationClient) UpdateTrustCorporation(ctx sirius.Context, trustCorporationId int, trustCorporation sirius.TrustCorporation) error {
	args := m.Called(ctx, trustCorporationId, trustCorporation)
	return args.Error(0)
}

func TestGetCreateTrustCorporation(t *testing.T) {
	flowCases := []struct {
		name      string
		flowQuery string
	}{
		{name: "Without flow"},
		{name: "Create flow", flowQuery: "&flow=create"},
	}
	for _, flowCase := range flowCases {
		t.Run(flowCase.name, func(t *testing.T) {
			for _, isReplacementAttorney := range []string{"false", "true"} {
				t.Run("Is Replacement Attorney: "+isReplacementAttorney, func(t *testing.T) {
					expectedData := createTrustCorporationData{
						FlowQuery: flowCase.flowQuery,
						DonorId:   1,
						CaseId:    2,
						TrustCorporation: sirius.TrustCorporation{
							IsReplacementAttorney: isReplacementAttorney == "true",
							Attorney:              sirius.Attorney{SystemStatus: shared.BoolPtr(true)},
						},
						Title:    "Add a trust corporation",
						HtmxPost: "/create-trust-corporation?id=1&caseId=2&replacement=" + isReplacementAttorney + flowCase.flowQuery,
					}

					if isReplacementAttorney == "true" {
						expectedData.AppointedAs = "Replacement attorney"
					} else {
						expectedData.AppointedAs = "Attorney"
					}

					template := &mockTemplate{}
					template.
						On("Func", mock.Anything, expectedData).
						Return(nil)

					r, _ := http.NewRequest(http.MethodGet, "/?id=1&caseId=2&replacement="+isReplacementAttorney+flowCase.flowQuery, nil)
					w := httptest.NewRecorder()

					err := CreateTrustCorporation(nil, template.Func)(w, r)
					resp := w.Result()

					assert.Nil(t, err)
					assert.Equal(t, http.StatusOK, resp.StatusCode)
					mock.AssertExpectationsForObjects(t, template)
				})
			}
		})
	}
}

func TestGetEditTrustCorporation(t *testing.T) {
	flowCases := []struct {
		name      string
		flowQuery string
	}{
		{name: "Without flow"},
		{name: "Create flow", flowQuery: "&flow=create"},
	}
	for _, flowCase := range flowCases {
		t.Run(flowCase.name, func(t *testing.T) {
			for _, isReplacementAttorney := range []string{"false", "true"} {
				t.Run("Is Replacement Attorney: "+isReplacementAttorney, func(t *testing.T) {
					expectedTrustCorporation := sirius.TrustCorporation{
						Attorney:              sirius.Attorney{Person: sirius.Person{ID: 3}},
						IsReplacementAttorney: isReplacementAttorney == "true",
					}
					expectedData := createTrustCorporationData{
						FlowQuery:        flowCase.flowQuery,
						DonorId:          1,
						CaseId:           2,
						TrustCorporation: expectedTrustCorporation,
						Title:            "Update trust corporation details",
						HtmxPost:         "/create-trust-corporation?id=1&caseId=2&trustCorporationId=3&replacement=" + isReplacementAttorney + flowCase.flowQuery,
						IsEditing:        true,
					}

					if isReplacementAttorney == "true" {
						expectedData.AppointedAs = "Replacement attorney"
					} else {
						expectedData.AppointedAs = "Attorney"
					}

					client := &mockCreateTrustCorporationClient{}
					client.
						On("Lpa", mock.Anything, 2).
						Return(sirius.Lpa{Case: sirius.Case{
							TrustCorporations: []sirius.TrustCorporation{expectedTrustCorporation},
						}}, nil)

					template := &mockTemplate{}
					template.
						On("Func", mock.Anything, expectedData).
						Return(nil)

					r, _ := http.NewRequest(http.MethodGet, "/?id=1&caseId=2&trustCorporationId=3&replacement="+isReplacementAttorney+flowCase.flowQuery, nil)
					w := httptest.NewRecorder()

					err := CreateTrustCorporation(client, template.Func)(w, r)
					resp := w.Result()

					assert.Nil(t, err)
					assert.Equal(t, http.StatusOK, resp.StatusCode)
					mock.AssertExpectationsForObjects(t, client, template)
				})
			}
		})
	}
}

func TestGetCreateTrustCorporationBadQuery(t *testing.T) {
	testCases := map[string]string{
		"no-id":                    "/",
		"bad-id":                   "/?id=test",
		"no-case-id":               "/?id=123",
		"bad-case-id":              "/?id=123&caseId=test",
		"bad-trust-corporation-id": "/?id=123&caseId=123&trustCorporationId=test",
	}

	for name, query := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, query, nil)
			w := httptest.NewRecorder()

			err := CreateTrustCorporation(nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestGetEditTrustCorporationWhenLpaFails(t *testing.T) {
	client := &mockCreateTrustCorporationClient{}
	client.
		On("Lpa", mock.Anything, 2).
		Return(sirius.Lpa{}, errExample)

	r, _ := http.NewRequest(http.MethodGet, "/create-trust-corporation/?id=1&caseId=2&trustCorporationId=3", nil)
	w := httptest.NewRecorder()

	err := CreateTrustCorporation(client, nil)(w, r)
	assert.Equal(t, errExample, err)
}

func TestPostCreateTrustCorporationTransitions(t *testing.T) {
	transitions := []struct {
		name                  string
		isReplacementAttorney bool
		addAnother            bool
		redirect              string
		fragment              string
		attorneyType          string
		systemStatus          *bool
	}{
		{
			name:         "Save attorney and return to LPA",
			redirect:     "/create-lpa?id=1&caseId=2",
			fragment:     "#scroll-to-attorneys-corporation",
			attorneyType: "Attorney",
			systemStatus: shared.BoolPtr(true),
		},
		{
			name:                  "Save replacement attorney and return to LPA",
			isReplacementAttorney: true,
			redirect:              "/create-lpa?id=1&caseId=2",
			fragment:              "#scroll-to-attorneys-corporation",
			attorneyType:          "Replacement Attorney",
			systemStatus:          shared.BoolPtr(false),
		},
		{
			name:         "Save attorney and add another",
			addAnother:   true,
			redirect:     "/create-attorney?id=1&caseId=2&caseType=lpa",
			attorneyType: "Attorney",
			systemStatus: shared.BoolPtr(true),
		},
		{
			name:                  "Save replacement attorney and add another",
			isReplacementAttorney: true,
			addAnother:            true,
			redirect:              "/create-replacement-attorney?id=1&caseId=2",
			attorneyType:          "Replacement Attorney",
			systemStatus:          shared.BoolPtr(false),
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
					expectedTrustCorporation := sirius.TrustCorporation{
						Attorney: sirius.Attorney{
							Person: sirius.Person{
								CompanyName:       "ACME",
								Email:             "test@test.com",
								PhoneNumber:       "1234",
								AddressLine1:      "221B",
								AddressLine2:      "Baker Street",
								AddressLine3:      "Marylebone",
								Town:              "London",
								County:            "Greater London",
								Postcode:          "NW1 6XE",
								Country:           "United Kingdom",
								IsAirmailRequired: false,
							},
							CompanyNumber: "123",
							SystemStatus:  transition.systemStatus,
						},
						IsReplacementAttorney:       transition.isReplacementAttorney,
						TrustCorporationAppointedAs: transition.attorneyType,
					}

					client := &mockCreateTrustCorporationClient{}
					client.
						On("CreateTrustCorporation", mock.Anything, 2, expectedTrustCorporation).
						Return(nil)

					form := url.Values{
						"companyName":              {"ACME"},
						"companyNumber":            {"123"},
						"email":                    {"test@test.com"},
						"phoneNumber":              {"1234"},
						"addressLine1":             {"221B"},
						"addressLine2":             {"Baker Street"},
						"addressLine3":             {"Marylebone"},
						"town":                     {"London"},
						"county":                   {"Greater London"},
						"postcode":                 {"NW1 6XE"},
						"country":                  {"United Kingdom"},
						"isAirmailRequired":        {"false"},
						"isReplacementAttorney":    {strconv.FormatBool(transition.isReplacementAttorney)},
						"isTrustCorporationActive": {"true"},
					}
					if transition.addAnother {
						form.Set("add-another", "true")
					}

					r, _ := http.NewRequest(http.MethodPost, "create-trust-corporation/?id=1&caseId=2&replacement="+strconv.FormatBool(transition.isReplacementAttorney)+flowCase.flowQuery, strings.NewReader(form.Encode()))
					r.Header.Add("Content-Type", formUrlEncoded)
					w := httptest.NewRecorder()

					err := CreateTrustCorporation(client, nil)(w, r)

					assert.Equal(t, RedirectError(transition.redirect+flowCase.flowQuery+transition.fragment), err)
					assert.Equal(t, http.StatusOK, w.Result().StatusCode)
					mock.AssertExpectationsForObjects(t, client)
				})
			}
		})
	}
}

func TestPostCreateTrustCorporationWhenCreateFails(t *testing.T) {
	expectedTrustCorporation := sirius.TrustCorporation{
		Attorney: sirius.Attorney{
			Person: sirius.Person{
				CompanyName: "ACME",
			},
			SystemStatus: shared.BoolPtr(true),
		},
		IsReplacementAttorney:       false,
		TrustCorporationAppointedAs: "Attorney",
	}

	client := &mockCreateTrustCorporationClient{}
	client.
		On("CreateTrustCorporation", mock.Anything, 2, expectedTrustCorporation).
		Return(errExample)

	form := url.Values{
		"companyName":              {"ACME"},
		"isReplacementAttorney":    {"false"},
		"isTrustCorporationActive": {"true"},
	}

	r, _ := http.NewRequest(http.MethodPost, "create-trust-corporation/?id=1&caseId=2&replacement=false", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateTrustCorporation(client, nil)(w, r)

	assert.Equal(t, errExample, err)
}

func TestPostEditTrustCorporationWhenUpdateFails(t *testing.T) {
	existingTrustCorporation := sirius.TrustCorporation{
		Attorney:              sirius.Attorney{Person: sirius.Person{ID: 3}},
		IsReplacementAttorney: false,
	}
	expectedTrustCorporation := sirius.TrustCorporation{
		Attorney: sirius.Attorney{
			Person: sirius.Person{
				CompanyName: "ACME",
			},
			SystemStatus: shared.BoolPtr(true),
		},
		IsReplacementAttorney:       false,
		TrustCorporationAppointedAs: "Attorney",
	}

	client := &mockCreateTrustCorporationClient{}
	client.
		On("Lpa", mock.Anything, 2).
		Return(sirius.Lpa{Case: sirius.Case{TrustCorporations: []sirius.TrustCorporation{existingTrustCorporation}}}, nil)
	client.
		On("UpdateTrustCorporation", mock.Anything, 3, expectedTrustCorporation).
		Return(errExample)

	form := url.Values{
		"companyName":              {"ACME"},
		"isReplacementAttorney":    {"false"},
		"isTrustCorporationActive": {"true"},
	}

	r, _ := http.NewRequest(http.MethodPost, "create-trust-corporation/?id=1&caseId=2&trustCorporationId=3&replacement=false", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateTrustCorporation(client, nil)(w, r)

	assert.Equal(t, errExample, err)
}

func TestPostCreateTrustCorporationWhenValidationError(t *testing.T) {
	expectedError := sirius.ValidationError{
		Field: sirius.FieldErrors{"field": {"": "problem"}},
	}

	client := &mockCreateTrustCorporationClient{}
	client.
		On("CreateTrustCorporation", mock.Anything, 2, mock.Anything).
		Return(expectedError)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createTrustCorporationData{
			AppointedAs: "Attorney",
			DonorId:     1,
			CaseId:      2,
			TrustCorporation: sirius.TrustCorporation{
				IsReplacementAttorney:       false,
				TrustCorporationAppointedAs: "Attorney",
				Attorney: sirius.Attorney{
					Person: sirius.Person{
						CompanyName: "ACME",
					},
					SystemStatus: shared.BoolPtr(true),
				},
			},
			Title:    "Add a trust corporation",
			HtmxPost: "/create-trust-corporation?id=1&caseId=2&replacement=false",
			Error:    expectedError,
		}).
		Return(nil)

	form := url.Values{
		"companyName":              {"ACME"},
		"isReplacementAttorney":    {"false"},
		"isTrustCorporationActive": {"true"},
	}

	r, _ := http.NewRequest(http.MethodPost, "create-trust-corporation/?id=1&caseId=2&replacement=false", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateTrustCorporation(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostEditTrustCorporationTransitions(t *testing.T) {
	transitions := []struct {
		name                  string
		lpa                   sirius.Lpa
		isReplacementAttorney bool
		formKey               string
		redirect              string
		fragment              string
	}{
		{
			name: "Save and return to LPA",
			lpa: sirius.Lpa{Case: sirius.Case{TrustCorporations: []sirius.TrustCorporation{
				{Attorney: sirius.Attorney{Person: sirius.Person{ID: 3, PersonType: "Trust Corporation"}}},
			}}},
			redirect: "/create-lpa?id=1&caseId=2",
			fragment: "#scroll-to-attorneys-corporation",
		},
		{
			name:    "Update next trust corporation",
			formKey: "update-next-attorney",
			lpa: sirius.Lpa{Case: sirius.Case{TrustCorporations: []sirius.TrustCorporation{
				{Attorney: sirius.Attorney{Person: sirius.Person{ID: 3, PersonType: "Trust Corporation"}}},
				{Attorney: sirius.Attorney{Person: sirius.Person{ID: 4, PersonType: "Trust Corporation"}}},
			}}},
			redirect: "/create-trust-corporation?id=1&caseId=2&trustCorporationId=4&replacement=false",
		},
		{
			name:    "Update next attorney",
			formKey: "update-next-attorney",
			lpa: sirius.Lpa{Case: sirius.Case{
				TrustCorporations: []sirius.TrustCorporation{{Attorney: sirius.Attorney{Person: sirius.Person{ID: 3, PersonType: "Trust Corporation"}}}},
				Attorneys:         []sirius.Attorney{{SystemStatus: shared.BoolPtr(true), Person: sirius.Person{ID: 4, PersonType: "Attorney"}}},
			}},
			redirect: "/create-attorney?id=1&caseId=2&caseType=lpa&attorneyId=4",
		},
		{
			name:                  "Update next replacement attorney",
			formKey:               "update-next-attorney",
			isReplacementAttorney: true,
			lpa: sirius.Lpa{Case: sirius.Case{
				TrustCorporations:    []sirius.TrustCorporation{{Attorney: sirius.Attorney{Person: sirius.Person{ID: 3, PersonType: "Trust Corporation"}}, IsReplacementAttorney: true}},
				ReplacementAttorneys: []sirius.Attorney{{SystemStatus: shared.BoolPtr(false), Person: sirius.Person{ID: 4, PersonType: "Replacement Attorney"}}},
			}},
			redirect: "/create-replacement-attorney?id=1&caseId=2&attorneyId=4",
		},
		{
			name:     "Update next with no following actor",
			formKey:  "update-next-attorney",
			lpa:      sirius.Lpa{Case: sirius.Case{TrustCorporations: []sirius.TrustCorporation{{Attorney: sirius.Attorney{Person: sirius.Person{ID: 3, PersonType: "Trust Corporation"}}}}}},
			redirect: "/create-lpa?id=1&caseId=2",
			fragment: "#scroll-to-attorneys-corporation",
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
					isReplacementAttorney := transition.isReplacementAttorney
					attorneyType := "Attorney"
					systemStatus := shared.BoolPtr(true)
					if isReplacementAttorney {
						attorneyType = "Replacement Attorney"
						systemStatus = shared.BoolPtr(false)
					}
					expectedTrustCorporation := sirius.TrustCorporation{
						Attorney: sirius.Attorney{
							Person: sirius.Person{
								CompanyName: "ACME",
							},
							SystemStatus: systemStatus,
						},
						IsReplacementAttorney:       isReplacementAttorney,
						TrustCorporationAppointedAs: attorneyType,
					}

					client := &mockCreateTrustCorporationClient{}
					client.
						On("Lpa", mock.Anything, 2).
						Return(transition.lpa, nil)
					client.
						On("UpdateTrustCorporation", mock.Anything, 3, expectedTrustCorporation).
						Return(nil)

					form := url.Values{
						"companyName":              {"ACME"},
						"isReplacementAttorney":    {strconv.FormatBool(isReplacementAttorney)},
						"isTrustCorporationActive": {"true"},
					}
					if transition.formKey != "" {
						form.Set(transition.formKey, "true")
					}

					r, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("create-trust-corporation/?id=1&caseId=2&trustCorporationId=3&replacement=%t%s", isReplacementAttorney, flowCase.flowQuery), strings.NewReader(form.Encode()))
					r.Header.Add("Content-Type", formUrlEncoded)
					w := httptest.NewRecorder()

					err := CreateTrustCorporation(client, nil)(w, r)

					assert.Equal(t, RedirectError(transition.redirect+flowCase.flowQuery+transition.fragment), err)
					assert.Equal(t, http.StatusOK, w.Result().StatusCode)
					mock.AssertExpectationsForObjects(t, client)
				})
			}
		})
	}
}

func TestGetIdForNextAttorneyTrustCorpsAndActiveAttorneys(t *testing.T) {
	trustCorporations := []sirius.TrustCorporation{
		{
			Attorney: sirius.Attorney{Person: sirius.Person{ID: 1, PersonType: "Trust Corporation"}},
		},
		{
			Attorney: sirius.Attorney{Person: sirius.Person{ID: 3, PersonType: "Trust Corporation"}},
		},
		{
			Attorney:              sirius.Attorney{Person: sirius.Person{ID: 4, PersonType: "Trust Corporation"}},
			IsReplacementAttorney: true,
		},
		{
			Attorney: sirius.Attorney{Person: sirius.Person{ID: 5, PersonType: "Trust Corporation"}},
		},
	}

	attorneys := []sirius.Attorney{
		{
			Person:       sirius.Person{ID: 2, PersonType: "Attorney"},
			SystemStatus: shared.BoolPtr(true),
		},
		{
			Person:       sirius.Person{ID: 6, PersonType: "Attorney"},
			SystemStatus: shared.BoolPtr(true),
		},
		{
			Person:       sirius.Person{ID: 7, PersonType: "Attorney"},
			SystemStatus: shared.BoolPtr(true),
		},
	}

	tests := []struct {
		name           string
		currentID      int
		expectedNextID int
		personType     string
	}{
		{
			name:           "Current is attorney, next is trust corp",
			currentID:      2,
			expectedNextID: 3,
			personType:     "Trust Corporation",
		},
		{
			name:           "Current is attorney, next is attorney",
			currentID:      6,
			expectedNextID: 7,
			personType:     "Attorney",
		},
		{
			name:           "Current is trust corp, next is attorney",
			currentID:      1,
			expectedNextID: 2,
			personType:     "Attorney",
		},
		{
			name:           "Current is trust corp, next is trust corp",
			currentID:      3,
			expectedNextID: 5,
			personType:     "Trust Corporation",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, personType := GetIdForNextAttorney(tc.currentID, false, trustCorporations, attorneys)
			assert.Equal(t, tc.expectedNextID, result)
			assert.Equal(t, tc.personType, personType)
		})
	}
}
