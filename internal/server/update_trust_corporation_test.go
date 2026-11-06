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

type mockUpdateTrustCorporationClient struct {
	mock.Mock
}

func (m *mockUpdateTrustCorporationClient) Lpa(ctx sirius.Context, id int) (sirius.Lpa, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Lpa), args.Error(1)
}

func (m *mockUpdateTrustCorporationClient) UpdateTrustCorporation(ctx sirius.Context, trustCorporationId int, trustCorporation sirius.TrustCorporation) error {
	args := m.Called(ctx, trustCorporationId, trustCorporation)
	return args.Error(0)
}

func TestGetUpdateTrustCorporationBadQuery(t *testing.T) {
	testCases := map[string]string{
		"bad-trust-corporation-id": "/?id=123&caseId=123&trustCorporationId=test",
	}

	for name, query := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, query, nil)
			w := httptest.NewRecorder()

			err := UpdateTrustCorporation(nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestGetUpdateTrustCorporationWhenLpaFails(t *testing.T) {
	client := &mockUpdateTrustCorporationClient{}
	client.
		On("Lpa", mock.Anything, 2).
		Return(sirius.Lpa{}, errExample)

	r, _ := http.NewRequest(http.MethodGet, "/update-trust-corporation/?id=1&caseId=2&trustCorporationId=3", nil)
	w := httptest.NewRecorder()

	err := UpdateTrustCorporation(client, nil)(w, r)
	assert.Equal(t, errExample, err)
}

func TestPostUpdateTrustCorporationWhenUpdateFails(t *testing.T) {
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

	client := &mockUpdateTrustCorporationClient{}
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

	r, _ := http.NewRequest(http.MethodPost, "update-trust-corporation/?id=1&caseId=2&trustCorporationId=3&replacement=false", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := UpdateTrustCorporation(client, nil)(w, r)

	assert.Equal(t, errExample, err)
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
			redirect: "/update-trust-corporation?id=1&caseId=2&trustCorporationId=4&replacement=false",
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

					client := &mockUpdateTrustCorporationClient{}
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

					r, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("update-trust-corporation/?id=1&caseId=2&trustCorporationId=3&replacement=%t%s", isReplacementAttorney, flowCase.flowQuery), strings.NewReader(form.Encode()))
					r.Header.Add("Content-Type", formUrlEncoded)
					w := httptest.NewRecorder()

					err := UpdateTrustCorporation(client, nil)(w, r)

					assert.Equal(t, RedirectError(transition.redirect+flowCase.flowQuery+transition.fragment), err)
					assert.Equal(t, http.StatusOK, w.Result().StatusCode)
					mock.AssertExpectationsForObjects(t, client)
				})
			}
		})
	}
}
