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

func TestCreateCertificateProvider(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name                string
		certificateProvider Person
		setup               func()
		expectedError       func(int) error
	}{
		{
			name: "OK",
			certificateProvider: Person{
				Salutation:   "Prof",
				Firstname:    "Melanie",
				Middlenames:  "Josefina",
				Surname:      "Vanvolkenburg",
				AddressLine1: "29737 Andrew Plaza",
				AddressLine2: "Apt. 814",
				AddressLine3: "Gislasonside",
				Town:         "Hirthehaven",
				County:       "Saskatchewan",
				Postcode:     "S7R 9F9",
				Country:      "United Kingdom",
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending case assigned").
					UponReceiving("A request to create a certificate provider").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPost,
						Path:   matchers.String("/lpa-api/v1/persons"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: []map[string]interface{}{{
							"caseId":                800,
							"personType":            "CertificateProvider",
							"salutation":            "Prof",
							"firstname":             "Melanie",
							"middlenames":           "Josefina",
							"surname":               "Vanvolkenburg",
							"addressLine1":          "29737 Andrew Plaza",
							"addressLine2":          "Apt. 814",
							"addressLine3":          "Gislasonside",
							"town":                  "Hirthehaven",
							"county":                "Saskatchewan",
							"postcode":              "S7R 9F9",
							"country":               "United Kingdom",
							"companyName":           "",
							"companyReference":      "",
							"correspondenceByEmail": false,
							"correspondenceByPhone": false,
							"correspondenceByPost":  false,
							"correspondenceByWelsh": false,
							"dateOfDeath":           nil,
							"dob":                   nil,
							"email":                 "",
							"isAirmailRequired":     false,
							"otherNames":            "",
							"phoneNumber":           "",
							"previousNames":         "",
							"researchOptOut":        false,
							"sageId":                "",
						}},
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusCreated,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				err := client.CreateCertificateProvider(Context{Context: context.Background()}, 800, tc.certificateProvider)
				if (tc.expectedError) == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(config.Port), err)
				}
				return nil
			}))
		})
	}
}

func TestUpdateCertificateProvider(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name                string
		certificateProvider Person
		setup               func()
		expectedError       func(int) error
	}{
		{
			name: "OK",
			certificateProvider: Person{
				ID:           123,
				Salutation:   "Prof",
				Firstname:    "Melanie",
				Middlenames:  "Josefina",
				Surname:      "Vanvolkenburg",
				AddressLine1: "29737 Andrew Plaza",
				AddressLine2: "Apt. 814",
				AddressLine3: "Gislasonside",
				Town:         "Hirthehaven",
				County:       "Saskatchewan",
				Postcode:     "S7R 9F9",
				Country:      "Canada",
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("I have an lpa with a certificate provider").
					UponReceiving("A request to update that certificate provider").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPut,
						Path:   matchers.String("/lpa-api/v1/certificate-providers/123"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]interface{}{
							"id":                    123,
							"salutation":            "Prof",
							"firstname":             "Melanie",
							"middlenames":           "Josefina",
							"surname":               "Vanvolkenburg",
							"addressLine1":          "29737 Andrew Plaza",
							"addressLine2":          "Apt. 814",
							"addressLine3":          "Gislasonside",
							"otherNames":            "",
							"companyName":           "",
							"companyReference":      "",
							"correspondenceByEmail": false,
							"correspondenceByPhone": false,
							"correspondenceByPost":  false,
							"correspondenceByWelsh": false,
							"town":                  "Hirthehaven",
							"county":                "Saskatchewan",
							"postcode":              "S7R 9F9",
							"country":               "Canada",
							"dob":                   nil,
							"dateOfDeath":           nil,
							"email":                 "",
							"isAirmailRequired":     false,
							"phoneNumber":           "",
							"previousNames":         "",
							"researchOptOut":        false,
							"sageId":                "",
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				err := client.UpdateCertificateProvider(Context{Context: context.Background()}, 123, tc.certificateProvider)
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
