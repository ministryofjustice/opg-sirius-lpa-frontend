package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

func TestCase(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		setup            func()
		expectedResponse Case
		expectedError    func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending case assigned").
					UponReceiving("A request for the case").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/cases/800"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.Like(map[string]interface{}{
							"id":       800,
							"uId":      matchers.String("7000-0000-0000"),
							"caseType": matchers.String("LPA"),
							"status":   matchers.String("Pending"),
							"donor": matchers.Like(map[string]interface{}{
								"id": matchers.Like(189),
							}),
						}),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: Case{ID: 800, UID: "7000-0000-0000", CaseType: "LPA", Status: "Pending", Donor: &Person{ID: 189}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				caseitem, err := client.Case(Context{Context: context.Background()}, 800)

				assert.Equal(t, tc.expectedResponse, caseitem)
				if tc.expectedError == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(config.Port), err)
				}
				return nil
			}))
		})
	}
}

func TestCaseNoPayments(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		setup            func()
		expectedResponse Case
		expectedError    func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending case with no payment assigned").
					UponReceiving("A request for the case").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/cases/801"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.Like(map[string]interface{}{
							"uId":      matchers.String("7000-0000-0001"),
							"caseType": matchers.String("LPA"),
							"status":   matchers.String("Pending"),
							"donor": matchers.Like(map[string]interface{}{
								"id": matchers.Like(189),
							}),
						}),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: Case{UID: "7000-0000-0001", CaseType: "LPA", Status: "Pending", Donor: &Person{ID: 189}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				caseitem, err := client.Case(Context{Context: context.Background()}, 801)

				assert.Equal(t, tc.expectedResponse, caseitem)
				if tc.expectedError == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(config.Port), err)
				}
				return nil
			}))
		})
	}
}

func TestCaseWithFeeReduction(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		setup            func()
		expectedResponse Case
		expectedError    func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending case with a fee reduction assigned").
					UponReceiving("A request for the case").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/cases/802"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.Like(map[string]interface{}{
							"uId":      matchers.String("7000-0000-0002"),
							"caseType": matchers.String("LPA"),
							"status":   matchers.String("Pending"),
							"donor": matchers.Like(map[string]interface{}{
								"id": matchers.Like(189),
							}),
						}),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: Case{UID: "7000-0000-0002", CaseType: "LPA", Status: "Pending", Donor: &Person{ID: 189}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				caseitem, err := client.Case(Context{Context: context.Background()}, 802)

				assert.Equal(t, tc.expectedResponse, caseitem)
				if tc.expectedError == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(config.Port), err)
				}
				return nil
			}))
		})
	}
}

func TestCaseFiltersInactiveActors(t *testing.T) {
	actor1 := Person{ID: 1, SystemStatus: true}
	actor2 := Person{ID: 2, SystemStatus: true}
	inactiveActor1 := Person{ID: 3, SystemStatus: false}
	inactiveActor2 := Person{ID: 4, SystemStatus: false}

	caseItem := Case{ID: 1, Attorneys: []Person{actor1, inactiveActor1}, TrustCorporations: []Person{actor2, inactiveActor2}}
	filteredCase := caseItem.FilterInactiveAttorneys()

	assert.Equal(t, []Person{actor1}, filteredCase.Attorneys)
	assert.Equal(t, []Person{actor2}, filteredCase.TrustCorporations)
}
