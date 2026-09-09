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

func TestCreateTrustCorporation(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		trustCorporation TrustCorporation
		setup            func()
		expectedError    func(int) error
	}{
		{
			name: "OK",
			trustCorporation: TrustCorporation{
				Attorney: Attorney{
					Person: Person{
						CompanyName:       "Trust Ltd.",
						AddressLine1:      "29737 Andrew Plaza",
						AddressLine2:      "Apt. 814",
						AddressLine3:      "Gislasonside",
						Town:              "Hirthehaven",
						County:            "Saskatchewan",
						Postcode:          "S7R 9F9",
						Country:           "Canada",
						Email:             "test@test.com",
						PhoneNumber:       "072345678",
						IsAirmailRequired: true,
					},
					CompanyNumber: "123",
				},
				TrustCorporationAppointedAs: "Attorney",
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending case assigned").
					UponReceiving("A request to create a trust corporation").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPost,
						Path:   matchers.String("/lpa-api/v1/trust-corporation"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]interface{}{
							"companyName":                 "Trust Ltd.",
							"addressLine1":                "29737 Andrew Plaza",
							"addressLine2":                "Apt. 814",
							"addressLine3":                "Gislasonside",
							"town":                        "Hirthehaven",
							"county":                      "Saskatchewan",
							"postcode":                    "S7R 9F9",
							"country":                     "Canada",
							"email":                       "test@test.com",
							"phoneNumber":                 "072345678",
							"isAirmailRequired":           true,
							"companyNumber":               "123",
							"trustCorporationAppointedAs": "Attorney",
							"companyReference":            "",
							"correspondenceByEmail":       false,
							"correspondenceByPhone":       false,
							"correspondenceByPost":        false,
							"correspondenceByWelsh":       false,
							"dateOfDeath":                 nil,
							"dob":                         nil,
							"firstname":                   "",
							"middlenames":                 "",
							"otherNames":                  "",
							"surname":                     "",
							"previousNames":               "",
							"researchOptOut":              false,
							"sageId":                      "",
							"salutation":                  "",
							"caseId":                      800,
						},
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

				err := client.CreateTrustCorporation(Context{Context: context.Background()}, 800, tc.trustCorporation)
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

func TestUpdateTrustCorporation(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		trustCorporation TrustCorporation
		setup            func()
		expectedError    func(int) error
	}{
		{
			name: "OK",
			trustCorporation: TrustCorporation{
				Attorney: Attorney{
					Person: Person{
						CompanyName:       "Trust Ltd.",
						AddressLine1:      "29737 Andrew Plaza",
						AddressLine2:      "Apt. 814",
						AddressLine3:      "Gislasonside",
						Town:              "Hirthehaven",
						County:            "Saskatchewan",
						Postcode:          "S7R 9F9",
						Country:           "Canada",
						Email:             "test@test.com",
						PhoneNumber:       "072345678",
						IsAirmailRequired: true,
					},
					CompanyNumber: "123",
				},
				TrustCorporationAppointedAs: "Attorney",
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("I have an lpa with a trust corporation").
					UponReceiving("A request to update that trust corporation").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPut,
						Path:   matchers.String("/lpa-api/v1/trust-corporation/123"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]interface{}{
							"companyName":                 "Trust Ltd.",
							"addressLine1":                "29737 Andrew Plaza",
							"addressLine2":                "Apt. 814",
							"addressLine3":                "Gislasonside",
							"town":                        "Hirthehaven",
							"county":                      "Saskatchewan",
							"postcode":                    "S7R 9F9",
							"country":                     "Canada",
							"email":                       "test@test.com",
							"phoneNumber":                 "072345678",
							"isAirmailRequired":           true,
							"companyNumber":               "123",
							"trustCorporationAppointedAs": "Attorney",
							"companyReference":            "",
							"correspondenceByEmail":       false,
							"correspondenceByPhone":       false,
							"correspondenceByPost":        false,
							"correspondenceByWelsh":       false,
							"dateOfDeath":                 nil,
							"dob":                         nil,
							"firstname":                   "",
							"middlenames":                 "",
							"otherNames":                  "",
							"surname":                     "",
							"previousNames":               "",
							"researchOptOut":              false,
							"sageId":                      "",
							"salutation":                  "",
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

				err := client.UpdateTrustCorporation(Context{Context: context.Background()}, 123, tc.trustCorporation)
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
