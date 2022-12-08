package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestCase(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

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
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/lpa-api/v1/cases/800"),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: dsl.Like(map[string]interface{}{
							"id":       800,
							"uId":      dsl.String("7000-0000-0000"),
							"caseType": dsl.String("LPA"),
							"status":   dsl.String("Pending"),
							"donor": dsl.Like(map[string]interface{}{
								"id":           dsl.Like(771),
								"uId":          dsl.String("7000-0290-0192"),
								"salutation":   dsl.String("Prof"),
								"firstname":    dsl.String("Melanie"),
								"surname":      dsl.String("Vanvolkenburg"),
								"addressLine1": dsl.String("29737 Andrew Plaza"),
								"addressLine2": dsl.String("Apt. 814"),
								"addressLine3": dsl.String("Gislasonside"),
								"town":         dsl.String("Hirthehaven"),
								"county":       dsl.String("Saskatchewan"),
								"postcode":     dsl.String("S7R 9F9"),
								"personType":   dsl.String("Donor"),
							}),
						}),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
			},
			expectedResponse: Case{ID: 800, UID: "7000-0000-0000", CaseType: "LPA", Status: "Pending", Donor: &Person{
				ID:           771,
				UID:          "7000-0290-0192",
				Salutation:   "Prof",
				Firstname:    "Melanie",
				Surname:      "Vanvolkenburg",
				AddressLine1: "29737 Andrew Plaza",
				AddressLine2: "Apt. 814",
				AddressLine3: "Gislasonside",
				Town:         "Hirthehaven",
				County:       "Saskatchewan",
				Postcode:     "S7R 9F9",
				PersonType:   "Donor",
			}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				caseitem, err := client.Case(Context{Context: context.Background()}, 800)

				assert.Equal(t, tc.expectedResponse, caseitem)
				if tc.expectedError == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(pact.Server.Port), err)
				}
				return nil
			}))
		})
	}
}

func TestCaseNoPayments(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

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
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/lpa-api/v1/cases/801"),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: dsl.Like(map[string]interface{}{
							"uId":      dsl.String("7000-0000-0001"),
							"caseType": dsl.String("LPA"),
							"status":   dsl.String("Pending"),
						}),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
			},
			expectedResponse: Case{UID: "7000-0000-0001", CaseType: "LPA", Status: "Pending"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				caseitem, err := client.Case(Context{Context: context.Background()}, 801)

				assert.Equal(t, tc.expectedResponse, caseitem)
				if tc.expectedError == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(pact.Server.Port), err)
				}
				return nil
			}))
		})
	}
}

func TestCaseWithFeeReduction(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

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
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/lpa-api/v1/cases/802"),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: dsl.Like(map[string]interface{}{
							"uId":      dsl.String("7000-0000-0002"),
							"caseType": dsl.String("LPA"),
							"status":   dsl.String("Pending"),
						}),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
			},
			expectedResponse: Case{UID: "7000-0000-0002", CaseType: "LPA", Status: "Pending"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				caseitem, err := client.Case(Context{Context: context.Background()}, 802)

				assert.Equal(t, tc.expectedResponse, caseitem)
				if tc.expectedError == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(pact.Server.Port), err)
				}
				return nil
			}))
		})
	}
}
