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

func TestCreateReplacementAttorney(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		setup            func()
		attorney         Attorney
		expectedResponse Attorney
	}{
		{
			name: "OK",
			attorney: Attorney{
				Person: Person{
					Salutation:   "Mrs",
					Firstname:    "Rosalind",
					Surname:      "Achebe",
					Email:        "r.achebe@example.test",
					AddressLine1: "14 Kestrel Way",
					Town:         "Hirthehaven",
					County:       "Saskatchewan",
					Postcode:     "S7R 9F9",
					Country:      "Canada",
				},
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending case assigned").
					UponReceiving("A request to create a replacement attorney").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPost,
						Path:   matchers.String("/lpa-api/v1/persons"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: []interface{}{
							map[string]interface{}{
								"personType":            "ReplacementAttorney",
								"caseId":                800,
								"salutation":            "Mrs",
								"firstname":             "Rosalind",
								"middlenames":           "",
								"surname":               "Achebe",
								"dob":                   nil,
								"dateOfDeath":           nil,
								"previousNames":         "",
								"otherNames":            "",
								"addressLine1":          "14 Kestrel Way",
								"addressLine2":          "",
								"addressLine3":          "",
								"town":                  "Hirthehaven",
								"county":                "Saskatchewan",
								"postcode":              "S7R 9F9",
								"country":               "Canada",
								"sageId":                "",
								"isAirmailRequired":     false,
								"phoneNumber":           "",
								"email":                 "r.achebe@example.test",
								"correspondenceByPost":  false,
								"correspondenceByEmail": false,
								"correspondenceByPhone": false,
								"correspondenceByWelsh": false,
								"researchOptOut":        false,
								"companyName":           "",
								"companyNumber":         "",
								"companyReference":      "",
							},
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusCreated,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: []interface{}{
							map[string]interface{}{
								"id":         matchers.Like(321),
								"firstname":  matchers.String("Rosalind"),
								"surname":    matchers.String("Achebe"),
								"personType": matchers.String("Replacement Attorney"),
							},
						},
					})
			},
			expectedResponse: Attorney{Person: Person{ID: 321, Firstname: "Rosalind", Surname: "Achebe", PersonType: "Replacement Attorney"}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				attorney, err := client.CreateReplacementAttorney(Context{Context: context.Background()}, 800, tc.attorney)
				assert.Equal(t, tc.expectedResponse, attorney)
				assert.Nil(t, err)
				return nil
			}))
		})
	}
}
